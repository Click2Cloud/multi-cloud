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

package migration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opensds/multi-cloud/api/pkg/common"
	"github.com/opensds/multi-cloud/backend/proto"
	"github.com/opensds/multi-cloud/dataflow/pkg/model"
	flowtype "github.com/opensds/multi-cloud/dataflow/pkg/model"
	"github.com/opensds/multi-cloud/datamover/pkg/db"
	. "github.com/opensds/multi-cloud/datamover/pkg/utils"
	pb "github.com/opensds/multi-cloud/datamover/proto"
	s3utils "github.com/opensds/multi-cloud/s3/pkg/utils"
	osdss3 "github.com/opensds/multi-cloud/s3/proto"
	log "github.com/sirupsen/logrus"
)

var simuRoutines = 10
var PART_SIZE int64 = 5 * 1024 * 1024 //The max object size that can be moved directly, default is 16M.
var JOB_RUN_TIME_MAX = 86400          //seconds, equals 1 day
var s3client osdss3.S3Service
var bkendclient backend.BackendService
var MiniSpeed int64 = 5 // 5KByte/Sec
var jobstate = make(map[string]string)

const WT_MOVE = 96
const WT_DELETE = 4
const JobType = "migration"

var (
	PENDING    = "pending"
	STARTED    = "started"
	VALIDATING = "validating"
	RUNNING    = "running"
	FAILED     = "failed"
	ABORTED    = "aborted"
	COMPLETED  = "completed"
	CANCELLED  = "cancelled"
	PAUSED     = "paused"
	RESUME     = "resumed"
)

type Migration interface {
	Init()
	HandleMsg(msg string)
}

func Init() {
	log.Infof("Migration init.")
	s3client = osdss3.NewS3Service("s3", client.DefaultClient)
	bkendclient = backend.NewBackendService("backend", client.DefaultClient)
}

func HandleMsg(msgData []byte) error {
	var job pb.RunJobRequest
	err := json.Unmarshal(msgData, &job)
	if err != nil {
		log.Infof("unmarshal failed, err:%v\n", err)
		return err
	}
	if jobstate[job.Id] == ABORTED || jobstate[job.Id] == CANCELLED {
		return nil
	} else {
		jobstate[job.Id] = PENDING
	}

	//Check the status of job, and run it if needed
	status := db.DbAdapter.GetJobStatus(job.Id)

	if status != flowtype.JOB_STATUS_PENDING {
		if status == flowtype.JOB_STATUS_RESUME {
			log.Print("***********************RESUMING***********************")

			log.Printf("job[id#%s] is %s.\n", job.Id, flowtype.JOB_STATUS_RESUME)
		} else {
			log.Printf("job[id#%s] is not in %s status.\n", job.Id, flowtype.JOB_STATUS_PENDING)
			return nil //No need to consume this message again}
		}
	}

	log.Infof("HandleMsg:job=%+v\n", job)
	go runjob(&job)
	return nil
}

