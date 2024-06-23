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
	"time"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func confirmemail(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})

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

	var userExist Users

	rows := db.Conn.Where(`uid = ?`, t.UID).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "0000031", "There is no such user")
	}

	//Check if user already request for confirmatin email

	//check for hash lifetime
	ctime := int(time.Now().Unix())
	waittime := 300
	realtime := ctime - userExist.Mailsent
	//SetErrorLog("ctime: " + fmt.Sprintf("%d", ctime))
	//SetErrorLog("realtime: " + fmt.Sprintf("%d", userExist.Mailsent))
	if realtime < waittime {
		return ErrorReturn(t, 400, "000012", "You already asked for confirmation email")
	}

	db.Conn.Table("users").Where("uid = ?", t.UID).Updates(map[string]interface{}{"mailsent": int(time.Now().Unix()), "completed": 0})

	go SendConfiramtion(t, userExist.UID, userExist.Mail)

	ans["response"] = "100201"
	ans["message"] = "Confirmation message was sent to your email"

	response = Interfacetoresponse(t, ans)
	return response
}
