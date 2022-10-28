// Copyright © 2021 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package inttests

import (
	"fmt"
	"log"
	"os"

	"github.com/dell/goscaleio"
	"github.com/joho/godotenv"
)

const envVarsFile = "GOSCALEIO_TEST.env"

// C is global goscaleio Client instance for testing
var C *goscaleio.Client

func initClient() {
	err := godotenv.Load(envVarsFile)
	if err != nil {
		log.Printf("%s file not found.", envVarsFile)
	}

	C, err = goscaleio.NewClient()
	if err != nil {
		panic(err)
	}

	if C.GetToken() == "" {
		_, err := C.Authenticate(&goscaleio.ConfigConnect{
			Endpoint: os.Getenv("GOSCALEIO_ENDPOINT"),
			Username: os.Getenv("GOSCALEIO_USERNAME"),
			Password: os.Getenv("GOSCALEIO_PASSWORD"),
		})
		if err != nil {
			panic(fmt.Errorf("unable to login to VxFlexOS Gateway: %s", err.Error()))
		}
	}
}

func init() {
	initClient()
}
