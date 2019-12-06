package nw_capture

import (
	"encoding/json"
	"github.com/watermint/toolbox/infra/api/api_request"
	"github.com/watermint/toolbox/infra/api/api_response"
	"github.com/watermint/toolbox/infra/control/app_root"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type Capture interface {
	WithResponse(req api_request.Request, res api_response.Response, resErr error, latency int64)
	NoResponse(req api_request.Request, resErr error, latency int64)
}

func currentImpl(cap *zap.Logger) Capture {
	return &captureImpl{
		capture: cap,
	}
}

func Current() Capture {
	cap := app_root.Capture()
	return currentImpl(cap)
}

type Record struct {
	Timestamp      time.Time         `json:"timestamp"`
	RequestMethod  string            `json:"req_method"`
	RequestUrl     string            `json:"req_url"`
	RequestParam   string            `json:"req_param,omitempty"`
	RequestHeaders map[string]string `json:"req_headers"`
	ResponseCode   int               `json:"res_code"`
	ResponseBody   string            `json:"res_body,omitempty"`
	ResponseError  string            `json:"res_error,omitempty"`
	Latency        int64             `json:"latency"`
}

type mockImpl struct {
}

func (mockImpl) WithResponse(req api_request.Request, res api_response.Response, resErr error, latency int64) {
	// ignore
}

var (
	tokenMatcher = regexp.MustCompile(`\w`)
)

func NewCapture(cap *zap.Logger) Capture {
	return &captureImpl{
		capture: cap,
	}
}

type captureImpl struct {
	capture *zap.Logger
}
type Req struct {
	RequestMethod  string            `json:"method"`
	RequestUrl     string            `json:"url"`
	RequestParam   string            `json:"param,omitempty"`
	RequestHeaders map[string]string `json:"headers"`
}

func (z *Req) Apply(req api_request.Request) {
	z.RequestMethod = req.Method()
	z.RequestUrl = req.Url()
	z.RequestParam = req.ParamString()
	z.RequestHeaders = make(map[string]string)
	for k, v := range req.Headers() {
		// Anonymize token
		if k == api_request.ReqHeaderAuthorization {
			z.RequestHeaders[k] = "Bearer <secret>"
		} else {
			z.RequestHeaders[k] = v
		}
	}
}

type Res struct {
	ResponseCode    int               `json:"code"`
	ResponseBody    string            `json:"body,omitempty"`
	ResponseHeaders map[string]string `json:"headers"`
	ResponseJson    json.RawMessage   `json:"json,omitempty"`
	ResponseError   string            `json:"error,omitempty"`
}

func (z *Res) Apply(res api_response.Response, resErr error) {
	z.ResponseCode = res.StatusCode()
	resBody, _ := res.Result()
	if len(resBody) == 0 {
		z.ResponseBody = ""
	} else if resBody[0] == '[' || resBody[0] == '{' {
		z.ResponseJson = []byte(resBody)
	} else {
		z.ResponseBody = resBody
	}
	if resErr != nil {
		z.ResponseError = resErr.Error()
	}
	z.ResponseHeaders = res.Headers()
}

func (z *captureImpl) WithResponse(req api_request.Request, res api_response.Response, resErr error, latency int64) {
	// request
	rq := Req{}
	rq.Apply(req)

	// response
	rs := Res{}
	rs.Apply(res, resErr)

	z.capture.Debug("",
		zap.Any("req", rq),
		zap.Any("res", rs),
		zap.Int64("latency", latency),
	)
}

func (z *captureImpl) NoResponse(req api_request.Request, resErr error, latency int64) {
	// request
	rq := Req{}
	rq.Apply(req)

	// response
	rs := Res{}
	rs.ResponseError = resErr.Error()

	z.capture.Debug("",
		zap.Any("req", rq),
		zap.Any("res", rs),
		zap.Int64("latency", latency),
	)
}
