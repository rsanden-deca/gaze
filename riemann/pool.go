package riemann

import (
    "log"
    "sync"
    "time"
)

type RiemannWorkerPool struct {
    endpoint string
    workers  []*RiemannWorker
    jobChan  chan Job
    wg       *sync.WaitGroup
}

func NewRiemannWorkerPool(riemannEndpoint string, jobChan chan Job, size uint) *RiemannWorkerPool {
    var wg sync.WaitGroup

    workers := []*RiemannWorker{}
    for i := uint(0); i < size; i++ {
        worker := NewRiemannWorker(riemannEndpoint, jobChan, i, &wg)
        workers = append(workers, worker)
    }

    pool := &RiemannWorkerPool{
        endpoint: riemannEndpoint,
        workers:  workers,
        jobChan:  jobChan,
        wg:       &wg,
    }
    return pool
}

func (self *RiemannWorkerPool) Start() {
    for _, worker := range self.workers {
        worker.Start()
    }
}

func (self *RiemannWorkerPool) Stop() {
    close(self.jobChan)
}

func (self *RiemannWorkerPool) Wait() {
    self.wg.Wait()
}

func (self *RiemannWorkerPool) Submit(line string, timestamp time.Time) {
    job := Job{timestamp, line}
    select {
    case self.jobChan <- job:
    default:
        log.Printf("WARN  : Channel for riemann server %s full! Discarding line: %s\n", self.endpoint, line)
    }
}
