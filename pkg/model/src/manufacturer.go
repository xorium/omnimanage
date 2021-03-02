package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
)

type Manufacturer struct {
	ID   int            `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	Name string         `omni:"Name"`
	Info datatypes.JSON `omni:"Info;src:JSON2src;web:JSON2web"`
}

type Manufacturers []*Manufacturer

//func (m *Manufacturer) GetModelMapper() []*mapper.ModelMap {
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
//	}
//}

func (m *Manufacturer) ToWeb(mapper *mapper.ModelMapper) (*webmodels.Manufacturer, error) {
	web := new(webmodels.Manufacturer)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Manufacturer) ScanFromWeb(web *webmodels.Manufacturer, mapper *mapper.ModelMapper) (*Manufacturer, error) {
	m := new(Manufacturer)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Manufacturers) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.Manufacturer, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Manufacturer, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m Manufacturers) ScanFromWeb(web []*webmodels.Manufacturer, mapper *mapper.ModelMapper) (Manufacturers, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Manufacturer)
	res := make(Manufacturers, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u, mapper)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
