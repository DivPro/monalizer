package db

import (
	"sync"
)

type Module struct {
	ID        uint8
	Name      string
	GoVersion string
	mx        sync.Mutex
	Requires  map[string]*Require
	UsedBy    map[string]*Module
}

func (m *Module) Require(mod *Module, ver string, indirect bool) {
	m.mx.Lock()
	m.Requires[mod.Name] = &Require{
		Module:   mod,
		Version:  ver,
		Indirect: indirect,
	}
	m.mx.Unlock()
	mod.Use(m)
}

func (m *Module) Use(mod *Module) {
	m.mx.Lock()
	m.UsedBy[mod.Name] = mod
	m.mx.Unlock()
}

type Require struct {
	*Module
	Version  string
	Indirect bool
}

type DB struct {
	mx      sync.Mutex
	Modules map[string]*Module
	idSeq   uint8
}

func New() *DB {
	return &DB{
		Modules: map[string]*Module{},
	}
}

func (db *DB) Module(name string) *Module {
	db.mx.Lock()
	if m, ok := db.Modules[name]; ok {
		db.mx.Unlock()
		return m
	}
	db.idSeq++
	m := &Module{
		ID:       db.idSeq,
		Name:     name,
		Requires: make(map[string]*Require),
		UsedBy:   make(map[string]*Module),
	}
	db.Modules[name] = m
	db.mx.Unlock()

	return m
}
