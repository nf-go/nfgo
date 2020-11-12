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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthTicketVerifySignature(t *testing.T) {
	signKey := "5f0fa825de41a7d3fd000002"
	ts := "1594800736455"
	subject := "1"
	requestID := "xxx"
	token := "eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MjYzOTc2MDUuNjIyNTY0LCJzdWIiOiIxIn0.4S2w33nYU5rzOjPePP5t4YUnNVmRGJKlUxu8_ioogoU"
	sig := "d3dee5e6260b6a2ce59bc7a6a1a14d024756770cc7a6b5ec9a7ee4a9c5dc82d7"

	ticket := &AuthTicket{
		Timestamp: ts,
		Subject:   subject,
		RequestID: requestID,
		Signature: sig,
		Token:     token,
	}

	assert.True(t, ticket.VerifySignature(signKey), "the signature should be valid")

	assert.NoError(t, ticket.VerifyToken("VoCra8#GEBAbRl*+vos9UF@??gi8Oy"))

}
