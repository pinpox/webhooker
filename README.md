# Webhooker

Listens on a specified port for (authenticated) incoming webhooks or requests
and runs a specified command.

## Configuration

The configuration should mostly be self-explanatory, below is an example. Hooks
will be authenticated using the `Token` header, if either the hook-specific
`Token` or the `GlobalToken` match, the command below will be run using 
`/bin/sh -c`. 

Individual hooks will be matched by path using the name in the `Hooks` map.

### Example

```yaml
Host: "127.0.0.1"
Port: "9999"
GlobalToken: "the-defaulttoken"
Hooks:
  hello:
    Command: "echo 'hello world'"
    Token: "token1"
  ping:
    Command: "ping -c 4 google.de"
    Token: "token2"
  list:
    Command: "ls -l"
```

With the configuration above:

- `curl 127.0.0.1:9999/hello -H 'Token: token1'` would execute `echo 'hello world'`
- `curl 127.0.0.1:9999/ping -H 'Token: the-defaulttoken'` would execute `ping -c 4 google.de`
- `curl 127.0.0.1:9999/list -H 'Token: the-defaulttoken'` would execute `ls -l`

The path to the configuration can be specified using the `HOOKER_CONFIG`
environment variable. The `GlobalToken` can alternatively specified using the
`HOOKER_TOKEN` environment variable.

