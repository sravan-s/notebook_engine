{
  description = "Flake for notbeook engine";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  };

  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    devShells.${system}.default =
      pkgs.mkShell
        {
          buildInputs = [
            pkgs.go
            pkgs.firecracker
          ];

          shellHook = ''
            echo hello
          '';
        };
  };
}
