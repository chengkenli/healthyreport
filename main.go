package main

import (
	_ "healthyreport/init"
	"healthyreport/run"
	"healthyreport/util"
)

func main() {
	util.Parms()
	util.Loggers()
	run.Run()
}
