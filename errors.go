package modelsdev

import "errors"

// ErrNotFound is returned by lookup methods when no provider or model matches.
var ErrNotFound = errors.New("modelsdev: not found")

// ErrInvalidID is returned when a fully qualified model ID is not in the
// "provider:model" form.
var ErrInvalidID = errors.New("modelsdev: invalid model id, want \"provider:model\"")
