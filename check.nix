{ pkgs, system,  nixpkgs, self, ... }: {
  # A VM test of the NixOS module.
  vmTest = with import (nixpkgs+ "/nixos/lib/testing-python.nix") { inherit system; };

    (makeTest {
      name = "webhooker-test";
      nodes = {
        server =
          let
            configEnv = pkgs.writeTextFile {
              name = "envfile";
              text = ''
                HOOKER_TOKEN="global-testtoken"
                HOOKER_TOKEN_HELLO="hello-token"
              '';
            };
          in
          {
            imports = [ self.nixosModules.webhooker ];
            services.webhooker = {
              enable = true;
              envFile = "${configEnv}";
              settings = {
                host = "0.0.0.0";
                port = "8888";
                hooks = {
                  hello.command = "echo 'hello'";
                  world.command = "echo 'world'";
                };
              };
            };

            # Open firewall for testing
            networking.firewall.allowedTCPPorts = [ 8888 ];
          };
      };

      testScript = /*python*/ ''
        start_all()
        server.wait_for_unit("multi-user.target")
        server.wait_for_unit("webhooker.service")
        # machine.wait_for_open_port(8888)

        machine.succeed("curl http://localhost:8888/hello")
        machine.succeed("curl http://localhost:8888/world")

        # Fail authentication, when no token is provided
        machine.fail("curl --fail http://localhost:8888/world")

        # Fail authentication, when wrong token is provided
        machine.fail("curl --fail http://localhost:8888/world -H 'Token: hello-token'")

        # Succeed when authenticating wiht the global token
        machine.succeed("curl --fail http://localhost:8888/hello -H 'Token: global-testtoken'")

        # Succeed when authenticating wiht the hook-specific token
        machine.succeed("curl --fail http://localhost:8888/hello -H 'Token: hello-token'")
      '';
    }).test;

}
