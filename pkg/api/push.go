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
	"github.com/wingify/vwo-go-sdk/pkg/event"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

const push = "push.go"

// Push function
/*
This API method: Pushes the key-value tag pair for a particular user
1. Validates the arguments being passed
2. Checks the length of tag Key and Value
3. Sends a call to VWO push api
*/
func (vwo *VWOInstance) Push(tagKey, tagValue, userID string) bool {
	/*
		Args:
			tagKey: Key of the corresponding tag
			tagValue: Value of the corresponding tag
			userID: Unique identification of user

		Returns:
			bool: true if the push api call is done, else false
	*/
	vwoInstance := schema.VwoInstance{
		SettingsFile:      vwo.SettingsFile,
		UserStorage:       vwo.UserStorage,
		Logger:            vwo.Logger,
		IsDevelopmentMode: vwo.IsDevelopmentMode,
		UserID:            userID,
		API:               "Push",
		Integrations:      vwo.Integrations,
	}

	if !utils.ValidatePush(tagKey, tagValue, userID) {
		message := fmt.Sprintf(constants.ErrorMessagePushAPIMissingParams, vwoInstance.API)
		utils.LogMessage(vwo.Logger, constants.Error, push, message)
		return false
	}

	if len(tagKey) > constants.PushAPITagKeyLength {
		message := fmt.Sprintf(constants.ErrorMessageTagKeyLengthExceeded, vwoInstance.API, tagKey, userID)
		utils.LogMessage(vwo.Logger, constants.Error, push, message)
		return false
	}
	if len(tagValue) > constants.PushAPITagValueLength {
		message := fmt.Sprintf(constants.ErrorMessageTagValueLengthExceeded, vwoInstance.API, tagValue, tagKey, userID)
		utils.LogMessage(vwo.Logger, constants.Error, push, message)
		return false
	}

	impression := utils.CreateImpressionForPush(vwoInstance, tagKey, tagValue, userID)
	if vwo.IsBatchingEnabled {
		vwo.AddToBatch(impression)
	} else {
		go event.Dispatch(vwoInstance, impression)
	}

	message := fmt.Sprintf(constants.InfoMessageMainKeysForPushAPI, vwoInstance.API, vwoInstance.SettingsFile.AccountID, userID, impression.U, impression.URL)
	utils.LogMessage(vwo.Logger, constants.Info, push, message)

	return true
}
