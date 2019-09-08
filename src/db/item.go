package db

var ItemObj item

type item struct{}

type ItemModel struct {
	ID          string //规则id
	Operator    string //运算符
	FactorType  string //规则类型
	FactorValue string //规则取值
}

func (i item) ListItem() ([]*ItemModel, error) {

	sql := `SELECT IFNULL(id, 0), IFNULL(operator, ""), IFNULL(factor_type, ""), IFNULL(factor_val, "")
			FROM rule_item
			WHERE flag=0`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*ItemModel, 0)
	for rows.Next() {
		ruleItem := &ItemModel{}
		if err = rows.Scan(&ruleItem.ID, &ruleItem.Operator, &ruleItem.FactorType, &ruleItem.FactorValue); err != nil {
			return nil, err
		}
		result = append(result, ruleItem)
	}
	return result, nil
}
