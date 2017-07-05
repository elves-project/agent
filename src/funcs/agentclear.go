package funcs

import (
	"github.com/elves-project/agent/src/g"
)

func ClearApps() {
	g.Config().Apps = map[string]string{}
	g.SaveConfig()
}
