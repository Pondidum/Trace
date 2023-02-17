# Design

## Notes

Due to needing to store some information about a trace, and then use that later to complete the trace, we need to have some temporary storage.  Something like a temp directory path, and a file per span, which could be named after tha span ids:

```sh
> set -x TRACEPARENT (trace generate "build")
# 00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01

> set -x span (trace span start "download")
# 00-7107538ee3f6bc77ada1b2d34a412e1d-6eb172bfeefb767c-01

> ls /tmp/tracing/
# 00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01
# 00-7107538ee3f6bc77ada1b2d34a412e1d-6eb172bfeefb767c-01

> cat /tmp/tracing/00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01
name=build
start=<some date format>
```

By default the path would be `$TMP/tracing` but overridable with an environment variable (`set -x TRACING_STORE "/tmp/tracing/build-5"`)


## Commands

### Create a trace ID

- Just generates a valid traceID so that things can be attached to it.
- Sent later with [finish] command.

```sh
> trace generate
# 00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01
```

```sh
trace generate --export
> set -x TRACEPARENT "00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01"

# fish shell
> trace generate --export | source

# bash
> eval $(trace generate --export)

# general
export TRACEPARENT=$(trace generate)
```

Flags:

- `--export` generate a statement to export the `TRACEPARENT` env var
- `--shell <fish|bash|sh>` force shell detection

### Finish the trace

- Given a trace id, sends a span to the backend

```sh
> trace finish "00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01"
```

```sh
> set -x TRACEPARENT "00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01"
> trace finish
```


### Start a span

- takes 1 or 2 arguments: `name` and `trace_parent`
- if `trace_parent` isn't specified, defaults to `TRACEPARENT`
- prints a traceid which can be captured in a variable or exported
- if there is no traceid, it _will_ complain

```sh
span=$(trace span start "fetch_tools")
# ...
trace span finish "$span"
```

Child spans:

```sh
span=$(trace span start "fetch_tools")

  child=$(trace span start "jq" "$span")
    # download jq
  trace span finish "$child"


  child=$(trace span start "docker" "$span")
    # download docker
  trace span finish "$child"

trace span finish "$span"
```

### Capture a process

- takes 1 or 2 arguments: `name` and `trace_parent`
- if `trace_parent` isn't specified, defaults to `TRACEPARENT`
- everything after the `--` is passed to `$SHELL` for running
- captures exit code into the span, and then exits with same exit code

```sh
trace exec "docker_pull" -- docker pull alpine:3.18 --no-cache
```

```sh
trace exec "docker_pull" "$span" -- docker pull alpine:3.18 --no-cache
```


## Complete Example

```shell
export TRACEPARENT=$(trace generate)

  span=$(trace span start "pull_containers")

    trace exec "alpine" "$span" -- docker pull alpine:3.18 &
    trace exec "main" "$span" -- docker pull $app:main &
    trace exec "branch" "$span" -- docker pull $app:$branch &

    wait

  trace span finish "$span"

  span=$(trace span start "build")

    trace exec "docker_build" "$span" -- docker build -t $app:$branch .

  trace span finish "$span"

trace finish
