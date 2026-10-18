{
  description = "The website of Arkadia, https://arkadia-lang.github.io/";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { nixpkgs, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      # `nix develop`: Go and just build the ark command of the Arkadia
      # checkout and the standard library reference, pandoc the pages.
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = [ pkgs.go pkgs.just pkgs.pandoc ];
          CGO_ENABLED = "0";
        };
      });
    };
}
