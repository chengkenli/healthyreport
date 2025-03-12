/*
 *@author  chengkenli
 *@project SchemaHealthy
 *@package tools
 *@file    object
 *@date    2024/7/24 17:11
 */

package tools

import (
    "encoding/json"
    "fmt"
    "gorm.io/gorm"
    "healthyreport/util"
    "io"
    "io/ioutil"
    "net/http"
    "regexp"
    "sort"
    "strconv"
    "strings"
    "time"
)

func Chversion(db *gorm.DB) float64 {
    var versions map[string]interface{}
    r := db.Raw("select current_version() as version").Scan(&versions)
    if r.Error != nil {
        util.Logger.Error(r.Error.Error())
        return 0
    }
    version, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", strings.Split(strings.Split(versions["version"].(string), " ")[0], ".")[0], strings.Split(strings.Split(versions["version"].(string), " ")[0], ".")[1]), 64)
    return version
}

/*返回mb大小*/
func SizeJuetByMb(s string) float64 {
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

func FindMin(arrs []int) int {
    if len(arrs) == 0 {
        return -1 // 如果数组为空则返回-1表示没有最小值
    } else if len(arrs) == 1 {
        return arrs[0] // 如果只有一个元素，直接返回该元素
    }
    // 使用sort包对切片进行排序
    sort.Ints(arrs)
    return arrs[0] // 返回排序后的第一个元素作为最小值
}

func PartitionsJuet(partitions []map[string]interface{}, tablename string) (int, int, string) {
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
        if len(NulsName) > 100 {
            str = fmt.Sprintf("%s~%s", NulsName[0], NulsName[len(NulsName)-5])
        } else {
            str = ""
        }
    }
    return len(NulsName), len(NosName), str
}

func DomainGroup(leader string) []string {
    body := fmt.Sprintf(`{
  "appKey": "xxxxxxxxx",
  "sign": "xxxxxxxxx==",
  "callbackURL": "",
  "OwnerID": "%s"
}`, leader)
    util.Logger.Info(body)
    r := Post("POST", "https://xxxxxxxxx", strings.NewReader(body))

    type group struct {
        ID string `json:"id"`
    }
    var g group
    err := json.Unmarshal(r, &g)
    if err != nil {
        util.Logger.Error(err.Error())
        return nil
    }
    if len(g.ID) == 0 {
        return nil
    }
    var email []string
    for _, s := range strings.Split(g.ID, ",") {
        email = append(email, s+"@xxxxxxxxx.com")
    }
    return email
}

func Post(method, url string, body io.Reader) []byte {
    request, err := http.NewRequest(method, url, body)
    if err != nil {
        util.Logger.Error(err.Error())
        return nil
    }
    request.Header.Set("Content-Type", "application/json;charset=utf-8")
    client := &http.Client{Timeout: time.Second * 30}
    respone, err := client.Do(request)
    if err != nil {
        util.Logger.Error(err.Error())
        return nil
    }
    defer respone.Body.Close()
    b, err := ioutil.ReadAll(respone.Body)
    if err != nil {
        util.Logger.Error(err.Error())
        return nil
    }
    return b
}

// RemoveDuplicateStrings /*数组去重*/
func RemoveDuplicateStrings(strs []string) []string {
    result := []string{}
    tempMap := map[string]byte{} // 存放不重复字符串
    for _, e := range strs {
        l := len(tempMap)
        tempMap[e] = 0
        if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
            result = append(result, e)
        }
    }
    return result
}

// Size 返回bytes
func Size(s string) float64 {
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

// StringInSlice 检查数组中是否存在某个元素
func StringInSlice(str string, list []string) bool {
    for _, v := range list {
        if v == str {
            return true
        }
    }
    return false
}

//// PrintProgress 用于在一行内打印进度条
//func PrintProgress(current, total int) {
//	// 计算进度百分比
//	percent := int(float64(current) / float64(total) * 100)
//	// 创建进度条
//	bar := strings.Repeat("-", percent) + strings.Repeat(" ", 100-percent)
//	// 使用ANSI转义序列将光标移动到行首
//	fmt.Printf("\033[2K\r%d%% [%s](%d/%d)", percent, bar, current, total)
//}

// PrintProgress 用于在一行内打印进度条
func PrintProgress(current, total int, comment ...string) {
    percent := current * 100 / total                                                  // 计算进度百分比
    length := 20                                                                      // 设定进度条的长度
    s := strings.Repeat("■", current*length/total)                                    // 根据完成度填充进度条
    e := strings.Repeat("□", length-len(s)/3)                                         // 填充剩余的空格(这里一个□占用了3个字符，所以除以3)
    bar := s + e                                                                      // 拼接
    fmt.Printf("\033[2K\r%d%% [%s](%d/%d) %v", percent, bar, current, total, comment) // 使用ANSI转义序列将光标移动到行首
}
