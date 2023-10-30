# dcrd v1.5.2

This is a patch release of dcrd to address a potential denial of service vector.

## Changelog

This patch release consists of 5 commits from 2 contributors which total to 4 files changed, 114 additional lines of code, and 20 deleted lines of code.

All commits since the last release may be viewed on GitHub [here](https://github.com/sebitt27/dcrd/compare/release-v1.5.1...release-v1.5.2).

### Protocol and network:

- blockmanager: handle notfound messages from peers ([sebitt27/dcrd#2344](https://github.com/sebitt27/dcrd/pull/2344))
- blockmanager: limit the requested maps ([sebitt27/dcrd#2344](https://github.com/sebitt27/dcrd/pull/2344))
- server: increase ban score for notfound messages ([sebitt27/dcrd#2344](https://github.com/sebitt27/dcrd/pull/2344))
- server: return whether addBanScore disconnected the peer ([sebitt27/dcrd#2344](https://github.com/sebitt27/dcrd/pull/2344))

### Misc:

- release: Bump for 1.5.2([sebitt27/dcrd#2345](https://github.com/sebitt27/dcrd/pull/2345))

### Code Contributors (alphabetical order):

- Dave Collins
- David Hill
