
# king

> Protobuf RPC using JSON.

For the general public we recommend better public & documented projects like [Twitch Twirp](https://github.com/twitchtv/twirp) or [Improbable Engineering GRPC Web](https://github.com/improbable-eng/grpc-web).


### Install

```shell
curl https://tools.altipla.consulting/bin/king > ~/bin/king && chmod +x ~/bin/king
```


### Usage

```shell
king auth api.altipla.consulting FOO_TOKEN
king auth api.altipla.consulting
king call foo.bar.FooService.List project=shs foo.bar=3 foo.baz=one foo.baz=two numeric:=3 boolean=true
```
