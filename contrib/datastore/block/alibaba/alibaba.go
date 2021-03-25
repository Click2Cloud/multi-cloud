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

	//var tags []*block.Tag
	//for _, tag := range volumeParsed.Tags.Tag {
	//	tags = append(tags, &block.Tag{
	//		Key:   tag.TagKey,
	//		Value: tag.TagValue,
	//	})
	//}

	meta := make(map[string]interface{})
	meta = map[string]interface{}{
		VolumeId: &volumeParsed.DiskId,
		//CreationTimeAtBackend: &volumeParsed.CreationTime,
		//PerformanceLevel:      volumeParsed.PerformanceLevel,
	}

	metadata, err := utils.ConvertMapToStruct(meta)
	if err != nil {
		log.Errorf("failed to convert metadata = [%+v] to struct", metadata, err)
		return nil, err
	}

	volume := &block.Volume{

		Name:             volumeParsed.DiskName,
		Description:      volumeParsed.Description,
		Size:             (int64)(volumeParsed.Size),
		Type:             volumeParsed.Category,
		Region:           volumeParsed.RegionId,
		AvailabilityZone: volumeParsed.ZoneId,
		Status:           volumeParsed.Status,
		Iops:             (int64)(volumeParsed.IOPS),
		SnapshotId:       volumeParsed.SourceSnapshotId,
		Encrypted:        volumeParsed.Encrypted,
		Metadata:         metadata,
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

// volume functions
func (ad *AlibabaAdapter) CreateVolume(ctx context.Context, volume *block.CreateVolumeRequest) (*block.CreateVolumeResponse, error) {

	// Create a alibaba client.
	svc := ad.Client

	// we made tag of this type because alibaba request take this type of tag
	tags := []ecs.CreateDiskTag{}
	for _, tag := range volume.Volume.Tags {
		tags = append(tags, ecs.CreateDiskTag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	request := ecs.CreateCreateDiskRequest()

	request.DiskName = volume.Volume.Name
	request.Description = volume.Volume.Description
	request.Size = requests.NewInteger(int(volume.Volume.Size))
	request.DiskCategory = volume.Volume.Type
	request.RegionId = volume.Volume.Region
	request.ZoneId = volume.Volume.AvailabilityZone
	request.SnapshotId = volume.Volume.SnapshotId
	request.Tag = &tags
	//request.PerformanceLevel = strconv.Itoa(int(volume.Volume.Iops)) this parameter can be used for essd disk but default iops is in int64.
	request.Encrypted = requests.NewBoolean(volume.Volume.Encrypted)

	if request.Encrypted == "true" {
		request.KMSKeyId = volume.Volume.EncryptionSettings[KmsKeyId]
	}

	//if request.DiskCategory == "cloud_essd" {
	//	request.PerformanceLevel = volume.Volume.Metadata.Fields[PerformanceLevel].GetStringValue()
	//}else if request.DiskCategory == "cloud_ssd" || request.DiskCategory == "cloud_efficiency" {
	//	request.PerformanceLevel = ""
	//}

	response, err := svc.CreateDisk(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Create Volume response = %+v", response)
	vol, err := ad.ParseVolume(&ecs.Disk{
		DiskId:       response.DiskId,
		Size:         int(volume.Volume.Size),
		IOPS:         (int)(volume.Volume.Iops),
		RegionId:     volume.Volume.Region,
		Description:  volume.Volume.Description,
		Type:         volume.Volume.Type,
		CreationTime: time.Now().Format(time.RFC1123),
		ZoneId:       volume.Volume.AvailabilityZone,
		DiskName:     volume.Volume.Name,
		Status:       volume.Volume.Status,
		Encrypted:    volume.Volume.Encrypted,
		KMSKeyId:     volume.Volume.EncryptionSettings[KmsKeyId],
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
	volumeId := volume.Volume.Metadata.Fields[VolumeId].GetStringValue()
	//volumeIds := []string{volumeId}

	request.DiskIds = volumeId

	response, err := svc.DescribeDisks(request)

	log.Infof("Get Volume response = %+v", response)

	vol, err := ad.ParseVolume(&ecs.Disk{
		DiskId:       request.DiskIds,
		Size:         int(volume.Volume.Size),
		IOPS:         (int)(volume.Volume.Iops),
		RegionId:     volume.Volume.Region,
		Description:  volume.Volume.Description,
		Type:         volume.Volume.Type,
		CreationTime: time.Now().Format(time.RFC1123),
		ZoneId:       volume.Volume.AvailabilityZone,
		DiskName:     volume.Volume.Name,
		Status:       volume.Volume.Status,
		Encrypted:    volume.Volume.Encrypted,
		KMSKeyId:     volume.Volume.EncryptionSettings[KmsKeyId],
	})

	if err != nil {
		log.Error(err)
		return nil, err
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
	request := ecs.CreateResizeDiskRequest()
	request.RegionId = volume.Volume.Region
	request.NewSize = requests.NewInteger(int(volume.Volume.Size))
	request.DiskId = volume.Volume.Metadata.Fields[VolumeId].GetStringValue()
	request.Type = volume.Volume.Type

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

	vol, err := ad.ParseVolume(&ecs.Disk{
		Size:         int(volume.Volume.Size),
		CreationTime: time.Now().Format(time.RFC1123),
	})
	if err != nil {
		log.Error(err)
		return nil, err
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
