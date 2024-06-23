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

package model

import (
	"gorm.io/gorm"
)

type SignInCred struct {
	Password string `db:"password"`
	Username string `db:"name"`
	Mail     string `db:"mail"`
	Lang     string `db:"lang"`
}

type Users struct {
	gorm.Model
	UID           string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;" json:"uid"` //userID
	Pass          string `gorm:"column:pass;type:varchar(128);NOT NULL;DEFAULT ''" json:"pass"`
	Name          string `gorm:"column:name;type:varchar(60);NOT NULL;DEFAULT '';UNIQUE" json:"name"`
	Mail          string `gorm:"column:mail;type:varchar(254);DEFAULT '';UNIQUE"  json:"mail"`
	Mailsent      int    `gorm:"column:mailsent;type:int;DEFAULT '0'" json:"mailsent"`
	Mailconfirmed int    `gorm:"column:mailconfirmed;:int;DEFAULT '0'" json:"mailconfirmed"`
	Created       int    `gorm:"column:created;type:int;DEFAULT '0'" json:"created"`
	Access        int    `gorm:"column:access;type:int;DEFAULT '0'" json:"access"`
	Login         int    `gorm:"column:login;type:int;DEFAULT '0'" json:"login"`
	IP            string `gorm:"column:ip;type:varchar(128); DEFAULT ''" json:"ip"`
	Status        bool   `gorm:"column:status;type:bool;DEFAULT 'false'" json:"status"`
	Completed     bool   `gorm:"column:completed;type:bool;DEFAULT 'false'" json:"completed"`
	IsAdmin       bool   `gorm:"column:is_admin;type:bool;DEFAULT 'false'" json:"isadmin"`
	Readonly      bool   `gorm:"column:readonly;type:bool;DEFAULT 'false'" json:"readonly"`
	TFA           bool   `gorm:"column:tfa;type:bool;DEFAULT false;" json:"tfa"`
	TFAType       string `gorm:"column:tfatype;type:varchar(60);DEFAULT '';" json:"tfatype"`
}

type AuthHistory struct {
	gorm.Model
	UID       string `gorm:"column:uid;type:varchar(60);DEFAULT '';" json:"uid"` //userID
	Login     int    `gorm:"column:login;type:int;DEFAULT '0'" json:"login"`
	IP        string `gorm:"column:ip;type:varchar(60);DEFAULT ''" json:"ip"`
	UserAgent string `gorm:"column:user_agent;type:varchar(254);DEFAULT ''" json:"user_agent"`
}
