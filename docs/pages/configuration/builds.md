# Builds

The builds section is an optional section and can be used to generate binaries, dlls, etc..

**example**

```yaml
builds:
  - name: 'go'
  - command: go build -o app .
```

### Parameter reference

* `command`: command that should be executed
* `name`: Identification name (optional)
* `env`: Extra envs for this specific build (optional)
* `dir`: This field should be set if the build is not in the root (default `.`)
* `skip`: Skipt this specific build (default `false`)
