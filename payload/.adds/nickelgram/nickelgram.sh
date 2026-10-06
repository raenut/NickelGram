#!/bin/sh
ROOT=${0%/*}
if [ "$ROOT" = "$0" ]; then ROOT=.; fi

case "$(uname -m)" in
 arm*) ;;
 *) printf '%s\n' '此发布包只支持 32 位 ARM Linux。未启动。'; exit 0 ;;
esac

if [ ! -x "$ROOT/nickelgram" ]; then
 printf '%s\n' '缺少 NickelGram 可执行程序；请重新安装。'
 exit 0
fi

"$ROOT/nickelgram" "$@"
exit 0
