# Webhooker

Listens on a specified port for (authenticated) incoming webhooks or requests
and runs a specified command.

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

- rebuilding static blog
- reibliding nixos host

## NixOS Module
-->



