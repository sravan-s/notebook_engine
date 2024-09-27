## To install CNI (Container Network Interface)

- on your local machine, you'll need to download and set up the CNI plugins. The CNI project provides a collection of networking plugins and libraries to configure networking interfaces in Linux containers.

Here's a step-by-step guide to install CNI on a local machine:

1. Install CNI Plugins

The CNI project provides a standard set of plugins (e.g., bridge, host-local, loopback, etc.) that you can download.
Steps:

    Download the CNI Plugins:

    Go to the CNI Plugins Releases page and download the latest release tarball.
```bash
curl -LO https://github.com/containernetworking/plugins/releases/download/v1.3.0/cni-plugins-linux-amd64-v1.3.0.tgz
```
You can adjust the link to reflect the latest release or the architecture (e.g., ARM, etc.).

Extract the Tarball:

Extract the contents of the downloaded tarball.

```bash
tar -zxvf cni-plugins-linux-amd64-v1.3.0.tgz
```
Move Plugins to the CNI Directory:

Typically, CNI expects plugins to be installed in /opt/cni/bin. You can create that directory and move the extracted binaries there.

```bash
    sudo mkdir -p /opt/cni/bin
    sudo mv * /opt/cni/bin
```
    The /opt/cni/bin directory is the default directory where most CNI-compatible container runtimes (like Kubernetes) look for CNI plugins.

2. Verify CNI Plugins

To ensure that the CNI plugins are installed correctly, check the contents of the /opt/cni/bin directory:

```bash
ls /opt/cni/bin
```
You should see a list of CNI plugins like:

```
bridge
host-local
ipvlan
loopback
macvlan
portmap
tuning
```
These are the basic CNI plugins used for container networking.

3. Install tc-redirect-tap
This one is a plugin, but is not present in previous list, need download from
https://github.com/awslabs/tc-redirect-tap

```bash
mkdir -p ~/src/github.com/awslabs/tc-redirect-tap
cd ~/src/github.com/awslabs/tc-redirect-tap
git clone https://github.com/awslabs/tc-redirect-tap.git .
sudo make install
```
