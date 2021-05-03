package s3

import (
	"github.com/emicklei/go-restful"
	"github.com/opensds/multi-cloud/api/pkg/filters/signature"
	"github.com/prometheus/common/log"
)

func (s *APIService) RouteMoveToGlacier(request *restful.Request, response *restful.Response) {
	err := signature.PayloadCheck(request, response)
	if err != nil {
		WriteErrorResponse(response, request, err)
		return
	}
	log.Info("Route to move to Glacier")
	s.MoveToGlacier(request, response)

}
