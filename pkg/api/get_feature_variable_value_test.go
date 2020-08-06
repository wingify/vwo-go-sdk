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
	"strconv"
	"github.com/wingify/vwo-go-sdk/pkg/core"
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

func TestGetFeatureVariableValue(t *testing.T) {
	assertOutput := assert.New(t)

	var userExpectation map[string]map[string]interface{}
	data, err := ioutil.ReadFile("../testdata/user_expectations2.json")
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

	for settingsFileName, settingsFile := range settingsFiles {
		vwoInstance := schema.VwoInstance{
			Logger: logs,
		}
		settingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, settingsFile.Campaigns[0].Variations)

		instance.SettingsFile = settingsFile

		if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureRollout && settingsFileName != "FR_WRONG_VARIABLE_TYPE" && settingsFileName != "NEW_SETTINGS_FILE" {
			userID := testdata.GetFeatureDummyUser
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				testCases := userExpectation["ROLLOUT_VARIABLES"]

				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.Nil(err, "Error encuntered")

				actual := instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_VARIABLE", userID, nil)
				assertOutput.Equal(testCases["STRING_VARIABLE"], actual.(string), settingsFileName + " : STRING_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_VARIABLE", userID, nil)
				assertOutput.Equal(testCases["INTEGER_VARIABLE"], actual.(float64), settingsFileName + " : INTEGER_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "FLOAT_VARIABLE", userID, nil)
				assertOutput.Equal(testCases["FLOAT_VARIABLE"], actual.(float64), settingsFileName + " : FLOAT_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "BOOLEAN_VARIABLE", userID, nil)
				assertOutput.Equal(testCases["BOOLEAN_VARIABLE"], actual.(bool), settingsFileName + " : BOOLEAN_VARIABLE")
			}
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureRollout && settingsFileName == "FR_WRONG_VARIABLE_TYPE" && settingsFileName != "NEW_SETTINGS_FILE" {
			userID := testdata.GetFeatureDummyUser
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.Nil(err, "Error encuntered")
				
				actual := instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_VARIABLE", userID, nil)
				assertOutput.Nil(actual, settingsFileName + " : STRING_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_VARIABLE", userID, nil)
				assertOutput.Nil(actual, settingsFileName + " : INTEGER_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "FLOAT_VARIABLE", userID, nil)
				assertOutput.Nil(actual, settingsFileName + " : FLOAT_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "BOOLEAN_VARIABLE", userID, nil)
				assertOutput.Nil(actual, settingsFileName + " : BOOLEAN_VARIABLE")

				value := instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_TO_INTEGER", userID, nil)
				actual, err := strconv.ParseInt(value.(string), 10, 64)
				assertOutput.Nil(err, "Error")
				assertOutput.Equal(int64(123), actual, settingsFileName + " : STRING_TO_INTEGER")

				value = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_TO_FLOAT", userID, nil)
				actual, err = strconv.ParseFloat(value.(string), 64)
				assertOutput.Equal(float64(123.456), actual, settingsFileName + " : STRING_TO_FLOAT")

				value = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "BOOLEAN_TO_STRING", userID, nil)
				actual = strconv.FormatBool(value.(bool))
				assertOutput.Equal("true", actual, settingsFileName + " : BOOLEAN_TO_STRING")

				// value = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_TO_STRING", userID, nil)
				// actual = strconv.FormatFloat(value.(float64), 'E', -1, 64)
				// assertOutput.Equal("24", actual, settingsFileName + " : INTEGER_TO_STRING")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_TO_FLOAT", userID, nil)
				assertOutput.Equal(float64(24), actual.(float64), settingsFileName + " : INTEGER_TO_FLOAT")

				// value = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "FLOAT_TO_STRING", userID, nil)
				// value = strconv.FormatFloat(value.(float64), 'E', -1, 64)
				// assertOutput.Equal("24.24", actual, settingsFileName + " : FLOAT_TO_STRING")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "FLOAT_TO_INTEGER", userID, nil)
				assertOutput.Equal(float64(24), actual.(float64), settingsFileName + " : FLOAT_TO_INTEGER")

				value = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "WRONG_BOOLEAN", userID, nil)
				actual, err = strconv.ParseBool(value.(string))
				assertOutput.Nil(err, "Error : ")
				assertOutput.Equal(true, actual, settingsFileName + " : WRONG_BOOLEAN")
			}
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureRollout && settingsFileName != "FR_WRONG_VARIABLE_TYPE" && settingsFileName == "NEW_SETTINGS_FILE" {
			userID := testdata.GetFeatureDummyUser
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				testCases := userExpectation["ROLLOUT_VARIABLES"]

				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.Nil(err, "Error encuntered")

				actual := instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_VARIABLE", userID, nil)
				assertOutput.Equal("d1", actual.(string), settingsFileName + " : STRING_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_VARIABLE", userID, nil)
				assertOutput.Equal(testCases["INTEGER_VARIABLE"], actual.(float64), settingsFileName + " : INTEGER_VARIABLE")
			}
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeVisualAB {
			userID := testdata.GetFeatureDummyUser
			actual := instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_VARIABLE", userID, nil)
			assertOutput.Nil(actual, "Wrong Campaign Type")

			actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_VARIABLE", userID, nil)
			assertOutput.Nil(actual, "Wrong Campaign Type")

			actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "FLOAT_VARIABLE", userID, nil)
			assertOutput.Nil(actual, "Wrong Campaign Type")

			actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "BOOLEAN_VARIABLE", userID, nil)
			assertOutput.Nil(actual, "Wrong Campaign Type")
		} else if instance.SettingsFile.Campaigns[0].Type == constants.CampaignTypeFeatureTest {
				userID := testdata.GetFeatureDummyUser
			if variation, storedGoalIdentifiers, err := core.GetVariation(vwoInstance, userID, instance.SettingsFile.Campaigns[0], "", schema.Options{}); variation.Name != "" {
				assertOutput.Empty(storedGoalIdentifiers, "Incorrect Assertion for storedGoalIdentifiers ")
				assertOutput.Nil(err, "Error encuntered")
				
				actual := instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "STRING_VARIABLE", userID, nil)
				expected := userExpectation["STRING_VARIABLE"][variation.Name]
				assertOutput.Equal(expected, actual.(string), settingsFileName + " : STRING_VARIABLE")

				actual = instance.GetFeatureVariableValue(instance.SettingsFile.Campaigns[0].Key, "INTEGER_VARIABLE", userID, nil)
				expected = userExpectation["INTEGER_VARIABLE"][variation.Name]
				assertOutput.Equal(expected, actual.(float64), settingsFileName + " : INTEGER_VARIABLE")
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
	variableKey := ""
	value := instance.GetFeatureVariableValue(campaignKey, variableKey, userID, nil)
	assertOutput.Nil(value, "Invalid params")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NonExistingCampaign
	variableKey = testdata.GetFeatureDummyVariable
	value = instance.GetFeatureVariableValue(campaignKey, variableKey, userID, nil)
	assertOutput.Nil(value, "Campaign does not exist")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.NotRunningCampaign
	variableKey = testdata.GetFeatureDummyVariable
	value = instance.GetFeatureVariableValue(campaignKey, variableKey, userID, nil)
	assertOutput.Nil(value, "Campaign Not running")

	userID = testdata.GetRandomUser()
	campaignKey = testdata.FeatureRolloutCampaign
	variableKey = testdata.GetFeatureDummyVariable
	value = instance.GetFeatureVariableValue(campaignKey, variableKey, userID, nil)
	assertOutput.Nil(value, "Variation Not alloted as none exist")
}
