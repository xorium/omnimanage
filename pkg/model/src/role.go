package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/domain"
)

type Role struct {
	ID         int    `gorm:"primaryKey" omni:"ID;src:ID2src;domain:ID2web"`
	Name       string `omni:"Name"`
	Slug       string `omni:"Slug"`
	Assigned   bool   `omni:"Assigned"`
	Persistent bool   `omni:"Persistent"`
	Users      Users  `gorm:"many2many:user_role;joinForeignKey:RolesID;JoinReferences:users_id"`
	Info       datatypes.JSON
	CompanyID  int
	Company    *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
}

type Roles []*Role

//func (m *Role) GetModelMapper() []*mapper.ModelMap {
//	return []*mapper.ModelMap{
//		&mapper.ModelMap{SrcName: "ID", WebName: "ID",
//			ConverterToSrc: func(domain interface{}) (interface{}, error) {
//				id, err := converters.IDWebToSrc(domain)
//				if err != nil {
//					return nil, fmt.Errorf("ID: %v. %v", domain, err)
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
//		&mapper.ModelMap{SrcName: "Name", WebName: "Name"},
//		&mapper.ModelMap{SrcName: "Slug", WebName: "Slug"},
//		&mapper.ModelMap{SrcName: "Assigned", WebName: "Assigned"},
//		&mapper.ModelMap{SrcName: "Persistent", WebName: "Persistent"},
//		&mapper.ModelMap{SrcName: "Info", WebName: "Info",
//			ConverterToSrc: func(domain interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(domain)
//				if err != nil {
//					return nil, fmt.Errorf("Info: %v. %v", domain, err)
//				}
//				return j, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				w, err := converters.JSONSrcToWeb(src)
//				if err != nil {
//					return nil, fmt.Errorf("Info: %v. %v", src, err)
//				}
//				return w, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
//	}
//}

func (m *Role) ToWeb() (*webmodels.Role, error) {
	web := new(webmodels.Role)
	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (Role) ScanFromWeb(web *webmodels.Role) (*Role, error) {
	m := new(Role)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m Roles) ToWeb() ([]*webmodels.Role, error) {
	if m == nil {
		return nil, nil
	}

	omniM := make([]*webmodels.Role, 0, len(m))
	for _, u := range m {
		webUser, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}

func (m Roles) ScanFromWeb(web []*webmodels.Role) (Roles, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Role)
	res := make(Roles, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
