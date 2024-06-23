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
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func SendTimeHash(t *pb.Request, hash string, uid string, parametr string, email string, lifetime int) {

	host := viper.GetString("server.internal_host")
	port := viper.GetString("server.grpc_port")

	s := &pb.Request{}
	module := "rights"
	method := "POST"
	inf := "settimehash"
	s.Module = &module
	s.Param = &inf
	s.Sign = t.Sign
	s.Method = &method
	args := make(map[string]interface{})
	args["hash"] = hash
	args["uid"] = uid
	args["parametr"] = parametr
	args["email"] = email
	args["lifetime"] = lifetime
	argst := ToMapStringAny(args)
	s.Args = argst

	GRPCConnect(host, port, s)

}
