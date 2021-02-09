package model

import (
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"strconv"
)

type Role struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Slug       string
	Assigned   bool
	Persistent bool
	Users      []*User `gorm:"many2many:user_role;joinForeignKey:RolesID;JoinReferences:users_id"`
	//Info       map[string]interface{} `jsonapi:"attr,info"`
	CompanyID int
	Company   *Company `gorm:"foreignKey:CompanyID"`
}

func (m *Role) ToWeb() *omnimodels.Role {
	web := new(omnimodels.Role)

	web.ID = strconv.Itoa(m.ID)
	web.Name = m.Name
	web.Slug = m.Slug
	web.Assigned = m.Assigned
	web.Persistent = m.Persistent
	//Users...
	if m.Company != nil {
		web.Company = m.Company.ToWeb()
	}
	return web
}

func RolesToWeb(mSl []*Role) []*omnimodels.Role {
	if mSl == nil {
		return nil
	}

	omniM := make([]*omnimodels.Role, 0, 5)
	for _, u := range mSl {
		omniM = append(omniM, u.ToWeb())
	}
	return omniM
}
