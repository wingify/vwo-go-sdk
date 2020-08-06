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

package utils

import (
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGetVariationBucketingRange(t *testing.T) {
	var weight float64

	weight = 0
	actual := GetVariationBucketingRange(weight)
	expected := 0
	assert.Equal(t, expected, actual, "Expected and Actual Ranges should be same")

	weight = 33.333
	actual = GetVariationBucketingRange(weight)
	expected = 3334
	assert.Equal(t, expected, actual, "Expected and Actual Ranges should be same")

	weight = 102
	actual = GetVariationBucketingRange(weight)
	expected = constants.MaxTrafficValue
	assert.Equal(t, expected, actual, "Expected and Actual Ranges should be same")
}

func TestGetCampaign(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")

	campaignKey := testdata.ValidCampaignKey
	campaign, err := GetCampaign("", vwoInstance.SettingsFile, campaignKey)
	assert.Nil(t, err)
	assert.Equal(t, vwoInstance.SettingsFile.Campaigns[0], campaign, "Expected and Actual Campaign IDs should be same")

	campaignKey = testdata.InvalidCampaignKey
	campaign, err = GetCampaign("", vwoInstance.SettingsFile, campaignKey)
	assert.NotNil(t, err)
	assert.Empty(t, campaign, "Expected campaign should be empty")
}

func TestGetCampaignVariation(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	campaign := vwoInstance.SettingsFile.Campaigns[0]

	variationName := testdata.ValidVariationName
	variation, err := GetCampaignVariation("", campaign, variationName)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Variations[0], variation, "Expected and Actual Variation IDs should be same")

	variationName = testdata.IncorrectNariationName
	variation, err = GetCampaignVariation("", campaign, variationName)
	assert.NotNil(t, err)
	assert.Empty(t, variation, "Expected Variation should be empty")

	vwoInstance = testdata.GetInstanceWithCustomSettings("SettingsFile1")
	campaign = vwoInstance.SettingsFile.Campaigns[0]
	variationName = testdata.InvalidVariationName
	variation, err = GetCampaignVariation("", campaign, variationName)
	assert.NotNil(t, err)
	assert.Empty(t, variation, "No Variations in the Campaign")
}

func TestGetCampaignGoal(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	campaign := vwoInstance.SettingsFile.Campaigns[0]

	goalName := testdata.ValidGoal
	goal, err := GetCampaignGoal("", campaign, goalName)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Goals[0], goal, "Expected and Actual Goal IDs should be same")

	goalName = testdata.InvalidGoal
	goal, err = GetCampaignGoal("", campaign, goalName)
	assert.NotNil(t, err)
	assert.Empty(t, goal, "Expected Goal should be empty")
}

func TestGetControlVariation(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")

	campaign := vwoInstance.SettingsFile.Campaigns[0]
	variation := GetControlVariation(campaign)
	assert.Equal(t, campaign.Variations[0], variation, "Expected variation should be present in the campaign")

	vwoInstance = testdata.GetInstanceWithCustomSettings("SettingsFile1")
	campaign = vwoInstance.SettingsFile.Campaigns[0]
	variation = GetControlVariation(campaign)
	assert.Empty(t, variation, "Expected variation should be empty")
}

func TestScaleVariations(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")

	variations := vwoInstance.SettingsFile.Campaigns[0].Variations
	variations = ScaleVariations(variations)
	assert.Equal(t, vwoInstance.SettingsFile.Campaigns[0].Variations, variations, "List of variations did not match")

	vwoInstance = testdata.GetInstanceWithCustomSettings("SettingsFile1")

	variations = vwoInstance.SettingsFile.Campaigns[1].Variations
	variations = ScaleVariations(variations)
	assert.Equal(t, vwoInstance.SettingsFile.Campaigns[1].Variations, variations, "List of variations did not match")
}

