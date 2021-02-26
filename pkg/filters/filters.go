package filters

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	omniErr "omnimanage/pkg/error"
	"omnimanage/pkg/mapper"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var queryFilterRegex = regexp.MustCompile(`^filter\[(.*)\]\[(.*)\]=(.*)$`) //MustCompile(`^filters\[(\w+)\]$`)

type Filter struct {
	Relation        string
	Field           string
	CompareOperator string
	Value           interface{}
	//LogicalOperator string
}

const (
	JSONAPIRelationTagPrefix  = "relation,"
	JSONAPIAttributeTagPrefix = "attr,"
	JSONAPIIdFieldName        = "id"
)

const (
	FilterOperatorIN = "in"
)

var OperatorsMap = map[string]string{
	"eq": "=",
	"ne": "<>",
	"gt": ">",
	"ge": ">=",
	"lt": "<",
	"le": "<=",
	"in": "in",
}

// filter schema: filter[relation.relation_field][operator]=value
func ParseFiltersFromQueryToSrcModel(queryStr string, modelWeb interface{}, modelSrc mapper.ISrcModel) ([]*Filter, error) {
	filtersStrings, err := ParseQueryString(queryStr, modelWeb)
	if err != nil {
		return nil, err
	}

	srcFilters, err := TransformWebToSrc(filtersStrings, modelWeb, modelSrc)
	if err != nil {
		return nil, err
	}

	return srcFilters, nil
}

