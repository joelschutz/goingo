package component

import "github.com/yohamta/donburi"

type ConfigurationData struct {
	BoardSize uint
}

var Configuration = donburi.NewComponentType[ConfigurationData](
	ConfigurationData{BoardSize: 9},
)
