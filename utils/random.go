package utils

import (
	"math/rand"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomSubset 随机选择字符串数组的子集并返回
func RandomSubset(arr []string) []string {
	// 创建一个字符串切片来存储选择的子集
	subset := make([]string, 0)

	// 如果输入的数组为空，直接返回空子集
	if len(arr) == 0 {
		return subset
	}

	// 创建一个映射来跟踪已选择的元素
	alreadySelected := make(map[int]bool)

	// 至少选择一个元素，因此子集大小至少为1，最大为输入数组的长度
	subsetSize := 1 + r.Intn(len(arr)-1)

	// 从输入的数组中随机选择子集
	for i := 0; i < subsetSize; {
		// 随机选择数组中的一个索引
		index := r.Intn(len(arr))
		// 检查是否已经选择了该元素
		if !alreadySelected[index] {
			// 将选择的元素添加到子集中
			subset = append(subset, arr[index])
			// 标记该索引已被选择
			alreadySelected[index] = true
			i++ // 只有在选择了一个不重复的元素时才增加i
		}
	}

	return subset
}

func GenerateRandomArray(length int) []string {
	l := make([]string, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		f := rand.NormFloat64()*0.566 + 3
		num := int(f + 0.5) //加0.5实现四舍五入
		if num < 1 {
			num = 1
		} else if num > 5 {
			num = 5
		}

		l[i] = strconv.Itoa(num) //将数字转换为字符串
	}

	return l
}

// RandomString2 随机返回字符串数组中的一个字符串
func RandomString2(arr []string) string {

	// 如果输入的数组为空，则返回空字符串
	if len(arr) == 0 {
		return ""
	}

	// 随机选择一个数组索引
	index := r.Intn(len(arr))

	// 返回选择的字符串
	return arr[index]
}

func RandomString(arr []string) string {
	// 如果输入的数组为空，则返回空字符串
	if len(arr) == 0 {
		return ""
	}

	// 随机生成满足正态分布的值
	norm := r.NormFloat64()

	// 限定生成的随机数在[-3,3]范围内，这对应正态分布中约99.7%的范围
	if norm < -3 {
		norm = -3
	} else if norm > 3 {
		norm = 3
	}

	// 将[-3,3]离散化为len(arr)个区间，找出norm所属的区间
	index := int((norm + 3) / 6 * float64(len(arr)))

	if index >= len(arr) {
		index = len(arr) - 1
	}

	// 返回选择的字符串
	return arr[index]
}
