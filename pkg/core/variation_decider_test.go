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
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/wingify/vwo-go-sdk/pkg/tests"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

type TestCase struct {
	User      string `json:"user"`
	Variation string `json:"variation"`
}

type SegmentsHelper struct {
	Or ORHelper `json:"or"`
}

type ORHelper struct {
	CustomVariables custVar `json:"custom_variable"`
}

type custVar struct {
	Chrome  string `json:"chrome"`
	Browser string `json:"browser"`
}

func TestPreEvaluateSegment(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_100_W_33_33_33")

	segments := vwoInstance.SettingsFile.Campaigns[0].Segments
	options := schema.Options{
		VariationTargetingVariables: nil,
	}
	value := PreEvaluateSegment(vwoInstance, segments, options, "")
	assert.False(t, value, "Expected False as no segments")
}

func TestEvaluateSegment(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("T_50_W_50_50_WS")

	segments := vwoInstance.SettingsFile.Campaigns[0].Segments
	options := schema.Options{
		CustomVariables: map[string]interface{}{"a": "123", "hello": "world"},
	}
	value := EvaluateSegment(vwoInstance, segments, options, false)
	assert.True(t, value, "Expected True as mismatch")
}

func TestGetWhiteListedVariationsList(t *testing.T) {
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_100_W_0_100")

	options := schema.Options{}
	userID := testdata.GetRandomUser()
	campaign := vwoInstance.SettingsFile.Campaigns[0]
	actual := GetWhiteListedVariationsList(vwoInstance, userID, campaign, options, false)
	assert.Empty(t, actual, "No WhiteListed Variations Found")

	vwoInstance = testdata.GetInstanceWithCustomSettings("SettingsFile3")
	options = schema.Options{
		VariationTargetingVariables: map[string]interface{}{"a": "123"},
		RevenueValue:                12,
	}
	userID = testdata.GetRandomUser()
	campaign = vwoInstance.SettingsFile.Campaigns[0]
	actual = GetWhiteListedVariationsList(vwoInstance, userID, campaign, options, false)
	expected := campaign.Variations[0:2]
	assert.Equal(t, expected, actual, "No WhiteListed Variations Found")
}

func TestFindTargetedVariation(t *testing.T) {
	assertOutput := assert.New(t)

	// CORNER CASES

	instance := testdata.GetInstanceWithCustomSettings("SettingsFile3")

	campaign := instance.SettingsFile.Campaigns[0]
	options := schema.Options{
		VariationTargetingVariables: map[string]interface{}{"a": "789"},
	}
	actual, _ := FindTargetedVariation(instance, testdata.ValidUser, campaign, options, false)
	assertOutput.Equal("", actual.Name, "Variations should match")
}

