#!/bin/bash

# https://stanislas.blog/2021/08/firecracker/

# this set of commands creates a  200Mb ext4 filesystem in a file (rootfs.ext4),
# formats it, and mounts it to the /tmp/my-rootfs directory so that it can be
# used like a real filesystem

dd if=/dev/zero of=rootfs.ext4 bs=1M count=200
mkfs.ext4 rootfs.ext4

loop_device=$(sudo losetup --find --show rootfs.ext4)
mkdir -p /tmp/my-rootfs
sudo mount "$loop_device" /tmp/my-rootfs

echo "Filesystem mounted at /tmp/my-rootfs"

docker run -i --rm \
    -v /tmp/my-rootfs:/my-rootfs \
    -v "$(pwd)/agent/agent:/usr/local/bin/agent" \
    -v "$(pwd)/openrc-service.sh:/etc/init.d/agent" \
    alpine sh <setup_alpine.sh

mv rootfs.ext4 ./assets/rootfs.ext4

sudo umount /tmp/my-rootfs
sudo losetup -d "$loop_device"
