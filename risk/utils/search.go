package utils

import "github.com/bcdannyboy/dgws/risk"

func FindEvent(ID int, events []*risk.Event) *risk.Event {
	for _, event := range events {
		if event.ID == ID {
			return event
		}
	}
	return nil
}
