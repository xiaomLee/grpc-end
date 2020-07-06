package grpc_end

import (
	"context"
	"math"
	"strconv"

	json "github.com/json-iterator/go"
)

const abortIndex int8 = math.MaxInt8 / 2

type GRpcContext struct {
	handlers HandlersChain // Execution chain
	index    int8          // index of current execution handler

	ctx    context.Context
	engine *GRpcEngine

	req  *Request
	resp *Response

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	// appName, for log
	appName string
}

// ---------------------------------------------------------------------------------------------------------------------

func (c *GRpcContext) reset() {
	c.handlers = nil
	c.index = -1
	c.ctx = nil
	c.engine = nil
	c.req = nil
	c.resp = &Response{}
	c.Keys = nil
}

func (c *GRpcContext) GetContext() context.Context {
	return c.ctx
}

// GetRequest returns GRpc Request
func (c *GRpcContext) GetRequest() *Request {
	return c.req
}

// GetResponse returns GRpc Response of this request
func (c *GRpcContext) GetResponse() *Response {
	return c.resp
}

// GetAppName returns the appName
func (c *GRpcContext) GetAppName() string {
	return c.appName
}

// GetFiles returns the files of this request
func (c *GRpcContext) GetFiles() map[string][]byte {
	return c.req.Files
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes c.Keys if it was not used previously.
func (c *GRpcContext) Set(key string, val interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = val
}

// Get returns the value for the given key.
// If the value does not exists it returns nil
func (c *GRpcContext) Get(key string) interface{} {
	if c.Keys == nil {
		return nil
	}
	return c.Keys[key]
}

// GetStringMap returns string for the given key
// If the value does not exists it returns empty string
func (c *GRpcContext) GetString(key string) string {
	if c.Keys == nil {
		return ""
	}
	if val, ok := c.Keys[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}

	return ""
}

// GetStringMap returns the map[string]string for the given key
// If the value does not exists it returns nil
func (c *GRpcContext) GetStringMap(key string) map[string]string {
	if c.Keys == nil {
		return nil
	}
	if val, ok := c.Keys[key]; ok {
		if m, ok := val.(map[string]string); ok {
			return m
		}
	}

	return nil
}

// IsParamExist returns true is param exist
func (c *GRpcContext) IsParamExist(key string) bool {
	_, ok := c.req.Params[key]
	return ok
}

// StringParam returns string val from request's Params for the given key
// and returns empty string if val not exists
func (c *GRpcContext) StringParam(key string) string {
	val := c.req.Params[key]
	if val == "" {
		val = c.req.Header[key]
	}
	return val
}

// StringParamDefault returns string val from request's Params for the given key
// and returns defVal if val not exists.
func (c *GRpcContext) StringParamDefault(key string, defVal string) string {
	if val := c.StringParam(key); val != "" {
		return val
	}
	return defVal
}

// IntParam returns int val from request's Params for the given key
// and returns zero if val not exists.
func (c *GRpcContext) IntParam(key string) int {
	if val := c.StringParam(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return 0
}

// IntParamDefault returns int val from request's Params for the given key
// and returns defVal if val not exists.
func (c *GRpcContext) IntParamDefault(key string, defVal int) int {
	if val := c.StringParam(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return defVal
}

// Int64Param returns int64 val from request's Params for the given key
// and returns zero if val not exists.
func (c *GRpcContext) Int64Param(key string) int64 {
	if val := c.StringParam(key); val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

// Int64ParamDefault returns int64 val from request's Params for the given key
// and returns defVal if val not exists.
func (c *GRpcContext) Int64ParamDefault(key string, defVal int64) int64 {
	if val := c.StringParam(key); val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			return n
		}
	}
	return defVal
}

// Float64Param returns float64 val from request's Params for the given key
// and returns zero if val not exists.
func (c *GRpcContext) Float64Param(key string) float64 {
	if val := c.StringParam(key); val != "" {
		if n, err := strconv.ParseFloat(val, 64); err == nil {
			return n
		}
	}
	return 0
}

// Float64ParamDefault returns float64 val from request's Params for the given key
// and returns defVal if val not exists.
func (c *GRpcContext) Float64ParamDefault(key string, defVal float64) float64 {
	if val := c.StringParam(key); val != "" {
		if n, err := strconv.ParseFloat(val, 64); err == nil {
			return n
		}
	}
	return defVal
}

// ---------------------------------------------------------------------------------------------------------------------

// HeaderString returns string val from request's Header for the given key
// and returns empty string if val not exists
//
// The GateWay will fill in some val to request's Header, like:
// --------------------------------------------------------------------
// | key     | desc                                                   |
// --------------------------------------------------------------------
// | ip      | client's ip address                                    |
// | lang    | the Language client use, ie: 'zh', 'en', 'ko'...       |
// | device  | the Device client use, ie: 'iphone 7 Plus', 'chrome'   |
// | dt      | 'N' mean web, 'a' mean android, 'i' mean ios           |                         |
// | host    | the host of this request belong                        |
// --------------------------------------------------------------------

var (
	HeaderKeyIp     = "ip"
	HeaderKeyLang   = "lang"
	HeaderKeyDevice = "device"
	HeaderKeyDt     = "dt"
	HeaderKeyHost   = "host"
)

func (c *GRpcContext) StringHeader(key string) string {
	return c.req.Header[key]
}

// StringHeaderDefault returns string val from request's Header for the given key
// and returns defVal if val not exists
func (c *GRpcContext) StringHeaderDefault(key string, defVal string) string {
	if _, ok := c.req.Header[key]; !ok {
		return defVal
	}

	return c.StringHeader(key)
}

// IntHeader returns int val from request's Header for the given key
// and returns zero if val not exists
func (c *GRpcContext) IntHeader(key string) int {
	if val := c.StringHeader(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return 0
}

// Int64Header returns int64 val from request's Header for the given key
// and returns zero if val not exists
func (c *GRpcContext) Int64Header(key string) int64 {
	if val := c.StringHeader(key); val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

// ---------------------------------------------------------------------------------------------------------------------

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (c *GRpcContext) Next() {
	c.index++
	for s := int8(len(c.handlers)); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *GRpcContext) Abort() {
	c.index = abortIndex
}

// IsAborted returns true if the current context was aborted.
func (c *GRpcContext) IsAbort() bool {
	return c.index >= abortIndex
}

// ---------------------------------------------------------------------------------------------------------------------

type SResponse struct {
	Success bool        `json:"success"`
	PayLoad interface{} `json:"payload"`
}

type EResponse struct {
	Success bool  `json:"success"`
	Err     Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *GRpcContext) SuccessResponse(v interface{}) {
	if v == nil {
		v = make(map[string]interface{})
	}
	jsonStr, _ := json.Marshal(&SResponse{
		Success: true,
		PayLoad: v,
	})
	c.resp.Data = jsonStr
}

func (c *GRpcContext) ErrResponse(code int, err error) {
	jsonStr, _ := json.Marshal(&EResponse{
		Success: false,
		Err: Error{
			Code:    code,
			Message: err.Error(),
		},
	})
	c.resp.Data = jsonStr
}
