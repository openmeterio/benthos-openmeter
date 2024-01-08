{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devenv.url = "github:cachix/devenv";
    dagger.url = "github:dagger/nix";
    dagger.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [ "x86_64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: rec {
        _module.args.pkgs = import inputs.nixpkgs {
          inherit system;

          overlays = [
            (final: prev: {
              dagger = inputs'.dagger.packages.dagger;
            })
          ];
        };

        devenv.shells = {
          default = {
            languages = {
              go.enable = true;
            };

            pre-commit.hooks = {
              nixpkgs-fmt.enable = true;
              commitizen.enable = true;

              commitizen-branch = {
                enable = true;
                name = "commitizen-branch check";
                description = ''
                  Check whether commit messages on the current HEAD follows committing rules.
                '';
                entry = "${pkgs.commitizen}/bin/cz check --allow-abort --rev-range origin/HEAD..HEAD";
                pass_filenames = false;
                stages = [ "manual" ];
              };
            };

            packages = with pkgs; [
              golangci-lint
              dagger

              benthos
            ];

            env = {
              DAGGER_MODULE = "ci";
            };

            # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
            containers = pkgs.lib.mkForce { };
          };

          ci = devenv.shells.default;
        };
      };
    };
}
