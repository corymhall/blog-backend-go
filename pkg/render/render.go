package render

import (
	"reflect"
	"net/http"
	"github.com/rs/zerolog"
)

// Renderer interface for mapping response payloads
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// Binder type for managing request payloads.
type Binder interface {
	Bind(r *http.Request) error
}

func Bind(r *http.Request, v Binder) error {
	if err := Decode(r, v); err != nil {
		return err
	}

	return binder(r, v)
}

func Render(w http.ResponseWriter, r *http.Request, v Renderer, l zerolog.Logger) error {
	if err := renderer(w, r, v); err != nil {
		return err
	}
	Respond(w, r, v, l)
	return nil
}

// RenderList renders a slice of payloads and responds to the client request.
func RenderList(w http.ResponseWriter, r *http.Request, l []Renderer, log zerolog.Logger) error {
	for _, v := range l {
		if err := renderer(w, r, v); err != nil {
			return err
		}
	}
	Respond(w, r, l, log)
	return nil
}

func isNil(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return f.IsNil()
	default:
		return false
	}
}

// renderer is a helper function that just calls any Render method of the Renderer
// that is passed in. If the Renderer is a struct, it will also check each field of
// the struct to see if that field also implements Renderer in which case it will
// perform these steps again for that field. This allows for you to perform some
// processing or validation on the request before you Respond.
func renderer(w http.ResponseWriter, r *http.Request, v Renderer) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// we call it top down
	// this is calling the render method that has been defined for your
	// renderer type. This method will be executed before the response is made.
	if err := v.Render(w, r); err != nil {
		return err
	}

	// we're done if the renderer isn't a struct object.
	if rv.Kind() != reflect.Struct {
		return nil
	}

	// if it is a struct then we will repeat the process for each field in the struct
	// that way if any of the fields also implements renderer then it's Render method
	// will be caused before we Respond
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if f.Type().Implements(rendererType) {

			if isNil(f) {
				continue
			}

			fv := f.Interface().(Renderer)
			if err := renderer(w, r, fv); err != nil {
				return err
			}

		}
	}

	return nil
}

// binder is a helper function that just calls any Bind method of the Binder
// that is passed in. If the Binder is a struct, it will also check each field of
// the struct to see if that field also implements Binder in which case it will
// perform these steps again for that field. This allows for you to perform some
// processing or validation on the request before you Decode.
func binder(r *http.Request, v Binder) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// Call Binder on non-struct types right away
	// If the binder is not a struct then we don't have
	// to check each field in the struct to see if it implements
	// binder
	if rv.Kind() != reflect.Struct {
		return v.Bind(r)
	}

	// For structs, we call Bind on each field that implements Binder
	// we are esentially repeating the process for each field so that
	// if any fields on the struct implement binder we call it's Bind method
	// after we decode
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if f.Type().Implements(binderType) {

			if isNil(f) {
				continue
			}

			fv := f.Interface().(Binder)
			if err := binder(r, fv); err != nil {
				return err
			}
		}
	}

	// We call it bottom-up
	if err := v.Bind(r); err != nil {
		return err
	}

	return nil
}

var (
	rendererType = reflect.TypeOf(new(Renderer)).Elem()
	binderType   = reflect.TypeOf(new(Binder)).Elem()
)


type contextKey struct {
	name string
}


func (k *contextKey) String() string {
	return "chi render context value " + k.name
}
