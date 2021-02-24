package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	"omnimanage/pkg/utils/converters"
)

/////////////////////////////// Device
type Device struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Slug        string
	Title       string
	Description string
	Kind        string
	Info        datatypes.JSON
	Image       string
	CompanyID   int
	Company     *Company `gorm:"foreignKey:CompanyID"`
	ModelID     int
	Model       *DeviceModel `gorm:"foreignkey:ModelID"`
	LocationID  int
	Location    *Location    `gorm:"foreignKey:LocationID"`
	Groups      DeviceGroups `gorm:"many2many:device_group;joinForeignKey:DevicesID;JoinReferences:groups_id"`
	ParentID    int
	Parent      *Device `gorm:"foreignkey:ParentID"`
	Rules       Rules   `gorm:"many2many:rules_devices;joinForeignKey:DeviceID;JoinReferences:rule_id"`
}

type Devices []*Device

func (m *Device) GetModelMapper() []*mapper.ModelMapper {
	return []*mapper.ModelMapper{
		&mapper.ModelMapper{SrcName: "ID", WebName: "ID",
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
		&mapper.ModelMapper{SrcName: "Name", WebName: "Name"},
		&mapper.ModelMapper{SrcName: "Slug", WebName: "Slug"},
		&mapper.ModelMapper{SrcName: "Title", WebName: "Title"},
		&mapper.ModelMapper{SrcName: "Description", WebName: "Description"},
		&mapper.ModelMapper{SrcName: "Info", WebName: "Info",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				j, err := converters.JSONWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("Settings: %v. %v", web, err)
				}
				return j, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				w, err := converters.JSONSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("Settings: %v. %v", src, err)
				}
				return w, nil
			},
		},
		&mapper.ModelMapper{SrcName: "Image", WebName: "Image"},
		&mapper.ModelMapper{SrcName: "Company", WebName: "Company"},
		&mapper.ModelMapper{SrcName: "Model", WebName: "Model"},
		&mapper.ModelMapper{SrcName: "Location", WebName: "Location"},
		&mapper.ModelMapper{SrcName: "Parent", WebName: "Parent"},
		&mapper.ModelMapper{SrcName: "Rules", WebName: "Rules"},
	}
}

func (m *Device) ToWeb() (*omnimodels.Device, error) {
	web := new(omnimodels.Device)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (m Devices) ToWeb() ([]*omnimodels.Device, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.Device, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

///////////////////////////// DeviceGroup
type DeviceGroup struct {
	ID          int
	Name        string
	Description string
	Type        string
	Filters     datatypes.JSON
	CompanyID   int
	Company     *Company `gorm:"foreignKey:CompanyID"`
	Devices     Devices  `gorm:"many2many:device_group;joinForeignKey:GroupsID;JoinReferences:devices_id"`
	UserID      int
	User        *User `gorm:"foreignKey:UserID"`
}

type DeviceGroups []*DeviceGroup

func (m *DeviceGroup) GetModelMapper() []*mapper.ModelMapper {
	return []*mapper.ModelMapper{
		&mapper.ModelMapper{SrcName: "ID", WebName: "ID",
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
		&mapper.ModelMapper{SrcName: "Name", WebName: "Name"},
		&mapper.ModelMapper{SrcName: "Description", WebName: "Description"},
		&mapper.ModelMapper{SrcName: "Type", WebName: "Type"},
		&mapper.ModelMapper{SrcName: "Filters", WebName: "Filters",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				j, err := converters.JSONWebToSrc(web)
				if err != nil {
					return nil, fmt.Errorf("Filters: %v. %v", web, err)
				}
				return j, nil
			},
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				w, err := converters.JSONSrcToWeb(src)
				if err != nil {
					return nil, fmt.Errorf("Filters: %v. %v", src, err)
				}
				return w, nil
			},
		},
		&mapper.ModelMapper{SrcName: "Company", WebName: "Company"},
		&mapper.ModelMapper{SrcName: "Devices", WebName: "Devices"},
		&mapper.ModelMapper{SrcName: "User", WebName: "User"},
	}
}

func (m *DeviceGroup) ToWeb() (*omnimodels.DeviceGroup, error) {
	web := new(omnimodels.DeviceGroup)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (m DeviceGroups) ToWeb() ([]*omnimodels.DeviceGroup, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.DeviceGroup, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

/////////////////////////////// DeviceModel
type DeviceModel struct {
	ID             int `gorm:"primaryKey"`
	Name           string
	Title          string
	Description    string
	ManufacturerID int
	Manufacturer   *Manufacturer `gorm:"foreignkey:ManufacturerID"`
}

type DeviceModels []*DeviceModel

func (DeviceModel) TableName() string {
	return "device_models"
}

func (m *DeviceModel) GetModelMapper() []*mapper.ModelMapper {
	return []*mapper.ModelMapper{
		&mapper.ModelMapper{SrcName: "ID", WebName: "ID",
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
		&mapper.ModelMapper{SrcName: "Name", WebName: "Name"},
		&mapper.ModelMapper{SrcName: "Title", WebName: "Title"},
		&mapper.ModelMapper{SrcName: "Description", WebName: "Description"},
		&mapper.ModelMapper{SrcName: "Manufacturer", WebName: "Manufacturer"},
	}
}

func (m *DeviceModel) ToWeb() (*omnimodels.DeviceModel, error) {
	web := new(omnimodels.DeviceModel)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (m DeviceModels) ToWeb() ([]*omnimodels.DeviceModel, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*omnimodels.DeviceModel, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}
