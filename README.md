# EC4 Config

_:warning: WORK IN PROGRESS - this cannot actually configure an EC4 yet :warning:_

Configuration as code for the [Faderfox EC4](https://faderfox.de/ec4.html).

## Usage

Your configuration must be piped to `stdin` as JSON. The expectation is that
configuration would be maintained with something like
[Jsonnet](https://jsonnet.org/), [Cue](https://cuelang.org/), or anything that
can produce JSON.

```sh
cat config.json | go run .
```

If you're using Jsonnet:

```sh
jsonnet config.jsonnet | go run .
```

## Configuration

Configuration is unified and validated with [schema.cue](schema.cue). This file
dictates object structure, field constraints, and default values.

When developing your configuration, the `cue vet` command is helpful for
determining if validation will pass.

```sh
cue vet schema.cue config.json
```

Jsonnet:

```sh
jsonnet config.jsonnet | cue vet schema.cue -
```

### Example

I maintain my own configuration here: https://gist.github.com/trotttrotttrott/fc4a74d8bd7d395cbf82431c467b77ef
