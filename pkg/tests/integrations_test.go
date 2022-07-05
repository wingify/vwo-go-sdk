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

package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/wingify/vwo-go-sdk/pkg/api"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
	"testing"
)

func GetVWOClientInstance(settingsFileIdentifier string) api.VWOInstance {
	instance := api.VWOInstance{}
	instance.SettingsFile = schema.SettingsFile{}

	vwoInstance := testdata.GetInstanceWithSettings(settingsFileIdentifier)
	instance.SettingsFile = vwoInstance.SettingsFile
	instance.Logger = vwoInstance.Logger
	instance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, instance.SettingsFile.Campaigns[0].Variations)

	return instance
}

func TestIntegrationsABTest(t *testing.T) {
	assertOutput := assert.New(t)
	instance := GetVWOClientInstance("AB_T_100_W_33_33_33")
	instance.Integrations.CallBack = func(integrationsMap map[string]interface{}) {
		assertOutput.Equal(len(integrationsMap), 16)
		assertOutput.Equal(integrationsMap["fromUserStorageService"], false)
		assertOutput.Equal(integrationsMap["isFeatureEnabled"], nil)
		assertOutput.Equal(integrationsMap["isUserWhitelisted"], false)
		assertOutput.Equal(integrationsMap["campaignKey"], "AB_T_100_W_33_33_33")
	}

	userID, campaignKey := testdata.GetRandomUser(), "AB_T_100_W_33_33_33"

	//getVariation and Activate API when the variation is fetched using the murmur logic.
	instance.Activate(campaignKey, userID, nil)
	instance.GetVariationName("AB_T_100_W_33_33_33", userID, nil)

	//with whitelisting
	instance = GetVWOClientInstance("T_100_W_33_33_33_WS_WW")
	instance.Integrations.CallBack = func(integrationsMap map[string]interface{}) {
		assertOutput.Equal(integrationsMap["fromUserStorageService"], false)
		assertOutput.Equal(integrationsMap["isFeatureEnabled"], nil)
		assertOutput.Equal(integrationsMap["isUserWhitelisted"], true)
		assertOutput.Equal(integrationsMap["variationName"], "Control")
		assertOutput.Equal(integrationsMap["campaignKey"], "T_100_W_33_33_33_WS_WW")
	}
	options := make(map[string]interface{})
	options["variationTargetingVariables"] = map[string]interface{}{"safari": true}
	instance.Activate("T_100_W_33_33_33_WS_WW", campaignKey, options)
	instance.GetVariationName("T_100_W_33_33_33_WS_WW", userID, options)
}

func TestIntegrationFeatureRollout(t *testing.T) {
	assertOutput := assert.New(t)
	instance := GetVWOClientInstance("FR_T_100_W_100")
	userID := testdata.GetRandomUser()
	instance.Integrations.CallBack = func(integrationsMap map[string]interface{}) {
		assertOutput.Equal(len(integrationsMap), 15)
		assertOutput.Equal(integrationsMap["fromUserStorageService"], false)
		assertOutput.Equal(integrationsMap["isFeatureEnabled"], true)
		assertOutput.Equal(integrationsMap["isUserWhitelisted"], false)
		assertOutput.Equal(integrationsMap["campaignKey"], "FR_T_100_W_100")
	}
	instance.IsFeatureEnabled("FR_T_100_W_100", userID, nil)
	instance.Integrations.CallBack = func(integrationsMap map[string]interface{}) {
		assertOutput.Equal(integrationsMap["fromUserStorageService"], false)
		assertOutput.Equal(integrationsMap["isFeatureEnabled"], true)
		assertOutput.Equal(integrationsMap["isUserWhitelisted"], false)
		assertOutput.Equal(integrationsMap["variationName"], nil)
		assertOutput.Equal(integrationsMap["variationId"], nil)
	}
	instance.GetFeatureVariableValue("FR_T_100_W_100", "STRING_VARIABLE", userID, nil)
}

func TestIntegrationFeatureTestCampaign(t *testing.T) {
	assertOutput := assert.New(t)
	instance := GetVWOClientInstance("FT_T_100_W_10_20_30_40")

	instance.Integrations.CallBack = func(integrationsMap map[string]interface{}) {
		assertOutput.Equal(integrationsMap["fromUserStorageService"], false)
		assertOutput.Equal(integrationsMap["isFeatureEnabled"], true)
		assertOutput.Equal(integrationsMap["isUserWhitelisted"], false)
		assertOutput.Equal(integrationsMap["campaignKey"], "FT_T_100_W_10_20_30_40")
	}
	instance.IsFeatureEnabled("FT_T_100_W_10_20_30_40", "Ashley", nil)

	options := make(map[string]interface{})
	options["variationTargetingVariables"] = map[string]interface{}{"chrome": false}
	instance = GetVWOClientInstance("FT_100_W_33_33_33_WS_WW")

	instance.Integrations.CallBack = func(integrationsMap map[string]interface{}) {
		assertOutput.Equal(integrationsMap["fromUserStorageService"], false)
		assertOutput.Equal(integrationsMap["isFeatureEnabled"], false)
		assertOutput.Equal(integrationsMap["isUserWhitelisted"], true)
		assertOutput.Equal(integrationsMap["variationName"], "Variation-2")
		assertOutput.Equal(integrationsMap["variationId"], 3)
	}

	instance.GetFeatureVariableValue("FT_100_W_33_33_33_WS_WW", "STRING_VARIABLE", "Ashley", options)
}
