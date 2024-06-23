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

package grpc_requests

import (
	"fmt"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func UpdateSession(t *pb.Request, uid string, isadm int, iscomp int, readon int) (access_token string, at_lifetime int, refresh_token string, rt_lifetime int) {

	host := viper.GetString("server.internal_host")
	port := viper.GetString("server.grpc_port")

	s := &pb.Request{}
	module := "session"
	inf := "setsession"
	method := "POST"
	s.Module = &module
	s.Param = &inf
	s.Sign = t.Sign
	s.Method = &method
	args := make(map[string]interface{})
	args["uid"] = uid
	args["isadm"] = isadm
	args["iscomp"] = iscomp
	args["readon"] = readon
	argst := ToMapStringAny(args)
	s.Args = argst

	ans := GRPCConnect(host, port, s)
	access_token = fmt.Sprintf("%v", ans["access_token"])
	refresh_token = fmt.Sprintf("%v", ans["refresh_token"])
	at_lifetime = int(ans["at_lifetime"].(float64))
	rt_lifetime = int(ans["rt_lifetime"].(float64))

	return access_token, at_lifetime, refresh_token, rt_lifetime
}
