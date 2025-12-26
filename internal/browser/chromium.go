package browser

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ulikunitz/xz"
)

// ChromiumManager mengelola instalasi dan konfigurasi Chrome
type ChromiumManager struct {
	installDir  string
	execPath    string
	versionFile string
	logger      func(key, value string)
	config      Config
}

// NewChromiumManager membuat instance baru ChromiumManager
func NewChromiumManager(logger func(key, value string)) *ChromiumManager {
	return NewChromiumManagerWithConfig(DefaultConfig(), logger)
}

// Setup mengecek dan mengunduh Chrome jika belum ada atau versi berbeda
func (cm *ChromiumManager) Setup() error {
	cm.logger("chrome_install_dir", cm.installDir)
	cm.logger("chrome_target_version", cm.config.Version)

	// Cek apakah sudah terinstall dengan versi yang sama
	if cm.isInstalledWithCorrectVersion() {
		cm.logger("chrome_status", "already_installed_correct_version")
		cm.logger("chrome_version", cm.config.Version)
		return nil
	}

	// Cek apakah ada instalasi lama dengan versi berbeda
	if cm.isInstalled() {
		cm.logger("chrome_status", "found_different_version_will_reinstall")
	} else {
		cm.logger("chrome_status", "not_found_downloading")
	}

	// Download dan ekstrak
	if err := cm.downloadAndExtract(); err != nil {
		return fmt.Errorf("failed to setup chrome: %w", err)
	}

	// Simpan versi yang terinstall
	if err := cm.saveVersion(); err != nil {
		cm.logger("chrome_version_save_error", err.Error())
	}

	cm.logger("chrome_status", "installation_complete")
	return nil
}

// isInstalledWithCorrectVersion mengecek apakah Chrome sudah terinstall dengan versi yang benar
func (cm *ChromiumManager) isInstalledWithCorrectVersion() bool {
	// Cek executable ada
	if !cm.isInstalled() {
		return false
	}

	// Cek versi file
	versionFile := filepath.Join(cm.installDir, ".version")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		cm.logger("chrome_version_check", "version_file_not_found")
		return false
	}

	installedVersion := strings.TrimSpace(string(data))
	cm.logger("chrome_installed_version", installedVersion)

	if installedVersion != cm.config.Version {
		cm.logger("chrome_version_mismatch", fmt.Sprintf("installed=%s, required=%s", installedVersion, cm.config.Version))
		return false
	}

	return true
}

// saveVersion menyimpan informasi versi yang terinstall
func (cm *ChromiumManager) saveVersion() error {
	versionFile := filepath.Join(cm.installDir, ".version")
	return os.WriteFile(versionFile, []byte(cm.config.Version), 0644)
}

// isInstalled mengecek apakah Chrome sudah terinstall
func (cm *ChromiumManager) isInstalled() bool {
	// Cari executable chrome di direktori instalasi
	chromePath := filepath.Join(cm.installDir, "chrome")

	if _, err := os.Stat(chromePath); err == nil {
		cm.execPath = chromePath
		cm.logger("chrome_executable", chromePath)
		return true
	}

	// Cari di subdirektori (Chrome for Testing zip biasanya punya folder chrome-linux64)
	entries, err := os.ReadDir(cm.installDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			chromePath := filepath.Join(cm.installDir, entry.Name(), "chrome")
			if _, err := os.Stat(chromePath); err == nil {
				cm.execPath = chromePath
				cm.logger("chrome_executable", chromePath)
				return true
			}
		}
	}

	return false
}

