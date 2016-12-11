package context

//ErrorFunc is a func that returns error
type ErrorFunc func() error

//SetupFunc is a func that initializes and returns error if unexpected results happens
type SetupFunc ErrorFunc
