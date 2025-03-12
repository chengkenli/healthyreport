package run

import (
	"encoding/csv"
	"fmt"
	"healthyreport/conn"
	"healthyreport/util"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SchemaTablet() {
	db, err := conn.StarRocks(util.P.App)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	/*information_schema.tables*/
	var Schema []map[string]interface{}
	all := db.Raw("select `TABLE_CATALOG`,`TABLE_SCHEMA`,`TABLE_NAME`,`TABLE_TYPE`,`ENGINE`,`CREATE_TIME`,`TABLE_COMMENT`,`TABLE_ROWS` from information_schema.tables").Scan(&Schema)
	if all.Error != nil {
		util.Logger.Error(all.Error.Error())
		return
	}
	var Procs []map[string]interface{}
	r := db.Raw("show proc '/dbs'").Scan(&Procs)
	if r.Error != nil {
		util.Logger.Error(r.Error.Error())
		return
	}

	util.Logger.Info(fmt.Sprintf("scan:%d，thread:%d\n正在读取...", all.RowsAffected, util.P.Thread))

	//var csvData = make([][]string, all.RowsAffected*100)
	var csvData [][]string
	var TableC int
	/*start*/
	var doneC = make(chan int)
	// 进度条
	tic := time.Tick(3 * time.Second)
	go func(c chan int) {
		for {
			select {
			case <-doneC:
				return
			case <-tic:
				processRate := float64(TableC) / float64(all.RowsAffected) * 100
				rate := strconv.FormatFloat(processRate, 'f', 2, 64)
				util.Logger.Info(fmt.Sprintf("(%5d,%5d) %s", TableC, all.RowsAffected, strings.Repeat("#", int(processRate/2))+rate+"%"))
				if TableC >= int(all.RowsAffected) {
					return
				}
			}
		}
	}(doneC)

	ch := make(chan struct{}, util.P.Thread)

	var wg sync.WaitGroup
	for i, schema := range Schema {
		if schema == nil {
			continue
		}
		wg.Add(1)
		go func(i int, schema map[string]interface{}) {
			defer func() {
				TableC++
				<-ch
				wg.Done()
			}()

			ch <- struct{}{}

			pool.Put(&informa{Schema: schema})
			p := pool.Get().(*informa)

			/*分析表类型*/
			var Type string
			for _, proc := range Procs {
				if p.Schema["TABLE_SCHEMA"].(string) == strings.ReplaceAll(proc["DbName"].(string), "default_cluster:", "") {
					/*判断表类型*/
					var m []map[string]interface{}
					r := db.Raw(fmt.Sprintf("show proc '/dbs/%s'", proc["DbId"].(string))).Scan(&m)
					if r.Error != nil {
						util.Logger.Error(r.Error.Error())
						break
					}
					for _, m2 := range m {
						if m2["TableName"].(string) == p.Schema["TABLE_NAME"].(string) {
							Type = m2["Type"].(string)
							break
						}
					}
					break
				}
			}

			if p.Schema["TABLE_TYPE"].(string) != "BASE TABLE" || Type != "OLAP" {
				return
			}

			/*容量*/
			var data []map[string]interface{}
			x := db.Raw(fmt.Sprintf("show data from %s.%s", p.Schema["TABLE_SCHEMA"].(string), p.Schema["TABLE_NAME"].(string))).Scan(&data)
			if x.Error != nil {
				util.Logger.Error(fmt.Sprintf("%s.%s %s -> %s", p.Schema["TABLE_SCHEMA"].(string), p.Schema["TABLE_NAME"].(string), p.Schema["TABLE_TYPE"].(string), x.Error.Error()))
				return
			}
			/*分区*/
			var partitions []map[string]interface{}
			r = db.Raw(fmt.Sprintf("show partitions from %s.%s order by LastConsistencyCheckTime,DataSize desc", p.Schema["TABLE_SCHEMA"].(string), p.Schema["TABLE_NAME"].(string))).Scan(&partitions)
			if r.Error != nil {
				util.Logger.Error(r.Error.Error())
				return
			}
			for _, partition := range partitions {
				atoi, _ := strconv.Atoi(partition["ReplicationNum"].(string))
				if atoi != util.P.ReplicaNum {
					continue
				}
				csvData = append(csvData, [][]string{{
					p.Schema["TABLE_SCHEMA"].(string),
					p.Schema["TABLE_NAME"].(string),
					data[0]["Size"].(string),
					partition["PartitionName"].(string),
					partition["DataSize"].(string),
					partition["ReplicationNum"].(string),
					"",
				}}...)
			}

		}(i, schema)
	}
	wg.Wait()

	util.Logger.Info(fmt.Sprintf("(%5d,%5d) %s", TableC, all.RowsAffected, strings.Repeat("#", 50)+"100%"))
	util.Logger.Info("数组数据写入文件.")

	// 创建一个新的CSV文件
	file, err := os.Create(nickname)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	/*写入UTF-8 BOM，防止中文乱码*/
	_, err = file.WriteString("\xEF\xBB\xBF")
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	defer file.Close()
	// 创建一个新的CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	arr := []string{"库名", "表名", "总容量", "分区名", "分区容量", "副本数", "备注"}
	err = writer.Write(arr)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}

	fmt.Println(len(csvData))
	for _, csvdata := range csvData {
		err := writer.Write(csvdata)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}
	}
	util.Logger.Info(nickname)
	util.Logger.Info("done.")
}
