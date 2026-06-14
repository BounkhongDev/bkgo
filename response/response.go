package response

// Response is the standard HTTP JSON response envelope used across all endpoints.
type Response struct {
	Success   bool   `json:"success"`
	Data      any    `json:"data,omitempty"`
	ErrorCode string `json:"error,omitempty"` // renamed from Error to avoid clash with Error() function
	Message   string `json:"message"`
	Meta      *Meta  `json:"meta,omitempty"`
}

// Meta holds pagination information for list responses.
type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// Success returns a successful response with data.
func Success(data any) Response {
	return Response{Success: true, Data: data, Message: "ok"}
}

// SuccessMessage returns a successful response with a custom message and no data.
func SuccessMessage(message string) Response {
	return Response{Success: true, Message: message}
}

// Paginated returns a successful list response with pagination meta.
func Paginated(data any, page, limit int, total int64) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: "ok",
		Meta:    &Meta{Page: page, Limit: limit, Total: total},
	}
}

// Error returns a failure response with an error code and message.
func Error(code, message string) Response {
	return Response{Success: false, ErrorCode: code, Message: message}
}
