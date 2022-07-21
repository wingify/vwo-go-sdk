/*
 * Copyright 2020-2022 Wingify Software Pvt. Ltd.
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

package schema

type VwoInstance struct {
	SettingsFile             SettingsFile
	UserStorage              interface{}
	Logger                   interface{}
	IsDevelopmentMode        bool
	UserID                   string
	Campaign                 Campaign
	API                      string
	GoalTypeToTrack          interface{}
	ShouldTrackReturningUser interface{}
	BatchEventQueue          BatchEventQueue
	IsBatchingEnabled        bool
	Integrations             Integrations
}
