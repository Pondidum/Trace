#!/bin/sh

export "PATH=${PATH}:../../"

# start of build
export TRACEPARENT="$(trace start "basic-build")"

# not needed normally!
# just done so we can query the trace after this has run in the github actions build
echo "${TRACEPARENT}" > .traceid

clone=$(trace group start "clone_artifacts")
  # git pull
  sleep 1s
  # install tools
trace group finish "${clone}"

docker=$(trace group start "docker")

  pulls=$(trace group start "pull" "${docker}")

    trace task "$pulls" -- docker pull alpine:latest

  trace group finish "${pulls}"

  build=$(trace group start "build" "${docker}")

    trace task "$build" -- docker build . --no-cache
    trace task "$build" -- sleep 5s

  trace group finish "${build}"


  push=$(trace group start "push" "${docker}")

    trace task "$push" --  sleep 2s &
    trace task "$push" --  sleep 10s &
    trace task "$push" --  sleep 5s &

    wait

  trace group finish "${push}"

trace group finish "${docker}"

# end of build
trace finish
