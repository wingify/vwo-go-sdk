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

package schema

import (
	"github.com/wingify/vwo-go-sdk/pkg/constants"
)

type Integrations struct {
	CallBack func(map[string]interface{})
}

func (integrations *Integrations) ExecuteCallBack(integrationsMap map[string]interface{}, fromUserStorage bool, campaign Campaign, variation Variation, isUserWhitelisted bool) {
	if integrations.CallBack != nil && integrationsMap != nil && len(variation.Name) > 0 {
		integrationsMap["fromUserStorageService"] = fromUserStorage
		integrationsMap["isUserWhitelisted"] = isUserWhitelisted
		if campaign.Type == constants.CampaignTypeFeatureRollout {
			integrationsMap["isFeatureEnabled"] = true
		} else {
			if campaign.Type == constants.CampaignTypeFeatureTest {
				integrationsMap["isFeatureEnabled"] = variation.IsFeatureEnabled
			}
			integrationsMap["variationName"] = variation.Name
			integrationsMap["variationId"] = variation.ID
		}
		integrations.CallBack(integrationsMap)
	}
}
