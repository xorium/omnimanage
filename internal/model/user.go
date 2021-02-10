package model

import (
	"fmt"
	"github.com/fatih/structs"
	omnimodels "gitlab.omnicube.ru/libs/omnilib/models"
	"gorm.io/datatypes"
	"reflect"
	"strconv"
)

type User struct {
	ID            int    `gorm:"primaryKey" omni:"func:FillID"`
	Name          string `gorm:"column:user_name" omni:"attr,Name"`
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

type Mapper struct {
	SrcName        string
	WebName        string
	ConverterToSrc func(web interface{}) (interface{}, error)
	ConverterToWeb func(web interface{}) (interface{}, error)
}

func (m *User) GetMapper() []*Mapper {
	return []*Mapper{
		&Mapper{SrcName: "ID", WebName: "",
			ConverterToSrc: func(web interface{}) (interface{}, error) {
				us, ok := web.(*omnimodels.User)
				if !ok {
					return nil, fmt.Errorf("Wrong type...")
				}
				id, err := strconv.Atoi(us.ID)
				if err != nil {
					return nil, fmt.Errorf("Wrong User ID: %v", us.ID)
				}
				return id, nil
			},
		},
		&Mapper{SrcName: "", WebName: "ID",
			ConverterToWeb: func(src interface{}) (interface{}, error) {
				us, ok := src.(*User)
				if !ok {
					return nil, fmt.Errorf("Wrong type...")
				}
				id := strconv.Itoa(us.ID)
				return id, nil
			},
		},
		&Mapper{SrcName: "Name", WebName: "Name"},
		&Mapper{SrcName: "Company", WebName: "Company"},
		&Mapper{SrcName: "Location", WebName: "Location"},
	}
}

func (m *User) ToWebM() (*omnimodels.User, error) {
	web := new(omnimodels.User)

	webS := structs.New(web)
	srcS := structs.New(m)
	srcM := srcS.Map()

	mapper := m.GetMapper()
	for _, val := range mapper {
		if val.WebName == "" {
			continue
		}

		webField := webS.Field(val.WebName)
		if val.SrcName != "" {
			srcField, ok := srcM[val.SrcName]
			if !ok {
				return nil, fmt.Errorf("unknown src field %v", val.SrcName)
			}
			if srcField == nil {
				continue
			}

			//Relation
			typeKind := webField.Kind()
			if typeKind == reflect.Ptr {
				//i2s(srcField, webField.Value())

				srcField := srcS.Field(val.SrcName)
				resInv, _ := Invoke(srcField.Value(), "ToWeb")
				webField.Set(resInv)
			} else { //Simple Attribute
				webField.Set(srcField)
			}
			continue
		}

		if val.ConverterToWeb != nil {
			srcField, err := val.ConverterToWeb(m)
			if err != nil {
				return nil, err
			}
			webField.Set(srcField)
			continue
		}

		return nil, fmt.Errorf("wrong mapper line %v", val)
	}
	return web, nil
}

//func (m *User) FillIDFromWeb(web interface{}) (interface{}, error) {
//	s := structs.New(web)
//	s.Map()
//
//}

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

func Invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}
	return method.Call(in)[0], nil
}
