#!/bin/bash
#
# scripts/notify.sh: pre-push 背景任務共用的 macOS 通知 (格式與 `inf deploy` 相同)。
#
# 用法: scripts/notify.sh <repo-dir> <message>
#   標題: <repo 名> <路徑 (~ 縮寫)>
#   內文: <message>, 例如 `Installed locally`
# terminal-notifier 優先 (全部 inf 通知共用 group `inf`, 新通知取代舊的), 否則 osascript; 非 macOS 略過。

set -uo pipefail

dir="${1:?usage: notify.sh <repo-dir> <message>}"
message="${2:?usage: notify.sh <repo-dir> <message>}"
[[ "$(uname -s)" == "Darwin" ]] || exit 0

tilde="~" # 變數取代: /bin/bash 3.2 會把 \~ 的反斜線一起印出來
title="$(basename "$dir") ${dir/#"$HOME"/$tilde}"
icon="$(cd "$(dirname "$0")" && pwd)/notify-icon.icns"
app="$HOME/.config/inf/InfNotifier.app"
notifier="$app/Contents/MacOS/terminal-notifier"

if [[ ! -x "$notifier" ]] && [[ -f "$icon" ]] && command -v terminal-notifier >/dev/null 2>&1; then
	real_bin="$(readlink -f "$(command -v terminal-notifier)" 2>/dev/null || true)"
	src_app="$(dirname "$real_bin")/../terminal-notifier.app"
	if [[ -d "$src_app" ]]; then
		mkdir -p "$HOME/.config/inf"
		rm -rf "$app"
		cp -R "$src_app" "$app"
		cp "$icon" "$app/Contents/Resources/Inf.icns"
		rm -f "$app/Contents/Resources/Terminal.icns"
		plutil -replace CFBundleIconFile -string "Inf" "$app/Contents/Info.plist" >/dev/null 2>&1 || true
		plutil -replace CFBundleIdentifier -string "dev.shuks.inf.notifier" "$app/Contents/Info.plist" >/dev/null 2>&1 || true
		plutil -replace CFBundleName -string "inf" "$app/Contents/Info.plist" >/dev/null 2>&1 || true
		codesign --force --deep --sign - "$app" >/dev/null 2>&1 || true
	fi
fi

if [[ -x "$notifier" ]]; then
	"$notifier" -title "$title" -message "$message" -group "inf" >/dev/null 2>&1
elif command -v terminal-notifier >/dev/null 2>&1; then
	terminal-notifier -title "$title" -message "$message" -group "inf" >/dev/null 2>&1
else
	osascript -e "display notification \"${message//\"/\\\"}\" with title \"${title//\"/\\\"}\"" >/dev/null 2>&1
fi
exit 0
