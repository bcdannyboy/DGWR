package utils

import (
	"errors"

	"github.com/bcdannyboy/montecargo/dgws/types"
)

type FilteredEvent struct {
	ID               uint64              `json:"id"`
	Event            *types.Event        `json:"event"`
	Independent      bool                `json:"independent"`
	DependencyType   uint64              `json:"dependency_type"`
	DependentEventID *uint64             `json:"dependent_event_hash,omitempty"`
	DependencyValue  *types.SingleNumber `json:"dependency_value,omitempty"`
	DependencyRange  *types.Range        `json:"dependency_range,omitempty"`
	DependencyDecomp *types.Decomposed   `json:"dependency_decomposed,omitempty"`
}

// FilterDependencies filters out dependent events from the list of events
func FilterDependencies(Events []*types.Event, Risks []*types.Risk, Mitigations []*types.Mitigation) ([]*FilteredEvent, []*types.BadEvent) {
	EventMap := make(map[uint64]*types.Event)
	RiskMap := make(map[uint64]*types.Risk)
	MitigationMap := make(map[uint64]*types.Mitigation)

	BadEvents := []*types.BadEvent{}
	FilteredEvents := []*FilteredEvent{}

	for _, event := range Events {
		EventMap[event.ID] = event
	}

	for _, risk := range Risks {
		RiskMap[risk.ID] = risk
	}

	for _, mitigation := range Mitigations {
		MitigationMap[mitigation.ID] = mitigation
	}

	for eventKey, event := range EventMap {
		if isIndependentEvent(event) {
			FilteredEvents = append(FilteredEvents, &FilteredEvent{
				ID:          eventKey,
				Event:       event,
				Independent: true,
			})
			continue
		}

		processEventDependencies(event, eventKey, EventMap, &FilteredEvents, &BadEvents)
		processProbabilityDependencies(event, eventKey, EventMap, &FilteredEvents, &BadEvents)
		processImpactDependencies(event, eventKey, EventMap, &FilteredEvents, &BadEvents)
		processCostDependencies(event, eventKey, EventMap, MitigationMap, &FilteredEvents, &BadEvents)
		processRiskDependencies(event, eventKey, EventMap, RiskMap, &FilteredEvents, &BadEvents)
		processMitigationDependencies(event, eventKey, EventMap, MitigationMap, &FilteredEvents, &BadEvents)
	}

	return FilteredEvents, BadEvents
}

func processEventDependencies(event *types.Event, eventKey uint64, EventMap map[uint64]*types.Event, FilteredEvents *[]*FilteredEvent, BadEvents *[]*types.BadEvent) {
	for _, dep := range event.DependsOnEvent {
		if _, exists := EventMap[dep.DependentEventID]; !exists {
			*BadEvents = append(*BadEvents, &types.BadEvent{
				Event: event,
				Err:   "Dependent event does not exist",
			})
			continue
		}

		*FilteredEvents = append(*FilteredEvents, &FilteredEvent{
			ID:               eventKey,
			Event:            event,
			Independent:      false,
			DependencyType:   dep.Type,
			DependentEventID: &dep.DependentEventID,
		})
	}
}

func processProbabilityDependencies(event *types.Event, eventKey uint64, EventMap map[uint64]*types.Event, FilteredEvents *[]*FilteredEvent, BadEvents *[]*types.BadEvent) {
	for _, dep := range event.DependsOnProbability {
		if dep.DependentEventID != nil && !eventExists(*dep.DependentEventID, EventMap) {
			appendBadEvent(event, "Dependent event in probability does not exist", BadEvents)
			continue
		}

		fe := createFilteredEvent(eventKey, event, dep.Type, dep.DependentEventID)
		fe.DependencyValue = dep.SingleValue
		fe.DependencyRange = dep.Range
		fe.DependencyDecomp = dep.Decomposed

		*FilteredEvents = append(*FilteredEvents, fe)
	}
}

