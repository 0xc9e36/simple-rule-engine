package utils

//StringIsIn 判断字符串是否在输入串组中
func StringIsIn(in string, arg ...string) bool {
	for _, v := range arg {
		if in == v {
			return true
		}
	}
	return false
}

