# trusted_proxies cloudfront module for `Caddy`

The module auto trusted_proxies `AWS CloudFront EDGE servers` from https://d7uri8nf7uskq.cloudfront.net/tools/list-cloudfront-ips

## Requirements

- [Go installed](https://golang.org/doc/install)
- [xcaddy](https://github.com/caddyserver/xcaddy)

## Build

```bash
$ xcaddy build --with github.com/xcaddyplugins/caddy-trusted-cloudfront
```

## `Caddyfile` Syntax

```caddyfile
trusted_proxies cloudfront {
    interval <duration>
}
```

- `interval` How often to fetch the latest IP list. format is [Golang Duration](https://pkg.go.dev/time#ParseDuration). For example `12h` represents **12 hours**, and "1d" represents **one day**. default value `1d`.

## `Caddyfile` Example

```caddyfile
trusted_proxies cloudfront {
    interval 1d
}
```
