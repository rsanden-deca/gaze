package riemann

import (
    "encoding/json"
    "log"
    "time"

    "github.com/riemann/riemann-go-client"

    . "gaze/common"
)

type TagMessage struct {
    Raw string `json:"raw"`
}

type Job struct {
    Timestamp time.Time
    Line      string
}

func (self Job) ToEvent() (*riemanngo.Event, error) {
    cfg := GetConfig()

    msg := &TagMessage{self.Line}
    msgBytes, err := json.Marshal(msg)
    if err != nil {
        log.Printf("ERROR : Failed to marshal line : %s : %s\n", self.Line, err.Error())
        return nil, err
    }
    msgJson := string(msgBytes)
    tags := []string{msgJson}

    event := &riemanngo.Event{
        Time:    self.Timestamp,
        Host:    cfg.Hostname,
        Service: cfg.Service,
        TTL:     cfg.TTL,
        Tags:    tags,
    }
    return event, nil
}

func (self *RiemannWorker) ProcessJob(job Job) {
    event, err := job.ToEvent()
    if err != nil { return }

    log.Printf("INFO  : [%s] Sending event to %s : %v\n", self.Id, self.endpoint, event)
    err = self.SendEventRelentlessly(event)
    if err == nil {
        log.Printf("INFO  : [%s] Successfully processed event to %s: %v\n", self.Id, self.endpoint, event)
    } else {
        log.Printf("ERROR : [%s] Unsuccessfully processed event to %s: %v\n", self.Id, self.endpoint, event)
    }
}
