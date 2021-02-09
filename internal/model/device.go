package model

type Device struct {
	ID          int                    `jsonapi:"primary,devices"`
	Name        string                 `jsonapi:"attr,name"`
	Slug        string                 `jsonapi:"attr,slug"`
	Title       string                 `jsonapi:"attr,title"`
	Description string                 `jsonapi:"attr,desc"`
	Kind        string                 `jsonapi:"attr,kind"`
	Info        map[string]interface{} `jsonapi:"attr,info"`
	Image       string                 `jsonapi:"attr,image"`
	Company     *Company               `jsonapi:"relation,company"`
	DevModel    *DeviceModel           `jsonapi:"relation,model"`
	Location    *Location              `jsonapi:"relation,location"`
	Groups      []*DeviceGroup         `jsonapi:"relation,groups"`
	Parent      *Device                `jsonapi:"relation,parent"`
	Rules       []*Rule                `jsonapi:"relation,rules"`
}

type DeviceGroup struct {
	ID          int                    `jsonapi:"primary,deviceGroups"`
	Name        string                 `jsonapi:"attr,name"`
	Description string                 `jsonapi:"attr,desc"`
	Type        string                 `jsonapi:"attr,type"`
	Filters     map[string]interface{} `jsonapi:"attr,filters"`
	Company     *Company               `jsonapi:"relation,company"`
	Devices     []*Device              `jsonapi:"relation,devices"`
	User        *User                  `jsonapi:"relation,user"`
}

type DeviceModel struct {
	ID           int           `jsonapi:"primary,deviceModels"`
	Name         string        `jsonapi:"attr,name"`
	Title        string        `jsonapi:"attr,title"`
	Description  string        `jsonapi:"attr,desc"`
	Manufacturer *Manufacturer `jsonapi:"relation,manufacturer"`
}
