/*
 *@author  chengkenli
 *@project healthyreport
 *@package run
 *@file    SchemaOneStop
 *@date    2024/7/5 17:42
 */

package run

import (
    "bytes"
    "encoding/json"
    "fmt"
    "healthyreport/util"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

func SchemaOneStopZty() {
    type fileter struct {
        ControlId  string `json:"controlId"`
        DataType   int    `json:"dataType"`
        SpliceType int    `json:"spliceType"`
        FilterType int    `json:"filterType"`
        Value      string `json:"value"`
    }
    type test1 struct {
        AppKey      string    `json:"appKey"`
        Sign        string    `json:"sign"`
        WorksheetId string    `json:"worksheetId"`
        ViewId      string    `json:"viewId"`
        PageSize    int       `json:"pageSize"`
        PageIndex   int       `json:"pageIndex"`
        SortId      string    `json:"sortId"`
        IsAsc       bool      `json:"isAsc"`
        Filters     []fileter `json:"filters"`
        //Filters     struct{}
        NotGetTotal      bool `json:"notGetTotal"`
        UseControlId     bool `json:"useControlId"`
        GetSystemControl bool `json:"getSystemControl"`
    }
    type data struct {
        Data struct {
            Rows []struct {
                ID    string `json:"_id"`
                Rowid string `json:"rowid"`
                Ctime string `json:"ctime"`
                Caid  struct {
                    AccountID string `json:"accountId"`
                    Fullname  string `json:"fullname"`
                    Avatar    string `json:"avatar"`
                    IsPortal  bool   `json:"isPortal"`
                    Status    int    `json:"status"`
                } `json:"caid"`
                Uaid struct {
                    AccountID string `json:"accountId"`
                    Fullname  string `json:"fullname"`
                    Avatar    string `json:"avatar"`
                    IsPortal  bool   `json:"isPortal"`
                    Status    int    `json:"status"`
                } `json:"uaid"`
                Ownerid struct {
                    AccountID string `json:"accountId"`
                    Fullname  string `json:"fullname"`
                    Avatar    string `json:"avatar"`
                    IsPortal  bool   `json:"isPortal"`
                    Status    int    `json:"status"`
                } `json:"ownerid"`
                Utime   string `json:"utime"`
                Ztmc    string `json:"ztmc"`
                Ztowner []struct {
                    AccountID string `json:"accountId"`
                    Fullname  string `json:"fullname"`
                    Avatar    string `json:"avatar"`
                    Status    int    `json:"status"`
                } `json:"ztowner"`
                Ztymc              string `json:"ztymc"`
                Ownerid1           string `json:"ownerid1"`
                Autoid             int    `json:"autoid"`
                Ztjc               string `json:"ztjc"`
                Allowdelete        bool   `json:"allowdelete"`
                Controlpermissions string `json:"controlpermissions"`
            } `json:"rows"`
            Total int `json:"total"`
        } `json:"data"`
        Success   bool `json:"success"`
        ErrorCode int  `json:"error_code"`
    }

    test := test1{
        AppKey:           "1a358123246c0ab6",
        Sign:             "YjRiZTMxOTg1ZTE3ZjgxNDVjNGE0NWUxODViYjg4N2EzZWYwODJjYWNmMWExZDkwMWI5MmY1ZTBkOTkwNDIxYg==",
        WorksheetId:      "64a515452d3d4fc5eb0c1514",
        ViewId:           "",
        PageSize:         1000,
        PageIndex:        1,
        SortId:           "",
        IsAsc:            false,
        Filters:          nil,
        NotGetTotal:      false,
        UseControlId:     false,
        GetSystemControl: false,
    }

    marshal, errM := json.Marshal(test)
    if errM != nil {
        fmt.Println(errM)
    }
    request1, err1 := http.NewRequest("POST", "https://xxxxxxxxx", bytes.NewBuffer(marshal))
    if err1 != nil {
        log.Fatal("1", err1)
    }
    request1.Header.Set("Content-Type", "application/json")
    do1, _ := (&http.Client{}).Do(request1)
    defer do1.Body.Close()
    all1, _ := ioutil.ReadAll(do1.Body)

    var d data
    replaceAll := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(all1), "\\", ""), "\"[", "["), "]\"", "]")
    json.Unmarshal([]byte(replaceAll), &d)
    for _, r := range d.Data.Rows {
        util.Domain = append(util.Domain, map[string]string{r.Ztmc: r.Ownerid1})
    }
}
