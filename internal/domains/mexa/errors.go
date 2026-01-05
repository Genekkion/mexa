package mexadomain

type MexaError struct {
	message string
}

func NewMexaError(message string) MexaError {
	return MexaError{message: message}
}

func (e MexaError) Error() string {
	return e.message
}

var (
	NoCasualtyFoundError = NewMexaError("no casualty found")
)
