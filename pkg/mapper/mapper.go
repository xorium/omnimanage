package mapper

import (
	"fmt"
	"github.com/fatih/structs"
	"reflect"
)

const (
	MethodNameToWeb       = "ToWeb"
	MethodNameModelMapper = "GetModelMapper"
)

type ModelMapper struct {
	SrcName        string
	WebName        string
	ConverterToSrc func(web interface{}) (interface{}, error)
	ConverterToWeb func(src interface{}) (interface{}, error)
}

type ISrcModel interface {
	GetModelMapper() []*ModelMapper
}

func GetModelMapBySrcName(name string, m []*ModelMapper) *ModelMapper {
	for _, val := range m {
		if val.SrcName == name {
			return val
		}
	}
	return nil
}

func GetModelMapByWebName(name string, m []*ModelMapper) *ModelMapper {
	for _, val := range m {
		if val.WebName == name {
			return val
		}
	}
	return nil
}

func ConvertSrcToWeb(src ISrcModel, web interface{}) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	webS := structs.New(web)
	srcS := structs.New(src)
	srcM := srcS.Map()

	modelMaps := src.GetModelMapper()
	for _, val := range modelMaps {
		if val.WebName == "" {
			continue
		}

		srcFieldValue, ok := srcM[val.SrcName]
		if !ok {
			return fmt.Errorf("unknown src field %v", val.SrcName)
		}

		webField, ok := webS.FieldOk(val.WebName)
		if !ok {
			return fmt.Errorf("unknown web field %v", val.WebName)
		}

		if val.ConverterToWeb == nil { // no converter function -> simple conversion
			typeKind := webField.Kind()

			// Relation
			if typeKind == reflect.Ptr || typeKind == reflect.Slice {
				srcField := srcS.Field(val.SrcName)
				if srcField.IsZero() {
					continue
				}
				resInv, err := CallMethodWith2Output(srcField.Value(), MethodNameToWeb)
				if err != nil {
					return fmt.Errorf("Error in converting %v : %v", val.SrcName, err)
				}
				err = webField.Set(resInv.Interface())
				if err != nil {
					return fmt.Errorf("Error in converting %v: %v", val.SrcName, err)
				}
			} else { // Simple Attribute
				err := webField.Set(srcFieldValue)
				if err != nil {
					return fmt.Errorf("Error in converting %v: %v", val.SrcName, err)
				}
			}
			continue
		}

		if val.ConverterToWeb != nil { //with converter function
			srcField, err := val.ConverterToWeb(srcFieldValue)
			if err != nil {
				return err
			}
			webField.Set(srcField)
			continue
		}

		return fmt.Errorf("wrong mapper line %v", val)
	}
	return nil
}

// CallMethodWith2Output calls method of structure Any. Must be 2 output params - result value and error
func CallMethodWith2Output(any interface{}, name string, args ...interface{}) (out reflect.Value, errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	method := reflect.ValueOf(any).MethodByName(name)
	if !method.IsValid() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %v not exists", name)
	}
	methodType := method.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}
	results := method.Call(in)
	if len(results) != 2 {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have 2 output parameters", name)
	}
	err, ok := results[1].Interface().(error)
	if ok {
		return reflect.ValueOf(nil), err
	}
	return results[0], nil
}

func GetMapperDynamic(t reflect.Type) (out []*ModelMapper, errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	ptr := reflect.New(t)

	method := ptr.Elem().MethodByName(MethodNameModelMapper)
	if !method.IsValid() {
		return nil, fmt.Errorf("Method %v not exists", MethodNameModelMapper)
	}

	results := method.Call([]reflect.Value{})
	maps, ok := results[0].Interface().([]*ModelMapper)
	if !ok {
		return nil, fmt.Errorf("Internal error in method %v", MethodNameModelMapper)
	}
	return maps, nil
}
