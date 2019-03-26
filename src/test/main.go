package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(isPalindrome(121))
}

func isPalindrome(x int) bool {
	if x >= 0 {
		var i int
		var t int64
		var y = x
		for {
			r := x % 10
			x = x / 10
			t += int64(r) * int64(math.Pow10(9-i))
			if x == 0 {
				t = t / int64(math.Pow10(9-i))
				return int(t) == y
			}
			i++
		}
	}
	return false

}
