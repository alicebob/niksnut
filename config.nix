{ pkgs, ... }:
let
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
			post = ''
				echo that was it!.
				echo pwd: $(pwd)
				echo readlink: $(readlink -f ./result/)
				echo result: $(ls ./result/)
				echo ENV: $(printenv)
				$(${pkgs.openssh}/bin/ssh -V)
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
