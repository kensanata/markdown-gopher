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

The code of the server is released under the AGPL 3.

The code of the renderer is released under the GPL 3. It is based on
[the renderer in gmnhg by Timur Demin](https://github.com/tdemin/gmnhg/tree/v0.4.2/internal/renderer).

