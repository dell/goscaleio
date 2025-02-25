package goscaleio

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRfCacheParams_MarshalJSON(t *testing.T) {
	// Test case: successful execution
	params := PDRfCacheParams{
		RfCacheOperationalMode: PDRfCacheOpMode("test"),
		RfCachePageSizeKb:      100,
		RfCacheMaxIoSizeKb:     200,
	}

	expectedJSON := map[string]string{
		"rfcacheOperationMode": "test",
		"pageSizeKb":           "100",
		"maxIOSizeKb":          "200",
	}

	obj, err := params.MarshalJSON()
	assert.NoError(t, err)

	// unmarshall json obj back to map[string]string
	recoveredJSON := make(map[string]string)
	err = json.Unmarshal(obj, &recoveredJSON)
	assert.NoError(t, err)
	assert.Equal(t, expectedJSON, recoveredJSON)
}

func TestParamsFromPD(t *testing.T) {
	dp := &ProtectionDomain{
		RfCacheOperationalMode:                          PDRfCacheOpMode("test"),
		RfCachePageSizeKb:                               100,
		RfCacheMaxIoSizeKb:                              200,
		RebuildNetworkThrottlingInKbps:                  111,
		RebalanceNetworkThrottlingInKbps:                222,
		VTreeMigrationNetworkThrottlingInKbps:           333,
		ProtectedMaintenanceModeNetworkThrottlingInKbps: 444,
		OverallIoNetworkThrottlingInKbps:                555,
	}

	// Test case: get RfCache params
	rfcParams := dp.GetRfCacheParams()
	assert.NotNil(t, rfcParams)
	assert.Equal(t, "test", string(rfcParams.RfCacheOperationalMode))
	assert.Equal(t, 100, rfcParams.RfCachePageSizeKb)
	assert.Equal(t, 200, rfcParams.RfCacheMaxIoSizeKb)

	// Test case: get SdsNetworkLimit params
	nwlParams := dp.GetNwLimitParams()
	assert.NotNil(t, nwlParams)
	assert.NotNil(t, nwlParams.RebuildNetworkThrottlingInKbps)
	assert.NotNil(t, nwlParams.RebalanceNetworkThrottlingInKbps)
	assert.NotNil(t, nwlParams.VTreeMigrationNetworkThrottlingInKbps)
	assert.NotNil(t, nwlParams.ProtectedMaintenanceModeNetworkThrottlingInKbps)
	assert.NotNil(t, nwlParams.OverallIoNetworkThrottlingInKbps)
	assert.Equal(t, 111, *nwlParams.RebuildNetworkThrottlingInKbps)
	assert.Equal(t, 222, *nwlParams.RebalanceNetworkThrottlingInKbps)
	assert.Equal(t, 333, *nwlParams.VTreeMigrationNetworkThrottlingInKbps)
	assert.Equal(t, 444, *nwlParams.ProtectedMaintenanceModeNetworkThrottlingInKbps)
	assert.Equal(t, 555, *nwlParams.OverallIoNetworkThrottlingInKbps)
}

func TestSdsNetworkLimitParams_MarshalJSON(t *testing.T) {
	// Test case: successful execution
	RebuildNetworkThrottlingInKbps := 123
	RebalanceNetworkThrottlingInKbps := 111
	VTreeMigrationNetworkThrottlingInKbps := 222
	ProtectedMaintenanceModeNetworkThrottlingInKbps := 192
	OverallIoNetworkThrottlingInKbps := 32

	params := &SdsNetworkLimitParams{
		RebuildNetworkThrottlingInKbps:                  &RebuildNetworkThrottlingInKbps,
		RebalanceNetworkThrottlingInKbps:                &RebalanceNetworkThrottlingInKbps,
		VTreeMigrationNetworkThrottlingInKbps:           &VTreeMigrationNetworkThrottlingInKbps,
		ProtectedMaintenanceModeNetworkThrottlingInKbps: &ProtectedMaintenanceModeNetworkThrottlingInKbps,
		OverallIoNetworkThrottlingInKbps:                &OverallIoNetworkThrottlingInKbps,
	}

	expectedJSON := map[string]string{
		"rebuildLimitInKbps":                  "123",
		"rebalanceLimitInKbps":                "111",
		"vtreeMigrationLimitInKbps":           "222",
		"protectedMaintenanceModeLimitInKbps": "192",
		"overallLimitInKbps":                  "32",
	}

	obj, err := params.MarshalJSON()
	assert.NoError(t, err)

	// unmarshall json obj back to map[string]string
	recoveredJSON := make(map[string]string)
	err = json.Unmarshal(obj, &recoveredJSON)
	assert.NoError(t, err)
	assert.Equal(t, expectedJSON, recoveredJSON)
}
