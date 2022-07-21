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

package testdata

import "math/rand"

// GetRandomUser function returns a random user from the array of pre defined users
func GetRandomUser() string {
	var users = []string{
		"Ashley",
		"Bill",
		"Chris",
		"Dominic",
		"Emma",
		"Faizan",
		"Gimmy",
		"Harry",
		"Ian",
		"John",
		"King",
		"Lisa",
		"Mona",
		"Nina",
		"Olivia",
		"Pete",
		"Queen",
		"Robert",
		"Sarah",
		"Tierra",
		"Una",
		"Varun",
		"Will",
		"Xin",
		"You",
		"Zeba",
	}

	ind := rand.Intn(24)
	return users[ind]
}
