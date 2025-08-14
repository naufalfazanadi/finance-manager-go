package utils

import "fmt"

const (
	// DefaultPageSize is the default number of items per page for pagination.
	DefaultPageSize          = 10
	MsgErrInvalidID          = "Invalid ID format"
	ErrInvalidIDFormat       = "ID must be a valid UUID"
	MsgErrIDRequired         = "ID is required"
	ErrIDRequired            = "ID parameter is missing"
	MsgErrFormValidation     = "Form validation failed"
	MsgErrInvalidQueryParams = "Invalid query parameters"
	MsgErrReqBody            = "Invalid request body"
)

// SuccessMsg
func SuccessRetrieveMsg(entity string) string {
	return fmt.Sprintf("%s retrieved successfully", entity)
}
func SuccessCreateMsg(entity string) string {
	return fmt.Sprintf("%s created successfully", entity)
}
func SuccessUpdateMsg(entity string) string {
	return fmt.Sprintf("%s updated successfully", entity)
}
func SuccessDeleteMsg(entity string) string {
	return fmt.Sprintf("%s deleted successfully", entity)
}

// FailedMsg
func FailedGetMsg(entity string) string {
	return fmt.Sprintf("Failed to get %s", entity)
}
func FailedCreateMsg(entity string) string {
	return fmt.Sprintf("Failed to create %s", entity)
}
func FailedUpdateMsg(entity string) string {
	return fmt.Sprintf("Failed to update %s", entity)
}
func FailedDeleteMsg(entity string) string {
	return fmt.Sprintf("Failed to delete %s", entity)
}
func FailedRestoreMsg(entity string) string {
	return fmt.Sprintf("Failed to restore %s", entity)
}
