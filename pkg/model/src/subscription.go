package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
)

type Subscription struct {
	ID              int            `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	Title           string         `omni:"Title"`
	ContactChannels datatypes.JSON `omni:"ContactChannels;src:JSON2src;web:JSON2web"`
	Options         datatypes.JSON `omni:"Options;src:JSON2src;web:JSON2web"`
	CompanyID       int
	Company         *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
	UserID          string
	User            *User `gorm:"foreignKey:UserID" omni:"User"`
	Rules           Rules `gorm:"many2many:rules_subscriptions;joinForeignKey:SubscriptionID;JoinReferences:rules_group_id" omni:"Rules"`
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

func (m *Subscription) ToWeb() (*webmodels.Subscription, error) {
	web := new(webmodels.Subscription)

	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (*Subscription) ScanFromWeb(web *webmodels.Subscription) (*Subscription, error) {
	m := new(Subscription)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Subscriptions) ToWeb() ([]*webmodels.Subscription, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Subscription, 0, 5)
	for _, u := range m {
		webUser, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}

func (m Subscriptions) ScanFromWeb(web []*webmodels.Subscription) (Subscriptions, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Subscription)
	res := make(Subscriptions, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
