package mo

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/oesand/mo/internal"
	"github.com/oesand/mo/validate"
)

// FormParam creates a ParamProvider that extracts a form parameter from the HTTP request.
// It supports validation and can be made optional.
func FormParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParamProvider[T] {
	return &formParameter[T]{
		name:       name,
		validators: validators,
		post:       false,
	}
}

// PostFormParam creates a ParamProvider that extracts a POST form parameter from the HTTP request.
// It supports validation and can be made optional.
func PostFormParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParamProvider[T] {
	return &formParameter[T]{
		name:       name,
		validators: validators,
		post:       true,
	}
}

var _ internal.ParamSchema = (*formParameter[string])(nil)

type formParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	post       bool
	validators []validate.Validator[T]
}

func (param *formParameter[T]) Name() string {
	return param.name
}

func (param *formParameter[T]) ParamType() internal.ParamType {
	if param.post {
		return internal.PostFormParamType
	}
	return internal.FormParamType
}

func (param *formParameter[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (param *formParameter[T]) IsRequired() bool {
	return !param.optional
}

func (param *formParameter[T]) Optional() ParamProvider[T] {
	param.optional = true
	return param
}

func (param *formParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	var str string
	if param.post {
		str = request.PostFormValue(param.name)
	} else {
		str = request.FormValue(param.name)
	}

	if str == "" {
		if !param.optional {
			errs = []string{fmt.Sprintf("form param '%s' is required", param.name)}
		}
		return
	}

	val, err := parseBasicTypes[T](str)
	if err != "" {
		errs = []string{fmt.Sprintf("form param '%s' %s", param.name, err)}
		return
	}

	for _, validator := range param.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("form param '%s': %s", param.name, err))
		}
	}
	return val, errs
}

// MultipartFormParam creates a ParamProvider that parses and returns multipart form data.
// The maxMemory parameter controls how much of the form data is stored in memory before
// spilling to temporary files on disk.
func MultipartFormParam(maxMemory int64) ParamProvider[*multipart.Form] {
	return &multipartFormParameter{maxMemory: maxMemory}
}

var _ internal.ParamSchema = (*multipartFormParameter)(nil)

type multipartFormParameter struct {
	maxMemory int64
	optional  bool
}

func (param *multipartFormParameter) Name() string {
	return ""
}

func (param *multipartFormParameter) ParamType() internal.ParamType {
	return internal.MultipartFormParamType
}

func (param *multipartFormParameter) Type() reflect.Type {
	return nil
}

func (param *multipartFormParameter) IsRequired() bool {
	return !param.optional
}

func (param *multipartFormParameter) Optional() ParamProvider[*multipart.Form] {
	param.optional = true
	return param
}

func (param *multipartFormParameter) GetParamValue(request *http.Request) (*multipart.Form, validate.Errors) {
	if request.MultipartForm != nil {
		return request.MultipartForm, nil
	}

	err := request.ParseMultipartForm(param.maxMemory)
	if err != nil {
		if !param.optional {
			return nil, []string{err.Error()}
		}
		return nil, nil
	}
	return request.MultipartForm, nil
}

var DefaultMaxMemory int64 = 32 << 20 // 32 MB

// FileParam creates a ParamProvider that extracts a single file from a multipart form upload.
// The name parameter specifies the form field name to extract the file from. If multiple files
// are uploaded under the same field name, only the first file is returned.
//
// The provider automatically parses the multipart form data using the DefaultMaxMemory limit
// (32 MB). Files smaller than this limit are kept in memory; larger files are written to disk.
//
// Returns an error if the file is required but not provided. Use Optional() to allow
// missing files.
func FileParam(name string) ParamProvider[*multipart.FileHeader] {
	return &multipartFormFileParameter{name: name}
}

var _ internal.ParamSchema = (*multipartFormFileParameter)(nil)

type multipartFormFileParameter struct {
	name     string
	optional bool
}

func (param *multipartFormFileParameter) Name() string {
	return param.name
}

func (param *multipartFormFileParameter) ParamType() internal.ParamType {
	return internal.FileParamType
}

func (param *multipartFormFileParameter) Type() reflect.Type {
	return nil
}

func (param *multipartFormFileParameter) IsRequired() bool {
	return !param.optional
}

// Optional marks the file parameter as optional. If the file is not provided in the request,
// no error is returned. The method returns nil for the file header and an empty error slice.
func (param *multipartFormFileParameter) Optional() ParamProvider[*multipart.FileHeader] {
	param.optional = true
	return param
}

// GetParamValue extracts the file from the multipart form. It parses the multipart form if
// not already parsed. Returns the first file uploaded under the specified field name, or an
// error if the file is required but missing.
func (param *multipartFormFileParameter) GetParamValue(request *http.Request) (*multipart.FileHeader, validate.Errors) {
	form, err := param.getMultipartForm(request)
	if err != nil {
		return nil, validate.Errors{err.Error()}
	}
	if form == nil {
		if !param.optional {
			return nil, validate.Errors{"multipart form is required"}
		}
		return nil, nil
	}

	files, has := form.File[param.name]
	if !has || len(files) == 0 {
		return nil, validate.Errors{fmt.Sprintf("multipart form file '%s' is required", param.name)}
	}
	return files[0], nil
}

func (param *multipartFormFileParameter) getMultipartForm(request *http.Request) (*multipart.Form, error) {
	if request.MultipartForm != nil {
		return request.MultipartForm, nil
	}

	err := request.ParseMultipartForm(DefaultMaxMemory)
	if err != nil {
		if !param.optional {
			return nil, err
		}
		return nil, nil
	}
	return request.MultipartForm, nil
}
