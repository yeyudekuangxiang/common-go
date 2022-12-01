package types

type EventHookType string

const (
	EventDeploymentHook   EventHookType = "Deployment Hook"
	EventMergeRequestHook EventHookType = "Merge Request Hook"
)
