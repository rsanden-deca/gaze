package riemann

import (
    "time"

    . "gaze/common"
)

var workerPools []*RiemannWorkerPool

func StartAll() {
    cfg := GetConfig()
    workerPools = []*RiemannWorkerPool{}
    for _, riemannEndpoint := range cfg.RiemannEndpoints {
        jobChan := make(chan Job, QUEUE_LENGTH)
        pool := NewRiemannWorkerPool(riemannEndpoint, jobChan, POOL_SIZE)
        pool.Start()
        workerPools = append(workerPools, pool)
    }
}

func StopAll() {
    for _, pool := range workerPools {
        pool.Stop()
    }
}

func WaitAll() {
    for _, pool := range workerPools {
        pool.Wait()
    }
}

func SendToAll(line string) {
    timestamp := time.Now().UTC()
    for _, pool := range workerPools {
        pool.Submit(line, timestamp)
    }
}
