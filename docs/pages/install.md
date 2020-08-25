# Install

### Docker

```sh
$ docker run --rm --privileged rainproj/rain build
```

### Compiling from source

**clone**

```sh
$ git clone git@github.com:rainproj/rain.git

$ cd rain
```

**download dependencies**

```sh
$ go mod download
```

**build**

```sh
$ go build -o rain main.go
```

**verify it works**

```sh
$ rain --help
```
