package run

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/grd/statistics"
	"gorm.io/gorm"
	"healthyreport/conn"
	"healthyreport/tools"
	"healthyreport/util"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Datas struct {
	Db                *gorm.DB
	Version           float64
	Dbname, TableName string
}

var csvDatas = make(map[string][][]string)
var ztylists []string

func SchemaFile() {
	db, err := conn.StarRocks(util.P.App)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	// 打开文件
	file, err := os.Open(util.P.File)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	defer file.Close()

	// 创建一个scanner
	scanner := bufio.NewScanner(file)
	var i int
	var wg sync.WaitGroup
	ch := make(chan struct{}, util.P.Thread)
	for scanner.Scan() {
		i++
		ss := strings.Split(strings.NewReplacer(" ", "").Replace(scanner.Text()), ".")
		if len(ss) <= 1 {
			continue
		}
		wg.Add(1)
		go func(i int, ss []string) {
			defer func() {
				<-ch
				wg.Done()
			}()

			ch <- struct{}{}

			util.Logger.Info(fmt.Sprintf("#%00005d <%s> %d", i, ss, len(ss)))
			ScanIsData(
				&Datas{
					Db:        db,
					Version:   chversion(db),
					Dbname:    ss[0],
					TableName: ss[1],
				})
			//ztylist = append(ztylist, z...)
			//csvData
		}(i, ss)
	}
	wg.Wait()

	util.Logger.Info("数组数据写入文件.")
	// 创建一个新的CSV文件
	filename := fmt.Sprintf("diy-%s.csv", util.P.File)
	util.Logger.Info(filename)
	file, err = os.Create(filename)
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

	arr := []string{"库名", "表名", "主题域", "类型", "建表模型", "注释", "创建时间", "数据量", "总容量/MB", "总容量/单位", "责任人", "总分区数", "实际分区数(有数据)", "空存储分区数(无数据)", "空存储分区范围", "总副本数", "副本数", "分桶数", "总副本数平均差", "总副本数标准差", "近一个月访问次数", "近三个月访问次数", "备注"}
	err = writer.Write(arr)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}

	for _, zty := range tools.RemoveDuplicateStrings(ztylists) {
		aa := csvDatas[zty]
		fmt.Println(len(aa))
		for _, csvs := range aa {
			err := writer.Write(csvs)
			if err != nil {
				util.Logger.Error(err.Error())
				return
			}
		}
	}
	util.Logger.Info("done.")
}

