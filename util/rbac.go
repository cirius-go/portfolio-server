package util

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
)

// NewRBACModel initializes the RBAC casbin model
func NewRBACModel() model.Model {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")`)
	return m
}

// NewRBACWithLevelInheritanceModel initializes the RBAC with level inheritance model
func NewRBACWithLevelInheritanceModel() model.Model {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, lvl, obj, act")
	m.AddDef("p", "p", "sub, lvl, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("g", "g2", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `((r.sub != p.sub && g(r.sub, p.sub)) || (r.sub == p.sub && g2(p.lvl, r.lvl))) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")`)
	return m
}

// NewRBACWithDomainModel initializes the RBAC with domain model
func NewRBACWithDomainModel() model.Model {
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act")
	m.AddDef("p", "p", "sub, dom, obj, act")
	m.AddDef("g", "g", "_, _, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`)
	return m
}
