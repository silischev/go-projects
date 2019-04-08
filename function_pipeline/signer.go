package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	//data := []int{0, 1}
	data := []int{0}

	jobs := []job{
		func(in, out chan interface{}) {
			for _, val := range data {
				in <- strconv.Itoa(int(val))
			}
		},
		SingleHash,
		MultiHash,
		CombineResults,
	}

	ExecutePipeline(jobs...)

	fmt.Scanln()
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})

	for _, j := range jobs {
		go j(in, out)
	}
}

func SingleHash(in, out chan interface{}) {
	val := <-in
	tmp := DataSignerCrc32(val.(string)) + "~" + DataSignerCrc32(DataSignerMd5(val.(string)))
	log.Println(tmp)
	out <- tmp
}

func MultiHash(in, out chan interface{}) {
	inVal := <-out
	steps := []string{"0", "1", "2", "3", "4", "5"}
	res := ""

	for _, val := range steps {
		res += DataSignerCrc32(val + inVal.(string))
	}

	log.Println(res)
	out <- res
}

func CombineResults(in, out chan interface{}) {
	log.Println(<-out)
}
