package src

import (
	"fmt"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	omnimodels "omnimanage/pkg/model/web"
	"strconv"
)

type Subscription struct {
	ID              int `gorm:"primaryKey`
	Title           string
	ContactChannels datatypes.JSON
	Options         datatypes.JSON
	CompanyID       int
	Company         *Company `gorm:"foreignKey:CompanyID"`
	UserID          string
	User            *User `gorm:"foreignKey:UserID"`
	Rules           Rules `gorm:"many2many:rules_subscriptions;joinForeignKey:SubscriptionID;JoinReferences:rules_group_id"`
}

type Subscriptions []*Subscription

//func (m *Subscription) GetModelMapper() []*mapper.ModelMap {
//	return []*mapper.ModelMap{
//		&mapper.ModelMap{SrcName: "ID", WebName: "ID",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				id, err := converters.IDWebToSrc(web)
//				if err != nil {
//					return nil, fmt.Errorf("ID: %v. %v", web, err)
//				}
//				return id, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				id, err := converters.IDSrcToWeb(src)
//				if err != nil {
//					return nil, fmt.Errorf("ID: %v. %v", src, err)
//				}
//				return id, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Title", WebName: "Title"},
//		&mapper.ModelMap{SrcName: "ContactChannels", WebName: "ContactChannels",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(web)
//				if err != nil {
//					return nil, fmt.Errorf("Settings: %v. %v", web, err)
//				}
//				return j, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				w, err := converters.JSONSrcToWeb(src)
//				if err != nil {
//					return nil, fmt.Errorf("Settings: %v. %v", src, err)
//				}
//				return w, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Options", WebName: "Options",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(web)
//				if err != nil {
//					return nil, fmt.Errorf("Settings: %v. %v", web, err)
//				}
//				return j, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				w, err := converters.JSONSrcToWeb(src)
//				if err != nil {
//					return nil, fmt.Errorf("Settings: %v. %v", src, err)
//				}
//				return w, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
//		&mapper.ModelMap{SrcName: "User", WebName: "User"},
//		&mapper.ModelMap{SrcName: "Rules", WebName: "Rules"},
//	}
//}

func (m *Subscription) ToWeb(mapper *mapper.ModelMapper) (*omnimodels.Subscription, error) {
	web := new(omnimodels.Subscription)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (m *Subscription) ScanFromWeb(us *omnimodels.Subscription, mapper *mapper.ModelMapper) error {
	var err error
	m.ID, err = strconv.Atoi(us.ID)
	if err != nil {
		return fmt.Errorf("Wrong User ID: %v", us.ID)
	}
	//
	//m.UserName = us.Name
	//m.Password = us.Password
	//m.FirstName = us.FirstName
	//m.LastName = us.LastName
	//m.PhoneNumber = us.PhoneNumber
	//m.Email = us.Email
	//m.Image = us.Image
	////....

	return nil
}

func (m Subscriptions) ToWeb(mapper *mapper.ModelMapper) ([]*omnimodels.Subscription, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.Subscription, 0, 5)
	for _, u := range m {
		webUser, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}
