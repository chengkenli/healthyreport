package run

import (
	"fmt"
	"healthyreport/util"
	"sync"
	"time"
)

var nickname string

func _init() {
	var nick string

	for _, m := range util.MetaLink {
		if m["app"].(string) == util.P.App {
			nick = m["nickname"].(string)
			break
		}
	}
	nickname = fmt.Sprintf("%s-%s.csv", time.Now().Format("2006-01-02"), nick)

	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Pool ...")
			return new(informa)
		},
	}
}

func Run() {
	_init()
	/*初始化主题域*/
	SchemaOneStopZty()
	if util.P.Bucket > 0 {
		SchemaBucket()
		return
	}
	if util.P.ReplicaNum > 0 {
		SchemaTablet()
		return
	}
	if util.P.File != "" {
		SchemaFile()
		return
	}
	SchemaAll()
}
