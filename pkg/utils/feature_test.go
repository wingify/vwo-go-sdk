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

package utils

import (
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGetVariableValueForVariation(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("FT_T_0_W_10_20_30_40")
	campaign := vwoInstance.SettingsFile.Campaigns[0]
	userID := ""

	variation := campaign.Variations[0]
	variableKey := testdata.InvalidVariableKey
	variable := GetVariableValueForVariation(vwoInstance, campaign, variation, variableKey, userID)
	assert.Empty(t, variable, "Expected object should be empty")

	variation = campaign.Variations[0]
	variableKey = testdata.ValidVariableKey2
	variable = GetVariableValueForVariation(vwoInstance, campaign, variation, variableKey, userID)
	assert.Equal(t, campaign.Variations[0].Variables[0], variable, "Expected and Actual IDs should be same")

	variation = campaign.Variations[1]
	variableKey = testdata.ValidVariableKey1
	variable = GetVariableValueForVariation(vwoInstance, campaign, variation, variableKey, userID)
	assert.Equal(t, campaign.Variations[1].Variables[1], variable, "Expected and Actual IDs should be same")

}

func TestGetVariableForFeature(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("FR_T_0_W_100")

	variables := vwoInstance.SettingsFile.Campaigns[0].Variables
	variableKey := testdata.ValidVariableKey2
	variable := GetVariableForFeature(variables, variableKey)
	assert.Equal(t, variables[0], variable, "Expected and Actual IDs should be same")

	variableKey = testdata.InvalidVariableKey
	variable = GetVariableForFeature(variables, variableKey)
	assert.Empty(t, variable, "Expected variable should be empty")
}
