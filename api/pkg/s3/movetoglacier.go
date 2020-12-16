package s3

import (
	"github.com/emicklei/go-restful"
	"github.com/opensds/multi-cloud/api/pkg/common"
	"github.com/opensds/multi-cloud/s3/pkg/utils"
	osdss3 "github.com/opensds/multi-cloud/s3/proto"
	log "github.com/sirupsen/logrus"
)

func (s *APIService) MoveToGlacier(request *restful.Request, response *restful.Response) {
	bucketName := request.PathParameter(common.REQUEST_PATH_BUCKET_NAME)
	objectKey := request.PathParameter(common.REQUEST_PATH_OBJECT_KEY)
	log.Println(bucketName, objectKey, "this is the request data")
	var tier int32 = 999
	ctx := common.InitCtxWithAuthInfo(request)
	object, _, _, err := s.getObjectMeta(ctx, bucketName, objectKey, "", false)
	log.Println(object, "this is the object data")
	if err != nil {
		WriteErrorResponse(response, request, err)
		log.Error(err)
		return
	}

	req := &osdss3.MoveObjectRequest{
		SrcObject:        objectKey,
		SrcObjectVersion: object.VersionId,
		SrcBucket:        bucketName,
		TargetTier:       tier,
		MoveType:         utils.MoveType_ChangeStorageTier,
	}
	_, err2 := s.s3Client.MoveObject(ctx, req)
	if err2 != nil {
		// if failed, it will try again in the next round schedule
		log.Errorf("Transition of %s failed:%v\n", objectKey, err)
		WriteErrorResponse(response, request, err)
	} else {
		log.Infof("Transition of %s succeed.\n", objectKey)
		response.WriteAsJson("Object Moved To Glacier")
	}

}
