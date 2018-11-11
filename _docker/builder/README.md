``` sh
$ docker build -t mmpc-builder .
$ docker run --name mmpc-builder mmpc-builder
$ docker cp mmpc-builder:/go/bin/more-minimal-plasma ./
$ docker cp mmpc-builder:/go/bin/plasma ./
```
