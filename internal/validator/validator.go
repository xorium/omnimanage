package validator

type Validator struct {
	//validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{}
}

// Validate implements the echo framework validator interface.
func (val *Validator) Validate(i interface{}) error {
	//err := val.validator.Struct(i)
	//if err == nil {
	//	return nil
	//}
	//err = errors.New(strings.Replace(err.Error(), "\n", ", ", -1))
	//return err
	return nil
}
