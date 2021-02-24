package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
)

type Role struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Slug       string
	Assigned   bool
	Persistent bool
	Users      Users `gorm:"many2many:user_role;joinForeignKey:RolesID;JoinReferences:users_id"`
	Info       datatypes.JSON
	CompanyID  int
	Company    *Company `gorm:"foreignKey:CompanyID"`
}

type Roles []*Role

func (m *Role) GetModelMapper() []*mapper.ModelMapper {
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
		&mapper.ModelMapper{SrcName: "Slug", WebName: "Slug"},
		&mapper.ModelMapper{SrcName: "Assigned", WebName: "Assigned"},
		&mapper.ModelMapper{SrcName: "Persistent", WebName: "Persistent"},
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
		&mapper.ModelMapper{SrcName: "Company", WebName: "Company"},
	}
}

func (m *Role) ToWeb() (*omnimodels.Role, error) {
	web := new(omnimodels.Role)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (m Roles) ToWeb() ([]*omnimodels.Role, error) {
	if m == nil {
		return nil, nil
	}

	omniM := make([]*omnimodels.Role, 0, 5)
	for _, u := range m {
		webUser, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}
