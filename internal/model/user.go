package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
)

type User struct {
	ID            int    `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	UserName      string `omni:"UserName"`
	Password      string `omni:"Password"`
	FirstName     string
	LastName      string
	PhoneNumber   string
	Email         string
	Image         string
	Settings      datatypes.JSON `omni:"Settings;src:JSON2src;web:JSON2web"`
	CompanyID     int
	Company       *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
	LocationID    int
	Location      *Location       `gorm:"foreignKey:LocationID"`
	Roles         Roles           `gorm:"many2many:user_role;joinForeignKey:UsersID;JoinReferences:roles_id"`
	Subscriptions []*Subscription `gorm:"foreignKey:UserID"`
}

type Users []*User

func (m *User) GetModelMapper() []*mapper.ModelMap {
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
		&mapper.ModelMap{SrcName: "UserName", WebName: "Name"},
		&mapper.ModelMap{SrcName: "Password", WebName: "Password"},
		&mapper.ModelMap{SrcName: "FirstName", WebName: "FirstName"},
		&mapper.ModelMap{SrcName: "LastName", WebName: "LastName"},
		&mapper.ModelMap{SrcName: "PhoneNumber", WebName: "PhoneNumber"},
		&mapper.ModelMap{SrcName: "Email", WebName: "Email"},
		&mapper.ModelMap{SrcName: "Image", WebName: "Image"},
		&mapper.ModelMap{SrcName: "Settings", WebName: "Settings",
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
		&mapper.ModelMap{SrcName: "Location", WebName: "Location"},
		&mapper.ModelMap{SrcName: "Roles", WebName: "Roles"},
	}
}

func (m *User) ToWeb() (*omnimodels.User, error) {
	web := new(omnimodels.User)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (*User) ScanFromWeb(web *omnimodels.User) (*User, error) {
	m := new(User)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Users) ToWeb() ([]*omnimodels.User, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.User, 0, 5)
	for _, u := range m {
		webUser, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}
