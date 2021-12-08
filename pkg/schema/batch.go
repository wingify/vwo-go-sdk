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

package schema

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/logger"
	"github.com/wingify/vwo-go-sdk/pkg/request"
)

type BatchEventQueue struct {
	AccountID           int
	impressions         []Impression
	Logger              interface{}
	ch                  chan Impression
	cancel              chan bool
	RequestTimeInterval int
	EventsPerRequest    int
	SDKKey              string
	IsDevelopmentMode   bool
	FlushCallBack       func(error, []map[string]interface{})
}

func (batch *BatchEventQueue) AddToBatch(message Impression, vwoInstance VwoInstance) {
	if batch.cancel == nil {
		batch.cancel = make(chan bool)
	}
	if batch.ch == nil {
		batch.ch = make(chan Impression)
		interval := time.Duration(batch.RequestTimeInterval) * time.Second

		go func() {
			timer := time.NewTimer(interval)
			open := true
			for open {
				select {
				case <-timer.C:
					batch.FlushBatch(vwoInstance)
					timer.Reset(interval)
				case data := <-batch.ch:
					batch.impressions = append(batch.impressions, data)
					timer.Reset(interval)
					if len(batch.impressions) >= batch.EventsPerRequest {
						batch.FlushBatch(vwoInstance)
					}
				case <-batch.cancel:
					open = false
					batch.FlushBatch(vwoInstance)
				}
			}
			timer.Stop()
		}()
	}
	batch.ch <- message
}

func (batch *BatchEventQueue) Flush() {
	batch.cancel <- true
	var vwoInstance VwoInstance
	batch.FlushBatch(vwoInstance)
	if batch.ch != nil {
		close(batch.ch)
		batch.ch = nil
	}
}

func (batch *BatchEventQueue) getBatchMinifiedPayload(impressions []Impression) []map[string]interface{} {
	eventTypeMapping := constants.EventTypeMapping
	events := make([]map[string]interface{}, 0)
	for _, impression := range impressions {
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

func (batch *BatchEventQueue) FlushBatch(vwoInstance VwoInstance) {
	defer batch.clear()
	if batch.IsDevelopmentMode || batch.impressions == nil || len(batch.impressions) == 0 {
		return
	}
	log := batch.Logger.(*logger.Logger)
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf(constants.ErrorMessageBatchFlushError, err))
		}
	}()

	headers := map[string]string{"Authorization": batch.SDKKey}
	url := constants.HTTPSProtocol + constants.EndPointsBaseURL + constants.BatchEndPoint
	body := map[string]interface{}{"ev": batch.getBatchMinifiedPayload(batch.impressions)}
	queryParams := map[string]string{
		"a":   strconv.Itoa(batch.AccountID),
		"sd":  constants.SDKName,
		"sv":  constants.SDKVersion,
		"env": batch.SDKKey,
	}
	for key, element := range GetUsageStatsObject(vwoInstance) {
		queryParams[key]= element
    }
	log.Debug(fmt.Sprintf(constants.DebugBeforeBatchFlush, strconv.Itoa(len(batch.impressions)), strconv.Itoa(batch.AccountID)))
	_, status, err := request.PostRequest(url, body, headers, queryParams)
	log.Debug(fmt.Sprintf(constants.DebugAfterBatchFlush, strconv.Itoa(len(batch.impressions))))
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
		batch.FlushCallBack(err, batch.getBatchMinifiedPayload(batch.impressions))
	}
}

func (batch *BatchEventQueue) clear() {
	batch.impressions = nil
}

func (batch *BatchEventQueue) GetBatchImpressions() []Impression {
	return batch.impressions
}
