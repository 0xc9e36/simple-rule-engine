package rule

import (
	"unicode"
	. "utils"
)

const (
	Space        = " "
	LeftBracket  = "("
	RightBracket = ")"
	BitwiseAnd   = "&"
	BitwiseOr    = "|"
	LogicAnd     = "&&"
	LogicOr      = "||"
	LogicNon     = "!"
	LogicTrue  = "T" //规则值为 true
	LogicFalse = "F" //规则值为 false
)

//规则决策树
type DecisionTree struct {
	Exp      string //表达式
	Priority string //优先级
	Value    string //值
}

/**

	支持功能: 仅 支持 ! () && || 四种逻辑运算, 优先级从高到低

	实现思路:
	1.分词(tokenize)
 	输入: "(1||2)&&!3"	         输出: [(, 1, ||, 2, ), &&, !, 3]

	2.规则项匹配，转化为逻辑表达式, T 代表 true, F 代表 false
	输入: [( 1, ||, 2, ), &&, !, 3]   输出: [(, T, ||, F, ), &&, ! T]

	3.表达式求值
	输入: [(, T, ||, F, ), &&, ! T]  输出: false

 */
func (d DecisionTree) match(f ...Factor) (string, bool) {

	//1.分词
	exp := tokenize(d.Exp)
	if len(exp) == 0 {
		return "", false
	}

	//2.匹配规则
	matchRuleItem(exp, f...)

	//3.表达式求值
	matched := calculateExp(exp)

	if matched {
		return d.Value, true
	}

	return "", false
}

//分词: 表达式字符串转数组
func tokenize(exp string) []string {
	result := []string{}
	l := len(exp)
	for i := 0; i < l; i++ {
		char := string(exp[i])
		switch char {

		//空格跳过
		case Space:
			continue

		//单个字节的操作符
		case LeftBracket, RightBracket, LogicNon:
			result = append(result, char)

		//数字
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// 数字
			j := i
			digit := ""
			for ; j < l && unicode.IsDigit(rune(exp[j])); j++ {
				digit += string(exp[j])
			}
			result = append(result, digit)
			i = j - 1 // i 向前跨越一个整数，由于执行了一步多余的 j++，需要减 1

		default:
			// 多字节操作符
			j := i
			operator := ""
			for ; j < l && (string(exp[j]) == BitwiseAnd || string(exp[j]) == BitwiseOr); j++ {
				operator += string(char)
			}
			i = j - 1 // i 向前跨越一个整数，由于执行了一步多余的 j++，需要减 1

			//不是支持的操作符直接返回
			if !isSign(operator) {
				return []string{}
			}

			result = append(result, operator)
		}
	}
	return result
}

//计算表达式
func calculateExp(exp []string) bool {

	//存放操作数， 类型为 bool
	valStack := NewStack()
	//存放操作符，类型为 string
	opStack := NewStack()
	//# 表示栈为空
	opStack.Push("#")

	//逻辑计算，计算成功返回 true，失败返回 false
	solve := func() bool {
		if opStack.Peek().(string) != LogicNon {

			//表达式不合法， 直接返回 false
			if valStack.Size() < 2 || opStack.Peek() == "#" {
				return false
			}

			v1 := valStack.Pop().(bool)
			v2 := valStack.Pop().(bool)
			op := opStack.Pop().(string)
			v := calculate(v1, v2, op)
			valStack.Push(v)

		} else {

			if valStack.Size() < 1 || opStack.Peek() == "#" {
				return false
			}

			v := valStack.Pop().(bool)
			valStack.Push(!v)
			opStack.Pop()
		}

		return true
	}

	expLen := len(exp)

	// 遍历整个表达式
	for i := 0; i < expLen; i++ {
		char := exp[i]

		switch char {

		case LogicTrue, LogicFalse:
			valStack.Push(char == LogicTrue)

		case LeftBracket:
			// 左括号直接入栈
			opStack.Push(LeftBracket)

		case RightBracket:
			// 右括号则弹出元素直到遇到左括号
			for opStack.Peek() != LeftBracket {
				if success := solve(); !success {
					return false
				}
			}
			opStack.Pop()
		default:
			//对比栈顶元素和当前元素优先级， 先计算高优先级的
			for priority(opStack.Peek().(string)) >= priority(char) {
				if success := solve(); !success {
					return false
				}
			}
			opStack.Push(char)
		}
	}

	// 栈不空则全部输出
	for opStack.Peek() != "#" {
		if success := solve(); !success {
			return false
		}
	}

	return !valStack.IsEmpty() && valStack.Peek().(bool)
}

//是否为操作符
func isSign(s string) bool {
	signList := []string{LeftBracket, RightBracket, LogicAnd, LogicOr, LogicNon}
	return StringIsIn(s, signList...)
}

//判断每一条规则项的的逻辑取值
func matchRuleItem(postFixExp []string, f ...Factor) {
	for i := range postFixExp {
		e := postFixExp[i]

		//符号直接跳过
		if isSign(e) {
			continue
		}

		//默认设置为 false
		postFixExp[i] = LogicFalse

		//取每一条规则项
		ruleItem, exist := DefaultRuleEngine.Items[e]
		//不存在认为没有匹配上
		if !exist {
			continue
		}

		//规则项匹配函数
		matcher, exist := operatorMatcher[ruleItem.Operator]
		//匹配函数不存在
		if !exist {
			continue
		}

		//遍历输入规则因子
		for j := range f {
			if f[j].FactorType != ruleItem.FactorType {
				continue
			}

			matchSuccess := matcher.match(f[j].FactorValue, ruleItem.FactorValue)
			if matchSuccess {
				postFixExp[i] = LogicTrue
				break
			}
		}
	}
}

//获取当前操作符运算优先级
func priority(op string) int {
	switch op {
	case LogicNon:
		return 3
	case LogicAnd:
		return 2
	case LogicOr:
		return 1
	default:
		return 0
	}
}

//逻辑运算
func calculate(v1, v2 bool, op string) bool {
	if op == LogicAnd {
		return v1 && v2
	}

	if op == LogicOr {
		return v1 || v2
	}

	//不可能走到这里
	return false
}
