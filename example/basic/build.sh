#!/bin/sh

export "PATH=${PATH}:../../"

# start of build
export TRACEPARENT="$(trace start "basic-build")"


clone=$(trace group start "clone_artifacts")
  # git pull
  sleep 1s
  # install tools
trace group finish "${clone}"

docker=$(trace group start "docker")

  pulls=$(trace group start "pull" "${docker}")
    # docker pull
    sleep 2s
  trace group finish "${pulls}"

  build=$(trace group start "build" "${docker}")
    # docker builds
    sleep 2s
  trace group finish "${build}"


  push=$(trace group start "push" "${docker}")
    # docker push
    sleep 2s
  trace group finish "${push}"

trace group finish "${docker}"

# end of build
trace finish
