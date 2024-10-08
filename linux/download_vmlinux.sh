ARCH="$(uname -m)"

latest=$(wget "http://spec.ccfc.min.s3.amazonaws.com/?prefix=firecracker-ci/v1.9/x86_64/vmlinux-5.10&list-type=2" -O - 2>/dev/null | grep "(?<=<Key>)(firecracker-ci/v1.9/x86_64/vmlinux-5\.10\.[0-9]{3})(?=</Key>)" -o -P)

FILE="./assets/vmlinux"

# keep vmlinux file in user root and uncomment to avoid caching
# cp ~/vmlinux ./assets/vmlinux

# Check if the file exists
if [ -f "$FILE" ]; then
  echo "$FILE exists."
else
  echo "$FILE does not exist. Downloading..."
  # Download a linux kernel binary
  wget https://s3.amazonaws.com/spec.ccfc.min/${latest} -O assets/vmlinux
fi
