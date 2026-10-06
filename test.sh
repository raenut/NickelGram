#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
TEST_DIR="$ROOT/local-test"
DB="$TEST_DIR/KoboReader.sqlite"
SOURCE=${KOBO_DB_SOURCE:-}
SQLITE=${SQLITE3:-$(command -v sqlite3 || true)}
GO=${GO:-go}
GOCACHE="$TEST_DIR/cache"
GOPATH="$TEST_DIR/gopath"
GOTOOLCHAIN=local
export GOCACHE GOPATH GOTOOLCHAIN

if [ "${1:-}" = setup ]; then
  if [ -z "$SOURCE" ] || [ ! -f "$SOURCE" ]; then
    printf '请设置 KOBO_DB_SOURCE=/path/to/KoboReader.sqlite 后运行 ./test.sh setup\n' >&2
    exit 1
  fi
  if [ -z "$SQLITE" ] || [ ! -x "$SQLITE" ]; then
    printf '找不到电脑可运行的 sqlite3；可设置 SQLITE3=/path/to/sqlite3。\n' >&2
    exit 1
  fi
  mkdir -p "$TEST_DIR/output" "$TEST_DIR/bin"
  cp "$SOURCE" "$DB"
  chmod 600 "$DB"
  "$SQLITE" "file:$DB?immutable=1" 'SELECT count(*) FROM Bookmark;' >/dev/null
  if [ ! -e "$TEST_DIR/footer.txt" ]; then : > "$TEST_DIR/footer.txt"; fi
  printf '本地数据库快照：%s\n' "$DB"
  exit 0
fi

if [ ! -f "$DB" ]; then
  printf '先运行 ./test.sh setup\n' >&2
  exit 1
fi
if [ -z "$SQLITE" ] || [ ! -x "$SQLITE" ]; then
  printf '找不到电脑可运行的 sqlite3：%s\n' "$SQLITE" >&2
  exit 1
fi
if [ ! -x "$TEST_DIR/bin/nickelgram-debug" ] || find "$ROOT/cmd/nickelgram" -name '*.go' ! -name '*_test.go' -newer "$TEST_DIR/bin/nickelgram-debug" | grep -q .; then
  if ! command -v "$GO" >/dev/null 2>&1; then
    printf '源码已更新或缺少本地程序；设置 GO=/path/to/go 后重试。\n' >&2
    exit 1
  fi
  (cd "$ROOT" && "$GO" build -o "$TEST_DIR/bin/nickelgram-debug" ./cmd/nickelgram)
fi

NICKELGRAM_TEST_DB="file:$DB?immutable=1"
NICKELGRAM_TEST_SQLITE="$SQLITE"
NICKELGRAM_TEST_OUTPUT="$TEST_DIR/output"
NICKELGRAM_TEST_FOOTER=$(cat "$TEST_DIR/footer.txt")
export NICKELGRAM_TEST_DB NICKELGRAM_TEST_SQLITE NICKELGRAM_TEST_OUTPUT NICKELGRAM_TEST_FOOTER NICKELGRAM_TEST_CONFIG
exec "$TEST_DIR/bin/nickelgram-debug" debug "$@"
