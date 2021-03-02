package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"gorm.io/datatypes"
	omniErr "omnimanage/pkg/error"
	"reflect"
	"strconv"
	"strings"
)

const (
	MethodNameToWeb       = "ToWeb"
	MethodNameScanFromWeb = "ScanFromWeb"
)

const (
	ConverterToSrcTag = "src"
	ConverterToWebTag = "web"
)

type ModelMap struct {
	SrcName        string
	WebName        string
	ConverterToSrc func(web interface{}) (interface{}, error)
	ConverterToWeb func(src interface{}) (interface{}, error)
}

type ModelMapper struct {
	customFunc map[string]func(model interface{}) (interface{}, error)
}

func NewModelMapper() *ModelMapper {
	m := &ModelMapper{customFunc: make(map[string]func(model interface{}) (interface{}, error))}

	// TODO перенести куда-нибудь
	// default converters
	{
		m.RegisterCustomConverter(
			"ID2src",
			func(web interface{}) (interface{}, error) {
				w, ok := web.(string)
				if !ok {
					return 0, fmt.Errorf("ID: Wrong type '%T', value %v", web, web)
				}
				id, err := strconv.Atoi(w)
				if err != nil {
					return 0, err
				}
				return id, nil

			},
		)
		m.RegisterCustomConverter(
			"ID2web",
			func(src interface{}) (interface{}, error) {
				s, ok := src.(int)
				if !ok {
					return "", fmt.Errorf("ID: Wrong type '%T', value %v", src, src)
				}
				id := strconv.Itoa(s)
				return id, nil
			},
		)
		m.RegisterCustomConverter(
			"JSON2src",
			func(web interface{}) (interface{}, error) {
				w, ok := web.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("Wrong type '%T'", web)
				}

				j, err := json.Marshal(w)
				if err != nil {
					return nil, err
				}

				return j, nil
			},
		)
		m.RegisterCustomConverter(
			"JSON2web",
			func(src interface{}) (interface{}, error) {
				s, ok := src.(datatypes.JSON)
				if !ok {
					return nil, fmt.Errorf("Wrong type '%T'", src)
				}
				if len(s) == 0 {
					s = []byte("{}")
				}
				w := map[string]interface{}{}
				err := json.Unmarshal(s, &w)
				if err != nil {
					return nil, err
				}
				return w, nil
			},
		)
	}
	return m
}

func (m *ModelMapper) RegisterCustomConverter(tagName string, f func(model interface{}) (interface{}, error)) {
	m.customFunc[tagName] = f
}

func (m *ModelMapper) GetModelMaps(srcModel interface{}) (modelMaps []*ModelMap, errOut error) {
	srcType := reflect.TypeOf(srcModel)

	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input parameter is not a struct")
	}

	modelMaps = make([]*ModelMap, 0, 5)
	for i := 0; i < srcType.NumField(); i++ {
		fieldType := srcType.Field(i)
		tag := fieldType.Tag.Get("omni")
		if tag == "" {
			continue
		}

		args := strings.Split(tag, ";")
		if len(args) < 1 {
			return nil, fmt.Errorf("field %v: bad tag value: %v", fieldType.Name, tag)
		}

		// add new map
		newMap := &ModelMap{
			SrcName: fieldType.Name,
			WebName: args[0],
		}

		// fill converters
		for j := 1; j < len(args); j++ {
			converterSl := strings.Split(args[j], ":")
			if len(converterSl) != 2 {
				return nil, fmt.Errorf("field %v: bad tag value: %v", fieldType.Name, tag)
			}

			conv, ok := m.customFunc[converterSl[1]]
			if !ok {
				return nil, fmt.Errorf("field %v: unknown converter %v", fieldType.Name, converterSl[1])
			}

			switch converterSl[0] {
			case ConverterToSrcTag:
				newMap.ConverterToSrc = conv
			case ConverterToWebTag:
				newMap.ConverterToWeb = conv
			default:
				return nil, fmt.Errorf("field %v: bad converter tag value: %v", fieldType.Name, tag)
			}
		}

		modelMaps = append(modelMaps, newMap)
	}

	return modelMaps, nil
}

//type ISrcModel interface {
//	//GetModelMapper() []*ModelMap
//}

func GetModelMapBySrcName(name string, m []*ModelMap) *ModelMap {
	for _, val := range m {
		if val.SrcName == name {
			return val
		}
	}
	return nil
}