func TestGetVariation(t *testing.T) {
	assertOutput := assert.New(t)

	var userExpectation map[string][]TestCase
	data, err := ioutil.ReadFile("../testdata/user_expectations1.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &userExpectation); err != nil {
		logger.Info("Error: " + err.Error())
	}

	var settingsFiles map[string]schema.SettingsFile
	data, err = ioutil.ReadFile("../testdata/settings.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &settingsFiles); err != nil {
		logger.Info("Error: " + err.Error())
	}

	logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	defer logger.Close()

	instance := schema.VwoInstance{
		Logger: logs,
	}

	goalIdentifier := ""

	for settingsFileName, settingsFile := range settingsFiles {
		vwoInstance := schema.VwoInstance{
			Logger: logs,
		}
		settingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(vwoInstance, settingsFile.Campaigns[0].Variations)

		instance.SettingsFile = settingsFile

		testCases := userExpectation[settingsFileName]
		for i := range testCases {
			campaign, err := utils.GetCampaign("", instance.SettingsFile, settingsFile.Campaigns[0].Key)
			assertOutput.Nil(err, "Incorrect Get Campaign Call")
			actual, _, _ := GetVariation(instance, testCases[i].User, campaign, goalIdentifier, schema.Options{})
			expected := testCases[i].Variation
			assertOutput.Equal(expected, actual.Name, settingsFileName+" "+testCases[i].User)
		}
	}

	// CORNER CASES

	var customSettingsFiles map[string]schema.SettingsFile
	data, err = ioutil.ReadFile("../testdata/custom_settings.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &customSettingsFiles); err != nil {
		logger.Info("Error: " + err.Error())
	}

	settings := customSettingsFiles["SettingsFile3"]
	instance.SettingsFile = settings

	campaign := instance.SettingsFile.Campaigns[0]
	options := schema.Options{
		VariationTargetingVariables: map[string]interface{}{"a": "123"},
	}

	actual, storedGoalIdentifier, err := GetVariation(instance, testdata.ValidUser, campaign, goalIdentifier, options)
	assertOutput.Nil(err, "No Variation Will Be Allcoated")
	expected := testdata.ValidVariationControl
	assert.Empty(t, storedGoalIdentifier, "No Stored variation")
	assertOutput.Equal(expected, actual.Name, "Variations should match")

	options = schema.Options{
		VariationTargetingVariables: map[string]interface{}{"b": "456"},
	}
	actual, storedGoalIdentifier, err = GetVariation(instance, testdata.ValidUser, campaign, goalIdentifier, options)
	assertOutput.Nil(err, "No Variation Will Be Allcoated")
	expected = testdata.ValidVariationVariation2
	assert.Empty(t, storedGoalIdentifier, "No Stored variation")
	assertOutput.Equal(expected, actual.Name, "Variations should match")

	instance = testdata.GetInstanceWithStorage("AB_T_50_W_50_50")
	actual, storedGoalIdentifier, err = GetVariation(instance, testdata.TempUser, instance.SettingsFile.Campaigns[0], goalIdentifier, schema.Options{})
	assertOutput.Nil(err, "No Variation Will Be Allcoated")
	expected = instance.SettingsFile.Campaigns[0].Variations[0].Name
	assert.Equal(t, testdata.DummyGoal, storedGoalIdentifier, "No Stored variation")
	assertOutput.Equal(expected, actual.Name, "Variations should match")

	instance = testdata.GetInstanceWithStorage("AB_T_100_W_20_80")
	userID := testdata.GetRandomUser()
	actual, storedGoalIdentifier, err = GetVariation(instance, userID, instance.SettingsFile.Campaigns[0], goalIdentifier, schema.Options{})
	assertOutput.NotNil(err, "No Variation Will Be Allcoated")
	assert.Empty(t, storedGoalIdentifier, "No Stored variation")
	assertOutput.Empty(actual, "Variations should be empty : "+userID)

	instance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(instance, instance.SettingsFile.Campaigns[0].Variations)
	userID = testdata.GetRandomUser()
	actual, storedGoalIdentifier, err = GetVariation(instance, userID, instance.SettingsFile.Campaigns[0], goalIdentifier, schema.Options{})
	assertOutput.Equal(nil, err, "No error expected")
	assert.Empty(t, storedGoalIdentifier, "No Stored variation")
	assertOutput.NotEmpty(actual, "Variations should match : "+userID)

	instance = testdata.GetInstanceWithIncorrectStorage("AB_T_100_W_20_80")
	instance.SettingsFile.Campaigns[0].Variations = utils.GetVariationAllocationRanges(instance, instance.SettingsFile.Campaigns[0].Variations)
	userID = testdata.GetRandomUser()
	actual, storedGoalIdentifier, err = GetVariation(instance, userID, instance.SettingsFile.Campaigns[0], goalIdentifier, schema.Options{})
	assertOutput.Equal(nil, err, "No error expected")
	assert.Empty(t, storedGoalIdentifier, "No Stored variation")
	assertOutput.NotEmpty(actual, "Variations should match : "+userID)
}

