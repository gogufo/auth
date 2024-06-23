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
	"time"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func Forgot(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	//1. Check for need  data

	p := bluemonday.UGCPolicy()
	email := p.Sanitize(fmt.Sprintf("%s", args["email"]))
	lang := "eng"
	if t.Language != nil {
		lang = p.Sanitize(*t.Language)
	}

	if email == "" {
		return ErrorReturn(t, 400, "000001", "Missing email")
	}

	//2. Check if user exist
	var userExist Users

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

	rows := db.Conn.Where(`mail = ?`, email).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "User is not exist")
	}

	hashedkey := ""
	if args["key"] != nil {
		hashedkey = p.Sanitize(fmt.Sprintf("%s", args["key"]))
	}

	if hashedkey == "" {
		//If no confirmation code - just send this code to email
		hashkey := Numgen(6)

		//Write key to key Table
		go SendTimeHash(t, hashkey, userExist.UID, "forgot", email, 172800)

		//send email
		go SendForgot(t, userExist.Mail, lang, hashkey)

		ans["response"] = "100201" // sent email with confirmation code
		ans["email"] = userExist.Mail

		response = Interfacetoresponse(t, ans)
		return response

	}

	//check for key

	lifetime, _, errstr := CheckTimeHash(t, hashedkey, email)
	if errstr != "" {
		// return error. Hash is not exist in db
		return ErrorReturn(t, 400, "000008", "Hash is not exist in db")
	}

	//Check lifetime
	if int(time.Now().Unix()) > lifetime {
		go DeleteTimeHash(t, hashedkey, email)

		return ErrorReturn(t, 400, "000009", "Hash expired")

	}

	go DeleteTimeHash(t, hashedkey, email)
	// Create a new password
	userpass := RandomString(12)

	//2.1 generete pass passhash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userpass), 8)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog("forgot.go: " + err.Error())
		}
	}

	//6. Write data to users table
	db.Conn.Table("users").Where("mail = ?", email).Updates(map[string]interface{}{"pass": hashedPassword})

	go SendForgot(t, userExist.Mail, lang, userpass)
	//return data

	ans["response"] = "100202" // Password changed and sent email

	response = Interfacetoresponse(t, ans)
	return response
}