func GetModelMapByWebName(name string, m []*ModelMap) *ModelMap {
	for _, val := range m {
		if val.WebName == name {
			return val
		}
	}
	return nil
}

func (m *ModelMapper) ConvertSrcToWeb(src interface{}, web interface{}) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	webS := structs.New(web)
	srcS := structs.New(src)
	srcM := srcS.Map()

	//modelMaps := src.GetModelMapper()
	modelMaps, err := m.GetModelMaps(src)
	if err != nil {
		return err
	}

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
				resInv, err := CallMethodWith2Output(srcField.Value(), MethodNameToWeb, m)
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

func (m *ModelMapper) ConvertWebToSrc(web interface{}, src interface{}) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	srcS := structs.New(src)
	webS := structs.New(web)
	webM := webS.Map()

	//modelMaps := src.GetModelMapper()
	modelMaps, err := m.GetModelMaps(src)
	if err != nil {
		return err
	}
	for _, val := range modelMaps {
		if val.SrcName == "" {
			continue
		}

		webFieldValue, ok := webM[val.WebName]
		if !ok {
			return fmt.Errorf("unknown web field %v", val.WebName)
		}

		srcField, ok := srcS.FieldOk(val.SrcName)
		if !ok {
			return fmt.Errorf("unknown src field %v", val.SrcName)
		}

		if val.ConverterToSrc == nil { // no converter function -> simple conversion
			typeKind := srcField.Kind()

			// Relation
			if typeKind == reflect.Ptr || typeKind == reflect.Slice {
				webField := webS.Field(val.WebName)
				if webField.IsZero() {
					continue
				}

				resInv, err := CallMethodWith2Output(srcField.Value(), MethodNameScanFromWeb, webField.Value(), m)
				if err != nil {
					return fmt.Errorf("Error in converting %v : %v", val.WebName, err)
				}
				err = srcField.Set(resInv.Interface())
				if err != nil {
					return fmt.Errorf("Error in converting %v: %v", val.WebName, err)
				}
			} else { // Simple Attribute
				err := srcField.Set(webFieldValue)
				if err != nil {
					return fmt.Errorf("Error in converting %v: %v", val.WebName, err)
				}
			}
			continue
		}

		if val.ConverterToSrc != nil { //with converter function
			field, err := val.ConverterToSrc(webFieldValue)
			if err != nil {
				return err
			}
			srcField.Set(field)
			continue
		}

		return fmt.Errorf("wrong mapper line %v", val)
	}
	return nil
}

//func GetMapperDynamic(t reflect.Type) (out []*ModelMap, errOut error) {
//	defer func() {
//		if r := recover(); r != nil {
//			errOut = fmt.Errorf("panic: %v", r)
//		}
//	}()
//
//	ptr := reflect.New(t)
//
//	method := ptr.Elem().MethodByName(MethodNameModelMapper)
//	if !method.IsValid() {
//		return nil, fmt.Errorf("Method %v not exists", MethodNameModelMapper)
//	}
//
//	results := method.Call([]reflect.Value{})
//	maps, ok := results[0].Interface().([]*ModelMap)
//	if !ok {
//		return nil, fmt.Errorf("Internal error in method %v", MethodNameModelMapper)
//	}
//	return maps, nil
//}

func (m *ModelMapper) GetSrcID(webID string, srcModel interface{}) (idOut int, errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("%w: panic - %v", omniErr.ErrInternal, r)
		}
	}()

	modelMaps, err := m.GetModelMaps(srcModel)
	if err != nil {
		return 0, err
	}

	idMap := GetModelMapBySrcName("ID", modelMaps) //srcModel.GetModelMapper())
	if idMap == nil {
		return -1, fmt.Errorf("%w: map ID not found", omniErr.ErrInternal)
	}

	if idMap.ConverterToSrc == nil {
		return -1, fmt.Errorf("%w: Converter not found", omniErr.ErrInternal)
	}
	srcID, err := idMap.ConverterToSrc(webID)
	if idMap == nil {
		return -1, fmt.Errorf("%w: %v", omniErr.ErrBadRequest, err)
	}
	idOut, ok := srcID.(int)
	if !ok {
		return -1, fmt.Errorf("%w: wrong ID type", omniErr.ErrInternal)
	}
	return idOut, nil
}
