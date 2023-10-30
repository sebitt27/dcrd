# dcrd v1.7.7

This is a patch release of dcrd that includes the following changes:

- Use the latest network protocol version
- Reduce bandwidth usage in certain scenarios by avoiding requests for inventory that is already known
- Mitigate excessive CPU usage in some rare scenarios specific to the test network
- Improve best address candidate selection efficiency

## Changelog

This patch release consists of 19 commits from 3 contributors which total to 92
files changed, 1357 additional lines of code, and 1191 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/sebitt27/dcrd/compare/release-v1.7.5...release-v1.7.7).

### Protocol and network:

- peer: Use latest pver by default ([sebitt27/dcrd#3083](https://github.com/sebitt27/dcrd/pull/3083))
- peer: Correct known inventory check ([sebitt27/dcrd#3083](https://github.com/sebitt27/dcrd/pull/3083))

### Documentation:

- peer: Go 1.19 doc comment formatting ([sebitt27/dcrd#3083](https://github.com/sebitt27/dcrd/pull/3083))
- addrmgr: Go 1.19 doc comment formatting ([sebitt27/dcrd#3084](https://github.com/sebitt27/dcrd/pull/3084))
- multi: Go 1.19 doc comment formatting ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))
- docs: Update README.md to required Go 1.19/1.20 ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))

### Developer-related package and module changes:

- peer: Support module graph prune and lazy load ([sebitt27/dcrd#3083](https://github.com/sebitt27/dcrd/pull/3083))
- main: Use backported peer updates ([sebitt27/dcrd#3083](https://github.com/sebitt27/dcrd/pull/3083))
- addmrgr: Use TempDir to create temp test dirs ([sebitt27/dcrd#3084](https://github.com/sebitt27/dcrd/pull/3084))
- addrmgr: Support module graph prune and lazy load ([sebitt27/dcrd#3084](https://github.com/sebitt27/dcrd/pull/3084))
- addrmgr: Break after selecting random address ([sebitt27/dcrd#3084](https://github.com/sebitt27/dcrd/pull/3084))
- addrmgr: Set min value and optimize address chance ([sebitt27/dcrd#3084](https://github.com/sebitt27/dcrd/pull/3084))
- main: Use backported addrmgr updates ([sebitt27/dcrd#3084](https://github.com/sebitt27/dcrd/pull/3084))
- main: Update to use latest sys module ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))

### Testing and Quality Assurance:

- build: Enable run_tests.sh to work with go.work ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))
- build: Update to latest action versions ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))
- build: Update golangci-lint to v1.51.1 ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))
- build: Test against Go 1.20 ([sebitt27/dcrd#3087](https://github.com/sebitt27/dcrd/pull/3087))

### Misc:

- release: Bump for 1.7.7 ([sebitt27/dcrd#3085](https://github.com/sebitt27/dcrd/pull/3085))

### Code Contributors (alphabetical order):

- Dave Collins
- Eng Zer Jun
- Jonathan Chappelow
