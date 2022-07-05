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
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {
	url := testdata.IncorrectURL1
	content, err := GetRequest(url)
	assert.Nil(t, err, "Could not make the Get Request")
	assert.NotEmpty(t, content, "Recieved no content")

	url = testdata.IncorrectURL2
	content, err = GetRequest(url)
	assert.NotNil(t, err, "Could not make the Get Request")
	assert.Empty(t, content, "Recieved no content")

	url = testdata.IncorrectURL3
	content, err = GetRequest(url)
	assert.NotNil(t, err, "Could not make the Get Request")
	assert.Empty(t, content, "Recieved no content")
}
