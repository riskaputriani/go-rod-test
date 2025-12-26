package main

import (
    "github.com/go-rod/rod"
    "github.com/go-rod/rod/lib/launcher"
)

func main() {
    // Rod akan otomatis download browser ke user directory
    path, _ := launcher.LookPath()
    u := launcher.New().Bin(path).MustLaunch()
    
    browser := rod.New().ControlURL(u).MustConnect()
    defer browser.MustClose()
    
    // Gunakan browser
    page := browser.MustPage("https://example.com")
    println(page.MustInfo().Title)
}
