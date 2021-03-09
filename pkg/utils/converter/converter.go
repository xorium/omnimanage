package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/jsonapi"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"reflect"
	"strings"
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

func ModelToOutput(model interface{}) (res interface{}, errOut error) {
	// set example values
	err := defaults.Set(model)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = jsonapi.MarshalPayload(&b, model)
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(b.Bytes(), &jsonMap)
	if err != nil {
		return nil, err
	}

	dynStruct := MapToDyn(jsonMap)
	//r := dynamicstruct.NewReader(jsonMap)
	//
	//dyn := dynamicstruct.NewStruct()
	//_ = parseDynamicReader(r, dyn, "Root")
	//val := dyn.Build().New()
	//instance := dynamicstruct.NewStruct().
	//	AddField("Integer", 0, `json:"int"`).
	//	AddField("Text", "", `json:"someText"`).
	//	AddField("Float", 0.0, `json:"double"`).
	//	AddField("Boolean", false, "").
	//	AddField("Slice", []int{}, "").
	//	AddField("Anonymous", "", `json:"-"`).
	//
	//	Build().
	//	New()

	return dynStruct, nil
}

func MapToDyn(input map[string]interface{}) interface{} {
	dyn := dynamicstruct.NewStruct()

	for key, val := range input {
		dyn.AddField(strings.Title(key), valToDyn(val), "")
	}

	val := dyn.Build().New()
	return val
}

func valToDyn(val interface{}) interface{} {
	typeKind := reflect.ValueOf(val).Kind()

	switch typeKind {
	case reflect.Bool, reflect.String, reflect.Int, reflect.Float64:
		return val
	case reflect.Map:
		dyn := dynamicstruct.NewStruct()

		iter := reflect.ValueOf(val).MapRange()
		for iter.Next() {
			tmpIn := iter.Value().Elem().Interface()
			dyn.AddField(strings.Title(iter.Key().String()), valToDyn(tmpIn), "")
		}
		instance := dyn.Build().New()
		return instance
	case reflect.Slice:
		valRef := reflect.ValueOf(val)

		newSlice := make([]interface{}, 1)

		for i := 0; i < valRef.Len(); i++ {
			tmpIn := valRef.Index(i).Elem().Interface()
			rec := valToDyn(tmpIn)

			newSlice = append(newSlice, rec)
		}
		return newSlice

	default:
		fmt.Errorf("unexpected in - %v ", val)
	}

	return nil
}

//func parseDynamicReader(r dynamicstruct.Reader, dyn dynamicstruct.Builder, name string) interface{} {
//
//	//fields := r.GetAllFields()
//	//if fields != nil {
//	//	fields := r.GetAllFields()
//	//	for _, f := range fields {
//	//		dyn.AddField(f.Name(), f.Interface(), "")
//	//	}
//	//} else {
//	dyn.AddField(strings.Title(name), "", "")
//
//	rs := r.ToMapReaderOfReaders()
//	if rs != nil {
//		for key, rItem := range rs {
//
//			parseDynamicReader(rItem, dyn, strings.Title(key.(string)))
//		}
//	} else {
//		//rs := r.ToSliceOfReaders()
//		//if rs != nil {
//		//	for key, rItem := range rs {
//		//		parseDynamicReader(rItem, dyn, )
//		//	}
//		//} else {
//		//	//dyn.AddField(r.V, f.Interface(), "")
//		//}
//	}
//
//	return nil
//}
