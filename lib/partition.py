from __future__ import annotations

import subprocess
from enum import Enum

from typing_extensions import List, Optional


class Partition:
    def __init__(
        self,
        drive: str,
        size: PartitionSize,
        type: GPTPartitionType,
        mount_point: Optional[str] = None,
    ) -> None:
        self.drive = drive  # NOTE: used to group partitions according to drive (sfdisk only allows creating partitions for one device at a time)
        self.size = size
        self.type = type
        self.mount_point = mount_point

    def to_sfdisk_format(self) -> str:
        partition_string = f"uuid={self.type}"

        if (
            not self.size.takeRemaining
            and self.size.unit is not None
            and self.size.amount is not None
        ):
            partition_string += f", size={self.size.amount + self.size.unit}"

        if self.size.takeRemaining:
            partition_string += ', size="+"'

        return partition_string


class PartitionSize:
    def __init__(
        self,
        amount: Optional[str] = None,
        unit: Optional[str] = None,
        takeRemaining: bool = False,
    ) -> None:
        self.amount = amount
        self.unit = unit
        self.takeRemaining = takeRemaining


class GPTPartitionType(Enum):
    EFI = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
    SWAP = "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"
    ROOT = "4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709"
    FILE_SYSTEM = "0FC63DAF-8483-4772-8E79-3D69D8477DE4"
    HOME = "933AC7E1-2EB4-4F13-B844-0E14E2AEF915"


def create_partitions(partitions: List[Partition]) -> None:
    # will use _create_partition_string to create files (sfdisk named-fields format) for each drive (sfdisk only allows to create partitions for one drive at a time)
    pass
