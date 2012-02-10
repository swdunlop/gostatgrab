// Copyright 2012 Scott Dunlop. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gostatgrab 

/*  
#cgo pkg-config: libstatgrab
#include <statgrab.h>
static sg_disk_io_stats* disk_stats_ref(sg_disk_io_stats* ptr, int idx) {
   return &ptr[idx];
}
static sg_fs_stats* fs_stats_ref(sg_fs_stats* ptr, int idx) {
   return &ptr[idx];
}
static sg_process_stats* process_stats_ref(sg_process_stats* ptr, int idx) {
   return &ptr[idx];
}
*/
import "C"

import (
    "time"
)

func init() {
    if C.sg_init() != 0 {
        panic(getError())
    }
    if C.sg_drop_privileges() != 0 {
        panic(getError())
    }
    //int sg_snapshot();
    //int sg_shutdown();
}

func getError() error {
    return &Error{int(C.sg_get_error_errno())}
}

type Error struct {
    Errno int
}

func (err *Error) Error() string {
    return C.GoString(C.sg_str_error(C.sg_error(err.Errno)))
}

func Shutdown() error {
    if C.sg_shutdown() != 0 {
        return getError()
    }
    return nil
}

func GetCpuStats() (*CpuStats, error) {
    c := C.sg_get_cpu_stats()
    return &CpuStats{
        int64(c.user),
        int64(c.kernel),
        int64(c.idle),
        int64(c.iowait),
        int64(c.swap),
        int64(c.nice),
        int64(c.total),
        time.Unix(int64(c.systime), 0)}, nil
    // WARNING: modern unixes only; will horribly confuse your PDP.
}

type CpuStats struct {
    User    int64
    Kernel  int64
    Idle    int64
    Iowait  int64
    Swap    int64
    Nice    int64
    Total   int64
    Systime time.Time
}

//TODO: thread-safety by massive mutex

func GetCpuPercents() (*CpuPercents, error) {
    c := C.sg_get_cpu_percents()
    if c == nil {
        return nil, getError()
    }
    return &CpuPercents{
        float32(c.user),
        float32(c.kernel),
        float32(c.idle),
        float32(c.iowait),
        float32(c.swap),
        float32(c.nice),
        time.Unix(int64(c.time_taken), 0)}, nil
    // WARNING: modern unixes only; will horribly confuse your PDP.
}

type CpuPercents struct {
    User    float32
    Kernel  float32
    Idle    float32
    Iowait  float32
    Swap    float32
    Nice    float32
    Systime time.Time
}

func GetDiskIoStats() ([]*DiskIoStats, error) {
    var ct C.int
    c := C.sg_get_disk_io_stats(&ct)
    if c == nil {
        return nil, getError()
    }
    return asDiskIoStatsArray(c, int(ct)), nil
}

func GetDiskIoStatsDiff() ([]*DiskIoStats, error) {
    var ct C.int
    c := C.sg_get_disk_io_stats_diff(&ct)
    if c == nil {
        return nil, getError()
    }
    return asDiskIoStatsArray(c, int(ct)), nil
}

func asDiskIoStatsArray(d *C.sg_disk_io_stats, ct int) []*DiskIoStats {
    r := make([]*DiskIoStats, ct)
    for i := 0; i < ct; i++ {
        cc := C.disk_stats_ref(d, C.int(i))
        r[i] = &DiskIoStats{
            C.GoString(cc.disk_name),
            uint64(cc.read_bytes),
            uint64(cc.write_bytes),
            time.Unix(int64(cc.systime), 0)}
    }
    return r
}

type DiskIoStats struct {
    DiskName   string
    ReadBytes  uint64
    WriteBytes uint64
    Systime    time.Time
}

func GetMemStats() (*MemStats, error) {
    c := C.sg_get_mem_stats()
    if c == nil {
        return nil, getError()
    }
    return &MemStats{
        int64(c.total),
        int64(c.free),
        int64(c.used),
        int64(c.cache)}, nil
}

type MemStats struct {
    Total int64
    Free  int64
    Used  int64
    Cache int64
}

func GetSwapStats() (*SwapStats, error) {
    c := C.sg_get_swap_stats()
    if c == nil {
        return nil, getError()
    }
    return &SwapStats{
        int64(c.total),
        int64(c.free),
        int64(c.used)}, nil
}

