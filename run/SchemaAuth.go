/*
 *@author  chengkenli
 *@project SchemaHealthy
 *@package run
 *@file    SchemaAuth
 *@date    2024/9/26 10:55
 */

package run

import (
	"gorm.io/gorm"
	"healthyreport/util"
)

func SchemaAuth(db *gorm.DB) int {
	var m []map[string]interface{}
	r := db.Raw("show backends").Scan(&m)
	if r.Error != nil {
		util.Logger.Error(r.Error.Error())
		return -1
	}
	return len(m)
}
