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

package api

import (
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	vwoInstance, err := getInstance("./testdata/testdata.json")
	assert.Nil(t, err, "error fetching instance")
	userID := testdata.GetRandomUser()

	tagKey := ""
	tagValue := ""
	pushed := vwoInstance.Push(tagKey, tagValue, userID)
	assert.False(t, pushed, "Invalid params")

	tagKey = testdata.ValidTagKey
	tagValue = testdata.ValidTagValue
	pushed = vwoInstance.Push(tagKey, tagValue, userID)
	assert.True(t, pushed, "Unable to Push")

	tagKey = testdata.ValidTagKey
	tagValue = testdata.InvalidTagValue
	pushed = vwoInstance.Push(tagKey, tagValue, userID)
	assert.False(t, pushed, "Unable to Push")

	tagKey = testdata.InvalidTagKey
	tagValue = testdata.ValidTagValue
	pushed = vwoInstance.Push(tagKey, tagValue, userID)
	assert.False(t, pushed, "Unable to Push")
}
