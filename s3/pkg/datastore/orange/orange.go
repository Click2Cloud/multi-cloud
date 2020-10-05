package orange

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	backendpb "github.com/opensds/multi-cloud/backend/proto"
	. "github.com/opensds/multi-cloud/s3/error"
	dscommon "github.com/opensds/multi-cloud/s3/pkg/datastore/common"
	"github.com/opensds/multi-cloud/s3/pkg/model"
	osdss3 "github.com/opensds/multi-cloud/s3/pkg/service"
	"github.com/opensds/multi-cloud/s3/pkg/utils"
	pb "github.com/opensds/multi-cloud/s3/proto"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

type OrangeAdapter struct {
	backend *backendpb.BackendDetail
	session *session.Session
}

func (or *OrangeAdapter) Put(ctx context.Context, stream io.Reader, object *pb.Object) (dscommon.PutResult, error) {
	bucket := or.backend.BucketName
	objectId := object.BucketName + "/" + object.ObjectKey
	result := dscommon.PutResult{}
	userMd5 := dscommon.GetMd5FromCtx(ctx)
	size := object.Size
	log.Infof("put object[OBS], objectId:%s, bucket:%s, size=%d, userMd5=%s\n", objectId, bucket, size, userMd5)

	// Limit the reader to its provided size if specified.
	var limitedDataReader io.Reader
	if size > 0 { // request.ContentLength is -1 if length is unknown
		limitedDataReader = io.LimitReader(stream, size)
	} else {
		limitedDataReader = stream
	}
	md5Writer := md5.New()
	dataReader := io.TeeReader(limitedDataReader, md5Writer)

	if object.Tier == 0 {
		// default
		object.Tier = utils.Tier1
	}
	storClass, err := osdss3.GetNameFromTier(object.Tier, utils.OSTYPE_ORANGE)
	if err != nil {
		log.Infof("translate tier[%d] to aws storage class failed\n", object.Tier)
		return result, ErrInternalError
	}

	uploader := s3manager.NewUploader(or.session)
	input := &s3manager.UploadInput{
		Body:         dataReader,
		Bucket:       aws.String(bucket),
		Key:          aws.String(objectId),
		StorageClass: aws.String(storClass),
	}
	if userMd5 != "" {
		md5Bytes, err := hex.DecodeString(userMd5)
		if err != nil {
			log.Warnf("user input md5 is abandoned, cause decode md5 failed, err:%v\n", err)
		} else {
			input.ContentMD5 = aws.String(base64.StdEncoding.EncodeToString(md5Bytes))
			log.Debugf("input.ContentMD5=%s\n", *input.ContentMD5)
		}
	}
	log.Infof("upload object[ORANGE OSS] start, objectId:%s\n", objectId)
	ret, err := uploader.Upload(input)
	if err != nil {
		log.Errorf("put object[ORANGE OSS] failed, objectId:%s, err:%v\n", objectId, err)
		return result, ErrPutToBackendFailed
	}
	log.Infof("put object[ORANGE OSS] end, objectId:%s\n", objectId)

	calculatedMd5 := hex.EncodeToString(md5Writer.Sum(nil))
	log.Debug("calculatedMd5:", calculatedMd5, ", userMd5:", userMd5)
	if userMd5 != "" && userMd5 != calculatedMd5 {
		log.Error("### MD5 not match, calculatedMd5:", calculatedMd5, ", userMd5:", userMd5)
		return result, ErrBadDigest
	}

	if ret.VersionID != nil {
		result.Meta = *ret.VersionID
	}
	result.UpdateTime = time.Now().Unix()
	result.ObjectId = objectId
	result.Etag = calculatedMd5
	result.Written = size
	log.Infof("put object[ORANGE OSS] successfully, objectId:%s, UpdateTime is:%v\n", objectId, result.UpdateTime)

	return result, nil
}

