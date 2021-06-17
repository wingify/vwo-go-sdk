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

package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	User      string `json:"user"`
	Variation string `json:"variation"`
}

func TestActivate(t *testing.T) {
	assertOutput := assert.New(t)

	var userExpectation map[string][]TestCase
	data, err := ioutil.ReadFile("../testdata/user_expectations1.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &userExpectation); err != nil {
		logger.Info("Error: " + err.Error())
	}

	var settingsFiles map[string]schema.SettingsFile
	data, err = ioutil.ReadFile("../testdata/settings.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &settingsFiles); err != nil {
		logger.Info("Error: " + err.Error())
	}

	logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	defer logger.Close()

	instance := VWOInstance{}
	instance.Logger = logs

	for settingsFileName, settingsFile := range settingsFiles {
		vwoInstance := schema.VwoInstance{
			Logger: logs,
		}
		settingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, settingsFile.Campaigns[0].Variations)

		instance.SettingsFile = settingsFile

		testCases := userExpectation[settingsFileName]
		for i := range testCases {
			actual := instance.Activate(settingsFile.Campaigns[0].Key, testCases[i].User, nil)
			expected := testCases[i].Variation
			assertOutput.Equal(expected, actual, settingsFileName+" "+testCases[i].User)
		}
	}

	// CORNER CASES

	var customSettingsFiles map[string]schema.SettingsFile
	data, err = ioutil.ReadFile("../testdata/custom_settings.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &customSettingsFiles); err != nil {
		logger.Info("Error: " + err.Error())
	}

	settings := customSettingsFiles["SettingsFile2"]
	instance.SettingsFile = settings

	userID := ""
	campaignKey := ""
	value := instance.Activate(campaignKey, userID, nil)
	assertOutput.Empty(value, "Invalid params")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NonExistingCampaign
	value = instance.Activate(campaignKey, userID, nil)
	assertOutput.Empty(value, "Campaign does not exist")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NotRunningCampaign
	value = instance.Activate(campaignKey, userID, nil)
	assertOutput.Empty(value, "Campaign Not running")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.FeatureRolloutCampaign
	value = instance.Activate(campaignKey, userID, nil)
	assertOutput.Empty(value, "Campaign Not Valid")
}
