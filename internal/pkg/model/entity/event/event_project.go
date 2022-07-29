package event

type EventProject struct {
	ID       int64
	EventId  string
	Code     string
	Executor string
	Payee    string
	Time     string
}

func (EventProject) TableName() string {
	return "event_project"
}
