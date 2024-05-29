package main

type ForbiddenError struct {
	Message string
}

func (err *ForbiddenError) Error() string {
	return "Forbidden"
}

type UnauthorizedError struct {
	Message string
}

func (err *UnauthorizedError) Error() string {
	return "Unauthorized"
}
