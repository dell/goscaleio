# Copyright © 2021 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#      http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

# The tests rely upon a few environment variables
# GOSCALEIO_SDC_GUID: the locally installed SDC GUID
# GOSCALEIO_NUMBER_SYSTEMS: The number of connected MDM clusters/systems
# GOSCALEIO_SYSTEMID: the system (MDM cluster) ID of a connected clusters/system

# exit on failure
set -e

DRVCFG="/opt/emc/scaleio/sdc/bin/drv_cfg"
SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
cd "${SCRIPTDIR}"

# if the sdc is installed and the rv_cfg executable is available
if [ -f "${DRVCFG}" ]; then
    # check if the SDC kernel module is actually loaded
    SCINI=$(lsmod | grep scini | wc -l)
    if [ "${SCINI}" != "0" ]; then 
        # Get the SDC GUID
        GUID=$("${DRVCFG}" --query_guid)
        export GOSCALEIO_SDC_GUID="${GUID}"

        # get the number of systems connected
        COUNT=$("${DRVCFG}" --query_mdm | grep ^MDM-ID | wc -l)
        export GOSCALEIO_NUMBER_SYSTEMS="${COUNT}"

        # get the system ID (use the last to force iteration through all systems)
        MDM=$("${DRVCFG}" --query_mdm | grep ^MDM-ID | tail -n 1 | awk '{print $2}')
        export GOSCALEIO_SYSTEMID="${MDM}"
    else
        echo "The SDC is installed but the kernel module is not loaded"
        echo "All of the drv_cfg tests cannot be run"
    fi
else
    echo "The SDC is not installed"
    echo "All of the drv_cfg tests cannot be run"
fi

# Run the integration tests
echo "Starting tests"
go test -v -coverprofile=c.out -coverpkg github.com/dell/goscaleio .
