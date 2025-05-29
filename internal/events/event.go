package events

import "encoding/json"

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"` // Raw JSON kept intact
}

func NewEvent(name string, data []byte) *Event {
	return &Event{Type: name, Data: data}
}

func MarshalEvent(e *Event) ([]byte, error) {
	return json.Marshal(e)
}

func UnmarshalEvent(data []byte) (Event, error) {
	var e Event
	err := json.Unmarshal(data, &e)
	return e, err
}
