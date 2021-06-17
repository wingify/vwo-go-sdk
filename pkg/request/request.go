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

package request

import (
	"bytes"
	"encoding/json"
	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func PostRequest(uri string, body interface{}, headers map[string]string, queryParams map[string]string) ([]byte, int, error) {
	u, _ := url.Parse(uri)
	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(constants.HttpPostMethod, u.String(), bytes.NewBuffer(jsonBody))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	response, err := Client.Do(req)
	defer response.Body.Close()
	if err == nil {
		responseBody, err := ioutil.ReadAll(response.Body)
		return responseBody, response.StatusCode, err
	}
	return nil, response.StatusCode, err
}
