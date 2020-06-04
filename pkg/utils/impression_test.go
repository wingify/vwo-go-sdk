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

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

func TestCreateImpressionForPush(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	URL := "https://dev.visualwebsiteoptimizer.com/server-side/push"

	userID := testdata.GetRandomUser()
	tagValue := ""
	tagKey := ""
	DemoImpression := CreateImpressionForPush(vwoInstance, tagKey, tagValue, userID)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")

	userID = testdata.GetRandomUser()
	tagValue = testdata.TestValue1
	tagKey = testdata.TestKey1
	DemoImpression = CreateImpressionForPush(vwoInstance, tagKey, tagValue, userID)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	assert.Equal(t, `{"u":{"testKey":"testVal"}}`, DemoImpression.Tags, "Non Matching Parameters")

	userID = testdata.GetRandomUser()
	tagValue = testdata.TestValue2
	tagKey = testdata.TestKey2
	DemoImpression = CreateImpressionForPush(vwoInstance, tagKey, tagValue, userID)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	assert.Equal(t, `{"u":{"test+Key":"test+Val"}}`, DemoImpression.Tags, "Non Matching Parameters")
}

func TestCreateImpressionTrackingUser(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	variationID := 1
	campaignID := 283
	userID := testdata.GetRandomUser()
	URL := "https://dev.visualwebsiteoptimizer.com/server-side/track-user"

	DemoImpression := CreateImpressionTrackingUser(vwoInstance, campaignID, variationID, userID)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, variationID, DemoImpression.Combination, "Non Matching VariationIDs")
	assert.Equal(t, campaignID, DemoImpression.ExperimentID, "Non Matching CampaignIDs")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")

}

