package former

/**

给定两个大小为 m 和 n 的有序数组 nums1 和 nums2 。

请找出这两个有序数组的中位数。要求算法的时间复杂度为 O(log (m+n)) 。

你可以假设 nums1 和 nums2 不同时为空。

示例 1:

nums1 = [1, 3]
nums2 = [2]

中位数是 2.0
示例 2:

nums1 = [1, 2]
nums2 = [3, 4]

中位数是 (2 + 3)/2 = 2.5
*/

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	nums1Index, nums2Index := 0, 0
	target := make([]int, 0)
	for {
		if nums1Index == len(nums1) {
			for ; nums2Index < len(nums2); nums2Index++ {
				target = append(target, nums2[nums2Index])
			}
			break
		}
		if nums2Index == len(nums2) {
			for ; nums1Index < len(nums1); nums1Index++ {
				target = append(target, nums1[nums1Index])
			}
			break
		}
		if nums1[nums1Index] < nums2[nums2Index] {
			target = append(target, nums1[nums1Index])
			nums1Index++
		} else {
			target = append(target, nums2[nums2Index])
			nums2Index++
		}
	}
	if len(target)%2 == 1 {
		return float64(target[len(target)/2])
	} else {
		return float64(target[len(target)/2]+target[len(target)/2-1]) / 2
	}
}
