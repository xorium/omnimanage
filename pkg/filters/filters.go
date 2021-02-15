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
	Relation        string
	Field           string
	CompareOperator string
	Value           interface{}
	LogicalOperator string
}

const (
	JSONAPIRelationTagPrefix  = "relation,"
	JSONAPIAttributeTagPrefix = "attr,"
	JSONAPIIdFieldName        = "id"
)

const (
	FilterOperatorIN = "IN"
)

func ParseFiltersFromQueryToSrcModel(queryStr string, modelWeb interface{}, modelSrc mapper.ISrcModel) ([]*Filter, error) {
	filtersStrings, err := GetFiltersFromQueryString(queryStr, modelWeb)
	if err != nil {
		return nil, err
	}

	srcFilters, err := TransformWebToSrc(filtersStrings, modelWeb, modelSrc)
	if err != nil {
		return nil, err
	}

	return srcFilters, nil
}

func GetFiltersFromQueryString(queryStr string, modelWeb interface{}) ([]*Filter, error) {
	if !strings.HasPrefix(queryStr, "filter=") {
		return nil, nil
	}
	query := strings.TrimPrefix(queryStr, "filter=")

	tagMap, err := buildFieldsMapByTag("jsonapi", reflect.TypeOf(modelWeb))
	if err != nil {
		return nil, err
	}

	filters := make([]*Filter, 0, 1)

	partsOR := strings.Split(query, "||")
	for indexOR, pOR := range partsOR {
		partsAND := strings.Split(pOR, "&&")

		filtersAND := make([]*Filter, 0, 1)
		for indexAND, pAND := range partsAND {

			fNew := &Filter{}

			i := strings.Index(pAND, ".")
			if i > 0 {
				fNew.Relation = pAND[:i]
				pAND = pAND[i+1:]
			}

			operator, operIndex := GetOperator(pAND)
			if operIndex < 0 {
				return nil, fmt.Errorf("Wrong query parameter - cant find operator in %v", pAND)
			}
			fNew.CompareOperator = operator

			fNew.Field = pAND[:operIndex]
			_, isRelation := tagMap[JSONAPIRelationTagPrefix+fNew.Field]
			if isRelation {
				fNew.Relation = fNew.Field
				fNew.Field = JSONAPIIdFieldName
			}

			if fNew.CompareOperator == FilterOperatorIN {
				valueRange := strings.Split(pAND[operIndex+len(operator):], ",")
				fNew.Value = valueRange
			} else {
				fNew.Value = pAND[operIndex+len(operator):]
			}

			fNew.LogicalOperator = "AND"

			//Clear last in parts
			if indexAND == len(partsAND)-1 {
				fNew.LogicalOperator = ""
			}

			filtersAND = append(filtersAND, fNew)
		}

		logicalOperOR := "OR"
		if indexOR == len(partsOR)-1 {
			logicalOperOR = ""
		}
		for indexAND, f := range filtersAND {
			if indexAND == len(filtersAND)-1 {
				f.LogicalOperator = logicalOperOR
			}
			filters = append(filters, f)
		}
	}

	return filters, nil
}

