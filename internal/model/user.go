package model

import (
	"fmt"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"strconv"
)

type User struct {
	ID            int    `gorm:"primaryKey"`
	Name          string `gorm:"column:user_name"`
	Password      string
	FirstName     string
	LastName      string
	PhoneNumber   string
	Email         string
	Image         string
	Settings      datatypes.JSON
	CompanyID     int
	Company       *Company `gorm:"foreignKey:CompanyID"`
	LocationID    int
	Location      *Location       `gorm:"foreignKey:LocationID"`
	Roles         []*Role         `gorm:"many2many:user_role;joinForeignKey:UsersID;JoinReferences:roles_id"`
	Subscriptions []*Subscription `gorm:"foreignKey:UserID"`
}

func (m *User) ToWeb() *omnimodels.User {
	web := new(omnimodels.User)

	web.ID = strconv.Itoa(m.ID)
	web.Name = m.Name
	web.Password = m.Password
	web.FirstName = m.FirstName
	web.LastName = m.LastName
	web.PhoneNumber = m.PhoneNumber
	web.Email = m.Email
	web.Image = m.Image
	//m.Settings.S
	//json.Unmarshal()

	if m.Company != nil {
		web.Company = m.Company.ToWeb()
	}
	if m.Location != nil {
		web.Location = m.Location.ToWeb()
	}

	web.Roles = RolesToWeb(m.Roles)
	//web.Subscriptions
	return web
}

func (m *User) ScanFromWeb(us *omnimodels.User) error {
	var err error
	m.ID, err = strconv.Atoi(us.ID)
	if err != nil {
		return fmt.Errorf("Wrong User ID: %v", us.ID)
	}

	m.Name = us.Name
	m.Password = us.Password
	m.FirstName = us.FirstName
	m.LastName = us.LastName
	m.PhoneNumber = us.PhoneNumber
	m.Email = us.Email
	m.Image = us.Image
	//....

	return nil
}

func UsersToWeb(mSl []*User) []*omnimodels.User {
	if mSl == nil {
		return nil
	}
	omniM := make([]*omnimodels.User, 0, 5)
	for _, u := range mSl {
		omniM = append(omniM, u.ToWeb())
	}
	return omniM
}
