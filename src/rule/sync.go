package rule

import (
	"log"
	"sync"
	"time"
)

var once sync.Once

//默认刷新时间
const refreshTime = 10 * time.Minute

//定时更新规则数据
func syncRule() {
	timer := time.NewTicker(refreshTime)
	for range timer.C {
		if err := loadRuleFromDB(); err != nil {
			log.Printf("从数据库同步规则信息失败: ", err)
		} else {
			log.Printf("从数据库同步规则信息成功...")
		}
	}
}