func TransformWebToSrc(filtersIn []*Filter, modelWeb interface{}, modelSrc mapper.ISrcModel) (out []*Filter, errOut error) {
	if filtersIn == nil {
		return nil, nil
	}

	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	webModelRefType := reflect.TypeOf(modelWeb)
	if webModelRefType.Kind() == reflect.Ptr {
		webModelRefType = webModelRefType.Elem()
	}

	out = make([]*Filter, len(filtersIn), cap(filtersIn))

	tagMapJSONAPI, err := buildFieldsMapByTag("jsonapi", reflect.TypeOf(modelWeb))
	if err != nil {
		return nil, err
	}

	srcStruct := structs.New(modelSrc)
	srcStructMap := srcStruct.Map()

	modelSrcMaps := modelSrc.GetModelMapper()

	for i, filterRow := range filtersIn {
		var newFilterRow *Filter

		if filterRow.Relation != "" {
			relWebFieldName, err := getFieldNameByFilterName(filterRow.Relation, true, tagMapJSONAPI)
			if err != nil {
				return nil, err
			}

			relationRef, ok := webModelRefType.FieldByName(relWebFieldName)
			if !ok {
				return nil, fmt.Errorf("cant find Relation name %v", filterRow.Relation)
			}

			relTagMapJSONAPI, err := buildFieldsMapByTag("jsonapi", relationRef.Type)
			if err != nil {
				return nil, err
			}

			// get model mapper by relation field name
			modelRelFieldMapper := mapper.GetModelMapByWebName(relWebFieldName, modelSrcMaps)
			if modelRelFieldMapper == nil {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			srcRelModel, ok := srcStructMap[modelRelFieldMapper.SrcName]
			if !ok {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			srcModel, ok := srcRelModel.(mapper.ISrcModel)
			if !ok {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			modelSrcRelMaps := srcModel.GetModelMapper()

			//process attr field
			srcVal, srcFieldName, err := processAttrField(filterRow.Field, filterRow.Value, relTagMapJSONAPI, relationRef.Type.Elem(), modelSrcRelMaps)
			if err != nil {
				return nil, err
			}

			newFilterRow = &Filter{
				Relation:        modelRelFieldMapper.SrcName,
				CompareOperator: filterRow.CompareOperator,
				LogicalOperator: filterRow.LogicalOperator,
				Field:           srcFieldName,
				Value:           srcVal,
			}

		} else {

			srcVal, srcFieldName, err := processAttrField(filterRow.Field, filterRow.Value, tagMapJSONAPI, webModelRefType, modelSrcMaps)
			if err != nil {
				return nil, err
			}

			newFilterRow = &Filter{
				CompareOperator: filterRow.CompareOperator,
				LogicalOperator: filterRow.LogicalOperator,
				Field:           srcFieldName,
				Value:           srcVal,
			}
		}

		out[i] = newFilterRow
	}

	return out, nil
}

func SetGormFilters(db *gorm.DB, model interface{}, filtersIn []*Filter) (*gorm.DB, error) {

	if filtersIn == nil {
		return db, nil
	}

	err := db.Statement.Parse(model)
	if err != nil {
		return nil, err
	}

	for _, filter := range filtersIn {
		dbFieldName := db.Statement.NamingStrategy.ColumnName("", filter.Field)
		if filter.Relation != "" {
			rel, ok := db.Statement.Schema.Relationships.Relations[filter.Relation]
			if !ok {
				return nil, errors.New("wrong filter")
			}

			var joinType string
			//if filter.LogicalOperator == "AND" {
			joinType = "JOIN"
			//} else {
			//	joinType = "LEFT JOIN"
			//}

			joinStr := fmt.Sprintf("%v %v on %v.%v = %v.%v and %v.%v %v ?",
				joinType,
				rel.FieldSchema.Table,
				rel.FieldSchema.Table,
				rel.References[0].PrimaryKey.DBName,
				rel.Schema.Table,
				rel.References[0].ForeignKey.DBName,
				rel.FieldSchema.Table,
				dbFieldName,
				filter.CompareOperator,
			)

			db = db.Joins(joinStr, filter.Value)
		} else {
			whereStr := fmt.Sprintf("%v %v ?", dbFieldName, filter.CompareOperator) //TODO security? sql injection!
			db = db.Where(whereStr, filter.Value)
		}
	}
	return db, nil

}

func buildFieldsMapByTag(tagKey string, srcType reflect.Type) (map[string]string, error) {
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

func processAttrField(filterName string, filterValue interface{}, tagMap map[string]string, modelRefType reflect.Type, modelSrcMaps []*mapper.ModelMapper) (srcValue interface{}, scrName string, errOut error) {
	webFieldName, err := getFieldNameByFilterName(filterName, false, tagMap)
	if err != nil {
		return nil, "", err
	}

	modelFieldMapper := mapper.GetModelMapByWebName(webFieldName, modelSrcMaps)
	if modelFieldMapper == nil {
		return nil, "", fmt.Errorf("Cannot find map for field %v", webFieldName)
	}

	// write filterRow.Value to tmpValRef(type interface{})
	webFieldRef, ok := modelRefType.FieldByName(webFieldName)
	if !ok {
		return nil, "", fmt.Errorf("cant find Field %v", filterName)
	}

	// if filter is Slice of values -> make new slice and convert to Src type
	if reflect.TypeOf(filterValue).Kind() == reflect.Slice {
		filterValueSlice, ok := filterValue.([]string)
		if !ok {
			return nil, "", fmt.Errorf("error with Field %v : wrong filter type", filterName)
		}

		srcValSlice := make([]interface{}, len(filterValueSlice))

		for index, filterVal := range filterValueSlice {
			srcVal, err := convertFilterValueToSrc(webFieldRef.Type, filterVal, modelFieldMapper)
			if err != nil {
				return nil, "", fmt.Errorf("error with Field %v : %v", filterName, err)
			}
			srcValSlice[index] = srcVal
		}
		return srcValSlice, modelFieldMapper.SrcName, nil
	} else {
		srcVal, err := convertFilterValueToSrc(webFieldRef.Type, filterValue.(string), modelFieldMapper)
		if err != nil {
			return nil, "", fmt.Errorf("error with Field %v : %v", filterName, err)
		}

		return srcVal, modelFieldMapper.SrcName, nil
	}
	return nil, "", nil
}

func convertFilterValueToSrc(webFieldType reflect.Type, filterValStr string, modelFieldMapper *mapper.ModelMapper) (interface{}, error) {
	var tmpValIntf interface{}
	tmpValRef := reflect.ValueOf(&tmpValIntf)

	err := setFilterVal(tmpValRef, webFieldType, filterValStr)
	if err != nil {
		return nil, err
	}

	// Web field name and value -> Src field name and value
	var srcVal interface{}
	if modelFieldMapper.ConverterToSrc != nil {
		srcVal, err = modelFieldMapper.ConverterToSrc(tmpValRef.Elem().Interface())
		if err != nil {
			return nil, err
		}
	} else {
		srcVal = tmpValRef.Elem().Interface()
	}
	return srcVal, nil
}

func getFieldNameByFilterName(filterName string, isRelation bool, tagMap map[string]string) (string, error) {
	if filterName == JSONAPIIdFieldName {
		return "ID", nil
	}

	var (
		fieldName string
		ok        bool
	)

	if isRelation {
		fieldName, ok = tagMap[JSONAPIRelationTagPrefix+filterName]
		if !ok {
			return "", fmt.Errorf("cant find Field %v", filterName)
		}
		return fieldName, nil
	}

	fieldName, ok = tagMap[JSONAPIAttributeTagPrefix+filterName]
	if !ok {
		return "", fmt.Errorf("cant find Field %v", filterName)
	}

	return fieldName, nil
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

func GetOperator(s string) (string, int) {
	operators := []string{
		"<=", ">=", "!=", "IN", "in", ">", "<", "=",
	}
	for _, o := range operators {
		index := strings.Index(s, o)
		if index > 0 {
			if o == "IN" || o == "in" {
				return FilterOperatorIN, index
			}
			return o, index
		}
	}

	return "", -1
}
