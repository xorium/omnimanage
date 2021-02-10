package filters

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
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

func SetGormFilters(db *gorm.DB, model interface{}, f []*Filter) (*gorm.DB, error) {

	if f == nil {
		return db, nil
	}

	err := db.Statement.Parse(model)
	if err != nil {
		return nil, err
	}

	for _, filter := range f {
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
				filter.Field,
				filter.Operator,
			)

			db = db.Joins(joinStr, filter.Value)
		} else {
			whereStr := fmt.Sprintf("%v%v?", filter.Field, filter.Operator) //TODO security? sql injection!
			db = db.Where(whereStr, filter.Value)
		}
	}
	return db, nil

	//var tx *gorm.DB
	//tx = db
	//for _, val := range f {
	//	if val.Relation != "" {
	//		as := tx.Relation(val.Relation)
	//
	//	}
	//
	//}
	//
	//return tx
	//if filters != nil {
	//	for col, val := range filters {
	//		if tx == nil {
	//			tx = db.Where(fmt.Sprintf("%v = ?", col), val)
	//			continue
	//		}
	//		tx = tx.Where(fmt.Sprintf("%v = ?", col), val)
	//	}
	//} else {
	//	tx = db
	//}
	//return tx
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

func Transform(f []*Filter, modelWeb interface{}, modelDest interface{}) ([]*Filter, error) {
	if f == nil {
		return nil, nil
	}

	refWeb := reflect.TypeOf(modelWeb)
	//refDest := reflect.TypeOf(modelDest)

	res := make([]*Filter, len(f), cap(f))

	tagMap, err := buildFieldByTagMap("jsonapi", reflect.TypeOf(modelWeb))
	if err != nil {
		return nil, err
	}

	for i, filt := range f {
		fNew := &Filter{
			Operator: filt.Operator,
		}

		filtVal := reflect.ValueOf(&fNew.Value)

		if filt.Relation != "" {

			tagName := "relation," + filt.Relation
			relFieldName, ok := tagMap[tagName]
			if !ok {
				return nil, fmt.Errorf("cant find field %v", filt.Field)
			}

			rel, ok := refWeb.FieldByName(relFieldName)
			if !ok {
				return nil, fmt.Errorf("cant find Relation name %v", filt.Relation)
			}

			relTagMap, err := buildFieldByTagMap("jsonapi", rel.Type)
			if err != nil {
				return nil, err
			}

			tagName = "attr," + filt.Field
			fieldName, ok := relTagMap[tagName]
			if !ok {
				return nil, fmt.Errorf("cant find field %v", filt.Field)
			}

			field, ok := rel.Type.Elem().FieldByName(fieldName)
			if !ok {
				return nil, fmt.Errorf("cant find field %v", filt.Field)
			}

			err = setFilterVal(filtVal, field.Type, filt.Value.(string))
			if err != nil {
				return nil, fmt.Errorf("error with field %v : %v", filt.Field, err)
			}

		} else {
			tagName := "attr," + filt.Field
			fieldName, ok := tagMap[tagName]
			if !ok {
				return nil, fmt.Errorf("cant find field %v", filt.Field)
			}
			field, ok := refWeb.FieldByName(fieldName)
			if !ok {
				return nil, fmt.Errorf("cant find field %v", filt.Field)
			}

			err := setFilterVal(filtVal, field.Type, filt.Value.(string))
			if err != nil {
				return nil, fmt.Errorf("error with field %v : %v", filt.Field, err)
			}
		}

		res[i] = fNew
	}

	return res, nil
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

func GetFiltersFromQueryString(queryStr string) ([]*Filter, error) {

	if !strings.HasPrefix(queryStr, "filter=") {
		return nil, nil
	}
	query := strings.TrimPrefix(queryStr, "filter=")

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
			fNew.Value = p2[operIndex+len(operator):]
			filters = append(filters, fNew)

		}
	}

	return filters, nil
	//filter = strings.ReplaceAll(filter, "&&", " and ")
	//filter = strings.ReplaceAll(filter, "||", " or ")
	//return filter, nil

	//params, err := ctx.FormParams()
	//if err != nil {
	//	return "", err
	//}
	//
	////filt = strings.ReplaceAll(" and ", "&&", filt)

	//filters := make(map[string]string)
	//for key, val := range params {
	//match := queryFilterRegex.FindStringSubmatch(key)
	//if len(match) > 1 {
	//	filters[match[1]] = val[0]
	//}

	//}
	//if len(filters) == 0 {
	//	return nil, nil
	//}
	//return filters, nil
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
