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

package core

import (
	"fmt"
	"math"
	"strconv"

	"github.com/spaolacci/murmur3"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

const (
	umax32Bit = 0xFFFFFFFF
	bucketer  = "bucketer.go"
)

// GetBucketerVariation function returns the Variation by checking the Start and End Bucket Allocations of each Variation
func GetBucketerVariation(vwoInstance schema.VwoInstance, variations []schema.Variation, bucketValue int, userID, campaignKey string) (schema.Variation, error) {
	/*
		Args:
			variations : list of variations (schema.Variation)
			bucketValue: the bucket value of the user

		Returns:
			schema.Variation: variation  allotted to the user
			error: if no variation found, else nil
	*/

	for _, variation := range variations {
		if variation.StartVariationAllocation <= bucketValue && variation.EndVariationAllocation >= bucketValue {
			message := fmt.Sprintf(constants.DebugMessageGotVariationForUser, vwoInstance.API, userID, campaignKey, vwoInstance.Campaign.Type, variation.Name, "GetBucketerVariation")
			utils.LogMessage(vwoInstance.Logger, constants.Debug, bucketer, message)
			return variation, nil
		}
	}
	return schema.Variation{}, fmt.Errorf(constants.ErrorMessageNoVariationForBucketValue, vwoInstance.API, userID, campaignKey, bucketValue)
}

// GetBucketValueForUser returns Bucket Value of the user by hashing the userId with murmur hash and scaling it down.
func GetBucketValueForUser(vwoInstance schema.VwoInstance, userID string, maxValue,
	multiplier float64, campaign schema.Campaign) (uint32, int) {
	/*
		Args:
			vwoInstance: vwo Instance for logger implementation
			userID: the unique ID assigned to User
			maxValue: maximum value that can be alloted to the bucket value
			multiplier: value for distributing ranges slightly

		Returns:
			int: the bucket value allotted to User (between 1 to MAX_TRAFFIC_PERCENT)
	*/
	if campaign.IsBucketingSeedEnabled {
		var campaignId = strconv.Itoa(campaign.ID) //to convert campaign Id to string to append to userId
		userID = campaignId + "_" + userID
	}

	hashValue := hash(userID) & umax32Bit
	ratio := float64(hashValue) / math.Pow(2, 32)
	multipliedValue := (maxValue*ratio + 1) * multiplier
	bucketValue := int(math.Floor(multipliedValue))

	return hashValue, bucketValue
}

// IsUserPartOfCampaign calculates if the provided userID should become part of the campaign or not
func IsUserPartOfCampaign(vwoInstance schema.VwoInstance, userID string, campaign schema.Campaign) bool {
	/*
		Args:
			userID: the unique ID assigned to a user
			campaign: for getting traffic allotted to the campaign

		Returns:
			bool: if User is a part of Campaign or not
	*/

	if len(campaign.Variations) == 0 {
		return false
	}
	hashValue, valueAssignedToUser := GetBucketValueForUser(vwoInstance, userID, constants.MaxTrafficPercent, 1, campaign)

	message := fmt.Sprintf(constants.DebugMessageUserHashBucketValue, vwoInstance.API, userID, hashValue, valueAssignedToUser)
	utils.LogMessage(vwoInstance.Logger, constants.Debug, bucketer, message)

	isUserPart := valueAssignedToUser != 0 && valueAssignedToUser <= campaign.PercentTraffic

	message = fmt.Sprintf(constants.InfoMessageUserEligibilityForCampaign, vwoInstance.API, userID, isUserPart)
	utils.LogMessage(vwoInstance.Logger, constants.Info, bucketer, message)

	return isUserPart
}

// BucketUserToVariation validates the User ID and returns Variation into which the User is bucketed in.
func BucketUserToVariation(vwoInstance schema.VwoInstance, userID string, campaign schema.Campaign, disableLogs bool) (schema.Variation, error) {
	/*
		Args:
		    userID: the unique ID assigned to User
		    campaign: the Campaign of which User is a part of
			disableLogs : flag which when set to true nothing will be logged

		Returns:
			schema.Variation: variation data into which user is bucketed in
			error: if no variation found, else nil
	*/

	if len(campaign.Variations) == 0 {
		return schema.Variation{}, fmt.Errorf(constants.ErrorMessageNoVariationInCampaign, vwoInstance.API, campaign.Key)
	}
	multiplier := (float64(constants.MaxTrafficValue) / float64(campaign.PercentTraffic)) / 100
	_, bucketValue := GetBucketValueForUser(vwoInstance, userID, constants.MaxTrafficValue, multiplier, campaign)

	message := fmt.Sprintf(constants.DebugMessageVariationHashBucketValue, vwoInstance.API, userID, campaign.Key, campaign.PercentTraffic, bucketValue)
	utils.LogMessage(vwoInstance.Logger, constants.Debug, bucketer, message)

	return GetBucketerVariation(vwoInstance, campaign.Variations, bucketValue, userID, campaign.Key)
}

// hash function generates hash value for given string using murmur hash
func hash(s string) uint32 {
	hasher := murmur3.New32WithSeed(uint32(constants.SeedValue))
	hasher.Write([]byte(s))
	return hasher.Sum32()
}

// addRangesToVariations function helps to calculate range of every variation
func addRangesToVariations(variations []schema.Variation) []schema.Variation {
	/*
		Args:
			variations : array of type schema.Variation
		Return:
			variations : array of type schema.Variation after adding ranges for every variation in the array passed
	*/
	offset := 0
	for idx := range variations {
		limit := int(math.Floor(variations[idx].Weight * constants.MaxTrafficValue / 100))
		maxRange := offset + limit
		variations[idx].StartVariationAllocation = offset + 1
		variations[idx].EndVariationAllocation = maxRange
		offset = maxRange
	}
	return variations
}

// addRangesToCampaigns function helps to calculate range of every campaign
func addRangesToCampaigns(campaigns []schema.Campaign) []schema.Campaign {
	/*
		Args:
			campaigns : array of type schema.Campaign
		Return:
			campaigns : array of type schema.Campaign after adding ranges for every Campaign in the array passed
	*/
	offset := 0
	for idx := range campaigns {
		limit := int(math.Floor(campaigns[idx].Weight * constants.MaxTrafficValue / 100))
		maxRange := offset + limit
		campaigns[idx].MinRange = offset + 1
		campaigns[idx].MaxRange = maxRange
		offset = maxRange
	}
	return campaigns
}

// getCamapignsUsingRange function Returns a campaign by checking the Start and End Bucket Allocations of each campaign.
func getCampaignUsingRange(rangeForCampaigns int, campaigns []schema.Campaign) (schema.Campaign, error) {
	/*
		Args:
			rangeForCampaigns : the bucket value of the user
			campaigns         : array of campaigns of type schema.campaign
		Return:
			if ranges of the campaign are well within the bucket value --> correspoding campaign
			else -->nil
	*/
	rangeForCampaigns = rangeForCampaigns * constants.MaxTrafficValue
	for _, currentCampaign := range campaigns {
		if currentCampaign.MaxRange != 0 && currentCampaign.MaxRange >= rangeForCampaigns && currentCampaign.MinRange <= rangeForCampaigns {
			return currentCampaign, nil
		}
	}
	return schema.Campaign{}, fmt.Errorf(constants.ErrorMessageNoCampaignIsInRange, rangeForCampaigns) //to add some error message
	//fmt.Errorf(constants.ErrorMessageNoVariationForBucketValue, vwoInstance.API, userID, campaignKey, bucketValue)
}
