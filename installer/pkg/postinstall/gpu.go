package postinstall

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// AMD: mesa, lib32-mesa, vulkan-radeon, lib32-vulkan-radeon
// Intel: mesa, lib32-mesa, vulkan-intel, lib32-vulkan-intel
// NVIDIA: TUXXX, GAXXX, ADXXX: nvidia-open ---- GMXXX, GPXXX, GVXXX: nvidia-580xx-dkms (aur)

var amdGPUPackages []string = []string{"mesa", "lib32-mesa", "vulkan-radeon", "lib32-vulkan-radeon"}
var intelGPUPackages []string = []string{"mesa", "lib32-mesa", "vulkan-intel", "lib32-vulkan-intel"}

const nvidiaOpenGPUPackage string = "nvidia-open"
const nvidiaProprietaryGPUPackage string = "nvidia-580xx-dkms"

func BestEffortGPUDrivers() error {
	brand, err := getGPUBrand()
	if err != nil {
		return err
	}

	var officialPackages []string = make([]string, 0)
	var AURPackages []string = make([]string, 0)

	switch brand {
	case "Intel":
		officialPackages = slices.Concat(officialPackages, amdGPUPackages)
	case "AMD":
		officialPackages = slices.Concat(officialPackages, intelGPUPackages)
	default:
		if brand == "" {
			// ERROR LOL
			// maybe move to getGPUBrand() and handle above
		}
		if strings.Contains(brand, "NVIDIA") {
			// TODO
		}
	}

	fileOfficialPackages, err := os.OpenFile(packageFilePath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return &PostInstallError{
			err: fmt.Errorf("error writing gpu packages to '%s': error=%s", packageFilePath, err.Error()),
		}
	}
	defer fileOfficialPackages.Close()
	// TODO
	//if fileOfficialPackages.WriteString()
	return nil
}

func getGPUBrand() (string, error) {
	cmd := exec.Command("lspci", "|", "grep", "-i", "'VGA compatible controller'")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", &PostInstallError{
			err: fmt.Errorf("error piping stdout: error=%s", err.Error()),
		}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", &PostInstallError{
			err: fmt.Errorf("error piping stderr: error=%s", err.Error()),
		}
	}
	if err := cmd.Start(); err != nil {
		stderrOutput, _ := io.ReadAll(stderr)
		return "", &PostInstallError{
			err: fmt.Errorf("error getting GPU information: error=%s", string(stderrOutput)),
		}
	}
	var stdoutOutput []byte
	if stdoutOutput, err = io.ReadAll(stdout); err != nil {
		return "", &PostInstallError{
			err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
		}
	}
	if err := cmd.Wait(); err != nil {
		return "", &PostInstallError{
			err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
		}
	}

	stdoutOutputString := string(stdoutOutput)
	if strings.Contains(stdoutOutputString, "Intel") {
		return "Intel", nil
	}
	if strings.Contains(stdoutOutputString, "AMD") {
		return "AMD", nil
	}
	if strings.Contains(stdoutOutputString, "NVIDIA") {
		// TODO
		return "NVIDIA XXXXXXXXXX", nil
	}

	return "", nil
}
