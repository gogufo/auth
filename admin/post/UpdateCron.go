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
	"auth/cron"
	. "auth/global"
	"fmt"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func UpdateCron(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["action"] == nil {
		fmt.Printf("Missing important data")
		return ErrorReturn(t, 404, "000012", "Missing important data")
	}

	action := p.Sanitize(fmt.Sprintf("%v", args["action"]))
	setingskey := fmt.Sprintf("%s.cron", MicroServiceName)

	if action == "true" {
		viper.Set(setingskey, true)
		/// Run Cron
		go cron.Init()
	} else {
		viper.Set(setingskey, false)
	}

	ans["answer"] = action
	return Interfacetoresponse(t, ans)
}
