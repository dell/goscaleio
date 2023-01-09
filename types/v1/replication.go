package goscaleio

// PauseMode states in which the ConsistencyGroup can be set to when Paused.
type PauseMode string

// List of pause modes.
const (
	StopDataTransfer PauseMode = "StopDataTransfer"
	OnlyTrackChanges PauseMode = "OnlyTrackChanges"
)

// PauseReplicationConsistencyGroup defines struct for PauseReplicationConsistencyGroup.
type PauseReplicationConsistencyGroup struct {
	PauseMode string `json:"pauseMode"`
}

// SynchronizationResponse defines struct for SynchronizationResponse.
type SynchronizationResponse struct {
	SyncNowKey string `json:"syncNowKey"`
}

// QuerySyncNowRequest defines struct for QuerySyncNowRequest.
type QuerySyncNowRequest struct {
	SyncNowKey string `json:"syncNowKey"`
}
