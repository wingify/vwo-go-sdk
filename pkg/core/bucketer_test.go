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

package core

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

// BucketTestCase struct
type BucketTestCase struct {
	User        string `json:"user"`
	BucketValue int    `json:"bucket_value"`
}

func TestBucketUserToVariation(t *testing.T) {
	assertOutput := assert.New(t)
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_100_W_50_50")
	vwoInstance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, vwoInstance.SettingsFile.Campaigns[0].Variations)

	campaign := vwoInstance.SettingsFile.Campaigns[0]
	userID := testdata.GetRandomUser()
	actual, err := BucketUserToVariation(vwoInstance, userID, campaign)
	assertOutput.Nil(err, "Variations did not match")
	assertOutput.NotEmpty(actual, "Variations did not match")

	vwoInstance = testdata.GetInstanceWithCustomSettings("SettingsFile2")
	vwoInstance.SettingsFile.Campaigns[1].Variations = utils.GetVariationAllocationRanges(vwoInstance, vwoInstance.SettingsFile.Campaigns[1].Variations)

	campaign = vwoInstance.SettingsFile.Campaigns[1]
	userID = testdata.GetRandomUser()
	actual, err = BucketUserToVariation(vwoInstance, userID, campaign)
	assertOutput.NotNil(err, "Variation expected to be empty")
	assertOutput.Empty(actual, "Variation expected to be empty")
}

func TestGetBucketerVariation(t *testing.T) {
	assertOutput := assert.New(t)
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_100_W_33_33_33")
	vwoInstance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, vwoInstance.SettingsFile.Campaigns[0].Variations)

	variations := vwoInstance.SettingsFile.Campaigns[0].Variations
	bucketValue := testdata.ValidBucketValue
	actual, err := GetBucketerVariation(vwoInstance, variations, bucketValue, "", "")
	expected := variations[0]
	assertOutput.Nil(err, "Expected Variation do not match with Actual")
	assertOutput.Equal(expected, actual, "Expected Variation do not match with Actual")

	bucketValue = testdata.InvalidBucketValue
	actual, err = GetBucketerVariation(vwoInstance, variations, bucketValue, "", "")
	assertOutput.NotNil(err, "Variation should be empty")
	assertOutput.Empty(actual, "Variation should be empty")
}

func TestIsUserPartOfCampaign(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_100_W_33_33_33")
	vwoInstance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, vwoInstance.SettingsFile.Campaigns[0].Variations)

	userID := testdata.ValidUser
	campaign := vwoInstance.SettingsFile.Campaigns[0]
	actual := IsUserPartOfCampaign(vwoInstance, userID, campaign)
	assert.True(t, actual, "User should be part of the campaign")

	vwoInstance = testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	vwoInstance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, vwoInstance.SettingsFile.Campaigns[0].Variations)

	userID = testdata.InvalidUser
	campaign = vwoInstance.SettingsFile.Campaigns[0]
	actual = IsUserPartOfCampaign(vwoInstance, userID, campaign)
	assert.False(t, actual, "User should not be part of the campaign")

	vwoInstance = testdata.GetInstanceWithCustomSettings("SettingsFile1")
	userID = testdata.ValidUser
	campaign = vwoInstance.SettingsFile.Campaigns[0]
	actual = IsUserPartOfCampaign(vwoInstance, userID, campaign)
	assert.False(t, actual, "User should not be part of the campaign")
}

func TestGetBucketValueForUser(t *testing.T) {
	var settings map[string][]BucketTestCase
	data, err := ioutil.ReadFile("../testdata/bucket_value_expectations.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &settings); err != nil {
		logger.Info("Error: " + err.Error())
	}

	TestCases := settings["USER_AND_BUCKET_VALUES"]

	logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	defer logger.Close()

	vwoInstance := schema.VwoInstance{
		Logger: logs,
	}

	for _, testCase := range TestCases {
		expected := testCase.BucketValue
		_, actual := GetBucketValueForUser(vwoInstance, testCase.User, 10000, 1, vwoInstance.Campaign)
		assert.Equal(t, expected, actual, "Failed for: "+testCase.User)
	}
}

func GetBucketValueForUser1(t *testing.T) {
	var campaign schema.Campaign
	var vwoInstance schema.VwoInstance
	campaign.ID = 1
	campaign.IsBucketingSeedEnabled = true
	userID := "someone@mail.com"
	ExpectedBucketVal := 2444
	_, CalculatedVal := GetBucketValueForUser(vwoInstance, userID, 10000, 1, campaign)
	assert.Equal(t, ExpectedBucketVal, CalculatedVal, "Failed when userID is "+userID+" and bucketing seed is true")

	campaign.IsBucketingSeedEnabled = false
	ExpectedBucketVal = 6361
	_, CalculatedVal = GetBucketValueForUser(vwoInstance, userID, 10000, 1, campaign)
	assert.Equal(t, ExpectedBucketVal, CalculatedVal, "Failed when userID is "+userID+" and bucketing seed is false")
}

func GetBucketValueForUser1111111111111111(t *testing.T) {
	var campaign schema.Campaign
	var vwoInstance schema.VwoInstance
	campaign.ID = 1
	campaign.IsBucketingSeedEnabled = true
	userID := "1111111111111111"
	ExpectedBucketVal := 8177
	_, CalculatedVal := GetBucketValueForUser(vwoInstance, userID, 10000, 1, campaign)
	assert.Equal(t, ExpectedBucketVal, CalculatedVal, "Failed when userID is "+userID+" and bucketing seed is true")

	campaign.IsBucketingSeedEnabled = false
	ExpectedBucketVal = 4987
	_, CalculatedVal = GetBucketValueForUser(vwoInstance, userID, 10000, 1, campaign)
	assert.Equal(t, ExpectedBucketVal, CalculatedVal, "Failed when userID is "+userID+" and bucketing seed is false")
}

func TestHash(t *testing.T) {
	actual := hash(testdata.GetRandomUser())
	assert.NotNil(t, actual, "Hash values do not match")
}
