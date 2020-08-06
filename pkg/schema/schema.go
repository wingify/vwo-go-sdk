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

package schema

// VwoInstance struct utils
type VwoInstance struct {
	SettingsFile             SettingsFile
	UserStorage              interface{}
	Logger                   interface{}
	IsDevelopmentMode        bool
	UserID                   string
	Campaign                 Campaign
	API                      string
	GoalTypeToTrack          interface{}
	ShouldTrackReturningUser interface{}
}

// SettingsFile struct
type SettingsFile struct {
	SDKKey    string     `json:"sdkKey"`
	Campaigns []Campaign `json:"campaigns"`
	AccountID int        `json:"accountId"`
}

// Campaign struct
type Campaign struct {
	ID                int                    `json:"id"`
	Segments          map[string]interface{} `json:"segments"`
	Status            string                 `json:"status"`
	PercentTraffic    int                    `json:"percentTraffic"`
	Goals             []Goal                 `json:"goals"`
	Variations        []Variation            `json:"variations"`
	Variables         []Variable             `json:"variables"`
	IsForcedVariation bool                   `json:"isForcedVariationEnabled"`
	Key               string                 `json:"key"`
	Type              string                 `json:"type"`
}

// Goal struct
type Goal struct {
	Identifier string `json:"identifier"`
	ID         int    `json:"id"`
	Type       string `json:"type"`
}

// Variation struct
type Variation struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Changes  interface{}            `json:"changes"`
	Weight   float64                `json:"weight"`
	Segments map[string]interface{} `json:"segments"`

	Variables        []Variable `json:"variables"`
	IsFeatureEnabled bool       `json:"isFeatureEnabled"`

	StartVariationAllocation int
	EndVariationAllocation   int
}

// Variable struct
type Variable struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	ID    int         `json:"id"`
}

// Options struct
type Options struct {
	CustomVariables             map[string]interface{} `json:"customVariables"`
	VariationTargetingVariables map[string]interface{} `json:"variationTargetingVariables"`
	RevenueValue                interface{}
	GoalTypeToTrack             interface{}
	ShouldTrackReturningUser    interface{}
}

// UserData  struct
type UserData struct {
	UserID         string
	CampaignKey    string
	VariationName  string
	GoalIdentifier string
}

// VariationAllocationRange struct
type VariationAllocationRange struct {
	StartRange int
	EndRange   int
}

// Impression struct
type Impression struct {
	AccountID    int     `json:"account_id"`
	UID          string  `json:"uId"`
	Random       float32 `json:"random"`
	SID          string  `json:"sId"`
	U            string  `json:"u"`
	Sdk          string  `json:"sdk"`
	SdkV         string  `json:"sdk-v"`
	Ap           string  `json:"ap"`
	URL          string  `json:"url"`
	ExperimentID int     `json:"experiment_id"`
	Combination  int     `json:"combination"`
	ED           string  `json:"ed"`
	GoalID       int     `json:"goal_id"`
	R            string  `json:"r"`
	Tags         string  `json:"tags"`
}

// TrackResult struct
type TrackResult struct {
	CampaignKey string
	TrackValue bool
}