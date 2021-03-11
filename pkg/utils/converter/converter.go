package converter

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/jsonapi"
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

func GetExampleModelListSwagOutput(model interface{}) (res interface{}, errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	// set example values
	err := defaults.Set(model)
	if err != nil {
		return nil, err
	}

	modelList := make([]interface{}, 0, 1)
	modelList = append(modelList, model)

	p, err := jsonapi.Marshal(modelList)
	if err != nil {
		return nil, err
	}
	manyPay, ok := p.(*jsonapi.ManyPayload)
	if !ok {
		return nil, fmt.Errorf("wrong model input")
	}
	return manyPay, nil
}

func GetExampleModelSwagOutput(model interface{}) (res interface{}, errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	// set example values
	err := defaults.Set(model)
	if err != nil {
		return nil, err
	}

	p, err := jsonapi.Marshal(model)
	if err != nil {
		return nil, err
	}
	onePay, ok := p.(*jsonapi.OnePayload)
	if !ok {
		return nil, fmt.Errorf("wrong model input")
	}
	return onePay, nil
}

//func MapToDyn(input map[string]interface{}) interface{} {
//	dyn := dynamicstruct.NewStruct()
//
//	for key, val := range input {
//		tmpTag := `json:"` + strings.ToLower(key) + `"`
//		tmpVal := valToDyn(val)
//
//		//jsDef, err := json.Marshal(tmpVal)
//		//if err != nil {
//		//	return nil
//		//}
//		//strDef := strings.Replace(string(jsDef), `"`, `\"`, -1)
//		//
//		//tmpTag = tmpTag + ` default:"` + strDef + `"`
//
//		dyn.AddField(strings.Title(key), tmpVal, tmpTag)
//	}
//
//	val := dyn.Build().New()
//
//	r := dynamicstruct.NewReader(val)
//	r.GetAllFields()
//
//	err := defaults.Set(val)
//	if err != nil {
//		return nil
//	}
//
//	fmt.Printf("%+v", val)
//
//	return val
//}
//
//func valToDyn(val interface{}) interface{} {
//	typeKind := reflect.ValueOf(val).Kind()
//
//	switch typeKind {
//	case reflect.Bool, reflect.String, reflect.Int, reflect.Float64:
//		return val
//	case reflect.Map:
//		dyn := dynamicstruct.NewStruct()
//
//		iter := reflect.ValueOf(val).MapRange()
//		for iter.Next() {
//			if iter.Value().IsZero() {
//				continue
//			}
//			tmpIn := iter.Value().Elem().Interface()
//
//			tmpTag := `json:"` + strings.ToLower(iter.Key().String()) + `" `
//			tmpVal := valToDyn(tmpIn)
//			if reflect.TypeOf(tmpVal).Kind() == reflect.Map {
//
//			}
//			dyn.AddField(strings.Title(iter.Key().String()), tmpVal, tmpTag)
//		}
//
//		instance := dyn.Build().New()
//		return instance
//
//	case reflect.Slice:
//		if val == nil {
//			return nil
//		}
//		valRef := reflect.ValueOf(val)
//
//		newSlice := make([]interface{}, 1)
//
//		for i := 0; i < valRef.Len(); i++ {
//			tmpIn := valRef.Index(i).Elem().Interface()
//			rec := valToDyn(tmpIn)
//
//			newSlice = append(newSlice, rec)
//		}
//		return newSlice
//
//	default:
//		fmt.Errorf("unexpected in - %v ", val)
//	}
//
//	return nil
//}
//
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
