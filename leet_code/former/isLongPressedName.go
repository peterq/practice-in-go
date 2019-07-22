package former

func isLongPressedName(name string, typed string) bool {
	if len(typed) < len(name) {
		return false
	}
	nameLeft := name
	typedLeft := typed
	var nameCh, typedCh string
	var typedNum, nameNum int
	for len(nameLeft) > 0 {
		nameLeft, nameCh, nameNum = sameString(nameLeft)
		typedLeft, typedCh, typedNum = sameString(typedLeft)
		if nameCh != typedCh {
			return false
		}
		if typedNum < nameNum {
			return false
		}
		if len(typedLeft) < len(nameLeft) {
			return false
		}
	}
	return true
}

func sameString(str string) (left, ch string, num int) {
	for i := 0; i < len(str); i++ {
		if i == 0 {
			num = 1
			ch = string(str[0])
		} else {
			if ch == string(str[i]) {
				num++
			} else {
				break
			}
		}
	}
	left = str[num:]
	return
}
