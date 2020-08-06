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

package api

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	"github.com/wingify/vwo-go-sdk/pkg/utils"
)

const fileVWO = "vwo.go"

// VWOInstance is used to customize and construct an instance of VWO
type VWOInstance schema.VwoInstance

// VWOOption is used to provide custom instance configuration.
type VWOOption func(*VWOInstance)

// Init instantiates instance with the given options
func (vwo VWOInstance) Init(vwoOption ...VWOOption) (*VWOInstance, error) {
	// extracting options
	for _, option := range vwoOption {
		option(&vwo)
	}

	if vwo.Logger != nil {
		logger.Warning(constants.DebugMessageCustomLoggerFound)

		if !utils.ValidateLogger(vwo.Logger) {
			return &vwo, fmt.Errorf(constants.ErrorMessageCustomLoggerMisconfigured)
		}

		utils.LogMessage(vwo.Logger, constants.Debug, fileVWO, constants.DebugMessageCustomLoggerUsed)
	}

	if vwo.Logger == nil {
		logs := logger.Init(constants.SDKName, true, false, ioutil.Discard)
		logger.SetFlags(log.LstdFlags)
		message := fmt.Sprintf(constants.DebugMessageNoCustomLoggerFound)
		utils.LogMessage(logs, constants.Debug, fileVWO, message)
		vwo.Logger = logs
		defer logger.Close()
	}

	if !utils.ValidateStorage(vwo.UserStorage) {
		return &vwo, fmt.Errorf(constants.ErrorMessageInvalidLoggerStorage, "")
	}

	message := fmt.Sprintf(constants.DebugMessageDevelopmentMode+constants.DebugMessageSDKInitialized, vwo.IsDevelopmentMode)
	utils.LogMessage(vwo.Logger, constants.Debug, fileVWO, message)

	return &vwo, nil
}

// WithStorage sets user storage
func WithStorage(storage interface{}) VWOOption {
	return func(vwo *VWOInstance) {
		vwo.UserStorage = storage
	}
}

// WithLogger sets user custom logger
func WithLogger(logger interface{}) VWOOption {
	return func(vwo *VWOInstance) {
		vwo.Logger = logger
	}
}

// WithDevelopmentMode sets development mode true
func WithDevelopmentMode() VWOOption {
	return func(vwo *VWOInstance) {
		vwo.IsDevelopmentMode = true
	}
}

// WithGoalAttributes sets GoalTypeToTrack to the passed type and ShouldTrackReturningUser to false
func WithGoalAttributes(goalTypeToTrack interface{}, shouldTrackReturningUser interface{}) VWOOption {
	return func(vwo *VWOInstance) {
		vwo.GoalTypeToTrack = goalTypeToTrack
		vwo.ShouldTrackReturningUser = shouldTrackReturningUser
	}
}