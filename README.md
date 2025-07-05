# EC4 Config

Configuration as code for the [Faderfox EC4](https://faderfox.de/ec4.html).

## Usage

Your configuration must be piped to `stdin` as JSON.

```sh
cat config.json | go run .
```

The expectation is that configuration would be maintained with something like
[Jsonnet](https://jsonnet.org/), [Cue](https://cuelang.org/), or anything that
can produce JSON.

If you're using Jsonnet:

```sh
jsonnet config.jsonnet | go run .
```

## Configuration

Configuration is validated with [schema.cue](schema.cue). This file dictates
object structure, field constraints, and also sets default values.

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

I maintain my own configuration in [config.jsonnet](config.jsonnet). Send it to
an EC4 with `make send`.

### Ableton Live

[live](./live) contains a remote script that is intended to work with the above
configuration.

See [documentation](https://help.ableton.com/hc/en-us/articles/209072009-Installing-third-party-remote-scripts)
for installation info. `make live-symlink` should do the trick.
