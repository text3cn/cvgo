package mathkit

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 取随机数
// 范围: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// 保留指定位随机小数
func RandDecimals(min, max float64, precision ...int) float64 {
	prec := 2
	if len(precision) > 0 {
		prec = precision[0]
	}
	result := min + rand.Float64()*(max-min)
	scale := math.Pow(10, float64(prec))
	return math.Round(result*scale) / scale
}

// 对浮点数进保留几位小数
func Floor(value float64, decimal int) float64 {
	decimalStr := strconv.Itoa(decimal)
	finalValue, _ := strconv.ParseFloat(fmt.Sprintf("%."+decimalStr+"f", value), 64)
	return finalValue
}

// 对浮点数进保留几位小数，0 填充小数位，0 填充只能是 string 类型才能显示完整
func FloorWithZeroPad(value float64, decimal int) string {
	decimalStr := strconv.Itoa(decimal)
	finalValue, _ := strconv.ParseFloat(fmt.Sprintf("%."+decimalStr+"f", value), 64)
	// 小数不足 decimal 位长度时 0 填充
	return fmt.Sprintf("%.4f", finalValue)
}
