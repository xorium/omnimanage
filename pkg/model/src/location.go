package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
)

type Location struct {
	ID        int            `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	Name      string         `omni:"Name"`
	Timezone  string         `omni:"Timezone"`
	Info      datatypes.JSON `omni:"Info;src:JSON2src;web:JSON2web"`
	CompanyID int
	Company   *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
	ParentID  int
	Children  Locations `gorm:"foreignkey:ParentID" omni:"Children"`
	Users     Users     `gorm:"foreignkey:LocationID" omni:"Users"`
}

type Locations []*Location

//func (m *Location) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Name", WebName: "Name"},
//		&mapper.ModelMap{SrcName: "Timezone", WebName: "Timezone"},
//		&mapper.ModelMap{SrcName: "Info", WebName: "Info",
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
//		&mapper.ModelMap{SrcName: "Children", WebName: "Children"},
//		&mapper.ModelMap{SrcName: "Users", WebName: "Users"},
//	}
//}

func (m *Location) ToWeb() (*webmodels.Location, error) {
	web := new(webmodels.Location)

	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Location) ScanFromWeb(web *webmodels.Location) (*Location, error) {
	m := new(Location)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m Locations) ToWeb() ([]*webmodels.Location, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Location, 0, len(m))
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m Locations) ScanFromWeb(web []*webmodels.Location) (Locations, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Location)
	res := make(Locations, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
