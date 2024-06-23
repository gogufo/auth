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

// Generated with GRPC Microservice Creator v1.3.0

package main

import (
	"auth/cron"
	"auth/entrypoint"
	. "auth/global"
	. "auth/version"
	"fmt"
	"net"
	"os"
	"time"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	. "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/certifi/gocertifi"
	"github.com/getsentry/sentry-go"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("./config/")

	//port := viper.GetInt("server.port")
	//	portsrting := fmt.Sprintf(":%d", port)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error with Settings \n")
		os.Exit(3)
	}

	if viper.GetBool("server.sentry") {

		SetLog("Connect to Setry...")

		sentryClientOptions := sentry.ClientOptions{
			Dsn:              viper.GetString("sentry.dsn"),
			EnableTracing:    viper.GetBool("sentry.tracing"),
			Debug:            viper.GetBool("sentry.debug"),
			TracesSampleRate: viper.GetFloat64("sentry.trace"),
		}

		rootCAs, err := gocertifi.CACerts()
		if err != nil {
			SetLog("Could not load CA Certificates for Sentry: " + err.Error())

		} else {
			sentryClientOptions.CaCerts = rootCAs
		}

		err = sentry.Init(sentryClientOptions)

		if err != nil {
			SetLog("Error with sentry.Init: " + err.Error())
		}

		flushsec := viper.GetDuration("sentry.flush")

		defer sentry.Flush(flushsec * time.Second)

	}

	portvar := fmt.Sprintf("microservices.%s.port", MicroServiceName)
	getport := viper.GetString(portvar)
	port := ":5300"
	if getport != "" {
		port = fmt.Sprintf(":%s", getport)
	}

	listener, err := net.Listen("tcp", port)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s := &Server{}

	setingskey := fmt.Sprintf("microservices.%s.entrypointversion", MicroServiceName)
	lastentrypointversion := viper.GetString(setingskey)

	if lastentrypointversion != VERSIONPLUGIN {
		go entrypoint.Init()
	}

	// Run Cron
	setingskey = fmt.Sprintf("microservices.%s.cron", MicroServiceName)
	isCron := viper.GetBool(setingskey)
	if isCron {
		go cron.Init()
	}

	pb.RegisterReverseServer(grpcServer, s)

	grpcServer.Serve(listener)

}

type Server struct {
}

func (s *Server) Do(c context.Context, request *pb.Request) (response *pb.Response, err error) {

	//Check for Sign
	sign := viper.GetString("server.sign")
	if sign != *request.Sign {
		return ErrorReturn(request, 401, "00001", "You are not authorized"), nil
	}

	if viper.GetBool("server.checktoken") {
		//Check Authorisation
		if request.UID == nil {
			return ErrorReturn(request, 401, "00002", "You are not authorized"), nil
		}

	}

	//Check connection

	response = Init(request)

	return response, nil
}
