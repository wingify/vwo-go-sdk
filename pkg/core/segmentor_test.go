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

package core

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/testdata"
	"github.com/stretchr/testify/assert"
)

type SegmentorTestCase struct {
	DSL                         map[string]interface{} `json:"dsl"`
	Expected                    bool                   `json:"expectation"`
	CustomVariable              map[string]interface{} `json:"custom_variables"`
	VariationTargetingVariables map[string]interface{} `json:"variation_targeting_variables"`
}

func TestSegmentEvaluator(t *testing.T) {
	var TestData map[string]map[string]SegmentorTestCase
	data, err := ioutil.ReadFile("../testdata/test_segment.json")
	if err != nil {
		logger.Info("Error: " + err.Error())
	}

	if err = json.Unmarshal(data, &TestData); err != nil {
		logger.Info("Error: " + err.Error())
	}

	for parent, v := range TestData {
		for child, value := range v {
			var actual bool
			if value.CustomVariable != nil {
				actual = SegmentEvaluator(value.DSL, value.CustomVariable)
			} else {
				actual = SegmentEvaluator(value.DSL, value.VariationTargetingVariables)
			}
			expected := value.Expected
			assert.Equal(t, expected, actual, parent+" "+child)
		}
	}

	// CORNER CASES

	vwoInstance := testdata.GetInstanceWithCustomSettings("SettingsFile4")
	segments := vwoInstance.SettingsFile.Campaigns[0].Segments
	actual := SegmentEvaluator(segments, nil)
	assert.True(t, actual, "No Case for operator hit")
}

func TestEvaluate(t *testing.T) {
	operator := testdata.InvalidOperator
	var res []bool
	actual := evaluate(operator, res)
	assert.False(t, actual, "No Case for operator hit")
}
