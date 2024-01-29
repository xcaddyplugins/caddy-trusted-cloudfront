#!/usr/bin/env bash
function make_latest() {
  local ver=`curl -s https://api.github.com/repos/caddyserver/caddy/releases/latest | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
  local nver=${ver:1}
  # local body=`curl -s https://api.github.com/repos/caddyserver/caddy/releases/latest | grep "body" | head -n 1 | awk '{print substr($0, 12)}' | awk '{print substr($0, 0, length($0)-3)}' | sed 's/\\r//g' | sed 's/\\n/\n/g'`
  # echo "$body" > body.txt
  local html_url=`curl -s https://api.github.com/repos/caddyserver/caddy/releases/latest | grep "html_url" | head -n 1 | awk '{print substr($0, 16)}' | awk '{print substr($0, 0, length($0)-2)}'`
  local body="## Caddy [$ver]($html_url) with trusted CloudFront plugin"
  local self_ver=`curl -s https://api.github.com/repos/xcaddyplugins/caddy-trusted-cloudfront/releases/latest | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
  if [ "$ver" == "$self_ver" ]; then
    echo "::set-output name=updated::0"
    return
  fi
  echo "::set-output name=updated::1"
  echo "::set-output name=tag::$ver"
  echo "::set-output name=ntag::$nver"
  echo "::set-output name=body::$body"
}

make_latest