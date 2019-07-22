package former

func minFlipsMonoIncr(S string) int {
	leftNeedFlip := make([]int, len(S)+1)
	rightNeedFlip := make([]int, len(S)+1)
	// 先统计当 firstOne 为 i 时 左边为1的个数
	for i := 0; i < len(leftNeedFlip); i++ {
		if i == 0 {
			continue
		}
		leftNeedFlip[i] = leftNeedFlip[i-1]
		if string(S[i-1]) == "1" {
			leftNeedFlip[i]++
		}
	}
	// 统计右边为1的个数
	for i := len(S); i >= 0; i-- {
		if i == len(S) {
			continue
		}
		rightNeedFlip[i] = rightNeedFlip[i+1]
		if string(S[i]) == "0" {
			rightNeedFlip[i]++
		}
	}
	flip := len(S)
	for i := 0; i < len(leftNeedFlip); i++ {
		fl := leftNeedFlip[i] + rightNeedFlip[i]
		if fl < flip {
			flip = fl
		}
	}
	return flip
}
