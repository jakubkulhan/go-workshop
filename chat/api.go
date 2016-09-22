package chat

const (
	ErrBadRequest            = 400001
	ErrUserNameMissing       = 400002
	ErrUserNameAlreadyExists = 400003
	ErrFromUserIDMissing     = 400004
	ErrToUserIDMissing       = 400005
	ErrCannotChat            = 400006
	ErrTextMissing           = 400007
)

type OkResponse struct {
	OK bool `json:"ok"`
}

type MeResponse struct {
	OK   bool  `json:"ok"`
	Data *User `json:"data"`
}

type UserListResponse struct {
	OK   bool    `json:"ok"`
	Data []*User `json:"data"`
}

type ErrorResponse struct {
	OK         bool   `json:"ok"`
	StatusCode int    `json:"status_code"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}
