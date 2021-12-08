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

package constants

// constants
const (
	MaxTrafficPercent = 100
	MaxTrafficValue   = 10000
	StatusRunning     = "RUNNING"
	SDKVersion        = "1.16.0"
	SDKName           = "vwo-go-sdk"
	Platform          = "server"
	SeedValue         = 1

	CampaignTypeVisualAB       = "VISUAL_AB"
	CampaignTypeFeatureTest    = "FEATURE_TEST"
	CampaignTypeFeatureRollout = "FEATURE_ROLLOUT"

	GoalTypeRevenue         = "REVENUE_TRACKING"
	GoalTypeCustom          = "CUSTOM_GOAL"
	GoalTypeAll             = "ALL"
	GoalIdentifierSeperator = "_vwo_"

	PushAPITagValueLength = 255
	PushAPITagKeyLength   = 255

	OperatorTypeAnd = "and"
	OperatorTypeOr  = "or"
	OperatorTypeNot = "not"

	OperandTypesCustomVariable = "custom_variable"
	OperandTypesUser           = "user"

	HTTPSProtocol            = "https://"
	EndPointsBaseURL         = "dev.visualwebsiteoptimizer.com"
	EndPointsAccountSettings = "/server-side/settings"
	EndPointsTrackUser       = "/server-side/track-user"
	EndPointsTrackGoal       = "/server-side/track-goal"
	EndPointsPush            = "/server-side/push"

	BaseURL                = "dev.visualwebsiteoptimizer.com"
	AccountSettings        = "/server-side/settings"
	TrackUser              = "/server-side/track-user"
	TrackGoal              = "/server-side/track-goal"
	Push                   = "/server-side/push"
	WebHookAccountSettings = "/server-side/pull"
	BatchEndPoint          = "/server-side/batch-events"

	Boolean = "boolean"
	Double  = "double"
	Integer = "integer"
	String  = "string"

	LowerMatch    = `^lower\((.*)\)`
	WildcardMatch = `^wildcard\((.*)\)`
	RegexMatch    = `^regex\((.*)\)`
	StartingStar  = `^\*`
	EndingStar    = `\*$`

	LowerValue              = 1
	StartingEndingStarValue = 2
	StartingStarValue       = 3
	EndingStarValue         = 4
	RegexValue              = 5
	EqualValue              = 6

	Info    = "INFO"
	Debug   = "DEBUG"
	Error   = "ERROR"
	Warning = "WARN"

	EventsTrackUser = "EVENTS_TRACK_USER"
	EventsTrackGoal = "EVENTS_TRACK_GOAL"
	EventsPush      = "EVENTS_PUSH"

	HttpPostMethod = "POST"

	BatchMinEventsPerRequest = 1
	BatchMaxEventsPerRequest = 5000
	BatchMinRequestInterval  = 1

	BatchDefaultEventsPerRequest = 100
	BatchDefaultRequestInterval  = 600

	CampaignDecisionType = "CAMPAIGN_DECISION"
)

var EventTypeMapping = map[string]int{
	EventsTrackUser: 1,
	EventsTrackGoal: 2,
	EventsPush:      3,
}
