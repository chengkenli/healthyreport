package util

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Printf("\nUsage: %s [-s starrocks] [-h]\n%sstarrocks 表属性扫描工具1.2\n\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Println()
}

func init() {
	c := color.New()
	flag.IntVar(&P.ReplicaNum, "r", -1, c.Add(color.FgHiWhite).Sprint("REPLICA NUM"))
	flag.IntVar(&P.Bucket, "b", -1, c.Add(color.FgHiWhite).Sprint("BUCKET NUM"))
	flag.IntVar(&P.Thread, "t", 50, "THREAD")
	flag.StringVar(&P.App, "s", "", fmt.Sprintf("<%s>", c.Add(color.FgHiWhite).Sprint("APP")))
	flag.StringVar(&P.File, "f", "", "FILE")
	flag.StringVar(&P.Database, "d", "", "DATABASE")
	flag.BoolVar(&P.Visits, "v", false, c.Add(color.FgHiRed).Sprint("VISITS"))
	flag.BoolVar(&P.Help, "h", false, "HELP")

	flag.Parse()
	flag.Usage = usage

	if P.Help || P.App == "" {
		flag.Usage()
		os.Exit(-1)
	}

}

func Parms() {
}
