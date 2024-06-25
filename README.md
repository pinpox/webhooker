# Webhooker

Listens on a specified port for (authenticated) incoming webhooks or requests
and runs a specified command.

<p align="center">
  <a href="https://github.com/pinpox/webhooker/actions/workflows/go.yml">
      <img src="https://github.com/pinpox/webhooker/actions/workflows/go.yml/badge.svg" alt="Go Unit Tests" />
  </a>

  <a href="https://github.com/pinpox/webhooker/actions/workflows/nix.yml">
      <img src="https://github.com/pinpox/webhooker/actions/workflows/nix.yml/badge.svg" alt="Nix Integration Tests" />
  </a>

  <a href="https://buildbot.thalheim.io/#/projects/22">
      <img src="https://img.shields.io/badge/Buildbot-Nix-blue?link=https%3A%2F%2Fbuildbot.thalheim.io%2F%23%2Fprojects%2F22" alt="Buildbot Nix" />
  </a>
</p>

## Configuration

The configuration should mostly be self-explanatory, below is an example. Hooks
will be authenticated using the `Token` header, if either the hook-specific
`Token` or the `globalToken` match, the command below will be run using 
`/bin/sh -c`. 

Individual hooks will be matched by path using the name in the `hooks` map.

### Example

```yaml
host: "127.0.0.1"
port: "9999"
globalToken: "the-defaulttoken"
hooks:
  hello:
    command: "echo 'hello world'"
    token: "token1"
  ping:
    command: "ping -c 4 google.de"
    token: "token2"
  list:
    Command: "ls -l"
```

With the configuration above:

- `curl 127.0.0.1:9999/hello -H 'Token: token1'` would execute `echo 'hello world'`
- `curl 127.0.0.1:9999/ping -H 'Token: the-defaulttoken'` would execute `ping -c 4 google.de`
- `curl 127.0.0.1:9999/list -H 'Token: the-defaulttoken'` would execute `ls -l`

#### Environment Variables

The path to the configuration should be specified using the `HOOKER_CONFIG`
environment variable. 

To avoid having tokens in the configuration, it is also possible to specify them
via environment variables

- The `globalToken` can be specified using the `HOOKER_TOKEN` environment variable.
- Tokens for individual hooks will be searched in environment variables with the
  same name and the `HOOKER_TOKEN_` prefix (`HOOKER_TOKEN_<name of the token in uppercase>`).
  e.g.: for the last hook of the configuration example above the token would be
  specified in `HOOKER_TOKEN_LIST`.

<!-- 
## Use-Cases 
TODO add examples, e.g.
- building static blog
- rebuilding nixos host
-->

## NixOS Module

The `flake.nix` file includes a module for NixOS users. Import and use as shown
below:

```nix
# Add to flake inputs
inputs.webhooker.url = "github:pinpox/webhooker";

# Import the module in your configuration.nix
imports = [ self.inputs.webhooker.nixosModules.webhooker ];

services.webhooker = {
  enable = true;
  envFile = "/path/to/your/envfile";
  settings = {
    host = "0.0.0.0";
    port = "8888";
    hooks = {
      hello.command = "echo 'hello'";
      world.command = "echo 'world'";
    };
  };
};
```

For develoment the flake also provides NixOS tests. Use `nix flake check -v` to
run them or `nix run .\#checks.x86_64-linux.vmTest.driverInteractive` for
interactive testing and debugging.
