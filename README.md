## Setup
1. Setup nix
We use nix https://nixos.org/ for reproducible envs.
Install it.
Enable flakes on nix -> https://nixos.wiki/wiki/Flakes
```
Add the following to ~/.config/nix/nix.conf or /etc/nix/nix.conf: 
experimental-features = nix-command flakes
```

2. Fill .env
Fill enviornment variables to `.env` (see `.env.example` for list of env variables)
- We setup an example_webhook in ./example_webhook/
If you want to use it for local development, setup its URL(http://localhost:8080/event) in .env


3. Setup Container Network Interface
Install CNI plugins see `./CNI.README.MD`
This is for setting up networking between virtual machines
See: https://github.com/firecracker-microvm/firecracker-go-sdk/blob/10626d6b3f442d6b4460357ef38a110e8ca5fb4a/README.md#cni

4. Tmux(optional)
---
If you use tmux for development ->
Install tmux and tmuxp, then, in project root -
```
tmuxp load .
```

----

## Development

1. Setup enviornment(nix shell) ->
```
nix devlop
```

2. Write code/make changes

3. To build VM

You need a firecracker VM with preinstalled agent
Download linux binary(alpine), setup a file system, mount it on docker
install and setup some services, open some port

```
cd linux; ./build.sh
```
^ this needs root permission

4. To build the software

Make sure you are in nix shell; run the following code:

```
go build -o ./notebook-engine.bin .
sudo ./notebook-engine.bin
```
> sudo is required to setup networking

### Thanks

https://k-jingyang.github.io/firecracker/2024/06/15/firecracker-bridge.html
https://stanislas.blog/2021/08/firecracker/