func TestGetVariationFromUserStorage(t *testing.T) {
	assertOutput := assert.New(t)
	vwoInstance := testdata.GetInstanceWithSettings("AB_T_50_W_50_50")

	campaign := vwoInstance.SettingsFile.Campaigns[0]
	userID := testdata.ValidUser
	actual, storedGoalIdentifier := GetVariationFromUserStorage(vwoInstance, userID, campaign, false)
	assertOutput.Empty(actual, "Actual and Expected Variation Name mismatch")

	vwoInstance = testdata.GetInstanceWithStorage("AB_T_50_W_50_50")

	campaign = vwoInstance.SettingsFile.Campaigns[0]
	userID = testdata.ValidUser
	expected := testdata.DummyVariation
	actual, storedGoalIdentifier = GetVariationFromUserStorage(vwoInstance, userID, campaign, false)
	assertOutput.Equal(testdata.DummyGoal, storedGoalIdentifier, "Actual and Expected goalIdentifier did not match")
	assertOutput.Equal(expected, actual, "Actual and Expected Variation Name mismatch")

	campaign = vwoInstance.SettingsFile.Campaigns[0]
	userID = testdata.InvalidUser
	expected = ""
	actual, storedGoalIdentifier = GetVariationFromUserStorage(vwoInstance, userID, campaign, false)
	assertOutput.Equal(storedGoalIdentifier, "", "Actual and Expected goalIdentifier did not match")
	assertOutput.Equal(expected, actual, "Actual and Expected Variation Name mismatch")

}

// functions for testing Mutually Exclusive Groups

// function CheckIfKeyExists checks whether key is actually present in the array or not
func CheckIfKeyExists(TrackResult []schema.TrackResult, key string) bool {
	for _, currentTrackResult := range TrackResult {
		if currentTrackResult.CampaignKey == key {
			return currentTrackResult.TrackValue //if key exists return that corresponding track value
		}
	}
	return false //if key does not exist returning false
}

func TestVariationReturnAsWhitelisting(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	//instance := testdata.GetInstanceWithSettings("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[2].Key
	options := make(map[string]interface{})
	options["variationTargetingVariables"] = map[string]interface{}{"chrome": "false"}

	// called campaign satisfies the whitelisting
	Variation := instance.Activate(campaignKey, "Ashley", options)
	TrackResult := instance.Track(campaignKey, "Ashley", "CUSTOM", options)
	VariationName := instance.GetVariationName(campaignKey, "Ashley", options)
	isGoalTracked := CheckIfKeyExists(TrackResult, campaignKey)
	assertOutput.Equal("Variation-1", Variation, "Actual and expected variations did not match")
	assertOutput.Equal("Variation-1", VariationName, "Actual and expected variation Name mismatch ")
	assertOutput.Equal(true, isGoalTracked, "Goal is not tracked correctly")
}

func TestNullVariationAsOtherCampaignSatisfiesWhitelisting(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	//instance := testdata.GetInstanceWithSettings("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[3].Key
	options := make(map[string]interface{})
	options["variationTargetingVariables"] = map[string]interface{}{"chrome": "false"}

	Variation := instance.Activate(campaignKey, "Ashley", options)
	TrackResult := instance.Track(campaignKey, "Ashley", "CUSTOM", options)
	VariationName := instance.GetVariationName(campaignKey, "Ashley", options)
	isGoalTracked := CheckIfKeyExists(TrackResult, campaignKey)
	assertOutput.Equal("", Variation, "Actual and expected variations did not match")
	assertOutput.Equal("", VariationName, "Actual and expected variation Name mismatch ")
	assertOutput.Equal(false, isGoalTracked, "Goal is not tracked correctly")
}

func TestVariationForCalledCampaign(t *testing.T) {

}

func TestNullVariationAsOtherCampaignSatisfiesStorage(t *testing.T) {

}

func TestVariationForCalledCampaignInStorageAndOtherCampaignSatisfiesWhitelisting(t *testing.T) {

}

