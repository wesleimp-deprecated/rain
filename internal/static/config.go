package static

// InitConfig for command init
var InitConfig = `# .rain.yml
project_name: 'app'

version: '0.1.0'

dockers:
  - dockerfile: 'Dockerfile'
    image_templates:
      - 'app:latest'
      - 'app:v{{ .Version }}'
`