func doMigrate(ctx context.Context, objs []*osdss3.Object, capa chan int64, th chan int, req *pb.RunJobRequest,
	job *flowtype.Job) {
	status := jobstate[job.Id.Hex()]

	if status == ABORTED {
		if job.Status != ABORTED {
			job.EndTime = time.Now()
			job.Status = ABORTED
			db.DbAdapter.UpdateJob(job)
		}
		return
	}
	if status == PAUSED {
		if job.Status != PAUSED {
			job.TimeRequired = 0
			job.EndTime = time.Now() //TODO Need to check
			job.Msg = "Migration Paused"
			job.Status = PAUSED
			db.DbAdapter.UpdateJob(job)
		}
		return
	}
	for i := 0; i < len(objs); i++ {
		db.DbAdapter.UpdateJob(job)
		if jobstate[job.Id.Hex()] == ABORTED {
			if job.Status != flowtype.JOB_STATUS_ABORTED {
				job.TimeRequired = 0
				job.EndTime = time.Now()
				job.Msg = "Migration Aborted"
				job.Status = flowtype.JOB_STATUS_ABORTED
				db.DbAdapter.UpdateJob(job)
			}
			break
		}
		if jobstate[job.Id.Hex()] == PAUSED {
			if job.Status != flowtype.JOB_STATUS_HOLD {
				job.TimeRequired = 0
				job.EndTime = time.Now() //TODO Need to check
				job.Msg = "Migration Paused"
				job.Status = flowtype.JOB_STATUS_HOLD
				db.DbAdapter.UpdateJob(job)
			}
			break
		}
		if objs[i].Tier == s3utils.Tier999 {
			// archived object cannot be moved currently
			log.Warnf("Object(key:%s) is archived, cannot be migrated.\n", objs[i].ObjectKey)
			continue
		}
		log.Infof("************Begin to move obj(key:%s)\n", objs[i].ObjectKey)
		//Create one routine
		var isMigrate = false
		if job.TimesResumed > 0 {
			isMigrate = CheckStatus(job, objs[i].ObjectKey)
		}
		if isMigrate == false {
			go migrate(ctx, objs[i], capa, th, req, job)
			th <- 1
			log.Infof("doMigrate: produce 1 routine, len(th):%d.\n", len(th))
		}
		if jobstate[job.Id.Hex()] == ABORTED && job.Status != flowtype.JOB_STATUS_ABORTED {
			db.DbAdapter.UpdateJob(job)

		}
		if jobstate[job.Id.Hex()] == PAUSED && job.Status != flowtype.JOB_STATUS_HOLD {

			db.DbAdapter.UpdateJob(job)

		}
	}
}

func CopyObj(ctx context.Context, obj *osdss3.Object, destLoca *LocationInfo, job *flowtype.Job) error {
	log.Infof("*****Move object, size is %d.\n", obj.Size)
	if obj.Size <= 0 {
		return nil
	}
	status := jobstate[job.Id.Hex()]

	if status == ABORTED {
		if job.Status != ABORTED {
			job.Status = ABORTED
			job.EndTime = time.Now()
			db.DbAdapter.UpdateJob(job)
		}
		return errors.New(job.Status)
	}
	if status == PAUSED {
		if job.Status != PAUSED {
			job.TimeRequired = 0
			job.EndTime = time.Now() //TODO Need to check
			job.Msg = "Migration Paused"
			job.Status = PAUSED
			db.DbAdapter.UpdateJob(job)
		}
		return errors.New(job.Status)
	}

	req := &osdss3.CopyObjectRequest{
		SrcObjectName:    obj.ObjectKey,
		SrcBucketName:    obj.BucketName,
		TargetBucketName: destLoca.BucketName,
		TargetObjectName: obj.ObjectKey,
	}
	tmoutSec := obj.Size / MiniSpeed
	opt := client.WithRequestTimeout(time.Duration(tmoutSec) * time.Second)
	_, err := s3client.CopyObject(ctx, req, opt)
	if err != nil {
		log.Errorf("copy object[%s] failed, err:%v\n", obj.ObjectKey, err)
	}

	progress(job, obj.Size, WT_MOVE)

	return err
}

