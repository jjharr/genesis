package validate

import (
	"bytes"
	"errors"
	"fmt"
)

/*
TODO Validation upgrades from asaskevich package:

- permanent custom messages per validation key (currently bails after first error on a field)
- To preserve generic error interface, offer different validation method that returns that
- add support for db-linked validations: exists, unique
- add other laravel validations: in,not in, before/after, active URL
- file validations? dimensions, image
- bigger variety of built-in char set validations (name, title)
- conditional validations?
- nicer array validations for simple methods
- I don't like the set required func at all. It should take a varargs of fields or something
*/

// Error encapsulates a name, an error
type Error struct {
	Name  string
	Field string
	Err   error
}

func (e Error) Error() string {
	if len(e.Name) == 0 {
		return e.Err.Error()
	}
	return e.Name + ": " + e.Err.Error()
}

type ErrorBag struct {
	errors map[string][]Error
}

func NewErrorBag() *ErrorBag {

	eb := make(map[string][]Error)

	return &ErrorBag{eb}
}

func (eb *ErrorBag) Errors() []Error {

	es := make([]Error, 0)

	for _, v := range (*eb).errors {
		es = append(es, v...)
	}

	return es
}

func (eb *ErrorBag) Error() string {
	var err string
	for _, slc := range (*eb).errors {
		for _, e := range slc {
			err += e.Error() + ";"
		}
	}
	return err
}

func (eb *ErrorBag) ErrorMap() map[string][]Error {
	return (*eb).errors
}

func (eb *ErrorBag) Add(key, err, field string) *ErrorBag {
	return eb.addError(Error{
		Name:  key,
		Err:   errors.New(err),
		Field: field,
	})
}

func (eb *ErrorBag) addError(err Error) *ErrorBag {

	slc, exists := (*eb).errors[err.Name]
	if !exists {
		arr := make([]Error, 1)
		arr[0] = err
		(*eb).errors[err.Name] = arr
	} else {
		(*eb).errors[err.Name] = append(slc, err)
	}

	return eb
}

func (eb *ErrorBag) HasErrors() bool {
	return len((*eb).errors) > 0
}

func (eb *ErrorBag) String() string {
	res := new(bytes.Buffer)
	for _, v := range eb.errors {
		for _, e := range v {
			if res.Len() > 0 {
				_, _ = res.WriteString("; ")
			}
			if e.Field == "" {
				_, _ = res.WriteString(fmt.Sprintf("[%s] %s", e.Field, e.Error()))
			} else {
				_, _ = res.WriteString(e.Error())
			}
		}
	}
	return res.String()
}

func (eb *ErrorBag) HasErrorFor(key string) bool {
	_, ok := (*eb).errors[key]
	return ok
}

func (eb *ErrorBag) GetErrorsFor(key string) []Error {
	val, ok := (*eb).errors[key]
	if ok {
		return val
	}

	return nil
}
