package pack

import (
	"fmt"
	"math/rand"
)

func Fact(n int) uint64 {
	if n > 1 {
		var res uint64 = 1
		for i := 1; i < n+1; i++ {
			res = res * uint64(i)
		}
		fmt.Println("Факториал чила ", n, "! = ", res)
		return res
	}
	return 0
}

func Rand() int {
	r := rand.Intn(100)
	fmt.Println("Случайное число: ", r)
	return r

}

func Sum(a []int) int {
	sum := 0
	for i := 0; i < len(a); i++ {
		sum += a[i]
	}
	fmt.Println("Сумма ряда: ", sum)
	return sum
}

func CreateSeries(n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(100)
	}
	return a
}
