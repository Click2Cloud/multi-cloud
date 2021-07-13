package hws

import (
	"github.com/opensds/multi-cloud/backend/pkg/utils/constants"
	backendpb "github.com/opensds/multi-cloud/backend/proto"
	"github.com/opensds/multi-cloud/s3/pkg/datastore/aws"
	"github.com/opensds/multi-cloud/s3/pkg/datastore/driver"
)

type HWObsDriverFactory struct {
}

func (factory *HWObsDriverFactory) CreateDriver(backend *backendpb.BackendDetail) (driver.StorageDriver, error) {
	awss3Fac := &aws.AwsS3DriverFactory{}
	return awss3Fac.CreateDriver(backend)
}

func init() {
	driver.RegisterDriverFactory(constants.BackendTypeObs, &HWObsDriverFactory{})
	driver.RegisterDriverFactory(constants.BackendFusionStorage, &HWObsDriverFactory{})
	driver.RegisterDriverFactory(constants.BackendTypeHpcBlock, &HWObsDriverFactory{})
}
