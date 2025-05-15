# GOBIN Updater

Currently a work in progress.

## What is this for?

This is designed to fix a specific problem. I install and keep my locally installed Go language up to date by using [webinstall.dev](https://webinstall.dev/). I don't know if its the best, but it lets my install easily and change the current version easily.

When installing Go tools/programs, it will install them directly to the GOBIN path. Well, when I update to a new Go version, the tool/programs are left behind and I can no longer use them without installing them again.

The goal is to detect the previous versions of go installed, list out the tools, and both install them for the current version, and remove the previous old installs.
