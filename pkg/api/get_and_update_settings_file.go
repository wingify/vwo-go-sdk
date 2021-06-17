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

package api

import (
	"fmt"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/service"
	"strconv"
)

func (vwoInstance *VWOInstance) GetAndUpdateSettingsFile() {
	log := vwoInstance.Logger.(*logger.Logger)
	settingsFileManager := service.SettingsFileManager{}
	accountId := vwoInstance.SettingsFile.AccountID
	sdkKey := vwoInstance.SettingsFile.SDKKey
	err := settingsFileManager.FetchSettingsFile(strconv.Itoa(accountId), sdkKey, true)
	if err != nil {
		log.Error(fmt.Sprintf(constants.ErrorMessageSettingsFileUpdateFailed, accountId, err))
		return
	}
	settingsFileManager.Process()
	vwoInstance.SettingsFile = settingsFileManager.GetSettingsFile()
	log.Info(fmt.Sprintf(constants.InfoSDKInstanceUpdated, accountId))
}
