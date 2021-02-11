package filters

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"gorm.io/gorm"
	"omnimanage/pkg/mapper"
	"reflect"
	"strconv"
	"strings"
)

//var queryFilterRegex = regexp.MustCompile(`^filters\[(\w+)\]$`)

type Filter struct {
	Relation string
	Field    string
	Operator string
	Value    interface{}
}

const (
	JSONAPIRelationTagPrefix  = "relation,"
	JSONAPIAttributeTagPrefix = "attr,"
	JSONAPIIdFieldName        = "id"
)

func SetGormFilters(db *gorm.DB, model interface{}, f []*Filter) (*gorm.DB, error) {

	if f == nil {
		return db, nil
	}

	err := db.Statement.Parse(model)
	if err != nil {
		return nil, err
	}

	for _, filter := range f {
		dbFieldName := db.Statement.NamingStrategy.ColumnName("", filter.Field)
		if filter.Relation != "" {
			rel, ok := db.Statement.Schema.Relationships.Relations[filter.Relation]
			if !ok {
				return nil, errors.New("wrong filter")
			}

			joinStr := fmt.Sprintf("JOIN %v on %v.%v=%v.%v and %v.%v%v?",
				rel.FieldSchema.Table,
				rel.FieldSchema.Table,
				rel.References[0].PrimaryKey.DBName,
				rel.Schema.Table,
				rel.References[0].ForeignKey.DBName,
				rel.FieldSchema.Table,
				dbFieldName,
				filter.Operator,
			)

			db = db.Joins(joinStr, filter.Value)
		} else {
			whereStr := fmt.Sprintf("%v%v?", dbFieldName, filter.Operator) //TODO security? sql injection!
			db = db.Where(whereStr, filter.Value)
		}
	}
	return db, nil

}

func buildFieldByTagMap(tagKey string, srcType reflect.Type) (map[string]string, error) {
	src := srcType
	if src.Kind() == reflect.Ptr {
		src = srcType.Elem()
	}
	if src.Kind() != reflect.Struct {
		return nil, fmt.Errorf("bad type")
	}
	fieldsByTag := make(map[string]string)

	for i := 0; i < src.NumField(); i++ {
		f := src.Field(i)
		v := f.Tag.Get(tagKey) //strings.Split(f.Tag.Get(tagKey), ",")[0]
		if v == "" || v == "-" {
			continue
		}
		fieldsByTag[v] = f.Name
	}
	return fieldsByTag, nil
}

