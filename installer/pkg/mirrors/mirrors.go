package mirrors

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Absolute path to the mirrorlist file.
const mirrorlistFile string = "/etc/pacman.d/mirrorlist"

// Sets the mirrorlist file with only the servers for the
// given countries and removes all the unused ones.
//
// Can return error types:
//   - MirrorListError
func SetMirrorList(countries []string) error {
	mirrorMap, err := getMirrors()
	if err != nil {
		return MirrorListError{
			err: err,
		}
	}

	if err := saveMirrorlist(countries, mirrorMap); err != nil {
		return MirrorListError{
			err: err,
		}
	}

	return nil
}

// Checks if the given country is present inside mirrorlist.
//
// Can return errors of types:
//   - MirrorListError
func ValidateCountry(country string) error {
	command := fmt.Sprintf("cat %s | grep %s", mirrorlistFile, country)
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return MirrorListError{
			err: err,
		}
	}

	if err := cmd.Start(); err != nil {
		return MirrorListError{
			err: err,
		}
	}

	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return MirrorListError{
			err: err,
		}
	}

	if err := cmd.Wait(); err != nil {
		return MirrorListError{
			err: err,
		}
	}

	if !strings.Contains(string(stdoutBytes), country) {
		return MirrorListError{
			err: errors.New("Invalid country"),
		}
	}

	return nil
}

// Saves all the servers of the given countries inside the
// mirrorlist file.
func saveMirrorlist(countries []string, mirrorMap map[string][]string) error {
	file, err := os.Create(mirrorlistFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, country := range countries {
		for _, server := range mirrorMap[country] {
			if _, err := file.WriteString(server + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

// Reads the mirrorlist file and returns a map
// that has the country name as the key and a slice of
// all the servers as the value.
func getMirrors() (map[string][]string, error) {
	file, err := os.Open(mirrorlistFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var countryMap map[string][]string = make(map[string][]string)
	var lastCountry string = ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == " " {
			continue
		}

		if country, found := strings.CutPrefix(line, "## "); found {
			lastCountry = country
			countryMap[country] = make([]string, 0)
		} else {
			countryMap[lastCountry] = append(countryMap[lastCountry], strings.TrimPrefix(line, "#"))
		}
	}

	return countryMap, nil
}
