#!/bin/bash

# Change to the root of the workspace
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." >/dev/null 2>&1 && pwd )"
cd "$DIR"

exec bazelisk run \
--ui_event_filters=-info,-warning,-stdout,-stderr \
--noshow_progress \
--noshow_loading_progress \
--show_result=0 \
//cmd/autoreview/mcp -- --cwd="$DIR"