<p align="center"><img width="300px" src="https://bradley.codes/static/img/tcolors/logo.png" alt="tcolors"/></p>

**<p align="center">Commandline color picker and palette builder</p>**

<p align="center"><img src="https://bradley.codes/static/img/tcolors/screencap.png" alt="tcolors"/></p>

#

## Installing

Go get with:

```bash
go get github.com/bcicen/tcolors@v0.3.0
```

Or download the [latest release](https://github.com/bcicen/tcolors/releases) for your platform:

#### Linux / OSX

```bash
curl -Lo tcolors https://github.com/bcicen/tcolors/releases/download/v0.2/tcolors-0.2-$(uname -s)-amd64
chmod +x tcolors
sudo mv tcolors /usr/local/bin/
```
#### AUR

`tcolors` is also available for Arch in the [AUR](https://aur.archlinux.org/packages/tcolors)

#### Docker

```bash
docker run --rm -ti --name=tcolors \
  quay.io/vektorlab/tcolors:latest
```

## Usage

Simply run `tcolors` to view and modify the default palette. Changes are automatically saved and will persist across sessions.

### Keybindings

Key | Description
--- | ---
`↑, k` | navigate up
`↓, j` | navigate down
`←, h` | decrease selected value
`→, l` | increase selected value
`<shift> + ←/→/h/l` | more quickly increase/decrease selected value
`a, <ins>` | add a new palette color
`x, <del>` | remove the selected palette color
`q, <esc>` | exit tcolors
`?` | show help menu

### Palette files

To create a new palette or use a specific palette, use the `-f` option:

```bash
tcolors -f logo-palette.toml
```

Palette colors are stored in a human-readable TOML format and all changes are saved on exit.

### Output

In addition to a stored TOML palette file, `tcolors` provides several output options for parsing and using defined colors

#### All

Default output option providing a formatted table of colors

```bash
# tcolors -p
+----+--------+-------------+-------------+------------------------------------+
| #  |  HEX   |     HSV     |     RGB     |                TERM                |
+----+--------+-------------+-------------+------------------------------------+
| bg | 141414 | 000 000 008 | 020 020 020 | \033[38;2;020;020;020m$@\033[0;00m |
|  0 | FF7733 | 020 080 100 | 255 119 051 | \033[38;2;255;119;051m$@\033[0;00m |
|  1 | FFDD33 | 050 080 100 | 255 221 051 | \033[38;2;255;221;051m$@\033[0;00m |
|  2 | C8FF59 | 080 065 100 | 200 255 089 | \033[38;2;200;255;089m$@\033[0;00m |
|  3 | 55FF33 | 110 080 100 | 085 255 051 | \033[38;2;085;255;051m$@\033[0;00m |
|  4 | 33FF77 | 140 080 100 | 051 255 119 | \033[38;2;051;255;119m$@\033[0;00m |
|  5 | 33FFDD | 170 080 100 | 051 255 221 | \033[38;2;051;255;221m$@\033[0;00m |
|  6 | 33BBFF | 200 080 100 | 051 187 255 | \033[38;2;051;187;255m$@\033[0;00m |
+----+--------+-------------+-------------+------------------------------------+
```

#### Hex, RGB, HSV

Each of these output options provide all colors in a single comma-delimited line; e.g:

```bash
# tcolors -p -o hex
141414, FF7733, FFDD33, C8FF59, 55FF33, 33FF77, 33FFDD, 33BBFF
```

#### Term

The `term` output option provides a series of named functions for easy importing and terminal use
```bash
# tcolors -p -o term
_colorbg() { echo -ne "\033[38;2;020;020;020m$@\033[0;00m"; }
_color0() { echo -ne "\033[38;2;255;119;051m$@\033[0;00m"; }
_color1() { echo -ne "\033[38;2;255;221;051m$@\033[0;00m"; }
_color2() { echo -ne "\033[38;2;200;255;089m$@\033[0;00m"; }
_color3() { echo -ne "\033[38;2;085;255;051m$@\033[0;00m"; }
_color4() { echo -ne "\033[38;2;051;255;119m$@\033[0;00m"; }
_color5() { echo -ne "\033[38;2;051;255;221m$@\033[0;00m"; }
_color6() { echo -ne "\033[38;2;051;187;255m$@\033[0;00m"; }
```

Sourcing:
```bash
source <(tcolors -p -o term)
echo "my $(_color2 what) a $(_color4 bright) $(_color6 day)"
```

### Options

Option | Description
--- | ---
-f | specify palette file to load/save changes to
-p | output current palette contents
-o | color format to output (hex, rgb, hsv, term, all) (default "all")
-v | print version info
