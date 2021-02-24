package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"omnimanage/pkg/mapper"
	"strconv"
)

type Company struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func (Company) TableName() string {
	return "companies"
}

func (m *Company) GetModelMapper() []*mapper.ModelMapper {
	return []*mapper.ModelMapper{
		&mapper.ModelMapper{SrcName: "ID", WebName: "ID",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				w, ok := web.(string)
				if !ok {
					return nil, fmt.Errorf("ID: Wrong type. Value %v, type %T", web, web)
				}
				id, err := strconv.Atoi(w)
				if err != nil {
					return nil, fmt.Errorf("Wrong Company ID: %v", w)
				}
				return id, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				s, ok := src.(int)
				if !ok {
					return nil, fmt.Errorf("ID: Wrong type. Value %v, type %T", src, src)
				}
				id := strconv.Itoa(s)
				return id, nil
			},
		},
		&mapper.ModelMapper{SrcName: "Name", WebName: "Name"},
	}
}

func (m *Company) ToWeb() (*omnimodels.Company, error) {
	web := new(omnimodels.Company)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Company) ScanFromWeb(web *omnimodels.Company) (*Company, error) {
	m := new(Company)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return m, err
	}

	return m, nil
}
