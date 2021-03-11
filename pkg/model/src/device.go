package src

import (
	"gorm.io/datatypes"
	"omnimanage/pkg/mapper"
	webmodels "omnimanage/pkg/model/domain"
)

/////////////////////////////// Device
type Device struct {
	ID          int            `gorm:"primaryKey"`
	Name        string         `omni:"Name"`
	Slug        string         `omni:"Slug"`
	Title       string         `omni:"Title"`
	Description string         `omni:"Description"`
	Kind        string         `omni:"Kind"`
	Info        datatypes.JSON `omni:"Info;src:JSON2src;domain:JSON2web"`
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
//		&mapper.ModelMap{SrcName: "Slug", WebName: "Slug"},
//		&mapper.ModelMap{SrcName: "Title", WebName: "Title"},
//		&mapper.ModelMap{SrcName: "Description", WebName: "Description"},
//		&mapper.ModelMap{SrcName: "Info", WebName: "Info",
//			ConverterToSrc: func(domain interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(domain)
//				if err != nil {
//					return nil, fmt.Errorf("Settings: %v. %v", domain, err)
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

func (m *Device) ToWeb() (*webmodels.Device, error) {
	web := new(webmodels.Device)

	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*Device) ScanFromWeb(web *webmodels.Device) (*Device, error) {
	m := new(Device)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Devices) ToWeb() ([]*webmodels.Device, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.Device, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m Devices) ScanFromWeb(web []*webmodels.Device) (Devices, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(Device)
	res := make(Devices, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}

///////////////////////////// DeviceGroup
type DeviceGroup struct {
	ID          int            `omni:"ID;src:ID2src;domain:ID2web"`
	Name        string         `omni:"Name"`
	Description string         `omni:"Description"`
	Type        string         `omni:"Type"`
	Filters     datatypes.JSON `omni:"Filters;src:JSON2src;domain:JSON2web"`
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
//		&mapper.ModelMap{SrcName: "Filters", WebName: "Filters",
//			ConverterToSrc: func(domain interface{}) (interface{}, error) {
//				j, err := converters.JSONWebToSrc(domain)
//				if err != nil {
//					return nil, fmt.Errorf("Filters: %v. %v", domain, err)
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

func (m *DeviceGroup) ToWeb() (*webmodels.DeviceGroup, error) {
	web := new(webmodels.DeviceGroup)

	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*DeviceGroup) ScanFromWeb(web *webmodels.DeviceGroup) (*DeviceGroup, error) {
	m := new(DeviceGroup)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m DeviceGroups) ToWeb() ([]*webmodels.DeviceGroup, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.DeviceGroup, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m DeviceGroups) ScanFromWeb(web []*webmodels.DeviceGroup) (DeviceGroups, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(DeviceGroup)
	res := make(DeviceGroups, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}

/////////////////////////////// DeviceModel
type DeviceModel struct {
	ID             int    `gorm:"primaryKey" omni:"ID;src:ID2src;domain:ID2web"`
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
//		&mapper.ModelMap{SrcName: "Title", WebName: "Title"},
//		&mapper.ModelMap{SrcName: "Description", WebName: "Description"},
//		&mapper.ModelMap{SrcName: "Manufacturer", WebName: "Manufacturer"},
//	}
//}

func (m *DeviceModel) ToWeb() (*webmodels.DeviceModel, error) {
	web := new(webmodels.DeviceModel)

	err := mapper.Get().ConvertSrcToWeb(m, &web)
	if err != nil {
		return nil, err
	}

	return web, nil
}

func (*DeviceModel) ScanFromWeb(web *webmodels.DeviceModel) (*DeviceModel, error) {
	m := new(DeviceModel)
	err := mapper.Get().ConvertWebToSrc(web, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m DeviceModels) ToWeb() ([]*webmodels.DeviceModel, error) {
	if m == nil {
		return nil, nil
	}
	omniM := make([]*webmodels.DeviceModel, 0, 5)
	for _, u := range m {
		webU, err := u.ToWeb()
		if err != nil {
			return nil, err
		}
		omniM = append(omniM, webU)
	}
	return omniM, nil
}

func (m DeviceModels) ScanFromWeb(web []*webmodels.DeviceModel) (DeviceModels, error) {
	if len(web) == 0 {
		return nil, nil
	}

	srcPoint := new(DeviceModel)
	res := make(DeviceModels, 0, len(web))
	for _, u := range web {
		srcRec, err := srcPoint.ScanFromWeb(u)
		if err != nil {
			return nil, err
		}
		res = append(res, srcRec)
	}
	return res, nil
}
