{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    flake-compat.url = "github:edolstra/flake-compat";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          lib,
          pkgs,
          system,
          ...
        }:
        let
          overlays = [ inputs.gomod2nix.overlays.default ];

          goku = pkgs.buildGoApplication {
            name = "goku";
            src = lib.cleanSource ./.;
            modules = ./gomod2nix.toml;
          };
        in
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system overlays;
          };

          treefmt = {
            projectRootFile = ".git/config";

            # Nix
            programs.nixfmt.enable = true;

            # Go
            programs.gofmt.enable = true;

            # GitHub Actions
            programs.actionlint.enable = true;

            # Markdown
            programs.mdformat.enable = true;

            # ShellScript
            programs.shellcheck.enable = true;
            programs.shfmt.enable = true;
          };

          packages = {
            inherit goku;
            default = goku;
          };

          checks = {
            inherit goku;
          };

          devShells.default = pkgs.mkShell {
            nativeBuildInputs = [
              pkgs.go # Golang
              pkgs.nil # Nix LSP
              pkgs.gomod2nix # gomod2nix for creating Hashes (./gomod2nix.toml)
            ];

            inputsFrom = [ config.treefmt.build.devShell ];
          };
        };
    };
}
