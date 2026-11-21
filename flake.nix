{
  description = "fleet — TUI for running Claude Code, Codex, OpenCode and Copilot sessions in tmux";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "aarch64-darwin"
        "x86_64-darwin"
        "aarch64-linux"
        "x86_64-linux"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
      fleetFor =
        pkgs:
        pkgs.callPackage ./nix/package.nix {
          src = self;
          version = self.shortRev or self.dirtyShortRev or "dev";
        };
    in
    {
      packages = forAllSystems (pkgs: rec {
        fleet = fleetFor pkgs;
        default = fleet;
      });

      overlays.default = final: _prev: {
        fleet = fleetFor final;
      };

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            (fleetFor pkgs).go
            gopls
            golangci-lint
            tmux
            git
          ];
        };
      });
    };
}
