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
	"github.com/wingify/vwo-go-sdk/pkg/schema"
)

// CheckCampaignType matches campaign type
func CheckCampaignType(campaign schema.Campaign, campaignType string) bool {
	/*
		Args:
			campaign : Campaign object
			campaignType : Type of campaign to be matched
		Return:
			bool: true if the type matches else false
	*/
	return campaign.Type == campaignType
}

// GetKeyValue returns first key value pair of the given map
func GetKeyValue(obj map[string]interface{}) (string, interface{}) {
	/*
		Args:
			obj: map whose firsr key value pair is needed

		Return:
			string: Key
			interface: value
	*/
	for k, v := range obj {
		return k, v
	}
	return "", nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
