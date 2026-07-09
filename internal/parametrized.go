package internal

import (
	"net/http"
	"reflect"
)

type ParamType int

const (
	UndefinedParamType ParamType = iota
	PathParamType
	HeaderParamType
	CookieParamType
	FormParamType
	PostFormParamType
	MultipartFormParamType
	FileParamType
	RawBodyParamType
	JsonBodyParamType
)

type ParamSchema interface {
	Name() string
	ParamType() ParamType
	Type() reflect.Type
	IsRequired() bool
}

type ParamHandler interface {
	http.Handler
	Params() []ParamSchema
}

type CompiledRoute interface {
	PathParams() []string
	ClearedPattern() string
}
