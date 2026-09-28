/*
 *
 * Copyright © 2021-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package goscaleio

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetaData(t *testing.T) {
	vp := &VolumeParam{}
	assert.NotNil(t, vp.MetaData())
}

func TestIntString_MarshalJSON(t *testing.T) {
	is := IntString(123)
	obj, err := is.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"123"`), obj)
}

func TestGetBoolType(t *testing.T) {
	assert.Equal(t, "TRUE", GetBoolType(true))
	assert.Equal(t, "FALSE", GetBoolType(false))
}

func TestError(t *testing.T) {
	// Test case: error with message in details
	e := Error{
		Message:      errorWithDetails,
		ErrorDetails: []ErrorMessageDetails{{Error: "error1", ErrorMessage: "message1"}},
	}
	assert.EqualError(t, e, "message1")

	// Test case: error with untranslatable error in details
	e = Error{
		Message:      errorWithDetails,
		ErrorDetails: []ErrorMessageDetails{{Error: "error1"}},
	}
	assert.EqualError(t, e, errorWithDetails)

	// Test case: error with translatable error in details
	e = Error{
		Message:      errorWithDetails,
		ErrorDetails: []ErrorMessageDetails{{Error: "ALREADY_EXISTS"}},
	}
	assert.EqualError(t, e, "Already exists")

	// Test case: error without details
	e = Error{
		Message: "message2",
	}
	assert.EqualError(t, e, "message2")
}

// U-GS-01: Statistics struct JSON unmarshal includes compressionRatio
func TestStatistics_UnmarshalCompressionRatio(t *testing.T) {
	raw := `{
		"numOfVolumes":5,
		"capacityInUseInKb":1048576,
		"compressionRatio":2.0,
		"compressedDataCompressionRatio":1.5,
		"netFglCompressedDataSizeInKb":524288
	}`

	var stats Statistics
	err := json.Unmarshal([]byte(raw), &stats)
	require.NoError(t, err)

	assert.Equal(t, 5, stats.NumOfVolumes)
	assert.InDelta(t, 2.0, stats.CompressionRatio, 0.001)
	assert.InDelta(t, 1.5, stats.CompressedDataCompressionRatio, 0.001)
	assert.Equal(t, int64(524288), stats.NetFglCompressedDataSizeInKb)
}

// U-GS-02: ReplicationConsistencyGroupStatistics JSON unmarshal — all fields
func TestReplicationConsistencyGroupStatistics_Unmarshal(t *testing.T) {
	raw := `{
		"lagReceivedInMillis": 5000,
		"lagAppliedInMillis": 4000,
		"lagPersistentInMillis": 3000,
		"lagReceivedSkew": true,
		"lagAppliedSkew": false,
		"lagPersistentSkew": true,
		"rplTransmitBwc":    {"totalWeightInKb": 10240, "numSeconds": 10, "numOccured": 5},
		"rplReceiveBwc":     {"totalWeightInKb": 8192,  "numSeconds": 10, "numOccured": 4},
		"rplRemoteApplyBwc": {"totalWeightInKb": 4096,  "numSeconds": 10, "numOccured": 3},
		"rcgLocalWriteBwc":  {"totalWeightInKb": 2048,  "numSeconds": 5,  "numOccured": 2},
		"rcgLocalReadBwc":   {"totalWeightInKb": 1024,  "numSeconds": 5,  "numOccured": 1},
		"rcgRemoteWriteBwc": {"totalWeightInKb": 512,   "numSeconds": 5,  "numOccured": 1},
		"rcgRemoteReadBwc":  {"totalWeightInKb": 256,   "numSeconds": 5,  "numOccured": 1},
		"rplTransmitLatency":{"totalWeightInKb": 100,   "numSeconds": 10, "numOccured": 5},
		"rplReceiveLatency": {"totalWeightInKb": 80,    "numSeconds": 10, "numOccured": 5},
		"rplApplyLatency":   {"totalWeightInKb": 60,    "numSeconds": 10, "numOccured": 5},
		"rplPairIds": ["pair-1", "pair-2"],
		"numOfRplPairs": 2,
		"initialCopyProgress": 75.5,
		"freezeTransmit": true,
		"isInSlimMode": false
	}`

	var stats ReplicationConsistencyGroupStatistics
	err := json.Unmarshal([]byte(raw), &stats)
	require.NoError(t, err)

	// Lag fields
	assert.Equal(t, int64(5000), stats.LagReceivedInMillis)
	assert.Equal(t, int64(4000), stats.LagAppliedInMillis)
	assert.Equal(t, int64(3000), stats.LagPersistentInMillis)
	assert.True(t, stats.LagReceivedSkew)
	assert.False(t, stats.LagAppliedSkew)
	assert.True(t, stats.LagPersistentSkew)

	// Bandwidth BWC fields
	assert.Equal(t, 10240, stats.RplTransmitBwc.TotalWeightInKb)
	assert.Equal(t, 10, stats.RplTransmitBwc.NumSeconds)
	assert.Equal(t, 8192, stats.RplReceiveBwc.TotalWeightInKb)
	assert.Equal(t, 4096, stats.RplRemoteApplyBwc.TotalWeightInKb)
	assert.Equal(t, 2048, stats.RcgLocalWriteBwc.TotalWeightInKb)
	assert.Equal(t, 1024, stats.RcgLocalReadBwc.TotalWeightInKb)
	assert.Equal(t, 512, stats.RcgRemoteWriteBwc.TotalWeightInKb)
	assert.Equal(t, 256, stats.RcgRemoteReadBwc.TotalWeightInKb)

	// Latency BWC fields
	assert.Equal(t, 100, stats.RplTransmitLatency.TotalWeightInKb)
	assert.Equal(t, 80, stats.RplReceiveLatency.TotalWeightInKb)
	assert.Equal(t, 60, stats.RplApplyLatency.TotalWeightInKb)

	// Scalar fields
	assert.Equal(t, []string{"pair-1", "pair-2"}, stats.RplPairIDs)
	assert.Equal(t, 2, stats.NumOfRplPairs)
	assert.InDelta(t, 75.5, stats.InitialCopyProgress, 0.001)
	assert.True(t, stats.FreezeTransmit)
	assert.False(t, stats.IsInSlimMode)
}

// U-GS-05: BandwidthKBps with NumSeconds=0 returns 0.0 (no divide-by-zero)
func TestBandwidthKBps_ZeroNumSeconds_ReturnsZero(t *testing.T) {
	bwc := BWC{TotalWeightInKb: 1000, NumSeconds: 0}
	result := BandwidthKBps(bwc)
	assert.Equal(t, 0.0, result, "BandwidthKBps with NumSeconds=0 must return 0.0, not NaN or panic")
}

// U-GS-06: BandwidthKBps with valid BWC — TotalWeightInKb / NumSeconds
func TestBandwidthKBps_ValidBWC(t *testing.T) {
	bwc := BWC{TotalWeightInKb: 5000, NumSeconds: 10}
	result := BandwidthKBps(bwc)
	assert.InDelta(t, 500.0, result, 0.001, "5000 KB / 10 s = 500 KB/s")
}

// U-GS-07: RCG bandwidth fields calculation
func TestBandwidthKBps_RCGBandwidth(t *testing.T) {
	transmitBwc := BWC{TotalWeightInKb: 10240, NumSeconds: 10}
	transmitKBps := BandwidthKBps(transmitBwc)
	assert.InDelta(t, 1024.0, transmitKBps, 0.001, "10240 KB / 10 s = 1024 KB/s")

	receiveBwc := BWC{TotalWeightInKb: 8192, NumSeconds: 10}
	receiveKBps := BandwidthKBps(receiveBwc)
	assert.InDelta(t, 819.2, receiveKBps, 0.001, "8192 KB / 10 s = 819.2 KB/s")
}

// U-GS-10: LatencySeconds with NumOccured=0 returns 0.0 (no divide-by-zero)
func TestLatencySeconds_ZeroNumOccured_ReturnsZero(t *testing.T) {
	bwc := BWC{TotalWeightInKb: 50000, NumOccured: 0}
	result := LatencySeconds(bwc)
	assert.Equal(t, 0.0, result, "LatencySeconds with NumOccured=0 must return 0.0, not NaN or panic")
}

// U-GS-11: LatencySeconds computes average I/O latency: totalWeightInKb / NumOccured / 1_000_000
func TestLatencySeconds_ValidBWC(t *testing.T) {
	// 5,000,000 µs total / 1000 ops = 5000 µs avg = 0.005 seconds
	bwc := BWC{TotalWeightInKb: 5_000_000, NumOccured: 1000}
	result := LatencySeconds(bwc)
	assert.InDelta(t, 0.005, result, 0.0000001, "5,000,000 µs / 1000 ops / 1_000_000 = 0.005 seconds")
}

// U-GS-12: LatencySeconds with single op — latency equals TotalWeightInKb / 1_000_000
func TestLatencySeconds_SingleOccurence(t *testing.T) {
	bwc := BWC{TotalWeightInKb: 2000, NumOccured: 1}
	result := LatencySeconds(bwc)
	assert.InDelta(t, 0.002, result, 0.0000001, "2000 µs / 1 op / 1_000_000 = 0.002 seconds")
}

// U-GS-08: CompressionRatio zero value — should not panic
func TestStatistics_CompressionRatioZero(t *testing.T) {
	raw := `{"numOfVolumes":10,"capacityInUseInKb":2097152,"compressionRatio":0.0}`

	var stats Statistics
	err := json.Unmarshal([]byte(raw), &stats)
	require.NoError(t, err)
	assert.Equal(t, 0.0, stats.CompressionRatio)
}

// U-GS-09: CompressionRatio high value — should handle large ratios
func TestStatistics_CompressionRatioHighValue(t *testing.T) {
	raw := `{"numOfVolumes":5,"capacityInUseInKb":1048576,"compressionRatio":10.5}`

	var stats Statistics
	err := json.Unmarshal([]byte(raw), &stats)
	require.NoError(t, err)
	assert.InDelta(t, 10.5, stats.CompressionRatio, 0.001)
}