func TestNullVariationWhenCampaignNotInGroup(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	//instance := testdata.GetInstanceWithSettings("SettingsFileMeg")
	options := make(map[string]interface{})

	campaignKey := instance.SettingsFile.Campaigns[4].Key
	Variation := instance.Activate(campaignKey, "Ashley", options)
	TrackResult := instance.Track(campaignKey, "Ashley", "CUSTOM", options)
	VariationName := instance.GetVariationName(campaignKey, "Ashley", options)
	isGoalTracked := CheckIfKeyExists(TrackResult, campaignKey)
	assertOutput.Equal("", Variation, "Actual and expected variations did not match")
	assertOutput.Equal("", VariationName, "Actual and expected variation Name mismatch ")
	assertOutput.Equal(false, isGoalTracked, "Goal is not tracked correctly")
}

func TestNoCampaignsSatisfiesPresegmentation(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[0].Key
	options := make(map[string]interface{})
	options["customVariables"] = map[string]interface{}{"browser": "chrome"}

	segmentPassed := SegmentsHelper{}
	segmentPassed.Or.CustomVariables.Chrome = "false"

	instance.SettingsFile.Campaigns[0].Segments = segmentPassed
	instance.SettingsFile.Campaigns[1].Segments = segmentPassed
	isFeatureEnabled := instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue := instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, false, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, nil, "Variable value of feature is wrongly interpreted")

	// implementing the same condition with zero traffic percentage
	instance.SettingsFile.Campaigns[0].PercentTraffic = 0
	instance.SettingsFile.Campaigns[1].PercentTraffic = 0
	isFeatureEnabled = instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue = instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, false, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, nil, "Variable value of feature is wrongly interpreted")
}

func TestCalledCampaignNotSatisfyingPresegmentation(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[0].Key
	options := make(map[string]interface{})
	options["customVariables"] = map[string]interface{}{"browser": "chrome"}

	segmentPassed := SegmentsHelper{}
	segmentPassed.Or.CustomVariables.Browser = "chrome"

	segmentFailed := SegmentsHelper{}
	segmentFailed.Or.CustomVariables.Chrome = "false"

	instance.SettingsFile.Campaigns[0].Segments = segmentFailed
	instance.SettingsFile.Campaigns[1].Segments = segmentPassed
	isFeatureEnabled := instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue := instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, false, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, nil, "Variable value of feature is wrongly interpreted")

	// implementing the same condition with different traffic percentage
	instance.SettingsFile.Campaigns[0].PercentTraffic = 0
	instance.SettingsFile.Campaigns[1].PercentTraffic = 100
	isFeatureEnabled = instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue = instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, false, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, nil, "Variable value of feature is wrongly interpreted")
}

func TestOnlyCalledCampaignSatisfyPresegmentation(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[0].Key
	options := make(map[string]interface{})
	options["customVariables"] = map[string]interface{}{"browser": "chrome"}

	segmentPassed := SegmentsHelper{}
	segmentPassed.Or.CustomVariables.Browser = "chrome"

	segmentFailed := SegmentsHelper{}
	segmentFailed.Or.CustomVariables.Chrome = "false"

	instance.SettingsFile.Campaigns[0].Segments = segmentPassed
	instance.SettingsFile.Campaigns[1].Segments = segmentFailed
	isFeatureEnabled := instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue := instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, true, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, "Control string", "Variable value of feature is wrongly interpreted")

	// implementing the same condition with different traffic percentage
	instance.SettingsFile.Campaigns[0].PercentTraffic = 100
	instance.SettingsFile.Campaigns[1].PercentTraffic = 0
	isFeatureEnabled = instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue = instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, true, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, "Control string", "Variable value of feature is wrongly interpreted")
}

