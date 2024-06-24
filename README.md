# Webhooker

Listens on a specified port for (authenticated) incoming webhooks or requests
and runs a specified command.

## Configuration

The configuration should mostly be self-explanatory, below is an example. Hooks
will be authenticated using the `Token` header, if either the hook-specific
`Token` or the `GlobalToken` match, the command below will be run using 
`/bin/sh -c`. 

```yaml
Host: "127.0.0.1"
Port: "9999"
GlobalToken: "the-defaulttoken"
Hooks:
  hook1:
    Command: "command1"
    Token: "token1"
  hook2:
    Token: "token2"
  hook3:
    Command: "command3"
```

The path to the configuration can be specified using the `HOOKER_CONFIG`
environment variable. The `GlobalToken` can alternatively specified using the
`HOOKER_TOKEN` environment variable.

