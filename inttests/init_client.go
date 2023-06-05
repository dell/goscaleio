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
	"math"
	"os"

	"github.com/dell/goscaleio"
	"github.com/joho/godotenv"
)

const (
	envVarsFile        = "GOSCALEIO_TEST.env"
	mainEndpoint       = "GOSCALEIO_ENDPOINT"
	replicationEnpoint = "GOSCALEIO_ENDPOINT2"
)

// Global goscaleio Client instances for testing.
var (
	C  *goscaleio.Client
	C2 *goscaleio.Client
	GC *goscaleio.GatewayClient
)

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
			Endpoint: os.Getenv(mainEndpoint),
			Username: os.Getenv("GOSCALEIO_USERNAME"),
			Password: os.Getenv("GOSCALEIO_PASSWORD"),
			Insecure: os.Getenv("GOSCALEIO_INSECURE") == "true",
		})
		if err != nil {
			panic(fmt.Errorf("unable to login to VxFlexOS Gateway: %s", err.Error()))
		}
	}
}

// initClient2 initializes a second client for replication testing. Its use is optional.
// returns true if second client initialized
func initClient2() bool {
	var err error
	endpoint2 := os.Getenv(replicationEnpoint)
	if endpoint2 == "" {
		return false
	}

	C2, err = goscaleio.NewClientWithArgs(
		endpoint2,
		os.Getenv("GOSCALEIO_VERSION"),
		math.MaxInt64,
		os.Getenv("GOSCALEIO_INSECURE") == "true",
		os.Getenv("GOSCALEIO_INSECURE") == "true")

	if err != nil {
		panic(err)
	}

	if C2.GetToken() == "" {
		_, err := C2.Authenticate(&goscaleio.ConfigConnect{
			Endpoint: os.Getenv(replicationEnpoint),
			Username: os.Getenv("GOSCALEIO_USERNAME2"),
			Password: os.Getenv("GOSCALEIO_PASSWORD2"),
		})
		if err != nil {
			panic(fmt.Errorf("unable to login to VxFlexOS Gateway: %s", err.Error()))
		}
	}
	return true
}

func initGatewayClient() {
	err := godotenv.Load(envVarsFile)
	if err != nil {
		log.Printf("%s file not found.", envVarsFile)
	}

	GC, err = goscaleio.NewGateway(os.Getenv("GATEWAY_ENDPOINT"), os.Getenv("GATEWAY_USERNAME"), os.Getenv("GATEWAY_PASSWORD"), os.Getenv("GATEWAY_INSECURE") == "true", os.Getenv("GATEWAY_INSECURE") == "true")
	if err != nil {
		panic(err)
	}

}

func init() {
	initClient()
	initClient2()
	initGatewayClient()
}
