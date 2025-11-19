#!/usr/bin/env bash

docker ps --format "table {{.ID}}\t{{.Image}}\t{{.Names}}" \
  | awk '$3~/^reaper_/ { print $1 }'
