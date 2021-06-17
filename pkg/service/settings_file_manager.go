/*
 * Copyright 2020-2021 Wingify Software Pvt. Ltd.
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

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

// SettingsFileManager struct to implement SettingsFileM
type SettingsFileManager struct {
	SettingsFile schema.SettingsFile
}

// FetchSettingsFile function makes call to VWO server to fetch the settings file
func (sfm *SettingsFileManager) FetchSettingsFile(accountID, SDKKey string, isViaWebHook bool) error {
	/*
			Args:
				accountID: Config account ID
				SDKKey: Config SDK Key
	      isViaWebHook: specifies if the fetch operation is triggered by a webhook
			Returns:
				error: nil if the settings file id fetched else the error
	*/

	if accountID == "" {
		return fmt.Errorf(constants.ErrorMessageInvalidAccountID, "")
	}
	if SDKKey == "" {
		return fmt.Errorf(constants.ErrorMessageInvalidSDKKey, "")
	}

	endpoint := constants.AccountSettings
	if isViaWebHook {
		endpoint = constants.WebHookAccountSettings
	}

	protocol := constants.HTTPSProtocol
	hostname := constants.BaseURL
	path := endpoint + "?" +
		"a=" + accountID +
		"&i=" + SDKKey +
		"&r=" + strconv.FormatFloat(float64(rand.Float32()), 'f', -1, 64) +
		"&platform=server&" +
		"&sdk=" + constants.SDKName +
		"&sdk-v=" + constants.SDKVersion +
		"&api-version=1"

	resp, err := utils.GetRequest(protocol + hostname + path)
	if err != nil {
		return fmt.Errorf(constants.ErrorMessageSettingsFileCorrupted, "", err.Error())
	}
  resp = utils.JsonCleanUp(resp)
	if err = json.Unmarshal([]byte(resp), &sfm.SettingsFile); err != nil {
		return fmt.Errorf(constants.ErrorMessageInvalidSettingsFile, "", err.Error())
	}

	logger.Warningf(constants.DebugMessageValidConfiguration, "")
	return nil
}

// ProcessSettingsFile Processes the settings_file, assigns variation allocation range
func (sfm *SettingsFileManager) ProcessSettingsFile(settingsFileLocation string) error {
	/*
		Args:
			settingsFileLocation: Location of the settings file on system

		Returns:
			error: nil if the settings file id fetched else the error
	*/

	settingsFile, err := ioutil.ReadFile(settingsFileLocation)
	if err != nil {
		return fmt.Errorf(constants.ErrorMessageCannotReadSettingsFile, "", err.Error())
	}

	if err = json.Unmarshal(settingsFile, &sfm.SettingsFile); err != nil {
		return fmt.Errorf(constants.ErrorMessageInvalidSettingsFile, "", err.Error())
	}

	return nil
}

//Process function processes campaigns in the settings file and sets the variation allocation ranges to all variations
func (sfm *SettingsFileManager) Process() {
	logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	defer logger.Close()
	for i, campaign := range sfm.SettingsFile.Campaigns {
		var (
			currentAllocation         = 0
			variationAllocationRanges []schema.Variation
		)
		for _, variation := range campaign.Variations {
			stepFactor := utils.GetVariationBucketingRange(variation.Weight)
			if stepFactor != 0 {
				variation.StartVariationAllocation = currentAllocation + 1
				variation.EndVariationAllocation = currentAllocation + stepFactor
				currentAllocation += stepFactor
			} else {
				variation.StartVariationAllocation = -1
				variation.EndVariationAllocation = -1
			}
			logs.Infof(constants.InfoMessageVariationRangeAllocation, "", variation.Name, variation.Weight, variation.StartVariationAllocation, variation.EndVariationAllocation)
			variationAllocationRanges = append(variationAllocationRanges, variation)
		}
		sfm.SettingsFile.Campaigns[i].Variations = variationAllocationRanges
	}
}

// GetSettingsFile returns the settings file
func (sfm *SettingsFileManager) GetSettingsFile() schema.SettingsFile {
	return sfm.SettingsFile
}
