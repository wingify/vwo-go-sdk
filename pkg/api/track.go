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
	"strings"

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
func (vwo *VWOInstance) Track(campaignKeys interface{}, userID, goalIdentifier string, option interface{}) []schema.TrackResult {
	/*
		Args:
			campaignKeys: Key of the running campaign
			userID: Unique identification of user
			goalIdentifier: Unique identification of corresponding goal
			customVariables(In option): variables for pre-segmentation
			variationTargetingVariables(In option): variables for variation targeting
			revenueValue(In option): Value of revenue for the goal if the goal is revenue tracking

		Returns:
			schema.TrackResult : Array of Key value pair of CampaignKey and its tracking result
	*/

	vwoInstance := schema.VwoInstance{
		SettingsFile:             vwo.SettingsFile,
		UserStorage:              vwo.UserStorage,
		Logger:                   vwo.Logger,
		IsDevelopmentMode:        vwo.IsDevelopmentMode,
		API:                      "Track",
		GoalTypeToTrack:          vwo.GoalTypeToTrack,
		ShouldTrackReturningUser: vwo.ShouldTrackReturningUser,
	}

	options := utils.ParseOptions(option)

	isValid, err := utils.ValidateTrack(userID, goalIdentifier, options.GoalTypeToTrack, options.ShouldTrackReturningUser)
	if !isValid {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIMissingParams, vwoInstance.API, "Track API", err)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return []schema.TrackResult{}
	}

	isValid, err = utils.ValidateTrack(userID, goalIdentifier, vwo.GoalTypeToTrack, vwo.ShouldTrackReturningUser)
	if !isValid {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIMissingParams, vwoInstance.API, "VWO Instance", err)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return []schema.TrackResult{}
	}

	var Campaigns []schema.Campaign
	var goalTypeToTrack string
	var shouldTrackReturningUser bool

	if options.GoalTypeToTrack == nil {
		if vwo.GoalTypeToTrack == nil {
			goalTypeToTrack = constants.GoalTypeAll
		} else {
			goalTypeToTrack = vwo.GoalTypeToTrack.(string)
		}
	} else {
		goalTypeToTrack = options.GoalTypeToTrack.(string)
	}

	if options.ShouldTrackReturningUser == nil {
		if vwo.ShouldTrackReturningUser == nil {
			shouldTrackReturningUser = false
		} else {
			shouldTrackReturningUser = vwo.ShouldTrackReturningUser.(bool)
		}
	} else {
		shouldTrackReturningUser = options.ShouldTrackReturningUser.(bool)
	}

	if campaignKeys == nil {
		CampaignList, err := utils.GetCampaignForGoals(vwoInstance, goalIdentifier, goalTypeToTrack)
		if err != nil {
			utils.LogMessage(vwo.Logger, constants.Error, track, err.Error())
		} else {
			Campaigns = CampaignList
		}
	} else {
		switch Keys := campaignKeys.(type) {
		case []string:
			CampaignList, err := utils.GetCampaignForKeys(vwoInstance, Keys)
			if err != nil {
				utils.LogMessage(vwo.Logger, constants.Error, track, err.Error())
			} else {
				Campaigns = CampaignList
			}
		case string:
			CampaignList, err := utils.GetCampaign(vwoInstance.API, vwoInstance.SettingsFile, Keys)
			if err != nil {
				utils.LogMessage(vwo.Logger, constants.Error, track, err.Error())
			} else {
				Campaigns = append(Campaigns, CampaignList)
			}
		default:
			message := fmt.Sprintf(constants.InfoMessageIncorrectCampaignKeyType, vwoInstance.API, Keys)
			utils.LogMessage(vwo.Logger, constants.Info, track, message)
			return []schema.TrackResult{}
		}
	}

	var result []schema.TrackResult

	for _, campaign := range Campaigns {
		currResult := schema.TrackResult{
			CampaignKey: campaign.Key,
			TrackValue:  trackCampaignGoal(vwoInstance, campaign, userID, goalIdentifier, goalTypeToTrack, shouldTrackReturningUser, options),
		}
		result = append(result, currResult)
	}

	if len(result) == 0 {
		message := fmt.Sprintf(constants.ErrorMessageNoCampaignFoundForGoal, vwoInstance.API, goalIdentifier, goalTypeToTrack)
		utils.LogMessage(vwo.Logger, constants.Error, track, message)
		return []schema.TrackResult{}
	}

	return result
}

