#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
export GOMAXPROCS=1
go test -run TestRolePermissionSmoke -count=1 ./tests
