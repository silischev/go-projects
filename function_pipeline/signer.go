package main

import (
	"log"
	"sort"
	"strconv"
	"sync"
)

func main() {
	data := []int{0, 1}
	res := ""

	jobs := []job{
		func(in, out chan interface{}) {
			for _, val := range data {
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
	mu := &sync.Mutex{}
	wgCommon := &sync.WaitGroup{}

	for val := range in {
		wgCommon.Add(1)
		go func(wgCommon *sync.WaitGroup, mu *sync.Mutex, val int) {
			defer wgCommon.Done()

			wg := &sync.WaitGroup{}
			res1 := ""
			res2 := ""

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				res1 = DataSignerCrc32(strconv.Itoa(val))
			}(wg)

			wg.Add(1)
			go func(wg *sync.WaitGroup, mu *sync.Mutex) {
				defer wg.Done()

				mu.Lock()
				md5 := DataSignerMd5(strconv.Itoa(val))
				mu.Unlock()
				res2 = DataSignerCrc32(md5)
			}(wg, mu)

			wg.Wait()
			out <- res1 + "~" + res2
		}(wgCommon, mu, val.(int))
	}

	wgCommon.Wait()
}

func MultiHash(in, out chan interface{}) {
	wgCommon := &sync.WaitGroup{}

	for val := range in {
		wgCommon.Add(1)
		go func(wgCommon *sync.WaitGroup, val string, out chan interface{}) {
			defer wgCommon.Done()

			wg := &sync.WaitGroup{}
			mu := &sync.Mutex{}

			steps := []int{0, 1, 2, 3, 4, 5}
			var results []string = make([]string, 6)

			for _, step := range steps {
				wg.Add(1)
				go func(wg *sync.WaitGroup, mu *sync.Mutex, st int, results []string) {
					defer wg.Done()
					res := DataSignerCrc32(strconv.Itoa(st) + val)

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

			out <- res
		}(wgCommon, val.(string), out)
	}

	wgCommon.Wait()
}

func CombineResults(in, out chan interface{}) {
	var values []string
	res := ""

	for val := range in {
		values = append(values, val.(string))
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	for _, val := range values {
		if res == "" {
			res = val
			continue
		}

		res += "_" + val
	}

	out <- res
}
