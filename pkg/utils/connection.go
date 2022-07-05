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

package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
)

// GetRequest function to do a get call
func GetRequest(url string) (string, error) {
	/*
		Args:
			url: URL needed

		Return:
			string: stringified content recieved
			error: error encountered while Get rewuest, nil if no error
	*/
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf(constants.ErrorMessageURLNotFound, "", err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf(constants.ErrorMessageResponseNotParsed, "", url)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf(constants.ErrorMessageCouldNotGetURL, "", url)
	}
	return string(body), nil
}
