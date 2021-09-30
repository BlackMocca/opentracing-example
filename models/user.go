package models

import (
	"math/rand"
	"time"

	"git.innovasive.co.th/backend/helper"
	"github.com/gofrs/uuid"
	"github.com/spf13/cast"
)

const FIELD_FK_USER_TYPE = "UserType"

type User struct {
	TableName struct{} `json:"-" db:"users" pk:"Id"`

	Id         string     `json:"id" db:"id" type:"string"`
	Name       string     `json:"name" db:"name" type:"string"`
	Cover      string     `json:"cover" db:"-"`
	Address    string     `json:"address" db:"-"`
	UserTypeId *uuid.UUID `json:"user_type_id" db:"user_type_id" type:"uuid"`
	UserType   *UserType  `json:"user_type" db:"-" fk:"relation:one,fk_field1:UserTypeId,fk_field2:Id"`
}

type Users []*User

func NewUserWithParams(params map[string]interface{}, ptr *User) *User {
	if ptr == nil {
		ptr = new(User)
	}

	for key, v := range params {
		switch key {
		case "id":
			ptr.Id = cast.ToString(v)
		case "name":
			ptr.Name = cast.ToString(v)
		case "user_type_id":
			ptr.UserTypeId, _ = helper.ConvertToUUIDAndBinary(v)
		}
	}

	return ptr
}

func (u *User) GenUUID() {
	rand.Seed(time.Now().Unix())
	id := rand.Intn(100)
	u.Id = cast.ToString(id)
}
