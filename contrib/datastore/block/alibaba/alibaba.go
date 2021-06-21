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

package alibaba

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aws/aws-sdk-go/aws/awserr"
	block "github.com/opensds/multi-cloud/block/proto"
	"github.com/opensds/multi-cloud/contrib/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

type AlibabaAdapter struct {
	Client *ecs.Client
}

// utils
func (ad *AlibabaAdapter) ParseVolume(volumeParsed *ecs.Disk) (*block.Volume, error) {

	meta := make(map[string]interface{})
	meta = map[string]interface{}{
		VolumeId:              &volumeParsed.DiskId,
		CreationTimeAtBackend: &volumeParsed.CreationTime,
		//PerformanceLevel:      volumeParsed.PerformanceLevel,
	}

	metadata, err := utils.ConvertMapToStruct(meta)
	if err != nil {
		log.Errorf("failed to convert metadata = [%+v] to struct", metadata, err)
		return nil, err
	}

	volume := &block.Volume{

		Name:               volumeParsed.DiskName,
		Description:        volumeParsed.Description,
		Size:               (int64)(volumeParsed.Size),
		Type:               volumeParsed.Category,
		Region:             volumeParsed.RegionId,
		AvailabilityZone:   volumeParsed.ZoneId,
		Status:             volumeParsed.Status,
		Iops:               (int64)(volumeParsed.IOPS),
		SnapshotId:         volumeParsed.SourceSnapshotId,
		MultiAttachEnabled: false,
		Encrypted:          volumeParsed.Encrypted,
		Metadata:           metadata,
	}
	if volumeParsed.Encrypted {
		volume.EncryptionSettings = map[string]string{
			KmsKeyId: volumeParsed.KMSKeyId,
		}
	}

	return volume, nil
}

