package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup, results chan<- int) {
    defer wg.Done()
    time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
    result := rand.Intn(100)
    results <- result
    fmt.Printf("Worker %d: %d\n", id, result)
}

func fanIn(workers int, results chan int) []int {
    var wg sync.WaitGroup
    wg.Add(workers) 

    finalResults := make([]int, 0, workers)
    resultsChan := make(chan int, workers)

    for i := 0; i < workers; i++ {
        go worker(i, &wg, resultsChan)
    }

    go func() {
        wg.Wait()
        close(resultsChan)
    }()

    for {
        select {
        case result, ok := <-resultsChan:
            if !ok {
                return finalResults
            }
            finalResults = append(finalResults, result)
        
		}
    }
}

func main() {
    results := make(chan int, 10)
    finalResults := fanIn(10, results)
    fmt.Println("Final results:", finalResults)
}
