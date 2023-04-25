package component

import (
	"io/fs"

	"github.com/yohamta/donburi"
)

type ConfigurationData struct {
	BoardSize int
	AIEnabled bool
	Assets    fs.FS
}

var Configuration = donburi.NewComponentType[ConfigurationData](
	ConfigurationData{BoardSize: 9},
)
