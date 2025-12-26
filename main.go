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
	// TIDAK akan menggunakan chrome default atau auto-download dari rod
	// Hanya menggunakan Ungoogled Chromium yang didownload oleh aplikasi ini
	chromiumMgr := browser.NewChromiumManager(log.LogKV)

	// Dapatkan browser instance (akan otomatis download Chromium jika belum ada)
	browserInstance, err := chromiumMgr.GetBrowser()
	if err != nil {
		log.LogKV("browser_error", err.Error())
		fmt.Printf("\nError: %s\n", err.Error())
		fmt.Println("Chromium gagal disetup. Pastikan koneksi internet aktif.")
		panic(err)
	}
	defer browserInstance.MustClose()

	// Gunakan browser
	page := browserInstance.MustPage("https://example.com")

	// Tunggu sampai page selesai loading
	page.MustWaitStable()

	title := page.MustInfo().Title

	log.LogKV("page_title", title)
	log.LogKV("page_url", page.MustInfo().URL)
	fmt.Printf("\nPage Title: %s\n", title)
	fmt.Printf("Page URL: %s\n", page.MustInfo().URL)

	log.LogKV("status", "success")
}
