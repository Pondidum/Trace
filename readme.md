# Trace

*Trace your scripts in an opinionated way*

## General Usage

1. start a trace `export TRACEPARENT="$(trace start "some-name")`
2. start groups: `group=$(trace group start "some-name")`
3. run processes inside a group: `trace task "${group}" -- some command here`
4. finish the group `trace group finish "${group}"`
5. finish the trace `trace finish`

### Configuration

By default, traces are send to `localhost:4317` using `gRPC` in `OTLP` protocol.  This can be overridden using the `OTEL_EXPORTER_OTLP_ENDPOINT` or `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` environment variables.  Additionally, headers can be sent to the OTLP backend by setting the `OTEL_EXPORTER_OTLP_HEADERS` value (a csv of `key=value` pairs.)

For example, to send traces to Honeycomb's gRPC endpoint:

```shell
export OTEL_EXPORTER_OTLP_ENDPOINT=api.honeycomb.io:443
export OTEL_EXPORTER_OTLP_HEADERS="x-honeycomb-team=your-api-key"
```

### Grouping

Groups are optional, and can be as deeply nested as you like.  If you don't explicitly set the parent of a group (second argument), it will default to the span in the `TRACEPARENT` environment variable.

```shell
docker=$(trace group start "docker")
  push=$(trace group start "push" "${docker}")

    trace task "${push}" --  docker push --quiet "${mycontainer}:${mytag}"
    trace task "${push}" --  docker push --quiet "${mycontainer}:latest"

  trace group finish "${push}"
trace group finish "${docker}"
```

### Parallelism

Everything supports parallelism.  For example, these 3 tasks run in parallel, all parented to the same group:


```shell
push=$(trace group start "push")

  trace task "$push" --  sleep 2s &
  trace task "$push" --  sleep 10s &
  trace task "$push" --  sleep 5s &

  wait

trace group finish "${push}"
```

The only action to be careful with is adding attributes to a single group in parallel:

```shell
push=$(trace group start "push")

  ## the attributes on the "push" group are not guaranteed now:
  trace attr "${push}" first=true second=nope third=false &
  trace attr "${push}" first=one second=true third=false &
  trace attr "${push}" first=two second=two third=true &

trace group finish "${push}"
```

## Examples

For examples, start the `docker-compose.yml` file in the repo root, then go to one of the example directories and try a script:

```shell
docker-compose up -d

cd example/basic

./build.sh
```

