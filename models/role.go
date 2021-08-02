package models

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

// Role define
var (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleAnonymous = "anonymous"
	RolesId       = map[string]int{
		RoleAdmin:     -1,
		RoleUser:      -1,
		RoleAnonymous: -1,
	}
)

// Role ...
type Role struct {
	Id    int64   `orm:"auto;pk" description:"role serial number" json:"role_id"`
	Name  string  `orm:"unique" description:"role name" json:"role_name"`
	Users []*User `orm:"reverse(many)" description:"user list" json:"users"`
}

func init() {
	orm.RegisterModel(new(Role))
	RegisterRoles()
	AddRolesGroupPolicy()
}

// RegisterRoles register role model-initialization
func RegisterRoles() {
	o := orm.NewOrm()
	// Here I write to the database by traversing a dictionary constructed above
	// If you are not willing to use Sao operation, just write three ReadOrCreate directly
	// The GetRoleString method is required
	for key := range RolesId {
		_, id, err := o.ReadOrCreate(&Role{Name: GetRoleString(key)}, "Name")
		if err != nil {
			panic(err)
		}
		RolesId[key] = int(id)
	}
}

// GetRoleString method is mainly used to add a role_ prefix to the Name field
func GetRoleString(s string) string {
	if strings.HasPrefix(s, "role_") {
		return s
	}
	return fmt.Sprintf("role_%s", s)
}

// AddRolesGroupPolicy add role inheritance policy rules to Casbin
func AddRolesGroupPolicy() {
	// Ordinary administrator inherits the user
	_, _ = Enforcer.AddGroupingPolicy(GetRoleString(RoleAdmin), GetRoleString(RoleUser))
	// The user inherits the anonymous
	_, _ = Enforcer.AddGroupingPolicy(GetRoleString(RoleUser), GetRoleString(RoleAnonymous))
}
