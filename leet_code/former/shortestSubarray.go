package former

/**

返回 A 的最短的非空连续子数组的长度，该子数组的和至少为 K 。

如果没有和至少为 K 的非空子数组，返回 -1 。



示例 1：

输入：A = [1], K = 1
输出：1
示例 2：

输入：A = [1,2], K = 4
输出：-1
示例 3：

输入：A = [2,-1,2], K = 3
输出：3


提示：

1 <= A.length <= 50000
-10 ^ 5 <= A[i] <= 10 ^ 5
1 <= K <= 10 ^ 9
*/

func shortestSubarray1(A []int, K int) int {
	aLg := len(A)
	// 假设最小长度是1, 找出起点
	lastSum := make([]int, aLg)
	skip := make([]bool, aLg)
	for minLg := 1; minLg <= aLg; minLg++ {
		for leftIndex := 0; leftIndex <= aLg-minLg; leftIndex++ {
			if A[leftIndex] <= 0 || skip[leftIndex] {
				continue
			}
			sum := lastSum[leftIndex] + A[leftIndex+minLg-1]
			if sum >= K {
				return minLg
			}
			lastSum[leftIndex] = sum
			if sum <= 0 {
				skip[leftIndex] = true
			}
		}
	}
	return -1
}

func shortestSubarray(A []int, K int) int {
	aLg := len(A)
	sumMap := make([]int, aLg)
	lastSum := 0
	for i, v := range A {
		if A[i] >= K {
			return 1
		}
		sumMap[i] = lastSum + v
		lastSum = sumMap[i]
	}

	for leftIndex := 1; leftIndex < aLg-1; leftIndex++ {
		for rIdx := leftIndex + 1; rIdx < aLg; rIdx++ {
			if sumMap[rIdx]-sumMap[leftIndex-1] >= K {
				return rIdx - leftIndex + 1
			}
		}
	}

	return -1
}
