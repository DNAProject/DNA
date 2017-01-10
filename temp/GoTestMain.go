package main
import (
	"fmt"
)

func main() {

	//test slice append
	s1 := []int{1,2,3}
	s2 := []int{4,5,6}
	s1 = append(s1,s2[:]...)
	fmt.Println(s1)


}
