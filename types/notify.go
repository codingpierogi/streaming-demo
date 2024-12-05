package types

type Notification struct {
	Message string `json:"message"`
}

type NotifyMessage struct {
	Key   string       `json:"key"`
	Value Notification `json:"value"`
}