// downloadAndExtract mengunduh dan mengekstrak Chrome
func (cm *ChromiumManager) downloadAndExtract() error {
	// Buat direktori instalasi
	if err := os.MkdirAll(cm.installDir, 0o755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Download file
	cm.logger("chrome_download", "starting")
	cm.logger("chrome_url", cm.config.DownloadURL)
	resp, err := http.Get(cm.config.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download chrome: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download chrome: status %d", resp.StatusCode)
	}

	cm.logger("chrome_download", "extracting")

	// Deteksi format berdasarkan URL
	if strings.HasSuffix(cm.config.DownloadURL, ".zip") {
		// Chrome for Testing menggunakan ZIP
		return cm.extractZip(resp.Body)
	} else if strings.HasSuffix(cm.config.DownloadURL, ".tar.xz") {
		// Ungoogled Chromium menggunakan TAR.XZ
		return cm.extractTarXz(resp.Body)
	}

	return fmt.Errorf("unsupported file format: %s", cm.config.DownloadURL)
}

// extractZip mengekstrak file ZIP (untuk Chrome for Testing)
func (cm *ChromiumManager) extractZip(r io.Reader) error {
	cm.logger("chrome_extract", "format_zip")

	// Simpan ke file sementara (ZIP memerlukan random access)
	tmpFile := filepath.Join(cm.installDir, ".chrome-download.zip")
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	cm.logger("chrome_temp_file", tmpFile)

	// Copy data
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(tmpFile)
		return fmt.Errorf("failed to save zip: %w", err)
	}
	f.Close()

	cm.logger("chrome_extract", "opening_zip")

	// Buka ZIP
	zipReader, err := zip.OpenReader(tmpFile)
	if err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer zipReader.Close()

	// Ekstrak semua file
	for _, file := range zipReader.File {
		target := filepath.Join(cm.installDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(target, 0o755)
			continue
		}

		// Buat parent directory
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			continue
		}

		// Ekstrak file
		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			continue
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			continue
		}
	}

	// Hapus file temporary
	os.Remove(tmpFile)
	cm.logger("chrome_extract", "success")

	return nil
}

// extractTarXz mengekstrak file tar.xz menggunakan pure Go (tanpa dependency eksternal)
func (cm *ChromiumManager) extractTarXz(r io.Reader) error {
	cm.logger("chrome_extract", "format_tar_xz")

	// Decompress xz
	xzReader, err := xz.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	// Extract tar
	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		target := filepath.Join(cm.installDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
			cm.logger("chromium_extract_dir", header.Name)

		case tar.TypeReg:
			// Create file
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}

			outFile.Close()

			// Log progress setiap 100 file
			if header.Name[len(header.Name)-1] == 'e' {
				cm.logger("chromium_extract_file", header.Name)
			}

		case tar.TypeSymlink:
			// Create symlink
			if err := os.Symlink(header.Linkname, target); err != nil {
				// Ignore jika symlink sudah ada
				if !os.IsExist(err) {
					return fmt.Errorf("failed to create symlink %s: %w", target, err)
				}
			}
		}
	}

	cm.logger("chromium_extract", "success")
	return nil
}

// GetBrowser membuat dan mengembalikan instance browser Rod
func (cm *ChromiumManager) GetBrowser() (*rod.Browser, error) {
	// Setup Chrome jika belum (akan download jika perlu atau skip jika sudah ada dengan versi sama)
	if err := cm.Setup(); err != nil {
		return nil, fmt.Errorf("failed to setup chrome: %w", err)
	}

	// Pastikan execPath sudah ditemukan
	if cm.execPath == "" {
		return nil, fmt.Errorf("chrome executable not found after setup")
	}

	// Gunakan Chrome yang sudah didownload
	cm.logger("browser_using", "chrome_for_testing")
	cm.logger("browser_executable", cm.execPath)
	cm.logger("browser_version", cm.config.Version)

	// Set environment variable untuk library dependencies jika ada
	libPath := filepath.Join(cm.installDir, "lib")
	if _, err := os.Stat(libPath); err == nil {
		// Jika ada folder lib, tambahkan ke LD_LIBRARY_PATH
		currentLD := os.Getenv("LD_LIBRARY_PATH")
		if currentLD != "" {
			os.Setenv("LD_LIBRARY_PATH", libPath+":"+currentLD)
		} else {
			os.Setenv("LD_LIBRARY_PATH", libPath)
		}
		cm.logger("browser_ld_library_path", os.Getenv("LD_LIBRARY_PATH"))
	}

	u := launcher.New().
		Bin(cm.execPath).
		Headless(true).
		NoSandbox(true).
		Set("disable-gpu").
		MustLaunch()

	// Buat browser
	browser := rod.New().ControlURL(u).MustConnect()

	cm.logger("browser_status", "connected")
	return browser, nil
}

// GetExecutablePath mengembalikan path ke executable Chrome
func (cm *ChromiumManager) GetExecutablePath() string {
	return cm.execPath
}
