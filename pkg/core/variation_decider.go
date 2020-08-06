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

package core

import (
	"fmt"
	"strconv"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

const variationDecider = "variationDecider.go"

// VariationDecider struct
type VariationDecider struct {
	Bucketer         string
	SegmentEvaluator string
}

// GetVariation function
/*	Returns variation for the user for given campaign
    This method achieves the variation assignment in the following way:
    1. First get variation from UserStorage, if variation is found in user_storage_data,
    return from there
    2. Evaluates white listing users for each variation, and find a targeted variation.
    3. If no targeted variation is found, evaluate pre-segmentation result
    4. Evaluate percent traffic
    5. If user becomes part of campaign assign a variation.
	6. Store the variation found in the user_storage
*/
func GetVariation(vwoInstance schema.VwoInstance, userID string, campaign schema.Campaign, goalIdentifier string, options schema.Options) (schema.Variation, string, error) {
	/*
		Args:
			userId: the unique ID assigned to User
			campaign: campaign in which user is participating
			customVariables(In option): variables for pre-segmentation
			variationTargetingVariables(In option): variables for variation targeting
			revenueValue(In option): Value of revenue for the goal if the goal is revenue tracking

		Returns:
			schema.Variation: Struct object containing the information regarding variation assigned else empty object
			error: Error message
	*/
	vwoInstance.UserID = userID
	vwoInstance.Campaign = campaign

	_, ok := options.VariationTargetingVariables["_vwo_user_id"]
	if !ok {
		if options.VariationTargetingVariables == nil {
			options.VariationTargetingVariables = make(map[string]interface{})
		}
		options.VariationTargetingVariables["_vwo_user_id"] = userID
	}

	targettedVariation, err := FindTargetedVariation(vwoInstance, userID, campaign, options)
	if err != nil {
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, err.Error())
	} else {
		message := fmt.Sprintf(constants.InfoMessageGotVariationForUser, vwoInstance.API, userID, campaign.Key, campaign.Type, targettedVariation.Name)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)
		return targettedVariation, "", nil
	}

	variationName, storedGoalIdentifier := GetVariationFromUserStorage(vwoInstance, userID, campaign)
	if variationName != "" {
		message := fmt.Sprintf(constants.InfoMessageGotStoredVariation, vwoInstance.API, variationName, campaign.Key, userID)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)
		variation, err := utils.GetCampaignVariation(vwoInstance.API, campaign, variationName)
		return variation, storedGoalIdentifier, err 
	}

	if !IsUserPartOfCampaign(vwoInstance, userID, campaign) {
		return schema.Variation{}, "", fmt.Errorf(constants.DebugMessageUserNotPartOfCampaign, vwoInstance.API, userID, campaign.Key, campaign.Type, "IsUserPartOfCampaign")
	}

	if EvaluateSegment(vwoInstance, campaign.Segments, options) {
		variation, err := BucketUserToVariation(vwoInstance, userID, campaign)
		if err != nil {
			return schema.Variation{}, "", fmt.Errorf(constants.InfoMessageUserGotNoVariation, vwoInstance.API, userID, campaign.Key, err.Error())
		}

		if vwoInstance.UserStorage == nil {
			message := fmt.Sprintf(constants.DebugMessageNoUserStorageServiceSet, vwoInstance.API)
			utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
		} else {
			if storage, ok := vwoInstance.UserStorage.(interface{ Set(a, b, c, d string) }); ok {
				storage.Set(userID, campaign.Key, variation.Name, goalIdentifier)
				message := fmt.Sprintf(constants.InfoMessageSettingDataUserStorageService, vwoInstance.API, userID)
				utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)
			} else {
				message := fmt.Sprintf(constants.ErrorMessageSetUserStorageServiceFailed, vwoInstance.API, userID)
				utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
			}
		}

		message := fmt.Sprintf(constants.InfoMessageVariationAllocated, vwoInstance.API, userID, campaign.Key, variation.Name)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)

		return variation, "", nil
	}

	return schema.Variation{}, "", fmt.Errorf(constants.ErrorMessageNoVariationAlloted, vwoInstance.API, userID, campaign.Key, campaign.Type)
}

