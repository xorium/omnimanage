package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
	"strconv"
)

type Rule struct {
	ID                int `gorm:"primaryKey"`
	Title             string
	Slug              string
	Expression        datatypes.JSON
	Duration          int
	EventLevel        string
	EventSessionState string
	RuleGroup         string
	CompanyID         int
	Company           *Company   `gorm:"foreignKey:CompanyID"`
	Devices           Devices    `gorm:"many2many:rules_devices;joinForeignKey:RuleID;JoinReferences:device_id"`
	Params            Parameters `gorm:"many2many:rules_parameters;joinForeignKey:RuleID;JoinReferences:parameter_id"`
}
type Rules []*Rule

func (m *Rule) GetModelMapper() []*mapper.ModelMap {
	return []*mapper.ModelMap{
		&mapper.ModelMap{SrcName: "ID", WebName: "ID",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				id, err := converters.IDWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("ID: %v. %v", web, err)
				}
				return id, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				id, err := converters.IDSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("ID: %v. %v", src, err)
				}
				return id, nil
			},
		},
		&mapper.ModelMap{SrcName: "Title", WebName: "Title"},
		&mapper.ModelMap{SrcName: "Slug", WebName: "Slug"},
		&mapper.ModelMap{SrcName: "Expression", WebName: "Expression",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				j, err := converters.JSONWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("Expression: %v. %v", web, err)
				}
				return j, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				w, err := converters.JSONSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("Expression: %v. %v", src, err)
				}
				return w, nil
			},
		},
		&mapper.ModelMap{SrcName: "Duration", WebName: "Duration"},
		&mapper.ModelMap{SrcName: "EventLevel", WebName: "EventLevel"},
		&mapper.ModelMap{SrcName: "EventSessionState", WebName: "EventSessionState"},
		&mapper.ModelMap{SrcName: "RuleGroup", WebName: "RuleGroup"},

		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
		&mapper.ModelMap{SrcName: "Devices", WebName: "Devices"},
		&mapper.ModelMap{SrcName: "Params", WebName: "Params"},
	}
}

func (m *Rule) ToWeb() (*omnimodels.Rule, error) {
	web := new(omnimodels.Rule)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}
	return web, nil
}

func (m *Rule) ScanFromWeb(us *omnimodels.Rule) error {
	var err error
	m.ID, err = strconv.Atoi(us.ID)
	if err != nil {
		return fmt.Errorf("Wrong User ID: %v", us.ID)
	}

	return nil
}

func (m Rules) ToWeb() ([]*omnimodels.Rule, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.Rule, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}
