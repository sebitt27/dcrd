# dcrd v1.1.0

This release of dcrd primarily introduces a new consensus vote agenda which
allows the stakeholders to decide whether or not to activate the features needed
for providing full support for Lightning Network.  For those unfamiliar with the
voting process in Decred, this means that all code in order to support these
features is already included in this release, however its enforcement will
remain dormant until the stakeholders vote to activate it.

The following Decred Change Proposals (DCPs) describe the proposed changes in detail:
- [DCP0002](https://github.com/decred/dcps/blob/master/dcp-0002/dcp-0002.mediawiki)
- [DCP0003](https://github.com/decred/dcps/blob/master/dcp-0003/dcp-0003.mediawiki)

It is important for everyone to upgrade their software to this latest release
even if you don't intend to vote in favor of the agenda.

## Notable Changes

### Lightning Network Features Vote

In order to fully support many of the benefits that the Lightning Network will
bring, there are some primitives that involve changes to the current consensus
that need to be enabled.  A new vote with the id `lnfeatures` is now available
as of this release.  After upgrading, stakeholders may set their preferences
through their wallet or stake pool's website.

### Transaction Finality Policy

The standard policy for transaction relay has been changed to use the median
time of the past several blocks instead of the current network adjusted time
when examining lock times to determine if a transaction is final.  This provides
a more deterministic check across all peers and prevents the possibility of
miners attempting to game the timestamps in order to include more transactions.

Consensus enforcement of this change relies on the result of the aforementioned
`lnfeatures` vote.

### Relative Time Locks Policy

The standard policy for transaction relay has been modified to enforce relative
lock times for version 2 transactions via their sequence numbers and a new
`OP_CHECKSEQUENCEVERIFY` opcode.

Consensus enforcement of this change relies on the result of the aforementioned
`lnfeatures` vote.

### OP_SHA256 Opcode

In order to better support cross-chain interoperability, a new opcode to compute
the SHA-256 hash is being proposed.  Since this opcode is implemented as a hard
fork, it will not be available for use in scripts unless the aforementioned
`lnfeatures` vote passes.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/sebitt27/dcrd/compare/v1.0.7...v1.1.0).

### Protocol and network:
- chaincfg: update checkpoints for 1.1.0 release [sebitt27/dcrd#850](https://github.com/sebitt27/dcrd/pull/850)
- chaincfg: Introduce agenda for v5 lnfeatures vote [sebitt27/dcrd#848](https://github.com/sebitt27/dcrd/pull/848)
- txscript: Introduce OP_SHA256 [sebitt27/dcrd#851](https://github.com/sebitt27/dcrd/pull/851)
- wire: Decrease num allocs when decoding headers [sebitt27/dcrd#861](https://github.com/sebitt27/dcrd/pull/861)
- blockchain: Implement enforced relative seq locks [sebitt27/dcrd#864](https://github.com/sebitt27/dcrd/pull/864)
- txscript: Implement CheckSequenceVerify [sebitt27/dcrd#864](https://github.com/sebitt27/dcrd/pull/864)
- multi: Enable vote for DCP0002 and DCP0003 [sebitt27/dcrd#855](https://github.com/sebitt27/dcrd/pull/855)

### Transaction relay (memory pool):
- mempool: Use median time for tx finality checks [sebitt27/dcrd#860](https://github.com/sebitt27/dcrd/pull/860)
- mempool: Enforce relative sequence locks [sebitt27/dcrd#864](https://github.com/sebitt27/dcrd/pull/864)
- policy/mempool: Enforce CheckSequenceVerify opcode [sebitt27/dcrd#864](https://github.com/sebitt27/dcrd/pull/864)

### RPC:
- rpcserver: check whether ticketUtx was found [sebitt27/dcrd#824](https://github.com/sebitt27/dcrd/pull/824)
- rpcserver: return rule error on rejected raw tx [sebitt27/dcrd#808](https://github.com/sebitt27/dcrd/pull/808)

### dcrd command-line flags:
- config: Extend --profile cmd line option to allow interface to be specified [sebitt27/dcrd#838](https://github.com/sebitt27/dcrd/pull/838)

### Documentation
- docs: rpcapi format update [sebitt27/dcrd#807](https://github.com/sebitt27/dcrd/pull/807)
- config: export sampleconfig for use by dcrinstall [sebitt27/dcrd#834](https://github.com/sebitt27/dcrd/pull/834)
- sampleconfig: Add package README and doc.go [sebitt27/dcrd#835](https://github.com/sebitt27/dcrd/pull/835)
- docs: create entry for getstakeversions in rpcapi [sebitt27/dcrd#819](https://github.com/sebitt27/dcrd/pull/819)
- docs: crosscheck and update all rpc doc entries [sebitt27/dcrd#847](https://github.com/sebitt27/dcrd/pull/847)
- docs: update git commit messages section heading [sebitt27/dcrd#863](https://github.com/sebitt27/dcrd/pull/863)

### Developer-related package changes:
- Fix and regenerate precomputed secp256k1 curve [sebitt27/dcrd#823](https://github.com/sebitt27/dcrd/pull/823)
- dcrec: use hardcoded datasets in tests [sebitt27/dcrd#822](https://github.com/sebitt27/dcrd/pull/822)
- Use dchest/blake256  [sebitt27/dcrd#827](https://github.com/sebitt27/dcrd/pull/827)
- glide: use jessevdk/go-flags for consistency [sebitt27/dcrd#833](https://github.com/sebitt27/dcrd/pull/833)
- multi: Error descriptions are in lower case [sebitt27/dcrd#842](https://github.com/sebitt27/dcrd/pull/842)
- txscript: Rename OP_SHA256 to OP_BLAKE256 [sebitt27/dcrd#840](https://github.com/sebitt27/dcrd/pull/840)
- multi: Abstract standard verification flags [sebitt27/dcrd#852](https://github.com/sebitt27/dcrd/pull/852)
- chain: Remove memory block node pruning [sebitt27/dcrd#858](https://github.com/sebitt27/dcrd/pull/858)
- txscript: Add API to parse atomic swap contracts [sebitt27/dcrd#862](https://github.com/sebitt27/dcrd/pull/862)

### Testing and Quality Assurance:
- Test against go 1.9 [sebitt27/dcrd#836](https://github.com/sebitt27/dcrd/pull/836)
- dcrec: remove testify dependency [sebitt27/dcrd#829](https://github.com/sebitt27/dcrd/pull/829)
- mining_test: add edge conditions from btcd [sebitt27/dcrd#831](https://github.com/sebitt27/dcrd/pull/831)
- stake: Modify ticket tests to use chaincfg params [sebitt27/dcrd#844](https://github.com/sebitt27/dcrd/pull/844)
- blockchain: Modify tests to use chaincfg params [sebitt27/dcrd#845](https://github.com/sebitt27/dcrd/pull/845)
- blockchain: Cleanup various tests [sebitt27/dcrd#843](https://github.com/sebitt27/dcrd/pull/843)
- Ensure run_tests.sh local fails correctly when gometalinter errors [sebitt27/dcrd#846](https://github.com/sebitt27/dcrd/pull/846)
- peer: fix logic race in peer connection test [sebitt27/dcrd#865](https://github.com/sebitt27/dcrd/pull/865)

### Misc:
- glide: sync deps [sebitt27/dcrd#837](https://github.com/sebitt27/dcrd/pull/837)
- Update decred deps for v1.1.0 [sebitt27/dcrd#868](https://github.com/sebitt27/dcrd/pull/868)
- Bump for v1.1.0 [sebitt27/dcrd#867](https://github.com/sebitt27/dcrd/pull/867)

### Code Contributors (alphabetical order):

- Alex Yocom-Piatt
- Dave Collins
- David Hill
- Donald Adu-Poku
- Jason Zavaglia
- Jean-Christophe Mincke
- Jolan Luff
- Josh Rickmar
