package plugInModule

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
)

type PluginRegestry struct {
	CommandPlugins []ExtendCommandPlugin
	ThemePlugin    ThemePlugin
}

func getOSLibrary() string {
	if runtime.GOOS == "windows" {
		return ".dll"
	} else {
		return ".os"
	}
}

func (pr *PluginRegestry) loadPlugins(pluginFolder string) {
	err := filepath.Walk(pluginFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if "" == getOSLibrary() {
			log.Fatal("Your operation system is not supported")
		}
		if !info.IsDir() && filepath.Ext(path) == getOSLibrary() {
			fmt.Printf("")
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error until reading plugin. %s", err)
	}
}
