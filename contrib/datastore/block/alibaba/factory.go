// Copyright 2020 The SODA Authors.
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

package alibaba

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/opensds/multi-cloud/backend/pkg/utils/constants"
	backendpb "github.com/opensds/multi-cloud/backend/proto"
	"github.com/opensds/multi-cloud/contrib/datastore/block/driver"
)

type AlibabaBlockDriverFactory struct {
}

func (factory *AlibabaBlockDriverFactory) CreateBlockStorageDriver(backend *backendpb.BackendDetail) (driver.BlockDriver, error) {
	log.Infof("Entered to create Alibaba volume driver")

	ecsClient, err := ecs.NewClientWithAccessKey(string(backend.Region), string(backend.Access), string(backend.Security))
	if err != nil {
		log.Errorf("Error in creating the client")
		return nil, err
	}
	adapter := &AlibabaAdapter{Client: ecsClient}
	return adapter, nil
}

func init() {
	driver.RegisterDriverFactory(constants.BackendTypeAlibabaBlock, &AlibabaBlockDriverFactory{})
}
