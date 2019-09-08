package main

import (
	"db"
	"fmt"
	"rule"
)

func main()  {

	//连接数据库, 规则存在 db 中
	db.InitDB()


	//初始化规则引擎, 将规则信息加载到内存
	rule.InitRule()



	//A > = 90
	//B >= 80
	//C >= 70
	//D >= 60
	//E 60以下

	factors := []rule.Factor{
		{
			FactorType: "score",
			FactorValue: "59",
		},

		//这里还可以有其他规则因子...
	}

	val, matched := rule.MatchRule("score_level", factors...)
	fmt.Println(matched, "分数 59, 等级:", val)
}


