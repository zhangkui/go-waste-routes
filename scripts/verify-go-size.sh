#!/usr/bin/env sh
set -eu
files=$(find cmd internal -name '*.go' ! -name '*_test.go' | sort)
count=$(printf '%s\n' "$files" | sed '/^$/d' | wc -l)
lines=0
for file in $files; do
  file_lines=$(wc -l < "$file")
  lines=$((lines + file_lines))
done
printf 'files:\n%s\ncount:%s\nlines:%s\n' "$files" "$count" "$lines"
