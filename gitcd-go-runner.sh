#!/usr/bin/env bash

gitcd-go $*

if [ -f ~/.config/gitcd/change_dir.sh ]; then
  source ~/.config/gitcd/change_dir.sh
  rm ~/.config/gitcd/change_dir.sh
fi
