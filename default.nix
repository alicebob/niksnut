let
  sources = import ./build/default.nix;
  pkgs = import sources.nixpkgs { };
  gocache = pkgs.callPackage ./gocache.nix { };

  niksnut = pkgs.buildGoModule {
    name = "niksnut";
    buildInputs = [ gocache ];
    src = pkgs.lib.sourceByRegex ./. [
      "go.(mod|sum)"
      ".*\.go"
      "vendor"
      "vendor/.*"
      "httpd"
      "httpd/.*"
      "niks"
      "niks/.*"
      "static"
      "static/.*"
    ];
    vendorHash = null; # uses ./vendor/
    doCheck = false;
  };
in
{
  default = niksnut;
  niksnut = niksnut;
  shell = pkgs.mkShellNoCC {
    packages = [
      pkgs.nixfmt-rfc-style
      pkgs.npins
      pkgs.nodePackages.prettier
    ];
  };
}
