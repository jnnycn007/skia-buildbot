#!/bin/bash

VERSION_FILE="perf/VERSION.txt"
DATE_FILE="perf/DATE.txt"

# Function to exit with an error message
fail() {
  echo "ERROR: $1" >&2
  exit 1
}

# Check if the files exist
if [[ ! -f "$VERSION_FILE" ]]; then
  fail "VERSION.txt not found!"
fi
if [[ ! -f "$DATE_FILE" ]]; then
  fail "DATE.txt not found!"
fi

# Check if the files are not empty
if [[ ! -s "$VERSION_FILE" ]]; then
  fail "VERSION.txt is empty!"
fi
if [[ ! -s "$DATE_FILE" ]]; then
  fail "DATE.txt is empty!"
fi

# Check if the content looks like a git hash (40 hex chars) or "unknown" / "unversioned"
content=$(cat "$VERSION_FILE")
if [[ ! "$content" =~ ^[0-9a-f]{40}$ && \
     "$content" != "unknown" && \
     "$content" != "unversioned" ]]; then
  msg="VERSION.txt content '$content'"
  msg+=" doesn't look like a git hash or 'unknown' or 'unversioned'!"
  fail "$msg"
fi

# Check if the content looks like a date (YYYY-MM-DD)
date_content=$(cat "$DATE_FILE")
if [[ ! "$date_content" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
  fail "DATE.txt content '$date_content' doesn't look like a date (YYYY-MM-DD)!"
fi

echo "PASS"
exit 0