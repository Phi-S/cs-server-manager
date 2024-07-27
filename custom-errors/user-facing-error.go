package main

import "errors"

func New1(userFriendlyMsg string, internalError error) UserFriendlyError {
	return UserFriendlyError{
		userFriendlyMsg,
		internalError,
	}
}

func New2(userFriendlyMsg string, internalErrorMsg string) UserFriendlyError {
	return UserFriendlyError{
		userFriendlyMsg,
		errors.New(internalErrorMsg),
	}
}

func New3(userFriendlyMsg string) UserFriendlyError {
	return UserFriendlyError{
		userFriendlyMsg,
		errors.New(userFriendlyMsg),
	}
}

// The userFriendlyMsg can safely be returned to the end user
type UserFriendlyError struct {
	userFriendlyMsg string
	internalError   error
}

func (e UserFriendlyError) Error() string {
	return e.internalError.Error()
}

func (e UserFriendlyError) UserFriendlyMsg() string {
	return e.userFriendlyMsg
}