func MultipartCopyObj(ctx context.Context, obj *osdss3.Object, destLoca *LocationInfo, job *flowtype.Job) error {
	if job == nil || job.Id.Hex() == "" {
		return errors.New("Job cannot be nil")
	}
	log.Debugf("obj.Size=%d, PART_SIZE=%d\n", obj.Size, PART_SIZE)
	partCount := int64(obj.Size / PART_SIZE)
	if obj.Size%PART_SIZE != 0 {
		partCount++
	}
	status := jobstate[job.Id.Hex()]
	if status == ABORTED {
		if job.Status != ABORTED {
			job.EndTime = time.Now()
			job.Status = ABORTED
			job.EndTime = time.Now()
			db.DbAdapter.UpdateJob(job)
		}
		return errors.New(job.Status)
	}
	if status == PAUSED {
		if job.Status != PAUSED {
			job.EndTime = time.Now()
			job.Status = PAUSED
			job.EndTime = time.Now()
			db.DbAdapter.UpdateJob(job)
		}
		return errors.New(job.Status)
	}
	log.Infof("*****Copy object[%s] from #%s# to #%s#, size=%d, partCount=%d.\n", obj.ObjectKey, obj.BucketName,
		destLoca.BucketName, obj.Size, partCount)

	var i int64
	var err error
	var uploadId string
	var initSucceed bool = false
	var completeParts []*osdss3.CompletePart
	var partNo int64 = 1
	var resMultipart = false
	currPartSize := PART_SIZE
	for m := range job.ObjList {
		if job.ObjList[m].ObjKey == obj.ObjectKey && job.ObjList[m].PartNo != 0 {
			partNo = job.ObjList[m].PartNo
			uploadId = job.ObjList[m].UploadId
			log.Printf("[INFO] MIGRATION RESUMING objKey:%s \n", obj.ObjectKey)
			log.Print("GOT PART NO for ", job.ObjList[m].ObjKey, job.ObjList[m].UploadId, job.ObjList[m])
			resMultipart = true
			break
		}
	}
	for i = partNo - 1; i < partCount; i++ {

		partNumber := i + 1

		offset := int64(i) * PART_SIZE
		if i+1 == partCount {
			currPartSize = obj.Size - offset
		}
		status1 := jobstate[job.Id.Hex()]
		if status1 == ABORTED {
			err = errors.New(ABORTED)
			break
		}
		if status1 == PAUSED {
			err = nil
			break
		}

		if partNumber == 1 || resMultipart == true {
			// init upload
			rsp, err := s3client.InitMultipartUpload(ctx, &osdss3.InitMultiPartRequest{
				BucketName: destLoca.BucketName,
				ObjectKey:  obj.ObjectKey,
				Tier:       destLoca.Tier,
				Location:   destLoca.BakendName,
				Attrs:      obj.CustomAttributes,
				// TODO: add content-type
			})
			if err != nil {
				log.Errorf("init mulipart upload failed:%v\n", err)
				break
			}
			initSucceed = true
			uploadId = rsp.UploadID
			log.Debugln("**** init multipart upload succeed, uploadId=", uploadId)
		}

		// copy part
		copyReq := &osdss3.CopyObjPartRequest{SourceBucket: obj.BucketName, SourceObject: obj.ObjectKey,
			TargetBucket: destLoca.BucketName, TargetObject: obj.ObjectKey, PartID: partNumber,
			UploadID: uploadId, ReadOffset: offset, ReadLength: currPartSize, TargetLocation: destLoca.BakendName,
		}

		var rsp *osdss3.CopyObjPartResponse
		try := 0
		tmoutSec := currPartSize / MiniSpeed
		opt := client.WithRequestTimeout(time.Duration(tmoutSec) * time.Second)
		for try < 3 { // try 3 times in case network is not stable
			status2 := jobstate[job.Id.Hex()]
			log.Debugf("###copy object part, objkey=%s, uploadid=%s, offset=%d, lenth=%d\n", obj.ObjectKey, uploadId, offset, currPartSize)
			if status2 != ABORTED {
				rsp, err = s3client.CopyObjPart(ctx, copyReq, opt)
				completePart := &osdss3.CompletePart{PartNumber: partNumber, ETag: rsp.Etag}
				completeParts = append(completeParts, completePart)
				//tempArr = append(tempArr, model.PartDet{
				//	Etag: rsp.Etag,
				//	No:   partNumber,
				//})

			} else if status2 == ABORTED {
				if job.Status != ABORTED {
					job.EndTime = time.Now()
					job.Status = ABORTED
					job.EndTime = time.Now()
					db.DbAdapter.UpdateJob(job)
				}
				break
			} else if status2 == PAUSED {
				if job.Status != PAUSED {
					job.EndTime = time.Now()
					job.Status = PAUSED
					job.EndTime = time.Now()
					db.DbAdapter.UpdateJob(job)
				}
				break
			}

			if err == nil {
				log.Debugln("copy part succeed")
				break
			} else {
				log.Warnf("copy part failed, err:%v\n", err)
			}
			try++
			time.Sleep(time.Second * 1)
		}
		if try == 3 {
			log.Errorln("copy part failed too many times")
			break
		}

		log.Debugf("copy part[obj=%s, uploadId=%s, ReadOffset=%d, ReadLength=%d] succeed\n", obj.ObjectKey,
			uploadId, offset, currPartSize)
		log.Println(rsp, "  Etag  ", rsp.Etag)

		// update job progress
		if job != nil {
			log.Debugln("update job")
			progress(job, currPartSize, WT_MOVE)
		}
		resMultipart = false
	}
	for j := range job.ObjList {
		//logger.Printf("job update Inside FOR LOOP objKey:%s PART: %d  \n", obj.ObjectKey, partNo)
		if job.ObjList[j].ObjKey == obj.ObjectKey {
			log.Printf("job update OBJECT IDENTIFIED objKey:%s PART: %d  \n", obj.ObjectKey, partNo)
			job.ObjList[j].UploadId = uploadId
			job.ObjList[j].PartNo = partNo
			for m := range completeParts {
				if len(job.ObjList[j].PartTag) == int(completeParts[m].PartNumber)-1 {
					job.ObjList[j].PartTag = append(job.ObjList[j].PartTag, completeParts[m])
				}
			}
			completeParts = job.ObjList[j].PartTag
			log.Print("UPDATE STATUS FOR M", job.ObjList[j].ObjKey, job.ObjList[j].Migrated)
			db.DbAdapter.UpdateJob(job)
			break
		}
	}
	if jobstate[job.Id.Hex()] == PAUSED && i != partCount {
		job.TimeRequired = 0
		job.Msg = "Migration Paused"
		job.Status = flowtype.JOB_STATUS_HOLD
		db.DbAdapter.UpdateJob(job)
		log.Printf("JOB PAUSED RETURNING POINT-9 objKey:%s PART: %d  \n", obj.ObjectKey, partNo)
		return errors.New(job.Msg)
	}

	StatusCheck := jobstate[job.Id.Hex()]

	if err == nil && StatusCheck != PAUSED {
		// copy parts succeed, need to complete it
		completeReq := &osdss3.CompleteMultipartRequest{BucketName: destLoca.BucketName, ObjectKey: obj.ObjectKey,
			UploadId: uploadId, CompleteParts: completeParts, SourceVersionID: obj.VersionId}
		if job == nil {
			// this is for lifecycle management
			completeReq.RequestType = s3utils.RequestType_Lifecycle
		}
		_, err = s3client.CompleteMultipartUpload(ctx, completeReq)
		if err != nil {
			log.Errorf("complete multipart copy failed, err:%v\n", err)
		}
	}

	if err != nil && StatusCheck != PAUSED {
		if initSucceed == true {
			log.Debugf("abort multipart copy, bucket:%s, object:%s, uploadid:%s\n", destLoca.BucketName, obj.ObjectKey, uploadId)
			_, ierr := s3client.AbortMultipartUpload(ctx, &osdss3.AbortMultipartRequest{BucketName: destLoca.BucketName,
				ObjectKey: obj.ObjectKey, UploadId: uploadId,
			})
			if ierr != nil {
				// it shoud be cleaned by gc in s3 service
				log.Warnf("abort multipart copy failed, err:%v\n", ierr)
			}
		}

		return err
	}

	log.Infof("*****Copy object[%s] from #%s# to #%s# succeed.\n", obj.ObjectKey, obj.Location,
		destLoca.BakendName)
	return nil
}

