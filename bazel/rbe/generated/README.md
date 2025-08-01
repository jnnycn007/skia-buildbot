# Bazel toolchain configuration for RBE

---

**DO NOT EDIT THIS DIRECTORY BY HAND.**

All files in this directory (excluding this file) are generated with the `rbe_configs_gen` CLI
tool. Keep reading for details.

---

This directory contains a Bazel toolchain configuration for RBE. It is generated with the
`rbe_configs_gen` CLI tool from the
[bazel-toolchains](https://github.com/bazelbuild/bazel-toolchains) repository.

This directory is referenced from `//.bazelrc`.

## Upgrading to a new Bazel version or rebuilding the toolchain container

Take the following steps to upgrade to a new Bazel version or if rolling out a new version of the
toolchain, e.g. rebuilding to include updated tools.

### Step 1 - Making changes

#### Updating Bazel Version

If necessary, update the `//.bazelversion` file with the new Bazel version. This file is read by
[Bazelisk](https://github.com/bazelbuild/bazelisk) (for those engineers who use `bazelisk`
as a replacement for `bazel`).

If not updating Bazel, merely note what version is there.

#### Rebuild Toolchain Container

If necessary, make changes to `//bazel/rbe/gce_linux_container/Dockerfile` and commit those.
Use Louhi's Build infra-rbe-linux flow to create a new image.

Make note of the image's SHA256 hash, found in the Louhi page
[(example)](https://louhi.corp.goog/6316342352543744/execution-detail/5660155978317824).

If not rebuilding the container, merely note what SHA256 hash is specified in the `container-image`
exec_property of the platform defined in `//bazel/rbe/generated/config/BUILD`
([example](https://skia.googlesource.com/buildbot/+/bb3604fd9a57bb20d799341b50f616af9e0062d4/bazel/rbe/generated/config/BUILD#43)).

### Step 2 - Generate Bazel Files

Note: We frequently skip this step when making changes to the RBE toolchain container or when
upgrading to a new Bazel version, and things usually continue to work. If you wish to skip
regenerating these files, just update the `container-image` exec_property of the platform defined
in `//bazel/rbe/generated/config/BUILD`
([example](https://skia.googlesource.com/buildbot/+/bb3604fd9a57bb20d799341b50f616af9e0062d4/bazel/rbe/generated/config/BUILD#43)).

Regenerate the `//bazel/rbe/generated` directory with the `rbe_configs_gen` CLI tool:

```
# Replace the <PLACEHOLDERS> as needed.
$ bazelisk run //:rbe_configs_gen \
      --bazel_version=<BAZEL VERSION> \
      --toolchain_container=gcr.io/skia-public/infra-rbe-linux@sha256:<HASH OF MOST RECENT IMAGE> \
      --output_src_root=<PATH TO REPOSITORY CHECKOUT> \
      --output_config_path=bazel/rbe/generated \
      --exec_os=linux \
      --target_os=linux \
      --generate_java_configs=false
```

Example:

```
$ bazelisk run //:rbe_configs_gen -- \
      --bazel_version=$(cat .bazelversion) \
      --toolchain_container=gcr.io/skia-public/infra-rbe-linux@sha256:9d565deca2ec317a4c26403ba9d14cf4d3ed083632cd24870155db292eb4de6b \
      --output_src_root=$(pwd) \
      --output_config_path=bazel/rbe/generated \
      --exec_os=linux \
      --target_os=linux \
      --generate_java_configs=false
```

If `rbe_configs_gen` fails, try deleting all files under `//bazel/rbe/generated` (except for this
file) and re-run `rbe_configs_gen`.

### Step 3 - Cleanup

Run `bazel run //:buildifier` and fix any linting errors on the generated files (e.g. move load
statements at the top of the file, etc.)

### Step 4 - Manual Updates

If updating the Bazel version, update the
[bazel-toolchains](https://github.com/bazelbuild/bazel-toolchains) repository version imported from
`//WORKSPACE` to match the new Bazel version.

If updating the image version, the generated changes should have properly updated
`//bazel/rbe/generated/config/BUILD` to refer to the new image, so there is nothing else to do.

### Step 5 - Test out the changes

Try running various bazel commands with `--config remote` to make use of the new image.

When uploading the change as a CL, any tasks with the `-RBE` suffix will automatically use the
image specified in `//bazel/rbe/generated/config/BUILD`.
