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
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
)

const impressions = "impression.go"

// CreateImpressionForPush creates the impression from the arguments passed to push
func CreateImpressionForPush(vwoInstance schema.VwoInstance, tagKey, tagValue, userID string) schema.Impression {
	/*
		Args:
			tagKey : Campaign identifier
			tagValue : Variation identifier
			userId : User identifier

		Returns:
			schema.Impression: Imression struct with required values
	*/
	impression := getCommonProperties(vwoInstance, userID)
	impression.URL = constants.HTTPSProtocol + constants.EndPointsBaseURL + constants.EndPointsPush

	impression.Tags = `{"u":{"` + url.QueryEscape(tagKey) + `":"` + url.QueryEscape(tagValue) + `"}}`

	message := fmt.Sprintf(constants.DebugMessageImpressionForPush, vwoInstance.API, impression.AccountID, impression.UID, impression.SID, impression.URL, impression.Tags)
	LogMessage(vwoInstance.Logger, constants.Debug, impressions, message)

	return impression
}

// CreateImpressionTrackingGoal creates the impression from the arguments passed to track goal
func CreateImpressionTrackingGoal(vwoInstance schema.VwoInstance, variationID int, userID, goalType string, campaignID, goalID int, revenueValue interface{}) schema.Impression {
	/*
		Args:
		    variationID : Variation identifier
			userID : User identifier
			campaignID : Campaign identifier
		    goalID : Goal identifier
		    revenueValue : Revenue goal for the campaign

		Returns:
			schema.Impression: Imression struct with required values
	*/
	impression := getCommonProperties(vwoInstance, userID)

	impression.ExperimentID = campaignID
	impression.Combination = variationID

	impression.URL = constants.HTTPSProtocol + constants.EndPointsBaseURL + constants.EndPointsTrackGoal
	impression.GoalID = goalID

	if goalType == constants.GoalTypeRevenue {
		switch revenueValue.(type) {
		case int:
			impression.R = strconv.Itoa(revenueValue.(int))
		case float32:
			impression.R = strconv.FormatFloat(float64(revenueValue.(float32)), 'f', -1, 32)
		case float64:
			impression.R = strconv.FormatFloat(float64(revenueValue.(float64)), 'f', -1, 64)
		case string:
			impression.R = revenueValue.(string)
		}
	}

	if goalType == constants.GoalTypeRevenue {
		message := fmt.Sprintf(constants.DebugMessageImpressionForTrackRevenueGoal, vwoInstance.API, impression.AccountID, impression.UID, impression.SID, impression.URL, impression.ExperimentID, impression.Combination, impression.GoalID, revenueValue)
		LogMessage(vwoInstance.Logger, constants.Debug, impressions, message)
	} else {
		message := fmt.Sprintf(constants.DebugMessageImpressionForTrackCustomGoal, vwoInstance.API, impression.AccountID, impression.UID, impression.SID, impression.URL, impression.ExperimentID, impression.Combination, impression.GoalID)
		LogMessage(vwoInstance.Logger, constants.Debug, impressions, message)
	}

	return impression
}

// CreateImpressionTrackingUser creates the impression from the arguments passed to track user
func CreateImpressionTrackingUser(vwoInstance schema.VwoInstance, campaignID int, variationID int, userID string) schema.Impression {
	/*
		Args:
			variationID : Variation identifier
			userID : User identifier
			campaignID : Campaign identifier

		Returns:
			schema.Impression: Imression struct with required values
	*/
	impression := getCommonProperties(vwoInstance, userID)

	impression.ExperimentID = campaignID
	impression.Combination = variationID

	impression.ED = `{\"p\":\"` + constants.Platform + `\"}`
	impression.URL = constants.HTTPSProtocol + constants.EndPointsBaseURL + constants.EndPointsTrackUser

	message := fmt.Sprintf(constants.DebugMessageImpressionForTrackUser, vwoInstance.API, impression.AccountID, impression.UID, impression.SID, impression.URL, impression.ExperimentID, impression.Combination, impression.ED)
	LogMessage(vwoInstance.Logger, constants.Debug, impressions, message)

	return impression
}

// getCommonProperties returns commonly used params for making requests to our servers.
func getCommonProperties(vwoInstance schema.VwoInstance, userID string) schema.Impression {
	/*
		Args:
			userID : Unique identification of user

		Returns:
			schema.Impression: commonly used params for making call to our servers
	*/
	return schema.Impression{
		Random:    rand.Float32(),
		Sdk:       constants.SDKName,
		SdkV:      constants.SDKVersion,
		Ap:        constants.Platform,
		SID:       strconv.FormatInt(time.Now().Unix(), 10),
		U:         generateFor(vwoInstance, userID, vwoInstance.SettingsFile.AccountID),
		AccountID: vwoInstance.SettingsFile.AccountID,
		UID:       url.PathEscape(userID),
	}
}
