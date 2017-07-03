package funcs

import (
	"../g"
)

func ClearApps() {
	g.Config().Apps = map[string]string{}
	g.SaveConfig()
}
