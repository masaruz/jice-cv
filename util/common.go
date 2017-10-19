package util

// Absolute for int64
func Absolute(num int64) int64 {
	if num < 0 {
		num = -num
	}
	return num
}
