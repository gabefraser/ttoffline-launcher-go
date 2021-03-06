# Toontown Offline Launcher

[![release](https://github.com/gabefraser/ttoffline-launcher-go/actions/workflows/release.yml/badge.svg)](https://github.com/gabefraser/ttoffline-launcher-go/actions/workflows/release.yml)

This is an unofficial launcher I built as the current one struggles from downloading issues.

Feel free to take a look at the source code.

## Usage

Download the latest binary for your operating system of choice [here](https://github.com/gabefraser/ttoffline-launcher-go/releases).

*OPTIONAL* You can add the `--dedicated` flag to the executable to start the server straight from the launcher.

Windows
```
ttoffline-launcher-windows-amd64.exe [--dedicated]
```

Linux
```
chmod u+x ttoffline-launcher-linux-amd64 && ./ttoffline-launcher-linux-amd64 [--dedicated]
```

Mac
```
chmod u+x ttoffline-launcher-darwin-amd64 && ./ttoffline-launcher-darwin-amd64 [--dedicated]
```

## Credits

`go-humanize` - github.com/dustin/go-humanize

`archiver` - github.com/mholt/archiver/v3
