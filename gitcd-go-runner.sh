#!/usr/bin/env bash

/Users/mark/dev/thecheerfuldev/sources/gitcd-go/gitcd $*

if [ -f /tmp/gitcd-action.sh ]; then
  source /tmp/gitcd-action.sh
  rm /tmp/gitcd-action.sh
fi
