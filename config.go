package main

import (
	"github.com/bitmark-inc/logger"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

//StaticIdentity for save static node
type StaticIdentity struct {
	UseStatic bool   `gluamapper:"use_this" json:"use_this"`
	KeyFile   string `gluamapper:"private_key_file" json:"private_key_file"`
}
type config struct {
	DataDirectory  string               `gluamapper:"data_directory" json:"data_directory"`
	Logging        logger.Configuration `gluamapper:"logging" json:"loggin"`
	NodeType       int                  `gluamapper:"node_type" json:"node_type"`
	PublicIP       []string             `gluamapper:"public_ip" json:"public_ip"`
	Port           int                  `gluamapper:"port" json:"port"`
	StaticIdentity `gluamapper:"static_identity" json:"static_identity"`
}

// ParseConfigurationFile - read and execute a Lua files and assign
// the results to a configuration structure
func ParseConfigurationFile(fileName string, config interface{}) error {
	L := lua.NewState()
	defer L.Close()

	L.OpenLibs()
	arg := &lua.LTable{}
	arg.Insert(0, lua.LString(fileName))
	L.SetGlobal("arg", arg)

	// execute configuration
	if err := L.DoFile(fileName); err != nil {
		return err
	}

	mapperOption := gluamapper.Option{
		NameFunc: func(s string) string {
			return s
		},
		TagName: "gluamapper",
	}
	mapper := gluamapper.Mapper{Option: mapperOption}
	err := mapper.Map(L.Get(L.GetTop()).(*lua.LTable), config)
	return err
}
