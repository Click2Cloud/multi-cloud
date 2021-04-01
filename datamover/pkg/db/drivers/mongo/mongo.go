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

package mongo

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	backend "github.com/opensds/multi-cloud/backend/pkg/model"
	. "github.com/opensds/multi-cloud/dataflow/pkg/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var adap = &adapter{}

var DataBaseName = "multi-cloud"
var CollJob = "job"
var CollBackend = "backends"

func Init(host string) *adapter {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	adap.s = session
	adap.userID = "unknown"

	return adap
}

func Exit() {
	adap.s.Close()
}

type adapter struct {
	s      *mgo.Session
	userID string
}

func (ad *adapter) GetJobStatus(jobID string) string {
	job := Job{}
	ss := ad.s.Copy()
	defer ss.Close()
	c := ss.DB(DataBaseName).C(CollJob)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(jobID)}).One(&job)
	if err != nil {
		log.Errorf("Get job[ID#%s] failed:%v.\n", jobID, err)
		return ""
	}

	return job.Status
}

func (ad *adapter) UpdateJob(job *Job) error {
	ss := ad.s.Copy()
	defer ss.Close()

	c := ss.DB(DataBaseName).C(CollJob)
	j := Job{}
	err := c.Find(bson.M{"_id": job.Id}).One(&j)
	if err != nil {
		log.Errorf("Get job[id:%v] failed before update it, err:%v\n", job.Id, err)

		return errors.New("Get job failed before update it.")
	}

	if !job.StartTime.IsZero() {
		j.StartTime = job.StartTime
	}
	if !job.EndTime.IsZero() {
		j.EndTime = job.EndTime
	}
	if job.TotalCapacity != 0 {
		j.TotalCapacity = job.TotalCapacity
	}
	if job.TotalCount != 0 {
		j.TotalCount = job.TotalCount
	}
	if job.PassedCount != 0 {
		j.PassedCount = job.PassedCount
	}
	if job.PassedCapacity != 0 {
		j.PassedCapacity = job.PassedCapacity
	}
	if job.Status != "" {
		j.Status = job.Status
	}
	if job.Progress != 0 {
		j.Progress = job.Progress
	}

	err = c.Update(bson.M{"_id": j.Id}, &j)
	if err != nil {
		log.Errorf("Update job in database failed, err:%v\n", err)
		return errors.New("Update job in database failed.")
	}

	log.Info("Update job in database succeed.")
	return nil
}