func deleteObj(ctx context.Context, obj *osdss3.Object) error {
	delMetaReq := osdss3.DeleteObjectInput{Bucket: obj.BucketName, Key: obj.ObjectKey}
	_, err := s3client.DeleteObject(ctx, &delMetaReq)
	if err != nil {
		log.Infof("delete object[bucket:%s,objKey:%s] failed, err:%v\n", obj.BucketName,
			obj.ObjectKey, err)
	} else {
		log.Infof("Delete object[bucket:%s,objKey:%s] successfully.\n", obj.BucketName,
			obj.ObjectKey)
	}

	return err
}

func migrate(ctx context.Context, obj *osdss3.Object, capa chan int64, th chan int, req *pb.RunJobRequest, job *flowtype.Job) {
	log.Infof("Move obj[%s] from bucket[%s] to bucket[%s].\n",
		obj.ObjectKey, job.SourceLocation, job.DestLocation)
	succeed := true
	needMove := true
	status := jobstate[job.Id.Hex()]
	if status == ABORTED {
		if job.Status != ABORTED {
			job.EndTime = time.Now()
			job.Status = ABORTED
			db.DbAdapter.UpdateJob(job)
		}
	}
	if status == PAUSED {
		if job.Status != PAUSED {
			job.TimeRequired = 0
			job.EndTime = time.Now() //TODO Need to check
			job.Msg = "Migration Paused"
			job.Status = PAUSED
			db.DbAdapter.UpdateJob(job)
		}
	}
	if needMove {
		PART_SIZE = GetMultipartSize()
		db.DbAdapter.UpdateJob(job)
		if jobstate[job.Id.Hex()] == ABORTED {
			if job.Status != flowtype.JOB_STATUS_ABORTED {
				job.TimeRequired = 0
				job.EndTime = time.Now()
				job.Msg = "Migration Aborted"
				job.Status = flowtype.JOB_STATUS_ABORTED
				db.DbAdapter.UpdateJob(job)
			}
		}
		if jobstate[job.Id.Hex()] == PAUSED {
			db.DbAdapter.UpdateJob(job)
			if job.Status != flowtype.JOB_STATUS_HOLD {
				job.TimeRequired = 0
				job.EndTime = time.Now() //TODO Need to check
				job.Msg = "Migration Paused"
				job.Status = flowtype.JOB_STATUS_HOLD
				db.DbAdapter.UpdateJob(job)
			}

		}
		var err error

		destLoc := &LocationInfo{BucketName: req.DestConn.BucketName, Tier: obj.Tier}
		if obj.Size <= PART_SIZE {
			err = CopyObj(ctx, obj, destLoc, job)
		} else {
			err = MultipartCopyObj(ctx, obj, destLoc, job)
		}

		if err != nil {
			succeed = false
		}

	}
	// copy object

	if succeed && !req.RemainSource {
		deleteObj(ctx, obj)
		// TODO: Need to clean if delete failed.
	}

	if succeed {
		for i := range job.ObjList {
			if job.ObjList[i].ObjKey == obj.ObjectKey {
				job.ObjList[i].Migrated = true
				log.Print("UPDATE STATUS FOR %s %s %s", job.ObjList[i].ObjKey, job.ObjList[i].Migrated)
				break
			}
		}
		//If migrate success, update capacity
		log.Infof("  migrate object[key=%s,versionid=%s] succeed.", obj.ObjectKey, obj.VersionId)
		t := <-th
		log.Printf("[INFO] migrate: consume %d routine, alen(th)=%d\n", t, len(th))
		log.Printf("[INFO] CAPACITY  bcapa(capa)=%d\n", len(capa))
		capa <- obj.Size
		log.Printf("[INFO] CAPACITY  Acapa(capa)=%d\n", len(capa))
		if job.Type == "migration" {
			progress(job, obj.Size, WT_DELETE)
		}
	} else {
		var t int
		if job.Status != flowtype.JOB_STATUS_ABORTED && job.Status != flowtype.JOB_STATUS_HOLD {
			log.Printf("[ERROR] migrate object[%s] failed.", obj.ObjectKey)
			t = <-th
			capa <- -1
		} else if job.Status == flowtype.JOB_STATUS_HOLD {
			log.Printf("[INFO] migrate object[%s] paused.", obj.ObjectKey)

			log.Printf("[INFO] PAUSED:  object[%s] , len(th)=%d\n", obj.ObjectKey, len(th))
			t = <-th
			capa <- -1
			log.Printf("[INFO] migrate: consume  alen(th)=%d\n", len(th))

		} else {
			log.Printf("[INFO] migrate object[%s] aborted.", obj.ObjectKey)

			log.Printf("[INFO] migrate: consume  blen(th)=%d\n", len(th))
			t = <-th
			capa <- -1
		}
		log.Printf("[ERROR] CAPACITY  fAILED bcapa=%d\n", len(capa))

		log.Printf("[ERROR] CAPACITY FAILED Acapa=%d\n", len(capa))

		log.Printf("[ERROR] FAILED/PAUSED/ABORT: consume %d routine, len(th)=%d\n", t, len(th))
	}

	//log.Printf("[INFO] FULL MIGR

}