func TestGetVariationAllocationRanges(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")

	variations := vwoInstance.SettingsFile.Campaigns[0].Variations
	assert.NotEmpty(t, variations, "No Variations recieved")
	variations = GetVariationAllocationRanges(vwoInstance, variations)
	assert.Equal(t, 1, variations[0].StartVariationAllocation, "Value Mismatch")
	assert.Equal(t, 10000, variations[1].EndVariationAllocation, "Value Mismatch")

	vwoInstance = testdata.GetInstanceWithSettings("AB_T_100_W_0_100")
	variations = vwoInstance.SettingsFile.Campaigns[0].Variations
	assert.NotEmpty(t, variations, "No Variations recieved")
	variations = GetVariationAllocationRanges(vwoInstance, variations)
	assert.Equal(t, -1, variations[0].StartVariationAllocation, "Start Allocation range failed to match")
	assert.Equal(t, -1, variations[0].EndVariationAllocation, "End Allocation range failed to match")
	assert.Equal(t, 1, variations[1].StartVariationAllocation, "Start Allocation range failed to match")
	assert.Equal(t, 10000, variations[1].EndVariationAllocation, "End Allocation range failed to match")
}

func TestGetCampaignForKeys(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("NEW_SETTINGS_FILE")

	campaignKeys := []string{"FEATURE_ROLLOUT_KEY", "FEATURE_TEST", "NEW_RECOMMENDATION_AB_CAMPAIGN"}
	campaigns, err := GetCampaignForKeys(vwoInstance, campaignKeys)
	assert.Nil(t, err, "Encountered Error")
	assert.Equal(t, vwoInstance.SettingsFile.Campaigns, campaigns, "List of campaigns did not match")

	campaignKeys = []string{"FEATURE_ROLLOUT_KEY", "FEATURE_TEST_NT_EXIST"}
	campaigns, err = GetCampaignForKeys(vwoInstance, campaignKeys)
	assert.Nil(t, err, "Encountered Error")
	var expected []schema.Campaign
	expected = append(expected, vwoInstance.SettingsFile.Campaigns[0])
	assert.Equal(t, expected, campaigns, "List of campaigns did not match")

	campaignKeys = []string{"FEATURE_TEST_NT_EXIST"}
	campaigns, err = GetCampaignForKeys(vwoInstance, campaignKeys)
	assert.NotNil(t, err, "Encountered Error")
	assert.Empty(t, campaigns, "List of campaigns did not match")
}

func TestGetCampaignForGoals(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("NEW_SETTINGS_FILE")

	goalIdentifier := "FEATURE_TEST_GOAL"
	goalTypeToTrack := "CUSTOM_GOAL"
	campaigns, err := GetCampaignForGoals(vwoInstance, goalIdentifier, goalTypeToTrack)
	assert.Nil(t, err, "Encountered Error")
	expected := []schema.Campaign{vwoInstance.SettingsFile.Campaigns[1], vwoInstance.SettingsFile.Campaigns[2]}
	assert.Equal(t, expected, campaigns, "List of campaigns did not match")

	goalIdentifier = "FEATURE_TEST_GOAL"
	goalTypeToTrack = "REVENUE_TRACKING"
	campaigns, err = GetCampaignForGoals(vwoInstance, goalIdentifier, goalTypeToTrack)
	assert.NotNil(t, err, "Encountered Error")
	assert.Empty(t, campaigns, "List of campaigns did not match")
}

func TestMin(t *testing.T) {
	assert.Equal(t, 10, min(10, 20), "Incorrect")
	assert.Equal(t, 10, min(20, 10), "Incorrect")
	assert.NotEqual(t, 12, min(10, 20), "Incorrect")
}

func TestMax(t *testing.T) {
	assert.Equal(t, 20, max(10, 20), "Incorrect")
	assert.Equal(t, 20, max(20, 10), "Incorrect")
	assert.NotEqual(t, 12, max(10, 20), "Incorrect")
}
