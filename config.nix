let
  repo = "ssh://git@github.com/alicebob/gohello";
in
{
  users = {
    alice = { };
    bob = { };
    eve = { };
  };
  projects = [
    {
      id = "hello";
      name = "Hello!";
      git = repo;
      #nixfile = "/default.nix";
      attribute = "gohello";
      packages = [
        "which"
        "openssh"
        #"(google-cloud-sdk.withExtraComponents [google-cloud-sdk.components.gke-gcloud-auth-plugin])"
        #"kubectl"
      ];
      post = ''
        echo that was all.
        echo sha: $SHORT_SHA.
        echo pwd: $(pwd)
        echo readlink: $(readlink -f ./result/)
        echo result: $(ls ./result/)
        echo ENV: $(printenv)
        echo which ssh: $(which ssh)
        echo ssh version: $(ssh -V)
        echo result: $(./result/bin/gohello)
      '';
    }
    {
      id = "hello2";
      name = "Hello again - staging!";
      git = repo;
      attribute = "gohello";
      post = ''echo Same thing just to have more projects'';
    }
    {
      id = "recursive";
      name = "Team Builder";
      category = "More examples";
      git = "./";
      attribute = "default";
      post = ''ls result'';
    }
  ];
}
