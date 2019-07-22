package former

/**
给定一个整数数组和一个目标值，找出数组中和为目标值的两个数。

你可以假设每个输入只对应一种答案，且同样的元素不能被重复利用。

示例:

给定 nums = [2, 7, 11, 15], target = 9

因为 nums[0] + nums[1] = 2 + 7 = 9
所以返回 [0, 1]
*/

func twoSum1(nums []int, target int) []int {
	firstIndex := 0
	secondIndex := 1

	for firstIndex = 0; firstIndex < len(nums)-1; firstIndex++ {
		for secondIndex = firstIndex + 1; secondIndex < len(nums); secondIndex++ {
			if nums[firstIndex]+nums[secondIndex] == target {
				return []int{firstIndex, secondIndex}
			}
		}
	}
	return nil
}

func twoSum(nums []int, target int) []int {
	// 差值map
	mp := make(map[int]int)
	for index, v := range nums {
		need := target - v
		if idx, ok := mp[need]; ok {
			return []int{idx, index}
		}
		mp[v] = index
	}
	return nil
}
