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

package post

import (
	. "auth/grpc_requests"
	. "auth/model"
	"fmt"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

// POST only
func RefreshToken(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	p := bluemonday.UGCPolicy()

	refresh_token := ""

	if t.Args["refresh_token"] != nil {
		refresh_token = p.Sanitize(fmt.Sprintf("%s", t.Args["refresh_token"]))
	}

	uid := RefreshSession(t, refresh_token)

	var userExist Users
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}
	rows := db.Conn.Debug().Where(`uid = ?`, uid).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "There is no such user")
	}

	//2. If user active
	if !userExist.Status {

		return ErrorReturn(t, 400, "000013", "User blocked")
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

	access_token, at_lifetime, refresh_token, rt_lifetime := UpdateSession(t, uid, isadm, iscomp, readon)

	//4. Check if user confirmed his email
	ans["email_confirmed"] = true
	if !userExist.Completed {
		ans["email_confirmed"] = false //User blocked

	}

	if userExist.IsAdmin {
		var isa int32
		isa = 1
		t.IsAdmin = &isa
	}

	t.UID = &uid
	var sessionend int32
	sessionend = int32(at_lifetime)
	t.SessionEnd = &sessionend

	ans["access_token"] = access_token
	ans["refresh_token"] = refresh_token
	ans["at_lifetime"] = at_lifetime
	ans["rt_lifetime"] = rt_lifetime

	response = Interfacetoresponse(t, ans)
	return response

}
