package json_parser

import (
	"encoding/json"
	"fmt"

	"github.com/october-os/october-installer/pkg/hostname"
	"github.com/october-os/october-installer/pkg/locale"
	"github.com/october-os/october-installer/pkg/mirrors"
	"github.com/october-os/october-installer/pkg/partition"
	"github.com/october-os/october-installer/pkg/timezone"
	"github.com/october-os/october-installer/pkg/user"
)

// Installation represents all the parameters needed for the installation of October Linux
type Installation struct {
	Drives          []partition.Drive `json:"drives"`
	Users           []user.User       `json:"users"`
	MirrorCountries []string          `json:"mirrorCountries"`
	Timezone        string            `json:"timezone"`
	Locale          string            `json:"locale"`
	Hostname        string            `json:"hostname"`
	RootPassword    string            `json:"rootPassword"`
}

// Parses an installation JSON and validates it
//
// Can return errror types:
// - JsonParsingError
// - partition.ValidationError
// - user.NewUserError
// - mirrors.MirrorListError
// - timezone.TimezoneError
// - locale.LocaleGenError
// - hostname.HostnameError
func ParseJson(jsonString string) (*Installation, error) {
	jsonBytes := []byte(jsonString)
	var installation Installation
	if err := json.Unmarshal(jsonBytes, &installation); err != nil {
		return nil, &JsonParsingError{
			Err: fmt.Errorf("error parsing json: error=%s", err.Error()),
		}
	}

	for _, drive := range installation.Drives {
		if err := drive.Validate(); err != nil {
			return nil, err
		}
	}
	for _, user := range installation.Users {
		if err := user.Validate(); err != nil {
			return nil, err
		}
	}
	for _, country := range installation.MirrorCountries {
		if err := mirrors.ValidateCountry(country); err != nil {
			return nil, err
		}
	}
	if err := timezone.ValidateTimezone(installation.Timezone); err != nil {
		return nil, err
	}
	if err := locale.ValidateLocale(installation.Locale); err != nil {
		return nil, err
	}
	if err := hostname.ValidateHostname(installation.Hostname); err != nil {
		return nil, err
	}

	return &installation, nil
}
