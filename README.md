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

### 2. Chrome for Testing with Smart Version Checking + Auto Dependencies

Aplikasi akan:

- **Auto-detect OS** - Deteksi Debian vs Ubuntu untuk download yang tepat
- **Cek dependencies** - Cek library yang sudah ada di sistem/libs
- **Download .deb packages** - Download dependencies yang belum ada
- **Extract tanpa sudo** - Ekstrak .deb ke folder lokal (dpkg -x atau pure Go)
- **Set LD_LIBRARY_PATH** - Arahkan Chrome ke libraries yang di-extract
- **Cek versi Chrome** - Bandingkan dengan versi yang dibutuhkan
- **Skip download jika sama** - Hemat bandwidth dan waktu
- **Auto download jika berbeda/belum ada** - Download Chrome for Testing (~150MB)
- **Pure Go extraction** - Extract ZIP/DEB tanpa dependency eksternal
- **Save version info** - Simpan file `.version` untuk tracking

**Dependencies yang di-handle otomatis:**

- ✅ libnss3 (Network Security Services)
- ✅ libnspr4 (Netscape Portable Runtime)
- ✅ libatk-bridge2.0-0 (Accessibility Toolkit Bridge)
- ✅ libatk1.0-0 (ATK Library)
- ✅ libatspi2.0-0 (Assistive Technology SPI)
- ✅ libcups2 (CUPS printing)
- ✅ libdrm2 (Direct Rendering Manager)
- ✅ libgbm1 (Generic Buffer Management)
- ✅ libasound2 (ALSA sound)

**Keuntungan Chrome for Testing + Auto Deps:**

- ✅ Portable (tidak perlu install ke sistem)
- ✅ **Zero Sudo** - Tidak perlu root access untuk dependencies
- ✅ **Auto Dependencies** - Otomatis download & extract shared libraries
- ✅ **Pure Go** - Ekstrak .deb dengan pure Go (fallback ke dpkg -x)
- ✅ **OS Detection** - Auto pilih Debian/Ubuntu packages
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
- **✅ TIDAK PERLU** sudo/apt/yum atau package manager (dependencies di-download otomatis)
- **✅ TIDAK PERLU** tar, xz, dpkg atau utility eksternal (menggunakan pure Go)
- **✅ TIDAK PERLU** Chrome/Chromium terinstall di sistem
- **✅ TIDAK PERLU** install shared libraries (.so files) secara manual
- Minimal 250MB free space di `~/.local/share/` untuk Chrome
- Minimal 50MB free space untuk dependencies (9 packages ~5MB each)

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

### Perilaku Smart Version Check + Auto Dependencies

**Run Pertama** (Chrome dan dependencies belum ada):

```
os_detected: debian
chrome_target_version: 131.0.6778.204
chrome_status: not_found_downloading
chrome_download: starting
chrome_extract: format_zip
chrome_extract: success

dependency_check: starting
dependency: libnss3 - not_found
dependency_download: libnss3
dependency_extract: libnss3 - success
dependency: libnspr4 - not_found
dependency_download: libnspr4
dependency_extract: libnspr4 - success
... (7 dependencies lainnya)
dependency_setup: complete

chrome_version: 131.0.6778.204
browser_status: connected
```

**Run Kedua dan Selanjutnya** (Chrome dan dependencies sudah ada):

```
chrome_target_version: 131.0.6778.204
chrome_installed_version: 131.0.6778.204
chrome_status: already_installed_correct_version

dependency_check: starting
dependency: libnss3 - already_exists
dependency: libnspr4 - already_exists
... (semua dependencies sudah ada)
dependency_setup: skipped_all_exists

browser_status: connected
(INSTANT START - no download!)
```

**Jika Update Versi Chrome** (misalnya ganti ke v132):

```
chrome_target_version: 132.0.0.0
chrome_installed_version: 131.0.6778.204
chrome_version_mismatch: installed=131.0.6778.204, required=132.0.0.0
chrome_status: found_different_version_will_reinstall
dependency_check: skip (dependencies sudah ada)
(akan download Chrome versi baru, reuse dependencies)
```

## Environment Variables

- `RUNTIME_LOG_PATH`: Path custom untuk file log (opsional)

## Download Manual (Opsional)

### Chrome for Testing

Jika download Chrome otomatis gagal:

```bash
# Download Chrome for Testing
wget https://storage.googleapis.com/chrome-for-testing-public/131.0.6778.204/linux64/chrome-linux64.zip

# Ekstrak ke direktori yang benar
mkdir -p ~/.local/share/chrome-for-testing
unzip chrome-linux64.zip -d ~/.local/share/chrome-for-testing/

# Buat file version
echo "131.0.6778.204" > ~/.local/share/chrome-for-testing/.version
```

### Dependencies Manual

Jika auto-download dependencies gagal, download manual:

```bash
# Untuk Debian 12
cd ~/.local/share/chrome-for-testing/libs
wget http://ftp.debian.org/debian/pool/main/n/nss/libnss3_2%3a3.87.1-1_amd64.deb
dpkg -x libnss3_2:3.87.1-1_amd64.deb .

# Untuk Ubuntu 22.04
cd ~/.local/share/chrome-for-testing/libs
wget http://archive.ubuntu.com/ubuntu/pool/main/n/nss/libnss3_2%3a3.68.2-0ubuntu1.2_amd64.deb
dpkg -x libnss3_2:3.68.2-0ubuntu1.2_amd64.deb .

# Jalankan aplikasi
./go-rod-testing-browser-restrict
```

## Keuntungan Chrome for Testing + Auto Dependencies

1. **Zero Sudo**: Tidak memerlukan root access sama sekali
2. **Portable**: Chrome + libraries dalam folder lokal
3. **Official**: Chrome dari Google, libraries dari repos resmi
4. **Smart Caching**: Download sekali Chrome + deps, reuse selamanya
5. **Auto Dependencies**: Otomatis resolve shared libraries
6. **OS Detection**: Auto pilih Debian/Ubuntu packages
7. **Pure Go**: Ekstrak .deb tanpa dpkg (fallback ke dpkg -x jika ada)
8. **Headless**: Cocok untuk automation
9. **No Sandbox**: Bisa jalan di environment terbatas

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
