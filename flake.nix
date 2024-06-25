{
  description = "Listen for webhooks, run commands";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let

      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    {

      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          webhooker = pkgs.buildGoModule {
            pname = "webhooker";
            inherit version;
            src = ./.;
            vendorHash = "sha256-eqkfWhlp++hR/x39/0GSqUKugUb9X13+qDCrm46tSj4=";
            HOOKER_CONFIG = ./testdata/config.yaml;
            meta.mainProgram = "webhooker";
          };
          default = webhooker;
        });

      nixosModules.webhooker = { pkgs, lib, config, ... }: {
        imports = [ ./module.nix ];
        config.services.webhooker.package = self.packages.${pkgs.system}.webhooker;
      };

      # Tests run by 'nix flake check' and by Hydra.
      checks = forAllSystems
        (system:
          with nixpkgsFor.${system};

          lib.optionalAttrs stdenv.isLinux {
            # A VM test of the NixOS module.
            vmTest = with import (nixpkgs + "/nixos/lib/testing-python.nix") { inherit system; };

              (makeTest {
                name = "webhooker-test";
                nodes = {
                  server =
                    let
                      configEnv = pkgs.writeTextFile {
                        name = "envfile";
                        text = ''
                          HOOKER_TOKEN="global-testtoken"
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
                '';
              }).test;
          }
        );

      # Add dependencies that are only needed for development
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools gcc ];
          };
        });
    };
}
