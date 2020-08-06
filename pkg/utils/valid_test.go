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
	"io/ioutil"
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/stretchr/testify/assert"
)

type CLog interface {
	CustomLog(a, b string)
}
type CLogS struct{}

func (c *CLogS) CustomLog(a, b string) {}

type WLog interface {
	CustomLogger(a, b string)
}
type WLogS struct{}

func (w *WLogS) CustomLogger(a, b string) {}

func TestValidateLogger(t *testing.T) {
	logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
	actual := ValidateLogger(logs)
	assert.True(t, actual, "google logger not validated")

	correctLog := &CLogS{}
	actual = ValidateLogger(correctLog)
	assert.True(t, actual)

	wrongLog := &WLogS{}
	actual = ValidateLogger(wrongLog)
	assert.False(t, actual)

	actual = ValidateLogger(nil)
	assert.False(t, actual)
}

type CUserStorage interface {
	Get(userID, campaignKey string) schema.UserData
	Set(userID, campaignKey, variationName string)
}
type CUserStorageData struct{}

func (us *CUserStorageData) Get(userID, campaignKey string) schema.UserData {
	return schema.UserData{}
}
func (us *CUserStorageData) Set(userID, campaignKey, variationName string) {}

type WUserStorage interface {
	Getter(userID, campaignKey string) schema.UserData
	Setter(userID, campaignKey, variationName string)
}
type WUserStorageData struct{}

func (us *WUserStorageData) Getter(userID, campaignKey string) schema.UserData {
	return schema.UserData{}
}
func (us *WUserStorageData) Setter(userID, campaignKey, variationName string) {}

func TestValidateStorage(t *testing.T) {
	actual := ValidateStorage(nil)
	assert.True(t, actual)

	correctStorage := &CUserStorageData{}
	actual = ValidateStorage(correctStorage)
	assert.True(t, actual)

	wrongStorage := &WUserStorageData{}
	actual = ValidateStorage(wrongStorage)
	assert.False(t, actual)
}

func TestParseOptions(t *testing.T) {
	expected := schema.Options{}
	expected.CustomVariables = make(map[string]interface{})
	expected.VariationTargetingVariables = make(map[string]interface{})
	expected.ShouldTrackReturningUser = nil
	expected.GoalTypeToTrack = nil
	actual := ParseOptions(nil)
	assert.Equal(t, expected, actual)

	data := make(map[string]interface{})
	data["customVariables"] = map[string]interface{}{"a": "x"}
	data["variationTargetingVariables"] = map[string]interface{}{"a": "x"}
	data["revenueValue"] = 12
	data["goalTypeToTrack"] = "ALL"
	data["shouldTrackReturningUser"] = false
	expected = schema.Options{
		CustomVariables:             map[string]interface{}{"a": "x"},
		VariationTargetingVariables: map[string]interface{}{"a": "x"},
		RevenueValue:                12,
		GoalTypeToTrack:             "ALL",
		ShouldTrackReturningUser:    false,
	}
	actual = ParseOptions(data)
	assert.Equal(t, expected, actual)
}

func TestValidateActivate(t *testing.T) {
	actual := ValidateActivate("", "")
	assert.False(t, actual)

	actual = ValidateActivate("campaignKey", "userID")
	assert.True(t, actual)
}

func TestValidateGetFeatureVariableValue(t *testing.T) {
	actual := ValidateGetFeatureVariableValue("", "", "")
	assert.False(t, actual)

	actual = ValidateGetFeatureVariableValue("campaignKey", "variableKey", "userID")
	assert.True(t, actual)
}

func TestValidateGetVariationName(t *testing.T) {
	actual := ValidateGetVariationName("", "")
	assert.False(t, actual)

	actual = ValidateGetVariationName("campaignKey", "userID")
	assert.True(t, actual)
}

func TestValidateIsFeatureEnabled(t *testing.T) {
	actual := ValidateIsFeatureEnabled("", "")
	assert.False(t, actual)

	actual = ValidateIsFeatureEnabled("campaignKey", "userID")
	assert.True(t, actual)
}

func TestValidatePush(t *testing.T) {
	actual := ValidatePush("", "", "")
	assert.False(t, actual)

	actual = ValidatePush("tagKey", "tagValue", "userID")
	assert.True(t, actual)
}

func TestValidateTrack(t *testing.T) {
	actual, _ := ValidateTrack("", "", "", "")
	assert.False(t, actual)

	actual, _ = ValidateTrack("userID", "goalIdentifier", nil, nil)
	assert.True(t, actual)

	actual, _ = ValidateTrack("userID", "goalIdentifier", constants.GoalTypeRevenue, true)
	assert.True(t, actual)

	actual, _ = ValidateTrack("userID", "goalIdentifier", "Invalid_Type", true)
	assert.False(t, actual)

	actual, _ = ValidateTrack("userID", "goalIdentifier", "CUSTOM_GOAL", 123)
	assert.False(t, actual)

	actual, _ = ValidateTrack("userID", "goalIdentifier", 123, true)
	assert.False(t, actual)
}
