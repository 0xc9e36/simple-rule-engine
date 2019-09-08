package db

var RuleObj rule

type rule struct {}

//数据库中的rule记录
type RuleModel struct {
	Name        string //规则名
	DefautValue string //未匹配时，规则默认值
	Tree        string //规则表达式
	Value       string //表达式为真时，规则值
	Priority    string //规则优先级
}

//读取所有的规则
func (r rule) ListRule() ([]*RuleModel, error) {

	sql := `SELECT IFNULL(a.name,""), IFNULL(a.default_value, ""), IFNULL(b.tree,""), IFNULL(b.value,""), IFNULL(b.priority,1)
			FROM rule a 
			LEFT JOIN rule_decision_tree b 
			ON a.id=b.rule_id
			WHERE a.flag=0 and b.flag=0
			ORDER BY b.priority ASC`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*RuleModel, 0)
	for rows.Next() {
		rule := &RuleModel{}
		if err = rows.Scan(&rule.Name, &rule.DefautValue, &rule.Tree, &rule.Value, &rule.Priority); err != nil {
			return nil, err
		}
		result = append(result, rule)
	}
	return result, nil
}