func trackCampaignGoal(vwoInstance schema.VwoInstance, campaign schema.Campaign, userID, goalIdentifier, goalTypeToTrack string, shouldTrackReturningUser bool, options schema.Options) bool {
	/*
		Args:
			campaign: campaign whose user is to be tracked
			userID: Unique identification of user
			goalIdentifier: Unique identification of corresponding goal
			customVariables(In option): variables for pre-segmentation
			variationTargetingVariables(In option): variables for variation targeting
			revenueValue(In option): Value of revenue for the goal if the goal is revenue tracking

		Returns:
			bool: True if the track is successfull else false
	*/

	if campaign.Status != constants.StatusRunning {
		message := fmt.Sprintf(constants.ErrorMessageCampaignNotRunning, vwoInstance.API, campaign.Key)
		utils.LogMessage(vwoInstance.Logger, constants.Error, track, message)
		return false
	}

	if utils.CheckCampaignType(campaign, constants.CampaignTypeFeatureRollout) {
		message := fmt.Sprintf(constants.ErrorMessageInvalidAPI, vwoInstance.API, campaign.Key, campaign.Type, userID)
		utils.LogMessage(vwoInstance.Logger, constants.Error, track, message)
		return false
	}

	goal, err := utils.GetCampaignGoal(vwoInstance.API, campaign, goalIdentifier)
	if err != nil {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIGoalNotFound, vwoInstance.API, goalIdentifier, campaign.Key, userID, err.Error())
		utils.LogMessage(vwoInstance.Logger, constants.Error, track, message)
		return false
	}

	if goal.Type != goalTypeToTrack && goalTypeToTrack != constants.GoalTypeAll {
		message := fmt.Sprintf(constants.ErrorMessageInvalidGoalType, vwoInstance.API, goalTypeToTrack, goal.Type)
		utils.LogMessage(vwoInstance.Logger, constants.Error, track, message)
		return false
	}

	if goal.Type == constants.GoalTypeRevenue && options.RevenueValue == nil {
		message := fmt.Sprintf(constants.ErrorMessageTrackAPIRevenueNotPassedForRevenueValue, vwoInstance.API, goalIdentifier, campaign.Key, userID)
		utils.LogMessage(vwoInstance.Logger, constants.Error, track, message)
		return false
	}

	variation, storedGoalIdentifier, err := core.GetVariation(vwoInstance, userID, campaign, goalIdentifier, options)
	if err != nil {
		message := fmt.Sprintf(constants.InfoMessageInvalidVariationKey, vwoInstance.API, userID, campaign.Key, err.Error())
		utils.LogMessage(vwoInstance.Logger, constants.Info, track, message)
		return false
	}

	if variation.Name != "" {
		if storedGoalIdentifier != "" {
			identifiers := strings.Split(storedGoalIdentifier, constants.GoalIdentifierSeperator)
			flag := false

			for _, identifier := range identifiers {
				if identifier == goalIdentifier {
					flag = true
				}
			}

			if flag == false {
				storedGoalIdentifier = storedGoalIdentifier + constants.GoalIdentifierSeperator + goalIdentifier

				if storage, ok := vwoInstance.UserStorage.(interface{ Set(a, b, c, d string) }); ok {
					storage.Set(userID, campaign.Key, variation.Name, storedGoalIdentifier)
					message := fmt.Sprintf(constants.InfoMessageSettingDataUserStorageService, vwoInstance.API, userID)
					utils.LogMessage(vwoInstance.Logger, constants.Info, track, message)
				} else {
					message := fmt.Sprintf(constants.ErrorMessageSetUserStorageServiceFailed, vwoInstance.API, userID)
					utils.LogMessage(vwoInstance.Logger, constants.Debug, track, message)
				}
			} else if shouldTrackReturningUser == false {
				message := fmt.Sprintf(constants.InfoMessagesGoalAlreadyTracked, vwoInstance.API, goalIdentifier, campaign.Key, userID)
				utils.LogMessage(vwoInstance.Logger, constants.Info, track, message)
				return false
			}
		} else {
			if vwoInstance.UserStorage == nil {
				message := fmt.Sprintf(constants.DebugMessageNoUserStorageServiceSet, vwoInstance.API)
				utils.LogMessage(vwoInstance.Logger, constants.Debug, track, message)
			} else {
				if storage, ok := vwoInstance.UserStorage.(interface{ Set(a, b, c, d string) }); ok {
					storage.Set(userID, campaign.Key, variation.Name, goalIdentifier)
					message := fmt.Sprintf(constants.InfoMessageSettingDataUserStorageService, vwoInstance.API, userID)
					utils.LogMessage(vwoInstance.Logger, constants.Info, track, message)
				} else {
					message := fmt.Sprintf(constants.ErrorMessageSetUserStorageServiceFailed, vwoInstance.API, userID)
					utils.LogMessage(vwoInstance.Logger, constants.Debug, track, message)
				}
			}
		}

		impression := utils.CreateImpressionTrackingGoal(vwoInstance, variation.ID, userID, goal.Type, campaign.ID, goal.ID, options.RevenueValue)
		go event.DispatchTrackingGoal(vwoInstance, goal.Type, impression)

		message := fmt.Sprintf(constants.InfoMessageMainKeysForImpression, vwoInstance.API, vwoInstance.SettingsFile.AccountID, vwoInstance.UserID, campaign.ID, variation.ID)
		utils.LogMessage(vwoInstance.Logger, constants.Info, activate, message)

		return true
	}

	return false
}