// FindTargetedVariation function Identifies and retrives if there exists any targeted
// variation in the given campaign for given userID
func FindTargetedVariation(vwoInstance schema.VwoInstance, userID string, campaign schema.Campaign, options schema.Options) (schema.Variation, error) {
	/*
		Args:
			userId: the unique ID assigned to User
			campaign: campaign in which user is participating
			customVariables(In option): variables for pre-segmentation
			variationTargetingVariables(In option): variables for variation targeting
			revenueValue(In option): Value of revenue for the goal if the goal is revenue tracking

		Returns:
			schema.Variation: Struct object containing the information regarding variation assigned else empty object
			string: Log level
			error: Error message
	*/

	if campaign.IsForcedVariation == false {
		return schema.Variation{}, fmt.Errorf(constants.InfoMessageWhitelistingSkipped, vwoInstance.API, userID, campaign.Key)
	}
	whiteListedVariationsList := GetWhiteListedVariationsList(vwoInstance, userID, campaign, options)
	whiteListedVariationsLength := len(whiteListedVariationsList)
	var targettedVariation schema.Variation
	if whiteListedVariationsLength == 0 {
		return schema.Variation{}, fmt.Errorf(constants.InfoMessageNoWhitelistedVariation, vwoInstance.API, campaign.Key)
	} else if whiteListedVariationsLength == 1 {
		targettedVariation = whiteListedVariationsList[0]
	} else {
		whiteListedVariationsList = utils.ScaleVariations(whiteListedVariationsList)
		whiteListedVariationsList = utils.GetVariationAllocationRanges(vwoInstance, whiteListedVariationsList)
		_, bucketValue := GetBucketValueForUser(vwoInstance, userID, constants.MaxTrafficValue, 1)
		var err error
		targettedVariation, err = GetBucketerVariation(vwoInstance, whiteListedVariationsList, bucketValue, userID, campaign.Key)
		if err != nil {
			return schema.Variation{}, fmt.Errorf(constants.InfoMessageNoTargettedVariation, vwoInstance.API, err.Error())
		}

		message := fmt.Sprintf(constants.InfoMessageSegmentationStatusForVariation, vwoInstance.API, userID, campaign.Key, targettedVariation.Segments, options.VariationTargetingVariables, "True", "WhiteListing", targettedVariation.Name)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)

		message = fmt.Sprintf(constants.InfoMessageForcedvariationAllocated, vwoInstance.API, userID, campaign.Key, campaign.Type, targettedVariation.Name)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)
	}
	return targettedVariation, nil
}

// GetVariationFromUserStorage function tries retrieving variation from user_storage
func GetVariationFromUserStorage(vwoInstance schema.VwoInstance, userID string, campaign schema.Campaign) (string, string) {
	/*
		Args:
			userId: the unique ID assigned to User
			campaign: campaign in which user is participating

		Returns:
			variationName: Name of the found varaition in the user storage
			string: Log level
			error: Error message
	*/

	if vwoInstance.UserStorage == nil {
		message := fmt.Sprintf(constants.InfoMessageNoUserStorageServiceGet, vwoInstance.API)
		utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
		return "", ""
	}
	if storage, ok := vwoInstance.UserStorage.(interface {
		Get(a, b string) schema.UserData
	}); ok {
		userStorageFetch := storage.Get(userID, campaign.Key)
		message := fmt.Sprintf(constants.DebugMessageGettingStoredVariation, vwoInstance.API, userID, campaign.Key)
		utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
		if userStorageFetch.VariationName == "" {
			message := fmt.Sprintf(constants.DebugMessageNoStoredVariation, vwoInstance.API, userID, campaign.Key)
			utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
		}
		return userStorageFetch.VariationName, userStorageFetch.GoalIdentifier
	}
	
	message := fmt.Sprintf(constants.ErrorMessageGetUserStorageServiceFailed, vwoInstance.API, userID)
	utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
	return "", ""
}

