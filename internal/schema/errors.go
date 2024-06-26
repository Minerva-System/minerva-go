package schema

import (
	"net/http"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type ErrorMessage struct {
	Status int     `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (ErrorMessage) FromGrpcError(err error) ErrorMessage {
	var httpstatus int = http.StatusInternalServerError
	var message string = "An error has ocurred while communicating with the service"
	var details string = ""
	
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.AlreadyExists:
			httpstatus = http.StatusConflict
			message = "Entity already exists"
			details = s.Message()
		case codes.PermissionDenied:
			httpstatus = http.StatusForbidden
			message = "You do not have permission to access this resource."
			details = s.Message()
		case codes.InvalidArgument:
			httpstatus = http.StatusBadRequest
			message = "The service did not receive valid information."
			details = s.Message()
		case codes.OutOfRange:
			httpstatus = http.StatusBadRequest
			message = "The service received a request for a resource that is out of range."
			details = s.Message()
		case codes.Internal:
			httpstatus = http.StatusInternalServerError
			message = "There was an error on the service while processing your request."
			details = s.Message()
		case codes.NotFound:
			httpstatus = http.StatusNotFound
			message = "Fetched entity was not found"
			details = s.Message()
		}
	} else {
		message = "Unknown error while processing request"
	}

	return ErrorMessage{
		Status: httpstatus,
		Message: message,
		Details: details,
	}
}
