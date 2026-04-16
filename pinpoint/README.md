# Pinpoint Developer Documentation

## Run a try job locally

### 1. Run a local temporal dev server.

If Temporal CLI is not installed, install it from the official
[web site](https://temporal.io/setup/install-temporal-cli).

Note: don't forget to add Temporal CLI to your PATH. One option is to move
Temporal CLI to `/usr/local/bin` by running:

```
sudo mv temporal /usr/local/bin
```

Run Temporal local dev server.

```
temporal server start-dev
```

Provide a database file if you need persistent workflows.

```
temporal server start-dev --db-filename=temporal-db.db
```

Temporal Web UI is available on [localhost:8233](http://localhost:8233).

### 2. Run a temporal worker locally.

Feel free to use any `taskQueue`. In the example below it is `pptq`.

```
bazelisk run //pinpoint/go/workflows/worker -- \
  --taskQueue=pptq \
  --local
```

### 3. Create a try job workflow.

Make sure the task queue matches with the worker task queue.

```
bazelisk run //pinpoint/go/workflows/sample -- \
  --taskQueue=pptq \
  --pairwise \
  --configuration=win-11-perf \
  --benchmark=speedometer3 \
  --story=Speedometer3 \
  --start-git-hash=b2d27b144e4e4c5661bafc08f7b8494797f6ee1a \
  --end-git-hash=95b3180e9724995eb6d5a85ac3c93140e4506f7e
```
