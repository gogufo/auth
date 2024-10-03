// Copyright 2020 - 2024 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Sign In
// SignIn function authorisate user in Gufo.
//

package post

import (
	. "auth/grpc_requests"
	. "auth/model"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
)

// POST only
// func Signin(t *Request, r *http.Request) (map[string]interface{}, []ErrorMsg, *Request) {
func Signin(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	p := bluemonday.UGCPolicy()
	args := ToMapStringInterface(t.Args)

	ottoken := ""
	refreshtoken := ""

	if args["ot_token"] != nil {
		ottoken = p.Sanitize(fmt.Sprintf("%s", args["ot_token"]))
	}

	if ottoken != "" {
		return signinwithtoken(t)
	}

	if args["refresh_token"] != nil {
		refreshtoken = p.Sanitize(fmt.Sprintf("%s", args["refresh_token"]))
	}

	if refreshtoken != "" {
		return RefreshToken(t)
	}

	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	//1. Check for need  data
	creds := &SignInCred{}
	var userExist Users

	ptfa := ""
	if args["tfa"] != nil {
		ptfa = p.Sanitize(fmt.Sprintf("%s", args["tfa"]))
	}

	if ptfa != "" {
		// Check 2FA
		tfa := ptfa
		uname := p.Sanitize(fmt.Sprintf("%s", args["user"]))

		lifetime, _, errstr := CheckTimeHash(t, tfa, uname)

		if errstr != "" {
			// return error. user name is exist in db users
			return ErrorReturn(t, 400, "000021", "There is no data")
		}

		// Check for OTP livetime
		ctime := int(time.Now().Unix())
		SetErrorLog(fmt.Sprintf("ctime: %v", ctime))
		SetErrorLog(fmt.Sprintf("lifetime: %v", lifetime))

		if ctime > lifetime {
			//Delete OTP
			go DeleteTimeHash(t, tfa, uname)
			return ErrorReturn(t, 400, "000022", "OTP has expired")
		}

		//If right - return token

		rows := db.Conn.Debug().Where(`(name = ? OR mail = ?)`, uname, uname).First(&userExist)

		if rows.RowsAffected == 0 {
			// return error. user name is exist in db users
			return ErrorReturn(t, 400, "000003", "There is no such user")
		}

		//2. If user active
		if !userExist.Status {

			return ErrorReturn(t, 400, "000013", "User blocked")
		}

		//4. Check if user confirmed his email
		ans["email_confirmed"] = true
		if !userExist.Completed {
			ans["email_confirmed"] = false //User blocked

		}

		//5. Create token
		isadm := 0
		iscomp := 0
		readon := 0
		if userExist.IsAdmin {
			isadm = 1
		}
		if userExist.Completed {
			iscomp = 1
		}
		if userExist.Readonly {
			readon = 1
		}

		access_token, at_lifetime, refresh_token, rt_lifetime := UpdateSession(t, userExist.UID, isadm, iscomp, readon)

		//	*t.IsAdmin = int32(0)
		if userExist.IsAdmin {
			adm := int32(1)
			t.IsAdmin = &adm
		}

		//6. Write data to users table

		db.Conn.Table("users").Where("uid = ?", userExist.UID).Updates(map[string]interface{}{"access": int(time.Now().Unix()), "login": int(time.Now().Unix())})

		//8. TODO Write data to signin history table

		//return data
		userid := userExist.UID
		sesionend := int32(at_lifetime)
		t.UID = &userid
		t.SessionEnd = &sesionend

		ans["access_token"] = access_token
		ans["refresh_token"] = refresh_token
		ans["rt_lifetime"] = rt_lifetime
		ans["at_lifetime"] = at_lifetime
		//ans["uid"] = userExist.UID
		ans["username"] = userExist.Name
		ans["email"] = userExist.Mail
		//ans["session_expired"] = expecttime

		response = Interfacetoresponse(t, ans)
		return response

	}

	creds.Username = p.Sanitize(fmt.Sprintf("%s", args["user"]))
	creds.Password = p.Sanitize(fmt.Sprintf("%s", args["pass"]))

	if creds.Username == "" || creds.Password == "" {

		return ErrorReturn(t, 400, "000001", "Missing Name or Password")
	}

	//2. Check if user exist

	rows := db.Conn.Debug().Where(`(name = ? OR mail = ?)`, creds.Username, creds.Username).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "There is no such user")
	}

	//2. If user active
	if !userExist.Status {

		return ErrorReturn(t, 400, "000013", "User blocked")
	}

	//3. Check password

	if err := bcrypt.CompareHashAndPassword([]byte(userExist.Pass), []byte(creds.Password)); err != nil {
		// Password not matched

		return ErrorReturn(t, 400, "000008", "Password not matched")

	}

	//3.1 Check for 2FA
	if userExist.TFA {

		//1. Generate OTP and send email. Retrun user information about 2FA required

		otp := Numgen(6)
		lang := "eng"
		if t.Language != nil {
			lang = *t.Language
		}
		go SendOTP(t, userExist.Mail, lang, otp)
		go SendTimeHash(t, otp, userExist.Name, "tfa", userExist.Mail, 300)

		askedemail := maskemail(userExist.Mail)

		ans["tfa"] = true
		ans["tfatype"] = userExist.TFAType
		ans["sendto"] = askedemail
		response = Interfacetoresponse(t, ans)
		return response
	}

	//4. Check if user confirmed his email
	ans["email_confirmed"] = true
	if !userExist.Completed {
		ans["email_confirmed"] = false //User blocked

	}

	isadm := 0
	iscomp := 0
	readon := 0
	if userExist.IsAdmin {
		isadm = 1
	}
	if userExist.Completed {
		iscomp = 1
	}
	if userExist.Readonly {
		readon = 1
	}

	//5. Create token
	access_token, at_lifetime, refresh_token, rt_lifetime := UpdateSession(t, userExist.UID, isadm, iscomp, readon)

	//	*t.IsAdmin = int32(0)
	if userExist.IsAdmin {
		var isa int32
		isa = 1
		t.IsAdmin = &isa
	}

	//8. TODO Write data to signin history table
	AuthHist := AuthHistory{}
	AuthHist.UID = userExist.UID
	AuthHist.Login = int(time.Now().Unix())
	AuthHist.IP = *t.IP
	AuthHist.UserAgent = *t.UserAgent

	db.Conn.Create(&AuthHist)

	//6. Write data to users table
	db.Conn.Table("users").Where("uid = ?", userExist.UID).Updates(map[string]interface{}{"ip": *t.IP, "access": int(time.Now().Unix()), "login": int(time.Now().Unix())})

	//9. Compare current IP address and previous IP addres if it different, send notification email
	lastip := userExist.IP
	curip := *t.IP
	if lastip != curip {
		//Generate token for block account

		//Send notification, that Somebody was login from another device
		template := "unrecognized_sign"
		title := "Unrecognized device signed in to your Stripe account"
		allmessage := []string{}
		message1 := "We don't recognize the device that was just used to sign in to your Amy account. If this was you, you don't need to do anything. If you don't recognize it, please let us know."
		allmessage = append(allmessage, message1)
		allmessage = append(allmessage, *t.UserAgent)
		allmessage = append(allmessage, *t.IP)
		blocktoken := Stringen(64)
		allmessage = append(allmessage, blocktoken)

		//Send block token to timehash Table
		SendTimeHash(t, blocktoken, userExist.UID, "block", userExist.Mail, 86400)

		//Send Notification
		go SendNotification(t, title, allmessage, template, userExist.UID)
	}

	//return data
	t.UID = &userExist.UID
	var sessionend int32
	sessionend = int32(at_lifetime)
	t.SessionEnd = &sessionend

	ans["access_token"] = access_token
	ans["refresh_token"] = refresh_token
	ans["at_lifetime"] = at_lifetime
	ans["rt_lifetime"] = rt_lifetime
	//ans["uid"] = userExist.UID
	ans["username"] = userExist.Name
	ans["email"] = userExist.Mail
	//ans["session_expired"] = expecttime

	response = Interfacetoresponse(t, ans)
	return response

}

func maskemail(email string) string {
	mailarr := strings.Split(email, "@")
	domain := mailarr[1]
	milbody := mailarr[0]
	domainarr := strings.Split(domain, ".")
	domainbody := domainarr[0]
	mailbodymask := milbody[0:1] + "***"
	domainbodymask := domainbody[0:1] + "***" + domainbody[len(domainbody)-1:]
	maskedemail := mailbodymask + "@" + domainbodymask + "." + domainarr[1]
	return maskedemail
}
