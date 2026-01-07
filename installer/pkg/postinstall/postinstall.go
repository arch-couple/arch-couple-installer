package postinstall

import "github.com/october-os/october-installer/pkg/arch_chroot"

func InstallPostInstallPackages() error {
	packages, err := getPackageList(packageFilePath)
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := downloadAllPackages(packages, false); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := packageFlagParser(packages); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

func InstallAurHelperAndPackages() error {
	if err := activateBuilderAccount(); err != nil {
		return err
	}

	if err := installYay(); err != nil {
		return err
	}

	packages, err := getPackageList(aurFilePath)
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := downloadAllPackages(packages, true); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := deleteBuilderAccount(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := packageFlagParser(packages); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

func EnableMultilibRepo() error {
	command := "sed -i -e '/#\\[multilib\\]/,+1s/^#//' /etc/pacman.conf"
	return arch_chroot.Run(command)
}

func SetupSudoerFile() error {
	if err := addWheelGroup(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}
