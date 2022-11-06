{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils/v1.0.0";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  description = "Twitch chat archiving stuff I guess";

  outputs = { self, nixpkgs, flake-utils }:
  let
    systems = [ "x86_64-linux" ];
    outputs = flake-utils.lib.eachSystem systems (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in rec {
      packages = { 
        twitch = pkgs.buildGoModule {
          pname = "semver-go";
          version = "v0.1.0";
          modSha256 = pkgs.lib.fakeSha256;
          vendorSha256 = null;
          src = ./.;
        };
      };
      devShell = pkgs.mkShell {
        buildInputs = [
          pkgs.go
        ];
      };
    });
  in
    outputs
}
