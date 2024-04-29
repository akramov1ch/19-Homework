package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func taskExecutor(taskID int, taskGroup *sync.WaitGroup, resultPipe chan<- int) {
    defer taskGroup.Done()
    delay := time.Duration(rand.Intn(3)) * time.Second
    time.Sleep(delay)
    outcome := rand.Intn(100)
    resultPipe <- outcome
    fmt.Printf("Task %d: %d\n", taskID, outcome)
}

func pipelineMerger(workerCount int, resultPipe chan int) []int {
    var taskGroup sync.WaitGroup
    taskGroup.Add(workerCount)

    finalOutcomes := make([]int, 0, workerCount)
    intermediatePipe := make(chan int, workerCount)

    for i := 0; i < workerCount; i++ {
        go taskExecutor(i, &taskGroup, intermediatePipe)
    }

    go func() {
        taskGroup.Wait()
        close(intermediatePipe)
    }()

    for {
        select {
        case result, open := <-intermediatePipe:
            if !open {
                return finalOutcomes
            }
            finalOutcomes = append(finalOutcomes, result)
        }
    }
}

func main() {
    resultPipe := make(chan int, 10)
    finalOutcomes := pipelineMerger(10, resultPipe)
    fmt.Println("Final outcomes:", finalOutcomes)
}
