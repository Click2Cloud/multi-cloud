// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package cephs3mover

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	. "github.com/click2cloud-alpha/s3client"
	"github.com/click2cloud-alpha/s3client/models"
	"io/ioutil"
	"net/http"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/credentials"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/micro/go-log"
	. "github.com/opensds/multi-cloud/datamover/pkg/utils"
	pb "github.com/opensds/multi-cloud/datamover/proto"
	"io"
)

// DefaultDownloadPartSize is the default range of bytes to get at a time when
// using Download or Upload
const DefaultDownloadUploadPartSize = 1024 * 1024 * 5

// DefaultDownloadConcurrency is the default number of goroutines to spin up
// when using Download or Upload
const DefaultDownloadUploadConcurrency = 5

type Downloader struct {
	CephS3 *Client
}

func NewDownloader(c *Client, options ...func(*Downloader)) *Downloader {
	d := &Downloader{
		CephS3: c,
	}
	for _, option := range options {
		option(d)
	}

	return d
}

type S3Error struct {
	Code        int
	Description string
}

type FailPart struct {
	partNumber  int
	uploadId    string
	md5         string
	containType string
	body        io.ReadCloser
	length      int64
}

type CephS3Helper struct {
	accessKey string
	secretKey string
	endPoint  string
}

type Value struct {
	// AWS Access key ID
	AccessKeyID string

	// AWS Secret Access Key
	SecretAccessKey string

	// AWS Session Token
	Endpoint string
}

func (cephS3Helper *CephS3Helper) createClient() *Client {

	client := NewClient(cephS3Helper.endPoint, cephS3Helper.accessKey, cephS3Helper.secretKey)

	return client
}

type CreateMultipartUploadOutput struct {
	UploadID string
}

type CephS3Mover struct {
	downloader         *Client                      //for multipart download
	svc                *Uploads                     //for multipart upload
	multiUploadInitOut *CreateMultipartUploadOutput //for multipart upload
	//uploadId string //for multipart upload
	completeParts []*CompletePart //for multipart upload
}

func (myc *CephS3Helper) Retrieve() (Value, error) {
	cred := Value{
		AccessKeyID:     myc.accessKey,
		SecretAccessKey: myc.secretKey,
		Endpoint:        myc.endPoint,
	}
	return cred, nil
}

func (myc *CephS3Helper) IsExpired() bool {
	return false
}

func md5Content(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	value := base64.StdEncoding.EncodeToString(cipherStr)
	return value
}

//This function is used to get object contain type
func getFileContentTypeCephOrAWS(out []byte) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_ = len(out)

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func (mover *CephS3Mover) UploadObj(objKey string, destLoca *LocationInfo, buf []byte) error {
	log.Logf("[cephs3mover] UploadObj object, key:%s.", objKey)

	sess := NewClient(destLoca.EndPoint, destLoca.Access, destLoca.Security)
	bucket := sess.NewBucket()

	cephObject := bucket.NewObject(destLoca.BucketName)

	md5 := md5Content(buf)

	contentType, err := getFileContentTypeCephOrAWS(buf)

	length := int64(len(buf))

	body := ioutil.NopCloser(bytes.NewReader(buf))
	log.Logf("[cephs3mover] Try to upload, bucket:%s,obj:%s\n", destLoca.BucketName, objKey)

	for tries := 1; tries <= 3; tries++ {
		err = cephObject.Create(objKey, md5, string(contentType), length, body, models.PublicReadWrite)

		if err != nil {
			log.Logf("[cephs3mover] Upload object[%s] failed %d times, err:%v\n", objKey, tries, err)
			//s3error := S3Error{501, err.Error()}

			//return err
		} else {
			log.Logf("[cephs3mover] Upload object[%s] successfully.", objKey)
		}

	}
	log.Logf("[s3mover] Upload object, bucket:%s,obj:%s, should not be here.\n", destLoca.BucketName, objKey)
	return errors.New("internal error")
}

