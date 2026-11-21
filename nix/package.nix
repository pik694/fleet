{
  lib,
  buildGoModule,
  go_1_26,
  fetchurl,
  makeBinaryWrapper,
  tmux,
  git,
  src,
  version ? "dev",
  extraPackages ? [ ],
}:

let
  go = go_1_26.overrideAttrs (_: rec {
    version = "1.26.11";
    src = fetchurl {
      url = "https://go.dev/dl/go${version}.src.tar.gz";
      hash = "sha256-1EEBsP8n7SpMda1GwXucdDWSCeLrLKfJ32y0/baOGWk=";
    };
  });
  runtimePackages = lib.filter (p: p != null) [
    tmux
    git
  ]
  ++ extraPackages;
in
(buildGoModule.override { inherit go; }) {
  pname = "fleet";
  inherit version src;
  vendorHash = "sha256-H7DaxdTyerlzkY/uvFw/GBxSCIWwXgBUcdCxspYK++c=";
  subPackages = [ "cmd/fleet" ];
  env.CGO_ENABLED = 0;
  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];
  doCheck = false;

  nativeBuildInputs = [ makeBinaryWrapper ];
  postInstall = lib.optionalString (runtimePackages != [ ]) ''
    wrapProgram $out/bin/fleet --suffix PATH : ${lib.makeBinPath runtimePackages}
  '';

  passthru = { inherit go; };

  meta = {
    description = "TUI for running coding-agent sessions in parallel with tmux";
    homepage = "https://github.com/pik694/fleet";
    license = lib.licenses.asl20;
    mainProgram = "fleet";
    platforms = lib.platforms.darwin ++ lib.platforms.linux;
  };
}
