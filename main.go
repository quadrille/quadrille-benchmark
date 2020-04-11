package main

import (
	"fmt"
	"github.com/quadrille/quadgo"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func getRandomLocation(minLat, maxLat, minLong, maxLong float64) (float64, float64) {
	return minLat + rand.Float64()*(maxLat-minLat), minLong + rand.Float64()*(maxLong-minLong)
}

func main() {
	quadrilleAddr := "localhost:5679"
	if len(os.Args) > 1 {
		quadrilleAddr = os.Args[1]
	}
	noOfOps := 30000
	quadClient := quadgo.NewClient(quadrilleAddr)

	insertsPerSec := benchmarkInsert(noOfOps, quadClient)

	time.Sleep(time.Second * 2)

	getsPerSec := benchmarkGet(noOfOps, quadClient)
	nearbyPerSec := benchmarkNearby(noOfOps, quadClient)

	cleanup(noOfOps, quadClient)

	fmt.Printf("insert: %d qps\n", insertsPerSec)
	fmt.Printf("get: %d qps\n", getsPerSec)
	fmt.Printf("nearby: %d qps\n", nearbyPerSec)
}

func benchmarkInsert(noOfOps int, quadClient *quadgo.QuadrilleClient) int {
	locIDCh := make(chan string, 1000)
	go func() {
		for i := 0; i < noOfOps; i++ {
			locIDCh <- "loc" + strconv.Itoa(i)
		}
		close(locIDCh)
	}()
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for locID := range locIDCh {
				//fmt.Println(locID)
				lat, long := getRandomLocation(12.8623132, 13.0573792, 77.4743233, 77.7871233)
				err := quadClient.Insert(locID, lat, long, map[string]interface{}{})
				//_, err := quadClient.Get(locID)
				//fmt.Println(loc)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
	//fmt.Println("start wait")
	wg.Wait()
	sec := time.Since(startTime).Seconds()
	return int(float64(noOfOps) / sec)
}

func benchmarkGet(noOfOps int, quadClient *quadgo.QuadrilleClient) int {
	locIDCh := make(chan string, 1000)
	go func() {
		for i := 0; i < noOfOps; i++ {
			locIDCh <- "loc" + strconv.Itoa(i)
		}
		close(locIDCh)
	}()
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for locID := range locIDCh {
				_, err := quadClient.Get(locID)
				if err != nil {
					//fmt.Println(err)
				}
			}
		}()
	}
	//fmt.Println("start wait")
	wg.Wait()
	sec := time.Since(startTime).Seconds()
	return int(float64(noOfOps) / sec)
}

func benchmarkNearby(noOfOps int, quadClient *quadgo.QuadrilleClient) int {
	locIDCh := make(chan string, 1000)
	go func() {
		for i := 0; i < noOfOps; i++ {
			locIDCh <- "loc" + strconv.Itoa(i)
		}
		close(locIDCh)
	}()
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for locID := range locIDCh {
				_ = locID
				lat, long := getRandomLocation(12.8623132, 13.0573792, 77.4743233, 77.7871233)
				_, err := quadClient.Nearby(lat, long, 100, 10)
				if err != nil {
					//fmt.Println(err)
				}
			}
		}()
	}
	wg.Wait()
	sec := time.Since(startTime).Seconds()
	return int(float64(noOfOps) / sec)
}

func cleanup(noOfOps int, quadClient *quadgo.QuadrilleClient) {
	locIDCh := make(chan string, 1000)
	go func() {
		for i := 0; i < noOfOps; i++ {
			locIDCh <- "loc" + strconv.Itoa(i)
		}
		close(locIDCh)
	}()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for locID := range locIDCh {
				err := quadClient.Delete(locID)
				if err != nil {
					//fmt.Println(err)
				}
			}
		}()
	}
	wg.Wait()
}
