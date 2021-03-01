package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
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

//func (m *Parameter) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Description", WebName: "Description"},
//		&mapper.ModelMap{SrcName: "Type", WebName: "Type"},
//		&mapper.ModelMap{SrcName: "IsValuesSetFinite", WebName: "IsValuesSetFinite"},
//		&mapper.ModelMap{SrcName: "Info", WebName: "Info",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(web)
//				if err != nil {
//					return nil, fmt.Errorf("Info: %v. %v", web, err)
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
//	}
//}

func (m *Parameter) ToWeb(mapper *mapper.ModelMapper) (*webmodels.Parameter, error) {
	web := new(webmodels.Parameter)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (m Parameters) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.Parameter, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Parameter, 0, 5)
	for _, u := range m {
		webObj, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webObj)
	}
	return omniM, nil
}
