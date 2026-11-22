#!/bin/bash
set -eu

widget_dir=$(cd -- "$(dirname -- "$0")" && pwd -P)
widget_path="$widget_dir/free-disk-space-widget"
if [[ ! -x "$widget_path" ]]; then
  echo "Build the widget before creating the LaunchAgent: $widget_path" >&2
  exit 1
fi

label=free-disk-space-widget.autostart
agent_dir="$HOME/Library/LaunchAgents"
agent_path="$agent_dir/$label.plist"
mkdir -p "$agent_dir"
tmp_plist=$(mktemp "$agent_dir/.$label.XXXXXX")
trap 'rm -f "$tmp_plist"' EXIT

/usr/bin/plutil -create xml1 "$tmp_plist"
/usr/bin/plutil -insert Label -string "$label" "$tmp_plist"
/usr/bin/plutil -insert RunAtLoad -bool true "$tmp_plist"
/usr/bin/plutil -insert ProgramArguments -array "$tmp_plist"
index=0
for arg in "$widget_path" "$@"; do
  /usr/bin/plutil -insert "ProgramArguments.$index" -string "$arg" "$tmp_plist"
  index=$((index + 1))
done

chmod 644 "$tmp_plist"
mv -f "$tmp_plist" "$agent_path"
echo "Created $agent_path"
