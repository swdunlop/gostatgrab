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

// Shutdown ensures that any descriptors opened when gostatgrab was imported 
// are properly closed.
func Shutdown() error {
    if C.sg_shutdown() != 0 {
        return getError()
    }
    return nil
}

// GetCpuStats extracts the current CPU info from the system as ticks, see sg_get_cpu_stats(3).
func GetCpuStats() (*CpuStats, error) {
    c := C.sg_get_cpu_stats()
    if c == nil {
        return nil, error
    }
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
    User    int64     // number of ticks spent in user state
    Kernel  int64     // number of ticks spent in kernel state
    Idle    int64     // number of ticks spent in idle state
    Iowait  int64     // number of ticks spent in iowait state
    Swap    int64     // number of ticks spent in swap state
    Nice    int64     // number of ticks spent in nice state
    Total   int64     // total number of ticks spent
    Systime time.Time // current system time
}

// GetCpuPercents extracts the current CPU info from the system as percentages, see sg_get_cpu_percents(3).
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
        time.Duration(int64(c.time_taken)) * time.Second}, nil
    // WARNING: modern unixes only; will horribly confuse your PDP.
}

type CpuPercents struct {
    User    float32       // percentage of time spent in user state
    Kernel  float32       // percentage of time spent in kernel state
    Idle    float32       // percentage of time spent in idle state
    Iowait  float32       // percentage of time spent in iowait state
    Swap    float32       // percentage of time spent in swap state
    Nice    float32       // percentage of time spent in nice state
    Systime time.Duration // time since the last call to this API
}

// GetDiskIoStats returns information about I/O transfers since boot for all disks, see sg_get_disk_io_stats(3).
func GetDiskIoStats() ([]*DiskIoStats, error) {
    var ct C.int
    c := C.sg_get_disk_io_stats(&ct)
    if c == nil {
        return nil, getError()
    }
    return asDiskIoStatsArray(c, int(ct)), nil
}

// GetDiskIoStatsDiff returns information about I/O transfers the last call for all disks, see sg_get_disk_io_stats_diff(3).
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

// GetMemStats returns capacity and usage information about system memory, see sg_get_mem_stats(3).
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

// GetSwapStats returns capacity and usage information about system swap, see sg_get_swap_stats(3).
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

// GetFsStats returns information about mounted filesystems; see sg_get_fs_stats(3).
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
    DeviceName  string // the name of the device, as mounted
    FsType      string // the filesystem type (frequently misreported as "ext2" in Linux)
    MntPoint    string // where the filesystem is mounted
    Size        int64  // the size of the filesystem in bytes
    Used        int64  // the number of bytes used in the filesystem
    Avail       int64  // the number of bytes available for use
    TotalInodes int64  // the number of inodes in the filesystem
    UsedInodes  int64  // the number of inodes used by the filesystem
    FreeInodes  int64  // the number of free inodes, may be different from avail.
    AvailInodes int64  // the number of inodes available for use
    IoSize      int64  // the optimal size of a block for I/O
    BlockSize   int64  // the size of a block in the filesystem
    TotalBlocks int64  // the total number of blocks in the filesystem
    FreeBlocks  int64  // the number of blocks unused in the filesystem
    UsedBlocks  int64  // the number of blocks used in the filesystem
    AvailBlocks int64  // the number of blocks available to the filesystem
}

// GetHostInfo returns host identifying information; see sg_get_host_info(3).
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
    OsName    string        // equivalent to "uname -s"
    OsRelease string        // equivalent to "uname -r"
    OsVersion string        // equivalent to "uname -v"
    Platform  string        // equivalent to "uname -m"
    Hostname  string        // equivalent to "hostname"
    Uptime    time.Duration // time since last boot
}

// GetProcessStats returns process information similar to that found in ps(1); see sg_get_process_stats.
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
    ProcessName  string        // the name of the command
    ProcessTitle string        // the cmdline of the command (process-controlled)
    Pid          int64         // the process id
    Parent       int64         // the parent process id
    Pgid         int64         // the process group leader id
    Uid          int64         // the user id
    Euid         int64         // the effective user id
    Gid          int64         // the group id
    Egid         int64         // the efffective group id
    ProcSize     uint64        // the process size in bytes
    ProcResident uint64        // the size of the process in memory in bytes
    TimeSpent    time.Duration // time spent in running state
    CpuPercent   float32       // current cpu percentage utilized by the process
    Nice         int           // the "niceness" of the process, see nice(1) and/or nice(2)
    State        ProcessState  // the state of the process, see ProcessState
}

type ProcessState int

var (
    Running      = ProcessState(0) // the process is currently running
    Sleeping     = ProcessState(1) // the process is sleeping, waiting for an event
    Stopped      = ProcessState(2) // the process was stopped and must be resumed
    Zombie       = ProcessState(3) // the process is defunct and needs parent cleanup
    UnknownState = ProcessState(4) // the process state is not understood by statgrab
)

// GetProcessCount returns summary counts of processes in various states; see sg_get_process_count(3)
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
    Total    int // total number of processes in the system
    Running  int // number of processes that are Running
    Sleeping int // number of processes that are Sleeping
    Stopped  int // number of processes that are Stopped
    Zombie   int // number of processes that are Zombies
}

//TODO sg_get_load_stats
//TODO sg_get_page_stats
//TODO sg_get_user_stats
//TODO sg_get_network_io_stats
//TODO sg_get_network_iface_stats

//TODO: thread-safety by massive mutex
