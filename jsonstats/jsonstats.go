package main

import (
    "encoding/json"
    "os"
    "github.com/swdunlop/gostatgrab"
)

var outp *json.Encoder = json.NewEncoder(os.Stdout)
var errp *json.Encoder = json.NewEncoder(os.Stderr)

func printResult(v interface{}, e error) {
    if e != nil {
        errp.Encode(e.Error())
    } else {
        outp.Encode(v)
    }
}

func main() {
    defer gostatgrab.Shutdown()
    defer os.Stdout.Sync()
    defer os.Stderr.Sync()
    printResult(gostatgrab.GetCpuPercents())
    printResult(gostatgrab.GetDiskIoStats())
    printResult(gostatgrab.GetMemStats())
    printResult(gostatgrab.GetSwapStats())
    printResult(gostatgrab.GetFsStats())
    printResult(gostatgrab.GetHostInfo())
    printResult(gostatgrab.GetProcessStats())
    printResult(gostatgrab.GetProcessCount())

}
