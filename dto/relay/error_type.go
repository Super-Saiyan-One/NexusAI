package relay

type GeneralErrorContent struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    any    `json:"code"`
}

type GeneralErrorWithStatusCode struct {
	StatusCode int                 `json:"status_code"`
	Error      GeneralErrorContent `json:"error"`
	IsLocal    bool                `json:"is_local"`
}

type GeneralErrorResponse struct {
	Error    GeneralErrorContent `json:"error"`
	Message  string              `json:"message"`
	Msg      string              `json:"msg"`
	Err      string              `json:"err"`
	ErrorMsg string              `json:"error_msg"`
	Header   struct {
		Message string `json:"message"`
	} `json:"header"`
	Response struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	} `json:"response"`
}

func (g *GeneralErrorResponse) ToMessage() string {
	if g.Error.Message != "" {
		return g.Error.Message
	}
	if g.Message != "" {
		return g.Message
	}
	if g.Msg != "" {
		return g.Msg
	}
	if g.Err != "" {
		return g.Err
	}
	if g.ErrorMsg != "" {
		return g.ErrorMsg
	}
	if g.Header.Message != "" {
		return g.Header.Message
	}
	if g.Response.Error.Message != "" {
		return g.Response.Error.Message
	}
	return ""
}
