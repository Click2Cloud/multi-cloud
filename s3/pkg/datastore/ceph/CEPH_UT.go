package ceph

import (
	mocki "github.com/opensds/multi-cloud/s3/pkg/datastore/ceph/ceph adapter/mocks"
	pb "github.com/opensds/multi-cloud/s3/proto"
	"github.com/stretchr/testify/mock"
	"testing"
	"context"
	"time"
)

type CephAdapterInterface struct {
mock.Mock
}
func (_m *CephAdapterInterface) BucketDelete(ctx context.Context, in *pb.Bucket) error {
args := _m.Called()
result := args.Get(0)
return result.(error)
}


func TestDeleteBucket(t *testing.T){
var req = &pb.Bucket{
Name: "Name",
}
ctx := context.Background()
deadline := time.Now().Add(time.Duration(50) * time.Second)
ctx, cancel := context.WithDeadline(ctx, deadline)
defer cancel()
mockRepoClient := new(mocki.CephAdapterInterface)
mockRepoClient.On("DeleteBucket", ctx, "Name").Return(nil)

ad := &CephAdapter{}
ad.session = mockRepoClient
testService := &ad
err := (*testService).BucketDelete(ctx, req)
t.Log(err)
mockRepoClient.AssertExpectations(t)
}
