{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils/v1.0.0";
  };

  description = "twitch-chat-archiver, archives chats or something";

  outputs = { self, nixpkgs, flake-utils }:
  flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
      build = pkgs.buildGoModule {
        pname = "twitch-chat-archiver";
        version = "v0.1.0";
        modSha256 = pkgs.lib.fakeSha256;
        vendorSha256 = null;
        src = ./.;

        meta = {
          description = "twitch-chat-archiver, archives chats or something";
          homepage = "https://github.com/catouc/noita-chat-wars";
          license = pkgs.lib.licenses.mit;
          maintainers = [ "catouc" ];
          platforms = pkgs.lib.platforms.linux ++ pkgs.lib.platforms.darwin;
        };
      };
  in
    rec {
      packages = {
        twitch-chat-archiver = build;
        default = build;
      };

      devShells = {
        default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
          ];
        };
      };
    }
  );
}
