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

package cron

import (
	"fmt"

	. "auth/global"

	"time"

	"github.com/spf13/viper"
)

func Init() {

	n := 0
	for n != 1 {

		CronJob()

		time.Sleep(5 * time.Second)
		setingskey := fmt.Sprintf("%s.cron", MicroServiceName)
		isCron := viper.GetBool(setingskey)
		if !isCron {
			n = 1
		}

	}

}

func CronJob() {
	// Put your cron job codes here
}
