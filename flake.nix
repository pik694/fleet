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
      version = self.shortRev or self.dirtyShortRev or "dev";
      goFor =
        pkgs:
        pkgs.go_1_26.overrideAttrs (_: rec {
          version = "1.26.11";
          src = pkgs.fetchurl {
            url = "https://go.dev/dl/go${version}.src.tar.gz";
            hash = "sha256-1EEBsP8n7SpMda1GwXucdDWSCeLrLKfJ32y0/baOGWk=";
          };
        });
    in
    {
      packages = forAllSystems (pkgs: rec {
        fleet = (pkgs.buildGoModule.override { go = goFor pkgs; }) {
          pname = "fleet";
          inherit version;
          src = self;
          vendorHash = "sha256-H7DaxdTyerlzkY/uvFw/GBxSCIWwXgBUcdCxspYK++c=";
          subPackages = [ "cmd/fleet" ];
          env.CGO_ENABLED = 0;
          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];
          doCheck = false;
          meta = {
            description = "TUI for running coding-agent sessions in parallel with tmux";
            homepage = "https://github.com/pik694/fleet";
            license = pkgs.lib.licenses.asl20;
            mainProgram = "fleet";
            platforms = systems;
          };
        };
        default = fleet;
      });

      overlays.default = final: _prev: {
        fleet = self.packages.${final.stdenv.hostPlatform.system}.fleet;
      };

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            (goFor pkgs)
            gopls
            golangci-lint
            tmux
            git
          ];
        };
      });
    };
}