func TestCalledCampaignWinnerCampaign(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[0].Key
	options := make(map[string]interface{})

	instance.SettingsFile.Campaigns[0].PercentTraffic = 100
	instance.SettingsFile.Campaigns[1].PercentTraffic = 0
	isFeatureEnabled := instance.IsFeatureEnabled(campaignKey, "Ashley", options)
	variableValue := instance.GetFeatureVariableValue(campaignKey, "STRING_VARIABLE", "Ashley", options)
	assertOutput.Equal(isFeatureEnabled, true, "Feature is not enabled correctly")
	assertOutput.Equal(variableValue, "Control string", "Variable value of feature is wrongly interpreted")

	campaignKey = instance.SettingsFile.Campaigns[2].Key
	Variation := instance.Activate(campaignKey, "Ashley", options)
	TrackResult := instance.Track(campaignKey, "Ashley", "CUSTOM", options)
	VariationName := instance.GetVariationName(campaignKey, "Ashley", options)
	isGoalTracked := CheckIfKeyExists(TrackResult, campaignKey)
	assertOutput.Equal("Control", Variation, "Actual and expected variations did not match")
	assertOutput.Equal("Control", VariationName, "Actual and expected variation Name mismatch ")
	assertOutput.Equal(true, isGoalTracked, "Goal is not tracked correctly")
}

func TestCalledCampaignNotWinnerCampaign(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[0].Key
	options := make(map[string]interface{})

	instance.SettingsFile.Campaigns[0].PercentTraffic = 100
	instance.SettingsFile.Campaigns[1].PercentTraffic = 100
	isFeatureEnabled := instance.IsFeatureEnabled(campaignKey, "lisa", options)
	assertOutput.Equal(isFeatureEnabled, false, "Feature is not enabled correctly")

	campaignKey = instance.SettingsFile.Campaigns[2].Key
	Variation := instance.Activate(campaignKey, "lisa", options)
	assertOutput.Equal("", Variation, "Actual and expected variations did not match")
}

func TestWhenEqualTrafficAmongEligibleCampaigns(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[2].Key
	options := make(map[string]interface{})

	instance.SettingsFile.Campaigns[2].PercentTraffic = 80
	instance.SettingsFile.Campaigns[3].PercentTraffic = 50
	Variation := instance.Activate(campaignKey, "Ashley", options)
	assertOutput.Equal("Variation-1", Variation, "Actual and expected variations did not match")
}

func TestWhenBothCampaignsNewToUser(t *testing.T) {
	assertOutput := assert.New(t)
	instance := tests.GetVWOClientInstance("SettingsFileMeg")
	campaignKey := instance.SettingsFile.Campaigns[2].Key
	options := make(map[string]interface{})
	Variation := instance.Activate(campaignKey, "Ashley", options)
	assertOutput.Equal("Control", Variation, "Actual and expected variations did not match")

	campaignKey = instance.SettingsFile.Campaigns[3].Key
	Variation = instance.Activate(campaignKey, "Ashley", options)
	assertOutput.Equal("", Variation, "Actual and expected variations did not match")
}

func TestWhenUserAlreadyPartOfCampaignAndNewCampaignAddedToGroup(t *testing.T) {

}

func TestWhenViewedCampaignRemovedFromGroup(t *testing.T) {

}

