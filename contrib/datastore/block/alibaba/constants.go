// Copyright 2020 The OpenSDS Authors.
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

const (
	PerformanceLevel = "PerformanceLevel"

	// The time that the volume was created in Alibaba.

	CreationTimeAtBackend = "CreationTimeAtBackend"

	// Indicates whether the volume was created using fast snapshot restore in Alibaba.
	FastRestored = "FastRestored"

	// Information about the volume iops.
	Iops = "Iops"

	// The ID of an Alibaba Key Management Service (Alibaba KMS) customer master key (CMK)
	// that was used to protect the encrypted volume.
	KmsKeyId = "KmsKeyId"

	// The Alibab Resource Name (ARN) of the Outpost.
	OutpostArn = "OutpostArn"

	// The ID of the volume in Alibaba.
	VolumeId = "VolumeId"

	// The type of the volume in Alibaba.
	VolumeType = "VolumeType"

	// The type of the volume in Alibaba.
	Progress = "Progress"

	// The modification completion or failure time at Alibaba.
	StartTimeAtBackend = "StartTimeAtBackend"

	// The modification completion or failure time at Alibaba.
	EndTimeAtBackend = "EndTimeAtBackend"

	// A status message about the modification progress or failure in Alibaba.
	StatusMessage = "StatusMessage"
)
