# aliasd

**aliasd** is a go based symlink manager for aliasing utilities that normally would be locally installed to their docker equivalents.

The goal of **aliasd** is to make the transition from CLI utility to docker command as simple as possible, without introducing unneeded complexity. **aliasd** achieves this by making aliases simple to define, 

## Adding a resource

```shell
aliasd add -f config.yaml
cat config.yaml | aliasd add -
aliasd add fpm -i skandyla/fpm -m '$(pwd):/data'
```
