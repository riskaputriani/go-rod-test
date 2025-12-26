# Go Rod Testing Browser Restrict

Aplikasi Go untuk testing browser automation menggunakan Rod dengan dukungan Ungoogled Chromium Portable.

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

### 2. Automatic Ungoogled Chromium Setup

Aplikasi akan otomatis:

- Mengecek apakah Ungoogled Chromium sudah terinstall
- Download versi portable jika belum ada (~100MB, satu kali saja)
- Ekstrak ke direktori `~/.local/share/ungoogled-chromium`
- Menggunakan binary tersebut untuk menjalankan browser

**Penting**: Aplikasi ini **TIDAK** menggunakan Chrome default atau auto-download dari Rod. Hanya menggunakan Ungoogled Chromium yang dikelola oleh aplikasi sendiri. Jika Chromium gagal disetup, aplikasi akan error (tidak ada fallback ke browser lain).

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
- Linux OS (untuk Ungoogled Chromium)
- **TIDAK PERLU** sudo/apt/yum atau package manager
- **TIDAK PERLU** tar, xz, atau utility eksternal lainnya (menggunakan pure Go)

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

Aplikasi akan:

1. Membuat log file di `runtime-info.log` (atau sesuai `RUNTIME_LOG_PATH`)
2. Log semua informasi runtime
3. Setup Ungoogled Chromium (download jika perlu)
4. Buka browser dan navigasi ke example.com
5. Print title halaman

## Environment Variables

- `RUNTIME_LOG_PATH`: Path custom untuk file log (opsional)

## Download Manual Ungoogled Chromium

Jika download otomatis gagal, Anda bisa download manual:

```bash
# Download
wget https://github.com/macchrome/linchrome/releases/download/v142.7444.229-M142.0.7444.229-r1522585-portable-ungoogled-Lin64/ungoogled-chromium_142.0.7444.229_1.vaapi_linux.tar.xz

# Ekstrak ke direktori yang benar
mkdir -p ~/.local/share/ungoogled-chromium
tar -xf ungoogled-chromium_*.tar.xz -C ~/.local/share/ungoogled-chromium

# Jalankan aplikasi
./go-rod-testing-browser-restrict
```

## Keuntungan Ungoogled Chromium

1. **Portable**: Tidak memerlukan instalasi sistem
2. **Kompatibilitas**: Lebih baik di server terbatas
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
