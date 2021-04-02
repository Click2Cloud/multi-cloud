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

// DeleteInstances invokes the ecs.DeleteInstances API synchronously
// api document: https://help.aliyun.com/api/ecs/deleteinstances.html
func (client *Client) DeleteInstances(request *DeleteInstancesRequest) (response *DeleteInstancesResponse, err error) {
	response = CreateDeleteInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteInstancesWithChan invokes the ecs.DeleteInstances API asynchronously
// api document: https://help.aliyun.com/api/ecs/deleteinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteInstancesWithChan(request *DeleteInstancesRequest) (<-chan *DeleteInstancesResponse, <-chan error) {
	responseChan := make(chan *DeleteInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteInstances(request)
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

// DeleteInstancesWithCallback invokes the ecs.DeleteInstances API asynchronously
// api document: https://help.aliyun.com/api/ecs/deleteinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteInstancesWithCallback(request *DeleteInstancesRequest, callback func(response *DeleteInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteInstancesResponse
		var err error
		defer close(result)
		response, err = client.DeleteInstances(request)
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

// DeleteInstancesRequest is the request struct for api DeleteInstances
type DeleteInstancesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId       requests.Integer `position:"Query" name:"ResourceOwnerId"`
	InstanceId            *[]string        `position:"Query" name:"InstanceId"  type:"Repeated"`
	DryRun                requests.Boolean `position:"Query" name:"DryRun"`
	ResourceOwnerAccount  string           `position:"Query" name:"ResourceOwnerAccount"`
	ClientToken           string           `position:"Query" name:"ClientToken"`
	OwnerAccount          string           `position:"Query" name:"OwnerAccount"`
	TerminateSubscription requests.Boolean `position:"Query" name:"TerminateSubscription"`
	Force                 requests.Boolean `position:"Query" name:"Force"`
	OwnerId               requests.Integer `position:"Query" name:"OwnerId"`
}

// DeleteInstancesResponse is the response struct for api DeleteInstances
type DeleteInstancesResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteInstancesRequest creates a request to invoke DeleteInstances API
func CreateDeleteInstancesRequest() (request *DeleteInstancesRequest) {
	request = &DeleteInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DeleteInstances", "ecs", "openAPI")
	return
}

// CreateDeleteInstancesResponse creates a response to parse from DeleteInstances response
func CreateDeleteInstancesResponse() (response *DeleteInstancesResponse) {
	response = &DeleteInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
