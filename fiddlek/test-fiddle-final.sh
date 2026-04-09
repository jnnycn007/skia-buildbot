#!/bin/bash
set -e

# This script is used to build and test the final fiddler backend image,
# including any local changes in this repo, in order to experiment iteratively.
# You may set the following variables:
#
# SKIA_REPO_PATH: Path to a local checkout of the Skia repo. Set this to test
#     against any local changes to that repo. If unset, this script will clone
#     a fresh copy.
#
# PORT: Port number on which fiddler should run. Set this if the default
#     conflicts with other services on your machine.

# Configuration.
TEMP_DIR="$(mktemp -d)"
PORT="${PORT:=8080}"
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
cat <<EOF > $TEMP_DIR/request.json
{
  "code": "#include \"include/core/SkCanvas.h\"\n#include \"tools/fiddle/fiddle_main.h\"\nDrawOptions GetDrawOptions() {\n  return DrawOptions(128, 128, true, true, true, true, true, false, false, nullptr, skgpu::Mipmapped::kNo, 64, 64, 0, skgpu::Mipmapped::kNo);\n}\nvoid draw(SkCanvas* canvas) { canvas->clear(SkColorSetRGB(0, 255, 0)); }",
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
