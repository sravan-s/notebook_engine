#### Development
We use nix https://nixos.org/ for reproducible envs.
Install it.
Enable flakes on nix -> https://nixos.wiki/wiki/Flakes
```
Add the following to ~/.config/nix/nix.conf or /etc/nix/nix.conf: 
experimental-features = nix-command flakes
```
Fill enviornment variables to `.env` (see `.env.example` for list of env variables)
- We setup an example_webhook in ./example_webhook/
If you want to use it for local development, setup its URL(http://localhost:8080/event) in .env

Install CNI plugins see `./CNI.README.MD`
This is for setting up networking between virtual machines

---
If you use tmux for development ->
```
tmuxp load .
```

To start dev enviornment(nix shell) ->
```
nix devlop
```

To build -> Open nix shell ->
```
go build -o ./notebook-engine.bin .
sudo ./notebook-engine.bin
```
> sudo is required to setup networking
