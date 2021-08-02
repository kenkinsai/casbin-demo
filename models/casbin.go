package models

import (
	"casbin-demo/conf"
	"casbin-demo/pkg/database"
	"runtime"

	"github.com/beego/beego/v2/client/orm"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

// Enforcer note that this Enforcer is very important, Casbin is used to call this variable
var Enforcer *casbin.Enforcer

// CasbinRule ...
type CasbinRule struct {
	Id    int64  `orm:"auto"`      // Increase primary key
	Ptype string `orm:"size(128)"` // Policy Type-used to distinguish policy and group (role)
	V0    string `orm:"size(128)"` // subject
	V1    string `orm:"size(128)"` // object
	V2    string `orm:"size(128)"` // action
	V3    string `orm:"size(128)"` // This and the following fields are useless, only reserved, if yours is not
	V4    string `orm:"size(128)"` // Sub, obj, act will only be used
	V5    string `orm:"size(128)"` // Such as sub, obj, act, suf will use V3
}

func init() {
	orm.RegisterModel(new(CasbinRule))
	// init database
	err := database.InitConnection(conf.SQLConf())
	if err != nil {
		panic(err)
	}
	RegisterCasbin()
}

// Adapter ...
type Adapter struct {
	o orm.Ormer
}

// RegisterCasbin ...
func RegisterCasbin() {
	a := &Adapter{}
	a.o = orm.NewOrm()
	// I don't know why
	runtime.SetFinalizer(a, finalizer)
	// Enforcer initialization-that is passed into the Adapter object
	Enforcer, _ = casbin.NewEnforcer("conf/casbin.conf", a)
	// Enforcer reads Policy
	err := Enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
}

// Note that the specific code corresponding to the method should be copied from beego-ORM-Adapter / adapter.go
// The orm operation used in this method is still based on your own
// Make adjustments to the actual situation, do not copy blindly
func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.Ptype
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}
func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.Ptype = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// LoadPolicy ...
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	_, err := a.o.QueryTable("casbin_rule").All(&lines)
	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

// SavePolicy ...
func (a *Adapter) SavePolicy(model model.Model) error {
	var lines []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	_, err := a.o.InsertMulti(len(lines), lines)
	return err
}

// AddPolicy ...
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Insert(&line)
	return err
}

// RemovePolicy ...
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Delete(&line, "ptype", "v0", "v1", "v2", "v3", "v4", "v5")
	return err
}

// RemoveFilteredPolicy ...
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{}

	line.Ptype = ptype
	filter := []string{}
	filter = append(filter, "ptype")
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
		filter = append(filter, "v0")
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
		filter = append(filter, "v1")
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
		filter = append(filter, "v2")
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
		filter = append(filter, "v3")
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
		filter = append(filter, "v4")
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
		filter = append(filter, "v5")
	}

	_, err := a.o.Delete(&line, filter...)
	return err
}
