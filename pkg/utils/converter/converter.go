package converter

import (
	"fmt"
	"reflect"
)

// SliceI2SliceModel converts []interface{} to []model
func SliceI2SliceModel(srcSlice []interface{}, out interface{}) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	if reflect.TypeOf(out).Elem().Kind() != reflect.Slice {
		return fmt.Errorf("Out param is not a slice")
	}

	outRef := reflect.ValueOf(out)
	outSlice := reflect.MakeSlice(reflect.ValueOf(out).Elem().Type(), len(srcSlice), cap(srcSlice))

	for i, src := range srcSlice {
		outRow := outSlice.Index(i)
		outRow.Set(reflect.ValueOf(src))
	}
	outRef.Elem().Set(outSlice)

	return nil
}