func (ad *adapter) GetBackendByName(name string) (*backend.Backend, error) {
	log.Infof("Get backend by name:%s\n", name)
	session := ad.s.Copy()
	defer session.Close()

	var backend = &backend.Backend{}
	collection := session.DB(DataBaseName).C(CollBackend)
	err := collection.Find(bson.M{"name": name}).One(backend)
	if err != nil {
		return nil, err
	}

	return backend, nil
}
func (ad *adapter) CancelJob(id string) (*Job, error) {
	job := Job{}
	ss := ad.s.Copy()
	defer ss.Close()
	c := ss.DB(DataBaseName).C(CollJob)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&job)
	if err == mgo.ErrNotFound {
		log.Println("Job does not exist.")
		return nil, ERR_JOB_NOT_EXIST
	}
	if job.Msg == "" {
		job.Msg = "Migration Aborted"
	} else {
		return nil, ERR_JOB_COMPLETED
	}
	j := Job{}
	err = c.Find(bson.M{"_id": job.Id}).One(&j)
	if err != nil {
		log.Println("Get job[id:%v] failed before update it, err:%v\n", job.Id, err)

		return nil, errors.New("Get job failed before update it.")
	}

	if !job.StartTime.IsZero() {
		j.StartTime = job.StartTime
	}
	if !job.EndTime.IsZero() {
		j.EndTime = job.EndTime
	}
	if job.TotalCapacity != 0 {
		j.TotalCapacity = job.TotalCapacity
	}
	if job.TotalCount != 0 {
		j.TotalCount = job.TotalCount
	}
	if job.PassedCount != 0 {
		j.PassedCount = job.PassedCount
	}
	if job.PassedCapacity != 0 {
		j.PassedCapacity = job.PassedCapacity
	}
	if job.Status != "" {
		j.Status = job.Status
	}
	if job.Progress != 0 {
		j.Progress = job.Progress
	}
	if job.TimeRequired != 0 {
		j.TimeRequired = job.TimeRequired
	}
	if job.TimeRequired == 0 {
		j.TimeRequired = 0
	}
	if job.Msg != "" {
		j.Msg = job.Msg
	}
	for k := range job.ObjList {
		if len(j.ObjList) < len(job.ObjList) {
			j.ObjList = append(j.ObjList, ObjDet{
				ObjKey:   job.ObjList[k].ObjKey,
				UploadId: "",
				Migrated: false,
			})
		}
	}

	err = c.Update(bson.M{"_id": j.Id}, &j)
	if err != nil {
		log.Println("Update job in database failed, err:%v\n", err)
		return nil, errors.New("Update job in database failed.")
	}

	log.Println("Update job in database succeedRRRRRRRRRRRRRRRR.")
	return &j, nil

	//var query mgo.Query;
}
func (ad *adapter) UpdateObjectStatus(jobId string, objKey string) error {
	ss := ad.s.Copy()
	defer ss.Close()

	c := ss.DB(DataBaseName).C(CollJob)
	j := Job{}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(jobId)}).One(&j)
	if err != nil {
		log.Println("Get job[id:%v] failed before update it, err:%v\n", jobId, err)

		return errors.New("Get job failed before update it.")
	}

	for k := range j.ObjList {
		if j.ObjList[k].ObjKey == objKey {
			j.ObjList[k].Migrated = true
			log.Println("UPDATE STATUS FOR %s_%s", j.ObjList[k].ObjKey, objKey)
			break
		}
	}

	err = c.Update(bson.M{"_id": j.Id}, &j)
	if err != nil {
		log.Println("Update job in database failed, err:%v\n", err)
		return errors.New("Update job in database failed.")
	}
	log.Println("UPDATE STATUS FOR OBJECT")
	return nil
}
func (ad *adapter) UpdateObjectList(job *Job) error {
	ss := ad.s.Copy()
	defer ss.Close()

	c := ss.DB(DataBaseName).C(CollJob)
	j := Job{}
	err := c.Find(bson.M{"_id": job.Id}).One(&j)
	if err != nil {
		log.Println("Get job[id:%v] failed before update it, err:%v\n", job.Id, err)

		return errors.New("Get job failed before update it.")
	}
	for m := range job.ObjList {
		for k := range j.ObjList {
			if j.ObjList[k].ObjKey == job.ObjList[m].ObjKey {
				if job.ObjList[m].Migrated == true {
					j.ObjList[k].Migrated = true
				}
				if job.ObjList[m].UploadId != "" {
					j.ObjList[k].UploadId = job.ObjList[m].UploadId
				}
				if job.ObjList[m].PartNo != 0 {
					j.ObjList[k].PartNo = job.ObjList[m].PartNo
				}
				if len(j.ObjList[k].PartTag) < len(job.ObjList[m].PartTag) {
					j.ObjList[k].PartTag = job.ObjList[m].PartTag
				}

				log.Println("UPDATE STATUS FOR %s %s %s", j.ObjList[k].ObjKey, j.ObjList[k].Migrated, job.ObjList[m].ObjKey)
				break
			}
		}
	}
	err = c.Update(bson.M{"_id": j.Id}, &j)
	if err != nil {
		log.Println("Update job in database failed, err:%v\n", err)
		return errors.New("Update job in database failed.")
	}

	log.Println("UPDATE OBJECT LIST IN DATABASE SUCCEED")
	return nil
}
func (ad *adapter) GetJobDetails(j *Job) error {
	ss := ad.s.Copy()
	defer ss.Close()

	c := ss.DB(DataBaseName).C(CollJob)
	job := Job{}
	err := c.Find(bson.M{"_id": j.Id}).One(&job)
	if err != nil {
		log.Println("Get job[id:%v] failed before update it, err:%v\n", job.Id, err)

		return errors.New("Get job failed before update it.")
	}

	if !job.StartTime.IsZero() {
		j.StartTime = job.StartTime
	} else {
		j.StartTime = time.Now()
	}

	if !job.EndTime.IsZero() {
		j.EndTime = job.EndTime
	}
	if job.TotalCapacity != 0 {
		j.TotalCapacity = job.TotalCapacity
	}
	if job.TotalCount != 0 {
		j.TotalCount = job.TotalCount
	}
	if job.PassedCount != 0 {
		j.PassedCount = job.PassedCount
	}
	if job.PassedCapacity != 0 {
		j.PassedCapacity = job.PassedCapacity
	}
	if job.Status != "" {
		j.Status = job.Status
	}
	if job.Progress != 0 {
		j.Progress = job.Progress
	}
	if job.TimeRequired != 0 {
		j.TimeRequired = job.TimeRequired
	}
	if job.MigratedCapacity != 0 {
		j.MigratedCapacity = job.MigratedCapacity
	}
	if job.TimeRequired == 0 {
		j.TimeRequired = 0
	}
	if job.Msg != "" {
		j.Msg = job.Msg
	} else {
		job.Msg = j.Msg
	}
	if j.Msg != "" {
		job.Msg = j.Msg
	}
	for k := range job.ObjList {
		if len(j.ObjList) < len(job.ObjList) {
			j.ObjList = append(j.ObjList, ObjDet{
				ObjKey:   job.ObjList[k].ObjKey,
				UploadId: job.ObjList[k].UploadId,
				Migrated: job.ObjList[k].Migrated,
				PartNo:   job.ObjList[k].PartNo,
			})
			for i := range job.ObjList[k].PartTag {
				j.ObjList[k].PartTag = append(j.ObjList[k].PartTag, job.ObjList[k].PartTag[i])
			}
		}
	}

	//err = c.Update(bson.M{"_id": j.Id}, &j)
	//if err != nil {
	//	log.Logf("Update job in database failed, err:%v\n", err)
	//	return errors.New("Update job in database failed.")
	//}

	log.Println("Update job in database succeed.")
	return nil
}
func (ad *adapter) GetObjectList(jobID string) []ObjDet {
	session := ad.s.Copy()
	defer session.Close()
	job := Job{}
	collection := session.DB(DataBaseName).C(CollJob)
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(jobID)}).One(&job)
	if err != nil {
		log.Println("Get job[ID#%s] failed:%v.\n", jobID, err)
	}
	fmt.Println(job.ObjList, " Objectdetail")
	return job.ObjList

}