/*

  public function testVariationForCalledCampaign()
  {
      $campaignKey = $this->settingsFileMEG['campaigns'][2]['key'];
      $vwoInstance = TestUtil::instantiateSdk($this->settingsFileMEG, ['isUserStorage' => 1]);

      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $isGoalTracked = $vwoInstance->track($campaignKey, 'Ashley', 'CUSTOM');
      $variationName = $vwoInstance->getVariationName($campaignKey, 'Ashley');
      $this->assertEquals($variation, 'Control');
      $this->assertEquals($variationName, 'Control');
      $this->assertEquals($isGoalTracked, true);
  }

  public function testNullVariationAsOtherCampaignSatisfiesStorage()
  {
      $campaignKey = $this->settingsFileMEG['campaigns'][2]['key'];
      $vwoInstance = TestUtil::instantiateSdk($this->settingsFileMEG);

      $variationInfo = [
          'userId' => 'Ashley',
          'variationName' => 'Control',
          'campaignKey' => $campaignKey
      ];
      $vwoInstance->_userStorageObj = TestUtil::mockUserStorageInterface($this, $variationInfo);

      // called campaign satisfies the whitelisting
      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $variationName = $vwoInstance->getVariationName($campaignKey, 'Ashley');
      $isGoalTracked = $vwoInstance->track($campaignKey, 'Ashley', 'CUSTOM');
      $this->assertEquals($variation, 'Control');
      $this->assertEquals($variationName, 'Control');
      $this->assertEquals($isGoalTracked, true);

      $campaignKey = $this->settingsFileMEG['campaigns'][3]['key'];
      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $variationName = $vwoInstance->getVariationName($campaignKey, 'Ashley');
      $isGoalTracked = $vwoInstance->track($campaignKey, 'Ashley', 'CUSTOM');
      $this->assertEquals($variation, null);
      $this->assertEquals($variationName, null);
      $this->assertEquals($isGoalTracked, null);
  }

 public function testVariationForCalledCampaignInStorageAndOtherCampaignSatisfiesWhitelisting()
  {
      $campaignKey = $this->settingsFileMEG['campaigns'][2]['key'];
      $vwoInstance = TestUtil::instantiateSdk($this->settingsFileMEG, ['isUserStorage' => 1]);

      $options = [
          'variationTargetingVariables' => [
              'browser' => "chrome"
          ]
      ];

      $segmentPassed = [
          "or" => [
              [
                  "custom_variable" => [
                      'browser' => "chrome"
                  ]
              ]
          ]
      ];

      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $this->assertEquals($variation, 'Control');
      $vwoInstance->settings['campaigns'][3]['segments'] = $segmentPassed;
      $variation = $vwoInstance->activate($campaignKey, 'Ashley', $options);
      $this->assertEquals($variation, 'Control');
  }


  public function testWhenUserAlreadyPartOfCampaignAndNewCampaignAddedToGroup()
  {
      $campaignKey = $this->settingsFileMEG['campaigns'][2]['key'];
      $vwoInstance = TestUtil::instantiateSdk($this->settingsFileMEG);

      $variationInfo = [
          'userId' => 'Ashley',
          'variationName' => 'Control',
          'campaignKey' => $campaignKey
      ];
      $vwoInstance->_userStorageObj = TestUtil::mockUserStorageInterface($this, $variationInfo);

      // user is already a part of a campaign
      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $this->assertEquals($variation, 'Control');

      // new campaign is added to the group
      $vwoInstance->settings['campaignGroups'][164] = 2;
      $vwoInstance->settings['groups'][2]['campaigns'][] = 164;
      $variation = $vwoInstance->activate($this->settingsFileMEG['campaigns'][4]['key'], 'Ashley');
      $this->assertEquals($variation, null);
  }

  public function testWhenViewedCampaignRemovedFromGroup()
  {
      $campaignKey = $this->settingsFileMEG['campaigns'][2]['key'];
      $vwoInstance = TestUtil::instantiateSdk($this->settingsFileMEG, ['isUserStorage' => 1]);

      $variationInfo = [
          'userId' => 'Ashley',
          'variationName' => 'Control',
          'campaignKey' => $campaignKey
      ];
      $vwoInstance->_userStorageObj = TestUtil::mockUserStorageInterface($this, $variationInfo);

      // user is already a part of a campaign
      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $this->assertEquals($variation, 'Control');

      // old campaign is removed from the group
      $vwoInstance->settings['groups'][2]['campaigns'] = [163];

      // since user has already seen that campaign, they will continue to become part of that campaign
      $variation = $vwoInstance->activate($campaignKey, 'Ashley');
      $this->assertEquals($variation, 'Control');
  }*/
