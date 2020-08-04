package main

type Event interface {
	// Type tells us nature of the event.
	Type() string

	// Command returns command READY to be executed
	Command() string

	// Pattern returns the pattern this event is associated with.
	Pattern() string

	// AbsPath returns full path of file associated with this event
	AbsPath() string

	// BasePath returns base path of file associated with this event
	BasePath() string
}

type event struct {
	fullFilePath string
	baseFilePath string
	rule         Rule
}

func (e event) Type() string {
	return e.rule.Event
}

func (e event) Command() string {
	return e.rule.Command
}

func (e event) Pattern() string {
	return e.rule.Pattern
}

func (e event) AbsPath() string {
	return e.fullFilePath
}

func (e event) BasePath() string {
	return e.baseFilePath
}
