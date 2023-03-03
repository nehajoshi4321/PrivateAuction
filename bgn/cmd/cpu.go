package main

import (
    "fmt"
    "os"
    "time"

    "github.com/shirou/gopsutil/cpu"
)

func main() {
    pid := os.Getpid()
    fmt.Printf("PID: %d\n", pid)

    start := time.Now()
    // Run some CPU-intensive code
    // ...
    fmt.Println("Hello")
    end := time.Now()

    cpuTime := end.Sub(start)
    fmt.Printf("CPU time: %s\n", cpuTime)

    cpuPercent, err := cpu.Percent(time.Second, false)
    if err != nil {
        fmt.Printf("Error getting CPU usage: %v", err)
        return
    }

    fmt.Printf("CPU utilization: %f%%\n", cpuPercent[0])
}

