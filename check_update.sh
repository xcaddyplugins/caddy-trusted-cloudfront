#!/usr/bin/env bash
function make_latest() {
  local ver=`curl -s https://api.github.com/repos/caddyserver/caddy/releases/latest | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
  local self_ver=`curl -s https://api.github.com/repos/xcaddyplugins/caddy-trusted-cloudfront/releases/latest | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
  if [ "$ver" == "$self_ver" ]; then
    echo "::set-output name=updated::0"
    return
  fi
  echo "::set-output name=updated::1"
  echo "::set-output name=tag::$ver"
}

make_latest