func ParseQueryString(queryStr string, modelWeb interface{}) ([]*Filter, error) {
	parts := strings.Split(queryStr, "&")

	tagMap, err := buildFieldsMapByTag("jsonapi", reflect.TypeOf(modelWeb))
	if err != nil {
		return nil, err
	}

	filters := make([]*Filter, 0, 1)
	for _, p := range parts {
		regRes := queryFilterRegex.FindStringSubmatch(p)
		if len(regRes) == 0 {
			continue
		}
		fNew := &Filter{}

		fieldParts := strings.Split(regRes[1], ".")
		switch len(fieldParts) {
		case 1:
			fNew.Field = regRes[1]
		case 2:
			fNew.Relation = fieldParts[0]
			fNew.Field = fieldParts[1]
		default:
			return nil, fmt.Errorf("Wrong filter query: %v", p)
		}

		// only Relation
		_, isRelation := tagMap[JSONAPIRelationTagPrefix+fNew.Field]
		if isRelation {
			fNew.Relation = fNew.Field
			fNew.Field = JSONAPIIdFieldName
		}

		var ok bool
		fNew.CompareOperator, ok = OperatorsMap[regRes[2]]
		if !ok {
			return nil, fmt.Errorf("Wrong filter query: uknown operator %v", regRes[2])
		}

		if fNew.CompareOperator == FilterOperatorIN {
			valueRange := strings.Split(regRes[3], ",")
			fNew.Value = valueRange
		} else {
			fNew.Value = regRes[3]
		}

		filters = append(filters, fNew)
	}
	if len(filters) == 0 {
		return nil, nil
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

	//srcStruct := structs.New(modelSrc)
	//srcStructMap := srcStruct.Map()

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

			var relRefType reflect.Type
			if relationRef.Type.Kind() == reflect.Slice {
				relRefType = relationRef.Type.Elem()
			} else {
				relRefType = relationRef.Type
			}

			relTagMapJSONAPI, err := buildFieldsMapByTag("jsonapi", relRefType)
			if err != nil {
				return nil, err
			}

			// get model mapper by relation field name
			modelRelFieldMapper := mapper.GetModelMapByWebName(relWebFieldName, modelSrcMaps)
			if modelRelFieldMapper == nil {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			srcRefModel, ok := reflect.TypeOf(modelSrc).Elem().FieldByName(modelRelFieldMapper.SrcName)
			if !ok {
				return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			}

			var srcRefType reflect.Type
			if srcRefModel.Type.Kind() == reflect.Slice {
				srcRefType = srcRefModel.Type.Elem()
			} else {
				srcRefType = srcRefModel.Type
			}

			modelSrcRelMaps, err := mapper.GetMapperDynamic(srcRefType)
			if err != nil {
				return nil, fmt.Errorf("Field %v: %v", modelRelFieldMapper.SrcName, err)
			}
			//srcRelModel, ok := srcStructMap[modelRelFieldMapper.SrcName]
			//if !ok {
			//	return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			//}
			//
			//var srcModel mapper.ISrcModel
			//if reflect.TypeOf(srcRelModel).Kind() == reflect.Slice {
			//	srcModels, ok := srcRelModel.([]*mapper.ISrcModel)
			//	if !ok {
			//		return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			//	}
			//	srcModels = append(srcModels, nil)
			//	srcModel = *srcModels[0]
			//} else {
			//	srcModel, ok = srcRelModel.(mapper.ISrcModel)
			//	if !ok {
			//		return nil, fmt.Errorf("Cannot find map for field %v", relWebFieldName)
			//	}
			//}
			//
			//modelSrcRelMaps := srcModel.GetModelMapper()

			//process attr field
			srcVal, srcFieldName, err := processAttrField(filterRow.Field, filterRow.Value, relTagMapJSONAPI, relRefType.Elem(), modelSrcRelMaps)
			if err != nil {
				return nil, err
			}

			newFilterRow = &Filter{
				Relation:        modelRelFieldMapper.SrcName,
				CompareOperator: filterRow.CompareOperator,
				//LogicalOperator: filterRow.LogicalOperator,
				Field: srcFieldName,
				Value: srcVal,
			}

		} else {

			srcVal, srcFieldName, err := processAttrField(filterRow.Field, filterRow.Value, tagMapJSONAPI, webModelRefType, modelSrcMaps)
			if err != nil {
				return nil, err
			}

			newFilterRow = &Filter{
				CompareOperator: filterRow.CompareOperator,
				//LogicalOperator: filterRow.LogicalOperator,
				Field: srcFieldName,
				Value: srcVal,
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

			var joinStr string
			if rel.JoinTable == nil { // 1 to 1 relation
				joinStr = fmt.Sprintf(`
					JOIN %v 
						on %v.%v = %v.%v 
					AND %v.%v %v ?`,

					rel.FieldSchema.Table,
					rel.FieldSchema.Table,
					rel.References[0].PrimaryKey.DBName,
					rel.Schema.Table,
					rel.References[0].ForeignKey.DBName,
					rel.FieldSchema.Table,
					dbFieldName,
					filter.CompareOperator,
				)
			} else { // many to many relation
				// join with main table
				joinPart1, err := GetJoinConditionFromJoinTable(rel.JoinTable, rel.Schema.Table)
				if err != nil {
					return nil, err
				}
				joinPart1 = fmt.Sprintf("JOIN %v %v", rel.JoinTable.Table, joinPart1)

				// join with relation table
				joinPart2, err := GetJoinConditionFromJoinTable(rel.JoinTable, rel.FieldSchema.Table)
				if err != nil {
					return nil, err
				}
				joinPart2 = fmt.Sprintf("JOIN %v %v", rel.FieldSchema.Table, joinPart2)

				// total join
				joinStr = fmt.Sprintf(`
					%v
					%v
					AND %v.%v %v ?`,
					joinPart1,
					joinPart2,
					rel.FieldSchema.Table,
					dbFieldName,
					filter.CompareOperator,
				)

			}

			db = db.Joins(joinStr, filter.Value)
		} else {
			whereStr := fmt.Sprintf("%v %v ?", dbFieldName, filter.CompareOperator) //TODO security? sql injection!
			db = db.Where(whereStr, filter.Value)
		}
	}
	return db, nil

}

func GetJoinConditionFromJoinTable(joinTable *schema.Schema, tableName string) (string, error) {

	relTab := getRelationByTableName(tableName, joinTable.Relationships.Relations)
	if relTab == nil {
		return "", fmt.Errorf("Relation to %v not found", tableName)
	}
	joinTabField := relTab.References[0].ForeignKey.DBName

	joinStr := fmt.Sprintf("ON %v.%v = %v.%v",
		joinTable.Table,
		joinTabField,
		tableName,
		relTab.References[0].PrimaryKey.DBName,
	)

	return joinStr, nil
}

func getRelationByTableName(tableName string, relations map[string]*schema.Relationship) *schema.Relationship {
	for _, rel := range relations {
		if rel.FieldSchema.Table == tableName {
			return rel
		}
	}

	return nil
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

func processAttrField(filterName string, filterValue interface{}, tagMap map[string]string, modelRefType reflect.Type, modelSrcMaps []*mapper.ModelMap) (srcValue interface{}, scrName string, errOut error) {
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

func convertFilterValueToSrc(webFieldType reflect.Type, filterValStr string, modelFieldMapper *mapper.ModelMap) (interface{}, error) {
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

//func GetOperator(s string) (string, int) {
//	operators := []string{
//		"<=", ">=", "!=", "IN", "in", ">", "<", "=",
//	}
//	for _, o := range operators {
//		index := strings.Index(s, o)
//		if index > 0 {
//			if o == "IN" || o == "in" {
//				return FilterOperatorIN, index
//			}
//			return o, index
//		}
//	}
//
//	return "", -1
//}

func GetSrcFiltersFromRelationID(model interface{}) (res []*Filter, errOut error) {
	if model == nil {
		return nil, nil
	}
	defer func() {
		if r := recover(); r != nil {
			errOut = fmt.Errorf("panic: %v", r)
		}
	}()

	resVals := make([]int, 0, 10)

	kind := reflect.TypeOf(model).Kind()
	switch kind {
	case reflect.Slice:
		modelSlRef := reflect.ValueOf(model)

		for i := 0; i < modelSlRef.Len(); i++ {
			row := modelSlRef.Index(i)
			idIntf := row.Elem().FieldByName("ID").Interface()
			id, ok := idIntf.(int)
			if !ok {
				return nil, fmt.Errorf("%w: wrong type of ID", omniErr.ErrInternal)
			}
			resVals = append(resVals, id)
		}
		if len(resVals) == 0 {
			return nil, fmt.Errorf("%w: empty filter", omniErr.ErrInternal)
		}
		return []*Filter{
			&Filter{Field: "ID", CompareOperator: "in", Value: resVals},
		}, nil
	default:
		return nil, fmt.Errorf("w: wrong filter model type", omniErr.ErrInternal)
	}
	return nil, nil
}
