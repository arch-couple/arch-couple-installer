package postinstall

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

type pkg struct {
	name  string
	flags []string
}

const packageFilePath string = "/root/postinstall/packages"
const aurFilePath string = "/root/postinstall/aur"

func downloadAllPackages(packages []pkg, aur bool) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p.name)
		sb.WriteString(" ")
	}

	var command string
	if aur {
		command = fmt.Sprintf("sudo -u builder yay -S %s --noconfirm", sb.String())
	} else {
		command = fmt.Sprintf("pacman -Sy --noconfirm %s", sb.String())
	}

	return arch_chroot.Run(command)
}

func getPackageList(path string) ([]pkg, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	var packageList []pkg

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) != 0 && line[0] == '-' {
			splittedLine := strings.Split(line, " ")
			newPackage := pkg{
				name: splittedLine[1],
			}

			if len(splittedLine) > 2 {
				var flags []string
				for i := 2; i < len(splittedLine); i++ {
					flags = append(flags, splittedLine[i])
				}

				newPackage.flags = flags
			}
			packageList = append(packageList, newPackage)
		}
	}

	return packageList, nil
}

func packageFlagParser(packages []pkg) error {
	var systemdEnable []pkg
	var systemdUserEnable []pkg

	for _, p := range packages {
		if len(p.flags) == 0 {
			continue
		}

		for _, f := range p.flags {
			switch f {
			case "SYSTEMD_ENABLE":
				systemdEnable = append(systemdEnable, p)
			case "SYSTEMD_USER_ENABLE":
				systemdUserEnable = append(systemdUserEnable, p)
			}
		}
	}

	if len(systemdEnable) > 0 {
		if err := enableSystemd(systemdEnable); err != nil {
			return err
		}
	}

	if len(systemdUserEnable) > 0 {
		if err := enableUserSystemd(systemdUserEnable); err != nil {
			return err
		}
	}

	return nil
}
