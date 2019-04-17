package main

import (
	"fmt"
	"log"
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
	//wg.Add(4)
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
		//wg.Done()
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
		out <- res
	}

	log.Println("END IN...")
	//wg.Done()
	close(in)
}

func CombineResults(in, out chan interface{}) {
	var res interface{}

	for val := range in {
		log.Println("*3*" + " -> val " + val.(string))
		res = val.(string) + "_"
	}

	log.Println("Finish *3* => ")
	log.Println(res)

	/* for {
		//wg.Wait()
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
	} */
}
