package ecs

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeFleetInstances invokes the ecs.DescribeFleetInstances API synchronously
// api document: https://help.aliyun.com/api/ecs/describefleetinstances.html
func (client *Client) DescribeFleetInstances(request *DescribeFleetInstancesRequest) (response *DescribeFleetInstancesResponse, err error) {
	response = CreateDescribeFleetInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeFleetInstancesWithChan invokes the ecs.DescribeFleetInstances API asynchronously
// api document: https://help.aliyun.com/api/ecs/describefleetinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeFleetInstancesWithChan(request *DescribeFleetInstancesRequest) (<-chan *DescribeFleetInstancesResponse, <-chan error) {
	responseChan := make(chan *DescribeFleetInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeFleetInstances(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeFleetInstancesWithCallback invokes the ecs.DescribeFleetInstances API asynchronously
// api document: https://help.aliyun.com/api/ecs/describefleetinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeFleetInstancesWithCallback(request *DescribeFleetInstancesRequest, callback func(response *DescribeFleetInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeFleetInstancesResponse
		var err error
		defer close(result)
		response, err = client.DescribeFleetInstances(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeFleetInstancesRequest is the request struct for api DescribeFleetInstances
type DescribeFleetInstancesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	FleetId              string           `position:"Query" name:"FleetId"`
}

// DescribeFleetInstancesResponse is the response struct for api DescribeFleetInstances
type DescribeFleetInstancesResponse struct {
	*responses.BaseResponse
	RequestId  string                            `json:"RequestId" xml:"RequestId"`
	TotalCount int                               `json:"TotalCount" xml:"TotalCount"`
	PageNumber int                               `json:"PageNumber" xml:"PageNumber"`
	PageSize   int                               `json:"PageSize" xml:"PageSize"`
	Instances  InstancesInDescribeFleetInstances `json:"Instances" xml:"Instances"`
}

// CreateDescribeFleetInstancesRequest creates a request to invoke DescribeFleetInstances API
func CreateDescribeFleetInstancesRequest() (request *DescribeFleetInstancesRequest) {
	request = &DescribeFleetInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeFleetInstances", "ecs", "openAPI")
	return
}

// CreateDescribeFleetInstancesResponse creates a response to parse from DescribeFleetInstances response
func CreateDescribeFleetInstancesResponse() (response *DescribeFleetInstancesResponse) {
	response = &DescribeFleetInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
