package postinstall

import (
	"fmt"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

func enableSystemd(packages []pkg) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p.name)
		sb.WriteString(" ")
	}

	command := fmt.Sprintf("systemd enable %s", sb.String())
	return arch_chroot.Run(command)
}

func enableUserSystemd(packages []pkg) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p.name)
		sb.WriteString(" ")
	}

	command := fmt.Sprintf("systemd --user enable %s", sb.String())
	return arch_chroot.Run(command)
}
