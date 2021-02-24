package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
)

type Parameter struct {
	ID                int `gorm:"primaryKey"`
	Name              string
	Description       string
	Type              string
	IsValuesSetFinite bool
	Info              datatypes.JSON
}

type Parameters []*Parameter

func (m *Parameter) GetModelMapper() []*mapper.ModelMapper {
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
		&mapper.ModelMapper{SrcName: "Description", WebName: "Description"},
		&mapper.ModelMapper{SrcName: "Type", WebName: "Type"},
		&mapper.ModelMapper{SrcName: "IsValuesSetFinite", WebName: "IsValuesSetFinite"},
		&mapper.ModelMapper{SrcName: "Info", WebName: "Info",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				j, err := converters.JSONWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("Info: %v. %v", web, err)
				}
				return j, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				w, err := converters.JSONSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("Info: %v. %v", src, err)
				}
				return w, nil
			},
		},
	}
}

func (m *Parameter) ToWeb() (*omnimodels.Parameter, error) {
	web := new(omnimodels.Parameter)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (m Parameters) ToWeb() ([]*omnimodels.Parameter, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.Parameter, 0, 5)
	for _, u := range m {
		webObj, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webObj)
	}
	return omniM, nil
}
