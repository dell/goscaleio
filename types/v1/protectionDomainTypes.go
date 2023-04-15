package goscaleio

import (
	"fmt"
	"strconv"
)

type PflexParams interface {
	ToMap() (map[string]string, error)
}

// type PDRfCachePageSize int

// const (
// 	PDRfCachePageSize4  PDRfCachePageSize = 4
// 	PDRfCachePageSize8  PDRfCachePageSize = 8
// 	PDRfCachePageSize16 PDRfCachePageSize = 16
// 	PDRfCachePageSize32 PDRfCachePageSize = 32
// 	PDRfCachePageSize64 PDRfCachePageSize = 64
// )

// func GetPDRfCachePageSize(i int) (error, PDRfCachePageSize){
// 	switch i {
// 		4:
// 	}
// }

// type PDRfCacheIOSize int

// const (
// 	PDRfCacheIOSize32  PDRfCacheIOSize = 32
// 	PDRfCacheIOSize64  PDRfCacheIOSize = 64
// 	PDRfCacheIOSize126 PDRfCacheIOSize = 126
// 	PDRfCacheIOSize16  PDRfCacheIOSize = 16
// )

type PDRfCacheParams struct {
	RfCacheOperationalMode PDRfCacheOpMode `json:"rfcacheOperationMode"`
	RfCachePageSizeKb      int             `json:"pageSizeKb"`
	RfCacheMaxIoSizeKb     int             `json:"maxIOSizeKb"`
}

func (params PDRfCacheParams) ToMap() (map[string]string, error) {
	m := make(map[string]string)
	dict := map[int]bool{4: true, 8: true, 16: true, 32: true, 64: true}
	if params.RfCachePageSizeKb != 0 {
		if _, ok := dict[params.RfCachePageSizeKb]; !ok {
			return m, fmt.Errorf("")
		}
		m["pageSizeKb"] = strconv.Itoa(params.RfCachePageSizeKb)
	}
	if params.RfCacheMaxIoSizeKb != 0 {
		if size := params.RfCachePageSizeKb; !(size == 16 || size == 32 || size == 64 || size == 126) {
			return m, fmt.Errorf("")
		}
		m["maxIOSizeKb"] = strconv.Itoa(params.RfCacheMaxIoSizeKb)
	}
	if params.RfCacheOperationalMode != "" {
		m["rfcacheOperationMode"] = string(params.RfCacheOperationalMode)
	}
	return m, nil
}

func (pd *ProtectionDomain) GetRfCacheParams() PDRfCacheParams {
	return PDRfCacheParams{
		RfCacheOperationalMode: pd.RfCacheOperationalMode,
		RfCachePageSizeKb:      pd.RfCachePageSizeKb,
		RfCacheMaxIoSizeKb:     pd.RfCacheMaxIoSizeKb,
	}
}

type SdsNetworkLimitParams struct {
	RebuildNetworkThrottlingInKbps                  *int `json:"rebuildLimitInKbps"`
	RebalanceNetworkThrottlingInKbps                *int `json:"rebalanceLimitInKbps"`
	VTreeMigrationNetworkThrottlingInKbps           *int `json:"vtreeMigrationLimitInKbps"`
	ProtectedMaintenanceModeNetworkThrottlingInKbps *int `json:"protectedMaintenanceModeLimitInKbps"`
	OverallIoNetworkThrottlingInKbps                *int `json:"overallLimitInKbps"`
}

func (params SdsNetworkLimitParams) ToMap() (map[string]string, error) {
	m := make(map[string]string)
	if size := params.RebuildNetworkThrottlingInKbps; size != nil {
		m["rebuildLimitInKbps"] = strconv.Itoa(*size)
	}
	if size := params.RebalanceNetworkThrottlingInKbps; size != nil {
		m["rebalanceLimitInKbps"] = strconv.Itoa(*size)
	}
	if size := params.VTreeMigrationNetworkThrottlingInKbps; size != nil {
		m["vtreeMigrationLimitInKbps"] = strconv.Itoa(*size)
	}
	if size := params.ProtectedMaintenanceModeNetworkThrottlingInKbps; size != nil {
		m["protectedMaintenanceModeLimitInKbps"] = strconv.Itoa(*size)
	}
	if size := params.OverallIoNetworkThrottlingInKbps; size != nil {
		m["overallLimitInKbps"] = strconv.Itoa(*size)
	}
	return m, nil
}