func updateJob(j *flowtype.Job) {
	log.Println(j)
	for i := 1; i <= 3; i++ {
		err := db.DbAdapter.UpdateJob(j)
		log.Println(i)
		if err == nil {
			break
		}
		if i == 3 {
			log.Infof("update the finish status of job in database failed three times, no need to try more.")
		}
	}
}

func initJob(ctx context.Context, in *pb.RunJobRequest, j *flowtype.Job) error {
	j.Status = flowtype.JOB_STATUS_RUNNING
	j.SourceLocation = in.SourceConn.BucketName
	j.DestLocation = in.DestConn.BucketName
	j.Type = JobType

	// get total count and total size of objects need to be migrated
	totalCount, totalSize, err := countObjs(ctx, in)
	if err != nil || totalCount == 0 {
		if err != nil {
			j.Status = flowtype.JOB_STATUS_FAILED
		} else {
			j.Status = flowtype.JOB_STATUS_SUCCEED
		}
		j.EndTime = time.Now()
		updateJob(j)
		log.Infof("err:%v, totalCount=%d\n", err, totalCount)
		return errors.New("no need move")
	}

	j.TotalCount = totalCount
	j.TotalCapacity = totalSize
	updateJob(j)

	return nil
}

func runjob(in *pb.RunJobRequest) error {
	log.Infoln("Runjob is called in datamover service.")
	log.Infof("Request: %+v\n", in)
	log.Println(in.Id, "this is job ID ****************IN runjob")
	// set context tiemout
	ctx := metadata.NewContext(context.Background(), map[string]string{
		common.CTX_KEY_USER_ID:   in.UserId,
		common.CTX_KEY_TENANT_ID: in.TenanId,
	})
	// 60 means 1 minute, 2592000 means 30 days, 86400 means 1 day
	dur := GetCtxTimeout("JOB_MAX_RUN_TIME", 60, 2592000, 86400)
	_, ok := ctx.Deadline()
	if !ok {
		ctx, _ = context.WithTimeout(ctx, dur)
	}

	// init job
	j := flowtype.Job{Id: bson.ObjectIdHex(in.Id)}
	db.DbAdapter.GetJobDetails(&j)
	j.StartTime = time.Now()
	if j.Status == flowtype.JOB_STATUS_PENDING {
		j.TimesResumed = 0
	} else {
		log.Printf("[INFO] Migration Resumed")
		j.TimesResumed++
	}
	j.Status = flowtype.JOB_STATUS_RUNNING
	err := initJob(ctx, in, &j)
	if err != nil {
		return err
	}
	var limit int32 = 1000
	var marker string
	objs, err := getObjs(ctx, in, marker, limit)

	if len(j.ObjList) == 0 {

		for k := range objs {
			j.ObjList = append(j.ObjList, model.ObjDet{
				ObjKey:   objs[k].ObjectKey,
				UploadId: "",
				Migrated: false,
			})
			log.Println(j.ObjList, "objects ___________")
			updateJob(&j)
		}
	}

	totalObj := len(objs)
	if totalObj == 0 {
		log.Printf("[WARN] Bucket is empty.")
		j.Msg = "Bucket is empty"
		j.TimeRequired = int64(0)
	}
	var totalcount int64 = 0
	var totalcap int64 = 0
	for i := 0; i < totalObj; i++ {
		totalcount++
		totalcap += objs[i].Size
	}
	if j.TotalCount != totalcount {
		j.TotalCount = totalcount
	}
	if j.TotalCapacity != totalcap {
		j.TotalCapacity = totalcap
	}
	if err != nil || j.TotalCount == 0 || j.TotalCapacity == 0 {
		log.Printf("FAILED STATUS CHANGED-3")
		j.Status = flowtype.JOB_STATUS_FAILED
		j.EndTime = time.Now()
		j.TimeRequired = int64(0)
		updateJob(&j)
		return err
	}

	updateJob(&j)
	//hhhhh
	// used to transfer capacity(size) of objects
	capa := make(chan int64)
	// concurrent go routines is limited to be simuRoutines
	th := make(chan int, simuRoutines)

	for {

		//objs, err := getObjs(ctx, in, marker, limit)
		if err != nil {
			//update database
			j.Status = flowtype.JOB_STATUS_FAILED
			j.EndTime = time.Now()
			db.DbAdapter.UpdateJob(&j)
			return err
		}

		num := len(objs)
		if num == 0 {
			break
		}

		//Do migration for each object.
		go doMigrate(ctx, objs, capa, th, in, &j)

		if num < int(limit) {
			break
		}
		marker = objs[num-1].ObjectKey
	}

	var totalcapacity, capacity, count, passedCount, totalObjs int64 = j.TotalCapacity, j.PassedCapacity, j.PassedCount, j.PassedCount, j.TotalCount
	log.Println(capacity, count, passedCount, totalObjs, "<---------------------------------")
	tmout := false
	for {
		select {
		case c := <-capa:

			log.Debug(c, "line number 494")
			{ //if c is less than 0, that means the object is migrated failed.
				count++
				log.Println(count, "this is count of obect")
				if c >= 0 {
					passedCount++
					log.Println(count, "this is count of passedobect")
					capacity += c
				}

				//update database
				j.PassedCount = passedCount
				j.PassedCapacity = capacity
				Progress := int64((capacity / totalcapacity) * 100)
				log.Infof("ObjectMigrated:%d,TotalCapacity:%d Progress:%d\n", j.PassedCount, j.TotalCapacity, j.Progress)
				log.Println("This is progress=>", Progress, "This is totalcapacity=>", totalcapacity, "This is passedcapacity=>", j.PassedCapacity)
				updateJob(&j)

			}
		case <-time.After(dur):
			{
				tmout = true
				log.Warnln("Timout.")
			}
		}
		if count >= totalObjs || tmout {
			log.Debug(capa, th, "line number 515")
			log.Infof("break, capacity=%d, timout=%v, count=%d, passed count=%d\n", capacity, tmout, count, passedCount)
			close(capa)
			close(th)
			break
		}
	}

	var ret error = nil
	j.PassedCount = int64(passedCount)
	log.Println(in.Id, jobstate[in.Id], "this is status  ############################################################### IN run job")
	if passedCount < totalObjs {
		if jobstate[in.Id] == ABORTED {
			j.Status = ABORTED
		} else if j.Status == flowtype.JOB_STATUS_HOLD {
			log.Printf("[INFO] run job Paused: %s\n")
			log.Printf("[INFO] Migration Paused")
			j.Status = flowtype.JOB_STATUS_HOLD
		} else {
			errmsg := strconv.FormatInt(totalObjs, 10) + " objects, passed " + strconv.FormatInt(passedCount, 10)
			log.Infof("run job failed: %s\n", errmsg)
			ret = errors.New("failed")
			j.Status = flowtype.JOB_STATUS_FAILED
		}
	} else {
		j.Status = flowtype.JOB_STATUS_SUCCEED
		j.Progress = 100
	}

	j.EndTime = time.Now()
	for i := 1; i <= 3; i++ {
		err := db.DbAdapter.UpdateJob(&j)
		if err == nil {
			break
		}
		if i == 3 {
			log.Infof("update the finish status of job in database failed three times, no need to try more.")
		}
	}

	return ret
}

