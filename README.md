[![Build Status](https://github.com/xcaddyplugins/caddy-trusted-cloudfront/workflows/update/badge.svg)](https://github.com/xcaddyplugins/caddy-trusted-cloudfront)
[![Licenses](https://img.shields.io/github/license/xcaddyplugins/caddy-trusted-cloudfront)](LICENSE)
[![donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.buymeacoffee.com/illi)

# trusted_proxies cloudfront module for `Caddy`

The module auto trusted_proxies `AWS CloudFront EDGE servers` from https://d7uri8nf7uskq.cloudfront.net/tools/list-cloudfront-ips

## Requirements

- [Go installed](https://golang.org/doc/install)
- [xcaddy](https://github.com/caddyserver/xcaddy)

## Install

The simplest, cross-platform way to get started is to download Caddy from [GitHub Releases](https://github.com/xcaddyplugins/caddy-trusted-cloudfront/releases) and place the executable file in your PATH.

## Build from source

Requirements:

- [Go installed](https://golang.org/doc/install)
- [xcaddy](https://github.com/caddyserver/xcaddy)

Build:

```bash
$ xcaddy build --with github.com/xcaddyplugins/caddy-trusted-cloudfront
```

## `Caddyfile` Syntax

```caddyfile
trusted_proxies cloudfront {
	interval <duration>
}
```

- `interval` How often to fetch the latest IP list. format is [caddy.Duration](https://caddyserver.com/docs/conventions#durations). For example `12h` represents **12 hours**, and "1d" represents **one day**. default value `1d`.

## `Caddyfile` Example

```caddyfile
trusted_proxies cloudfront {
	interval 1d
}
```

### `Caddyfile` Use Default Settings Example

```Caddyfile
trusted_proxies cloudfront
```

## `Caddyfile` Global Trusted Example

Insert the following configuration of `Caddyfile` to apply it globally.

```Caddyfile
{
	servers {
		trusted_proxies cloudfront
	}
}
```
