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

package pkg

import (
	"context"
	"fmt"
	"github.com/opensds/multi-cloud/dataflow/pkg/model"
	"github.com/opensds/multi-cloud/dataflow/pkg/utils"
	"github.com/opensds/multi-cloud/datamover/pkg/db"
	migration "github.com/opensds/multi-cloud/datamover/pkg/drivers/https"
	"github.com/opensds/multi-cloud/datamover/pkg/kafka"
	pb "github.com/opensds/multi-cloud/datamover/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var dataMoverGroup = "datamover"

type datamoverService struct{}

func NewDatamoverService() pb.DatamoverHandler {
	return &datamoverService{}
}
func InitDatamoverService() error {
	host := os.Getenv("DB_HOST")
	dbstor := utils.Database{Credential: "unkonwn", Driver: "mongodb", Endpoint: host}
	db.Init(&dbstor)

	addrs := []string{}
	config := strings.Split(os.Getenv("KAFKA_ADVERTISED_LISTENERS"), ";")
	for i := 0; i < len(config); i++ {
		addr := strings.Split(config[i], "//")
		if len(addr) != 2 {
			log.Info("invalid addr:", config[i])
		} else {
			addrs = append(addrs, addr[1])
		}
	}
	topics := []string{"migration", "lifecycle"}
	err := kafka.Init(addrs, dataMoverGroup, topics)
	if err != nil {
		log.Info("init kafka consumer failed.")
		return nil
	}
	go kafka.LoopConsume()

	datamoverID := os.Getenv("HOSTNAME")
	log.Infof("init datamover[ID#%s] finished.\n", datamoverID)
	return nil
}
func (b *datamoverService) AbortJob(ctx context.Context, in *pb.AbortJobRequest, out *pb.AbortJobResponse) error {
	log.Info("*******************************AbortJob is called in datamover service******************************************")
	//actx := c.NewContextFromJson(in.GetContext())
	if in.Id == "" {
		errmsg := fmt.Sprint("No id specified.")
		out.Err = errmsg
		return nil
	}
	jobstatus := db.DbAdapter.GetJobStatus(in.Id)

	if jobstatus == model.JOB_STATUS_ABORTED {
		out.Err = "job already aborted"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_CANCELLED {
		out.Err = "job already cancelled"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_SUCCEED {
		out.Err = "job already completed"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_FAILED {
		out.Err = "job current status is failed"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}

	jobs, err := migration.Abort(in.Id)

	if err != nil {
		log.Println("Get job err:%d.", err)
		out.Err = err.Error()
		return nil
	}
	out.Id = in.Id
	out.Status = jobs

	return nil
}
func (b *datamoverService) PauseJob(ctx context.Context, in *pb.PauseJobRequest, out *pb.PauseJobResponse) error {
	log.Println("Pause job is called in datamover service.")
	//actx := c.NewContextFromJson(in.GetContext())
	if in.Id == "" {
		errmsg := fmt.Sprint("No id specified.")
		out.Err = errmsg
		return nil
	}
	jobstatus := db.DbAdapter.GetJobStatus(in.Id)
	if jobstatus == model.JOB_STATUS_ABORTED {
		out.Err = "job already aborted"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_CANCELLED {
		out.Err = "job already cancelled"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_SUCCEED {
		out.Err = "job already completed"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_FAILED {
		out.Err = "job current status is failed"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus == model.JOB_STATUS_HOLD {
		out.Err = "job is already Paused"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	if jobstatus != model.JOB_STATUS_RUNNING {
		out.Err = "Job is not in running state"
		out.Id = in.Id
		out.Status = jobstatus
		return nil
	}
	jobs, err := migration.Pause(in.Id)

	if err != nil {
		log.Println("Get job err:%d.", err)
		out.Err = err.Error()
		return nil
	}
	out.Id = in.Id
	out.Status = jobs
	return nil
}
func (d *datamoverService) Runjob(ctx context.Context, request *pb.RunJobRequest, response *pb.RunJobResponse) error {
	panic("implement me")
}
func (d *datamoverService) DoLifecycleAction(ctx context.Context, request *pb.LifecycleActionRequest, resonse *pb.LifecycleActionResonse) error {
	panic("implement me")
}
