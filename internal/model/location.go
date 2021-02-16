package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
)

type Location struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Timezone  string
	Info      datatypes.JSON
	CompanyID int
	Company   *Company `gorm:"foreignKey:CompanyID"`
	ParentID  int
	Children  Locations `gorm:"foreignkey:ParentID"`
	Users     Users     `gorm:"foreignkey:LocationID"`
}

type Locations []*Location

func (m *Location) GetModelMapper() []*mapper.ModelMapper {
	return []*mapper.ModelMapper{
		&mapper.ModelMapper{SrcName: "ID", WebName: "ID",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				id, err := converters.IDWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("ID: %v. %v", web, err)
				}
				return id, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				id, err := converters.IDSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("ID: %v. %v", src, err)
				}
				return id, nil
			},
		},
		&mapper.ModelMapper{SrcName: "Name", WebName: "Name"},
		&mapper.ModelMapper{SrcName: "Timezone", WebName: "Timezone"},
		&mapper.ModelMapper{SrcName: "Info", WebName: "Info",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				j, err := converters.JSONWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("Settings: %v. %v", web, err)
				}
				return j, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				w, err := converters.JSONSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("Settings: %v. %v", src, err)
				}
				return w, nil
			},
		},
		&mapper.ModelMapper{SrcName: "Company", WebName: "Company"},
		&mapper.ModelMapper{SrcName: "Children", WebName: "Children"},
		&mapper.ModelMapper{SrcName: "Users", WebName: "Users"},
	}
}

func (m *Location) ToWeb() (*omnimodels.Location, error) {
	web := new(omnimodels.Location)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (m Locations) ToWeb() ([]*omnimodels.Location, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.Location, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}
