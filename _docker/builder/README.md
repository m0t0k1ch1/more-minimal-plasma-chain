## Build

``` sh
$ docker build -t mmpc-builder .
```

## Run

```
$ docker run --name mmpc-builder mmpc-builder
```

## Copy binaries from container to local

``` sh
$ docker cp mmpc-builder:/go/bin/more-minimal-plasma ./
$ docker cp mmpc-builder:/go/bin/plasma ./
```
