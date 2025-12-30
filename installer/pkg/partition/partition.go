package partition

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Creates one file per drive containing its partitions in sfdisk named-fields syntax
// from a list of Drives
// Returns a map of the drives names and their files name
// Can also return one type of error:
//   - CreatePartitionsError:
//     when a file couldn't be created
//     or
//     when a file couldn't be modified
func createPartitioningFiles(drives []*Drive) (map[*Drive]string, error) {
	drivePartitionsFiles := make(map[*Drive]string)
	for _, drive := range drives {
		fileName := strings.ReplaceAll(drive.Path, "/", "")
		fullFilename := fmt.Sprintf("%s.txt", fileName)
		drivePartitionsFiles[drive] = fullFilename

		file, err := os.Create(fullFilename)
		if err != nil {
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("could not create file '%s' for '%s' drive partitioning: error=%s", fullFilename, drive.Path, err.Error()),
			}
		}
		defer file.Close()

		for _, partition := range drive.Partitions {
			partitionEntry := fmt.Sprintf("%s\n", partition.toSfdiskFormat())
			_, err := file.WriteString(partitionEntry)
			if err != nil {
				return nil, &SetupPartitionsError{
					Err: fmt.Errorf("could not edit file '%s' for partitioning: error=%s", fullFilename, err.Error()),
				}
			}
		}
	}
	return drivePartitionsFiles, nil
}

// Create Partitions from a list of Drives using sfdisk
//
// Can also return one type of error:
//   - CreatePartitionsError:
//     when the creation of partitions failed using sfdisk
//     or
//     stderr couldn't be piped
func createPartitions(drives []*Drive) ([]map[Partition]SfdiskJsonPartition, error) {
	partitioningFiles, err := createPartitioningFiles(drives)
	if err != nil {
		return nil, err
	}

	var mappings []map[Partition]SfdiskJsonPartition

	for drive, fileName := range partitioningFiles {
		sfdiskCommand := ""
		var initialState *SfdiskJsonDrive

		if drive.Append {
			initialState, err = getDriveStateWithSfdisk(drive.Path)
			if err != nil {
				return nil, &SetupPartitionsError{
					Err: fmt.Errorf("error getting initial state of drive '%s': error=%s", drive.Path, err.Error()),
				}
			}
			sfdiskCommand = fmt.Sprintf("sfdisk -a %s < %s", drive.Path, fileName)
		} else {
			sfdiskCommand = fmt.Sprintf("sfdisk %s < %s", drive.Path, fileName)
		}

		cmd := exec.Command("/bin/bash", "-c", sfdiskCommand)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("error piping stderr: error=%s", err.Error()),
			}
		}

		if err := cmd.Run(); err != nil {
			stderrOutput, _ := io.ReadAll(stderr)
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("error creating partitions on drive '%s' with file '%s' using sfdisk: error=%s", drive.Path, fileName, string(stderrOutput)),
			}
		}

		stateAfterCreatingPartitions, err := getDriveStateWithSfdisk(drive.Path)
		var newPartitions []SfdiskJsonPartition
		if err != nil {
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("error getting state after partitions creation of drive '%s': error=%s", drive.Path, err.Error()),
			}
		}
		if initialState != nil {
			newPartitions = stateAfterCreatingPartitions.PartitionTable.Partitions[len(initialState.PartitionTable.Partitions):]
		} else {
			newPartitions = stateAfterCreatingPartitions.PartitionTable.Partitions
		}

		partitionsMap := make(map[Partition]SfdiskJsonPartition)
		for i := 0; i < len(newPartitions) || i < len(drive.Partitions); i++ {
			partitionsMap[drive.Partitions[i]] = newPartitions[i]
		}

		mappings = append(mappings, partitionsMap)
	}

	return mappings, nil
}

func SetupPartitions(drives []*Drive) error {
	newPartitionsMappings, err := createPartitions(drives)
	if err != nil {
		return err
	}

	for _, mapping := range newPartitionsMappings {
		for partition, sfdiskPartition := range mapping {
			if err = formatPartition(&partition, sfdiskPartition.Node); err != nil {
				return err
			}
			if err = mountPartition(&partition, sfdiskPartition.Node); err != nil {
				return err
			}
		}
	}

	return nil
}

func formatPartition(partition *Partition, path string) error {
	var cmd *exec.Cmd
	switch partition.PartitionType {
	case gptPartitionTypeEfi:
		cmd = exec.Command("mkfs.fat", "-F", "32", path)
	case gptPartitionTypeSwap:
		cmd = exec.Command("mkswap", path)
	case gptPartitionTypeRoot, gptPartitionTypeHome, gptPartitionTypeFileSystem:
		switch *partition.FileSystem {
		case fileSystemExt4:
			cmd = exec.Command("mkfs.ext4", path)
		case fileSystemBtrfs:
			cmd = exec.Command("mkfs.btrfs", path)
		}
	}

	if cmd == nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error formatting partition '%s'", path),
		}
	}

	if err := cmd.Run(); err != nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error formatting partition '%s': error=%s", path, err.Error()),
		}
	}

	return nil
}

func mountPartition(partition *Partition, path string) error {
	var cmd *exec.Cmd
	switch partition.PartitionType {
	case gptPartitionTypeEfi:
		cmd = exec.Command("mount", "--mkdir", path, "/mnt/boot")
	case gptPartitionTypeSwap:
		cmd = exec.Command("swapon", path)
	case gptPartitionTypeRoot:
		cmd = exec.Command("mount", path, "/mnt")
	case gptPartitionTypeHome, gptPartitionTypeFileSystem:
		cmd = exec.Command("mount", "--mkdir", path, *partition.MountPoint)
	}

	if cmd == nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error mounting partition %s", path),
		}
	}

	if err := cmd.Run(); err != nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error mounting partition '%s': error=%s", path, err.Error()),
		}
	}

	return nil
}
