package model

import (
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"strconv"
)

type Company struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func (Company) TableName() string {
	return "companies"
}

func (m Company) ToWeb() *omnimodels.Company {
	web := new(omnimodels.Company)

	web.ID = strconv.Itoa(m.ID)
	web.Name = m.Name

	return web
}
