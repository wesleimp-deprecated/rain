# Pushes

Pushes section configures providers

**example**

```yml
pushes:
  - provider: hub
```

### Parameters reference

* `provider`: Image storage provider
* `env`: Set environment variables for the current push (optional)
* `name`: Identification name (optional)

**Avaliable providers**

* `hub`: Docker hub