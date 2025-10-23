package gui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gordonklaus/portaudio"
)

const (
	blackholePkgURL      = "https://existential.audio/downloads/BlackHole2ch.pkg"
	blackholeDeviceName  = "BlackHole 2ch"
	blackholeTmpFileName = "BlackHole2ch.pkg"
)

// isBlackHoleInstalled checks if BlackHole audio device is available
func isBlackHoleInstalled() (bool, error) {
	// Initialize PortAudio temporarily to check devices
	if err := portaudio.Initialize(); err != nil {
		return false, err
	}
	defer func() {
		_ = portaudio.Terminate() // Ignore error on cleanup
	}()

	devices, err := portaudio.Devices()
	if err != nil {
		return false, err
	}

	// Check if BlackHole 2ch device exists
	for _, device := range devices {
		if strings.Contains(device.Name, "BlackHole") {
			return true, nil
		}
	}

	return false, nil
}

// downloadBlackHole downloads the BlackHole installer package
func downloadBlackHole(progressCallback func(downloaded, total int64)) (string, error) {
	// Create temp directory for download
	tmpDir := os.TempDir()
	pkgPath := filepath.Join(tmpDir, blackholeTmpFileName)

	// Download the package
	resp, err := http.Get(blackholePkgURL)
	if err != nil {
		return "", fmt.Errorf("failed to download BlackHole: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Ignore error on cleanup
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download BlackHole: HTTP %d", resp.StatusCode)
	}

	// Create the file
	out, err := os.Create(pkgPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = out.Close() // Ignore error on cleanup
	}()

	// Copy with progress
	totalSize := resp.ContentLength
	downloaded := int64(0)
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return "", writeErr
			}
			downloaded += int64(n)
			if progressCallback != nil {
				progressCallback(downloaded, totalSize)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return pkgPath, nil
}

// installBlackHole installs the BlackHole driver package
// This requires admin privileges and will prompt for password
func installBlackHole(pkgPath string) error {
	// Use 'open' command to install the .pkg file
	// This will show the macOS installer UI and handle admin privileges
	cmd := exec.Command("open", pkgPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open installer: %w", err)
	}

	return nil
}

// cleanupTempFiles removes temporary installer files
func cleanupTempFiles() {
	tmpDir := os.TempDir()
	pkgPath := filepath.Join(tmpDir, blackholeTmpFileName)
	_ = os.Remove(pkgPath) // Ignore error if file doesn't exist
}
