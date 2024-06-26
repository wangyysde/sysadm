#!/bin/sh

export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
VERSION=${1#"v"}
if [ -z "$VERSION" ]; then
  echo "Please specify the Kubernetes version: e.g."
  echo "./download-deps.sh v1.21.0"
  exit 1
fi

set -euo pipefail

# Find out all the replaced imports, make a list of them.
MODS=($(
  curl -sS "https://raw.githubusercontent.com/kubernetes/kubernetes/v${VERSION}/go.mod" |
    sed -n 's|.*k8s.io/\(.*\) => ./staging/src/k8s.io/.*|k8s.io/\1|p'
))

# Now add those similar replace statements in the local go.mod file, but first find the version that
# the Kubernetes is using for them.
for MOD in "${MODS[@]}"; do
  V=$(
    /c/Users/10288/go19/go/bin/go mod download -json "${MOD}@kubernetes-${VERSION}" |
      sed -n 's|.*"Version": "\(.*\)".*|\1|p'
  )

  /c/Users/10288/go19/go/bin/go mod edit "-replace=${MOD}=${MOD}@${V}"
done
echo "go get"
/c/Users/10288/go19/go/bin/go get "k8s.io/kubernetes@v${VERSION}"
echo "go mod download"
/c/Users/10288/go19/go/bin/go mod download