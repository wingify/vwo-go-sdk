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
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestTrack(t *testing.T) {
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
	instance.SettingsFile = schema.SettingsFile{}
	instance.Logger = logs

	options := make(map[string]interface{})
	options["revenueValue"] = 12

	for settingsFileName, settingsFile := range settingsFiles {
		vwoInstance := schema.VwoInstance{
			Logger: logs,
		}
		settingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, settingsFile.Campaigns[0].Variations)

		instance.SettingsFile = settingsFile

		if instance.SettingsFile.Campaigns[0].Type != constants.CampaignTypeFeatureRollout {
			testCases := userExpectation[settingsFileName]
			for i := range testCases {
				if testCases[i].Variation != "" {
					actual := instance.Track(settingsFile.Campaigns[0].Key, testCases[i].User, settingsFile.Campaigns[0].Goals[0].Identifier, options)
					assertOutput.True(actual, settingsFileName+" "+testCases[i].User)
				} else {
					actual := instance.Track(settingsFile.Campaigns[0].Key, testCases[i].User, settingsFile.Campaigns[0].Goals[0].Identifier, options)
					assertOutput.False(actual, settingsFileName+" "+testCases[i].User)
				}
			}
		} else {
			testCases := userExpectation[settingsFileName]
			for i := range testCases {
				if testCases[i].Variation != "" {
					actual := instance.Track(settingsFile.Campaigns[0].Key, testCases[i].User, settingsFile.Campaigns[0].Goals[0].Identifier, options)
					assertOutput.True(actual, settingsFileName+" "+testCases[i].User)
				} else {
					actual := instance.Track(settingsFile.Campaigns[0].Key, testCases[i].User, settingsFile.Campaigns[0].Goals[0].Identifier, options)
					assertOutput.False(actual, settingsFileName+" "+testCases[i].User)
				}
			}
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
	goalIdentifier := ""
	value := instance.Track(campaignKey, userID, goalIdentifier, nil)
	assertOutput.False(value, "Invalid params")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NonExistingCampaign
	goalIdentifier = "GOAL_0"
	value = instance.Track(campaignKey, userID, goalIdentifier, nil)
	assertOutput.False(value, "Campaign does not exist")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NotRunningCampaign
	goalIdentifier = "GOAL_0"
	value = instance.Track(campaignKey, userID, goalIdentifier, nil)
	assertOutput.False(value, "Campaign Not running")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.FeatureRolloutCampaign
	goalIdentifier = "GOAL_0"
	value = instance.Track(campaignKey, userID, goalIdentifier, nil)
	assertOutput.False(value, "Campaign Not Valid")

	userID = testdata.GetRandomUser()
	campaignKey = "CAMPAIGN_3"
	goalIdentifier = "GOAL_0"
	value = instance.Track(campaignKey, userID, goalIdentifier, nil)
	assertOutput.False(value, "Goal Not Found")

	userID = testdata.GetRandomUser()
	campaignKey = "CAMPAIGN_3"
	goalIdentifier = "abcd"
	value = instance.Track(campaignKey, userID, goalIdentifier, nil)
	assertOutput.False(value, "Revenue Not defined")

	// userID = testdata.GetRandomUser()
	// campaignKey = "CAMPAIGN_2"
	// option = map[string]interface{}{
	// 	"revenueValue": 10,
	// }
	// goalIdentifier = "GOAL_2"
	// value = instance.Track(campaignKey, userID, goalIdentifier, option)
	// assertOutput.False(value, "No Variation in Campaign")
}