//To calculate Progress of migration process
func progress(job *flowtype.Job, size int64, wt int64) {
	// Migrated Capacity = Old_migrated capacity + WT(Process)*Size of Object/100
	log.Println(job.MigratedCapacity, "this is new log for migrated capacity", size, "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^", wt)
	MigratedCapacity := job.PassedCapacity + (size)*(wt/100)
	log.Println(MigratedCapacity, "new migration capacity")
	job.PassedCapacity = MigratedCapacity * 100 / 100
	// Progress = Migrated Capacity*100/ Total Capacity
	job.Progress = int64(job.PassedCapacity * 100 / job.TotalCapacity)
	log.Debugf("Progress %d, MigratedCapacity %d, TotalCapacity %d\n", job.Progress, job.MigratedCapacity, job.TotalCapacity)
	db.DbAdapter.UpdateJob(job)
}
func Abort(jobId string) (string, error) {
	j := flowtype.Job{Id: bson.ObjectIdHex(jobId)}

	if jobstate[jobId] == PAUSED {
		//	logger.Println("Migration Aborted Successfully.")
		log.Infof("Migration Aborted Successfully.")
	}
	log.Debug("i am here **************************************************************************************************************************************************")
	jobstate[jobId] = ABORTED
	j.Status = ABORTED
	db.DbAdapter.UpdateJob(&j)

	return j.Status, nil
}

