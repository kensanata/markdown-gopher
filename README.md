# Markdown Gopher

This is a Gopher server that serves the current directory and its
subdirectories. Markdown files are served as-is, but if the selector
doesn't use the ".md" extension, the Markdown is converted to plain
text.

When the selector is empty or points to a directory, the "index.md"
file is parsed for list items starting with an asterisk in order to
build a menu. The assumption is that these link to Markdown files
without the ".md" extension.

This convention is used by
[Oddmu](https://src.alexschroeder.ch/oddmu.git), for example.

## Installation

To install using systemd, for the current user:

Make necessary changes to `markdown-gopher.service`.

These keys need changing, in particular:

```
ExecStart=/home/alex/bin/markdown-gopher
WorkingDirectory=/home/alex/alexschroeder.ch/wiki
Environment="GOPHER_HOST=alexschroeder.ch"
```

The `GOPHER_PORT` environment variable can be used to change the port,
too.

Enable the unit:

```
sudo systemctl enable --now ./markdown-gopher.service
```

## License

The code of the server is released under the AGPL 3.

The code of the renderer.go started at one point with [the renderer in
gmnhg by Timur Demin](https://github.com/tdemin/gmnhg/tree/v0.4.2/internal/renderer),
released under the GPL 3.

Note section 13. Use with the GNU Affero General Public License.

>   Notwithstanding any other provision of this License, you have
> permission to link or combine any covered work with a work licensed
> under version 3 of the GNU Affero General Public License into a
> single combined work, and to convey the resulting work. The terms of
> this License will continue to apply to the part which is the covered
> work, but the special requirements of the GNU Affero General Public
> License, section 13, concerning interaction through a network will
> apply to the combination as such.

The word-wrapping code I added was heavily inspired by [the golinewrap
library by Karrick McDermott](https://godoc.org/github.com/karrick/golinewrap),
released under the MIT License.
