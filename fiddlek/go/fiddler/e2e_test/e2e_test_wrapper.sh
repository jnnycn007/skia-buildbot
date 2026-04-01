#!/bin/bash
set -e

# Bazel provides these paths as runfiles.
./fiddlek/load_fiddler_container-base.sh
./fiddlek/go/fiddler/e2e_test/e2e_test_/e2e_test
