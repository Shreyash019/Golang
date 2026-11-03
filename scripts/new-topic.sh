#!/usr/bin/env bash
# new-topic.sh creates topics/<NN-slug>/ from templates/topic/ as its own Go module:
#   README.md                 a short index (from templates/topic/README.md)
#   theory/README.md          the notes (from templates/topic-template.md)
#   theory/<slug>-qa.md       the Q&A log (from templates/topic-qa-template.md)
# {{NN}}, {{SLUG}} and {{TITLE}} are filled in. TITLE is derived from the slug
# ("arrays-and-slices" → "Arrays & slices"); fix acronyms by hand (REST, API).
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
mkdir -p "$dir/theory"
cp "$root/templates/topic-template.md" "$dir/theory/README.md"
cp "$root/templates/topic-qa-template.md" "$dir/theory/$slug-qa.md"

nn=${1%%-*}
title=$(echo "$slug" | sed -e 's/-and-/ \& /g' -e 's/-/ /g')
title="$(echo "${title:0:1}" | tr '[:lower:]' '[:upper:]')${title:1}"
for f in "$dir/README.md" "$dir/theory/README.md" "$dir/theory/$slug-qa.md"; do
	sed -i.bak -e "s#{{NN}}#$nn#g" -e "s#{{SLUG}}#$slug#g" -e "s#{{TITLE}}#${title//&/\\&}#g" "$f"
done
find "$dir" -type f \( -name '*.go' -o -name go.mod \) -exec sed -i.bak "s#$template_module#$module#g" {} +
find "$dir" -name '*.bak' -delete

cd "$dir"
go mod edit -go="$(go env GOVERSION | sed -E 's/^go([0-9.]+).*/\1/')"
go vet ./...
go test ./...

echo "created topics/$1 (module $module)"
