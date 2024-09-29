# ansisvg

Convert [ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code) to [SVG](https://en.wikipedia.org/wiki/Scalable_Vector_Graphics).

Pipe output from a program thru `ansisvg` and it will output a SVG file on stdout.

Can be used to produce nice looking example output for presentations, markdown files etc. Note that it
does not support programs that do cursor movements like ncurses programs etc.

```sh
./colortest | ansisvg > colortest.svg
 ```
Produces [colortest.svg](cli/testdata/colortest.svg):

![ansisvg output for colortest](cli/testdata/colortest.svg)

```
$ ansisvg -h
ansisvg - Convert ANSI to SVG
Usage: ansisvg [FLAGS]

Example usage:
  program | ansisvg > file.svg

--charboxsize WxH   Character box size (use pixel units instead of font units)
--colorscheme NAME  Color scheme
--fontfile PATH     Font file to use and embed
--fontname NAME     Font name
--fontref URL       External font URL to use
--fontsize NUMBER   Font size
--grid              Grid mode (sets position for each character)
--help, -h          Show help
--listcolorschemes  List color schemes
--marginsize WxH    Margin size (in either pixel or font units)
--transparent       Transparent background
--version, -v       Show version
--width, -w NUMBER  Terminal width (auto if not set)
```

Color themes are the ones from https://github.com/mbadolato/iTerm2-Color-Schemes

## Install

Pre-built binaries for Linux, macOS and Windows can be downloaded from [releases](https://github.com/wader/ansisvg/releases).

### macOS

For macOS you might have to allow to run the binary in security preferences. Alternatively run the below command:

```sh
xattr -d com.apple.quarantine ansisvg && spctl --add ansisvg
```

## Build

To build you will need at least go 1.18 or later.

Install latest master and copy it to `/usr/local/bin`:
```sh
go install github.com/wader/ansisvg@master
cp $(go env GOPATH)/bin/ansisvg /usr/local/bin
```

Build from cloned repo:
```
go build -o ansisvg .
```

## Fonts

`ansisvg` can either use system-installed fonts (`-fontname`), link to a webfont on a HTTP server (`-fontref`) or embed a webfont from the local filesystem (`-fontfile`).

### Compatibility issues

* Embedded and/or linked fonts might not be supported by some SVG viewers. At time of writing this is [not supported by Inkscape](https://gitlab.com/inkscape/inbox/-/issues/301).

* For SVGs that are intended to be included in websites via `<img>`, the only way to make a custom font work is [embedding it in the SVG](https://vecta.io/blog/how-to-use-fonts-in-svg).

### Variations of custom fonts (regular/bold/italic)

* System wide fonts (`-fontname`) get correctly rendered with variations, but when using external fonts with `-fontref` or `-fontfile` the SVG viewer knows only the regular variant and will try to render italic/bold text 'extrapolated' from it which may look different than the actual font variation. To use the actual bold/italic font variants, different woff2 files have to be used for the respective text styles which needs additional CSS code (currently not supported by `ansisvg`).

* Bold style 'extrapolated' from the regular font may even break monospace alignment. Use `-grid` option to mitigate that.

## Font-relative vs. pixel coordinates

By default, `ansisvg` uses font-relative `ch`/`em` coordinates. This should make SVG dimensions and line/character spacing consistent with font family/size. When SVG dimensions and/or text coordinates are off, it is possible to force explicit pixel units for coordinates by specifying `-charboxsize` in X/Y pixel units, e.g. `8x16`.

* Inkscape currently [cannot deal with SVG size expressed in font-relative units](https://gitlab.com/inkscape/inkscape/-/issues/4737), a quick workaround is Ctrl-Shift-R (resize page to content).

* Some SVG processing tools like [asciidoctor](https://docs.asciidoctor.org/pdf-converter/latest/image-paths-and-formats/#image-formats) require the presence of the `viewBox` attribute. Use `-charboxsize` option to enable this attribute (it only works with pixel dimensions).

## Margin size

With `--marginsize` a margin can be defined, so there is a bit of empty space (or "border") around the image. Default is zero margin size, i.e. the terminal characters are touching the edge of the image.
`--marginsize` is interpreted as X/Y in the currently selected units, i.e. `ch`/`em` by default, and `px` if `--charboxsize` is used.

## Consolidated text vs. grid mode

By default, `ansisvg` consolidates text to `<tspan>` chunks, leaving the X positioning of characters to the SVG renderer. This usually works well for monospace fonts. However if not all glyphs involved are monospace (e.g. when exotic characters are used, making the SVG renderer fall back to a different font for those characters) then the alignment will be off; this can be worked around with `-grid` mode which will make `ansisvg` put each character to explicit positions, making the SVG bigger and less readable but ensuring proper positioning/alignment for all characters.

## Tricks

### ANSI to PDF or PNG

```
... | ansisvg | inkscape --pipe --export-type=pdf -o file.pdf
... | ansisvg | inkscape --pipe --export-type=png -o file.png
```


### Use `bat` to produce source code highlighting

```
bat --color=always -p main.go | ansisvg
```

### Use `script` to run with a pty

```
script -q /dev/null <command> | ansisvg
```

### ffmpeg

```
TERM=a AV_LOG_FORCE_COLOR=1 ffmpeg ... 2>&1 | ansisvg
```

### jq
```
jq -C | ansisvg
```

### Make screenshots from a terminal

#### tmux

```
# <prefix>-H: Create a SVG screenshot of the current pane
bind H capture-pane -e \; run "tmux save-buffer - | $HOME/go/bin/ansisvg > $HOME/Pictures/tmux-$(date +%F_%T).svg"; delete-buffer
```

#### kitty

```
# F3: Create a SVG screenshot of the current selection
map f3 combine : copy_ansi_to_clipboard : launch sh -c 'kitty +kitten clipboard -g | $HOME/go/bin/ansisvg > $HOME/Pictures/kitty-$(date +%F_%T).svg'
```

## Development and release build

Run all tests and write new test output:
```
go test ./... -update
```

Manual release build with version can be done with:
```
go build -ldflags "-X main.version=1.2.3" -o ansisvg .
```

Visual inspect test output in browser:
```
for i in cli/testdata/*.svg; do echo "$i<br><img src=\"$i\"/><br>" ; done  > all.html
open all.html
```

Using [ffcat](https://github.com/wader/ffcat):
```
for i in cli/testdata/*.ansi; do echo $i ; cat $i | go run . | ffcat ; done
```

## Thanks

- Patrick Huesmann [@patrislav1](https://github.com/patrislav1) for better ANSI support and lots SVG output improvements.

## Licenses and thanks

Color themes from
https://github.com/mbadolato/iTerm2-Color-Schemes,
license https://github.com/mbadolato/iTerm2-Color-Schemes/blob/master/LICENSE

Uses colortest from https://github.com/pablopunk/colortest and terminal-colors from https://github.com/eikenb/terminal-colors.

 UbuntuMonoNerdFontMono-Regular.woff2 from https://github.com/ryanoasis/nerd-fonts license https://github.com/ryanoasis/nerd-fonts/blob/master/LICENSE

## TODO and ideas
- Underline overlaps a bit, sometimes causing weird blending
- Handle vertical tab and form feed (normalize into spaces?)
- Handle overdrawing
- More CSI, keep track of cursor?
- PNG output (embed nice fonts?)
