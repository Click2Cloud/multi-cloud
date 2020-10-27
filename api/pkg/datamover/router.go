package datamover

import (
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/v2/client"
)

func RegisterRouter(ws *restful.WebService) {
	handler := NewAPIService(client.DefaultClient)
	ws.Route(ws.POST("/jobs/{id}/abort").To(handler.AbortJob)).Doc("Cancel job")
	ws.Route(ws.POST("/jobs/{id}/pause").To(handler.PauseJob)).Doc("Pause job")
	//ws.Route(ws.PUT("/{tenantId}/objectlifecycle").To(handler.ObjectLifeCycle)).
	//	Doc("Pause job")
}
