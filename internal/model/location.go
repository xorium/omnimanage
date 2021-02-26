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

func (m *Location) GetModelMapper() []*mapper.ModelMap {
	return []*mapper.ModelMap{
		&mapper.ModelMap{SrcName: "ID", WebName: "ID",
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
		&mapper.ModelMap{SrcName: "Name", WebName: "Name"},
		&mapper.ModelMap{SrcName: "Timezone", WebName: "Timezone"},
		&mapper.ModelMap{SrcName: "Info", WebName: "Info",
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
		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
		&mapper.ModelMap{SrcName: "Children", WebName: "Children"},
		&mapper.ModelMap{SrcName: "Users", WebName: "Users"},
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

func (*Location) ScanFromWeb(web *omnimodels.Location) (*Location, error) {
	m := new(Location)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m Locations) ToWeb() ([]*omnimodels.Location, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.Location, 0, len(m))
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}
