package main

import (
	"fmt"

	"go-rod-testing-browser-restrict/internal/browser"
	"go-rod-testing-browser-restrict/internal/logger"
	"go-rod-testing-browser-restrict/internal/runtime"
)

func main() {
	// Inisialisasi logger
	log, err := logger.New()
	if err == nil && log.GetPath() != "" {
		log.LogKV("log_file", log.GetPath())
	} else if err != nil {
		fmt.Printf("log_file_error: %s\n", err.Error())
	}

	// Log semua informasi runtime
	runtimeInfo := runtime.NewInfo(log.LogKV)
	runtimeInfo.LogAll()

	// Setup dan jalankan browser dengan Ungoogled Chromium
	chromiumMgr := browser.NewChromiumManager(log.LogKV)

	// Setup akan otomatis download jika belum ada
	if err := chromiumMgr.Setup(); err != nil {
		log.LogKV("chromium_setup_error", err.Error())
		// Lanjutkan dengan browser default jika gagal
	}

	// Dapatkan browser instance
	browserInstance, err := chromiumMgr.GetBrowser()
	if err != nil {
		log.LogKV("browser_error", err.Error())
		panic(err)
	}
	defer browserInstance.MustClose()

	// Gunakan browser
	page := browserInstance.MustPage("https://example.com")
	title := page.MustInfo().Title

	log.LogKV("page_title", title)
	fmt.Printf("\nPage Title: %s\n", title)

	log.LogKV("status", "success")
}
