package types

type EventType string

const (
	EventTypeDeploymentHook EventType = "DeploymentHook"
	EventTypeJobHook        EventType = "JobHook"
)
