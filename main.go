package main

import (
    "bufio"
    "log"
    "os"
    "os/exec"

    . "gaze/common"
    "gaze/riemann"
)

func getFileScanner(fpath string) *bufio.Scanner {
    cmd := exec.Command("tail", "-n", "0", "-F", fpath)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Printf("ERROR : Failed to create command %v : %s\n", cmd, err.Error())
        os.Exit(1)
    }

    scanner := bufio.NewScanner(stdout)
    err = cmd.Start()
    if err != nil {
        log.Printf("ERROR : Failed to start command %v : %s\n", cmd, err.Error())
        os.Exit(1)
    }

    return scanner
}

func main() {
    cfg := GetConfig()

    var scanner *bufio.Scanner
    if cfg.Logfpath == "-" {
        scanner = bufio.NewScanner(os.Stdin)
    } else {
        scanner = getFileScanner(cfg.Logfpath)
    }

    riemann.StartAll()
    for scanner.Scan() {
        line := scanner.Text()
        riemann.SendToAll(line)
    }
    riemann.StopAll()
    riemann.WaitAll()
}
