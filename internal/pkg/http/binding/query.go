package binding

import (
	"net/http"

	"github.com/go-playground/form/v4"
)

var formDecoder = form.NewDecoder()

type queryBinding struct{}

func (queryBinding) Name() string {
	return "URL-QUERY"
}

func (queryBinding) Bind(req *http.Request, obj any) error {
	values := req.URL.Query()

	if err := formDecoder.Decode(obj, values); err != nil {
		return &bindingError{msg: "incorrect parameters"}
	}

	if err := validate(obj); err != nil {
		return &bindingError{msg: err.Error()}
	}

	return nil
}