//GetWhiteListedVariationsList function identifies all forced variations which are targeted by variation_targeting_variables
func GetWhiteListedVariationsList(vwoInstance schema.VwoInstance, userID string, campaign schema.Campaign, options schema.Options) []schema.Variation {
	/*
		Args:
			userId: the unique ID assigned to User
			campaign: campaign in which user is participating
			customVariables(In option): variables for pre-segmentation
			variationTargetingVariables(In option): variables for variation targeting
			revenueValue(In option): Value of revenue for the goal if the goal is revenue tracking

		Returns:
			schema.Variation: Struct object containing the information regarding variation assigned else empty object
	*/

	var whiteListedVariationsList []schema.Variation
	for _, variation := range campaign.Variations {
		if len(variation.Segments) == 0 {
			message := fmt.Sprintf(constants.DebugMessageNoSegmentsInVariation, vwoInstance.API, userID, campaign.Key, variation.Name)
			utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)

			message = fmt.Sprintf(constants.DebugMessageSegmentationStatusForVariation, vwoInstance.API, userID, campaign.Key, options.VariationTargetingVariables, variation.Segments, "False", "WhiteListing", variation.Name)
			utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)

			continue
		}

		status := PreEvaluateSegment(vwoInstance, variation.Segments, options, variation.Name)
		if status {
			whiteListedVariationsList = append(whiteListedVariationsList, variation)
		}

		message := fmt.Sprintf(constants.DebugMessageSegmentationStatusForVariation, vwoInstance.API, userID, campaign.Key, options.VariationTargetingVariables, variation.Segments, strconv.FormatBool(status), "WhiteListing", variation.Name)
		utils.LogMessage(vwoInstance.Logger, constants.Debug, variationDecider, message)
	}
	return whiteListedVariationsList
}

// EvaluateSegment function evaluates segmentation for the userID against the segments found inside the campaign.
func EvaluateSegment(vwoInstance schema.VwoInstance, segments map[string]interface{}, options schema.Options) bool {
	/*
		Args:
			segments: segments from campaign or variation
			options: options object containing CustomVariables, VariationTargertting variables and Revenue Goal

		Returns:
			bool: if the options falls in the segments criteria
	*/

	if len(segments) == 0 {
		message := fmt.Sprintf(constants.DebugMessageSegmentationSkipped, vwoInstance.API, vwoInstance.UserID, vwoInstance.Campaign.Key)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)

		return true
	}

	status := SegmentEvaluator(segments, options.CustomVariables)

	message := fmt.Sprintf(constants.InfoMessageSegmentationStatus, vwoInstance.API, vwoInstance.UserID, vwoInstance.Campaign.Key, segments, options.CustomVariables, strconv.FormatBool(status), "PreSegmentation")
	utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)

	return status
}

// PreEvaluateSegment function evaluates segmentation for the userID against the segments found inside the campaign.
func PreEvaluateSegment(vwoInstance schema.VwoInstance, segments map[string]interface{}, options schema.Options, variationName string) bool {
	/*
		Args:
			segments: segments from campaign or variation
			options: options object containing CustomVariables, VariationTargertting variables and Revenue Goal

		Returns:
			bool: if the options falls in the segments criteria
	*/

	if len(options.VariationTargetingVariables) == 0 {
		message := fmt.Sprintf(constants.DebugMessageSegmentationSkippedForVariation, vwoInstance.API, vwoInstance.UserID, vwoInstance.Campaign.Key, variationName)
		utils.LogMessage(vwoInstance.Logger, constants.Info, variationDecider, message)

		return false
	}
	return SegmentEvaluator(segments, options.VariationTargetingVariables)
}
