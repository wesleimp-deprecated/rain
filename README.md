# rain

Build and deliver docker images easily

# Usage

## Commands

```sh
$ rain --help

NAME:
   rain - Build and deliver docker images easily

USAGE:
   rain command [command options] [arguments...]

COMMANDS:
   init, i   Generates a .rain.yml file
   build, b  Builds the current project
   help, h   Shows a list of commands or help for one command
```

## Example config

```yml
# .rain.yml
project_name: 'rain'
version: '0.1.0'

# build application
builds:
  - name: 'go'
    command: 'go build -o rain main.go'

# build docker images
dockers:
  - dockerfile: 'Dockerfile'
    image_templates:
      - 'rainproj/rain:latest'
      - 'rainproj/rain:v{{ .Version }}'
    files:
      - glob: scripts/entrypoint.sh
      - glob: rain

# push to docker hub
pushes:
  - provider: hub
```