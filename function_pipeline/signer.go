package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	input := []int{0, 1}

	initiator := func(in, out chan interface{}) {
		for _, i := range input {
			out <- i
		}
	}

	jobs := []job{initiator, SingleHash, MultiHash, CombineResults}

	ExecutePipeline(jobs...)
}

func ExecutePipeline(workers ...job) {
	wg := &sync.WaitGroup{}
	outChannels := make([]chan interface{}, len(workers))

	for num, worker := range workers {
		var in chan interface{}

		out := make(chan interface{}, MaxInputDataLen)
		outChannels[num] = out

		if num > 0 {
			in = outChannels[num-1]
		} else {
			in = make(chan interface{}, MaxInputDataLen)
		}

		wg.Add(1)
		go func(w job) {
			defer wg.Done()
			w(in, out)
			close(out)
		}(worker)
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wgCommon := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for input := range in {
		wgCommon.Add(1)
		go func(val int) {
			defer wgCommon.Done()

			ch1 := make(chan interface{})
			ch2 := make(chan interface{})

			go func() {
				ch1 <- DataSignerCrc32(strconv.Itoa(val))
			}()

			go func() {
				mu.Lock()
				res := DataSignerMd5(strconv.Itoa(val))
				mu.Unlock()

				ch2 <- DataSignerCrc32(res)
			}()

			out <- (<-ch1).(string) + "~" + (<-ch2).(string)
		}(input.(int))
	}

	wgCommon.Wait()
}

func MultiHash(in, out chan interface{}) {
	wgCommon := &sync.WaitGroup{}

	for input := range in {
		wgCommon.Add(1)
		go func(val interface{}) {
			defer wgCommon.Done()
			result := ""
			tempRes := make([]string, 6)
			wg := &sync.WaitGroup{}

			wg.Add(6)
			for i := 0; i < 6; i++ {
				go func(num int) {
					defer wg.Done()
					ch := make(chan interface{})

					go func() {
						ch <- DataSignerCrc32(strconv.Itoa(num) + val.(string))
					}()

					tempRes[num] = (<-ch).(string)
				}(i)
			}

			wg.Wait()

			for i := 0; i < 6; i++ {
				result = result + tempRes[i]
			}

			out <- result
		}(input)
	}

	wgCommon.Wait()
}

func CombineResults(in, out chan interface{}) {
	var data []string

	for i := range in {
		data = append(data, i.(string))
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	out <- strings.Join(data, "_")
}
