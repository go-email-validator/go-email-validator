package mailboxvalidator

import (
	"encoding/json"
)

const (
	// MBVTrue is "true" in the return
	MBVTrue = "True"
	// MBVFalse is "False" in the return
	MBVFalse = "False"
	// MBVEmpty is "" in the return
	MBVEmpty = ""
)

// EmptyBool is a type for mailboxvalidator bool
type EmptyBool struct {
	*bool
}

// NewEmptyBool create EmptyBool with bool
func NewEmptyBool(boolean bool) EmptyBool {
	return EmptyBool{
		bool: &boolean,
	}
}

// NewEmptyBoolWithNil create EmptyBool with nil
func NewEmptyBoolWithNil() EmptyBool {
	return EmptyBool{
		bool: nil,
	}
}

// ToBool returns boolean value
func (e EmptyBool) ToBool() bool {
	if e.bool == nil {
		return false
	}

	return *e.bool
}

// UnmarshalJSON implements json.UnmarshalJSON
func (e *EmptyBool) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &e.bool)

	return err
}

// MarshalJSON implements json.Marshaler
func (e *EmptyBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.bool)
}

// ToString converts bool to string
func (e EmptyBool) ToString() string {
	if e.bool == nil {
		return ""
	}

	if *e.bool {
		return MBVTrue
	}
	return MBVFalse
}

// ToBool converts string to bool
func ToBool(value string) (result EmptyBool) {
	if value == MBVEmpty {
		return
	}

	boolean := value == MBVTrue
	return EmptyBool{&boolean}
}
