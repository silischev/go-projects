package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

//var mu sync.Mutex
//var wg sync.WaitGroup

func main() {
	data := []int{0, 1}
	//data := []int{0}

	jobs := []job{
		func(in, out chan interface{}) {
			for _, val := range data {
				log.Println("*** " + strconv.Itoa(int(val)))
				out <- strconv.Itoa(int(val))
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
	var outChannels []chan interface{}

	for key, job := range jobs {
		var in chan interface{}

		if key == 0 {
			in = make(chan interface{})
		} else {
			in = outChannels[key-1]
		}

		out := make(chan interface{})
		outChannels = append(outChannels, out)

		go job(in, out)
	}
}

func SingleHash(in, out chan interface{}) {
	for val := range in {
		log.Println("*1*" + " -> val " + val.(string))

		tmp := DataSignerCrc32(val.(string)) + "~" + DataSignerCrc32(DataSignerMd5(val.(string)))
		log.Println("*1* => " + tmp)
		out <- tmp
	}

	close(out)
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
		out <- res
	}

	close(out)
}

func CombineResults(in, out chan interface{}) {
	var values []string
	res := ""

	for val := range in {
		log.Println("*3*" + " -> val " + val.(string))
		//value, _ := strconv.Atoi(val.(string))
		//value, _ := strconv.ParseInt(val.(string), 10, 64)
		values = append(values, val.(string))
		sort.Slice(values, func(i, j int) bool {
			return values[i] < values[j]
		})
	}

	log.Println("Finish *3* => ")

	for _, val := range values {
		if res != "" {
			res += "_" + val
		} else {
			res = val
		}
	}

	log.Println(res)
}
