package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
)

type Rule struct {
	ID                int            `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	Title             string         `omni:"Title"`
	Slug              string         `omni:"Slug"`
	Expression        datatypes.JSON `omni:"Expression;src:JSON2src;web:JSON2web"`
	Duration          int
	EventLevel        string `omni:"EventLevel"`
	EventSessionState string `omni:"EventSessionState"`
	RuleGroup         string `omni:"RuleGroup"`
	CompanyID         int
	Company           *Company   `gorm:"foreignKey:CompanyID" omni:"Company"`
	Devices           Devices    `gorm:"many2many:rules_devices;joinForeignKey:RuleID;JoinReferences:device_id" omni:"Devices"`
	Params            Parameters `gorm:"many2many:rules_parameters;joinForeignKey:RuleID;JoinReferences:parameter_id" omni:"Parameters"`
}
type Rules []*Rule

//func (m *Rule) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Slug", WebName: "Slug"},
//		&mapper.ModelMap{SrcName: "Expression", WebName: "Expression",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(web)
//				if err != nil {
//					return nil, fmt.Errorf("Expression: %v. %v", web, err)
//				}
//				return j, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				w, err := converters.JSONSrcToWeb(src)
//				if err != nil {
//					return nil, fmt.Errorf("Expression: %v. %v", src, err)
//				}
//				return w, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Duration", WebName: "Duration"},
//		&mapper.ModelMap{SrcName: "EventLevel", WebName: "EventLevel"},
//		&mapper.ModelMap{SrcName: "EventSessionState", WebName: "EventSessionState"},
//		&mapper.ModelMap{SrcName: "RuleGroup", WebName: "RuleGroup"},
//
//		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
//		&mapper.ModelMap{SrcName: "Devices", WebName: "Devices"},
//		&mapper.ModelMap{SrcName: "Params", WebName: "Params"},
//	}
//}

func (m *Rule) ToWeb(mapper *mapper.ModelMapper) (*webmodels.Rule, error) {
	web := new(webmodels.Rule)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (*Rule) ScanFromWeb(web *webmodels.Rule, mapper *mapper.ModelMapper) (*Rule, error) {
	m := new(Rule)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m Rules) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.Rule, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Rule, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m Rules) ScanFromWeb(web []*webmodels.Rule, mapper *mapper.ModelMapper) (Rules, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Rule)
	res := make(Rules, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u, mapper)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
