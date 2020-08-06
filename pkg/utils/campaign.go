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
	"math"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
)

const campaign = "campaign.go"

// GetVariationAllocationRanges returns list of variation with set allocation ranges.
func GetVariationAllocationRanges(vwoInstance schema.VwoInstance, variations []schema.Variation) []schema.Variation {
	/*
		Args:
			variations: List of variations(schema.Variation)

		Returns:
			variations: List of variations(schema.Variation)
	*/

	var (
		currentAllocation         = 0
		variationAllocationRanges []schema.Variation
	)
	for _, variation := range variations {
		stepFactor := GetVariationBucketingRange(variation.Weight)
		if stepFactor != 0 {
			variation.StartVariationAllocation = currentAllocation + 1
			variation.EndVariationAllocation = currentAllocation + stepFactor
			currentAllocation += stepFactor
		} else {
			variation.StartVariationAllocation = -1
			variation.EndVariationAllocation = -1
		}

		message := fmt.Sprintf(constants.InfoMessageVariationRangeAllocation, vwoInstance.API, variation.Name, variation.Weight, variation.StartVariationAllocation, variation.EndVariationAllocation)
		LogMessage(vwoInstance.Logger, constants.Info, campaign, message)
		variationAllocationRanges = append(variationAllocationRanges, variation)
	}
	return variationAllocationRanges
}

// GetVariationBucketingRange returns the bucket size of variation.
func GetVariationBucketingRange(weight float64) int {
	/*
		Args:
			weight: Weight of variation

		Returns:
			int: Bucket start range of Variation
	*/

	if weight == 0 {
		return 0
	}
	startRange := int(math.Ceil(weight * 100))
	return min(startRange, constants.MaxTrafficValue)
}

// GetCampaign function finds and returns campaign from given campaign_key.
func GetCampaign(API string, settingsFile schema.SettingsFile, campaignKey string) (schema.Campaign, error) {
	/*
		Args:
			settingsFile  : Settings file for the project
			campaignKey: Campaign identifier key

		Returns:
			schema.Campaign: Campaign object
	*/
	for _, campaign := range settingsFile.Campaigns {
		if campaign.Key == campaignKey {
			return campaign, nil
		}
	}
	return schema.Campaign{}, fmt.Errorf(constants.ErrorMessageCampaignNotFound, API, campaignKey, "")
}

// GetCampaignForKeys function returns list of campaigns from the settings file that are in the list of CampaignKeys
func GetCampaignForKeys(vwoInstance schema.VwoInstance, campaignKeys []string) ([]schema.Campaign, error){
	/*
		Args:
			settingsFile  : Settings file for the project
			campaignKeys: Array of campaign keys to be searched

		Returns:
			[]schema.Campaign: Array of matching campaigns 
	*/

	var Campaigns []schema.Campaign
	for _, campaignKey := range campaignKeys {
		Campaign, err := GetCampaign(vwoInstance.API, vwoInstance.SettingsFile, campaignKey)
		if err != nil {
			LogMessage(vwoInstance.Logger, constants.Error, campaign, err.Error())
		} else {
			Campaigns = append(Campaigns, Campaign)
		}
	}

	if len(Campaigns) == 0 {
		return Campaigns, fmt.Errorf(constants.ErrorMessageNoCampaignInCampaignList, vwoInstance.API, campaignKeys, "")
	}
	return Campaigns, nil
}

// GetCampaignForGoals function returns list of campaigns from the settings file that are in the list of CampaignKeys
func GetCampaignForGoals(vwoInstance schema.VwoInstance, goalIdentifier, goalTypeToTrack string) ([]schema.Campaign, error){
	/*
		Args:
			settingsFile  : Settings file for the project
			goalidentifier : Goal to be searched in the campaigns
			goalTypeToTrack : Type the searched goal should be 

		Returns:
			[]schema.Campaign: Array of matching campaigns 
	*/

	var Campaigns []schema.Campaign
	for _, Campaign := range vwoInstance.SettingsFile.Campaigns {
		goal, err := GetCampaignGoal(vwoInstance.API, Campaign, goalIdentifier)
		if err != nil {
			LogMessage(vwoInstance.Logger, constants.Error, campaign, err.Error())
		} else {
			if goal.Type == goalTypeToTrack || goalTypeToTrack == constants.GoalTypeAll {
				Campaigns = append(Campaigns, Campaign)
			}
		}
	}
	
	if len(Campaigns) == 0 {
		return Campaigns, fmt.Errorf(constants.ErrorMessageNoCampaignInCampaignList, vwoInstance.API, goalIdentifier, goalTypeToTrack)
	}
	return Campaigns, nil
}



// ScaleVariations function It extracts the weights from all the variations inside the
// campaign and scales them so that the total sum of eligible variations weights become 100%
func ScaleVariations(variations []schema.Variation) []schema.Variation {
	/*
		Args:
			variations: List of variations(schema.Variartion) having weight as a property

		Return:
			variations: List of variations(schema.Variartion)
	*/
	weightSum := 0.0
	for _, variation := range variations {
		weightSum += variation.Weight
	}
	if weightSum == 0 {
		normalizedWeight := 100.0 / float64(len(variations))
		for i := range variations {
			variations[i].Weight = normalizedWeight
		}
	} else {
		for _, variation := range variations {
			variation.Weight = (variation.Weight / weightSum) * 100
		}
	}
	return variations
}

// GetCampaignGoal returns goal from given campaign and goal identifier.
func GetCampaignGoal(API string, campaign schema.Campaign, goalIdentifier string) (schema.Goal, error) {
	/*
		 Args:
			campaign: The running campaign
			goalIdentifier: Goal identifier

		Returns:
			schema.Goal: Goal corresponding to goal_identifer in respective campaign
	*/
	goals := campaign.Goals
	for _, goal := range goals {
		if goal.Identifier == goalIdentifier {
			return goal, nil
		}
	}
	return schema.Goal{}, fmt.Errorf(constants.ErrorMessageGoalNotFound, API, goalIdentifier)
}

// GetCampaignVariation returns variation from given campaign and variationName.
func GetCampaignVariation(API string, campaign schema.Campaign, variationName string) (schema.Variation, error) {
	/*
		 Args:
			campaign: The running campaign
			variationName: Variation identifier

		Returns:
			schema.Variation: Variation corresponding to variationName in respective campaign
	*/
	if len(campaign.Variations) == 0 {
		return schema.Variation{}, fmt.Errorf(constants.ErrorMessageNoVariationInCampaign, API, campaign.Key)
	}
	for _, variation := range campaign.Variations {
		if variation.Name == variationName {
			return variation, nil
		}
	}

	return schema.Variation{}, fmt.Errorf(constants.ErrorMessageVariationNotFound, API, variationName, campaign.Key)
}

// GetControlVariation returns control variation from a given campaign
func GetControlVariation(campaign schema.Campaign) schema.Variation {
	/*
		Args:
			campaign: Running campaign

		Returns:
			schema.Variation: Control variation from the campaign, ie having id = 1
	*/

	for _, variation := range campaign.Variations {
		if variation.ID == 1 {
			return variation
		}
	}
	return schema.Variation{}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
