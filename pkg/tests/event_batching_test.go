/*
 * Copyright 2020-2022 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wingify/vwo-go-sdk/pkg/api"
	"github.com/wingify/vwo-go-sdk/pkg/mocks"
	"github.com/wingify/vwo-go-sdk/pkg/request"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

func init() {
	request.Client = mocks.MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
		},
	}
}

func GetVWOInstance(batchSize int, batchInterval int) api.VWOInstance {
	instance := api.VWOInstance{}
	instance.SettingsFile = schema.SettingsFile{}

	instance.IsBatchingEnabled = true

	instance.BatchEventQueue = schema.BatchEventQueue{
		RequestTimeInterval: batchInterval,
		EventsPerRequest:    batchSize,
	}

	vwoInstance := testdata.GetInstanceWithSettings("AB_T_100_W_33_33_33")
	instance.SettingsFile = vwoInstance.SettingsFile
	instance.Logger = vwoInstance.Logger
	instance.BatchEventQueue.Logger = instance.Logger
	instance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, instance.SettingsFile.Campaigns[0].Variations)

	return instance
}

func TestEnqueueAndFlushEvents(t *testing.T) {
	assertOutput := assert.New(t)
	batchSize, batchInterval := 10, 5
	instance := GetVWOInstance(batchSize, batchInterval)
	instance.BatchEventQueue.FlushCallBack = func(err error, batch []map[string]interface{}) {
		assertOutput.Equal(len(batch), 1)
	}

	userID, campaignKey := testdata.GetRandomUser(), "AB_T_100_W_33_33_33"

	instance.Activate(campaignKey, userID, nil)
	assertOutput.Equal(len(instance.BatchEventQueue.GetBatchImpressions()), 1)
	instance.FlushEvents()

	assertOutput.Nil(instance.BatchEventQueue.GetBatchImpressions())
}

func TestFlushQueueOnMaxEvents(t *testing.T) {
	assertOutput := assert.New(t)
	batchSize, batchInterval := 10, 5
	instance := GetVWOInstance(batchSize, batchInterval)
	instance.BatchEventQueue.FlushCallBack = func(err error, batch []map[string]interface{}) {
		assertOutput.Equal(len(batch), batchSize)
	}

	userID, campaignKey, goalIdentifier := testdata.GetRandomUser(), "AB_T_100_W_33_33_33", "GOAL_2"

	for i := 0; i < batchSize; i++ {
		instance.Track(campaignKey, userID, goalIdentifier, nil)
	}
	time.Sleep(time.Duration(1) * time.Second)
	assertOutput.Nil(instance.BatchEventQueue.GetBatchImpressions())
}

func TestFlushQueueOnTimerExpire(t *testing.T) {
	assertOutput := assert.New(t)
	batchSize, batchInterval := 10, 2
	instance := GetVWOInstance(batchSize, batchInterval)
	instance.BatchEventQueue.FlushCallBack = func(err error, batch []map[string]interface{}) {
		assertOutput.Equal(len(batch), batchSize-1)
	}

	userID := testdata.GetRandomUser()

	for i := 0; i < batchSize-1; i++ {
		instance.Push(testdata.ValidTagKey, testdata.ValidTagValue, userID)
	}
	assertOutput.Equal(len(instance.BatchEventQueue.GetBatchImpressions()), batchSize-1)
	time.Sleep(time.Duration(batchInterval+1) * time.Second)
	assertOutput.Nil(instance.BatchEventQueue.GetBatchImpressions())
}
