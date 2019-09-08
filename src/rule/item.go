package rule

import (
	. "db"
)

//规则项
type Item struct {
	ID          string	//规则项 id
	Operator    string	//规则因子运算符
	FactorType  string	//规则因子类型
	FactorValue string	//规则因子值
}

//从数据库获取规则项
func listItem() (map[string]Item, error) {
	m := map[string]Item{}
	items, err := ItemObj.ListItem()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		m[item.ID] = Item{
			ID:          item.ID,
			Operator:    item.Operator,
			FactorType:  item.FactorType,
			FactorValue: item.FactorValue,
		}
	}
	return m, nil
}