// Copyright 2020 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nsecurity

import (
	"errors"
	"strconv"
	"time"

	"nfgo.ga/nfgo/nerrors"
	"nfgo.ga/nfgo/nutil"
	"nfgo.ga/nfgo/nutil/jwt"
)

// AuthTicket -
type AuthTicket struct {
	ClientType string
	RequestID  string
	Token      string
	Subject    string
	Timestamp  string
	Signature  string
}

// VerifyToken -
func (a *AuthTicket) VerifyToken(jwtSecret string) error {

	// Check the jwt token
	claims, err := jwt.ValidateToken(jwtSecret, a.Token)
	if err != nil {
		return nerrors.ErrUnauthorized
	}

	// Check the subject in the token
	if claims.Subject != a.Subject {
		return errors.New("the ticket's subject is not equal with the subject int the token")
	}
	return nil
}

// VerifySignature -
func (a *AuthTicket) VerifySignature(signKey string) bool {
	expectSig := nutil.Sha256(signKey + a.Timestamp + a.Subject + a.RequestID)
	return expectSig == a.Signature
}

// VerifyTimeWindow - check IsoverTimeWindow clientTs milliseconds since January 1, 1970 UTC.
func (a *AuthTicket) VerifyTimeWindow(timeWindow time.Duration) error {
	clientTs, err := strconv.ParseInt(a.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	clientTime := time.Unix(0, clientTs*int64(time.Millisecond))
	serverTime := time.Now()

	var duration time.Duration
	if serverTime.After(clientTime) {
		duration = serverTime.Sub(clientTime)
	} else {
		duration = clientTime.Sub(serverTime)
	}
	if duration > timeWindow {
		return nerrors.ErrUnauthorized
	}
	return nil
}
