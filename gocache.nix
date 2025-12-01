{ pkgs }:

pkgs.stdenv.mkDerivation {
  name = "gocache";
  buildInputs = [ pkgs.go_1_25 ];
  src = pkgs.lib.sourceByRegex ./. [
    "^go.(mod|sum)$"
    "vendor"
    "vendor/.*"
  ];
  configurePhase = "true";
  buildPhase = ''
    export GOFLAGS=-trimpath
    export GOPROXY=off
    export GOSUMDB=off
    export GOCACHE=$out/cache
    mkdir -p $out/cache;

    go build -v `go list ./vendor/...`
    mkdir -p $out/nix-support
    cat > $out/nix-support/setup-hook <<EOF
       cp --reflink=auto -r $out/cache $TMPDIR/go-cache
       chmod -R +w $TMPDIR/go-cache
    EOF
  '';
  installPhase = "true";
  fixupPhase = "true";
}
