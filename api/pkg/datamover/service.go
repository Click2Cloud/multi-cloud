package datamover

import (
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/v2/client"
	backend "github.com/opensds/multi-cloud/backend/proto"
	dataflow "github.com/opensds/multi-cloud/dataflow/proto"
	datamover "github.com/opensds/multi-cloud/datamover/proto"
	s3 "github.com/opensds/multi-cloud/s3/proto"
	log "github.com/sirupsen/logrus"
)

const (
	backendService   = "backend"
	s3Service        = "s3"
	dataflowService  = "dataflow"
	datamoverService = "datamover"
)

type APIService struct {
	backendClient   backend.BackendService
	s3Client        s3.S3Service
	dataflowClient  dataflow.DataFlowService
	datamoverClient datamover.DatamoverService
}

func NewAPIService(c client.Client) *APIService {
	return &APIService{
		backendClient:   backend.NewBackendService(backendService, c),
		s3Client:        s3.NewS3Service(s3Service, c),
		dataflowClient:  dataflow.NewDataFlowService(dataflowService, c),
		datamoverClient: datamover.NewDatamoverService(datamoverService, c),
	}
}

func (s *APIService) AbortJob(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("id")
	log.Info("Received request Abort migration################################################################# \n", id)
	ctx := context.Background()

	res, err := s.datamoverClient.AbortJob(ctx, &datamover.AbortJobRequest{Id: id})

	if err != nil {
		response.WriteEntity(err)
		return
	}
	//For debug -- begin
	log.Print("Abort jobs reponse:%v\n", res)
	jsons, errs := json.Marshal(res)
	if errs != nil {
		log.Error(errs.Error())
	} else {
		log.Info("res: %s.\n", jsons)
	}
	//For debug -- end

	log.Debug("Abort job successfully.")
	response.WriteEntity(res)
}
func (s *APIService) PauseJob(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("id")
	log.Info("Received request jobs [id=%s] details.\n", id)
	ctx := context.Background()
	res, err := s.datamoverClient.PauseJob(ctx, &datamover.PauseJobRequest{Id: id})
	if err != nil {
		response.WriteEntity(err)
		return
	}

	//For debug -- begin
	log.Print("Pause jobs reponse:%v\n", res)
	jsons, errs := json.Marshal(res)
	if errs != nil {
		log.Print(errs)
	} else {
		log.Print("res: %s.\n", jsons)
	}
	//For debug -- end

	log.Print("Paused job successfully.")
	response.WriteEntity(res)
}
