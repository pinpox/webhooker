{ config, pkgs, lib, ... }:
with lib; let
  cfg = config.services.webhooker;

  settingsFormat = pkgs.formats.yaml { };
  configFile = settingsFormat.generate "config.yaml" cfg.settings;
in
{

  options.services.webhooker = {
    enable = mkEnableOption "webhooker service";

    package = mkPackageOption pkgs "webhooker" { };

    envFile = mkOption {
      type = types.str;
      default = null;
      example = "/var/secrets/webhooker/envfile";
      description = "Environment file to pass tokens to the service";
    };

    settings = mkOption {
      description = ''
        YAML settings for webhooker. See the
        [example configuration](https://github.com/pinpox/webhooker/blob/main/README.md#example)
        for more info.

        Secrets should be passed in by using `envFile`.
      '';
      default = {};
      type = types.submodule {
        freeformType = settingsFormat.type;

        options = {

          host = mkOption {
            type = types.str;
            default = "127.0.0.1";
            description = "Address to listen on";
            example = "0.0.0.0";
          };

          port = mkOption {
            type = types.str;
            default = "9999";
            description = "Port to listen on";
            example = "1111";
          };

          hooks = mkOption {
            type = types.attrsOf (types.submodule {
              options = {
                command = mkOption {
                  type = types.str;
                  description = "Command to execute";
                };
              };
            });
            default = { };
          };

        };
      };
    };
  };

  config = mkIf cfg.enable {
    systemd.services.webhooker = {
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      description = "webhooker service";
      environment.HOOKER_CONFIG = configFile;
      serviceConfig = {
        EnvironmentFile = lib.mkIf (cfg.envFile != null) "${cfg.envFile}";
        ExecStart = "${lib.getExe cfg.package}";
        Restart = "on-failure";
        RestartSec = "5s";
      };
    };
  };
}
