package filters

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strings"
)

//var queryFilterRegex = regexp.MustCompile(`^filters\[(\w+)\]$`)

type Filter struct {
	Association string
	Field       string
	Operator    string
	Value       string
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
		if filter.Association != "" {
			rel, ok := db.Statement.Schema.Relationships.Relations[filter.Association]
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
			whereStr := fmt.Sprintf("%v%v?", filter.Field, filter.Operator)
			db = db.Where(whereStr, filter.Value)
		}
	}
	return db, nil

	//var tx *gorm.DB
	//tx = db
	//for _, val := range f {
	//	if val.Association != "" {
	//		as := tx.Association(val.Association)
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

func GetFiltersFromContext(ctx echo.Context) ([]*Filter, error) {

	if !strings.HasPrefix(ctx.Request().URL.RawQuery, "filter=") {
		return nil, nil
	}
	query := strings.TrimPrefix(ctx.Request().URL.RawQuery, "filter=")

	filters := make([]*Filter, 0, 1)

	parts1 := strings.Split(query, "||")
	for _, p1 := range parts1 {
		parts2 := strings.Split(p1, "&&")
		for _, p2 := range parts2 {
			fNew := &Filter{}

			i := strings.Index(p2, ".")
			if i > 0 {
				fNew.Association = p2[:i]
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
