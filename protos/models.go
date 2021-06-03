package protos

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	null "gopkg.in/guregu/null.v3/zero"
)

type User struct {
	UID        uint64       `json:"uid,omitempty" validate:"-" db:"uid"`
	TenantID   uint64       `json:"tenant_id,omitempty" validate:"-" db:"tenant_id"`
	Password   string       `json:"password,omitempty" validate:"required,min=6,max=256" db:"password"`
	Cellphone  *null.String `json:"cellphone,omitempty" validate:"omitempty,phone" db:"cellphone"`
	Email      *null.String `json:"email,omitempty" validate:"omitempty,email" db:"email"`
	Nickname   *null.String `json:"nickname,omitempty" validate:"omitempty,min=2,max=64" db:"nickname"`
	AvatarURL  *null.String `json:"avatarUrl,omitempty" db:"avatar_url"`
	Addr       *null.String `json:"addr,omitempty" db:"addr"`
	Gender     *null.Int    `json:"gender,omitempty" db:"gender"`
	AddTime    *time.Time   `json:"addTime,omitempty" validate:"-" db:"add_time"`
	UpdateTime *time.Time   `json:"updateTime,omitempty" validate:"-" db:"update_time"`
	DeleteTime *time.Time   `json:"deleteTime,omitempty" validate:"-" db:"delete_time"`
	LoginTime  *time.Time   `json:"loginTime,omitempty" validate:"-" db:"login_time"`

	Ext    MapStruct `json:"ext,omitempty" validate:"-" db:"ext"`
	Tenant *Tenant   `json:"tenant,omitempty" validate:"-" db:"tenant"`
}

// 租户
type Tenant struct {
	ID            uint64               `json:"id" validate:"-" db:"id"`
	UID           uint64               `json:"uid,omitempty" validate:"-" db:"uid"`
	TenantName    string               `json:"tenant_name" db:"tenant_name"`
	TenantType    string               `json:"tenant_type" db:"tenant_type"`
	AddTime       *time.Time           `json:"addTime,omitempty" validate:"-" db:"add_time"`
	UpdateTime    *time.Time           `json:"updateTime,omitempty" validate:"-" db:"update_time"`
	Info          MapStruct            `json:"info,omitempty" db:"info"`
	Configuration *TenantConfiguration `json:"configuration,omitempty" db:"configuration"`
}

// 租户配置字段
type TenantConfiguration struct {
	Roles []RoleStruct `json:"roles"` // 用户角色字典列表
	More  MapStruct    `json:"more"`
}

func (t *TenantConfiguration) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	if len(src.([]byte)) <= 2 {
		return nil
	}

	b, _ := src.([]byte)
	return json.Unmarshal(b, t)
}
func (t TenantConfiguration) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type UserReq struct {
	UID       uint64 `json:"uid,omitempty" validate:"-"`
	TenantID  uint64 `json:"tenant_id" validate:"-"`
	Cellphone string `json:"cellphone,omitempty" validate:"omitempty,phone"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	Nickname  string `json:"nickname,omitempty" validate:"omitempty,min=2,max=64"`
	Password  string `json:"password,omitempty" validate:"omitempty,min=6,max=256"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	Addr      string `json:"addr,omitempty"`
	Gender    int32  `json:"gender,omitempty"`
}

type MapStruct map[string]interface{}

func (t *MapStruct) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	if len(src.([]byte)) <= 2 {
		return nil
	}

	b, _ := src.([]byte)
	return json.Unmarshal(b, t)
}
func (t MapStruct) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type PolicyReq struct {
	Role string `json:"role" validate:"required,min=2,max=128"`
	Obj  string `json:"obj" validate:"required,min=2,max=128"`
	Act  string `json:"act" validate:"required,min=2,max=128"`
}

type RoleStruct struct {
	RoleTitle string `json:"title" validate:"max=10"`
	RoleValue string `json:"value" validate:"max=10"`

	UID uint64 `json:"uid,omitempty" validate:"-"`
}
