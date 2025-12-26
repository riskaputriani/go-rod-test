package browser

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ulikunitz/xz"
)

// DependencyManager mengelola download dan ekstraksi dependencies
type DependencyManager struct {
	libDir string
	logger func(key, value string)
}

// NewDependencyManager membuat instance baru DependencyManager
func NewDependencyManager(libDir string, logger func(key, value string)) *DependencyManager {
	return &DependencyManager{
		libDir: libDir,
		logger: logger,
	}
}

// Setup mengecek dan mengunduh dependencies jika belum ada
func (dm *DependencyManager) Setup(dependencies []Dependency) error {
	dm.logger("dependencies_lib_dir", dm.libDir)

	// Buat direktori lib jika belum ada
	if err := os.MkdirAll(dm.libDir, 0o755); err != nil {
		return fmt.Errorf("failed to create lib directory: %w", err)
	}

	// Deteksi OS
	osType := dm.detectOS()
	dm.logger("dependencies_os_detected", osType)
	dm.logger("dependencies_arch", runtime.GOARCH)

	// Cek dan download dependencies yang belum ada
	for _, dep := range dependencies {
		if dm.isLibraryInstalled(dep.LibraryName) {
			dm.logger("dependencies_skip", fmt.Sprintf("%s (already installed)", dep.Name))
			continue
		}

		dm.logger("dependencies_downloading", dep.Name)
		if err := dm.downloadAndExtractDep(dep, osType); err != nil {
			dm.logger("dependencies_error", fmt.Sprintf("%s: %v", dep.Name, err))
			// Continue dengan dependency lain
			continue
		}
		dm.logger("dependencies_installed", dep.Name)
	}

	return nil
}

// detectOS mendeteksi jenis OS (debian/ubuntu)
func (dm *DependencyManager) detectOS() string {
	// Coba baca /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "debian" // default
	}

	content := string(data)
	if strings.Contains(strings.ToLower(content), "ubuntu") {
		return "ubuntu"
	}
	if strings.Contains(strings.ToLower(content), "debian") {
		return "debian"
	}

	return "debian" // default
}

// isLibraryInstalled mengecek apakah library sudah terinstall
func (dm *DependencyManager) isLibraryInstalled(libName string) bool {
	// Cek di libDir
	pattern := filepath.Join(dm.libDir, "**", libName+"*")
	matches, _ := filepath.Glob(pattern)
	if len(matches) > 0 {
		return true
	}

	// Cek dengan find command (lebih akurat)
	libPath := filepath.Join(dm.libDir, "usr", "lib")
	if _, err := os.Stat(libPath); err == nil {
		cmd := exec.Command("find", libPath, "-name", libName+"*")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			return true
		}
	}

	return false
}

// downloadAndExtractDep download dan ekstrak .deb package
func (dm *DependencyManager) downloadAndExtractDep(dep Dependency, osType string) error {
	// Pilih URL berdasarkan OS
	url := dep.DebianURL
	if osType == "ubuntu" && dep.UbuntuURL != "" {
		url = dep.UbuntuURL
	}

	dm.logger("dependencies_url", url)

	// Download .deb file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	// Baca data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Ekstrak .deb menggunakan pure Go atau dpkg -x
	return dm.extractDeb(data, dep.Name)
}

// extractDeb mengekstrak .deb package
func (dm *DependencyManager) extractDeb(data []byte, name string) error {
	// Simpan ke file temporary
	tmpFile := filepath.Join(dm.libDir, fmt.Sprintf(".%s.deb", name))
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	dm.logger("dependencies_extracting", name)

	// Gunakan dpkg -x jika tersedia (lebih reliable)
	cmd := exec.Command("dpkg", "-x", tmpFile, dm.libDir)
	if err := cmd.Run(); err != nil {
		// Jika dpkg tidak tersedia, coba ekstrak manual
		dm.logger("dependencies_dpkg_not_found", "trying manual extraction")
		return dm.extractDebManual(data)
	}

	return nil
}

// extractDebManual mengekstrak .deb secara manual (pure Go)
func (dm *DependencyManager) extractDebManual(data []byte) error {
	// .deb adalah ar archive yang berisi:
	// - debian-binary
	// - control.tar.xz (metadata)
	// - data.tar.xz (actual files)

	// Kita perlu extract data.tar.xz
	reader := bytes.NewReader(data)

	// Skip ar header (8 bytes: "!<arch>\n")
	header := make([]byte, 8)
	if _, err := reader.Read(header); err != nil {
		return fmt.Errorf("failed to read ar header: %w", err)
	}

	// Parse ar entries
	for {
		// Read ar entry header (60 bytes)
		entryHeader := make([]byte, 60)
		n, err := reader.Read(entryHeader)
		if err == io.EOF {
			break
		}
		if err != nil || n != 60 {
			break
		}

		// Parse filename (first 16 bytes)
		filename := strings.TrimSpace(string(entryHeader[0:16]))

		// Parse file size (bytes 48-58)
		sizeStr := strings.TrimSpace(string(entryHeader[48:58]))
		var size int64
		fmt.Sscanf(sizeStr, "%d", &size)

		// Read file data
		fileData := make([]byte, size)
		if _, err := io.ReadFull(reader, fileData); err != nil {
			return fmt.Errorf("failed to read ar entry: %w", err)
		}

		// Align to 2 bytes
		if size%2 == 1 {
			reader.ReadByte()
		}

		// Process data.tar.xz
		if strings.Contains(filename, "data.tar") {
			return dm.extractDataTar(fileData, filename)
		}
	}

	return fmt.Errorf("data.tar not found in deb package")
}

// extractDataTar ekstrak data.tar.xz atau data.tar.gz
func (dm *DependencyManager) extractDataTar(data []byte, filename string) error {
	var tarReader *tar.Reader

	if strings.HasSuffix(filename, ".xz") || strings.Contains(filename, "tar.xz") {
		// Decompress xz
		xzReader, err := xz.NewReader(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to create xz reader: %w", err)
		}
		tarReader = tar.NewReader(xzReader)
	} else if strings.HasSuffix(filename, ".gz") || strings.Contains(filename, "tar.gz") {
		// TODO: Add gzip support if needed
		return fmt.Errorf("gzip format not yet supported")
	} else {
		// Plain tar
		tarReader = tar.NewReader(bytes.NewReader(data))
	}

	// Extract tar
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		target := filepath.Join(dm.libDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0o755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0o755)
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				continue
			}
			io.Copy(outFile, tarReader)
			outFile.Close()
		case tar.TypeSymlink:
			os.Symlink(header.Linkname, target)
		}
	}

	return nil
}

// GetLibraryPath mengembalikan path yang harus ditambahkan ke LD_LIBRARY_PATH
func (dm *DependencyManager) GetLibraryPath() string {
	// Kembalikan semua path yang mungkin berisi .so files
	paths := []string{
		filepath.Join(dm.libDir, "usr", "lib", "x86_64-linux-gnu"),
		filepath.Join(dm.libDir, "usr", "lib64"),
		filepath.Join(dm.libDir, "usr", "lib"),
		filepath.Join(dm.libDir, "lib", "x86_64-linux-gnu"),
		filepath.Join(dm.libDir, "lib64"),
		filepath.Join(dm.libDir, "lib"),
	}

	var existingPaths []string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			existingPaths = append(existingPaths, path)
		}
	}

	return strings.Join(existingPaths, ":")
}
