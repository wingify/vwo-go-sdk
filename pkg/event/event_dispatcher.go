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

package event

import (
	"fmt"
	"strconv"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

const eventDispatcher = "eventDispatcher.go"

// Dispatch function dispatches the event represented by the impression object to our servers
func Dispatch(vwoInstance schema.VwoInstance, impression schema.Impression) {
	/*
		Args:
			impression: impression to be dispatched
	*/

	if !vwoInstance.IsDevelopmentMode {
		URL := impression.URL + "?" +
			"random=" + strconv.FormatFloat(float64(impression.Random), 'f', -1, 64) +
			"&sdk=" + impression.Sdk +
			"&sdk-v=" + impression.SdkV +
			"&ap=" + impression.Ap +
			"&sId=" + impression.SID +
			"&u=" + impression.U +
			"&account_id=" + strconv.Itoa(impression.AccountID) +
			"&uId=" + impression.UID +
			"&tags=" + impression.Tags

		if vwoInstance.API != "Push" {
			URL = URL + "&ed=" + impression.ED +
				"&experiment_id=" + strconv.Itoa(impression.ExperimentID) +
				"&combination=" + strconv.Itoa(impression.Combination)
		}

		_, err := utils.GetRequest(URL)

		if err != nil {
			message := fmt.Sprintf(constants.ErrorMessageImpressionFailed, vwoInstance.API, err)
			utils.LogMessage(vwoInstance.Logger, constants.Error, eventDispatcher, message)
		} else {
			if vwoInstance.API == "Push" {
				message := fmt.Sprintf(constants.InfoMessageImpressionSuccess, vwoInstance.API, "Push", URL)
				utils.LogMessage(vwoInstance.Logger, constants.Info, eventDispatcher, message)
				} else {
				message := fmt.Sprintf(constants.InfoMessageImpressionSuccess, vwoInstance.API, "Tracking User", URL)
				utils.LogMessage(vwoInstance.Logger, constants.Info, eventDispatcher, message)
			}
		}
	}
}

// DispatchTrackingGoal function dispatches the event with goal tracking represented by
// the impression object to our servers
func DispatchTrackingGoal(vwoInstance schema.VwoInstance, goalType string, impression schema.Impression) {
	/*
		Args:
			impression: impression to be dispatched
	*/

	if !vwoInstance.IsDevelopmentMode {
		URL := impression.URL + "?" +
			"random=" + strconv.FormatFloat(float64(impression.Random), 'f', -1, 64) +
			"&sdk=" + impression.Sdk +
			"&sdk-v=" + impression.SdkV +
			"&ap=" + impression.Ap +
			"&sId=" + impression.SID +
			"&u=" + impression.U +
			"&account_id=" + strconv.Itoa(impression.AccountID) +
			"&uId=" + impression.UID +
			"&experiment_id=" + strconv.Itoa(impression.ExperimentID) +
			"&combination=" + strconv.Itoa(impression.Combination) +
			"&goal_id=" + strconv.Itoa(impression.GoalID)

		if goalType == constants.GoalTypeRevenue {
			URL = URL + "&r=" + impression.R
		}

		_, err := utils.GetRequest(URL)

		if err != nil {
			message := fmt.Sprintf(constants.ErrorMessageImpressionFailed, vwoInstance.API, err)
			utils.LogMessage(vwoInstance.Logger, constants.Error, eventDispatcher, message)
		} else {
			message := fmt.Sprintf(constants.InfoMessageImpressionSuccess, vwoInstance.API, "Tracking Goal", URL)
			utils.LogMessage(vwoInstance.Logger, constants.Info, eventDispatcher, message)
		}
	}
}
