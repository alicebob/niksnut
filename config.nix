let
  # pkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/9b5328b7f761a7bbdc0e332ac4cf076a3eedb89b.tar.gz") {};
  # sources = import ./build/default.nix;
  # pkgs = import sources.nixpkgs { };
  repo = "ssh://git@github.com/alicebob/gohello";
in
{
	users = {
		alice = {};
		bob = {};
		eve = {};
	};
	projects = [
		{
			id = "hello";
			name = "Hello!";
			git = repo;
			#nixfile = "/default.nix";
			attribute = "gohello";
			# buildInputs = [pkgs.openssh];
		    packages = [
		    "which"
			"openssh"
    		"(google-cloud-sdk.withExtraComponents [google-cloud-sdk.components.gke-gcloud-auth-plugin])"
			"kubectl"
];
			post = ''
				echo that was it!.
				echo pwd: $(pwd)
				echo readlink: $(readlink -f ./result/)
				echo result: $(ls ./result/)
				echo ENV: $(printenv)
				echo which ssh: $(which ssh)
				echo ssh version: $(ssh -V)
				echo which kubectl: $(which kubectl)
				echo kubectl version: $(kubectl -V)
			'';
		}
		{
			id = "hello2";
			name = "Hello again!";
			git = repo;
			attribute = "gohello";
			post = ''Same thing just to have more projects'';
		}
	];
}
