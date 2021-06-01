# aliasd

**aliasd** is a go based symlink manager for aliasing utilities that normally would be locally installed to their docker equivalents.

The goal of **aliasd** is to make the transition from CLI utility to docker command as simple as possible, without introducing unneeded complexity. **aliasd** achieves this by making aliases simple to define as yaml resources, with as few required changes to workflow as possible.

## Adding a proxy

```shell
aliasd add -f config.yaml
```

config.yaml:

```yaml
resources:
  fpm:
    image: skandyla/fpm
    volumeMounts:
      - mountPath: /data
        hostPath: $(pwd)
```

## Using a proxy

The hostPath above is specified as `$(pwd)` in order to generate a find and replace of all forwarded args to a proxy command. In other words, if your current working directory is `/tmp/test`, and you run the following command:

```shell
fpm -s dir -t rpm -n "aliasd" -v 0.1 -p $(pwd) -C $(pwd)/rpmbuild ./
```

`./` will be passed as part of the command to `docker run`, whereas `$(pwd)` will evaluate to `/tmp/test` and be replaced with `/data` per the config file. This isn't considered too diruptive, as the majority of the tools in most workflows (CMake, Ninja, Make) use absolute file paths by default (which in turn means gcc/g++ should be easy to support), and modifying scripts to use `$(pwd)` instead of `./` where applicable or alternatively changing directory structure to be compatible with both `docker run` and locally installed versions isn't too disruptive. Additional configuration is possible, but as **aliasd** grows new requirements will become apparent.

Alternatively, the above command could be achieved with the following equivalent command:

```shell
aliasd execute -n fpm -s dir -t rpm -n "aliasd" -v 0.1 -p $(pwd) -C $(pwd)/rpmbuild ./
```

which may be suitable in the event that `~/.aliasd/bin` does not have `$PATH` precedence.

## Upcoming Features

- client/server aliasd implementation
  - configuration storage in system locations
- kubernetes integration
  - agent configuration where aliasd serves as a proxy for kubernetes API request to create containers capable of executing specific commands
- api versioning
