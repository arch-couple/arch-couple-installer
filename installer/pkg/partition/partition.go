package partition

type GptPartitionType string

const (
	efi        GptPartitionType = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
	swap       GptPartitionType = "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"
	root       GptPartitionType = "4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709"
	fileSystem GptPartitionType = "0FC63DAF-8483-4772-8E79-3D69D8477DE4"
	home       GptPartitionType = "933AC7E1-2EB4-4F13-B844-0E14E2AEF915"
)

type Partition struct {
	Drive         string
	Size          PartitionSize
	PartitionType GptPartitionType
	MountPoint    string
}

type PartitionSize struct {
	Amount        *int64
	Unit          *string
	TakeRemaining bool
}

func NewPartition(drive string, size *PartitionSize, partitionType GptPartitionType, mountPoint string) *Partition {
	return &Partition{
		Drive:         drive,
		Size:          *size,
		PartitionType: partitionType,
		MountPoint:    mountPoint,
	}
}

func NewPartitionSize(amount *int64, unit *string, takeRemaining *bool) *PartitionSize {
	var takeRemainingValue bool
	if takeRemaining == nil {
		takeRemainingValue = false
	} else {
		takeRemainingValue = *takeRemaining
	}

	return &PartitionSize{
		Amount:        amount,
		Unit:          unit,
		TakeRemaining: takeRemainingValue,
	}
}