func ScanIsData(d *Datas) {
	var schema map[string]interface{}
	r := d.Db.Raw(fmt.Sprintf("select `TABLE_CATALOG`,`TABLE_SCHEMA`,`TABLE_NAME`,`TABLE_TYPE`,`ENGINE`,`CREATE_TIME`,`TABLE_COMMENT`,`TABLE_ROWS` from information_schema.tables where TABLE_SCHEMA='%s' and TABLE_NAME='%s'", d.Dbname, d.TableName)).Scan(&schema)
	if r.Error != nil {
		util.Logger.Error(r.Error.Error())
		return
	}
	var OwnerId []string
	var zty string
	var tableRows int64

	if r.RowsAffected == 0 {
		ztylists = append(ztylists, zty)
		csvDatas[zty] = append(csvDatas[zty], []string{
			d.Dbname,    /*库名*/
			d.TableName, /*表名*/
			"",          /*主题域*/
			"",          /*表类型*/
			"",          /*建表模型*/
			"",          /*表注释*/
			"",          /*创建时间*/
			"",          /*数据量*/
			"0",         /*总容量/MB*/
			"0",         /*总容量/单位*/
			"",          /*责任人*/
			"0",         /*总分区数*/
			"0",         /*实际分区数(有数据)*/
			"0",         /*空存储分区数(无数据)*/
			"0",         /*空存储分区范围*/
			"0",         /*总副本数*/
			"0",         /*副本数*/
			"0",         /*分桶数*/
			"0",         /*总副本数平均差*/
			"0",         /*总副本数标准差*/
			"",          /*备注*/
		})
		return
	}

	//tableName = fmt.Sprintf("%05d %s.%s", i, schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string))
	//util.Logger.Info(tableName)
	/*适配责任人*/
	if schema == nil {
		return
	}
	if schema["TABLE_SCHEMA"].(string) == "" {
		util.Logger.Info(fmt.Sprintf("%v", schema))
		return
	}
	/*过滤系统表*/
	//if schema["TABLE_SCHEMA"].(string) == "information_schema" ||
	//	schema["TABLE_SCHEMA"].(string) == "ops" ||
	//	schema["TABLE_SCHEMA"].(string) == "demo" ||
	//	schema["TABLE_SCHEMA"].(string) == "audit" {
	//	return
	//}
	/*适配责任人*/
	if schema["TABLE_SCHEMA"].(string) == "ods" ||
		schema["TABLE_SCHEMA"].(string) == "ods_app_dev" ||
		schema["TABLE_SCHEMA"].(string) == "ods_app_dev_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_app_test" ||
		schema["TABLE_SCHEMA"].(string) == "ods_app_test_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_dev" ||
		schema["TABLE_SCHEMA"].(string) == "ods_dev_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_rt" ||
		schema["TABLE_SCHEMA"].(string) == "ods_rt_dev" ||
		schema["TABLE_SCHEMA"].(string) == "ods_rt_dev_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_rt_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_sox" ||
		schema["TABLE_SCHEMA"].(string) == "ods_sox_app_dev" ||
		schema["TABLE_SCHEMA"].(string) == "ods_sox_app_test" ||
		schema["TABLE_SCHEMA"].(string) == "ods_sox_dev" ||
		schema["TABLE_SCHEMA"].(string) == "ods_sox_test" ||
		schema["TABLE_SCHEMA"].(string) == "ods_test" ||
		schema["TABLE_SCHEMA"].(string) == "ods_test_secure" ||
		schema["TABLE_SCHEMA"].(string) == "ods_secure_rt" {

		OwnerId = append(OwnerId, "i0l02bg")
		zty = "ods"
	} else if strings.Contains(schema["TABLE_SCHEMA"].(string), "flash_report") {
		OwnerId = append(OwnerId, "k0z01dc")
		zty = "fin"
	} else {
		ttl := strings.Split(schema["TABLE_NAME"].(string), "_")
		if len(ttl) > 2 {
			zty = ttl[1]
			for _, d := range util.Domain {
				if d[zty] == "" {
					continue
				}
				for _, u := range strings.Split(d[zty], ",") {
					OwnerId = append(OwnerId, u)
				}
			}
		}
	}
	if OwnerId == nil {
		zty = "non_standardized"
	}

	if schema["TABLE_ROWS"] != nil {
		tableRows = schema["TABLE_ROWS"].(int64)
	}
	if schema["TABLE_TYPE"].(string) != "BASE TABLE" {
		/*先加锁*/
		ztylists = append(ztylists, zty)
		csvDatas[zty] = append(csvDatas[zty], []string{
			schema["TABLE_SCHEMA"].(string),  /*库名*/
			schema["TABLE_NAME"].(string),    /*表名*/
			zty,                              /*主题域*/
			schema["TABLE_TYPE"].(string),    /*表类型*/
			"",                               /*建表模型*/
			schema["TABLE_COMMENT"].(string), /*表注释*/
			schema["CREATE_TIME"].(time.Time).Format("2006-01-02 15:04:05"), /*创建时间*/
			fmt.Sprintf("%d", tableRows),                                    /*数据量*/
			"0",                                                             /*总容量/MB*/
			"0",                                                             /*总容量/单位*/
			strings.Join(OwnerId, ","),                                      /*责任人*/
			"0",                                                             /*总分区数*/
			"0",                                                             /*实际分区数(有数据)*/
			"0",                                                             /*空存储分区数(无数据)*/
			"0",                                                             /*空存储分区范围*/
			"0",                                                             /*总副本数*/
			"0",                                                             /*副本数*/
			"0",                                                             /*分桶数*/
			"0",                                                             /*总副本数平均差*/
			"0",                                                             /*总副本数标准差*/
			"",                                                              /*备注*/
		})
		/*释放锁*/
		return
	}
	/*建表模型*/
	/**************************** 外表内表 **********************************/
	var olap map[string]interface{}
	sql := fmt.Sprintf("show create table %s.%s", schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string))
	r = d.Db.Raw(sql).Scan(&olap)
	if r.Error != nil {
		util.Logger.Error(fmt.Sprintf("%v %v", r.Error.Error(), sql))
		return
	}
	var model string
	crl := strings.Split(olap["Create Table"].(string), "\n")
	for _, s := range crl {
		if strings.Contains(s, " KEY(") {
			model = strings.Split(s, "(")[0]
		}
	}

	/*容量*/
	var data []map[string]interface{}
	x := d.Db.Raw(fmt.Sprintf("show data from %s.%s", schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string))).Scan(&data)
	if x.Error != nil {
		util.Logger.Error(fmt.Sprintf("%s.%s %s -> %s", schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string), schema["TABLE_TYPE"].(string), x.Error.Error()))
		return
	}
	/*datasize unit -- 这里的数值以MB为标准*/
	size := strings.Split(data[0]["Size"].(string), " ")[0]
	unit := strings.Split(data[0]["Size"].(string), " ")[1]
	float, _ := strconv.ParseFloat(size, 64)
	switch unit {
	case "B":
		float = float / 1024 / 1024
	case "KB":
		float = float / 1024
	case "MB":
	case "GB":
		float = float * 1024
	case "TB":
		float = float * 1024 * 1024
	case "PB":
		float = float * 1024 * 1024 * 1024
	default:
		float = 0
	}
	/*分区*/
	var partitions []map[string]interface{}
	r = d.Db.Raw(fmt.Sprintf("show partitions from %s.%s order by LastConsistencyCheckTime,DataSize desc", schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string))).Scan(&partitions)
	if r.Error != nil {
		util.Logger.Error(r.Error.Error())
		return
	}
	if partitions == nil {
		/*先加锁*/
		ztylists = append(ztylists, zty)
		csvDatas[zty] = append(csvDatas[zty], []string{
			schema["TABLE_SCHEMA"].(string),  /*库名*/
			schema["TABLE_NAME"].(string),    /*表名*/
			zty,                              /*主题域*/
			schema["TABLE_TYPE"].(string),    /*表类型*/
			"",                               /*建表模型*/
			schema["TABLE_COMMENT"].(string), /*表注释*/
			schema["CREATE_TIME"].(time.Time).Format("2006-01-02 15:04:05"), /*创建时间*/
			fmt.Sprintf("%d", tableRows),                                    /*数据量*/
			"0",                                                             /*总容量/MB*/
			"0",                                                             /*总容量/单位*/
			strings.Join(OwnerId, ","),                                      /*责任人*/
			"0",                                                             /*总分区数*/
			"0",                                                             /*实际分区数(有数据)*/
			"0",                                                             /*空存储分区数(无数据)*/
			"0",                                                             /*空存储分区范围*/
			"0",                                                             /*总副本数*/
			"0",                                                             /*副本数*/
			"0",                                                             /*分桶数*/
			"0",                                                             /*总副本数平均差*/
			"0",                                                             /*总副本数标准差*/
			"",                                                              /*备注*/
		})
		/*释放锁*/
		return
	}
	/*判断副本数*/
	var replicationArr []int
	for _, partition := range partitions {
		replicationNum, _ := strconv.Atoi(partition["ReplicationNum"].(string))
		replicationArr = append(replicationArr, replicationNum)
	}
	/*判断分桶数*/
	var Max []float64
	for _, m := range partitions {
		Max = append(Max, tools.Size(m["DataSize"].(string)))
	}

	maxDataSize := Max[0]
	for i := 0; i < len(Max); i++ {
		if Max[i] >= maxDataSize {
			maxDataSize = Max[i]
		}
	}
	var conservative, best int
	var maxds string
	Buckets := partitions[0]["Buckets"].(string)
	parseInt, _ := strconv.ParseInt(Buckets, 10, 64)

	if maxDataSize < 1073741824 {
		maxds = "不足1GB"
		conservative = 1
		best = 1
	} else {
		maxds = fmt.Sprintf("%dGB", int(maxDataSize/1073741824))
		conservative = int(maxDataSize / 1073741824)
		best = conservative + 2
	}

	/*tablet 平均值 -- 公式：[tablet的容量总和 / tablet个数 / 1024 / 1024]*/
	var tablet []map[string]interface{}
	r = d.Db.Raw(fmt.Sprintf("show tablet from %s.%s", schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string))).Scan(&tablet)
	if r.Error != nil {
		util.Logger.Error(r.Error.Error())
		return
	}
	var datasize float64
	var i2 statistics.Float64

	for _, m2 := range tablet {
		var ds float64
		if d.Version >= 2.5 {
			ds = sizeJuetByMb(m2["DataSize"].(string))
		} else {
			ds, _ = strconv.ParseFloat(m2["DataSize"].(string), 64)
		}
		datasize += ds
		i2 = append(i2, ds)
	}

	var mean, sqrt float64
	if d.Version >= 2.5 {
		mean = datasize / float64(len(tablet))
		/*tablet 标准差  --  公式：[平方根 (tablet数组方差) / 1024 / 1024]*/
		sqrt = math.Sqrt(statistics.Variance(&i2))
	} else {
		mean = datasize / float64(len(tablet)) / 1024 / 1024
		/*tablet 标准差  --  公式：[平方根 (tablet数组方差) / 1024 / 1024]*/
		sqrt = math.Sqrt(statistics.Variance(&i2)) / 1024 / 1024
	}

	/*生成备注*/
	/*生成备注*/
	var comment []string
	if tools.FindMin(replicationArr) < 3 && schema["TABLE_NAME"].(string) != "sams_cos__sams_user_action_d" {
		if tools.FindMin(replicationArr) == 1 {
			comment = append(comment, fmt.Sprintf(`💬副本异常，原副本数%d，当这个副本发生损坏时数据将丢失也无法修复，请按照标准"replication_num" = "3"调整，调整方法：创建新表,新表属性PROPERTIES ("replication_num" = "3") -> 旧表数据写入新表 -> 对比新表旧表数据量 -> 旧表重命名 -> 新表重命名 -> 删除旧表。`, tools.FindMin(replicationArr)))
		} else {
			comment = append(comment, fmt.Sprintf(`💬副本异常，原副本数%d，当某个副本发生损坏时将无法更好的自动修复，请按照标准"replication_num" = "3"调整，调整方法：创建新表,新表属性PROPERTIES ("replication_num" = "3") -> 旧表数据写入新表 -> 对比新表旧表数据量 -> 旧表重命名 -> 新表重命名 -> 删除旧表。`, tools.FindMin(replicationArr)))
		}
	}
	//util.Logger.Info(fmt.Sprintf("%s.%s 分桶：%d，大小：%f，保底：%d，建议：%d", schema["TABLE_SCHEMA"].(string), schema["TABLE_NAME"].(string), parseInt, maxDataSize, conservative, best))
	if parseInt > int64(conservative)+10 {
		comment = append(comment, fmt.Sprintf("💬分桶数过大，表总大小%s，其中最大的分区%s，按照分区1GB=1BUCKETS原则，原BUCKETS %d不足以更好均衡并发，调整方法：创建新表,新表属性BUCKETS %d -> 旧表数据写入新表 -> 对比新表旧表数据量 -> 旧表重命名 -> 新表重命名 -> 删除旧表。", data[0]["Size"].(string), maxds, parseInt, best))
	} else if parseInt < int64(conservative) {
		comment = append(comment, fmt.Sprintf("💬分桶数过小，表总大小%s，其中最大的分区%s，按照分区1GB=1BUCKETS原则，原BUCKETS %d将会导致底层数据分片均衡倾斜，查询性能降低，调整方法：创建新表,新表属性BUCKETS %d -> 旧表数据写入新表 -> 对比新表旧表数据量 -> 旧表重命名 -> 新表重命名 -> 删除旧表。", data[0]["Size"].(string), maxds, parseInt, best))
	}
	/*判断分区总数，空分区*/
	nilpt, stopt, nilptname := tools.PartitionsJuet(partitions, schema["TABLE_NAME"].(string))
	if nilpt >= 365 {
		comment = append(comment, fmt.Sprintf(`💬空存储分区过多，%d个分区，其中%d个实际分区(有数据的)，%d个空存储分区(无数据)，调整方法(二选一)：
1.查看旧表空闲分区 -> 创建新表,不需要带上空闲分区 -> 旧表数据写入新表 -> 对比新表旧表数据量 -> 旧表重命名 -> 新表重命名 -> 删除旧表。
2.查看旧表空闲分区 -> 删除旧表空闲分区。`, nilpt+stopt, stopt, nilpt))
	}
	if comment != nil {
		comment = append(comment, `

-- 无论删除表或分区,请带上force,跳过回收站.
【-建议-】
（建议）1.drop table force
（建议）2.alter table drop partition force

-- 删表重建过程中，切勿先rename正式表再创建
【-慎用-】
（慎用）1.alter table a rename b
（慎用）2.create table a
（慎用）3.insert a select  b
（慎用）4.drop b

【-建议-】
（建议）1.create table c
（建议）2.insert c select a
（建议）3.alter table a rename b
（建议）4.alter table c rename a
（建议）5.drop table b
`)
	}

	/*设置csv数组*/
	//"库名", "表名", "总容量/MB", "总容量/单位", "总分区数", "实际分区数(有数据)", "空存储分区数(无数据)", "空存储分区范围", "总副本数", "副本数", "分桶数", "总副本数平均差", "总副本数标准差", "备注"
	/*先加锁*/
	ztylists = append(ztylists, zty)
	csvDatas[zty] = append(csvDatas[zty], []string{
		schema["TABLE_SCHEMA"].(string),  /*库名*/
		schema["TABLE_NAME"].(string),    /*表名*/
		zty,                              /*主题域*/
		schema["TABLE_TYPE"].(string),    /*表类型*/
		model,                            /*建表模型*/
		schema["TABLE_COMMENT"].(string), /*表注释*/
		schema["CREATE_TIME"].(time.Time).Format("2006-01-02 15:04:05"), /*创建时间*/
		fmt.Sprintf("%d", tableRows),                                    /*数据量*/
		fmt.Sprintf("%0.2f", float),                                     /*总容量/MB*/
		data[0]["Size"].(string),                                        /*总容量/单位*/
		strings.Join(OwnerId, ","),                                      /*责任人*/
		fmt.Sprintf("%d", nilpt+stopt),                                  /*总分区数*/
		fmt.Sprintf("%d", stopt),                                        /*实际分区数(有数据)*/
		fmt.Sprintf("%d", nilpt),                                        /*空存储分区数(无数据)*/
		nilptname,                                                       /*空存储分区范围*/
		data[0]["ReplicaCount"].(string),                                /*总副本数*/
		fmt.Sprintf("%d", findMin(replicationArr)),                      /*副本数*/
		fmt.Sprintf("%d", parseInt),                                     /*分桶数*/
		fmt.Sprintf("%0.2f", mean),                                      /*总副本数平均差*/
		fmt.Sprintf("%0.2f", sqrt),                                      /*总副本数标准差*/
		"",                                                              /*近一个月访问次数*/
		"",                                                              /*近三个月访问次数*/
		strings.Join(comment, ",\n"),                                    /*备注*/
	})
}