type SwapStats struct {
    Total int64
    Free  int64
    Used  int64
}

func GetFsStats() ([]*FsStats, error) {
    var ct C.int
    c := C.sg_get_fs_stats(&ct)
    if c == nil {
        return nil, getError()
    }
    return asFsStatsArray(c, int(ct)), nil
}

func asFsStatsArray(d *C.sg_fs_stats, ct int) []*FsStats {
    r := make([]*FsStats, ct)
    for i := 0; i < ct; i++ {
        cc := C.fs_stats_ref(d, C.int(i))
        r[i] = &FsStats{
            C.GoString(cc.device_name),
            C.GoString(cc.fs_type),
            C.GoString(cc.mnt_point),

            int64(cc.size),
            int64(cc.used),
            int64(cc.avail),
            int64(cc.total_inodes),
            int64(cc.used_inodes),
            int64(cc.free_inodes),
            int64(cc.avail_inodes),
            int64(cc.io_size),
            int64(cc.block_size),
            int64(cc.total_blocks),
            int64(cc.free_blocks),
            int64(cc.used_blocks),
            int64(cc.avail_blocks)}
    }
    return r
}

type FsStats struct {
    DeviceName  string
    FsType      string
    MntPoint    string
    Size        int64
    Used        int64
    Avail       int64
    TotalInodes int64
    UsedInodes  int64
    FreeInodes  int64
    AvailInodes int64
    IoSize      int64
    BlockSize   int64
    TotalBlocks int64
    FreeBlocks  int64
    UsedBlocks  int64
    AvailBlocks int64
}

func GetHostInfo() (*HostInfo, error) {
    c := C.sg_get_host_info()
    if c == nil {
        return nil, getError()
    }
    return &HostInfo{
        C.GoString(c.os_name),
        C.GoString(c.os_release),
        C.GoString(c.os_version),
        C.GoString(c.platform),
        C.GoString(c.hostname),
        time.Duration(c.uptime) * time.Second}, nil
}

type HostInfo struct {
    OsName    string
    OsRelease string
    OsVersion string
    Platform  string
    Hostname  string
    Uptime    time.Duration
}

func GetProcessStats() ([]*ProcessStats, error) {
    var ct C.int
    c := C.sg_get_process_stats(&ct)
    if c == nil {
        return nil, getError()
    }
    return asProcessStatsArray(c, int(ct)), nil
}

func asProcessStatsArray(d *C.sg_process_stats, ct int) []*ProcessStats {
    r := make([]*ProcessStats, ct)
    for i := 0; i < ct; i++ {
        cc := C.process_stats_ref(d, C.int(i))
        r[i] = &ProcessStats{
            C.GoString(cc.process_name),
            C.GoString(cc.proctitle),

            int64(cc.pid),
            int64(cc.parent),
            int64(cc.pgid),
            int64(cc.uid),
            int64(cc.euid),
            int64(cc.gid),
            int64(cc.egid),
            uint64(cc.proc_size),
            uint64(cc.proc_resident),
            time.Duration(cc.time_spent) * time.Second,
            float32(cc.cpu_percent),
            int(cc.nice),
            ProcessState(cc.state)}
    }
    return r
}

type ProcessStats struct {
    ProcessName  string
    ProcessTitle string
    Pid          int64
    Parent       int64
    Pgid         int64
    Uid          int64
    Euid         int64
    Gid          int64
    Egid         int64
    ProcSize     uint64
    ProcResident uint64
    TimeSpent    time.Duration
    CpuPercent   float32
    Nice         int
    State        ProcessState
}

type ProcessState int

var (
    Running      = ProcessState(0)
    Sleeping     = ProcessState(1)
    Stopped      = ProcessState(2)
    Zombie       = ProcessState(3)
    UnknownState = ProcessState(4)
)

func GetProcessCount() (*ProcessCount, error) {
    c := C.sg_get_process_count()
    if c == nil {
        return nil, getError()
    }
    return &ProcessCount{
        int(c.total),
        int(c.running),
        int(c.sleeping),
        int(c.stopped),
        int(c.zombie)}, nil
}

type ProcessCount struct {
    Total    int
    Running  int
    Sleeping int
    Stopped  int
    Zombie   int
}

//TODO sg_get_load_stats
//TODO sg_get_page_stats
//TODO sg_get_user_stats
//TODO sg_get_network_io_stats
//TODO sg_get_network_iface_stats
