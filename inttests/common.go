/*
 *
 * Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package inttests

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

const (
	invalidIdentifier = "invalidIdentifier"
	testPrefix        = "inttest"
	letters           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	incrementingNumber = 0
	defaultVolumeSize  = (8 * 1024)
)

func checkAPIErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[randomInt.Int64()]
	}
	return string(b)
}

func getUniqueName() string {
	name := fmt.Sprintf("%s-%d", testPrefix, incrementingNumber)
	incrementingNumber++
	return name
}
