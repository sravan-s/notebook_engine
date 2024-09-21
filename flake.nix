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

          # shellHook = ''
          #   echo hello
          # '';
        };
    # packages.${system}.default = pkgs.stdenv.mkDerivation {
    #   name = "notebook-engine";
    #   version = "0.0.0";
    #   system = system;
    #   dontUnpack = true;
    #   buildInputs = [ pkgs.go ];
    #   buildPhase = ''
    #     go build -o ./notebook-engine.bin ./src/index.go
    #   '';
    # };
  };
}