func (or *OrangeAdapter) Get(ctx context.Context, object *pb.Object, start int64, end int64) (io.ReadCloser, error) {
	bucket := or.backend.BucketName
	objectId := object.ObjectId
	getObjectInput := awss3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objectId,
	}
	log.Infof("get object[ORANGE OSS], objectId:%s, start = %d, end = %d\n", objectId, start, end)
	if start != 0 || end != 0 {
		strStart := strconv.FormatInt(start, 10)
		strEnd := strconv.FormatInt(end, 10)
		rangestr := "bytes=" + strStart + "-" + strEnd
		getObjectInput.SetRange(rangestr)
	}

	svc := awss3.New(or.session)
	result, err := svc.GetObject(&getObjectInput)
	if err != nil {
		log.Errorf("get object[ORANGE OSS] failed, objectId:%s, err:%v", objectId, err)
		return nil, ErrGetFromBackendFailed
	}

	log.Infof("get object[ORANGE OSS] succeed, objectId:%s, ContentLength:%d\n", objectId, *result.ContentLength)
	return result.Body, nil
}

func (or *OrangeAdapter) Delete(ctx context.Context, input *pb.DeleteObjectInput) error {
	bucket := or.backend.BucketName
	objectId := input.Bucket + "/" + input.Key
	deleteInput := awss3.DeleteObjectInput{Bucket: &bucket, Key: &objectId}

	log.Infof("delete object[ORANGE OSS], objectId:%s.\n", objectId)
	svc := awss3.New(or.session)
	_, err := svc.DeleteObject(&deleteInput)
	if err != nil {
		log.Errorf("delete object[ORANGE OSS] failed, objectId:%s, err:%v.\n", objectId, err)
		return ErrDeleteFromBackendFailed
	}

	log.Infof("delete object[ORANGE OSS] succeed, objectId:%s.\n", objectId)

	return nil
}

func (or *OrangeAdapter) Copy(ctx context.Context, stream io.Reader, target *pb.Object) (result dscommon.PutResult, err error) {
	log.Errorf("copy[ORANGE OSS] is not supported.")
	err = ErrInternalError
	return
}

func (or *OrangeAdapter) InitMultipartUpload(ctx context.Context, object *pb.Object) (*pb.MultipartUpload, error) {
	bucket := or.backend.BucketName
	objectId := object.BucketName + "/" + object.ObjectKey
	log.Infof("init multipart upload[ORANGE OSS], bucket = %v,objectId = %v\n", bucket, objectId)

	storClass, err := osdss3.GetNameFromTier(object.Tier, utils.OSTYPE_ORANGE)
	if err != nil {
		log.Warnf("translate tier[%d] to aws storage class failed, use default value.\n", object.Tier)
		return nil, ErrInternalError
	}

	multipartUpload := &pb.MultipartUpload{}
	multiUpInput := &awss3.CreateMultipartUploadInput{
		Bucket:       &bucket,
		Key:          &objectId,
		StorageClass: aws.String(storClass),
	}

	svc := awss3.New(or.session)
	res, err := svc.CreateMultipartUpload(multiUpInput)
	if err != nil {
		log.Fatalf("init multipart upload[ORANGE OSS] failed, err:%v\n", err)
		return nil, ErrBackendInitMultipartFailed
	} else {
		log.Infof("init multipart upload[ORANGE OSS] succeed, UploadId:%s\n", *res.UploadId)
		multipartUpload.Bucket = object.BucketName
		multipartUpload.Key = object.ObjectKey
		multipartUpload.UploadId = *res.UploadId
		multipartUpload.ObjectId = objectId
		return multipartUpload, nil
	}
}

func (or *OrangeAdapter) UploadPart(ctx context.Context, stream io.Reader, multipartUpload *pb.MultipartUpload, partNumber int64, upBytes int64) (*model.UploadPartResult, error) {
	bucket := or.backend.BucketName
	bytess, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Errorf("read data failed, err:%v\n", err)
		return nil, ErrInternalError
	}
	upPartInput := &awss3.UploadPartInput{
		Body:          bytes.NewReader(bytess),
		Bucket:        &bucket,
		Key:           &multipartUpload.ObjectId,
		PartNumber:    aws.Int64(partNumber),
		UploadId:      &multipartUpload.UploadId,
		ContentLength: aws.Int64(upBytes),
	}
	log.Infof("upload part[ORANGE OSS], input:%v\n", *upPartInput)

	svc := awss3.New(or.session)
	upRes, err := svc.UploadPart(upPartInput)
	if err != nil {
		log.Errorf("upload part[ORANGE OSS] failed. err:%v\n", err)
		return nil, ErrPutToBackendFailed
	} else {
		log.Infof("upload object[ORANGE OSS], objectId:%s, part #%d succeed, ETag:%s\n", multipartUpload.ObjectId,
			partNumber, *upRes.ETag)
		result := &model.UploadPartResult{
			Xmlns:      model.Xmlns,
			ETag:       *upRes.ETag,
			PartNumber: partNumber}
		return result, nil
	}

	log.Error("upload part[ORANGE OSS]: should not be here.")
	return nil, ErrInternalError
}

