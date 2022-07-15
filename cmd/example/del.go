package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3, 4, 5}
	for i, val := range a {
		fmt.Printf("i=%v, val=%v, a=%v\n", i, val, a)
		if val%2 == 0 {
			a = append(a[:i], a[i+1:]...)
		}
	}
	fmt.Println("a=", a)
}

// 内存: a   -> 1 2 3 4
// i=1,val=2 -> 1 2 3 4   a=append(a[:1],a[2:4]) 
// i=2,val=4 -> 1 3 4 [4]  a=append(a[:2],a[3:3]) 
// i=3,val=4 -> 1 3 [4 4]  此时炸掉,a=append(a[:3],a[4:2]) 

// 内存: a   -> 1 2 3 4 5 
// i=1,val=2 -> 1 2 3 4 5   a=append(a[:1],a[2:5]) 
// i=2,val=4 -> 1 3 4 5 [5] a=append(a[:2],a[3:4]) 
// i=3,val=5 -> 1 3 5 [5 5] a不变
// i=4,val=5 -> 1 3 5 [5 5] a不变