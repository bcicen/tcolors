# tcolors

Commandline color picker and palette builder

## Installing
Download the latest release for your platform:

```bash
curl -Lo tcolors https://github.com/bcicen/tcolors/releases/download/v0.1/tcolors-0.1-$(uname -s)-amd64
sudo mv tcolors /usr/local/bin/
sudo chmod +x /usr/local/bin/tcolors
```

## Usage

Simply run `tcolors` to view and modify the default template. Changes are automatically saved and will persist across sessions.

### Palette files

To create a new palette or switch between multiple palettes, use the `-f` option:

```bash
tcolors -f logo-palette.toml
```

Palette colors are stored in a human-readable TOML format and any changes are saved on exit.

### Options

Option | Description
--- | ---
-f | specify palette file to load/save changes to
-p | output current palette contents
-o | color format to output (hex, rgb, hsv, all) (default "all")
