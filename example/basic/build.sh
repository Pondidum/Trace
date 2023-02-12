#!/bin/sh

export "PATH=${PATH}:../../"

# start of build
export TRACEPARENT="$(trace start "basic-build")"


clone=$(trace span start "clone_artifacts")
  # git pull
  sleep 1s
  # install tools
trace span finish "${clone}"

docker=$(trace span start "docker")

  pulls=$(trace span start "pull" "${docker}")
    # docker pull
    sleep 2s
  trace span finish "${pulls}"

  build=$(trace span start "build" "${docker}")
    # docker builds
    sleep 2s
  trace span finish "${build}"


  push=$(trace span start "push" "${docker}")
    # docker push
    sleep 2s
  trace span finish "${push}"

trace span finish "${docker}"

# end of build
trace finish
