// Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goscaleio

import (
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

// This test can be checked when NewGateway() function is fixed
func TestUploadCompliance(t *testing.T) {
	t.Skip("Skipping this test case")
	type testCase struct {
		ucParam  *types.UploadComplianceParam
		expected error
	}
	cases := []testCase{
		{
			ucParam: &types.UploadComplianceParam{
				SourceLocation: "https://10.10.10.1/artifactory/pfmp20/RCM/Denver/RCMs/SoftwareOnly/PowerFlex_Software_4.5.0.0_287_r1.zip",
			},
			expected: nil,
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc

		t.Run("", func(_ *testing.T) {
			GC, err := NewGateway(svr.URL, "", "", true, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, errFs = GC.UploadCompliance(tc.ucParam)
			if errFs != nil {
				if tc.expected == nil {
					t.Errorf("Uploading Compliance File did not work as expected, \n\tgot: %s \n\twant: %v", errFs, tc.expected)
				} else {
					if errFs.Error() != tc.expected.Error() {
						t.Errorf("Uploading Compliance File did not work as expected, \n\tgot: %s \n\twant: %s", errFs, tc.expected)
					}
				}
			}
		})
	}
}
