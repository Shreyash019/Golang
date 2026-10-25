#!/usr/bin/env bash
# new-topic.sh creates topics/<NN-slug>/ from templates/topic/ as its own Go module,
# with the notes file copied from templates/topic-template.md.
#
# Usage: scripts/new-topic.sh 06-arrays-and-slices
set -euo pipefail

if [[ $# -ne 1 || ! $1 =~ ^[0-9]{2}-[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
	echo "usage: $0 NN-kebab-slug   (e.g. 06-arrays-and-slices)" >&2
	exit 2
fi

root=$(cd "$(dirname "$0")/.." && pwd)
dir="$root/topics/$1"
slug=${1#*-}

if [[ -e $dir ]]; then
	echo "$dir already exists: resume it instead of overwriting it" >&2
	exit 1
fi

# Module path = the GitHub repo path + topics/<NN-slug>, from the origin remote
# (https://github.com/owner/repo.git or git@github.com:owner/repo.git).
remote=$(git -C "$root" remote get-url origin)
repo=${remote%.git}
repo=${repo#https://}
repo=${repo#git@}
repo=${repo/://}
module="$repo/topics/$1"
template_module="github.com/Shreyash019/Golang/templates/topic"

cp -R "$root/templates/topic" "$dir"
cp "$root/templates/topic-template.md" "$dir/$slug.md"
find "$dir" -type f \( -name '*.go' -o -name go.mod \) -exec sed -i.bak "s#$template_module#$module#g" {} +
find "$dir" -name '*.bak' -delete

cd "$dir"
go mod edit -go="$(go env GOVERSION | sed -E 's/^go([0-9.]+).*/\1/')"
go vet ./...
go test ./...

echo "created topics/$1 (module $module)"
