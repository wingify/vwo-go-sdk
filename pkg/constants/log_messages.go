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

package constants

// constants for logger
const (
	//Debug Messages
	DebugMessageCustomLoggerUsed                = "Custom logger used"
	DebugMessageDevelopmentMode                 = "Development mode is : %v "
	DebugMessageGettingStoredVariation          = "[%v] Got User Storage, Checking stored variation for User ID: %v of Campaign: %v "
	DebugMessageGotVariationForUser             = "[%v] User ID: %v of Campaign: %v campaignType: %v got variation: %v inside method: %v "
	DebugMessageImpressionForPush               = "[%v] impression built for pushing - AccountID: %v, UserID: %v, SID: %v, URL: %v, Tags: %v "
	DebugMessageImpressionForTrackCustomGoal    = "[%v] impression built for track goal -  AccountID: %v, UserID: %v, SID: %v, URL: %v, ExperimentID: %v, Combination: %v, GoalID: %v "
	DebugMessageImpressionForTrackRevenueGoal   = "[%v] impression built for track goal -  AccountID: %v, UserID: %v, SID: %v, URL: %v, ExperimentID: %v, Combination: %v, GoalID: %v, RevenueValue: %v "
	DebugMessageImpressionForTrackUser          = "[%v] impression built for track user - AccountID: %v, UserID: %v, SID: %v, URL: %v, ExperimentID: %v, Combination: %v, ED: %v"
	DebugMessageSDKInitialized                  = "SDK properly initialized"
	DebugMessageNoCustomLoggerFound             = "No custom logger found, using pre-defined google logger "
	DebugMessageNoStoredVariation               = "[%v] No stored variation for User ID: %v for Campaign: %v found in UserStorageService"
	DebugMessageNoUserStorageServiceSet         = "[%v] No UserStorageService to set data"
	DebugMessageSegmentationSkipped             = "[%v] For User ID: %v of CampaignKey: %v segments are missing, hence skipping Presegmentation"
	DebugMessageSegmentationSkippedForVariation = "[%v] For User ID: %v of CampaignKey: %v Variation Targeting variables are missing, hence skipping segmentation for variation %v "
	DebugMessageSegmentationStatusForVariation  = "[%v] For User ID: %v of Campaign: %v with Variation Targeting variables: %v, Segments: %v, %v, %v for variation %v "
	DebugMessageUserHashBucketValue             = "[%v] User ID: %v having hash: %v got bucketValue: %v "
	DebugMessageUserNotPartOfCampaign           = "[%v] User ID: %v for CampaignKey: %v type: %v did not become part of campaign method: %v "
	DebugMessageUUIDForUser                     = "[%v] Uuid generated for User ID: %v and accountId: %v is %v "
	DebugMessageVariationHashBucketValue        = "[%v] User ID: %v for CampaignKey: %v having percent traffic: %v got bucket value: %v "

	/*Extras*/
	DebugMessageCustomLoggerFound     = "[%v] Custom logger found"
	DebugMessageNoSegmentsInVariation = "[%v] For User ID: %v of Campaign: %v, segment was missing, hence skipping segmentation %v "
	DebugMessageSettingsFileProcessed = "[%v] Settings file processed"
	DebugMessageValidConfiguration    = "[%v] SDK configuration and account settings are valid"

	//Error Messages
	ErrorMessageActivateAPIMissingParams                = "[%v] activate API got bad parameters. It expects campaignKey(String) as first, User ID(String) as second and options(Optional) as third argument"
	ErrorMessageCampaignNotRunning                      = "[%v] API Campaign: %v is not RUNNING. Please verify from VWO App"
	ErrorMessageCustomLoggerMisconfigured               = "Custom logger is provided but seems to have misconfigured. Please check the API Docs. Using default logger."
	ErrorMessageGetFeatureVariableMissingParams         = "[%v] getFeatureVariableValue API got bad parameters. It expects campaignKey(String) as first, variableKey(String) as second, User ID(String) as third, and options as fourth argument"
	ErrorMessageGetUserStorageServiceFailed             = "[%v] Getting data from UserStorageService failed for User ID: %v "
	ErrorMessageGetVariationAPIMissingParams            = "[%v] getVariation API got bad parameters. It expects campaignKey(String) as first, User ID(String) as second and options(Optional) as third argument"
	ErrorMessageImpressionFailed                        = "[%v] Impression event could not be sent to VWO endpoint: %v "
	ErrorMessageInvalidAPI                              = "[%v] API is not valid for Campaign: %v of type: %v for User ID: %v "
	ErrorMessageIsFeatureEnabledAPIMissingParams        = "[%v] isFeatureEnabled API got bad parameters. It expects Campaign(String) as first, User ID(String) as second and options(Optional) as third argument"
	ErrorMessageNoCampaignInCampaignList                = "[%v] No campaign found as per the required attributes : %v %v "
	ErrorMessageNoCampaignFoundForGoal                  = "[%v] No campaign found for Goal Identifier: %v with goal type to track : %v. Please verify from VWO app."
	ErrorMessagePushAPIMissingParams                    = "[%v] push API got bad parameters. It expects tagKey(String) as first, tagKey(String) as second and User ID(String) as third argument"
	ErrorMessageSettingsFileCorrupted                   = "[%v] Settings file is corrupted. Please contact VWO Support for help : %v "
	ErrorMessageSetUserStorageServiceFailed             = "[%v] Error while saving data into UserStorage for User ID: %v."
	ErrorMessageTagKeyLengthExceeded                    = "[%v] Length of tagKey: %v for User ID: %v can not be greater than 255"
	ErrorMessageTagValueLengthExceeded                  = "[%v] Length of value: %v of tagKey: %v for User ID: %v can not be greater than 255"
	ErrorMessageTrackAPIGoalNotFound                    = "[%v] Goal: %v not found for Campaign: %v and User ID: %v : %v "
	ErrorMessageTrackAPIMissingParams                   = "[%v] %v got bad parameters. It expects campaignKey(String / array of string / nil) as first, User ID(String) as second, goalIdentifier(String) as third argument and options(Optional) as fourth parameter but got : %v"
	ErrorMessageTrackAPIRevenueNotPassedForRevenueValue = "[%v] Revenue value should be passed for revenue, Goal: %v for Campaign: %v and User ID: %v "
	ErrorMessageVariableNotFound                        = "[%v] Variable: %v not found for User ID: %v for campaign %v of type %v "

	/*Extras*/
	ErrorMessageCampaignNotFound                          = "[%v] Campaign key: %v not found : %v "
	ErrorMessageCannotProcessSettingsFile                 = "[%v] Error processing settings file err : %v "
	ErrorMessageCannotReadSettingsFile                    = "[%v] Settings file could not be read and processed. Please contact VWO Support for help : %v "
	ErrorMessageCouldNotGetURL                            = "[%v] Failed get request for URL: %v "
	ErrorMessageGoalNotFound                              = "[%v] Goal: %v not found"
	ErrorMessageInvalidAccountID                          = "[%v] AccountId is required for fetching account settings. Aborting"
	ErrorMessageInvalidCampaignKeyType                    = "[%v] Campaign key should only be of type nil, array of strings or string not : %T"
	ErrorMessageInvalidGoalType                           = "[%v] Invalid goal type, Goal type to track is : %v but recieved goal type : %v"
	ErrorMessageInvalidLoggerStorage                      = "[%v] Invalid storage object/Logger given. Refer documentation on how to pass custom storage."
	ErrorMessageInvalidSDKKey                             = "[%v] SDKKey is required for fetching account settings. Aborting"
	ErrorMessageInvalidSettingsFile                       = "[%v] Settings-file fetched is not proper : %v "
	ErrorMessageNoVariationAlloted                        = "[%v] User ID: %v of CampaignKey: %v type: %v did not get any variation "
	ErrorMessageNoVariationForBucketValue                 = "[%v] No variation found for user ID %v in campaignKey: %v having bucket value: %v "
	ErrorMessageNoVariationInCampaign                     = "[%v] No variations in campaign: %v "
	ErrorMessageResponseNotParsed                         = "[%v] Error parsing response for URL: %v "
	ErrorMessageTrackAPIEmptyParam                        = "Empty %v"
	ErrorMessageTrackAPIIncorrectParamType                = "Incorrect data type for %v"
	ErrorMessageTrackAPIIncorrectShouldTrackReturningUser = "ShouldTrackReturningUser should only have value true or false but got %v"
	ErrorMessageTrackAPIIncorrectGoalTypeToTrack          = "GoalTypeTotrack should only have value REVENUE_TRACKING, CUSTOM_GOAL or ALL but got %v"
	ErrorMessageURLNotFound                               = "[%v] URL not Found: %v "
	ErrorMessageVariationNotFound                         = "[%v] Variation : %v not found in campaign : %v "

	//Info Messages
	InfoMessageFeatureEnabledForUser            = "[%v] Campaign: %v for user ID: %v is enabled"
	InfoMessageFeatureNotEnabledForUser         = "[%v] Campaign: %v for user ID: %v is not enabled"
	InfoMessageForcedvariationAllocated         = "[%v] User ID: %v of CampaignKey: %v type: %v got forced-variation: %v "
	InfoMessagesGoalAlreadyTracked              = "[%v] Goal: %v of Campaign: %v for User ID:%v has already been tracked earlier. Skipping now."
	InfoMessageGettingDataUserStorageService    = "[%v] Getting data into UserStorageService for User ID: %v successful"
	InfoMessageGotStoredVariation               = "[%v] Got stored variation: %v of CampaignKey: %v for User ID: %v from UserStorage"
	InfoMessageGotVariationForUser              = "[%v] User ID: %v for CampaignKey: %v type: %v got variation_name: %v "
	InfoMessageImpressionSuccess                = "[%v] Impression event - %v was successfully received by VWO having keys: %v "
	InfoMessageIncorrectCampaignKeyType         = "[%v] Incorrect CampaignKey type passed : %T is incorrect, should be of type string, array of string or nil"
	InfoMessageInvalidVariationKey              = "[%v] Variation was not assigned to User ID: %v for Campaign: %v : %v "
	InfoMessageMainKeysForFeatureTestImpression = "[%v] Having main keys AccountID: %v, UserID: %v, CampaignID: %v, VariationID: %v"
	InfoMessageMainKeysForPushAPI               = "[%v] Having main keys: AccountID: %v User ID: %v U: %v and tags: %v "
	InfoMessageMainKeysForImpression            = "[%v] Having main keys: AccountID: %v User ID: %v campaignId: %v and VariationID: %v "
	InfoMessageNoUserStorageServiceGet          = "[%v] No UserStorageService to get stored data"
	InfoMessageSegmentationStatus               = "[%v] For User ID: %v of Campaign: %v with Segments: %v, Custom Variables: %v, %v, %v "
	InfoMessageSegmentationStatusForVariation   = "[%v] For User ID: %v of Campaign: %v with Segments: %v, Variation targeting Variables: %v, %v, %v for variation %v "
	InfoMessageSettingDataUserStorageService    = "[%v] Setting data into UserStorageService for User ID: %v successful"
	InfoMessageUserEligibilityForCampaign       = "[%v] Is User ID: %v part of campaign ? %v "
	InfoMessageUserGotNoVariation               = "[%v] User ID: %v for Campaign: %v did not allot any variation : %v "
	InfoMessageVariationAllocated               = "[%v] User ID: %v of Campaign: %v got variation: %v "
	InfoMessageVariationRangeAllocation         = "[%v] Variation: %v with weight: %v got range as: ( %v - %v )"
	InfoMessageWhitelistingSkipped              = "[%v] For User ID: %v of Campaign: %v, whitelisting was skipped"

	/*Extras*/
	InfoMessageNoTargettedVariation      = "[%v] No targetted variation found : %v "
	InfoMessageNoWhitelistedVariation    = "[%v] No whitelisting variation found in campaign: %v "
	InfoMessageUserRecievedVariableValue = "[%v] Value for variable: %v of feature flag: %v is: %v for user: %v "
)
