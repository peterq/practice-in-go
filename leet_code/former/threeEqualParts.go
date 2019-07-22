package former

func threeEqualParts(A []int) []int {
	// 算出左边为 1 的个数的map
	mp := make([]int, len(A))
	oneChangeMap := make([]int, 0)
	if A[0] == 1 {
		mp[0] = 1
		oneChangeMap = append(oneChangeMap, 0)
	}
	for i := 1; i < len(A); i++ {
		mp[i] = mp[i-1]
		if A[i] == 1 {
			mp[i]++
			oneChangeMap = append(oneChangeMap, i)
		}
	}
	ones := mp[len(A)-1]
	if ones%3 != 0 {
		return []int{-1, -1}
	}
	if ones == 0 {
		return []int{0, 2}
	}
	var iMin, iMax, jMin, jMax int
	// 找到分界线
	iMin = oneChangeMap[len(oneChangeMap)/3-1] - 1
	iMax = oneChangeMap[len(oneChangeMap)/3] - 1
	jMin = oneChangeMap[len(oneChangeMap)/3*2-1] - 1
	jMax = oneChangeMap[len(oneChangeMap)/3*2] - 1
	// 确定后面有多少个0
	zeros := len(A) - oneChangeMap[len(oneChangeMap)-1]
	i := iMin + zeros
	j := jMin + zeros
	if i <= iMax && j <= jMax {
		firstLength := i + 1
		secondLength := j - i
		thirdLength := len(A) - 1 - j
		minLength := firstLength
		if secondLength < minLength {
			minLength = secondLength
		}
		if thirdLength < minLength {
			minLength = thirdLength
		}
		for k := 0; k < minLength; k++ {
			if A[i-k] == A[len(A)-1-k] && A[i-k] == A[j-k] {
				continue
			} else {
				return []int{-1, -1}
			}
		}
		for k := minLength; k < firstLength; k++ {
			if A[i-k] != 0 {
				return []int{-1, -1}
			}
		}
		for k := minLength; k < secondLength; k++ {
			if A[j-k] != 0 {
				return []int{-1, -1}
			}
		}
		for k := minLength; k < thirdLength; k++ {
			if A[len(A)-1-k] != 0 {
				return []int{-1, -1}
			}
		}
		return []int{i, j + 1}
	}
	return []int{-1, -1}
}
