package main

import (
	"log"
	"sort"
	"strconv"
	"sync"
)

func main() {
	data := []int{0, 1}

	jobs := []job{
		func(in, out chan interface{}) {
			for _, val := range data {
				log.Println("*** " + strconv.Itoa(val))
				out <- strconv.Itoa(val)
			}
			close(out)
		},
		SingleHash,
		MultiHash,
		CombineResults,
	}

	ExecutePipeline(jobs...)

	//fmt.Scanln()
}

func ExecutePipeline(jobs ...job) {
	var outChannels []chan interface{}
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for key, j := range jobs {
		if key > 0 {
			in = outChannels[key-1]
		}

		out := make(chan interface{})
		outChannels = append(outChannels, out)

		wg.Add(1)
		go func(wg *sync.WaitGroup, j job, in, out chan interface{}) {
			j(in, out)
			close(out)
			wg.Done()
		}(wg, j, in, out)
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	for val := range in {
		log.Println("*1*" + " -> val " + strconv.Itoa(val.(int)))
		tmp := DataSignerCrc32(strconv.Itoa(val.(int))) + "~" + DataSignerCrc32(DataSignerMd5(strconv.Itoa(val.(int))))
		out <- tmp
	}
}

func MultiHash(in, out chan interface{}) {
	for val := range in {
		log.Println("*2*" + " -> val " + val.(string))
		steps := []string{"0", "1", "2", "3", "4", "5"}
		res := ""

		//log.Println("*2* ->" + "before cycle")
		for _, step := range steps {
			res += DataSignerCrc32(step + val.(string))
		}

		out <- res
	}
}

func CombineResults(in, out chan interface{}) {
	var values []string
	res := ""

	for val := range in {
		log.Println("*3*" + " -> val " + val.(string))
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

	//log.Println(res)

	out <- res
}
