package model

import (
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"strconv"
)

type Location struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Timezone string
	//Info     map[string]interface{}
	CompanyID int
	Company   *Company `gorm:"foreignKey:CompanyID"`
	ParentID  int
	Children  []*Location `gorm:"foreignkey:ParentID"`
	Users     []*User     `gorm:"foreignkey:LocationID"`
}

func (m Location) ToWeb() *omnimodels.Location {
	web := new(omnimodels.Location)

	web.ID = strconv.Itoa(m.ID)
	web.Name = m.Name
	web.Timezone = m.Timezone
	//info
	if m.Company != nil {
		web.Company = m.Company.ToWeb()
	}
	web.Children = LocationsToWeb(m.Children)
	web.Users = UsersToWeb(m.Users)

	return web
}

func LocationsToWeb(mSl []*Location) []*omnimodels.Location {
	if mSl == nil {
		return nil
	}
	omniRecs := make([]*omnimodels.Location, 0, 5)
	for _, s := range mSl {
		omniRecs = append(omniRecs, s.ToWeb())
	}
	return omniRecs
}
