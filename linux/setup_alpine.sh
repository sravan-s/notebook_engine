# maybe do this with nix?
apk add --no-cache openrc
apk add --no-cache util-linux
apk add --no-cache --update nodejs npm
apk add --no-cache gcompat

ln -s agetty /etc/init.d/agetty.ttyS0
echo ttyS0 > /etc/securetty
rc-update add agetty.ttyS0 default

echo "root:root" | chpasswd

echo "nameserver 1.1.1.1" >>/etc/resolv.conf

# Make sure special file systems are mounted on boot:
rc-update add devfs boot
rc-update add procfs boot
rc-update add sysfs boot
rc-update add agent boot

addgroup -g 1000 -S notebook && adduser -u 1000 -S notebook -G notebook

# Then, copy the newly configured system to the rootfs image:
for d in bin etc lib root sbin usr; do tar c "/$d" | tar x -C /my-rootfs; done

# The above command may trigger the following message:
# tar: Removing leading "/" from member names
# However, this is just a warning, so you should be able to
# proceed with the setup process.

for dir in dev proc run sys var; do mkdir /my-rootfs/${dir}; done

mkdir -p /my-rootfs/tmp
chmod 1777 /my-rootfs/tmp
mkdir -p /my-rootfs/home/notebook/
chown 1000:1000 /my-rootfs/home/notebook/

exit
