# EC4 Config

:warning: WORK IN PROGRESS :warning:

Configuration as code for the [Faderfox EC4](https://faderfox.de/ec4.html).

## Usage

Your configuration must be piped to `stdin` as JSON. The expectation is that
configuration would be maintained with something like
[jsonnet](https://jsonnet.org/), [Cue](https://cuelang.org/), or anything that
can produce JSON.

```sh
cat config.json | go run .
```

Or if you're using something like Jsonnet:

```sh
jsonnet config.jsonnet | go run .
```

## Configuration

Configuration is validated against [schema.cue](schema.cue).

When developing your configuration, the `cue vet` command is helpful for
determining if validation will pass.

```sh
cue vet schema.cue config.json
```

Jsonnet:

```sh
jsonnet config.jsonnet | cue vet schema.cue -
```
