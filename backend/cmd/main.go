// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	//"github.com/micro/go-micro/v2/registry/mdns"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"github.com/micro/go-plugins/registry/mdns/v2"
	"github.com/opensds/multi-cloud/api/pkg/utils/obs"
	"github.com/opensds/multi-cloud/backend/pkg/db"
	handler "github.com/opensds/multi-cloud/backend/pkg/service"
	"github.com/opensds/multi-cloud/backend/pkg/utils/config"
	pb "github.com/opensds/multi-cloud/backend/proto"
	"os"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	db.Init(&config.Database{
		Driver:   "mongodb",
		Endpoint: dbHost})
	defer db.Exit()

	obs.InitLogs()

	regType := os.Getenv("MICRO_REGISTRY")
	var reg registry.Registry
	if regType == "mdns" {
		reg = mdns.NewRegistry()
	}
	if regType == "kubernetes" {
		reg = kubernetes.NewRegistry()
	}

	service := micro.NewService(
		micro.Name("backend"),
		micro.Registry(reg),
	)

	service.Init()

	pb.RegisterBackendHandler(service.Server(), handler.NewBackendService())
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