func TestCreateImpressionTrackingGoal(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	variationID := 1
	campaignID := 283
	goalID := 281
	goalType := constants.GoalTypeRevenue
	revenueValueActual1 := "1000"
	revenueValueExpected := "1000"
	userID := testdata.GetRandomUser()
	URL := "https://dev.visualwebsiteoptimizer.com/server-side/track-goal"

	DemoImpression := CreateImpressionTrackingGoal(vwoInstance, variationID, userID, goalType, campaignID, goalID, revenueValueActual1)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, variationID, DemoImpression.Combination, "Non Matching VariationIDs")
	assert.Equal(t, goalID, DemoImpression.GoalID, "Non Matching GoalIDs")
	assert.Equal(t, campaignID, DemoImpression.ExperimentID, "Non Matching CampaignIDs")
	assert.Equal(t, revenueValueExpected, DemoImpression.R, "Non Matching Revenues")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	var temp1 string
	assert.IsType(t, temp1, DemoImpression.SID, "Incorrect SID type")
	var temp2 float32
	assert.IsType(t, temp2, DemoImpression.Random, "Incorrect Random type")

	variationID = 2
	campaignID = 284
	goalID = 283
	goalType = constants.GoalTypeRevenue
	revenueValueActual2 := 1234
	revenueValueExpected = "1234"
	userID = testdata.GetRandomUser()
	URL = "https://dev.visualwebsiteoptimizer.com/server-side/track-goal"

	DemoImpression = CreateImpressionTrackingGoal(vwoInstance, variationID, userID, goalType, campaignID, goalID, revenueValueActual2)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, variationID, DemoImpression.Combination, "Non Matching VariationIDs")
	assert.Equal(t, goalID, DemoImpression.GoalID, "Non Matching GoalIDs")
	assert.Equal(t, campaignID, DemoImpression.ExperimentID, "Non Matching CampaignIDs")
	assert.Equal(t, revenueValueExpected, DemoImpression.R, "Non Matching Revenues")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	assert.IsType(t, temp1, DemoImpression.SID, "Incorrect SID type")
	assert.IsType(t, temp2, DemoImpression.Random, "Incorrect Random type")

	variationID = 3
	campaignID = 284
	goalID = 283
	goalType = constants.GoalTypeRevenue
	revenueValueActual3 := 123.234
	revenueValueExpected = "123.234"
	userID = testdata.GetRandomUser()
	URL = "https://dev.visualwebsiteoptimizer.com/server-side/track-goal"

	DemoImpression = CreateImpressionTrackingGoal(vwoInstance, variationID, userID, goalType, campaignID, goalID, revenueValueActual3)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, variationID, DemoImpression.Combination, "Non Matching VariationIDs")
	assert.Equal(t, goalID, DemoImpression.GoalID, "Non Matching GoalIDs")
	assert.Equal(t, campaignID, DemoImpression.ExperimentID, "Non Matching CampaignIDs")
	assert.Equal(t, revenueValueExpected, DemoImpression.R, "Non Matching Revenues")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	assert.IsType(t, temp1, DemoImpression.SID, "Incorrect SID type")
	assert.IsType(t, temp2, DemoImpression.Random, "Incorrect Random type")

	variationID = 3
	campaignID = 284
	goalID = 283
	goalType = constants.GoalTypeRevenue
	revenueValueActual4 := float32(123.234)
	revenueValueExpected = "123.234"
	userID = testdata.GetRandomUser()
	URL = "https://dev.visualwebsiteoptimizer.com/server-side/track-goal"

	DemoImpression = CreateImpressionTrackingGoal(vwoInstance, variationID, userID, goalType, campaignID, goalID, revenueValueActual4)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, variationID, DemoImpression.Combination, "Non Matching VariationIDs")
	assert.Equal(t, goalID, DemoImpression.GoalID, "Non Matching GoalIDs")
	assert.Equal(t, campaignID, DemoImpression.ExperimentID, "Non Matching CampaignIDs")
	assert.Equal(t, revenueValueExpected, DemoImpression.R, "Non Matching Revenues")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	assert.IsType(t, temp1, DemoImpression.SID, "Incorrect SID type")
	assert.IsType(t, temp2, DemoImpression.Random, "Incorrect Random type")

	variationID = 3
	campaignID = 284
	goalID = 283
	goalType = constants.GoalTypeCustom
	revenueValueActual5 := "1000"
	userID = testdata.GetRandomUser()
	URL = "https://dev.visualwebsiteoptimizer.com/server-side/track-goal"

	DemoImpression = CreateImpressionTrackingGoal(vwoInstance, variationID, userID, goalType, campaignID, goalID, revenueValueActual5)
	assert.Equal(t, userID, DemoImpression.UID, "Non Matching UIDs")
	assert.Equal(t, variationID, DemoImpression.Combination, "Non Matching VariationIDs")
	assert.Equal(t, goalID, DemoImpression.GoalID, "Non Matching GoalIDs")
	assert.Equal(t, campaignID, DemoImpression.ExperimentID, "Non Matching CampaignIDs")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, DemoImpression.AccountID, "Non Matching Account IDs")
	assert.Equal(t, URL, DemoImpression.URL, "Non Matching URLs")
	assert.IsType(t, temp1, DemoImpression.SID, "Incorrect SID type")
	assert.IsType(t, temp2, DemoImpression.Random, "Incorrect Random type")
}

func TestGetCommonProperties(t *testing.T) {
	userID := testdata.GetRandomUser()
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")
	impression := getCommonProperties(vwoInstance, userID)

	assert.Equal(t, userID, impression.UID, "Non Matching UIDs")
	assert.Equal(t, vwoInstance.SettingsFile.AccountID, impression.AccountID, "Non Matching AccountIDs")
	assert.Equal(t, constants.Platform, impression.Ap, "Non Matching Platforms")
	assert.Equal(t, constants.SDKName, impression.Sdk, "Non Matching SDK")
	assert.Equal(t, constants.SDKVersion, impression.SdkV, "Non Matching SDK Version")
	var temp1 float32
	assert.IsType(t, temp1, impression.Random, "Random number should be of type float32")
	var temp2 string
	assert.IsType(t, temp2, impression.SID, "Random number should be of type float32")

}
