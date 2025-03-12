package run

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"healthyreport/util"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	NulsNum = 100
)

var DataBaseOwner map[string]string

/*初始化匹配库级别owner*/
func init() {

	//if dbname == "information_schema" || dbname == "starrocks_monitor" || dbname == "_statistics_" {

	switch util.P.App {
	case "adhoc":
		DataBaseOwner = map[string]string{
			"demo":                   "StarRocksAdmin(Chengken li)",
			"information_schema":     "StarRocks Metadata",
			"starrocks_monitor":      "StarRocks Metadata",
			"_statistics_":           "StarRocks Metadata",
			"ops":                    "StarRocksAdmin(Alden Dong)",
			"audit":                  "StarRocksAdmin(Alden Dong)",
			"cn_po_home_system":      "Zhiwen Jiang",
			"ods_dev":                "ives li",
			"ods_gray":               "ives li",
			"ods_migration_td_gray":  "ives li",
			"scm_network_secure":     "Yuheng Chen",
			"adhoc":                  "StarRocksAdmin(Alden Dong)",
			"ads":                    "Eric he",
			"ads_dev":                "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"ads_dev_secure":         "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"ads_rt":                 "Eric he",
			"ads_secure":             "Eric he",
			"ap_secure":              "Eason Xie",
			"bi_realty_secure":       "Huiyin Yang",
			"bi_sams_secure":         "Victor zhang",
			"bi_sc_secure":           "Eric he",
			"cloud_fcst_dm":          "Xi Liu",
			"cn_backup_secure":       "king zhang",
			"cn_ec_bi_secure":        "Shane shen",
			"cn_ec_wmdj_user_action": "Eric he",
			"cn_mdse_dm_dl_tables":   "Leap Leng",
			"cn_pricing_dl_tables":   "IVAN XIE",
			"cn_sams_dl_secure":      "Ender Wu",
			"cn_wid_dl_secure":       "Xi Liu",
			"dim":                    "Eric he",
			"dim_dev":                "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"dim_dev_secure":         "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"dim_secure":             "Eric he",
			"dm":                     "Eric he",
			"dm_dev":                 "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"dm_dev_secure":          "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"dm_secure":              "Eric he",
			"dw":                     "Eric he",
			"dw_dev":                 "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"dw_dev_secure":          "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"dw_rt":                  "Eric he",
			"dw_secure":              "Eric he",
			"fin_sox":                "king zhang",
			"fin_sox_dev":            "king zhang",
			"flash_report":           "king zhang",
			"flash_report_dev":       "king zhang",
			"flash_report_sit":       "king zhang",
			"hyper_ec_secure":        "Joyce Jiang",
			"hyper_mdse_dm_secure":   "Leap Leng",
			"mbrship_secure":         "Josh Xu",
			"mcfc_report_dev":        "Eric he",
			"o2o_datacubes_secure":   "Eric he",
			"ods":                    "ives li",
			"ods_archive":            "Eric he",
			"ods_rt":                 "ives li",
			"ods_rt_secure":          "ives li",
			"ods_secure":             "ives li",
			"ods_sox":                "ives li",
			"ods_test":               "Sarah xie",
			"pro_dgtmkt_data":        "Jeffrey jiang",
			"pro_scct_dev":           "Eric he",
			"sams_finance":           "Echo Cao",
			"scct_logis_dev":         "Eric he",
			"scm_dcqe_secure":        "Jiayou Pan",
			"scm_secure":             "Lynn Cao",
			"svcdordgtmkt":           "Bo Ma",
			"test":                   "chengken li",
		}
	case "app":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"ads":                "Eric he",
			"ads_rt":             "Eric he",
			"ads_secure":         "Eric he",
			"audit":              "StarRocksAdmin(Alden Dong)",
			"cn_sams_dl_secure":  "Ender Wu",
			"dim":                "Eric he",
			"dim_secure":         "Eric he",
			"dm":                 "Eric he",
			"dm_secure":          "Eric he",
			"dw":                 "Eric he",
			"dw_rt":              "Eric he",
			"dw_secure":          "Eric he",
			"fin_sox":            "king zhang",
			"ods":                "ives li",
			"ods_rt":             "Eric he",
			"ods_rt_secure":      "ives li",
			"ods_secure":         "ives li",
		}
	case "scct":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"audit":              "StarRocksAdmin(Alden Dong)",
			"cn_po_home_system":  "Eric he",
			"ads":                "Eric he",
			"dim":                "Eric he",
			"dwd":                "Eric he",
			"dws":                "Eric he",
			"mcfc_report":        "Eric he",
			"ods":                "ives li",
			"ods_secure":         "ives li",
			"scct_logis":         "Eric he",
			"svccn_logis":        "Eric he",
			"svccn_logis_query":  "Eric he",
		}
	case "sccts":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"audit":              "StarRocksAdmin(Alden Dong)",
			"cn_po_home_system":  "Eric he",
			"ads":                "Eric he",
			"dim":                "Eric he",
			"dwd":                "Eric he",
			"dws":                "Eric he",
			"mcfc_report":        "Eric he",
			"ods":                "ives li",
			"ods_secure":         "ives li",
			"scct_logis":         "Eric he",
			"svccn_logis":        "Eric he",
			"svccn_logis_query":  "Eric he",
		}
	case "cdp", "sr-cdp":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"cdp":                "Bryan FAN",
			"cdp_test":           "Bryan FAN",
			"ods_secure":         "ives li",
			"audit":              "StarRocksAdmin(Alden Dong)",
		}
	case "api", "sr-api":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"cdp":                "Bryan FAN",
			"cdp_test":           "Bryan FAN",
			"ma_test":            "Bryan FAN",
			"cdp_api":            "Bryan FAN",
			"ods_secure":         "ives li",
			"audit":              "StarRocksAdmin(Alden Dong)",
		}
	case "ma", "sr-ma":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"cdp":                "Bryan FAN",
			"ma_test":            "Bryan FAN",
			"audit":              "StarRocksAdmin(Alden Dong)",
		}
	case "sr-adhoc":
		DataBaseOwner = map[string]string{
			"demo":                     "StarRocksAdmin(Chengken li)",
			"information_schema":       "StarRocks Metadata",
			"starrocks_monitor":        "StarRocks Metadata",
			"_statistics_":             "StarRocks Metadata",
			"ops":                      "StarRocksAdmin(Alden Dong)",
			"adhoc":                    "StarRocksAdmin(Alden Dong)",
			"ads_dev_secure":           "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"ads_rt_dev":               "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"ads_rt_dev_secure":        "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"ads_secure":               "Eric he",
			"algo":                     "Xi Liu",
			"algo_dev":                 "Xi Liu",
			"ap_secure":                "Eason Xie",
			"audit":                    "StarRocksAdmin(Alden Dong)",
			"bi_realty":                "Huiyin Yang",
			"cn_core_dim_vm":           "James HUANG",
			"cn_po_home_system_dev":    "marco xi",
			"cn_wc_highsecure":         "James HUANG",
			"cn_wc_mb_secure":          "James HUANG",
			"cn_wc_mb_vm":              "James HUANG",
			"cn_wc_repl_vm":            "James HUANG",
			"cn_wc_vm":                 "James HUANG",
			"cn_wm_mb_secure":          "James HUANG",
			"cn_wm_mb_vm":              "James HUANG",
			"cn_wm_repl_vm":            "James HUANG",
			"cn_wm_vm":                 "James HUANG",
			"data_test":                "Sarah xie",
			"dim_dev_secure":           "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dim_rt_dev":               "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dim_rt_dev_secure":        "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dim_secure":               "Eric he",
			"dm_dev_secure":            "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dm_secure":                "Eric he",
			"dwd_dev_secure":           "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd_rt_dev":               "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd_rt_dev_secure":        "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd_secure":               "Eric he",
			"dws_dev_secure":           "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dws_rt_dev":               "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dws_rt_dev_secure":        "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dws_secure":               "Eric he",
			"euclid_scn_forecast_prod": "James HUANG",
			"fin_sox":                  "king zhang",
			"fin_sox_dev":              "king zhang",
			"finance_kettle":           "王玲",
			"flash_report_dev":         "king zhang",
			"hyper_bi_secure":          "Shane shen",
			"mcfc_report_dev":          "marco xi",
			"ods_dev":                  "ives li",
			"ods_dev_secure":           "ives li",
			"ods_rt_dev":               "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"ods_rt_dev_secure":        "Eric he,flank zhang,king zhang,Jeffrey jiang",
			"ods_secure":               "ives li",
			"ods_sox":                  "king zhang",
			"ods_sox_dev":              "king zhang",
			"ods_sox_test":             "Dylan Pang",
			"ods_test":                 "Dylan Pang",
			"ods_test_secure":          "Dylan Pang",
			"scct_logis_dev":           "Eric he",
			"wm_ad_hoc":                "James HUANG",
			"wm_common_vm":             "James HUANG",
			"ads":                      "nick zhu",
			"ads_dev":                  "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"ads_rt":                   "nick zhu",
			"dim":                      "nick zhu",
			"dim_dev":                  "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dm":                       "nick zhu",
			"dm_dev":                   "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd":                      "nick zhu",
			"dwd_dev":                  "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd_rt":                   "nick zhu",
			"dws":                      "nick zhu",
			"dws_dev":                  "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dws_rt":                   "nick zhu",
			"flash_report":             "king zhang",
			"ods":                      "ives li",
			"ods_rt":                   "ives li",
			"ods_rt_secure":            "ives li",
		}
	case "sr-app":
		DataBaseOwner = map[string]string{
			"demo":               "StarRocksAdmin(Chengken li)",
			"information_schema": "StarRocks Metadata",
			"starrocks_monitor":  "StarRocks Metadata",
			"_statistics_":       "StarRocks Metadata",
			"ops":                "StarRocksAdmin(Alden Dong)",
			"ads_rt_secure":      "James HUANG",
			"ads_secure":         "James HUANG",
			"audit":              "StarRocksAdmin(Alden Dong)",
			"dim_rt":             "James HUANG",
			"dim_rt_secure":      "James HUANG",
			"dim_secure":         "James HUANG",
			"dm":                 "Eric he",
			"dm_secure":          "Eric he",
			"dwd_rt_secure":      "James HUANG",
			"dwd_secure":         "James HUANG",
			"dws_rt_secure":      "James HUANG",
			"dws_secure":         "James HUANG",
			"ods_secure":         "ives li",
			"ods_sox":            "ives li",
			"ads":                "Eric he",
			"ads_dev":            "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"ads_rt":             "Eric he",
			"dim":                "Eric he",
			"dim_dev":            "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd":                "Eric he",
			"dwd_dev":            "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dwd_rt":             "Eric he",
			"dws":                "Eric he",
			"dws_dev":            "Eric he,king zhang,flank zhang,Jeffrey jiang",
			"dws_rt":             "Eric he",
			"fin_sox":            "king zhang",
			"flash_report":       "king zhang",
			"flash_report_dev":   "king zhang",
			"flash_report_sit":   "king zhang",
			"ods":                "ives li",
			"ods_dev":            "ives li",
			"ods_rt":             "Eric he",
			"ods_rt_secure":      "Eric he",
			"ods_test":           "ives li",
		}
	default:
	}
}