func processImpactDependencies(event *types.Event, eventKey uint64, EventMap map[uint64]*types.Event, FilteredEvents *[]*FilteredEvent, BadEvents *[]*types.BadEvent) {
	for _, dep := range event.DependsOnImpact {
		if dep.DependentEventID != nil && !eventExists(*dep.DependentEventID, EventMap) {
			appendBadEvent(event, "Dependent event in impact does not exist", BadEvents)
			continue
		}

		fe := createFilteredEvent(eventKey, event, dep.Type, dep.DependentEventID)
		fe.DependencyValue = dep.SingleValue
		fe.DependencyRange = dep.Range
		fe.DependencyDecomp = dep.Decomposed

		*FilteredEvents = append(*FilteredEvents, fe)
	}
}

func processCostDependencies(event *types.Event, eventKey uint64, EventMap map[uint64]*types.Event, MitigationMap map[uint64]*types.Mitigation, FilteredEvents *[]*FilteredEvent, BadEvents *[]*types.BadEvent) {
	for _, dep := range event.DependsOnCost {
		if dep.DependentEventID != nil && !eventExists(*dep.DependentEventID, EventMap) {
			appendBadEvent(event, "Dependent event in cost does not exist", BadEvents)
			continue
		}

		if dep.DependentMitigationOrRiskID != nil && !mitigationExists(*dep.DependentMitigationOrRiskID, MitigationMap) {
			appendBadEvent(event, "Dependent mitigation or risk in cost does not exist", BadEvents)
			continue
		}

		fe := createFilteredEvent(eventKey, event, dep.Type, dep.DependentEventID)
		fe.DependencyValue = dep.SingleValue
		fe.DependencyRange = dep.Range
		fe.DependencyDecomp = dep.Decomposed

		*FilteredEvents = append(*FilteredEvents, fe)
	}
}

func processRiskDependencies(event *types.Event, eventKey uint64, EventMap map[uint64]*types.Event, RiskMap map[uint64]*types.Risk, FilteredEvents *[]*FilteredEvent, BadEvents *[]*types.BadEvent) {
	for _, dep := range event.DependsOnRisk {
		if dep.DependentRiskID != nil && !riskExists(*dep.DependentRiskID, RiskMap) {
			appendBadEvent(event, "Dependent risk does not exist", BadEvents)
			continue
		}

		if dep.DependentEventID != nil && !eventExists(*dep.DependentEventID, EventMap) {
			appendBadEvent(event, "Dependent event in risk does not exist", BadEvents)
			continue
		}

		fe := createFilteredEvent(eventKey, event, dep.Type, dep.DependentEventID)
		*FilteredEvents = append(*FilteredEvents, fe)
	}
}

func processMitigationDependencies(event *types.Event, eventKey uint64, EventMap map[uint64]*types.Event, MitigationMap map[uint64]*types.Mitigation, FilteredEvents *[]*FilteredEvent, BadEvents *[]*types.BadEvent) {
	for _, dep := range event.DependsOnMitigation {
		if dep.DependentMitigationOrRiskID != nil && !mitigationExists(*dep.DependentMitigationOrRiskID, MitigationMap) {
			appendBadEvent(event, "Dependent mitigation or risk does not exist", BadEvents)
			continue
		}

		if dep.DependentEventID != nil && !eventExists(*dep.DependentEventID, EventMap) {
			appendBadEvent(event, "Dependent event in mitigation does not exist", BadEvents)
			continue
		}

		fe := createFilteredEvent(eventKey, event, dep.Type, dep.DependentEventID)
		*FilteredEvents = append(*FilteredEvents, fe)
	}
}

func isIndependentEvent(event *types.Event) bool {
	return len(event.DependsOnEvent) == 0 &&
		len(event.DependsOnProbability) == 0 &&
		len(event.DependsOnImpact) == 0 &&
		len(event.DependsOnRisk) == 0 &&
		len(event.DependsOnCost) == 0 &&
		len(event.DependsOnMitigation) == 0
}

// FindEventByID finds an event by its ID from a slice of FilteredEvent.
func FindEventByID(id uint64, events []*FilteredEvent) (*FilteredEvent, error) {
	for _, event := range events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, errors.New("event not found")
}
