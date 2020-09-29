package swe

// SweError represents an error reported by the Swiss Ephemeris library.
type SweError struct {
	msg string
	err error
}

func NewSweError(msg string) error {
	return &SweError{
		msg: msg,
		err: nil,
	}
}

func (e *SweError) Error() string {
	return "swisseph: " + e.msg
}

func (e *SweError) Unwrap() error {
	return e.err
}