package orange

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/opensds/multi-cloud/backend/pkg/utils/constants"
	backendpb "github.com/opensds/multi-cloud/backend/proto"
	"github.com/opensds/multi-cloud/s3/pkg/datastore/driver"
)

type OrangeCloudDriverFactory struct {
}

func (factory *OrangeCloudDriverFactory) CreateDriver(backend *backendpb.BackendDetail) (driver.StorageDriver, error) {
	endpoint := backend.Endpoint
	AccessKeyID := backend.Access
	AccessKeySecret := backend.Security
	region := backend.Region

	OssCred := OssCred{ak: AccessKeyID, sk: AccessKeySecret}
	creds := credentials.NewCredentials(&OssCred)

	disableSSL := true
	sess, err := session.NewSession(&aws.Config{
		Region:      &region,
		Endpoint:    &endpoint,
		Credentials: creds,
		DisableSSL:  &disableSSL,
	})
	if err != nil {
		return nil, err
	}

	adap := &OrangeAdapter{backend: backend, session: sess}

	return adap, nil
}

func init() {
	driver.RegisterDriverFactory(constants.BackendTypeOrange, &OrangeCloudDriverFactory{})
}
