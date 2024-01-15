#!/usr/bin/env bash

/Users/mark/dev/thecheerfuldev/sources/gitcd-go/gitcd-go $*

if [ -f /Users/mark/.config/gitcd/change_dir.sh ]; then
  source /Users/mark/.config/gitcd/change_dir.sh
  rm /Users/mark/.config/gitcd/change_dir.sh
fi