func (mover *CephS3Mover) DownloadObj(objKey string, srcLoca *LocationInfo, buf []byte) (size int64, err error) {
	log.Logf("[cephs3mover] DownloadObj object, key:%s.", objKey)
	sess := NewClient(srcLoca.EndPoint, srcLoca.Access, srcLoca.Security)
	bucket := sess.NewBucket()
	cephObject := bucket.NewObject(srcLoca.BucketName)
	res, err := cephObject.Get(objKey, nil)
	var numBytes int64
	if err != nil {
		log.Logf("[cephs3mover]download object[bucket:%s,key:%s] failed, err:%v\n", srcLoca.BucketName, objKey, err)
	} else {
		numBytes = res.ContentLength
		log.Logf("[cephs3mover]downlad object[bucket:%s,key:%s] succeed, bytes:%d\n", srcLoca.BucketName, objKey, numBytes)
	}

	return numBytes, err

}

func (mover *CephS3Mover) MultiPartDownloadInit(srcLoca *LocationInfo) error {

	sess := NewClient(srcLoca.EndPoint, srcLoca.Access, srcLoca.Security)
	mover.downloader = sess

	log.Logf("[cephs3mover] MultiPartDownloadInit succeed.")

	return nil
}

func (mover *CephS3Mover) DownloadRange(objKey string, srcLoca *LocationInfo, buf []byte, start int64, end int64) (size int64, err error) {
	log.Logf("[cephs3mover] Download object[%s] range[%d - %d]...\n", objKey, start, end)
	sess := NewClient(srcLoca.EndPoint, srcLoca.Access, srcLoca.Security)
	bucket := sess.NewBucket()

	//newObjectKey :=srcLoca.BakendName +"/"+objKey
	cephObject := bucket.NewObject(srcLoca.BucketName)
	err = cephObject.SetACL(objKey, models.PublicReadWrite)
	if err != nil {
		log.Logf("fail to assign acl")
	}

	var getObjectOption GetObjectOption
	rangeObj := Range{Begin: start, End: end}
	getObjectOption = GetObjectOption{
		Range: &rangeObj,
	}
	for tries := 1; tries <= 3; tries++ {
		resp, err := cephObject.Get(objKey, &getObjectOption)
		d, err := ioutil.ReadAll(resp.Body)
		data := []byte(d)
		size = int64(len(data))
		if err != nil {
			log.Logf("[cephs3mover] Download object[%s] range[%d - %d] faild %d times, err:%v\n",
				objKey, start, end, tries, err)
			if tries == 3 {
				return 0, err
			}
		} else {
			log.Logf("[cephs3mover] Download object[%s] range[%d - %d] succeed, bytes:%d\n", objKey, start, end, size)
			return size, err
		}
	}

	log.Logf("[cephs3mover] Download object[%s] range[%d - %d], should not be here.\n", objKey, start, end)
	return 0, errors.New("internal error")
}

func (mover *CephS3Mover) MultiPartUploadInit(objKey string, destLoca *LocationInfo) error {
	sess := NewClient(destLoca.EndPoint, destLoca.Access, destLoca.Security)
	bucket := sess.NewBucket()
	cephObject := bucket.NewObject(destLoca.BucketName)
	mover.svc = cephObject.NewUploads(objKey)

	//s3c := s3Cred{ak: destLoca.Access, sk: destLoca.Security}
	//creds := credentials.NewCredentials(&s3c)
	//sess, err := session.NewSession(&aws.Config{
	//	Region:      aws.String(destLoca.Region),
	//	Endpoint:    aws.String(destLoca.EndPoint),
	//	Credentials: creds,
	//})
	//if err != nil {
	//	log.Logf("[cephs3mover] New session failed, err:%v\n", err)
	//	return err
	//}
	//
	//mover.svc = s3.New(sess)
	//multiUpInput := &s3.CreateMultipartUploadInput{
	//	Bucket: aws.String(destLoca.BucketName),
	//	Key:    aws.String(objKey),
	//}
	//resp, err := mover.svc.CreateMultipartUpload(multiUpInput)
	log.Logf("[s3mover] Try to init multipart upload[objkey:%s].\n", objKey)
	for tries := 1; tries <= 3; tries++ {
		resp, err := mover.svc.Initiate(nil)
		if err != nil {
			log.Logf("[s3mover] Init multipart upload[objkey:%s] failed %d times.\n", objKey, tries)
			if tries == 3 {
				return err
			}
		} else {
			mover.multiUploadInitOut = &CreateMultipartUploadOutput{resp.UploadID}
			log.Logf("[s3mover] Init multipart upload[objkey:%s] successfully, UploadId:%s\n", objKey, resp.UploadID)
			return nil
		}
	}

	//log.Logf("[s3mover] Init multipart upload[objkey:%s], should not be here.\n", objKey)
	//return errors.New("internal error")
	//if err != nil {
	//	log.Logf("[cephs3mover] Init multipart upload[objkey:%s] failed, err:%v\n", objKey, err)
	//	return errors.New("[cephs3mover] Init multipart upload failed.")
	//} else {
	//	log.Logf("[cephs3mover] Init multipart upload[objkey:%s] succeed, UploadId:%s\n", objKey, resp.UploadID)
	//}
	//
	////mover.uploadId = *resp.UploadId
	//mover.multiUploadInitOut = &CreateMultipartUploadOutput{resp.UploadID}
	log.Logf("[s3mover] Init multipart upload[objkey:%s], should not be here.\n", objKey)
	return errors.New("internal error")

}

