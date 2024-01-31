package types

import "fmt"

type BadEvent struct {
	Event  *Event
	Err    string
	Reason string
}

func (e *BadEvent) Error() string {
	if e.Event == nil {
		return e.Err
	}
	if e.Reason != "" {
		errStr := fmt.Sprintf("got bad event '%s': %s: %s", e.Event.Name, e.Reason, e.Err)
		return errStr
	}
	errStr := fmt.Sprintf("got bad event '%s': %s", e.Event.Name, e.Err)
	return errStr
}
