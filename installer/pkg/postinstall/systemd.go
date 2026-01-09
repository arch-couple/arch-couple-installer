package postinstall

import (
	"fmt"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Takes a list of packages that needs to be
// enabled in systemd and enables them.
func systemdEnable(packages []pkg) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p.name + " ")
	}

	command := fmt.Sprintf("systemd enable %s", sb.String())
	return arch_chroot.Run(command)
}

// Takes a list of packages that needs to be
// enabled in user systemd and enables them.
//
// By "user systemd" it is meant:
//
//	systemd --user enable [package]
func systemdUserEnable(packages []pkg) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p.name + " ")
	}

	command := fmt.Sprintf("systemd --user enable %s", sb.String())
	return arch_chroot.Run(command)
}
