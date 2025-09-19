package handlers

type ContextKey string

const UidKey ContextKey = "uid"

type ExecuteAppRequest struct {
	AppName string                 `json:"app_name"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type ExecuteAppResponse struct {
	UID     string `json:"uid"`
	AppName string `json:"app_name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ProfileResponse struct {
	UID       string `json:"uid"`
	Email     string `json:"email,omitempty"`
	TwitterID string `json:"twitter_id,omitempty"`
	Message   string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
