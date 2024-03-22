package schema

type ErrorMessage struct {
	Status int     `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

