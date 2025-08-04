package cli

import (
	"fmt"
	"strings"
	"time"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	ErrorCategoryValidation  ErrorCategory = "validation"
	ErrorCategoryNotFound    ErrorCategory = "not_found"
	ErrorCategoryConflict    ErrorCategory = "conflict"
	ErrorCategoryAuth        ErrorCategory = "auth"
	ErrorCategoryNetwork     ErrorCategory = "network"
	ErrorCategoryTimeout     ErrorCategory = "timeout"
	ErrorCategoryUnsupported ErrorCategory = "unsupported"
	ErrorCategoryInternal    ErrorCategory = "internal"
)

// ExitCode represents CLI exit codes
type ExitCode int

const (
	ExitCodeSuccess     ExitCode = 0
	ExitCodeInternal    ExitCode = 1
	ExitCodeValidation  ExitCode = 2
	ExitCodeNotFound    ExitCode = 3
	ExitCodeConflict    ExitCode = 4
	ExitCodeAuth        ExitCode = 5
	ExitCodeNetwork     ExitCode = 6
	ExitCodeTimeout     ExitCode = 7
	ExitCodeUnsupported ExitCode = 8
)

// ErrorContract represents a structured error with remediation information
type ErrorContract struct {
	Code      string            `json:"code"`
	Category  ErrorCategory     `json:"category"`
	Message   string            `json:"message"`
	Hint      string            `json:"hint,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	RequestID string            `json:"requestId,omitempty"`
}

// ToExitCode returns the appropriate exit code for the error category
func (ec ErrorContract) ToExitCode() ExitCode {
	switch ec.Category {
	case ErrorCategoryValidation:
		return ExitCodeValidation
	case ErrorCategoryNotFound:
		return ExitCodeNotFound
	case ErrorCategoryConflict:
		return ExitCodeConflict
	case ErrorCategoryAuth:
		return ExitCodeAuth
	case ErrorCategoryNetwork:
		return ExitCodeNetwork
	case ErrorCategoryTimeout:
		return ExitCodeTimeout
	case ErrorCategoryUnsupported:
		return ExitCodeUnsupported
	case ErrorCategoryInternal:
		fallthrough
	default:
		return ExitCodeInternal
	}
}

// ToOutputError converts ErrorContract to ErrorData for output formatting
func (ec ErrorContract) ToOutputError() ErrorData {
	return ErrorData{
		Code:      ec.Code,
		Category:  string(ec.Category),
		Message:   ec.Message,
		Hint:      ec.Hint,
		Details:   ec.Details,
		Timestamp: ec.Timestamp,
		RequestID: ec.RequestID,
	}
}

// Common error codes
const (
	ErrorCodeSessionNotFound             = "session_not_found"
	ErrorCodeSessionAlreadyExists        = "session_already_exists"
	ErrorCodeSessionNotRunning           = "session_not_running"
	ErrorCodeSessionAlreadyRunning       = "session_already_running"
	ErrorCodeInvalidSessionName          = "invalid_session_name"
	ErrorCodeInvalidSessionID            = "invalid_session_id"
	ErrorCodeMissingRequiredFlag         = "missing_required_flag"
	ErrorCodeInvalidFlagValue            = "invalid_flag_value"
	ErrorCodeMutuallyExclusiveFlags      = "mutually_exclusive_flags"
	ErrorCodeMultiplexerNotAvailable     = "multiplexer_not_available"
	ErrorCodeMultiplexerConnectionFailed = "multiplexer_connection_failed"
	ErrorCodeMultiplexerTimeout          = "multiplexer_timeout"
	ErrorCodeAttachmentFailed            = "attachment_failed"
	ErrorCodePermissionDenied            = "permission_denied"
	ErrorCodeConfigurationError          = "configuration_error"
	ErrorCodeStorageError                = "storage_error"
	ErrorCodeNotInteractiveTerminal      = "not_interactive_terminal"
)

// ErrorMapper provides methods to convert various error types to ErrorContract
type ErrorMapper struct {
	requestID string
}

// NewErrorMapper creates a new ErrorMapper instance
func NewErrorMapper(requestID string) *ErrorMapper {
	return &ErrorMapper{
		requestID: requestID,
	}
}

// MapError converts a generic error to an ErrorContract
func (em *ErrorMapper) MapError(err error) ErrorContract {
	if err == nil {
		return ErrorContract{}
	}

	// Try to extract structured error information
	if contract, ok := err.(*StructuredError); ok {
		return ErrorContract{
			Code:      contract.Code,
			Category:  contract.Category,
			Message:   contract.Message,
			Hint:      contract.Hint,
			Details:   contract.Details,
			Timestamp: time.Now(),
			RequestID: em.requestID,
		}
	}

	// Map common error patterns from service layer
	errorMessage := err.Error()
	lowerMsg := strings.ToLower(errorMessage)

	// Session-related errors
	if strings.Contains(lowerMsg, "session not found") || strings.Contains(lowerMsg, "no such session") {
		return em.createErrorContract(
			ErrorCodeSessionNotFound,
			ErrorCategoryNotFound,
			errorMessage,
			"Use 'claude-pilot list' to see available sessions.",
			nil,
		)
	}

	if strings.Contains(lowerMsg, "session already exists") || strings.Contains(lowerMsg, "duplicate session") {
		return em.createErrorContract(
			ErrorCodeSessionAlreadyExists,
			ErrorCategoryConflict,
			errorMessage,
			"Use a different session name or attach to the existing session.",
			nil,
		)
	}

	if strings.Contains(lowerMsg, "session not running") || strings.Contains(lowerMsg, "session inactive") {
		return em.createErrorContract(
			ErrorCodeSessionNotRunning,
			ErrorCategoryConflict,
			errorMessage,
			"Start the session first before attempting to attach.",
			nil,
		)
	}

	// Validation errors
	if strings.Contains(lowerMsg, "invalid session name") || strings.Contains(lowerMsg, "malformed name") {
		return em.createErrorContract(
			ErrorCodeInvalidSessionName,
			ErrorCategoryValidation,
			errorMessage,
			"Session names must contain only alphanumeric characters, hyphens, and underscores.",
			nil,
		)
	}

	if strings.Contains(lowerMsg, "required") && strings.Contains(lowerMsg, "missing") {
		return em.createErrorContract(
			ErrorCodeMissingRequiredFlag,
			ErrorCategoryValidation,
			errorMessage,
			"Check the command usage with --help for required flags.",
			nil,
		)
	}

	// Network and connection errors
	if strings.Contains(lowerMsg, "connection refused") || strings.Contains(lowerMsg, "cannot connect") {
		return em.createErrorContract(
			ErrorCodeMultiplexerConnectionFailed,
			ErrorCategoryNetwork,
			errorMessage,
			"Ensure the terminal multiplexer (tmux) is running and accessible.",
			nil,
		)
	}

	if strings.Contains(lowerMsg, "timeout") || strings.Contains(lowerMsg, "timed out") {
		return em.createErrorContract(
			ErrorCodeMultiplexerTimeout,
			ErrorCategoryTimeout,
			errorMessage,
			"Try again or check if the multiplexer is responding.",
			nil,
		)
	}

	// Permission errors
	if strings.Contains(lowerMsg, "permission denied") || strings.Contains(lowerMsg, "access denied") {
		return em.createErrorContract(
			ErrorCodePermissionDenied,
			ErrorCategoryAuth,
			errorMessage,
			"Check file permissions and user access rights.",
			nil,
		)
	}

	// Configuration errors
	if strings.Contains(lowerMsg, "configuration") || strings.Contains(lowerMsg, "config") {
		return em.createErrorContract(
			ErrorCodeConfigurationError,
			ErrorCategoryValidation,
			errorMessage,
			"Check your configuration file or environment variables.",
			nil,
		)
	}

	// Storage errors
	if strings.Contains(lowerMsg, "storage") || strings.Contains(lowerMsg, "file system") {
		return em.createErrorContract(
			ErrorCodeStorageError,
			ErrorCategoryInternal,
			errorMessage,
			"Check disk space and file permissions.",
			nil,
		)
	}

	// TTY errors
	if strings.Contains(lowerMsg, "not a terminal") || strings.Contains(lowerMsg, "tty") {
		return em.createErrorContract(
			ErrorCodeNotInteractiveTerminal,
			ErrorCategoryUnsupported,
			errorMessage,
			"This command requires an interactive terminal. Use --yes flag for non-interactive mode.",
			nil,
		)
	}

	// Default to internal error for unrecognized errors
	return em.createErrorContract(
		"unknown_error",
		ErrorCategoryInternal,
		errorMessage,
		"If this problem persists, please report it as a bug.",
		nil,
	)
}

// createErrorContract is a helper to create ErrorContract instances
func (em *ErrorMapper) createErrorContract(code string, category ErrorCategory, message, hint string, details map[string]string) ErrorContract {
	return ErrorContract{
		Code:      code,
		Category:  category,
		Message:   message,
		Hint:      hint,
		Details:   details,
		Timestamp: time.Now(),
		RequestID: em.requestID,
	}
}

// StructuredError implements error interface with additional structure
type StructuredError struct {
	Code     string
	Category ErrorCategory
	Message  string
	Hint     string
	Details  map[string]string
}

// Error implements the error interface
func (se *StructuredError) Error() string {
	return se.Message
}

// NewStructuredError creates a new StructuredError
func NewStructuredError(code string, category ErrorCategory, message, hint string, details map[string]string) *StructuredError {
	return &StructuredError{
		Code:     code,
		Category: category,
		Message:  message,
		Hint:     hint,
		Details:  details,
	}
}

// Common error constructors for convenience
func NewValidationError(message, hint string) *StructuredError {
	return NewStructuredError(
		ErrorCodeInvalidFlagValue,
		ErrorCategoryValidation,
		message,
		hint,
		nil,
	)
}

func NewNotFoundError(resourceType, resourceID string) *StructuredError {
	return NewStructuredError(
		ErrorCodeSessionNotFound,
		ErrorCategoryNotFound,
		fmt.Sprintf("%s '%s' not found", resourceType, resourceID),
		fmt.Sprintf("Use 'claude-pilot list' to see available %ss.", strings.ToLower(resourceType)),
		map[string]string{
			"resourceType": resourceType,
			"resourceID":   resourceID,
		},
	)
}

func NewConflictError(message, hint string) *StructuredError {
	return NewStructuredError(
		ErrorCodeSessionAlreadyExists,
		ErrorCategoryConflict,
		message,
		hint,
		nil,
	)
}

func NewNetworkError(message, hint string) *StructuredError {
	return NewStructuredError(
		ErrorCodeMultiplexerConnectionFailed,
		ErrorCategoryNetwork,
		message,
		hint,
		nil,
	)
}

func NewTimeoutError(message, hint string) *StructuredError {
	return NewStructuredError(
		ErrorCodeMultiplexerTimeout,
		ErrorCategoryTimeout,
		message,
		hint,
		nil,
	)
}

func NewAuthError(message, hint string) *StructuredError {
	return NewStructuredError(
		ErrorCodePermissionDenied,
		ErrorCategoryAuth,
		message,
		hint,
		nil,
	)
}

func NewUnsupportedError(message, hint string) *StructuredError {
	return NewStructuredError(
		ErrorCodeNotInteractiveTerminal,
		ErrorCategoryUnsupported,
		message,
		hint,
		nil,
	)
}

// WrapError wraps a generic error with additional context
func WrapError(err error, code string, category ErrorCategory, hint string) *StructuredError {
	if err == nil {
		return nil
	}

	return NewStructuredError(code, category, err.Error(), hint, nil)
}

// IsErrorCode checks if an error matches a specific error code
func IsErrorCode(err error, code string) bool {
	if structErr, ok := err.(*StructuredError); ok {
		return structErr.Code == code
	}
	return false
}

// IsErrorCategory checks if an error matches a specific category
func IsErrorCategory(err error, category ErrorCategory) bool {
	if structErr, ok := err.(*StructuredError); ok {
		return structErr.Category == category
	}
	return false
}

// PrintError prints an error using the provided OutputWriter
func PrintError(contract ErrorContract, writer OutputWriter) error {
	return writer.WriteError(contract.ToOutputError())
}

// HandleError converts an error to ErrorContract and prints it, then returns the exit code
func HandleError(err error, writer OutputWriter, requestID string) ExitCode {
	if err == nil {
		return ExitCodeSuccess
	}

	mapper := NewErrorMapper(requestID)
	contract := mapper.MapError(err)

	PrintError(contract, writer)
	return contract.ToExitCode()
}
