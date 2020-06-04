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

package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wingify/vwo-go-sdk/pkg/constants"
	"github.com/wingify/vwo-go-sdk/pkg/schema"
	guuid "github.com/google/uuid"
	suuid "github.com/satori/go.uuid"
)

const uuid = "uuid.go"

// generateFor generates desired UUID
func generateFor(vwoInstance schema.VwoInstance, userID string, accountID int) string {
	/*
		Args:
		    userID : User identifier
		    accountID : Account identifier

		Returns:
			string : Desired Uuid
	*/
	NameSpaceURL, _ := guuid.Parse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	VWONamespace := suuid.NewV5(suuid.UUID(NameSpaceURL), "https://vwo.com")
	userIDNamespace := suuid.NewV5(VWONamespace, strconv.Itoa(accountID))
	uuidForAccountUserID := suuid.NewV5(userIDNamespace, userID)
	desiredUUID := strings.ToUpper(strings.Replace(uuidForAccountUserID.String(), "-", "", -1))

	message := fmt.Sprintf(constants.DebugMessageUUIDForUser, vwoInstance.API, userID, accountID, desiredUUID)
	LogMessage(vwoInstance.Logger, constants.Debug, uuid, message)

	return desiredUUID
}
