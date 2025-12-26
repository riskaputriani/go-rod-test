# Go Rod Testing Browser Restrict

Aplikasi Go untuk testing browser automation menggunakan Rod dengan dukungan **Chrome for Testing** (portable, tidak perlu sudo/apt).

## Fitur Utama

✅ **Auto Version Check** - Cek versi, skip download jika sudah ada  
✅ **Zero Sudo** - Tidak perlu root/admin access  
✅ **Pure Go Extraction** - Tidak perlu tar/xz command  
✅ **Portable Chrome** - Chrome for Testing dari Google  
✅ **Smart Caching** - Download sekali, reuse selamanya

## Struktur Proyek

```
.
├── main.go                          # Entry point aplikasi
├── internal/
│   ├── browser/
│   │   └── chromium.go             # Manajemen Ungoogled Chromium (download & setup)
│   ├── logger/
│   │   └── logger.go               # Sistem logging
│   └── runtime/
│       └── info.go                 # Runtime information logging
├── go.mod
├── Makefile
└── build_instructions.txt
```

## Fitur

### 1. Modular Architecture

- **internal/browser**: Mengelola browser Chromium, termasuk download otomatis Ungoogled Chromium
- **internal/logger**: Sistem logging yang fleksibel dengan output ke file dan stdout
- **internal/runtime**: Logging informasi runtime sistem yang lengkap

### 2. Chrome for Testing with Smart Version Checking

Aplikasi akan:

- **Cek versi yang terinstall** - Bandingkan dengan versi yang dibutuhkan
- **Skip download jika sama** - Hemat bandwidth dan waktu
- **Auto download jika berbeda/belum ada** - Download Chrome for Testing (~150MB)
- **Pure Go extraction** - Extract ZIP tanpa dependency eksternal
- **Save version info** - Simpan file `.version` untuk tracking

**Keuntungan Chrome for Testing:**

- ✅ Portable (tidak perlu install ke sistem)
- ✅ Tidak perlu sudo/apt
- ✅ Include dependencies minimal
- ✅ Official dari Google
- ✅ Stable dan terupdate

### 3. Comprehensive Runtime Logging

- Informasi Go runtime (versi, compiler, arch)
- Memory statistics
- Build information dan dependencies
- System information (hostname, PID, user)
- Environment variables
- Linux-specific info (OS release, kernel, cgroup)

## Instalasi

### Prerequisites

- Go 1.25.4 atau lebih baru
- Linux OS (untuk Chrome for Testing)
- **✅ TIDAK PERLU** sudo/apt/yum atau package manager
- **✅ TIDAK PERLU** tar, xz, atau utility eksternal (menggunakan pure Go)
- **✅ TIDAK PERLU** Chrome/Chromium terinstall di sistem
- Minimal 250MB free space di `~/.local/share/` untuk Chrome

### Build

```bash
# Initialize module (jika belum)
go mod tidy

# Build untuk Linux
make build
```

## Penggunaan

```bash
# Jalankan aplikasi
./go-rod-testing-browser-restrict
```

### Perilaku Smart Version Check

**Run Pertama** (Chrome belum ada):

```
chrome_target_version: 131.0.6778.204
chrome_status: not_found_downloading
chrome_download: starting
chrome_extract: format_zip
chrome_extract: success
chrome_version: 131.0.6778.204
browser_status: connected
```

**Run Kedua dan Selanjutnya** (Chrome sudah ada dengan versi sama):

```
chrome_target_version: 131.0.6778.204
chrome_installed_version: 131.0.6778.204
chrome_status: already_installed_correct_version
browser_status: connected
(SKIP DOWNLOAD - instant start!)
```

**Jika Update Versi** (misalnya ganti ke v132):

```
chrome_target_version: 132.0.0.0
chrome_installed_version: 131.0.6778.204
chrome_version_mismatch: installed=131.0.6778.204, required=132.0.0.0
chrome_status: found_different_version_will_reinstall
(akan download versi baru)
```

## Environment Variables

- `RUNTIME_LOG_PATH`: Path custom untuk file log (opsional)

## Download Manual Chrome for Testing

Jika download otomatis gagal, Anda bisa download manual:

```bash
# Download Chrome for Testing
wget https://storage.googleapis.com/chrome-for-testing-public/131.0.6778.204/linux64/chrome-linux64.zip

# Ekstrak ke direktori yang benar
mkdir -p ~/.local/share/chrome-for-testing
unzip chrome-linux64.zip -d ~/.local/share/chrome-for-testing/

# Buat file version
echo "131.0.6778.204" > ~/.local/share/chrome-for-testing/.version

# Jalankan aplikasi
./go-rod-testing-browser-restrict
```

## Keuntungan Chrome for Testing

1. **Portable**: Tidak memerlukan instalasi sistem (no sudo needed)
2. **Official**: Langsung dari Google, bukan third-party
3. **Dependencies Minimal**: Sudah include library yang dibutuhkan
4. **Smart Caching**: Download sekali, reuse selamanya dengan version check
3. **Privacy**: Tanpa tracking Google
4. **Headless**: Cocok untuk automation
5. **No Sandbox**: Bisa jalan di environment terbatas

## Catatan

- Aplikasi **TIDAK** akan fallback ke browser default jika Ungoogled Chromium gagal disetup
- Ekstraksi memerlukan command `tar` yang sudah terinstall
- Browser dijalankan dalam mode `--headless` dan `--no-sandbox`
- Download file (~100MB) disimpan di direktori instalasi, bukan `/tmp` (untuk menghindari "no space" error di tmpfs)
- Pastikan ada minimal 200MB free space di `~/.local/share` untuk instalasi

## Troubleshooting

### Error: "no space left on device"

Jika mendapat error ini padahal disk masih banyak space:

- `/tmp` kemungkinan mounted sebagai tmpfs (RAM) dengan size terbatas
- Solusi sudah diterapkan: file didownload ke `~/.local/share/ungoogled-chromium` bukan `/tmp`
- Pastikan home directory (`~`) punya minimal 200MB free space

## License

MIT
