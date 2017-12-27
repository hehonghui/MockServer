package main

import "fmt"

type User struct {
	Name string
}


func main() {
	arr := []int {3,4,6}
	arr = append(arr, 7)
	for _, val := range arr {
		fmt.Println(val)
	}

	arr2 := make([]*User, 0)
	u1 := new(User)
	u1.Name = "MrX"
	arr2 = append(arr2, u1)

	u2 := new(User)
	u2.Name = "MrA"
	arr2 = append(arr2, u2)

	u3 := new(User)
	u3.Name = "Mr3"
	arr2 = append(arr2, u3)

	for _, val := range arr2 {
		fmt.Println(val.Name)
	}
}
