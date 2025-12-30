# JSON payload for main parser

```json
{
  "drives" : [
    {
      "drive": "/dev/sda",
      "append": true
      "partitions" : []
    }
  ]
  "partitions": [
    {
      "drive": "[path to drive]",
      "size": {
        "size": 1234,
        "unit": "[unit]",
        "takeRemaining": true
      },
      "fileSystem": "btrfs/ext4"
      "partitionType": "[partition type]",
      "label": "[drive label]",
      "mountPoint": "[dir]"
    }
  ],
  "users": [
    {
      "username": "[username]",
      "password": "[password]",
      "homepath": "[path to home]",
      "sudoer": true,
    }
  ],
  "timezone": "[user timezone]",
  "locale": "[user locale]",
  "hostname": "[user hostname]",
  "rootpassword": "[root password]"
}
```
