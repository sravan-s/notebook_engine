rm -rf ./assets
mkdir -p assets

sudo rm -rf ./executer

./build_agent.sh
./download_vmlinux.sh
./make_rootfs.sh
