package rdb

import (
	"fmt"
)

var drivers = map[string]Driver{}

// Panics if called twice with the same name.
// Make the driver instance available clients.
func Register(name string, dr Driver) {
	_, found := drivers[name]
	if found {
		panic(fmt.Sprintf("Driver already present: %s", name))
	}
	drivers[name] = dr
}

func Open(config *Config) (Database, error) {
	dr, found := drivers[config.DriverName]
	if !found {
		return nil, fmt.Errorf("Driver name not found: %s", config.DriverName)
	}
	return dr.Open(config)
}
