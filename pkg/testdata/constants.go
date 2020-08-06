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

package testdata

//Constants for testing
const (
	//API
	NonExistingCampaign     = "notPresent"
	NotRunningCampaign      = "CAMPAIGN_1"
	FeatureRolloutCampaign  = "CAMPAIGN_2"
	GetFeatureDummyUser     = "Ashley"
	GetFeatureDummyVariable = "STRING_VARIABLE"
	ValidTagKey             = "demoTagKey"
	ValidTagValue           = "demoTagVal"
	InvalidTagKey           = "demoTagKey-Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis id tellus quis massa iaculis interdum. Morbi rutrum, lacus ac egestas lobortis, lectus lectus mollis sem, eget vehicula justo velit ut erat. Mauris ac ligula id nulla laoreet fringilla non at purus. Quisque eu risus quis mi convallis sagittis. Aliquam luctus posuere mollis. Nullam rhoncus mauris a lorem sagittis efficitur. Nulla quis risus sit amet tellus bibendum facilisis. Aliquam erat volutpat.In aliquam imperdiet nulla, sed consequat ex pharetra eget. Mauris eget vestibulum nunc. Morbi sem lectus, elementum sit amet laoreet at, euismod a purus. Aliquam ut tristique neque, tempor aliquet nisl. Aenean vestibulum lectus ut semper fringilla. Phasellus accumsan lorem at risus laoreet, non molestie est egestas. Fusce ac tellus vel nulla mollis auctor. Praesent ac laoreet lorem.Proin bibendum sodales nulla eget consectetur. Etiam auctor non lacus ac venenatis. Maecenas a magna dolor. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; In id ornare nunc, vel sodales purus. Integer ultricies dui at tortor bibendum facilisis. Vestibulum mollis porttitor ligula. Fusce odio tortor, imperdiet vel lectus id, rhoncus facilisis tortor. Ut sagittis purus non sapien condimentum, vitae iaculis ligula pharetra. Donec in metus id libero pellentesque mattis sed sed metus. Maecenas a nisi ut risus volutpat posuere. Nunc id semper quam, ac vehicula lacus. Aliquam erat volutpat.Aliquam cursus lacinia odio non pretium. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lectus ex, consectetur at augue pretium, iaculis cursus lacus. Aliquam nec porta erat. Aliquam blandit lobortis sapien, vitae maximus."
	InvalidTagValue         = "demoTagVal-Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis id tellus quis massa iaculis interdum. Morbi rutrum, lacus ac egestas lobortis, lectus lectus mollis sem, eget vehicula justo velit ut erat. Mauris ac ligula id nulla laoreet fringilla non at purus. Quisque eu risus quis mi convallis sagittis. Aliquam luctus posuere mollis. Nullam rhoncus mauris a lorem sagittis efficitur. Nulla quis risus sit amet tellus bibendum facilisis. Aliquam erat volutpat.In aliquam imperdiet nulla, sed consequat ex pharetra eget. Mauris eget vestibulum nunc. Morbi sem lectus, elementum sit amet laoreet at, euismod a purus. Aliquam ut tristique neque, tempor aliquet nisl. Aenean vestibulum lectus ut semper fringilla. Phasellus accumsan lorem at risus laoreet, non molestie est egestas. Fusce ac tellus vel nulla mollis auctor. Praesent ac laoreet lorem.Proin bibendum sodales nulla eget consectetur. Etiam auctor non lacus ac venenatis. Maecenas a magna dolor. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; In id ornare nunc, vel sodales purus. Integer ultricies dui at tortor bibendum facilisis. Vestibulum mollis porttitor ligula. Fusce odio tortor, imperdiet vel lectus id, rhoncus facilisis tortor. Ut sagittis purus non sapien condimentum, vitae iaculis ligula pharetra. Donec in metus id libero pellentesque mattis sed sed metus. Maecenas a nisi ut risus volutpat posuere. Nunc id semper quam, ac vehicula lacus. Aliquam erat volutpat.Aliquam cursus lacinia odio non pretium. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam lectus ex, consectetur at augue pretium, iaculis cursus lacus. Aliquam nec porta erat. Aliquam blandit lobortis sapien, vitae maximus."

	//CORE
	ValidBucketValue         = 2345
	InvalidBucketValue       = 0
	ValidUser                = "DummyUser"
	UserIsFeatureEnabled     = "John"
	InvalidUser              = "UserInvalid"
	InvalidOperator          = "InvalidOperator"
	DummyVariation           = "DummyVariation"
	DummyGoal                = "DummyGoal"
	TempUser                 = "TempUser"
	ValidVariationControl    = "Control"
	ValidVariationVariation2 = "Variation-2"

	//SERVICE
	DummyAccountID      = "accountID"
	DummySDKKey         = "SDKKey"
	InvalidAccountID    = ""
	InvalidSDKKey       = ""
	ValidSettingsFile   = "../testdata/dummy_settings_file.json"
	InvalidSettingsFile = "../testdata/invalid_settings_file.json"
	EmptySettingsFile   = "../testdata/settings_file.json"

	//UTILS
	ValidCampaignKey       = "AB_T_50_W_50_50"
	InvalidCampaignKey     = "notAvailable"
	ValidVariationName     = "Control"
	InvalidVariationName   = "NoVaritionInCampaign"
	IncorrectNariationName = "Variation-2"
	ValidGoal              = "GOAL_2"
	InvalidGoal            = "NotAvailable"
	IncorrectURL1          = "https://jsonplaceholder.typicode.com/todos/1"
	IncorrectURL2          = "https.com"
	IncorrectURL3          = "https://github.com/wrong-endpoint/abc"
	ValidVariableKey1      = "INTEGER_VARIABLE"
	ValidVariableKey2      = "STRING_VARIABLE"
	InvalidVariableKey     = "STRINGER_VARIABLE"
	TestKey1               = "testKey"
	TestKey2               = "test Key"
	TestValue1             = "testVal"
	TestValue2             = "test Val"
)
