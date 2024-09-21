#### Development
We use nix https://nixos.org/ for reproducible envs.
Install it.
Enable flakes on nix -> https://nixos.wiki/wiki/Flakes
```
Add the following to ~/.config/nix/nix.conf or /etc/nix/nix.conf: 
experimental-features = nix-command flakes
```

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
go build -o ./notebook-engine.bin ./src/index.go
```
