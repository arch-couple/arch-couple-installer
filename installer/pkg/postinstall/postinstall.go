package postinstall

func DownloadPostInstallPackages() error {
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

	if err := installAurHelperAndPackages(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

func installAurHelperAndPackages() error {
	if err := activateBuilderAccount(); err != nil {
		return err
	}

	if err := installYay(); err != nil {
		return err
	}

	packages, err := getPackageList(aurFilePath)
	if err != nil {
		return err
	}

	if err := downloadAllPackages(packages, true); err != nil {
		return err
	}

	return deletingBuilderAccount()
}
