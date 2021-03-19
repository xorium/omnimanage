package src

import (
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/domain"
)

type Company struct {
	ID   int    `gorm:"primaryKey" omni:"ID;src:ID2src;domain:ID2web"`
	Name string `omni:"Name"`
}

type Companys []*Company

func (Company) TableName() string {
	return "companies"
}

func (m *Company) ToWeb() (*webmodels.Company, error) {
	web := new(webmodels.Company)
	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Company) ScanFromWeb(web *webmodels.Company) (*Company, error) {
	m := new(Company)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m Companys) ToWeb() ([]*webmodels.Company, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Company, 0, 5)
	for _, u := range m {
		webM, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webM)
	}
	return omniM, nil
}
