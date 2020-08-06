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

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/core"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsFeatureEnabled(t *testing.T) {
	assertOutput := assert.New(t)

	var settingsFiles map[string]schema.SettingsFile
	data, err := ioutil.ReadFile("../testdata/settings.json")
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

	for settingsFileName, settingsFile := range settingsFiles {
		vwoInstance := schema.VwoInstance{
			Logger: logs,
		}
		settingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, settingsFile.Campaigns[0].Variations)

		instance.SettingsFile = settingsFile
		if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeVisualAB {
			actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, testdata.GetRandomUser(), nil)
			assertOutput.False(actual, "Wrong Campaign Type")
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureRollout && settingsFileName != "NEW_SETTINGS_FILE" {
			userID := testdata.GetRandomUser()
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, userID, nil)

				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.Nil(err, "Error encuntered")
				assertOutput.True(actual, "Feature Rollout Campaign")
			}
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureRollout && settingsFileName == "NEW_SETTINGS_FILE" {
			userID := testdata.GetRandomUser()
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, userID, nil)

				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.Nil(err, "Error encuntered")
				assertOutput.True(actual, "Feature Rollout Campaign")
			}
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[1], "", schema.Options{}); variation.Name != "" {
				if variation.Name == "Control" {
					actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[1].Key, userID, nil)

					assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
					assertOutput.NotNil(err, "Error encuntered")
					assertOutput.False(actual, "Feature Test Campaign")
				} else {
					actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[1].Key, userID, nil)

					assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
					assertOutput.Nil(err, "Error encuntered")
					assertOutput.True(actual, "Feature Test Campaign")
				}
			}
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[2], "", schema.Options{}); variation.Name != "" {
				actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[2].Key, userID, nil)

				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.NotNil(err, "Error encuntered")
				assertOutput.False(actual, "Visual AB Campaign")
			}
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureTest && settingsFileName != "FT_T_100_W_10_20_30_40_IFEF" {
			userID := testdata.GetRandomUser()
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				if variation.Name == "Control" {
					actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, userID, nil)
					
					assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
					assertOutput.Nil(err, "Error encuntered")
					assertOutput.False(actual, "Feature Test Campaign : " + variation.Name + userID)
				} else {
					actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, userID, nil)
					assertOutput.True(actual, "Feature Test Campaign : " + variation.Name + userID)
				}
			}
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureTest && settingsFileName == "FT_T_100_W_10_20_30_40_IFEF" {
			userID := testdata.GetRandomUser()
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				if variation.Name == "Variation-2" {
					actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, userID, nil)

					assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
					assertOutput.Nil(err, "Error encuntered")
					assertOutput.True(actual, "Feature Test Campaign : " + variation.Name + userID)
				} else {
					actual := instance.IsFeatureEnabled(instance.SettingsFile.Campaigns[0].Key, userID, nil)

					assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
					assertOutput.Nil(err, "Error encuntered")
					assertOutput.False(actual, "Feature Test Campaign : " + variation.Name + userID)
				}
			}
		}
	}

	// CORNER CASES

	instance.SettingsFile = settingsFiles["FT_T_100_W_10_20_30_40"]
	userID := testdata.UserIsFeatureEnabled
	campaignKey := instance.SettingsFile.Campaigns[0].Key
	value := instance.IsFeatureEnabled(campaignKey, userID, nil)
	assertOutput.False(value, "Control Variation")

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

	userID = ""
	campaignKey = ""
	value = instance.IsFeatureEnabled(campaignKey, userID, nil)
	assertOutput.False(value, "Invalid params")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NonExistingCampaign
	value = instance.IsFeatureEnabled(campaignKey, userID, nil)
	assertOutput.False(value, "Campaign does not exist")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NotRunningCampaign
	value = instance.IsFeatureEnabled(campaignKey, userID, nil)
	assertOutput.False(value, "Campaign Not Running")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.FeatureRolloutCampaign
	value = instance.IsFeatureEnabled(campaignKey, userID, nil)
	assertOutput.False(value, "No Variation from campaign Not alloted")
}