func TransformWebToSrc(f []*Filter, modelWeb interface{}, modelSrc mapper.ISrcModel) (out []*Filter, errOut error) {
	if f == nil {
		return nil, nil
	}

	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	refWeb := reflect.TypeOf(modelWeb)

	out = make([]*Filter, len(f), cap(f))

	tagMap, err := buildFieldByTagMap("jsonapi", reflect.TypeOf(modelWeb))
	if err != nil {
		return nil, err
	}

	srcS := structs.New(modelSrc)
	srcM := srcS.Map()

	modelSrcMaps := modelSrc.GetModelMapper()

	for i, filt := range f {
		var fNew *Filter

		fNewTmp := &Filter{}
		tmpVal := reflect.ValueOf(&fNewTmp.Value)

		var ok bool
		var webFieldName string
		var webField reflect.StructField

		if filt.Relation != "" {
			// get web relation field name by tag
			relWebFieldName, ok := tagMap[JSONAPIRelationTagPrefix+filt.Relation]
			if !ok {
				return nil, fmt.Errorf("cant find Field %v", filt.Field)
			}

			relationRef, ok := refWeb.FieldByName(relWebFieldName)
			if !ok {
				return nil, fmt.Errorf("cant find Relation name %v", filt.Relation)
			}

			relTagMap, err := buildFieldByTagMap("jsonapi", relationRef.Type)
			if err != nil {
				return nil, err
			}

			// process attr field
			if filt.Field != JSONAPIIdFieldName {
				webFieldName, ok = relTagMap[JSONAPIAttributeTagPrefix+filt.Field]
				if !ok {
					return nil, fmt.Errorf("cant find Field %v", filt.Field)
				}
			} else {
				webFieldName = "ID"
			}

			webField, ok = relationRef.Type.Elem().FieldByName(webFieldName)
			if !ok {
				return nil, fmt.Errorf("cant find Field %v", filt.Field)
			}

			err = setFilterVal(tmpVal, webField.Type, filt.Value.(string))
			if err != nil {
				return nil, fmt.Errorf("error with Field %v : %v", filt.Field, err)
			}

			// get model mapper by relation field name
			modelRelFieldMapper := mapper.GetModelMapByWebName(relWebFieldName, modelSrcMaps)
			if modelRelFieldMapper == nil {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			srcRelModel, ok := srcM[modelRelFieldMapper.SrcName]
			if !ok {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			srcModel, ok := srcRelModel.(mapper.ISrcModel)
			if !ok {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			modelSrcRelMaps := srcModel.GetModelMapper()

			// Web field name and value -> Src field name and value
			modelFieldMapper := mapper.GetModelMapByWebName(webFieldName, modelSrcRelMaps)
			if modelFieldMapper == nil {
				return nil, fmt.Errorf("Cannot find map for field %v", webFieldName)
			}

			var srcVal interface{}
			if modelFieldMapper.ConverterToSrc != nil {
				srcVal, err = modelFieldMapper.ConverterToSrc(tmpVal.Elem().Interface())
				if err != nil {
					return nil, fmt.Errorf("error with Field %v : %v", filt.Field, err)
				}
			} else {
				srcVal = tmpVal.Elem().Interface()
			}

			fNew = &Filter{
				Relation: modelRelFieldMapper.SrcName,
				Operator: filt.Operator,
				Field:    modelFieldMapper.SrcName,
				Value:    srcVal,
			}

		} else {
			// get web field name by tag
			if filt.Field != JSONAPIIdFieldName {
				webFieldName, ok = tagMap[JSONAPIAttributeTagPrefix+filt.Field]
				if !ok {
					return nil, fmt.Errorf("cant find Field %v", filt.Field)
				}
			} else {
				webFieldName = "ID"
			}

			// write filt.Value to tmpVal(type interface{})
			webField, ok = refWeb.FieldByName(webFieldName)
			if !ok {
				return nil, fmt.Errorf("cant find Field %v", filt.Field)
			}

			err = setFilterVal(tmpVal, webField.Type, filt.Value.(string))
			if err != nil {
				return nil, fmt.Errorf("error with Field %v : %v", filt.Field, err)
			}

			// Web field name and value -> Src field name and value
			modelFieldMapper := mapper.GetModelMapByWebName(webFieldName, modelSrcMaps)
			if modelFieldMapper == nil {
				return nil, fmt.Errorf("Cannot find map for field %v", webFieldName)
			}

			var srcVal interface{}
			if modelFieldMapper.ConverterToSrc != nil {
				srcVal, err = modelFieldMapper.ConverterToSrc(tmpVal.Elem().Interface())
				if err != nil {
					return nil, fmt.Errorf("error with Field %v : %v", filt.Field, err)
				}
			} else {
				srcVal = tmpVal.Elem().Interface()
			}

			fNew = &Filter{
				Operator: filt.Operator,
				Field:    modelFieldMapper.SrcName,
				Value:    srcVal,
			}

		}

		out[i] = fNew
	}

	return out, nil
}

func setFilterVal(filt reflect.Value, fieldType reflect.Type, newVal string) error {
	typeKind := fieldType.Kind()
	switch typeKind {
	case reflect.Bool:
		newVal, err := strconv.ParseBool(newVal)
		if err != nil {
			return fmt.Errorf("cant parse field to bool type")
		}
		filt.Elem().Set(reflect.ValueOf(newVal))

	case reflect.String:
		filt.Elem().Set(reflect.ValueOf(newVal))

	case reflect.Int:
		newVal, err := strconv.Atoi(newVal)
		if err != nil {
			return fmt.Errorf("cant parse field to int type")
		}
		filt.Elem().Set(reflect.ValueOf(newVal))

	case reflect.Float64:
		newVal, err := strconv.ParseFloat(newVal, 64)
		if err != nil {
			return fmt.Errorf("cant parse field to int type")
		}
		filt.Elem().Set(reflect.ValueOf(newVal))
	default:
		return fmt.Errorf("unexpected type ")
	}

	return nil
}

func GetFiltersFromQueryString(queryStr string, modelWeb interface{}) ([]*Filter, error) {
	if !strings.HasPrefix(queryStr, "filter=") {
		return nil, nil
	}
	query := strings.TrimPrefix(queryStr, "filter=")

	tagMap, err := buildFieldByTagMap("jsonapi", reflect.TypeOf(modelWeb))
	if err != nil {
		return nil, err
	}

	filters := make([]*Filter, 0, 1)

	parts1 := strings.Split(query, "||")
	for _, p1 := range parts1 {
		parts2 := strings.Split(p1, "&&")
		for _, p2 := range parts2 {
			fNew := &Filter{}

			i := strings.Index(p2, ".")
			if i > 0 {
				fNew.Relation = p2[:i]
				p2 = p2[i+1:]
			}

			operator, operIndex := GetOperator(p2)
			if operIndex < 0 {
				return nil, fmt.Errorf("Wrong query parameter - cant find operator in %v", p2)
			}
			fNew.Operator = operator

			fNew.Field = p2[:operIndex]
			_, isRelation := tagMap[JSONAPIRelationTagPrefix+fNew.Field]
			if isRelation {
				fNew.Relation = fNew.Field
				fNew.Field = JSONAPIIdFieldName
			}
			fNew.Value = p2[operIndex+len(operator):]
			filters = append(filters, fNew)

		}
	}

	return filters, nil
}

func GetOperator(s string) (string, int) {
	operators := []string{
		"<=", ">=", "!=", ">", "<", "=",
	}
	for _, o := range operators {
		index := strings.Index(s, o)
		if index > 0 {
			return o, index
		}
	}

	return "", -1
}
