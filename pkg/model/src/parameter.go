package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/domain"
)

type Parameter struct {
	ID                int            `gorm:"primaryKey" omni:"ID;src:ID2src;domain:ID2web"`
	Name              string         `omni:"Name"`
	Description       string         `omni:"Description"`
	Type              string         `omni:"Type"`
	IsValuesSetFinite bool           `omni:"IsValuesSetFinite"`
	Info              datatypes.JSON `omni:"Info;src:JSON2src;domain:JSON2web"`
}

type Parameters []*Parameter

//func (m *Parameter) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Description", WebName: "Description"},
//		&mapper.ModelMap{SrcName: "Type", WebName: "Type"},
//		&mapper.ModelMap{SrcName: "IsValuesSetFinite", WebName: "IsValuesSetFinite"},
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
//	}
//}

func (m *Parameter) ToWeb() (*webmodels.Parameter, error) {
	web := new(webmodels.Parameter)

	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (*Parameter) ScanFromWeb(web *webmodels.Parameter) (*Parameter, error) {
	m := new(Parameter)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Parameters) ToWeb() ([]*webmodels.Parameter, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Parameter, 0, 5)
	for _, u := range m {
		webObj, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webObj)
	}
	return omniM, nil
}

func (m Parameters) ScanFromWeb(web []*webmodels.Parameter) (Parameters, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Parameter)
	res := make(Parameters, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
