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

package utils

import (
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

func TestCheckCampaignType(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")

	campaign := vwoInstance.SettingsFile.Campaigns[0]
	campaignType := constants.CampaignTypeVisualAB
	value := CheckCampaignType(campaign, campaignType)
	assert.True(t, value, "Campaign should match")

	vwoInstance = testdata.GetInstanceWithSettings("FT_T_0_W_10_20_30_40")

	campaign = vwoInstance.SettingsFile.Campaigns[0]
	campaignType = constants.CampaignTypeFeatureTest
	value = CheckCampaignType(campaign, campaignType)
	assert.True(t, value, "Campaign should not match")

	vwoInstance = testdata.GetInstanceWithSettings("FR_T_0_W_100")

	campaign = vwoInstance.SettingsFile.Campaigns[0]
	campaignType = constants.CampaignTypeFeatureRollout
	value = CheckCampaignType(campaign, campaignType)
	assert.True(t, value, "Campaign should not match")

	campaign = vwoInstance.SettingsFile.Campaigns[0]
	campaignType = constants.CampaignTypeFeatureTest
	value = CheckCampaignType(campaign, campaignType)
	assert.False(t, value, "Campaign should not match")
}

func TestGetKeyValue(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("T_50_W_50_50_WS")

	segment := vwoInstance.SettingsFile.Campaigns[0].Segments
	actualKey, actualValue := GetKeyValue(segment)
	expectedKey := constants.OperatorTypeAnd
	assert.Equal(t, expectedKey, actualKey, "Expected and Actual Keys should be same")
	var Temp []interface{}
	assert.IsType(t, Temp, actualValue, "Type Mismatch")

	var tempSegment map[string]interface{}
	actualKey, actualValue = GetKeyValue(tempSegment)
	assert.Equal(t, "", actualKey, "Nil Value expected")
	assert.Nil(t, actualValue, "Nil Value expected")
}
