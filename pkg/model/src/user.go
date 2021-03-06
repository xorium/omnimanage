package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/model/domain"
)

type User struct {
	ID            int            `gorm:"primaryKey" omni:"ID;src:ID2src;domain:ID2web"`
	UserName      string         `omni:"Name"`
	Password      string         `omni:"Password"`
	FirstName     string         `omni:"FirstName"`
	LastName      string         `omni:"LastName"`
	PhoneNumber   string         `omni:"PhoneNumber"`
	Email         string         `omni:"Email"`
	Image         string         `omni:"Image"`
	Settings      datatypes.JSON `omni:"Settings;src:JSON2src;domain:JSON2web"`
	CompanyID     int
	Company       *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
	LocationID    int
	Location      *Location     `gorm:"foreignKey:LocationID" omni:"Location"`
	Roles         Roles         `gorm:"many2many:user_role;joinForeignKey:UsersID;JoinReferences:roles_id" omni:"Roles"`
	Subscriptions Subscriptions `gorm:"foreignKey:UserID" omni:"Subscriptions"`
}

type Users []*User

//func (m *User) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "UserName", WebName: "Name"},
//		&mapper.ModelMap{SrcName: "Password", WebName: "Password"},
//		&mapper.ModelMap{SrcName: "FirstName", WebName: "FirstName"},
//		&mapper.ModelMap{SrcName: "LastName", WebName: "LastName"},
//		&mapper.ModelMap{SrcName: "PhoneNumber", WebName: "PhoneNumber"},
//		&mapper.ModelMap{SrcName: "Email", WebName: "Email"},
//		&mapper.ModelMap{SrcName: "Image", WebName: "Image"},
//		&mapper.ModelMap{SrcName: "Settings", WebName: "Settings",
//			ConverterToSrc: func(domain interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(domain)
//				if err != nil {
//					return nil, fmt.Errorf("Settings: %v. %v", domain, err)
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
//		&mapper.ModelMap{SrcName: "Location", WebName: "Location"},
//		&mapper.ModelMap{SrcName: "Roles", WebName: "Roles"},
//	}
//}

func (m *User) ToWeb() (*domain.User, error) {
	web := new(domain.User)

	mr := mapper.Get()
	err := mr.ConvertSrcToWeb(m, &web)

	if err != nil {
		return nil, err
	}
	return web, nil
}

func (*User) ScanFromWeb(web *domain.User) (*User, error) {
	m := new(User)
	mr := mapper.Get()
	err := mr.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Users) ToWeb() ([]*domain.User, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*domain.User, 0, 5)
	for _, u := range m {
		webUser, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webUser)
	}
	return omniM, nil
}

func (m Users) ScanFromWeb(web []*domain.User) (Users, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(User)
	res := make(Users, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
