package identity

import "encoding/json"

type Event string

const (
	EventUserCreated = "UserCreated"
)

func PublishUserCreatedEvent(userID int64) {
	payload := struct {
		UserID int64 `json:"userID"`
	}{userID}

	str, err := json.Marshal(payload)
	if err != nil {
		return
	}

	EventPublisher.Publish(EventUserCreated, str)
}
