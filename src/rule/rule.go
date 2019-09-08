package rule

import (
	. "db"
	"log"
)

var (
	//默认规则引擎
	DefaultRuleEngine = &defaultRuleEngine
	defaultRuleEngine RuleEngine
)



type Factor struct {
	FactorType  string
	FactorValue string
}

//规则引擎,包含两部分,常驻内存
//1.所有的规则组成的 map
//2.所有规则项组成的 map
type RuleEngine struct {
	Rules map[string]Rule
	Items map[string]Item
}

//加载规则
func (r RuleEngine) SetRules() error {
	ruleMap := map[string]Rule{}

	ruleList, err := RuleObj.ListRule()
	if err != nil {
		return err
	}

	for _, rule := range ruleList {

		//决策树
		tree := &DecisionTree{
			Exp:      rule.Tree,
			Priority: rule.Priority,
			Value:    rule.Value,
		}
		m, exist := ruleMap[rule.Name]

		// 不存在规则， 新增一条规则，并初始化决策树数组
		if !exist {
			m = Rule{
				Name:         rule.Name,
				DefaultValue: rule.DefautValue,
				TreeList:     make([]*DecisionTree, 0),
			}
		}
		m.TreeList = append(m.TreeList, tree)
		ruleMap[rule.Name] = m
	}

	DefaultRuleEngine.Rules = ruleMap
	return nil
}

//加载规则项
func (r RuleEngine) SetItems() error {
	items, err := listItem()

	if err != nil {
		return err
	}
	DefaultRuleEngine.Items = items

	log.Printf("加载规则项成功，一共 %d 条", len(DefaultRuleEngine.Items))

	return nil
}


//根据规则名和规则因子匹配
func (r RuleEngine) Match(name string, f ...Factor) (string, bool) {

	//规则不存在， 直接返回不匹配
	m, exist := r.Rules[name];
	if !exist {
		return "", false
	}

	return m.match(f...)
}


type Rule struct {
	Name         string          //规则名
	DefaultValue string          //未匹配时默认值
	TreeList     []*DecisionTree //决策树
}

func (r Rule) match(f ...Factor) (string, bool) {

	if len(f) == 0 {
		//直接返回规则默认值
		return r.DefaultValue, false
	}

	//根据决策树优先级进行比较(优先级已经在 sql 中排序完毕)
	for _, tree := range r.TreeList {
		value, matched := tree.match(f...)
		if matched {
			return value, true
		}
	}

	//不匹配返回默认值
	return r.DefaultValue, false
}




func loadRuleFromDB() error {
	//加载规则项
	err := DefaultRuleEngine.SetItems()
	if err != nil {
		return err
	}

	//加载决策树信息，按优先级排序，并按规则分类
	err = DefaultRuleEngine.SetRules()
	if err != nil {
		return err
	}
	return nil
}


//初始化规则
func InitRule() error {

	//从 db 加载到内存
	err := loadRuleFromDB()
	if err != nil {
		return err
	}

	//用于自动刷新规则到内存
	//once.Do(func() {
	//	go syncRule()
	//})

	return nil
}

//规则匹配
func MatchRule(ruleName string, factor ...Factor) (string, bool) {
	return DefaultRuleEngine.Match(ruleName, factor...)
}