func (mover *CephS3Mover) UploadPart(objKey string, destLoca *LocationInfo, upBytes int64, buf []byte, partNumber int64, offset int64) error {
	log.Logf("[cephs3mover] Upload range[objkey:%s, partnumber#%d,offset#%d,upBytes#%d,uploadid#%s]...\n", objKey, partNumber,
		offset, upBytes, mover.multiUploadInitOut.UploadID)
	//tries := 1
	//sess := NewClient(destLoca.EndPoint, destLoca.Access, destLoca.Security)
	//bucket := sess.NewBucket()
	//cephObject := bucket.NewObject(destLoca.BucketName)
	//destUploader := cephObject.NewUploads(objKey)
	//upPartInput := &s3.UploadPartInput{
	//	Body:          bytes.NewReader(buf),
	//	Bucket:        aws.String(destLoca.BucketName),
	//	Key:           aws.String(objKey),
	//	PartNumber:    aws.Int64(partNumber),
	//	UploadId:      aws.String(*mover.multiUploadInitOut.UploadId),
	//	ContentLength: aws.Int64(upBytes),
	//}
	md5 := md5Content(buf)
	contentType, err := getFileContentTypeCephOrAWS(buf)
	if err != nil {
		log.Logf("[cephs3mover]. err:%v\n", err)
	}

	length := int64(len(buf))

	data := []byte(buf)
	body := ioutil.NopCloser(bytes.NewReader(data))

	for tries := 1; tries <= 3; tries++ {
		upRes, err := mover.svc.UploadPart(int(partNumber), mover.multiUploadInitOut.UploadID, md5, contentType, length, body)
		if err != nil {
			log.Logf("[s3mover] Upload range[objkey:%s, partnumber#%d, offset#%d] failed %d times, err:%v\n",
				objKey, partNumber, offset, tries, err)
			if tries == 3 {
				return err
			}
		} else {

			//part := s3client.CompletePart{Etag: upRes.Etag, PartNumber:upRes.PartNumber}

			//mover.completeParts = part

			mover.completeParts = append(mover.completeParts, upRes)
			log.Logf("[s3mover] Upload range[objkey:%s, partnumber#%d,offset#%d] successfully.\n", objKey, partNumber, offset)
			return nil
		}
	}
	log.Logf("[s3mover] Upload range[objkey:%s, partnumber#%d, offset#%d], should not be here.\n", objKey, partNumber, offset)
	return errors.New("internal error")
}

func (mover *CephS3Mover) AbortMultipartUpload(objKey string, destLoca *LocationInfo) error {
	log.Logf("[cephs3mover] Aborting multipart upload[objkey:%s] for uploadId#%s.\n", objKey, mover.multiUploadInitOut.UploadID)
	bucket := mover.downloader.NewBucket()
	cephObject := bucket.NewObject(destLoca.BucketName)
	uploader := cephObject.NewUploads(objKey)
	err := uploader.RemoveUploads(mover.multiUploadInitOut.UploadID)
	if err != nil {
		log.Logf("abortMultipartUpload failed, err:%v\n", err)
		return err
	} else {
		log.Logf("abortMultipartUpload successfully\n")
	}
	return nil
}

