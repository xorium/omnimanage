package src

import (
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
)

type Company struct {
	ID   int    `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	Name string `omni:"Name"`
}

type Companies []*Company

func (Company) TableName() string {
	return "companies"
}

//func (m *Company) GetModelMapper() []*mapper.ModelMap {
//	return []*mapper.ModelMap{
//		&mapper.ModelMap{SrcName: "ID", WebName: "ID",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				w, ok := web.(string)
//				if !ok {
//					return nil, fmt.Errorf("ID: Wrong type. Value %v, type %T", web, web)
//				}
//				id, err := strconv.Atoi(w)
//				if err != nil {
//					return nil, fmt.Errorf("Wrong Company ID: %v", w)
//				}
//				return id, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				s, ok := src.(int)
//				if !ok {
//					return nil, fmt.Errorf("ID: Wrong type. Value %v, type %T", src, src)
//				}
//				id := strconv.Itoa(s)
//				return id, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Name", WebName: "Name"},
//	}
//}

func (m *Company) ToWeb(mapper *mapper.ModelMapper) (*webmodels.Company, error) {
	web := new(webmodels.Company)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Company) ScanFromWeb(web *webmodels.Company, mapper *mapper.ModelMapper) (*Company, error) {
	m := new(Company)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m Companies) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.Company, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Company, 0, 5)
	for _, u := range m {
		webM, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webM)
	}
	return omniM, nil
}
