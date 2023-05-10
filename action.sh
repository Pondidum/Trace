#!/bin/sh

set -eu

version="${1:-""}"

if [ -z "${version}" ]; then
  version=$(curl -sSL https://api.github.com/repos/Pondidum/trace/releases/latest | sed -n 's/.*tag_name.*"\(.*\)".*/\1/p')
fi

binary_dir="${RUNNER_TOOL_CACHE}/trace/${version}"
binary_path="${binary_dir}/trace"

if ! [ -f "${binary_path}" ]; then
  echo "Downloading Trace ${version}"
  mkdir -p "${binary_dir}"
  curl -sSL "https://github.com/Pondidum/trace/releases/download/${version}/trace" -o "${binary_path}"

  echo "Done"
else
  echo "Trace ${version} found in cache"
fi

chmod +x "${binary_path}"

${binary_path} version

echo "${binary_dir}" >> "${GITHUB_PATH}"
echo "absolute_path=${binary_path}" >> "${GITHUB_OUTPUT}"
