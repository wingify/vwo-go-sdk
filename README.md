# VWO GO SDK

[![Build Status](https://img.shields.io/travis/wingify/vwo-go-sdk)](http://travis-ci.org/wingify/vwo-go-sdk)
![Size in Bytes](https://img.shields.io/github/languages/code-size/wingify/vwo-go-sdk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![codecov](https://codecov.io/gh/wingify/vwo-go-sdk/branch/master/graph/badge.svg)](https://codecov.io/gh/wingify/vwo-go-sdk)

This open source library allows you to A/B Test your Website at server-side.

## Requirements

- Works with Go 1.11.4+

## Installation

```go
go get "github.com/wingify/vwo-go-sdk"
```

## Basic usage

**Importing and Instantiation**

```go
import vwo "github.com/wingify/vwo-go-sdk"
import "github.com/wingify/vwo-go-sdk/pkg/api"

// Get SettingsFile
settingsFile := vwo.GetSettingsFile("accountID", "SDKKey")

// Default instance of VwoInstance
vwoClientInstance, err := vwo.Launch(settingsFile)
if err != nil {
	//handle err
}

// Instance with custom options
vwoClientInstance, err := vwo.Launch(settingsFile, api.WithDevelopmentMode())
if err != nil {
	//handle err
}

// Activate API
// With Custom Variables
options := make(map[string]interface{})
options["customVariables"] = map[string]interface{}{"a": "x"}
options["variationTargetingVariables"] = map[string]interface{}{"a": "x"}
options["revenueValue"] = 12
variationName = vwoClientInstance.Activate(campaignKey, userID, options)

// Without Custom Variables
variationName = vwoClientInstance.Activate(campaignKey, userID, nil)


// GetVariation
// With Custom Variables
options := make(map[string]interface{})
options["customVariables"] = map[string]interface{}{"a": "x"}
variationName = vwoClientInstance.GetVariationName(campaignKey, userID, options)

//Without Custom Variables
variationName = vwoClientInstance.GetVariationName(campaignKey, userID, nil)


// Track API
// With Custom Variables
options := make(map[string]interface{})
options["customVariables"] = map[string]interface{}{"a": "x"}
isSuccessful = vwoClientInstance.Track(campaignKey, userID, goalIdentifier, options)

// With Revenue Value
options := make(map[string]interface{})
options["revenueValue"] = 12
isSuccessful = vwoClientInstance.Track(campaignKey, userID, goalIdentifier, options)

// With Custom Variables, Revenue Value, GoalTypeToTrack and ShouldTrackreturningUser
options := make(map[string]interface{})
options["customVariables"] = map[string]interface{}{"a": "x"}
options["revenueValue"] = 12
//  Set specific goalType to Track
//  Available GoalTypes - constants.GoalTypeRevenue, constants.GoalTypeCustom, constants.GoalTypeAll (Default)
options["goalTypeToTrack"] = constants.GoalTypeAll
//  Set if a return user should be tracked, default false
options["ShouldTrackreturningUser"] = false
isSuccessful = vwoInstance.Track(campaignKey, userID, goalIdentifier, options)

// For Goal Conversion in Multiple Campaign
// campaignKeys = []string{"campaignKey1", "campaignKey2", "campaignKey3"}
// For Goal Conversion in All Possible Campaigns
// campaignKeys = nil
isSuccessful = vwoInstance.Track(campaignKeys, userID, goalIdentifier, options)


// FeatureEnabled API
// With Custom Varibles
options := make(map[string]interface{})
options["customVariables"] = map[string]interface{}{"a": "x"}
isSuccessful = vwoClientInstance.IsFeatureEnabled(campaignKey, userID, options)

// Without Custom Variables
isSuccessful = vwoClientInstance.IsFeatureEnabled(campaignKey, userID, nil)

// GetFeatureVariableValue API
// With Custom Variables
options := make(map[string]interface{})
options["customVariables"] = map[string]interface{}{"a": "x"}
variableValue = vwoClientInstance.GetFeatureVariableValue(campaignKey, variableKey, userID, options)

// Without Custom Variables
variableValue = vwoClientInstance.GetFeatureVariableValue(campaignKey, variableKey, userID, nil)

// Push API
isSuccessful = vwoClientInstance.Push(tagKey, tagValue, userID)
```

**User Storage**

```go
import vwo "github.com/wingify/vwo-go-sdk/"
import "github.com/wingify/vwo-go-sdk/pkg/api"
import "github.com/wingify/vwo-go-sdk/pkg/schema"

// declare UserStorage interface with the following Get & Set function signature
type UserStorage interface{
    Get(userID, campaignKey string) UserData
    Set(string, string, string, string)
}

// declare a UserStorageData struct to implement UserStorage interface
type UserStorageData struct{}

// Get method to fetch user variation from storage
func (us *UserStorageData) Get(userID, campaignKey string) schema.UserData {
    //Example code showing how to get userData  from DB
    userData, ok := userDatas[campaignKey]
    if ok {
		for _, userdata := range userData {
			if userdata.UserID == userID {
				return userdata
			}
		}
    }
    /*
    // UserData  struct
    type UserData struct {
		UserID         string
		CampaignKey    string
		VariationName  string
		GoalIdentifier string
	}
    */
	return schema.UserData{}
}

// Set method to save user variation to storage
func (us *UserStorageData) Set(userID, campaignKey, variationName, goalIdentifer string) {
    //Example code showing how to store userData in DB
    userdata := schema.UserData{
		UserID:        userID,
		CampaignKey:   campaignKey,
		VariationName: variationName,
		GoalIdentifier: goalIdentifier,
	}

	flag := false
	userData, ok := userDatas[userdata.CampaignKey]
	if ok {
		for _, user := range userData {
			if user.UserID == userdata.UserID {
				flag = true
			}
		}
		if !flag {
			userDatas[userdata.CampaignKey] = append(userDatas[userdata.CampaignKey], userdata)
		} else {
			for i, user := range userData {
				if user.UserID == userdata.UserID && user.CampaignKey == userdata.CampaignKey {
					userData[i].VariationName = userdata.VariationName
					userData[i].GoalIdentifier = userdata.GoalIdentifier
				}
			}
		}
	} else {
		userDatas[userdata.CampaignKey] = []schema.UserData{
			userdata,
		}
	}

func main() {
	settingsFile := vwo.GetSettingsFile("accountID", "SDKKey")
	// create UserStorageData object
	storage := &UserStorageData{}

	vwoClientInstance, err := vwo.Launch(settingsFile, api.WithStorage(storage))
	if err != nil {
		//handle err
	}
}

```

**Custom Logger**

```go
import vwo "github.com/wingify/vwo-go-sdk"
import "github.com/wingify/vwo-go-sdk/pkg/api"

// declare Log interface with the following CustomLog function signature
type Log interface {
	CustomLog(level, errorMessage string)
}

// declare a LogS struct to implement Log interface
type LogS struct{}

// Get function to handle logs
func (c *LogS) CustomLog(level, errorMessage string) {}

func main() {
	settingsFile := vwo.GetSettingsFile("accountID", "SDKKey")
	// create LogS object
	logger := &LogS{}

	vwoClientInstance, err := vwo.Launch(settingsFile, api.WithLogger(logger))
	if err != nil {
		//handle err
	}
}
```

## Demo App

[Example](https://github.com/wingify/vwo-go-sdk-example)

## Documentation

Refer [Official VWO FullStack Documentation](https://developers.vwo.com/reference#fullstack-introduction)

## Local Setup

1. Install dependencies

```bash
go get .
```

2. Configure the environment

```bash
bash start.sh
```

## Running Unit Tests

```shell
./test.sh
```

## Third-party Resources and Credits

Refer [third-party-attributions.txt](third-party-attribution.txt)

## Contributing

Please go through our [contributing guidelines](CONTRIBUTING.md)

## Code of Conduct

[Code of Conduct](CODE_OF_CONDUCT.md)

## License

[Apache License, Version 2.0](LICENSE)