func Pause(jobId string) (string, error) {
	j := flowtype.Job{Id: bson.ObjectIdHex(jobId)}

	jobstate[jobId] = PAUSED
	j.Status = flowtype.JOB_STATUS_HOLD
	//j.Msg = "Migration Paused"
	db.DbAdapter.UpdateJob(&j)
	log.Print("****************************************************************************PAUSED*********************************************************************************************************")
	return j.Status, nil
}
func CheckStatus(job *model.Job, objKey string) bool {
	var isMig = false
	for i := range job.ObjList {
		if job.ObjList[i].ObjKey == objKey {
			if job.ObjList[i].Migrated == true {
				log.Print("Object already Migrated %s %s %s", job.ObjList[i].ObjKey, job.ObjList[i].Migrated)
				isMig = true
				break
			}
			log.Print("MIGRATING OBJECT  %s %s %s", job.ObjList[i].ObjKey, job.ObjList[i].Migrated)
			break
		}
	}

	return isMig
}
func SizeCalulator(obj *osdss3.Object) string {

	objectsize := float64(obj.Size)
	kb := objectsize / 1024
	mb := kb / 1024
	gb := mb / 1024
	tb := gb / 1024

	if tb >= 1 {
		tbsize := math.Round(tb*100) / 100
		return fmt.Sprint(tbsize, " TB")
	} else if gb >= 1 {
		gbsize := math.Round(gb*100) / 100
		return fmt.Sprint(gbsize, " GB")
	} else if mb >= 1 {
		mbsize := math.Round(mb*100) / 100
		return fmt.Sprint(mbsize, " MB")
	} else if kb >= 1 {
		kbsize := math.Round(kb*100) / 100
		return fmt.Sprint(kbsize, " KB")
	} else {
		byteobjectsize := math.Round(objectsize*100) / 100
		return fmt.Sprint(byteobjectsize, " B")
	}

}
