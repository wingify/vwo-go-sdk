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

package api

import (
	"fmt"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/core"
	"github.com/wingify/vwo-go-sdk/pkg/event"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

const track = "track.go"

// Track function
/*
This API method: Marks the conversion of the campaign for a particular goal
1. Validates the arguments being passed
2. Finds the corresponding Campaign
3. Checks the Campaign Status
4. Validates the Campaign Type
5. Gets passed Goal
6. Validates Revenue and Goal type
7. Assigns the determinitic variation to the user(based on userId), if user becomes part of campaign
   If userStorageService is used, it will look into it for the variation and if found, no further processing is done
8. If feature enabled, sends a call to VWO server for tracking visitor
*/
func (vwo *VWOInstance) Track(campaignKey, userID, goalIdentifier string, option interface{}) bool {
	/*
		Args:
			campaignKey: Key of the running campaign
			userID: Unique identification of user
			goalIdentifier: Unique identification of corresponding goal
			customVariables(In option): variables for pre-segmentation
			variationTargetingVariables(In option): variables for variation targeting
			revenueValue(In option): Value of revenue for the goal if the goal is revenue tracking

		Returns:
			bool: True if the track is successfull else false
	*/

	vwoInstance := schema.VwoInstance{
		SettingsFile:      vwo.SettingsFile,
		UserStorage:       vwo.UserStorage,
		Logger:            vwo.Logger,
		IsDevelopmentMode: vwo.IsDevelopmentMode,
		API:               "Track",
	}

	if !utils.ValidateTrack(campaignKey, userID, goalIdentifier) {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIMissingParams, vwoInstance.API)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return false
	}

	options := utils.ParseOptions(option)

	campaign, err := utils.GetCampaign(vwoInstance.API, vwo.SettingsFile, campaignKey)
	if err != nil {
		message := fmt.Sprintf(constants.ErrorMessageCampaignNotFound, vwoInstance.API, campaignKey, err.Error())
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return false
	}

	if campaign.Status != constants.StatusRunning {
		message := fmt.Sprintf(constants.ErrorMessageCampaignNotRunning, vwoInstance.API, campaignKey)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return false
	}
	if utils.CheckCampaignType(campaign, constants.CampaignTypeFeatureRollout) {
		message := fmt.Sprintf(constants.ErrorMessageInvalidAPI, vwoInstance.API, campaignKey, campaign.Type, userID)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return false
	}

	goal, err := utils.GetCampaignGoal(vwoInstance.API, campaign, goalIdentifier)
	if err != nil {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIGoalNotFound, vwoInstance.API, goalIdentifier, campaignKey, userID, err.Error())
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return false
	}

	if goal.Type == constants.GoalTypeRevenue && options.RevenueValue == nil  {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIRevenueNotPassedForRevenueValue, vwoInstance.API, goalIdentifier, campaignKey, userID)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return false
	}

	variation, err := core.GetVariation(vwoInstance, userID, campaign, options)
	if err != nil {
		message := fmt.Sprintf(constants.InfoMessageInvalidVariationKey, vwoInstance.API, userID, campaignKey, err.Error())
		utils.LogMessage(vwo.Logger, constants.Info, track, message)
		return false
	}

	impression := utils.CreateImpressionTrackingGoal(vwoInstance, variation.ID, userID, goal.Type, campaign.ID, goal.ID, options.RevenueValue)
	go event.DispatchTrackingGoal(vwoInstance, goal.Type, impression)

	message := fmt.Sprintf(constants.InfoMessageMainKeysForImpression, vwoInstance.API, vwoInstance.SettingsFile.AccountID, vwoInstance.UserID, campaign.ID, variation.ID)
	utils.LogMessage(vwo.Logger, constants.Info, activate, message)

	return true
}
