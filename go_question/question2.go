package main

import "fmt"

//func main() {
//	s := "abc"
//	fmt.Println(isUniqueString(s))
//}
//func isUniqueString(s string) bool {
//	/*cntM := make(map[int32]int)
//	for _, v := range  s {
//		cntM[v]++
//	}
//	for _,v := range cntM {
//		if v > 1 {
//			return false
//		}
//	}
//	return true*/
//	for _, v := range  s {
//		if strings.Count(s,string(v)) > 1 {
//			return false
//		}
//	}
//	return true
//}

func main() {
	s := make([]int, 3, 9)
	fmt.Println(len(s))
	s2 := s[4:8]
	fmt.Println(len(s2))

}
