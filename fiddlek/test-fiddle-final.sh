#!/bin/bash
set -e

# This script is used to build and test the final fiddler backend image,
# including any local changes in this repo, in order to experiment iteratively.
# You may pass the following flags:
#
# --skia-repo-path <path>: Path to a local checkout of the Skia repo. Set this
#     to test against any local changes to that repo. If unset, this script
#     will clone a fresh copy.
#
# --port <port>: Port number on which fiddler should run. Set this if the
#     default conflicts with other services on your machine.

# Configuration.
SKIA_REPO_PATH=""
PORT="8080"
TEMP_DIR="$(mktemp -d)"

while [[ $# -gt 0 ]]; do
    case $1 in
        --skia-repo-path)
            SKIA_REPO_PATH="$2"
            shift 2
            ;;
        --port)
            PORT="$2"
            shift 2
            ;;
        *)
            echo "Unknown argument: $1"
            echo "Usage: $0 [--skia-repo-path <path>] [--port <port>]"
            exit 1
            ;;
    esac
done
if [[ -z "$SKIA_REPO_PATH" ]]; then
    pushd $TEMP_DIR
    git clone --depth 1 https://skia.googlesource.com/skia.git
    popd
    SKIA_REPO_PATH="$TEMP_DIR/skia"
fi

echo "Using Skia repo at: $SKIA_REPO_PATH"

# 1. Build the base image.
echo "--- Building fiddler-base image ---"
bazelisk run //fiddlek:load_fiddler_container-base
BASE_IMAGE="gcr.io/skia-public/fiddler-base:latest"

# 2. Create a modified Dockerfile for the final image
DOCKERFILE_ORIG="$SKIA_REPO_PATH/infra/fiddler-backend/Dockerfile"
DOCKERFILE_LOCAL="$TEMP_DIR/Dockerfile"
echo "--- Creating modified Dockerfile at $DOCKERFILE_LOCAL ---"
# Replace all occurrences of the base image with our local one.
sed "s|gcr.io/skia-public/fiddler-base@[^ ]*|$BASE_IMAGE|g" $DOCKERFILE_ORIG > $DOCKERFILE_LOCAL

# 3. Build the final fiddler image
echo "--- Building final fiddler image ---"
# We run docker build from the Skia repo root because the Dockerfile expects to COPY .
cd $SKIA_REPO_PATH
docker build -t fiddler:local -f $DOCKERFILE_LOCAL .
echo "--- Fiddler image built successfully as fiddler:local ---"

echo "--- Running reproduction test ---"

# 4. Start the container.
CONTAINER_NAME="fiddler_repro"
echo "Removing any existing container..."
docker rm -f $CONTAINER_NAME || true
echo "Starting container on port $PORT..."
docker run --name $CONTAINER_NAME -d -p "$PORT":"$PORT" --entrypoint /bin/bash fiddler:local -c "sleep infinity"

# 5. Start fiddler inside the container.
docker exec -d -u skia $CONTAINER_NAME bash -c "export PATH=/usr/local/bin:\$PATH && /usr/local/bin/fiddler --local --fiddle_root=/tmp --checkout=/tmp/skia --port=:$PORT > /tmp/fiddler.log 2>&1"
echo "Waiting for fiddler to start..."
MAX_RETRIES=30
COUNT=0
while ! curl -s http://localhost:$PORT/ > /dev/null; do
    sleep 1
    COUNT=$((COUNT+1))
    if [ $COUNT -ge $MAX_RETRIES ]; then
        echo "Fiddler failed to start. Logs:"
        docker exec $CONTAINER_NAME cat /tmp/fiddler.log
        exit 1
    fi
done
echo "Fiddler is up!"

# 6. Send a request to build and run a fiddle.
#    Edit the code below if desired.
cat <<EOF > $TEMP_DIR/draw.cpp
/////////////////////////////////////////
// The fiddle server adds this header. //
/////////////////////////////////////////
#include "skia.h"
#include "fiddle_main.h"
DrawOptions GetDrawOptions() {
  static const char *path = 0; // Either a string, or 0.
  return DrawOptions(256, 256, true, true, false, false, false, false, false, path, skgpu::Mipmapped::kNo, 64, 64, 0, skgpu::Mipmapped::kNo);
}
/////////////////////////////////////////

void draw(SkCanvas* canvas) {
    SkPaint p;
    p.setColor(SK_ColorRED);
    p.setAntiAlias(true);
    p.setStyle(SkPaint::kStroke_Style);
    p.setStrokeWidth(10);

    canvas->drawLine(20, 20, 100, 100, p);
}
EOF
CODE=$(jq -Rs '.' $TEMP_DIR/draw.cpp)
cat <<EOF > $TEMP_DIR/request.json
{
  "code": $CODE,
  "options": {
    "width": 128,
    "height": 128,
    "textOnly": true
  }
}
EOF

echo "Sending request..."
curl -X POST -H "Content-Type: application/json" -d @$TEMP_DIR/request.json http://localhost:$PORT/run > $TEMP_DIR/response.json

echo "Response received:"
cat $TEMP_DIR/response.json | jq .

COMPILE_ERRORS=$(jq -r '.compile.errors // empty' $TEMP_DIR/response.json)
EXECUTE_ERRORS=$(jq -r '.execute.errors // empty' $TEMP_DIR/response.json)

if [ -n "$COMPILE_ERRORS" ] || [ -n "$EXECUTE_ERRORS" ]; then
    echo "--- Fiddler Logs ---"
    docker exec $CONTAINER_NAME cat /tmp/fiddler.log
fi

docker rm -f $CONTAINER_NAME
rm -rf $TEMP_DIR
