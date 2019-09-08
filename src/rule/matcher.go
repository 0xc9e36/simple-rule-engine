package rule

//操作符比较
var operatorMatcher = map[string]Matcher{

	//字符串相等规则
	"=": MatcherFunc(Equal),

	">=": MatcherFunc(GreaterOrEqual),

	"<=": MatcherFunc(LessOrEqual),



	//TODO 自定义规则比较方法
}

type Matcher interface {
	match(string, string) bool
}

//规则适配器
type MatcherFunc func(string, string) bool


// in 为输入值, s2 为数据库配置值
func (m MatcherFunc) match(in, factorVal string) bool {
	return m(in, factorVal)
}


//== 相等规则
func Equal(in, factorVal string) bool {
	return in == factorVal
}

// >= 规则函数
func GreaterOrEqual(in, factorVal string) bool {
	return in >= factorVal
}

// <= 规则函数
func LessOrEqual(in, factorVal string) bool {
	return in <= factorVal
}