func findMin(arrs []int) int {
	if len(arrs) == 0 {
		return -1 // 如果数组为空则返回-1表示没有最小值
	} else if len(arrs) == 1 {
		return arrs[0] // 如果只有一个元素，直接返回该元素
	}
	// 使用sort包对切片进行排序
	sort.Ints(arrs)
	return arrs[0] // 返回排序后的第一个元素作为最小值
}

/*返回bytes大小*/
func sizeJuet(s string) float64 {
	s = strings.ToLower(s)
	//正则
	re := regexp.MustCompile(`\d+\.?\d*`)

	float, _ := strconv.ParseFloat(re.FindString(s), 64)

	if strings.Contains(s, "kb") {
		return float * 1024
	}
	if strings.Contains(s, "mb") {
		return float * 1024 * 1024
	}
	if strings.Contains(s, "gb") {
		return float * 1024 * 1024 * 1024
	}
	if strings.Contains(s, "tb") {
		return float * 1024 * 1024 * 1024 * 1024
	}
	return 0
}

/*返回mb大小*/
func sizeJuetByMb(s string) float64 {
	s = strings.ToLower(s)
	//正则
	re := regexp.MustCompile(`\d+\.?\d*`)

	float, _ := strconv.ParseFloat(re.FindString(s), 64)

	if strings.Contains(s, "kb") {
		return float / 1024
	}
	if strings.Contains(s, "mb") {
		return float
	}
	if strings.Contains(s, "gb") {
		return float * 1024 * 1024
	}
	if strings.Contains(s, "tb") {
		return float * 1024 * 1024 * 1024
	}
	return 0
}

