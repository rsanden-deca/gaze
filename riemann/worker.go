package riemann

import (
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/riemann/riemann-go-client"

    . "gaze/common"
)

type RiemannWorker struct {
    Id       string
    endpoint string
    jobChan  chan Job
    client   *riemanngo.TCPClient
    wg       *sync.WaitGroup
}

func NewRiemannWorker(riemannEndpoint string, jobChan chan Job, number uint, wg *sync.WaitGroup) *RiemannWorker {
    worker := &RiemannWorker{
        Id:       fmt.Sprintf("%s-%03d", riemannEndpoint, number),
        endpoint: riemannEndpoint,
        jobChan:  jobChan,
        client:   nil,
        wg:       wg,
    }
    return worker
}

func (self *RiemannWorker) run() {
    defer self.wg.Done()
    for job := range self.jobChan {
        self.ProcessJob(job)
    }
}

func (self *RiemannWorker) Start() {
    self.wg.Add(1)
    go self.run()
}

func (self *RiemannWorker) Connect() (*riemanngo.TCPClient, error) {
    if self.client == nil {
        cfg := GetConfig()
        c := riemanngo.NewTCPClient(self.endpoint, cfg.RiemannConnectTimeout)
        err := c.Connect()
        if err != nil { return nil, err }
        self.client = c
    }
    return self.client, nil
}

func (self *RiemannWorker) Disconnect() {
    if self.client != nil {
        self.client.Close()
        self.client = nil
    }
}

func (self *RiemannWorker) SendEvent(event *riemanngo.Event, attempt int) error {
    c, err := self.Connect()
    if err != nil {
        log.Printf("WARN  : [%s] Failed to Connect to %s (attempt %d): %s\n", self.Id, self.endpoint, attempt, err.Error())
        self.Disconnect()
        return err
    }

    result, err := riemanngo.SendEvent(c, event)
    if err != nil {
        log.Printf("WARN  : [%s] Failed to send event to %s (attempt %d): %s (result: %v)\n", self.Id, self.endpoint, attempt, err.Error(), result)
        self.Disconnect()
        return err
    }

    return nil
}

func (self *RiemannWorker) SendEventRelentlessly(event *riemanngo.Event) error {
    cfg := GetConfig()
    var err error

    attempt := 1
    for attempt <= int(cfg.RiemannRetryAttempts) {
        err = self.SendEvent(event, attempt)
        if err == nil { return nil }
        time.Sleep(cfg.RiemannRetryDelay)
        attempt += 1
    }
    log.Printf("ERROR : [%s] Failed to send event after %d attempts : %v\n", self.Id, cfg.RiemannRetryAttempts, event)
    self.Disconnect()
    return err
}
