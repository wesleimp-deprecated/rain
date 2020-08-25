# Dockers

Dockers section configures the docker images

**example**

```yml
dockers:
  - dockerfile: 'Dockerfile'
    image_templates:
      - 'myorg/app:latest'
    files:
      - glob: scripts/myscript.sh
```

It's possible to use evironment variables in `image_templates` and `build_flags_templates`

```yaml
env:
  - "FOO=BAR"

dockers:
  - dockerfile: 'Dockerfile'
    image_templates:
      - 'myorg/app:latest'
      - 'myorg/app:{{ .FOO }}'
    files:
      - glob: scripts/myscript.sh
```

### Parameter reference

* `dockerfile`: Dockerfile file
* `image_templates`: Image templates is a list of templates.
* `files`: Files used inside Dockerfile (optional)
* `build_flag_templates`: Image flags (optional)
* `skip_push`: Skip push to hub (default `false`)

