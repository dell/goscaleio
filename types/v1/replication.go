package goscaleio

type PauseMode string

const (
	STOP_DATA_TRANSFER PauseMode = "StopDataTransfer"
	ONLY_TRACK_CHANGES PauseMode = "OnlyTrackChanges"
)

type ActionType string

const (
	RS_ACTION_FAILOVER  ActionType = "failover"
	RS_ACTION_REPROTECT ActionType = "reprotect"
	RS_ACTION_RESTORE   ActionType = "restore"
	RS_ACTION_RESUME    ActionType = "resume"
	RS_ACTION_PAUSE     ActionType = "pause"
	RS_ACTION_SYNC      ActionType = "sync"
)

// failover params create failover request
type FailoverParams struct {
	// For DR failover.
	IsPlanned bool `json:"is_planned, omitempty"`
}

type PauseReplicationConsistencyGroup struct {
	PauseMode string `json:"pauseMode"`
}
