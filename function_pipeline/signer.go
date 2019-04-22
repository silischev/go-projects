package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

func main() {
	data := []int{0, 1}
	res := ""

	jobs := []job{
		func(in, out chan interface{}) {
			for _, val := range data {
				log.Println("*** " + strconv.Itoa(val))
				out <- val
			}
		},
		SingleHash,
		MultiHash,
		CombineResults,
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				log.Fatal("cant convert result data to string")
			}
			res = data
		}),
	}

	ExecutePipeline(jobs...)

	log.Println(res)
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
		start := time.Now()
		//log.Println("*1*" + " -> val " + strconv.Itoa(val.(int)))

		wg := &sync.WaitGroup{}
		res1 := ""
		res2 := ""

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			res1 = DataSignerCrc32(strconv.Itoa(val.(int)))
		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			res2 = DataSignerCrc32(DataSignerMd5(strconv.Itoa(val.(int))))
		}(wg)

		wg.Wait()

		end := time.Since(start)
		log.Println(fmt.Sprintf("******** 1 => %s", end))

		out <- res1 + "~" + res2
	}
}

func MultiHash(in, out chan interface{}) {
	for val := range in {
		start := time.Now()

		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}

		//log.Println("*2*" + " -> val " + val.(string))

		steps := []int{0, 1, 2, 3, 4, 5}
		var results []string = make([]string, 6)

		//log.Println("*2* ->" + "before cycle")
		for _, step := range steps {
			wg.Add(1)
			go func(wg *sync.WaitGroup, mu *sync.Mutex, st int, results []string) {
				defer wg.Done()
				res := DataSignerCrc32(strconv.Itoa(st) + val.(string))

				mu.Lock()
				results[st] = res
				mu.Unlock()
			}(wg, mu, step, results)
		}

		wg.Wait()

		res := ""
		for _, val := range results {
			res += val
		}

		end := time.Since(start)
		log.Println(fmt.Sprintf("******** 2 => %s", end))

		out <- res
	}
}

func CombineResults(in, out chan interface{}) {
	var values []string
	res := ""

	for val := range in {
		start := time.Now()

		//log.Println("*3*" + " -> val " + val.(string))
		values = append(values, val.(string))

		end := time.Since(start)
		log.Println(fmt.Sprintf("******** 3 => %s", end))
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	log.Println("Finish *3* => ")

	for _, val := range values {
		if res == "" {
			res = val
			continue
		}

		res += "_" + val
	}

	out <- res
}
