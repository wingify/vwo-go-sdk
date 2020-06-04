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

package vwo

import (
	"github.com/wingify/vwo-go-sdk/pkg/api"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/service"
)

// VWOInstance struct to store params
type VWOInstance schema.VwoInstance

const fileVWO = "vwo.go"

// Launch function to intialise sdk
func Launch(settingsFile schema.SettingsFile, vwoOption ...api.VWOOption) (*api.VWOInstance, error) {
	vwo := &api.VWOInstance{
		SettingsFile: settingsFile,
	}
	return vwo.Init(vwoOption...)
}

// GetSettingsFile function to fetch and parse settingsfile
func GetSettingsFile(accountID, SDKKey string) schema.SettingsFile {
	/*
		Args:
			accountID: Config account ID
			SDKKey: Config SDK Key

		Returns:
			schema.SettingsFile: settings file fetched
	*/
	settingsFileManager := service.SettingsFileManager{}
	if err := settingsFileManager.FetchSettingsFile(accountID, SDKKey); err != nil {
		logger.Warningf(fileVWO+" : "+constants.ErrorMessageCannotProcessSettingsFile, "", err.Error())
	}
	settingsFileManager.Process()
	logger.Warningf(fileVWO+" : "+constants.DebugMessageSettingsFileProcessed, "")
	return settingsFileManager.GetSettingsFile()
}

func SetLogLevel(lvl int) {
	logger.SetLogLevel(lvl)
}
