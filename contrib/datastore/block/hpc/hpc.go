// Copyright 2020 The SODA Authors.
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

// This uses the  golang SDK by Huaweicloud and uses EVS client for all the ops

package hpc

import (
	"context"
	"errors"
	evs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/evs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/evs/v2/model"
	model2 "github.com/opensds/multi-cloud/block/pkg/model"
	"github.com/opensds/multi-cloud/block/proto"
	pb "github.com/opensds/multi-cloud/block/proto"
	"github.com/opensds/multi-cloud/contrib/datastore/block/common"
	"github.com/opensds/multi-cloud/contrib/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

type HpcAdapter struct {
	evsc *evs.EvsClient
}

// This fucntion takes an object of VolumeDetail from the API result
// VolumeDetail is the base struct in the model for huaweicloud-sdk-go-v3 service model
func (ad *HpcAdapter) ParseVolume(evsVolResp *model.VolumeDetail) (*block.Volume, error) {
	//Create the map for metadata
	meta := make(map[string]interface{})
	meta = map[string]interface{}{
		common.VolumeId:              evsVolResp.Id,
		common.CreationTimeAtBackend: evsVolResp.CreatedAt,
	}
	metadata, err := utils.ConvertMapToStruct(meta)
	if err != nil {
		log.Errorf("failed to convert metadata = [%+v] to struct", metadata, err)
		return nil, err
	}

	size64 := (int64)(evsVolResp.Size)
	volume := &block.Volume{
		Name:     evsVolResp.Name,
		Size:     size64 * utils.GB_FACTOR,
		Status:   evsVolResp.Status,
		Type:     evsVolResp.VolumeType,
		Metadata: metadata,
	}
	return volume, nil
}

// This is the helper function which basically based upon the Volume Type provided by
// API request, will return the corresponding CreateVolumeOptionVolumeType data of the SDK
func getVolumeTypeForHPC(volType string) (model.CreateVolumeOptionVolumeType, error) {
	if volType == "SAS" {
		return model.GetCreateVolumeOptionVolumeTypeEnum().SAS, nil
	} else if volType == "SATA" {
		return model.GetCreateVolumeOptionVolumeTypeEnum().SATA, nil
	} else if volType == "SSD" {
		return model.GetCreateVolumeOptionVolumeTypeEnum().SSD, nil
	} else if volType == "GPSSD" {
		return model.GetCreateVolumeOptionVolumeTypeEnum().GPSSD, nil
	} else {
		err := "wrong volumetype provided"
		return model.CreateVolumeOptionVolumeType{}, errors.New(err)
	}
}

// Todo: custom code -> Start

// This is the helper function, which is used to get Status of Job and required
// Entities (like VolumeId, etc.) if jobType is createVolume, and will return the Job Response from SDK
func getJobStat(jobID string, evsClient *evs.EvsClient) (*model.ShowJobResponse, error) {
	status := model.GetShowJobResponseStatusEnum()
	jobReq := &model.ShowJobRequest{JobId: jobID}
	for c := 0; c < 6; {
		jobResponse, err := evsClient.ShowJob(jobReq)
		if err != nil {
			return nil, err
		}
		switch *jobResponse.Status {
		case status.INIT:
			{
				time.Sleep(time.Second / 2)
				c++
				continue
			}
		case status.FAIL:
			{
				return nil, errors.New(*jobResponse.FailReason)
			}
		}
		return jobResponse, nil
	}
	return nil, errors.New("Error: Unable to create Volume, request timeout")
}

// This is the helper function to calculate IOPS for specific size of specific volumeType
// Returns the IOPS value
func getVolumeIOPs(volType string, volSize int64) int64 {
	const (
		// Formula to get VOLTYPE_10_BASE_IOPS = iops_of_volType_10GB
		// Formula to get VOLTYPE_DIF_VOLUME_SIZE = iops_of_volType_11GB - iops_of_volType_10GB

		MIN_DISK_SIZE int64 = 10

		SSD_10GB_BASE_IOPS  = 2300
		SSD_DIF_VOLUME_SIZE = 50

		GPSSD_10GB_BASE_IOPS  = 1920
		GPSSD_DIF_VOLUME_SIZE = 12

		SAS_10GB_BASE_IOPS  = 1880
		SAS_DIF_VOLUME_SIZE = 8

		//Todo: Please evaluate & set me after I'm in stock of huawei cloud
		//SATA_10GB_BASE_IOPS = ?
		//SATA_DIF_VOLUME_SIZE = ?
	)
	var volIOPs int64

	difSize := volSize - MIN_DISK_SIZE
	switch volSize {
	case 10:
		{
			switch volType {

			case "SSD":
				{
					volIOPs = SSD_10GB_BASE_IOPS
				}

			case "GPSSD":
				{
					volIOPs = GPSSD_10GB_BASE_IOPS
				}

			case "SAS":
				{
					volIOPs = SAS_10GB_BASE_IOPS
				}

				//case "SATA": {
				//	volIOPs = SATA_10GB_BASE_IOPS
				//}
			}
		}
	default:
		{
			switch volType {
			case "SSD":
				{
					volIOPs = SSD_10GB_BASE_IOPS + (SSD_DIF_VOLUME_SIZE * difSize)
				}

			case "GPSSD":
				{
					volIOPs = GPSSD_10GB_BASE_IOPS + (GPSSD_DIF_VOLUME_SIZE * difSize)
				}

			case "SAS":
				{
					volIOPs = SAS_10GB_BASE_IOPS + (SAS_DIF_VOLUME_SIZE * difSize)
				}

				//case "SATA": {
				//	volIOPs = SATA_10GB_BASE_IOPS + (SATA_DIF_VOLUME_SIZE * difSize)
				//}
			}
		}
	}
	return volIOPs
}

func (ad *HpcAdapter) ParseVolumeTag(tags []model2.Tag) ([]*pb.Tag, map[string]bool, error) {
	var convtags []*pb.Tag
	keys := make(map[string]bool)
	for _, tag := range tags {
		if _, value := keys[tag.Key]; !value {
			keys[tag.Key] = true
			convtags = append(convtags, &pb.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}
	}
	return convtags, keys, nil
}

func (ad *HpcAdapter) CreateVolumeTag(evsClient *evs.EvsClient, volume *block.UpdateVolumeRequest) error {
	var tags []model.Tag
	log.Info("Updating Tag Request: ", volume.Volume)
	for _, tagData := range volume.Volume.Tags {
		tags = append(tags, model.Tag{
			Key:   tagData.Key,
			Value: tagData.Value,
		})
	}
	log.Info("Tags Converted ", tags)

	action := model.GetBatchCreateVolumeTagsRequestBodyActionEnum()
	tagReqBody := &model.BatchCreateVolumeTagsRequestBody{
		Action: action.CREATE,
		Tags:   tags,
	}
	volumeCloudID := volume.Volume.Metadata.Fields[common.VolumeId].GetStringValue()

	createTagReq := &model.BatchCreateVolumeTagsRequest{
		VolumeId: volumeCloudID,
		Body:     tagReqBody,
	}
	createTagResp, err := evsClient.BatchCreateVolumeTags(createTagReq)
	if err != nil {
		log.Errorf("failed to update Tags %+v", err)
		return err
	}
	log.Info("Updated Tag's Response: ", createTagResp.String())
	return nil
}

func (ad *HpcAdapter) ResizeVolume(evsClient *evs.EvsClient, volume *block.UpdateVolumeRequest) error {
	log.Info("Updating Tag Request: ", volume.Volume)
	volumeId := volume.Volume.Metadata.Fields[common.VolumeId].GetStringValue()
	resizeVolumeRequest := &model.ResizeVolumeRequest{
		VolumeId: volumeId,
		Body: &model.ResizeVolumeRequestBody{
			OsExtend: &model.OsExtend{NewSize: int32(volume.Volume.Size)},
		},
	}

	resizeResp, err := evsClient.ResizeVolume(resizeVolumeRequest)
	if err != nil {
		return err
	}
	log.Info("Updated Disk Size Response ", resizeResp)
	_, err = getJobStat(*resizeResp.JobId, evsClient)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// Todo: custom code -> End

// Create EVS volume
func (ad *HpcAdapter) CreateVolume(ctx context.Context, volume *block.CreateVolumeRequest) (*block.CreateVolumeResponse, error) {
	evsClient := ad.evsc

	availabilityZone := volume.Volume.AvailabilityZone

	// The SDK model expects the size to be int32,
	size := (int32)(volume.Volume.Size / utils.GB_FACTOR)

	volumeType, typeErr := getVolumeTypeForHPC(volume.Volume.Type)
	if typeErr != nil {
		log.Error(typeErr)
		return nil, typeErr
	}

	name := volume.Volume.Name

	// Todo: custom code -> Start
	tags := make(map[string]string)
	for _, tag := range volume.Volume.Tags {
		tags[tag.Key] = tag.Value
	}
	// Todo: custom code -> End

	createVolOption := &model.CreateVolumeOption{
		AvailabilityZone: availabilityZone,
		Name:             &name,
		// Todo: custom code -> Start
		Description: &volume.Volume.Description,
		Tags:        tags,
		// Todo: custom code -> End
		Size:       &size,
		VolumeType: volumeType,
	}

	volRequest := &model.CreateVolumeRequest{
		Body: &model.CreateVolumeRequestBody{
			Volume: createVolOption},
	}

	// CreateVolume fromt the EVS client takes CreateVolumeRequest and the
	// Type hieararchy is CreateVolumeRequest -> Volume (CreateVolumeOption) -> VolumeDetails
	result, err := evsClient.CreateVolume(volRequest)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Submitted job for volume creation. Job is %+v", result)

	// Todo: Custom Code -> Start
	// New Approach to get VolumeID
	var volumeID string

	if result.VolumeIds == nil {
		jobResponse, err := getJobStat(*result.JobId, evsClient)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		volumeID = *jobResponse.Entities.VolumeId

		if len(volumeID) < 1 {
			return nil, errors.New("Not able to get volume Id")
		}
	} else {
		volIDs := *result.VolumeIds
		volumeID = volIDs[0]
	}

	vol, err := ad.ParseVolume(&model.VolumeDetail{
		Id:               volumeID,
		Name:             name,
		Status:           volume.Volume.Status,
		AvailabilityZone: availabilityZone,
		SnapshotId:       volume.Volume.SnapshotId,
		Description:      volume.Volume.Description,
		VolumeType:       volume.Volume.Type,
		CreatedAt:        time.Now().Format(time.RFC3339),
		Size:             size,
		Encrypted:        &volume.Volume.Encrypted,
		UserId:           volume.Volume.UserId,
		Tags:             tags,
	})

	vol.Iops = getVolumeIOPs(volume.Volume.Type, volume.Volume.Size)
	return &block.CreateVolumeResponse{
		Volume: vol,
	}, nil
	// Todo: Custom Code -> End

	// Todo: Default code
	// TODO: START
	// The job for Create Volume is submitted.
	// As of now there is no option to get the volume ID from the response
	// It only provides Job ID
	// FIXME: Making it to call list volumes and get the volume details

	//listVolumes, err := ad.ListVolume(nil, nil)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//volumes := listVolumes.Volumes
	//
	//
	//for _, volume := range volumes {
	//	if volume.Name == name {
	//        log.Debugf("Create Volume response %+v", volume)
	//		return &block.CreateVolumeResponse{
	//			Volume: volume,
	//		}, nil
	//	}
	//}
	//// TODO: END
	//return nil, errors.New("Error in creating volume")
}

// This will list all the volumes
func (ad *HpcAdapter) ListVolume(ctx context.Context, volume *block.ListVolumeRequest) (*block.ListVolumeResponse, error) {
	log.Infof("Entered to List all volumes")
	evsClient := ad.evsc

	listVolRequest := &model.ListVolumesRequest{}
	response, volErr := evsClient.ListVolumes(listVolRequest)
	if volErr != nil {
		log.Error(volErr)
		return nil, volErr
	}

	var volumes []*block.Volume
	for _, resObj := range *response.Volumes {
		volume, parseErr := ad.ParseVolume(&resObj)
		if parseErr != nil {
			log.Error(parseErr)
			return nil, parseErr
		}
		volumes = append(volumes, volume)
	}

	log.Debugf("Listing volumes: %+v", volumes)
	return &block.ListVolumeResponse{
		Volumes: volumes,
	}, nil
}

// Get the details of a particular volume
func (ad *HpcAdapter) GetVolume(ctx context.Context, volume *block.GetVolumeRequest) (*block.GetVolumeResponse, error) {
	log.Infof("Entered to get the volume details of %s", volume.Volume.Name)
	evsClient := ad.evsc

	// Get the VolumeId from the metadata for the requested volume
	volumeId := volume.Volume.Metadata.Fields[common.VolumeId].GetStringValue()
	showRequest := &model.ShowVolumeRequest{VolumeId: volumeId}
	showResponse, showErr := evsClient.ShowVolume(showRequest)
	if showErr != nil {
		log.Error(showErr)
		return nil, showErr
	}

	log.Infof("Get volume response: %+v", showResponse)
	vol, parseErr := ad.ParseVolume(showResponse.Volume)
	if parseErr != nil {
		log.Error(parseErr)
		return nil, parseErr
	}
	// Todo: custom code -> Start
	var tags []*pb.Tag
	for key, value := range showResponse.Volume.Tags {
		tags = append(tags, &pb.Tag{
			Key:   key,
			Value: value,
		})
	}
	vol.Tags = tags
	vol.Description = showResponse.Volume.Description
	// Todo: custom code -> End
	return &block.GetVolumeResponse{
		Volume: vol,
	}, nil
}

// Update the name of a volume
func (ad *HpcAdapter) UpdateVolume(ctx context.Context, volume *block.UpdateVolumeRequest) (*block.UpdateVolumeResponse, error) {
	evsClient := ad.evsc

	log.Info("Entered to update volume with details: %+v", volume)
	name := volume.Volume.Name
	updateVolOption := &model.UpdateVolumeOption{
		Name: &name,
		// Todo: custom code -> Start
		Description: &volume.Volume.Description,
		// Todo: custom code -> End
	}

	volumeId := volume.Volume.Metadata.Fields[common.VolumeId].GetStringValue()
	updateVolumeRequest := &model.UpdateVolumeRequest{
		VolumeId: volumeId,
		Body: &model.UpdateVolumeRequestBody{
			Volume: updateVolOption,
		},
	}

	updateResp, updateErr := evsClient.UpdateVolume(updateVolumeRequest)
	if updateErr != nil {
		log.Error(updateErr)
		//return nil, updateErr
	}
	log.Info("Updated Name & Description Response ", updateResp)

	// Todo: custom code -> Start

	// Creating Tags
	if volume.Volume.Size == 0 {
		volume.Volume.Size = int64(*updateResp.Size)
	}
	if len(volume.Volume.Tags) > 0 {
		err := ad.CreateVolumeTag(evsClient, volume)
		if err != nil {
			log.Error("Error While Updating Volume Tags: ", err)
		} else {
			log.Info("Updated Volume Tags successfully.")
		}
	}
	// Resizing Volume
	if volume.Volume.Size > 0 {
		err := ad.ResizeVolume(evsClient, volume)
		if err != nil {
			log.Error("Error While Updating Volume Size: ", err)
		} else {
			log.Info("Updated Volume Size successfully.")
		}
	}

	size := volume.Volume.Size
	volIOPs := getVolumeIOPs(*updateResp.VolumeType, size)
	// Todo: custom code -> End

	meta := map[string]interface{}{
		common.VolumeId: *updateResp.Id,
	}

	metadata, err := utils.ConvertMapToStruct(meta)
	if err != nil {
		log.Errorf("failed to convert metadata = [%+v] to struct", metadata, err)
		return nil, err
	}

	// Todo: Uncomment me if not using custom code
	//size := (int64)(*updateResp.Size * utils.GB_FACTOR)

	vol := &block.Volume{
		Size:     size,
		Status:   *updateResp.Status,
		Type:     *updateResp.VolumeType,
		Metadata: metadata,

		// Todo: custom code -> Start
		Iops: volIOPs,
		// Todo: custom code -> End
	}

	log.Debugf("Volume updated details: %+v", vol)
	return &block.UpdateVolumeResponse{
		Volume: vol,
	}, nil
}

// Delete the volume with a particular volume ID
func (ad *HpcAdapter) DeleteVolume(ctx context.Context, volume *block.DeleteVolumeRequest) (*block.DeleteVolumeResponse, error) {
	volumeId := volume.Volume.Metadata.Fields[common.VolumeId].GetStringValue()

	delRequest := &model.DeleteVolumeRequest{VolumeId: volumeId}

	client := ad.evsc
	delResponse, delErr := client.DeleteVolume(delRequest)
	if delErr != nil {
		log.Error(delErr)
		return nil, delErr
	}
	log.Debugf("Delete volume response = %+v", delResponse)

	return &block.DeleteVolumeResponse{}, nil
}

func (ad *HpcAdapter) Close() error {
	// TODO:
	return nil
}
