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

package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/request"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
)

func AddToBatch(message schema.Impression, vwoInstance schema.VwoInstance, batch *schema.BatchEventQueue) {
	if batch.Cancel == nil {
		batch.Cancel = make(chan bool)
	}
	if batch.Ch == nil {
		batch.Ch = make(chan schema.Impression)
		interval := time.Duration(batch.RequestTimeInterval) * time.Second

		go func() {
			timer := time.NewTimer(interval)
			open := true
			for open {
				select {
				case <-timer.C:
					FlushBatch(vwoInstance, batch)
					timer.Reset(interval)
				case data := <-batch.Ch:
					batch.Impressions = append(batch.Impressions, data)
					timer.Reset(interval)
					if len(batch.Impressions) >= batch.EventsPerRequest {
						FlushBatch(vwoInstance, batch)
					}
				case <-batch.Cancel:
					open = false
					FlushBatch(vwoInstance, batch)
				}
			}
			timer.Stop()
		}()
	}
	batch.Ch <- message
}

func Flush(batch *schema.BatchEventQueue) {
	batch.Cancel <- true
	var vwoInstance schema.VwoInstance
	FlushBatch(vwoInstance, batch)
	if batch.Ch != nil {
		close(batch.Ch)
		batch.Ch = nil
	}
}

func getBatchMinifiedPayload(batch *schema.BatchEventQueue) []map[string]interface{} {
	eventTypeMapping := constants.EventTypeMapping
	events := make([]map[string]interface{}, 0)
	for _, impression := range batch.Impressions {
		event := make(map[string]interface{}, 0)
		sessionId, _ := strconv.Atoi(impression.SID)
		event["u"] = impression.U
		event["sId"] = sessionId
		eventName := impression.EventType
		event["eT"] = eventTypeMapping[eventName]
		if eventName == constants.EventsTrackGoal || eventName == constants.EventsTrackUser {
			event["e"] = impression.ExperimentID
			event["c"] = impression.Combination
		}

		if eventName == constants.EventsTrackGoal {
			event["g"] = impression.GoalID
			if len(impression.R) > 0 {
				event["r"] = impression.R
			}
		}

		if eventName == constants.EventsPush {
			event["t"] = impression.Tags
		}
		event["env"] = batch.SDKKey
		events = append(events, event)
	}
	return events
}

func FlushBatch(vwoInstance schema.VwoInstance, batch *schema.BatchEventQueue) {
	defer clear(batch)
	if batch.IsDevelopmentMode || batch.Impressions == nil || len(batch.Impressions) == 0 {
		return
	}
	log := batch.Logger.(*logger.Logger)
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf(constants.ErrorMessageBatchFlushError, err))
		}
	}()

	headers := map[string]string{"Authorization": batch.SDKKey}
	UpdatedBaseURL := GetDataLocation(vwoInstance.SettingsFile)
	url := constants.HTTPSProtocol + UpdatedBaseURL + constants.BatchEndPoint
	body := map[string]interface{}{"ev": getBatchMinifiedPayload(batch)}
	queryParams := map[string]string{
		"a":   strconv.Itoa(batch.AccountID),
		"sd":  constants.SDKName,
		"sv":  constants.SDKVersion,
		"env": batch.SDKKey,
	}
	for key, element := range schema.GetUsageStatsObject(vwoInstance) {
		queryParams[key] = element
	}
	log.Debug(fmt.Sprintf(constants.DebugBeforeBatchFlush, strconv.Itoa(len(batch.Impressions)), strconv.Itoa(batch.AccountID)))
	_, status, err := request.PostRequest(url, body, headers, queryParams)
	log.Debug(fmt.Sprintf(constants.DebugAfterBatchFlush, strconv.Itoa(len(batch.Impressions))))
	if status == http.StatusOK {
		log.Info(fmt.Sprintf(constants.InfoBatchImpressionSuccess, constants.BatchEndPoint))
	} else {
		var errString string
		if status == http.StatusRequestEntityTooLarge {
			errString = fmt.Sprintf(constants.DebugMessagePayloadTooLarge, constants.BatchEndPoint, batch.EventsPerRequest)
		} else if status == http.StatusBadRequest {
			errString = fmt.Sprintf(constants.ErrorMessageBatchImpressionFailed, constants.BatchEndPoint, strconv.Itoa(status))
		}
		log.Debug(errString)
		err = fmt.Errorf(errString)
	}

	if batch.FlushCallBack != nil {
		batch.FlushCallBack(err, getBatchMinifiedPayload(batch))
	}
}

func clear(batch *schema.BatchEventQueue) {
	batch.Impressions = nil
}

func GetBatchImpressions(batch *schema.BatchEventQueue) []schema.Impression {
	return batch.Impressions
}
