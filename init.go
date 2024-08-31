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

package main

import (
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	ad "auth/admin"
	gt "auth/get"
	pd "auth/post"
	. "auth/version"
)

func Init(t *pb.Request) (response *pb.Response) {

	switch *t.Param {
	case "admin":
		return admincheck(t)
	}

	if *t.Method == "GET" {

		switch *t.Param {
		case "info":
			response = info(t)
		case "health":
			response = health(t)
		default:
			response = gt.Init(t)
		}
	}
	if *t.Method == "POST" {
		response = pd.Init(t)
	}

	return response

}

func info(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["microservicename"] = "Authorisation"
	ans["version"] = VERSIONPLUGIN
	ans["description"] = "SignIn Microservice"
	response = Interfacetoresponse(t, ans)
	return response
}

func health(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["health"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}

func admincheck(t *pb.Request) (response *pb.Response) {

	if *t.IsAdmin != 1 {
		response = ErrorReturn(t, 401, "000012", "You have no admin rights")
	}

	return ad.Init(t)

}
