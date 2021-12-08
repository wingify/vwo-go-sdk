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
	"net/url"

	"github.com/wingify/vwo-go-sdk/pkg/logger"
)

func GetUsageStatsObject(vwoInstance VwoInstance) ( usageStats map[string]string ){
	usageStats = make(map[string]string)
	if vwoInstance.Integrations.CallBack != nil {
		usageStats["ig"] = "1"
	}
	if vwoInstance.IsBatchingEnabled {
		usageStats["eb"] = "1"
	}
	if vwoInstance.Logger != nil {
		usageStats["cl"] = "1"
	}
	if vwoInstance.UserStorage != nil {
		usageStats["ss"] = "1"
	}
	if vwoInstance.Logger != nil {
		if (logger.GetLogLevel() > 0 ){
			usageStats["ll"] = "1"
		}
	}
	if vwoInstance.ShouldTrackReturningUser != nil {
		usageStats["tr"] = "1"
	}
	if vwoInstance.GoalTypeToTrack != nil {
		usageStats["gt"] = "1"
	}
	usageStats["_l"] = "1"
	return 
}


func GetUsageStatsImpression(vwoInstance VwoInstance) (usageStats string) {	
	params := url.Values{}
	for key, element := range GetUsageStatsObject(vwoInstance) {
		params.Add(key, element)
    }
	usageStats = "&"+params.Encode()
	return 
}