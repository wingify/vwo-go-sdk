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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
)

func testGetUrlWithOutCollectionPrefix(t *testing.T) {
	var settingsFile schema.SettingsFile
	BaseUrlAchieved := GetDataLocation(settingsFile)
	BaseUrl := constants.BaseURL
	assert.Equal(t, BaseUrlAchieved, BaseUrl, "BaseUrl not matched when collection prefix is empty")
}

func testGetUrlWithCollectionPrefix(t *testing.T) {
	var settingsFile schema.SettingsFile
	settingsFile.CollectionPrefix = "eu"
	BaseUrlAchieved := GetDataLocation(settingsFile)
	BaseUrl := constants.BaseURL + "/" + settingsFile.CollectionPrefix
	assert.Equal(t, BaseUrlAchieved, BaseUrl, "BaseUrl not matched when collection prefix is set to eu")
}
