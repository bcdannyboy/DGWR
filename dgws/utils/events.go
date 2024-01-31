package utils

import "github.com/bcdannyboy/montecargo/dgws/types"

func eventExists(eventID uint64, EventMap map[uint64]*types.Event) bool {
	_, exists := EventMap[eventID]
	return exists
}

func riskExists(riskID uint64, RiskMap map[uint64]*types.Risk) bool {
	_, exists := RiskMap[riskID]
	return exists
}

func mitigationExists(mitigationID uint64, MitigationMap map[uint64]*types.Mitigation) bool {
	_, exists := MitigationMap[mitigationID]
	return exists
}

func appendBadEvent(event *types.Event, errMsg string, BadEvents *[]*types.BadEvent) {
	*BadEvents = append(*BadEvents, &types.BadEvent{
		Event: event,
		Err:   errMsg,
	})
}

func createFilteredEvent(eventKey uint64, event *types.Event, depType uint64, dependentEventID *uint64) *FilteredEvent {
	return &FilteredEvent{
		ID:               eventKey,
		Event:            event,
		Independent:      false,
		DependencyType:   depType,
		DependentEventID: dependentEventID,
	}
}
