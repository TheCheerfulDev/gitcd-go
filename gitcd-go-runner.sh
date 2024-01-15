#!/usr/bin/env bash

gitcd-go $*

if [ -f /Users/mark/.config/gitcd/change_dir.sh ]; then
  source /Users/mark/.config/gitcd/change_dir.sh
  rm /Users/mark/.config/gitcd/change_dir.sh
fi
