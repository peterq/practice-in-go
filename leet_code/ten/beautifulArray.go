package ten

import "log"

func beautifulArray(N int) []int {
	arr := make([]int, N)
	for i, _ := range arr {
		arr[i] = i + 1
	}
	stk := [][]int{arr}
	for len(stk) > 0 {
		child := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		lg := len(child)
		log.Println(arr, lg)
		if lg == 4 {
			t := child[1]
			child[1] = child[0]
			child[0] = t
			t = child[3]
			child[3] = child[2]
			child[2] = t
			continue
		} else if lg == 3 {
			t := child[2]
			child[2] = child[1]
			child[1] = t
			continue
		} else if lg < 3 {
			continue
		}
		sm := lg / 3
		tmp := append([]int{}, child[sm:lg-sm]...)
		_ = append(child[sm:sm], child[lg-sm:]...)
		_ = append(child[2*sm:2*sm], tmp...)
		a := child[:sm]
		b := child[sm : 2*sm]
		c := child[2*sm:]
		stk = append(stk, a, b, c)
	}
	return arr
}
