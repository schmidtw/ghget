# ghget

A simple zero-dependency tool for downloading GitHub release assets — most useful when you're SSH'd into a remote VM and need to grab a tool from GitHub fast.

## install

```sh
curl -fsSL https://schmidtw.github.io/ghget/install.sh | sh
```

Or from source:

```sh
go install github.com/schmidtw/ghget/cmd/ghget@latest
```

Prebuilt binaries are published for linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64.

## usage

```sh
ghget https://github.com/OWNER/REPO/releases/download/TAG/ASSET
```

The easiest way to get that URL: on any GitHub release page, right-click an asset in the Assets list and choose "Copy link address" — that's exactly the argument `ghget` expects.

For private repos, set `GITHUB_TOKEN` in the environment.
