# dcrd v1.7.2

This is a patch release of dcrd to resolve a rare and hard to hit case when
optional indexing is enabled.

## Changelog

This patch release consists of 4 commits from 2 contributors which total to 11
files changed, 158 additional lines of code, and 15 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/sebitt27/dcrd/compare/release-v1.7.1...release-v1.7.2) and
[here](https://github.com/sebitt27/dcrd/compare/blockchain/v4.0.0...blockchain/v4.0.1).

### Protocol and network:

- server: Fix syncNotified race ([sebitt27/dcrd#2931](https://github.com/sebitt27/dcrd/pull/2931))

### Developer-related package and module changes:

- indexers: fix subscribers race ([sebitt27/dcrd#2921](https://github.com/sebitt27/dcrd/pull/2921))
- main: Use backported blockchain updates ([sebitt27/dcrd#2935](https://github.com/sebitt27/dcrd/pull/2935))

### Misc:

- release: Bump for 1.7.2 ([sebitt27/dcrd#2936](https://github.com/sebitt27/dcrd/pull/2936))

### Code Contributors (alphabetical order):

- Dave Collins
- Donald Adu-Poku