func (or *OrangeAdapter) CompleteMultipartUpload(ctx context.Context, multipartUpload *pb.MultipartUpload, completeUpload *model.CompleteMultipartUpload) (*model.CompleteMultipartUploadResult, error) {
	bucket := or.backend.BucketName
	log.Infof("complete multipart upload[AWS S3], bucket:%s, objectId:%s.\n", bucket, multipartUpload.ObjectId)

	var completeParts []*awss3.CompletedPart
	for _, p := range completeUpload.Parts {
		completePart := &awss3.CompletedPart{
			ETag:       aws.String(p.ETag),
			PartNumber: aws.Int64(p.PartNumber),
		}
		completeParts = append(completeParts, completePart)
	}
	completeInput := &awss3.CompleteMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &multipartUpload.ObjectId,
		UploadId: &multipartUpload.UploadId,
		MultipartUpload: &awss3.CompletedMultipartUpload{
			Parts: completeParts,
		},
	}

	log.Infof("completeInput %v\n", *completeInput)
	svc := awss3.New(or.session)
	resp, err := svc.CompleteMultipartUpload(completeInput)
	if err != nil {
		log.Errorf("complete multipart upload[ORANGE OSS] failed, err:%v\n", err)
		return nil, ErrBackendCompleteMultipartFailed
	}
	result := &model.CompleteMultipartUploadResult{
		Xmlns:    model.Xmlns,
		Location: *resp.Location,
		Bucket:   multipartUpload.Bucket,
		Key:      multipartUpload.Key,
		ETag:     *resp.ETag,
	}

	log.Infof("complete multipart upload[ORANGE OSS] successfully, resp:%v\n", resp)
	return result, nil
}

func (or *OrangeAdapter) AbortMultipartUpload(ctx context.Context, multipartUpload *pb.MultipartUpload) error {
	bucket := or.backend.BucketName
	log.Infof("abort multipart upload[ORANGE OSS], bucket:%s, objectId:%s.\n", bucket, multipartUpload.ObjectId)

	abortInput := &awss3.AbortMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &multipartUpload.ObjectId,
		UploadId: &multipartUpload.UploadId,
	}

	svc := awss3.New(or.session)
	rsp, err := svc.AbortMultipartUpload(abortInput)
	if err != nil {
		log.Errorf("abort multipart upload[ORANGE OSS] failed, err:%v\n", err)
		return ErrBackendAbortMultipartFailed
	}

	log.Infof("abort multipart upload[ORANGE OSS] successfully, rsp:%v\n", rsp)
	return nil
}

func (or *OrangeAdapter) ListParts(ctx context.Context, multipartUpload *pb.ListParts) (*model.ListPartsOutput, error) {
	return nil, errors.New("not implemented yet.")
}

func (or *OrangeAdapter) Close() error {
	return nil
}

func (or *OrangeAdapter) ChangeStorageClass(ctx context.Context, object *pb.Object, newClass *string) error {
	objectId := object.ObjectId
	log.Infof("change storage class[ORANGE OSS] of object[%s] to %s .\n", objectId, *newClass)

	svc := awss3.New(or.session)
	input := &awss3.CopyObjectInput{
		Bucket:     aws.String(or.backend.BucketName),
		Key:        aws.String(objectId),
		CopySource: aws.String(or.backend.BucketName + "/" + objectId),
	}
	input.StorageClass = aws.String(*newClass)
	_, err := svc.CopyObject(input)
	if err != nil {
		log.Errorf("change storage class[ORANGE OSS] of object[%s] to %s failed: %v.\n", objectId, *newClass, err)
		return ErrPutToBackendFailed
	}

	log.Infof("change storage class[ORANGE OSS] of object[%s] to %s succeed.\n", objectId, *newClass)
	return nil
}

type OssCred struct {
	ak string
	sk string
}

func (myc *OssCred) IsExpired() bool {
	return false
}

func (myc *OssCred) Retrieve() (credentials.Value, error) {
	cred := credentials.Value{AccessKeyID: myc.ak, SecretAccessKey: myc.sk}
	return cred, nil
}
