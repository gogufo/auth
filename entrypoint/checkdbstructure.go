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

package entrypoint

import (
	. "auth/model"
	"fmt"
	"time"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func CheckDBStructure() {
	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		SetErrorLog("dbstructure.go:81: " + err.Error())
		//return "error with db"
	}

	dbtype := viper.GetString("database.type")

	//Check if table users and roles exist
	if !db.Conn.Migrator().HasTable(&Users{}) {
		SetErrorLog("dbstructure.go:94: " + "Table users do not exist. Create table Users")
		//db.Conn.Debug().AutoMigrate(&Users{})
		//Create users table
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&Users{})
		} else {
			db.Conn.Migrator().CreateTable(&Users{})
		}

		//Add admin user
		//1. generate user hash
		userhash := Hashgen(8)
		//2. generate users Password
		userpass := RandomString(12)
		//2.1 generete pass passhash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userpass), 8)
		if err != nil {

			SetErrorLog("dbstructure.go:108: " + err.Error())
		}

		//3. Admin User email
		useremail := viper.GetString("email.address")

		user := Users{
			UID:           userhash,
			Name:          "admin",
			Pass:          string(hashedPassword),
			Mail:          useremail,
			Mailsent:      int(time.Now().Unix()),
			Mailconfirmed: int(time.Now().Unix()),
			Created:       int(time.Now().Unix()),
			Status:        true,
			Completed:     true,
			IsAdmin:       true,
		}
		/*
			role := Roles{
				UID:   userhash,
				Admin: true,
			}
		*/
		db.Conn.Create(&user)
		//db.Conn.Create(&role)

		ans := fmt.Sprintf("Admin User created!\t\nname: admin\t\npass: %s\t\n", userpass)

		fmt.Printf(ans)

	}

	if !db.Conn.Migrator().HasTable(&AuthHistory{}) {
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&AuthHistory{})
		} else {
			db.Conn.Migrator().CreateTable(&AuthHistory{})
		}
	}

}
