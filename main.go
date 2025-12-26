package main

import (
    "fmt"
    "io"
    "os"
    "os/user"
    "path/filepath"
    "runtime"
    "runtime/debug"
    "sort"
    "strings"
    "time"

    "github.com/go-rod/rod"
    "github.com/go-rod/rod/lib/launcher"
)

var logWriter io.Writer = os.Stdout

func logKV(key, value string) {
    fmt.Fprintf(logWriter, "%s: %s\n", key, value)
}

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

func openLogFile() (*os.File, string, error) {
    if path := strings.TrimSpace(os.Getenv("RUNTIME_LOG_PATH")); path != "" {
        if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
            return nil, "", err
        }
        file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
        return file, path, err
    }

    if cwd, err := os.Getwd(); err == nil {
        path := filepath.Join(cwd, "runtime-info.log")
        file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
        if err == nil {
            return file, path, nil
        }
    }

    path := filepath.Join(os.TempDir(), "runtime-info.log")
    file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
    return file, path, err
}

func logRuntimeInfo() {
    logKV("start_time", time.Now().Format(time.RFC3339))
    logKV("args", sanitizeValue(strings.Join(os.Args, " ")))
    logKV("go_version", runtime.Version())
    logKV("go_compiler", runtime.Compiler)
    logKV("go_os_arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
    logKV("go_root", runtime.GOROOT())
    logKV("go_maxprocs", fmt.Sprintf("%d", runtime.GOMAXPROCS(0)))
    logKV("go_numcpu", fmt.Sprintf("%d", runtime.NumCPU()))

    var mem runtime.MemStats
    runtime.ReadMemStats(&mem)
    logKV("mem_alloc_bytes", fmt.Sprintf("%d", mem.Alloc))
    logKV("mem_total_alloc_bytes", fmt.Sprintf("%d", mem.TotalAlloc))
    logKV("mem_sys_bytes", fmt.Sprintf("%d", mem.Sys))
    logKV("mem_heap_alloc_bytes", fmt.Sprintf("%d", mem.HeapAlloc))
    logKV("mem_heap_sys_bytes", fmt.Sprintf("%d", mem.HeapSys))
    logKV("mem_heap_inuse_bytes", fmt.Sprintf("%d", mem.HeapInuse))
    logKV("mem_stack_inuse_bytes", fmt.Sprintf("%d", mem.StackInuse))
    logKV("mem_num_gc", fmt.Sprintf("%d", mem.NumGC))
    logKV("mem_pause_total_ns", fmt.Sprintf("%d", mem.PauseTotalNs))

    if info, ok := debug.ReadBuildInfo(); ok {
        if info.Main.Path != "" {
            logKV("module_path", info.Main.Path)
        }
        if info.Main.Version != "" {
            logKV("module_version", info.Main.Version)
        }
        for _, setting := range info.Settings {
            logKV("build_"+setting.Key, setting.Value)
        }
        if len(info.Deps) > 0 {
            for _, dep := range info.Deps {
                ver := dep.Version
                if dep.Replace != nil {
                    ver = dep.Replace.Version
                }
                logKV("dep", fmt.Sprintf("%s@%s", dep.Path, ver))
            }
        }
    }

    if hostname, err := os.Hostname(); err == nil {
        logKV("hostname", hostname)
    } else {
        logKV("hostname_error", err.Error())
    }

    logKV("pid", fmt.Sprintf("%d", os.Getpid()))
    logKV("ppid", fmt.Sprintf("%d", os.Getppid()))

    if exe, err := os.Executable(); err == nil {
        if abs, absErr := filepath.Abs(exe); absErr == nil {
            exe = abs
        }
        logKV("executable", exe)
    } else {
        logKV("executable_error", err.Error())
    }

    if cwd, err := os.Getwd(); err == nil {
        logKV("working_dir", cwd)
    } else {
        logKV("working_dir_error", err.Error())
    }

    logKV("temp_dir", os.TempDir())

    if u, err := user.Current(); err == nil {
        if u.Username != "" {
            logKV("user_name", u.Username)
        }
        if u.Uid != "" {
            logKV("user_uid", u.Uid)
        }
        if u.Gid != "" {
            logKV("user_gid", u.Gid)
        }
        if u.HomeDir != "" {
            logKV("user_home", u.HomeDir)
        }
    } else {
        logKV("user_error", err.Error())
    }

    if env := os.Environ(); len(env) > 0 {
        sort.Strings(env)
        for _, entry := range env {
            logKV("env", sanitizeValue(entry))
        }
    }

    if runtime.GOOS == "linux" {
        if osrel := readOSRelease(); osrel != nil {
            keys := make([]string, 0, len(osrel))
            for key := range osrel {
                keys = append(keys, key)
            }
            sort.Strings(keys)
            for _, key := range keys {
                logKV("os_release", fmt.Sprintf("%s=%s", key, osrel[key]))
            }
        }
        if val, err := readFirstLine("/proc/sys/kernel/osrelease"); err == nil {
            logKV("kernel_release", val)
        }
        if val, err := readFirstLine("/proc/version"); err == nil {
            logKV("kernel_version", val)
        }
        if cgroup, err := os.ReadFile("/proc/1/cgroup"); err == nil {
            lines := strings.Split(strings.TrimSpace(string(cgroup)), "\n")
            for _, line := range lines {
                if line == "" {
                    continue
                }
                logKV("cgroup", sanitizeValue(line))
            }
        }
    }
}

func main() {
    logFile, logPath, err := openLogFile()
    if err == nil {
        logWriter = io.MultiWriter(os.Stdout, logFile)
        defer logFile.Close()
        logKV("log_file", logPath)
    } else {
        logWriter = os.Stdout
        logKV("log_file_error", err.Error())
    }

    logRuntimeInfo()

    // Rod akan otomatis download browser ke user directory
    path, _ := launcher.LookPath()
    u := launcher.New().Bin(path).MustLaunch()
    
    browser := rod.New().ControlURL(u).MustConnect()
    defer browser.MustClose()
    
    // Gunakan browser
    page := browser.MustPage("https://example.com")
    println(page.MustInfo().Title)
}
