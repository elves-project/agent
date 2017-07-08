package main

import (
	"flag"
	"fmt"
	"github.com/elves-project/agent/src/cron"
	"github.com/elves-project/agent/src/funcs"
	"github.com/elves-project/agent/src/g"
	"github.com/elves-project/agent/src/http"
	"github.com/gy-games-libs/seelog"
	"os"
)

func main() {
	defer seelog.Flush()
	cfg := flag.String("r", g.GetRoot(), "elves-agent root directory")
	version := flag.Bool("v", false, "show version")
	clear := flag.Bool("clear", false, "clear this agent's apps")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	if *clear {
		funcs.ClearApps()
		fmt.Println("elves-agent apps cleared..")
		os.Exit(0)
	}

	g.InitStat()

	if g.Config().HeartBeat.Enabled {
		go cron.HearBeatCron(int64(g.Config().HeartBeat.Interval))
	}

	if g.Config().Devmode.Enabled {
		seelog.Info("elves agent init as dev mode")
	}

	if g.Config().Http.Enabled || g.Config().Devmode.Enabled {
		go http.Start()
	}

	if g.Config().Agentcron.Enabled {
		go cron.WatchCron()
	}

	funcs.ServerRun(g.Config().AuthIps)
}
