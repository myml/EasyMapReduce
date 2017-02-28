package easyMapReduce

import (
	"fmt"
	"strconv"
	"testing"
)

const goroutineLimit = 10

func TestMapReduce(t *testing.T) {
	var data []int
	for i := 0; i < 1000; i++ {
		data = append(data, i)
	}
	Map := func(i int) int {
		return i * i
	}
	var result []int
	Reduce := func(r int) error {
		result = append(result, r)
		//		log.Println(r)
		return nil
	}
	err := mapreduce(data, goroutineLimit, Map, Reduce)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestMapReduce2(t *testing.T) {
	source := make(chan int)
	go func() {
		for i := 0; i < 1000; i++ {
			source <- i
		}
		//！！！通道类型的源要自己关闭
		close(source)
	}()
	Map := func(i int) *string {
		if i%2 == 0 {
			return nil
		}
		s := strconv.Itoa(i)
		return &s
	}
	var result []string
	Reduce := func(r *string) error {
		if r != nil {
			result = append(result, *r)
		}
		return nil
	}
	err := mapreduce(source, goroutineLimit, Map, Reduce)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