func (mover *CephS3Mover) CompleteMultipartUpload(objKey string, destLoca *LocationInfo) error {
	//sess := NewClient(destLoca.EndPoint, destLoca.Access, destLoca.Security)
	//bucket := sess.NewBucket()
	//cephObject := bucket.NewObject(destLoca.BucketName)
	//destUploader := cephObject.NewUploads(objKey)

	//completeInput := &s3.CompleteMultipartUploadInput{
	//	Bucket:   aws.String(destLoca.BucketName),
	//	Key:      aws.String(objKey),
	//	UploadId: aws.String(*mover.multiUploadInitOut.UploadId),
	//	MultipartUpload: &s3.CompletedMultipartUpload{
	//		Parts: mover.completeParts,
	//	},
	//}

	log.Logf("[s3mover] Try to do CompleteMultipartUpload [objkey:%s].\n", objKey)
	for tries := 1; tries <= 3; tries++ {
		var completeParts []CompletePart
		for _, p := range mover.completeParts {
			completePart := CompletePart{
				Etag:       p.Etag,
				PartNumber: int(p.PartNumber),
			}
			completeParts = append(completeParts, completePart)
		}
		rsp, err := mover.svc.Complete(mover.multiUploadInitOut.UploadID, completeParts)
		if err != nil {
			log.Logf("[s3mover] completeMultipartUpload [objkey:%s] failed %d times, err:%v\n", objKey, tries, err)
			if tries == 3 {
				return err
			}
		} else {
			log.Logf("[s3mover] completeMultipartUpload successfully [objkey:%s], rsp:%v\n", objKey, rsp)
			return nil
		}
	}

	//completeInput := &s3.CompleteMultipartUploadInput{
	//	Bucket:   aws.String(destLoca.BucketName),
	//	Key:      aws.String(objKey),
	//	UploadId: aws.String(*mover.multiUploadInitOut.UploadId),
	//	MultipartUpload: &s3.CompletedMultipartUpload{
	//		Parts: mover.completeParts,
	//	},
	//}
	//
	//rsp, err := mover.svc.CompleteMultipartUpload(completeInput)
	//if err != nil {
	//	log.Logf("[cephs3mover] completeMultipartUploadS3 failed [objkey:%s], err:%v\n", objKey, err)
	//} else {
	//	log.Logf("[cephs3mover] completeMultipartUploadS3 successfully [objkey:%s], rsp:%v\n", objKey, rsp)
	//}
	//
	//return err
	log.Logf("[s3mover] completeMultipartUpload [objkey:%s], should not be here.\n", objKey)
	return errors.New("internal error")
}

func (mover *CephS3Mover) DeleteObj(objKey string, loca *LocationInfo) error {
	sess := NewClient(loca.EndPoint, loca.Access, loca.Security)
	bucket := sess.NewBucket()
	cephObject := bucket.NewObject(loca.BucketName)

	err := cephObject.Remove(objKey)

	if err != nil {
		log.Logf("[cephs3mover] Error occurred while waiting for object[%s] to be deleted.\n", objKey)
		return err
	}

	log.Logf("[cephs3mover] Delete Object[%s] successfully.\n", objKey)
	return nil
}

func ListObjs(loca *LocationInfo, filt *pb.Filter) ([]models.GetBucketResponseContent, error) {

	sess := NewClient(loca.EndPoint, loca.Access, loca.Security)
	bucket := sess.NewBucket()
	object := bucket.NewObject(loca.BakendName)

	//cephObject := bucket.NewObject(loca.BucketName)
	//bucket := connection.NewBucket()

	output, err := bucket.Get(string(loca.BucketName), "", "", "", 1000)
	if err != nil {
		log.Logf("[s3mover] List objects failed, err:%v\n", err)
		return nil, err
	}

	if filt != nil {
		output.Prefix = filt.Prefix
	}

	objs := output.Contents

	size := len(objs)
	for i := 0; i < size; i++ {
		objs = append(objs, output.Contents...)
		err := object.SetACL(objs[i].Key, models.PublicReadWrite)
		if err != nil {
			log.Logf("[s3mover] Set ACL failed, err:%v\n", err)
			return nil, err
		}
	}
	log.Logf("[s3mover] Number of objects in bucket[%s] is %d.\n", loca.BucketName, len(objs))
	return output.Contents, nil
}
