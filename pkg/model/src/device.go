package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/web"
)

/////////////////////////////// Device
type Device struct {
	ID          int            `gorm:"primaryKey"`
	Name        string         `omni:"Name"`
	Slug        string         `omni:"Slug"`
	Title       string         `omni:"Title"`
	Description string         `omni:"Description"`
	Kind        string         `omni:"Kind"`
	Info        datatypes.JSON `omni:"Info;src:JSON2src;web:JSON2web"`
	Image       string         `omni:"Image"`
	CompanyID   int
	Company     *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
	ModelID     int
	Model       *DeviceModel `gorm:"foreignkey:ModelID" omni:"Model"`
	LocationID  int
	Location    *Location    `gorm:"foreignKey:LocationID" omni:"Location"`
	Groups      DeviceGroups `gorm:"many2many:device_group;joinForeignKey:DevicesID;JoinReferences:groups_id" omni:"Groups"`
	ParentID    int
	Parent      *Device `gorm:"foreignkey:ParentID" omni:"Parent"`
	Rules       Rules   `gorm:"many2many:rules_devices;joinForeignKey:DeviceID;JoinReferences:rule_id" omni:"Rules"`
}

type Devices []*Device

//func (m *Device) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Slug", WebName: "Slug"},
//		&mapper.ModelMap{SrcName: "Title", WebName: "Title"},
//		&mapper.ModelMap{SrcName: "Description", WebName: "Description"},
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
//		&mapper.ModelMap{SrcName: "Image", WebName: "Image"},
//		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
//		&mapper.ModelMap{SrcName: "Model", WebName: "Model"},
//		&mapper.ModelMap{SrcName: "Location", WebName: "Location"},
//		&mapper.ModelMap{SrcName: "Parent", WebName: "Parent"},
//		&mapper.ModelMap{SrcName: "Rules", WebName: "Rules"},
//	}
//}

func (m *Device) ToWeb(mapper *mapper.ModelMapper) (*webmodels.Device, error) {
	web := new(webmodels.Device)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Device) ScanFromWeb(web *webmodels.Device, mapper *mapper.ModelMapper) (*Device, error) {
	m := new(Device)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Devices) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.Device, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Device, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m Devices) ScanFromWeb(web []*webmodels.Device, mapper *mapper.ModelMapper) (Devices, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Device)
	res := make(Devices, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u, mapper)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}

///////////////////////////// DeviceGroup
type DeviceGroup struct {
	ID          int            `omni:"ID;src:ID2src;web:ID2web"`
	Name        string         `omni:"Name"`
	Description string         `omni:"Description"`
	Type        string         `omni:"Type"`
	Filters     datatypes.JSON `omni:"Filters;src:JSON2src;web:JSON2web"`
	CompanyID   int
	Company     *Company `gorm:"foreignKey:CompanyID" omni:"Company"`
	Devices     Devices  `gorm:"many2many:device_group;joinForeignKey:GroupsID;JoinReferences:devices_id" omni:"Devices"`
	UserID      int
	User        *User `gorm:"foreignKey:UserID" omni:"User"`
}

type DeviceGroups []*DeviceGroup

//func (m *DeviceGroup) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Filters", WebName: "Filters",
//			ConverterToSrc: func(web interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(web)
//				if err != nil {
//					return nil, fmt.Errorf("Filters: %v. %v", web, err)
//				}
//				return j, nil
//			},
//			ConverterToWeb: func(src interface{}) (interface{}, error) {
//				w, err := converters.JSONSrcToWeb(src)
//				if err != nil {
//					return nil, fmt.Errorf("Filters: %v. %v", src, err)
//				}
//				return w, nil
//			},
//		},
//		&mapper.ModelMap{SrcName: "Company", WebName: "Company"},
//		&mapper.ModelMap{SrcName: "Devices", WebName: "Devices"},
//		&mapper.ModelMap{SrcName: "User", WebName: "User"},
//	}
//}

func (m *DeviceGroup) ToWeb(mapper *mapper.ModelMapper) (*webmodels.DeviceGroup, error) {
	web := new(webmodels.DeviceGroup)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*DeviceGroup) ScanFromWeb(web *webmodels.DeviceGroup, mapper *mapper.ModelMapper) (*DeviceGroup, error) {
	m := new(DeviceGroup)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m DeviceGroups) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.DeviceGroup, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.DeviceGroup, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m DeviceGroups) ScanFromWeb(web []*webmodels.DeviceGroup, mapper *mapper.ModelMapper) (DeviceGroups, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(DeviceGroup)
	res := make(DeviceGroups, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u, mapper)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}

/////////////////////////////// DeviceModel
type DeviceModel struct {
	ID             int    `gorm:"primaryKey" omni:"ID;src:ID2src;web:ID2web"`
	Name           string `omni:"Name"`
	Title          string `omni:"Title"`
	Description    string `omni:"Description"`
	ManufacturerID int
	Manufacturer   *Manufacturer `gorm:"foreignkey:ManufacturerID" omni:"Manufacturer"`
}

type DeviceModels []*DeviceModel

func (DeviceModel) TableName() string {
	return "device_models"
}

//func (m *DeviceModel) GetModelMapper() []*mapper.ModelMap {
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
//		&mapper.ModelMap{SrcName: "Title", WebName: "Title"},
//		&mapper.ModelMap{SrcName: "Description", WebName: "Description"},
//		&mapper.ModelMap{SrcName: "Manufacturer", WebName: "Manufacturer"},
//	}
//}

func (m *DeviceModel) ToWeb(mapper *mapper.ModelMapper) (*webmodels.DeviceModel, error) {
	web := new(webmodels.DeviceModel)

	err := mapper.ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*DeviceModel) ScanFromWeb(web *webmodels.DeviceModel, mapper *mapper.ModelMapper) (*DeviceModel, error) {
	m := new(DeviceModel)
	err := mapper.ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m DeviceModels) ToWeb(mapper *mapper.ModelMapper) ([]*webmodels.DeviceModel, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.DeviceModel, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb(mapper)
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m DeviceModels) ScanFromWeb(web []*webmodels.DeviceModel, mapper *mapper.ModelMapper) (DeviceModels, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(DeviceModel)
	res := make(DeviceModels, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u, mapper)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
