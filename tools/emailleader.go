/*
 *@author  chengkenli
 *@project SchemaHealthy
 *@package tools
 *@file    emailleader
 *@date    2024/7/26 13:37
 */

package tools

import (
    "encoding/json"
    "fmt"
    "healthyreport/util"
    "strings"
)

func SchemaDomainGroup(leader string) []string {
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
        email = append(email, s+"@xxxxxxxxx")
    }
    return email
}
