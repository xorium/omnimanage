package validator

import (
	"fmt"
	"github.com/go-playground/validator"
	"sync"
)

//type Validator struct {
//	//validator *validator.Validate
//}
//
//func NewValidator() *Validator {
//	return &Validator{}
//}
//
//// ValidateStruct implements the echo framework validator interface.
//func (val *Validator) ValidateStruct(i interface{}) error {
//	//err := val.validator.Struct(i)
//	//if err == nil {
//	//	return nil
//	//}
//	//err = errors.New(strings.Replace(err.Error(), "\n", ", ", -1))
//	//return err
//	return nil
//}

// use a single instance of validate, it caches struct info
var (
	validate *validator.Validate
	once     sync.Once
)

func ValidateStruct(s interface{}) error {
	once.Do(func() {
		validate = validator.New()
	})

	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("%w: %v", err.Error())
	}

	return nil
}
