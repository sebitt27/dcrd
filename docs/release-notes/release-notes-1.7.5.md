# dcrd v1.7.5

This is a patch release of dcrd that updates the utxo cache to improve its
robustness, optimize it, and correct some hard to hit corner cases that involve
a mix of manual block invalidation, conditional flushing, and successive unclean
shutdowns.

## Changelog

This patch release consists of 19 commits from 1 contributor which total to 13
files changed, 1118 additional lines of code, and 484 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/sebitt27/dcrd/compare/a2c3c656...release-v1.7.5).

### Developer-related package and module changes:

- blockchain: Misc consistency cleanup pass ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Pre-allocate in-flight utxoview tx map ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Remove unused utxo cache add entry err ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Fix rare unclean utxo cache recovery ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Don't fetch trsy{base,spend} inputs ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Don't add treasurybase utxos ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Separate utxo cache vs view state ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Improve utxo cache spend robustness ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Split regular/stake view tx connect ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Bypass utxo cache for zero conf spends ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- main: Use backported blockchain updates ([sebitt27/dcrd#3007](https://github.com/sebitt27/dcrd/pull/3007))

### Testing and Quality Assurance:

- blockchain: Address some linter complaints ([sebitt27/dcrd#3005](https://github.com/sebitt27/dcrd/pull/3005))
- blockchain: Allow tests to override cache flushing ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Improve utxo cache initialize tests ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Consolidate utxo cache test entries ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Rework utxo cache spend entry tests ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Rework utxo cache commit tests ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))
- blockchain: Rework utxo cache add entry tests ([sebitt27/dcrd#3006](https://github.com/sebitt27/dcrd/pull/3006))

### Misc:

- release: Bump for 1.7.5 ([sebitt27/dcrd#3008](https://github.com/sebitt27/dcrd/pull/3008))

### Code Contributors (alphabetical order):

- Dave Collins
