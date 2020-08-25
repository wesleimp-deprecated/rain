# rain

Build docker images easily

# Usage

## Commands

```sh
$ rain --help

NAME:
   rain - Deploy docker images and keep yout environments updated

USAGE:
   rain command [command options] [arguments...]

COMMANDS:
   init, i   Generates a .rain.yml file
   build, b  Builds the current project
   help, h   Shows a list of commands or help for one command
```

## Example config

```yml
project_name: 'app'

version: '0.1.0'

builds:
  - name: 'go'
    command: 'go build -o app main.go'

dockers:
  - dockerfile: 'Dockerfile'
    image_templates:
      - 'app:latest'
      - 'app:v{{ .Version }}'
    files:
      - glob: scripts/entrypoint.sh
      - glob: app
```