func bucketJuet(partitions []map[string]interface{}) int {
	Bucket := partitions[0]["Buckets"].(string)
	parseInt, _ := strconv.ParseInt(Bucket, 10, 64)

	var Max []float64
	var conservative int
	for _, mp := range partitions {
		Max = append(Max, sizeJuet(mp["DataSize"].(string)))
	}
	maxDataSize := Max[0]
	for i := 0; i < len(Max); i++ {
		if Max[i] >= maxDataSize {
			maxDataSize = Max[i]
		}
	}
	if maxDataSize <= 1073741824 && parseInt <= 3 {
	} else if maxDataSize < 1073741824 {
		conservative = 1
	} else {
		conservative = int(maxDataSize / 1073741824)
	}
	return conservative
}

func partitionsJuet(partitions []map[string]interface{}, tablename string) (int, int, string) {
	var NulsName, NosName []string
	for _, s := range partitions {
		if s["DataSize"].(string) == ".000 " {
			NulsName = append(NulsName, s["PartitionName"].(string))
			continue
		}
		if s["DataSize"].(string) == "0B" {
			NulsName = append(NulsName, s["PartitionName"].(string))
			continue
		}
		NosName = append(NosName, s["PartitionName"].(string))
	}

	var str string
	if NulsName == nil || strings.Contains(strings.Join(NulsName, ","), tablename) {
		str = ""
	} else {
		if len(NulsName) > NulsNum {
			str = fmt.Sprintf("%s~%s", NulsName[0], NulsName[len(NulsName)-5])
		} else {
			str = ""
		}
	}
	return len(NulsName), len(NosName), str
}

func chversion(db *gorm.DB) float64 {
	var versions map[string]interface{}
	r := db.Raw("select current_version() as version").Scan(&versions)
	if r.Error != nil {
		util.Logger.Error(r.Error.Error())
		return 0
	}
	version, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", strings.Split(strings.Split(versions["version"].(string), " ")[0], ".")[0], strings.Split(strings.Split(versions["version"].(string), " ")[0], ".")[1]), 64)
	return version
}

func writefile(fname, msg string) {
	fileHandle, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriterSize(fileHandle, len(msg))

	buf.WriteString(msg)

	err = buf.Flush()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