func (ad *AlibabaAdapter) ParseTag(tagParsed []ecs.CreateDiskTag) ([]*block.Tag, error) {

	var tags []*block.Tag
	for _, tag := range tagParsed {
		tags = append(tags, &block.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return tags, nil
}

// Todo: custom code -> Start

func (ad *AlibabaAdapter) ParseDescribeTag(tagParsed ecs.TagsInDescribeDisks) ([]*block.Tag, error) {

	var tags []*block.Tag
	for _, tag := range tagParsed.Tag {
		tags = append(tags, &block.Tag{
			Key:   tag.TagKey,
			Value: tag.TagValue,
		})
	}
	return tags, nil
}

// Todo: custom code -> End

// volume functions
func (ad *AlibabaAdapter) CreateVolume(ctx context.Context, volume *block.CreateVolumeRequest) (*block.CreateVolumeResponse, error) {

	svc := ad.Client

	tags := []ecs.CreateDiskTag{}
	for _, tag := range volume.Volume.Tags {
		tags = append(tags, ecs.CreateDiskTag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	request := ecs.CreateCreateDiskRequest()
	request.RegionId = volume.Volume.Region
	request.ZoneId = volume.Volume.AvailabilityZone
	request.SnapshotId = volume.Volume.SnapshotId
	request.DiskName = volume.Volume.Name
	request.Size = requests.NewInteger(int(volume.Volume.Size))
	request.DiskCategory = volume.Volume.Type
	request.Description = volume.Volume.Description
	request.Encrypted = requests.NewBoolean(volume.Volume.Encrypted)
	request.Tag = &tags
	if request.Encrypted == "true" {
		request.KMSKeyId = volume.Volume.EncryptionSettings[KmsKeyId]
	}

	// volumetype ESSD is not yet supported by alibaba api hence, performancelevel variable is commented
	//	PL1: A single ESSD can deliver up to 50,000 random read/write IOPS.
	//	PL2: A single ESSD can deliver up to 100,000 random read/write IOPS.
	//	PL3: A single ESSD can deliver up to 1,000,000 random read/write IOPS

	//if volume.Volume.Type == "cloud_essd"{
	//	if volume.Volume.Iops > 0 && volume.Volume.Iops <= 1000 {
	//		request.PerformanceLevel = "PL0"
	//	}
	//	if volume.Volume.Iops > 1000 && volume.Volume.Iops <= 50000 {
	//		request.PerformanceLevel = "PL1"
	//	}
	//	if volume.Volume.Iops > 50000 && volume.Volume.Iops <= 100000 {
	//		request.PerformanceLevel = "PL2"
	//	}
	//	if volume.Volume.Iops > 100000 && volume.Volume.Iops <= 1000000 {
	//		request.PerformanceLevel = "PL3"
	//	}
	//}

	response, err := svc.CreateDisk(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Create Volume response = %+v", response)
	vol, err := ad.ParseVolume(&ecs.Disk{
		RegionId:         volume.Volume.Region,
		ZoneId:           volume.Volume.AvailabilityZone,
		SourceSnapshotId: volume.Volume.SnapshotId,
		DiskName:         volume.Volume.Name,
		Size:             int(volume.Volume.Size),
		Type:             volume.Volume.Type,
		Description:      volume.Volume.Description,
		Encrypted:        volume.Volume.Encrypted,
		DiskId:           response.DiskId,
		Status:           volume.Volume.Status,
		KMSKeyId:         volume.Volume.EncryptionSettings[KmsKeyId],
		IOPS:             (int)(volume.Volume.Iops),
		CreationTime:     time.Now().Format(time.RFC1123),
	})
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		vol.Tags, _ = ad.ParseTag(tags)
	}

	return &block.CreateVolumeResponse{
		Volume: vol,
	}, nil

}

func (ad *AlibabaAdapter) GetVolume(ctx context.Context, volume *block.GetVolumeRequest) (*block.GetVolumeResponse, error) {

	svc := ad.Client
	request := ecs.CreateDescribeDisksRequest()
	request.DiskName = volume.Volume.Name
	request.RegionId = volume.Volume.Region
	request.ZoneId = volume.Volume.AvailabilityZone
	request.Category = volume.Volume.Type
	//id,_ := fmt.Println(volume.Volume.Metadata.Fields[VolumeId].GetStringValue())
	//request.DiskIds = string(id)

	response, err := svc.DescribeDisks(request)
	if err != nil {
		log.Error(err)
		//	return nil, err
	}
	// Todo: custom code -> Start
	var diskResp ecs.Disk
	diskId := volume.Volume.Metadata.Fields[VolumeId].GetStringValue()
	if response.TotalCount > 0 {
		for i := 0; i < response.TotalCount; i++ {
			if response.Disks.Disk[i].DiskId == diskId {
				diskResp = response.Disks.Disk[i]
				break
			}
		}
	}
	// Todo: custom code -> End
	log.Infof("Get Volume response = %+v", response)

	vol, err := ad.ParseVolume(&ecs.Disk{
		DiskId: volume.Volume.Metadata.Fields[VolumeId].GetStringValue(),
		Size:   int(volume.Volume.Size),
		//IOPS:         (int)(volume.Volume.Iops),
		IOPS:         diskResp.IOPS,
		RegionId:     volume.Volume.Region,
		Description:  volume.Volume.Description,
		Type:         volume.Volume.Type,
		CreationTime: time.Now().Format(time.RFC1123),
		ZoneId:       volume.Volume.AvailabilityZone,
		DiskName:     volume.Volume.Name,
		//Status:       volume.Volume.Status,
		Status:    diskResp.Status,
		Encrypted: volume.Volume.Encrypted,
		KMSKeyId:  volume.Volume.EncryptionSettings[KmsKeyId],
	})
	// Todo: custom code line
	vol.Tags, _ = ad.ParseDescribeTag(diskResp.Tags)

	if err != nil {
		log.Error(err)
		//return nil, err
	}

	return &block.GetVolumeResponse{
		Volume: vol,
	}, nil

}

func (ad *AlibabaAdapter) ListVolume(ctx context.Context, volume *block.ListVolumeRequest) (*block.ListVolumeResponse, error) {

	svc := ad.Client
	request := ecs.CreateDescribeDisksRequest()
	response, err := svc.DescribeDisks(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var volumes []*block.Volume
	for _, vol := range response.Disks.Disk {
		fs, err := ad.ParseVolume(&vol)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		volumes = append(volumes, fs)
	}

	log.Debugf("List Volumes = %+v", response)

	return &block.ListVolumeResponse{
		Volumes: volumes,
	}, nil

}

func (ad *AlibabaAdapter) UpdateVolume(ctx context.Context, volume *block.UpdateVolumeRequest) (*block.UpdateVolumeResponse, error) {

	svc := ad.Client

	tags := []ecs.CreateDiskTag{}
	for _, tag := range volume.Volume.Tags {
		tags = append(tags, ecs.CreateDiskTag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	//if volume.Volume.Name != "" || volume.Volume.Description != "" {
	//
	//	request := ecs.CreateModifyDiskAttributeRequest()
	//	request.DiskId = volume.Volume.Metadata.Fields[VolumeId].GetStringValue()
	//	request.DiskName = volume.Volume.Name
	//	request.Description = volume.Volume.Description
	//
	//	response, err := svc.ModifyDiskAttribute(request)
	//	if err != nil {
	//		log.Errorf("Error in updating volume at modify disk attributes: %+v", response, err)
	//		if aerr, ok := err.(awserr.Error); ok {
	//			switch aerr.Code() {
	//			default:
	//				log.Errorf(aerr.Error())
	//			}
	//		} else {
	//			log.Errorf(err.Error())
	//		}
	//		return nil, err
	//	}
	//	log.Debugf("Update Volume response = %+v", response)
	//
	//}

	if &volume.Volume.Size != nil {

		request := ecs.CreateResizeDiskRequest()
		request.NewSize = requests.NewInteger(int(volume.Volume.Size))
		request.DiskId = volume.Volume.Metadata.Fields[VolumeId].GetStringValue()

		response, err := svc.ResizeDisk(request)
		if err != nil {
			log.Errorf("Error in updating volume: %+v", response, err)
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					log.Errorf(aerr.Error())
				}
			} else {
				log.Errorf(err.Error())
			}
			return nil, err
		}
		log.Debugf("Update Volume response = %+v", response)

	}

	vol, err := ad.ParseVolume(&ecs.Disk{
		Size:         int(volume.Volume.Size),
		CreationTime: time.Now().Format(time.RFC1123),
		DiskId:       volume.Volume.Metadata.Fields[VolumeId].GetStringValue(),
		Category:     volume.Volume.Type,
		Description:  volume.Volume.Description,
		Encrypted:    volume.Volume.Encrypted,
		KMSKeyId:     volume.Volume.Metadata.Fields[VolumeId].GetStringValue(),
	})
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		vol.Tags, _ = ad.ParseTag(tags)
	}

	return &block.UpdateVolumeResponse{
		Volume: vol,
	}, nil
}

func (ad *AlibabaAdapter) DeleteVolume(ctx context.Context, volume *block.DeleteVolumeRequest) (*block.DeleteVolumeResponse, error) {
	svc := ad.Client
	request := ecs.CreateDeleteDiskRequest()
	request.RegionId = volume.Volume.Region
	request.DiskId = volume.Volume.Metadata.Fields[VolumeId].GetStringValue()

	response, err := svc.DeleteDisk(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debugf("Delete volume response = %+v", response)

	return &block.DeleteVolumeResponse{}, nil

}

func (ad *AlibabaAdapter) Close() error {
	return nil
}
