/*
 * Copyright 2020 Wingify Software Pvt. Ltd.
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

package api

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/service"
	"github.com/stretchr/testify/assert"
)

// GetInstance function creates and return a temporary VWO instance for testing
func getInstance(path string) (*VWOInstance, error) {
	settingsFileManager := service.SettingsFileManager{}
	if err := settingsFileManager.ProcessSettingsFile(path); err != nil {
		log.Println("Error Processing Settings File: ", err)
	}
	settingsFileManager.Process()
	settingsFile := settingsFileManager.GetSettingsFile()

	var vwo VWOInstance
	vwo.SettingsFile = settingsFile

	return vwo.Init(WithDevelopmentMode())
}

type WUserStorage interface {
	Getter(userID, campaignKey string) schema.UserData
	Setter(userID, campaignKey, variationName, goalIdentifier string)
}

type WUserStorageData struct{}

func (us *WUserStorageData) Getter(userID, campaignKey string) schema.UserData {
	return schema.UserData{}
}
func (us *WUserStorageData) Setter(userID, campaignKey, variationName, goalIdentifier string) {}

func TestInit(t *testing.T) {
	vwoInstance := VWOInstance{}
	_, err := vwoInstance.Init(WithDevelopmentMode())
	assert.Nil(t, err)

	vwoInstance = VWOInstance{}
	storage := &WUserStorageData{}
	_, err = vwoInstance.Init(WithStorage(storage), WithGoalAttributes(nil, false))
	assert.NotNil(t, err)

	logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	defer logger.Close()

	vwoInstance = VWOInstance{}
	_, err = vwoInstance.Init(WithLogger(logs))
	assert.Nil(t, err)
}
