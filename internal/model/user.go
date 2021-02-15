package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
	"strconv"
)

type User struct {
	ID            int `gorm:"primaryKey"`
	UserName      string
	Password      string
	FirstName     string
	LastName      string
	PhoneNumber   string
	Email         string
	Image         string
	Settings      datatypes.JSON
	CompanyID     int
	Company       *Company `gorm:"foreignKey:CompanyID"`
	LocationID    int
	Location      *Location       `gorm:"foreignKey:LocationID"`
	Roles         []*Role         `gorm:"many2many:user_role;joinForeignKey:UsersID;JoinReferences:roles_id"`
	Subscriptions []*Subscription `gorm:"foreignKey:UserID"`
}

func (m *User) GetModelMapper() []*mapper.ModelMapper {
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
		&mapper.ModelMapper{SrcName: "UserName", WebName: "Name"},
		&mapper.ModelMapper{SrcName: "Password", WebName: "Password"},
		&mapper.ModelMapper{SrcName: "FirstName", WebName: "FirstName"},
		&mapper.ModelMapper{SrcName: "LastName", WebName: "LastName"},
		&mapper.ModelMapper{SrcName: "PhoneNumber", WebName: "PhoneNumber"},
		&mapper.ModelMapper{SrcName: "Email", WebName: "Email"},
		&mapper.ModelMapper{SrcName: "Image", WebName: "Image"},
		&mapper.ModelMapper{SrcName: "Settings", WebName: "Settings",
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
		//&mapper.ModelMapper{SrcName: "Location", WebName: "Location"},
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

func (m *User) ScanFromWeb(us *omnimodels.User) error {
	var err error
	m.ID, err = strconv.Atoi(us.ID)
	if err != nil {
		return fmt.Errorf("Wrong User ID: %v", us.ID)
	}

	m.UserName = us.Name
	m.Password = us.Password
	m.FirstName = us.FirstName
	m.LastName = us.LastName
	m.PhoneNumber = us.PhoneNumber
	m.Email = us.Email
	m.Image = us.Image
	//....

	return nil
}

func UsersToWeb(mSl []*User) ([]*omnimodels.User, error) {
	if mSl == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.User, 0, 5)
	for _, u := range mSl {
		webUser, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}
