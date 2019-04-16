package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

//var mu sync.Mutex
var wg sync.WaitGroup

func main() {
	data := []int{0, 1}
	//data := []int{0}

	jobs := []job{
		func(in, out chan interface{}) {
			for _, val := range data {
				log.Println("*** " + strconv.Itoa(int(val)))
				out <- strconv.Itoa(int(val))
				//log.Println("after put in out")
			}
			close(out)
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
	wg.Add(4)

	for _, j := range jobs {
		go j(in, out)
	}
}

func SingleHash(in, out chan interface{}) {
	for val := range out {
		log.Println("*1*" + " -> val " + val.(string))

		tmp := DataSignerCrc32(val.(string)) + "~" + DataSignerCrc32(DataSignerMd5(val.(string)))
		log.Println("*1* => " + tmp)
		in <- tmp
		wg.Done()
	}

	//wg.Done()
}

func MultiHash(in, out chan interface{}) {
	for val := range in {
		log.Println("*2*" + " -> val " + val.(string))
		steps := []string{"0", "1", "2", "3", "4", "5"}
		res := ""

		log.Println("*2* ->" + "before cycle")
		for _, step := range steps {
			res += DataSignerCrc32(step + val.(string))
		}

		log.Println("*2* => " + res)
		wg.Done()
		in <- res
	}

	log.Println("END IN...")
	//wg.Done()
	close(in)
}

func CombineResults(in, out chan interface{}) {
	for {
		wg.Wait()
		var res []string
		val, opened := <-in
		if opened {
			log.Println("*3*" + " -> val " + val.(string))
			res = append(res, val.(string)+"_")
		} else {
			// sort
			log.Println("Finish *3* => ")
			log.Println(res)
			break
		}
	}
}
