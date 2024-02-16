package handler

type BadRequestError struct {
	reason string
}

func NewBadRequestError(reason string) BadRequestError {
	return BadRequestError{
		reason,
	}
}

func (self BadRequestError) Error() string {
	return self.reason
}

type UnprocessableEntityError struct {
	reason string
}

func NewUnprocessableEntityError(reason string) UnprocessableEntityError {
	return UnprocessableEntityError{
		reason,
	}
}

func (self UnprocessableEntityError) Error() string {
	return self.reason
}
