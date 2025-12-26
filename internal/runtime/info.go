package runtime

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

// LogFunc adalah tipe fungsi untuk logging
type LogFunc func(key, value string)

// Info berisi informasi runtime
type Info struct {
	logger LogFunc
}

// NewInfo membuat instance baru Info
func NewInfo(logger LogFunc) *Info {
	return &Info{
		logger: logger,
	}
}

// LogAll mencatat semua informasi runtime
func (ri *Info) LogAll() {
	ri.LogBasicInfo()
	ri.LogGoInfo()
	ri.LogMemoryInfo()
	ri.LogBuildInfo()
	ri.LogSystemInfo()
	ri.LogProcessInfo()
	ri.LogUserInfo()
	ri.LogEnvironment()
	ri.LogLinuxSpecific()
}

// LogBasicInfo mencatat informasi dasar
func (ri *Info) LogBasicInfo() {
	ri.logger("start_time", time.Now().Format(time.RFC3339))
	ri.logger("args", sanitizeValue(strings.Join(os.Args, " ")))
}

// LogGoInfo mencatat informasi Go runtime
func (ri *Info) LogGoInfo() {
	ri.logger("go_version", runtime.Version())
	ri.logger("go_compiler", runtime.Compiler)
	ri.logger("go_os_arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	ri.logger("go_root", runtime.GOROOT())
	ri.logger("go_maxprocs", fmt.Sprintf("%d", runtime.GOMAXPROCS(0)))
	ri.logger("go_numcpu", fmt.Sprintf("%d", runtime.NumCPU()))
}

// LogMemoryInfo mencatat informasi memori
func (ri *Info) LogMemoryInfo() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	ri.logger("mem_alloc_bytes", fmt.Sprintf("%d", mem.Alloc))
	ri.logger("mem_total_alloc_bytes", fmt.Sprintf("%d", mem.TotalAlloc))
	ri.logger("mem_sys_bytes", fmt.Sprintf("%d", mem.Sys))
	ri.logger("mem_heap_alloc_bytes", fmt.Sprintf("%d", mem.HeapAlloc))
	ri.logger("mem_heap_sys_bytes", fmt.Sprintf("%d", mem.HeapSys))
	ri.logger("mem_heap_inuse_bytes", fmt.Sprintf("%d", mem.HeapInuse))
	ri.logger("mem_stack_inuse_bytes", fmt.Sprintf("%d", mem.StackInuse))
	ri.logger("mem_num_gc", fmt.Sprintf("%d", mem.NumGC))
	ri.logger("mem_pause_total_ns", fmt.Sprintf("%d", mem.PauseTotalNs))
}

// LogBuildInfo mencatat informasi build
func (ri *Info) LogBuildInfo() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Path != "" {
			ri.logger("module_path", info.Main.Path)
		}
		if info.Main.Version != "" {
			ri.logger("module_version", info.Main.Version)
		}
		for _, setting := range info.Settings {
			ri.logger("build_"+setting.Key, setting.Value)
		}
		if len(info.Deps) > 0 {
			for _, dep := range info.Deps {
				ver := dep.Version
				if dep.Replace != nil {
					ver = dep.Replace.Version
				}
				ri.logger("dep", fmt.Sprintf("%s@%s", dep.Path, ver))
			}
		}
	}
}

// LogSystemInfo mencatat informasi sistem
func (ri *Info) LogSystemInfo() {
	if hostname, err := os.Hostname(); err == nil {
		ri.logger("hostname", hostname)
	} else {
		ri.logger("hostname_error", err.Error())
	}

	ri.logger("temp_dir", os.TempDir())
}

// LogProcessInfo mencatat informasi proses
func (ri *Info) LogProcessInfo() {
	ri.logger("pid", fmt.Sprintf("%d", os.Getpid()))
	ri.logger("ppid", fmt.Sprintf("%d", os.Getppid()))

	if exe, err := os.Executable(); err == nil {
		ri.logger("executable", exe)
	} else {
		ri.logger("executable_error", err.Error())
	}

	if cwd, err := os.Getwd(); err == nil {
		ri.logger("working_dir", cwd)
	} else {
		ri.logger("working_dir_error", err.Error())
	}
}

// LogUserInfo mencatat informasi user
func (ri *Info) LogUserInfo() {
	if u, err := user.Current(); err == nil {
		if u.Username != "" {
			ri.logger("user_name", u.Username)
		}
		if u.Uid != "" {
			ri.logger("user_uid", u.Uid)
		}
		if u.Gid != "" {
			ri.logger("user_gid", u.Gid)
		}
		if u.HomeDir != "" {
			ri.logger("user_home", u.HomeDir)
		}
	} else {
		ri.logger("user_error", err.Error())
	}
}

// LogEnvironment mencatat environment variables
func (ri *Info) LogEnvironment() {
	if env := os.Environ(); len(env) > 0 {
		sort.Strings(env)
		for _, entry := range env {
			ri.logger("env", sanitizeValue(entry))
		}
	}
}

// LogLinuxSpecific mencatat informasi khusus Linux
func (ri *Info) LogLinuxSpecific() {
	if runtime.GOOS != "linux" {
		return
	}

	// OS Release
	if osrel := readOSRelease(); osrel != nil {
		keys := make([]string, 0, len(osrel))
		for key := range osrel {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			ri.logger("os_release", fmt.Sprintf("%s=%s", key, osrel[key]))
		}
	}

	// Kernel info
	if val, err := readFirstLine("/proc/sys/kernel/osrelease"); err == nil {
		ri.logger("kernel_release", val)
	}
	if val, err := readFirstLine("/proc/version"); err == nil {
		ri.logger("kernel_version", val)
	}

	// Cgroup info
	if cgroup, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		lines := strings.Split(strings.TrimSpace(string(cgroup)), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			ri.logger("cgroup", sanitizeValue(line))
		}
	}
}

// Helper functions

func sanitizeValue(value string) string {
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\n", "\\n")
	return value
}

func readFirstLine(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	line := strings.TrimSpace(string(data))
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}
	return line, nil
}

func readOSRelease() map[string]string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil
	}
	result := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		result[key] = val
	}
	return result
}
