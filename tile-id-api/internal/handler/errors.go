package handler

type BadRequestError struct {
	reason string
}

func NewBadRequestError(reason string) BadRequestError {
	return BadRequestError{
		reason,
	}
}

func (err BadRequestError) Error() string {
	return err.reason
}

type UnprocessableEntityError struct {
	reason string
}

func NewUnprocessableEntityError(reason string) UnprocessableEntityError {
	return UnprocessableEntityError{
		reason,
	}
}

func (err UnprocessableEntityError) Error() string {
	return err.reason
}
