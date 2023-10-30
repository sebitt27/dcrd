# dcrd v1.7.0

This is a new major release of dcrd.  Some of the key highlights are:

* Four new consensus vote agendas which allow stakeholders to decide whether or
  not to activate support for the following:
  * Reverting the Treasury maximum expenditure policy
  * Enforcing explicit version upgrades
  * Support for automatic ticket revocations for missed votes
  * Changing the Proof-of-Work and Proof-of-Stake subsidy split from 60%/30% to 10%/80%
* Substantially reduced initial sync time
* Major performance enhancements to unspent transaction output handling
* Faster cryptographic signature validation
* Significant improvements to network synchronization
* Support for a configurable assumed valid block
* Block index memory usage reduction
* Asynchronous indexing
* Version 1 block filters removal
* Various updates to the RPC server:
  * Additional per-connection read limits
  * A more strict cross origin request policy
  * A new alternative client authentication mechanism based on TLS certificates
  * Availability of the scripting language version for transaction outputs
  * Several other notable updates, additions, and removals related to the JSON-RPC API
* New developer modules:
  * Age-Partitioned Bloom Filters
  * Fixed-Precision Unsigned 256-bit Integers
  * Standard Scripts
  * Standard Addresses
* Infrastructure improvements
* Quality assurance changes

For those unfamiliar with the
[voting process](https://docs.decred.org/governance/consensus-rule-voting/overview/)
in Decred, all code needed in order to support each of the aforementioned
consensus changes is already included in this release, however it will remain
dormant until the stakeholders vote to activate it.

For reference, the consensus change work for each of the four changes was
originally proposed and approved for initial implementation via the following
Politeia proposals:
- [Decentralized Treasury Spending](https://proposals-archive.decred.org/proposals/c96290a)
- [Explicit Version Upgrades Consensus Change](https://proposals.decred.org/record/3a98861)
- [Automatic Ticket Revocations Consensus Change](https://proposals.decred.org/record/e2d7b7d)
- [Change PoW/PoS Subsidy Split From 60/30 to 10/80](https://proposals.decred.org/record/427e1d4)

The following Decred Change Proposals (DCPs) describe the proposed changes in
detail and provide full technical specifications:
- [DCP0007](https://github.com/decred/dcps/blob/master/dcp-0007/dcp-0007.mediawiki)
- [DCP0008](https://github.com/decred/dcps/blob/master/dcp-0008/dcp-0008.mediawiki)
- [DCP0009](https://github.com/decred/dcps/blob/master/dcp-0009/dcp-0009.mediawiki)
- [DCP0010](https://github.com/decred/dcps/blob/master/dcp-0010/dcp-0010.mediawiki)

## Upgrade Required

**It is extremely important for everyone to upgrade their software to this
latest release even if you don't intend to vote in favor of the agenda.  This
particularly applies to PoW miners as failure to upgrade will result in lost
rewards after block height 635775.  That is estimated to be around Feb 21st,
2022.**

## Downgrade Warning

The database format in v1.7.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

The database migration typically takes around 40-50 minutes on HDDs and 20-30
minutes on SSDs.

## Notable Changes

### Four New Consensus Change Votes

Four new consensus change votes are now available as of this release.  After
upgrading, stakeholders may set their preferences through their wallet.

#### Revert Treasury Maximum Expenditure Policy Vote

The first new vote available as of this release has the id `reverttreasurypolicy`.

The primary goal of this change is to revert the currently active maximum
expenditure policy of the decentralized Treasury to the one specified in the
[original Politeia proposal](https://proposals-archive.decred.org/proposals/c96290a).

See [DCP0007](https://github.com/decred/dcps/blob/master/dcp-0007/dcp-0007.mediawiki) for
the full technical specification.

#### Explicit Version Upgrades Vote

The second new vote available as of this release has the id `explicitverupgrades`.

The primary goals of this change are to:

* Provide an easy, reliable, and efficient method for software and hardware to
  determine exactly which rules should be applied to transaction and script
  versions
* Further embrace the increased security and other desirable properties that
  hard forks provide over soft forks

See the following for more details:

* [Politeia proposal](https://proposals.decred.org/record/3a98861)
* [DCP0008](https://github.com/decred/dcps/blob/master/dcp-0008/dcp-0008.mediawiki)

#### Automatic Ticket Revocations Vote

The third new vote available as of this release has the id `autorevocations`.

The primary goals of this change are to:

* Improve the Decred stakeholder user experience by removing the requirement for
  stakeholders to manually revoke missed and expired tickets
* Enable the recovery of funds for users who lost their redeem script for the
  legacy VSP system (before the release of vspd, which removed the need for the
  redeem script)

See the following for more details:

* [Politeia proposal](https://proposals.decred.org/record/e2d7b7d)
* [DCP0009](https://github.com/decred/dcps/blob/master/dcp-0009/dcp-0009.mediawiki)

#### Change PoW/PoS Subsidy Split to 10/80 Vote

The fourth new vote available as of this release has the id `changesubsidysplit`.

The proposed modification to the subsidy split is intended to substantially
diminish the ability to attack Decred's markets with mined coins and improve
decentralization of the issuance process.

See the following for more details:

* [Politeia proposal](https://proposals.decred.org/record/427e1d4)
* [DCP0010](https://github.com/decred/dcps/blob/master/dcp-0010/dcp-0010.mediawiki)

### Substantially Reduced Initial Sync Time

The amount of time it takes to complete the initial chain synchronization
process has been substantially reduced.  With default settings, it is around 48%
faster versus the previous release.

### Unspent Transaction Output Overhaul

The way unspent transaction outputs (UTXOs) are handled has been significantly
reworked to provide major performance enhancements to both steady-state
operation as well as the initial chain sync process as follows:

* Each UTXO is now tracked independently on a per-output basis
* The UTXOs now reside in a dedicated database
* All UTXO reads and writes now make use of a cache

#### Unspent Transaction Output Cache

All reads and writes of unspent transaction outputs (utxos) now go through a
cache that sits on top of the utxo set database which drastically reduces the
amount of reading and writing to disk, especially during the initial sync
process when a very large number of blocks are being processed in quick
succession.

This utxo cache provides significant runtime performance benefits at the cost of
some additional memory usage.  The maximum size of the cache can be configured
with the new `--utxocachemaxsize` command-line configuration option.  The
default value is 150 MiB, the minimum value is 25 MiB, and the maximum value is
32768 MiB (32 GiB).

Some key properties of the cache are as follows:

* For reads, the UTXO cache acts as a read-through cache
  * All UTXO reads go through the cache
  * Cache misses load the missing data from the disk and cache it for future lookups
* For writes, the UTXO cache acts as a write-back cache
  * Writes to the cache are acknowledged by the cache immediately, but are only
    periodically flushed to disk
* Allows intermediate steps to effectively be skipped thereby avoiding the need
  to write millions of entries to disk
* On average, recent UTXOs are much more likely to be spent in upcoming blocks
  than older UTXOs, so only the oldest UTXOs are evicted as needed in order to
  maximize the hit ratio of the cache
* The cache is periodically flushed with conditional eviction:
  * When the cache is NOT full, nothing is evicted, but the changes are still
    written to the disk set to allow for a quicker reconciliation in the case of
    an unclean shutdown
  * When the cache is full, 15% of the oldest UTXOs are evicted

### Faster Cryptographic Signature Validation

Some aspects of the underlying crypto code has been updated to further improve
its execution speed and reduce the number of memory allocations resulting in
about a 1% reduction to signature verification time.

The primary benefits are:

* Improved vote times since blocks and transactions propagate more quickly
  throughout the network
* Approximately a 2% reduction to the duration of the initial sync process

### Significant Improvements to Network Synchronization

The method used to obtain blocks from other peers on the network is now guided
entirely by block headers.  This provides a wide variety of benefits, but the
most notable ones for most users are:

* Faster initial synchronization
* Reduced bandwidth usage
* Enhanced protection against attempted DoS attacks
* Percentage-based progress reporting
* Improved steady state logging

### Support for Configurable Assumed Valid Block

This release introduces a new model for deciding when several historical
validation checks may be skipped for blocks that are an ancestor of a known good
block.

Specifically, a new `AssumeValid` parameter is now used to specify the
aforementioned known good block.  The default value of the parameter is updated
with each release to a recent block that is part of the main chain.

The default value of the parameter can be overridden with the `--assumevalid`
command-line option by setting it as follows:

* `--assumevalid=0`: Disable the feature resulting in no skipped validation checks
* `--assumevalid=[blockhash]`:  Set `AssumeValid` to the specified block hash

Specifying a block hash closer to the current best chain tip allows for faster
syncing.  This is useful since the validation requirements increase the longer a
particular release build is out as the default known good block becomes deeper
in the chain.

### Block Index Memory Usage Reduction

The block index that keeps track of block status and connectivity now occupies
around 30MiB less memory and scales better as more blocks are added to the
chain.

### Asynchronous Indexing

The various optional indexes are now created asynchronously versus when
blocks are processed as was previously the case.

This permits blocks to be validated more quickly when the indexes are enabled
since the validation no longer needs to wait for the indexing operations to
complete.

In order to help keep consistent behavior for RPC clients, RPCs that involve
interacting with the indexes will not return results until the associated
indexing operation completes when the indexing tip is close to the current best
chain tip.

One side effect of this change that RPC clients should be aware of is that it is
now possible to receive sync timeout errors on RPCs that involve interacting
with the indexes if the associated indexing tip gets so far behind it would end
up delaying results for too long.  In practice, errors of this type are rare and
should only ever be observed during the initial sync process before the
associated indexes are current.  However, callers should be aware of the
possibility and handle the error accordingly.

The following RPCs are affected:

* `existsaddress`
* `existsaddresses`
* `getrawtransaction`
* `searchrawtransactions`

### Version 1 Block Filters Removal

The previously deprecated version 1 block filters are no longer available on the
peer-to-peer network.  Use
[version 2 block filters](https://github.com/decred/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#version-2-block-filters)
with their associated
[block header commitment](https://github.com/decred/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#block-header-commitments)
and [inclusion proof](https://github.com/decred/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#verifying-commitment-root-inclusion-proofs)
instead.

### RPC Server Changes

The RPC server version as of this release is 7.0.0.

#### Max Request Limits

The RPC server now imposes the following additional per-connection read limits
to help further harden it against potential abuse in non-standard configurations
on poorly-configured networks:

* 0 B / 8 MiB for pre and post auth HTTP connections
* 4 KiB / 16 MiB for pre and post auth WebSocket connections

In practice, these changes will not have any noticeable effect for the vast
majority of nodes since the RPC server is not publicly accessible by default and
also requires authentication.

Nevertheless, it can still be useful for scenarios such as authenticated fuzz
testing and improperly-configured networks that have disabled all other security
measures.

#### More Strict Cross Origin Request (CORS) Policy

The CORS policy for WebSocket clients is now more strict and rejects requests
from other domains.

In practice, CORS requests will be rejected before ever reaching that point due
to the use of a self-signed TLS certificate and the requirement for
authentication to issue any commands.  However, additional protection mechanisms
make it that much more difficult to attack by providing defense in depth.

#### Alternative Client Authentication Method Based on TLS Certificates

A new alternative method for TLS clients to authenticate to the RPC server by
presenting a client certificate in the TLS handshake is now available.

Under this authentication method, the certificate authority for a client
certificate must be added to the RPC server as a trusted root in order for it to
trust the client.  Once activated, clients will no longer be required to provide
HTTP Basic authentication nor use the `authenticate` RPC in the case of
WebSocket clients.

Note that while TLS client authentication has the potential to ultimately allow
more fine grained access controls on a per-client basis, it currently only
supports clients with full administrative privileges.  In other words, it is not
currently compatible with the `--rpclimituser` and `--rpclimitpass` mechanism,
so users depending on the limited user settings should avoid the new
authentication method for now.

The new authentication type can be activated with the `--authtype=clientcert`
configuration option.

By default, the trusted roots are loaded from the `clients.pem` file in dcrd's
application data directory, however, that location can be modified via the
`--clientcafile` option if desired.

#### Updates to Transaction Output Query RPC (`gettxout`)

The `gettxout` RPC has the following modifications:

* An additional `tree` parameter is now required in order to explicitly identify
  the exact transaction output being requested
* The transaction `version` field is no longer available in the primary JSON
  object of the results
* The child `scriptPubKey` JSON object in the results now includes a new
  `version` field that identifies the scripting language version

See the
[gettxout JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#gettxout)
for API details.

#### Removal of Stake Difficulty Notification RPCs (`notifystakedifficulty` and `stakedifficulty`)

The deprecated `notifystakedifficulty` and `stakedifficulty` WebSocket-only RPCs
are no longer available.  This notification is unnecessary since the difficulty
change interval is well defined.  Callers may obtain the difficulty via
`getstakedifficulty` at the appropriate difficulty change intervals instead.

See the
[getstakedifficulty JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getstakedifficulty)
for API details.

#### Removal of Version 1 Filter RPCs (`getcfilter` and `getcfilterheader`)

The deprecated `getcfilter` and `getcfilterheader` RPCs, which were previously
used to obtain version 1 block filters via RPC are no longer available. Use
`getcfilterv2` instead.

See the
[getcfilterv2 JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getcfilterv2)
for API details.

#### New Median Time Field on Block Query RPCs (`getblock` and `getblockheader`)

The verbose results of the `getblock` and `getblockheader` RPCs now include a
`mediantime` field that specifies the median block time associated with the
block.

See the following for API details:

* [getblock JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getblock)
* [getblockheader JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getblockheader)

#### New Scripting Language Version Field on Raw Transaction RPCs (`getrawtransaction`, `decoderawtransaction`, `searchrawtransactions`, and `getblock`)

The verbose results of the `getrawtransaction`, `decoderawtransaction`,
`searchrawtransactions`, and `getblock` RPCs now include a `version` field in
the child `scriptPubKey` JSON object that identifies the scripting language
version.

See the following for API details:

* [getrawtransaction JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getrawtransaction)
* [decoderawtransaction JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#decoderawtransaction)
* [searchrawtransactions JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#searchrawtransactions)
* [getblock JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getblock)

#### New Treasury Add Transaction Filter on Mempool Query RPC (`getrawmempool`)

The transaction type parameter of the `getrawmempool` RPC now accepts `tadd` to
only include treasury add transactions in the results.

See the
[getrawmempool JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#getrawmempool)
for API details.

#### New Manual Block Invalidation and Reconsideration RPCs (`invalidateblock` and `reconsiderblock`)

A new pair of RPCs named `invalidateblock` and `reconsiderblock` are now
available.  These RPCs can be used to manually invalidate a block as if it had
violated consensus rules and reconsider a block for validation and best chain
selection by removing any invalid status from it and its ancestors, respectively.

This capability is provided for development, testing, and debugging.  It can be
particularly useful when developing services that build on top of Decred to more
easily ensure edge conditions associated with invalid blocks and chain
reorganization are being handled properly.

These RPCs do not apply to regular users and can safely be ignored outside of
development.

See the following for API details:

* [invalidateblock JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#invalidateblock)
* [reconsiderblock JSON-RPC API Documentation](https://github.com/sebitt27/dcrd/blob/master/docs/json_rpc_api.mediawiki#reconsiderblock)

### Reject Protocol Message Deprecated (`reject`)

The `reject` peer-to-peer protocol message is now deprecated and is scheduled to
be removed in a future release.

This message is a holdover from the original codebase where it was required, but
it really is not a useful message and has several downsides:

* Nodes on the network must be trustless, which means anything relying on such a
  message is setting itself up for failure because nodes are not obligated to
  send it at all nor be truthful as to the reason
* It can be harmful to privacy as it allows additional node fingerprinting
* It can lead to security issues for implementations that don't handle it with
  proper sanitization practices
* It can easily give software implementations the fully incorrect impression
  that it can be relied on for determining if transactions and blocks are valid
* The only way it is actually used currently is to show a debug log message,
  however, all of that information is already available via the peer and/or wire
  logging anyway
* It carries a non-trivial amount of development overhead to continue to support
  it when nothing actually uses it

### No DNS Seeds Command-Line Option Deprecated (`--nodnsseed`)

The `--nodnsseed` command-line configuration option is now deprecated and will
be removed in a future release.  Use `--noseeders` instead.

DNS seeding has not been used since the previous release.

## Notable New Developer Modules

### Age-Partitioned Bloom Filters

A new `github.com/sebitt27/dcrd/container/apbf` module is now available that
provides Age-Partitioned Bloom Filters (APBFs).

An APBF is a probabilistic lookup device that can quickly determine if it
contains an element.  It permits tracking large amounts of data while using very
little memory at the cost of a controlled rate of false positives.  Unlike
classic Bloom filters, it is able to handle an unbounded amount of data by aging
and discarding old items.

For a concrete example of actual savings achieved in Decred by making use of an
APBF, the memory to track addresses known by 125 peers was reduced from ~200 MiB
to ~5 MiB.

See the
[apbf module documentation](https://pkg.go.dev/github.com/sebitt27/dcrd/container/apbf)
for full details on usage, accuracy under workloads, expected memory usage, and
performance benchmarks.

### Fixed-Precision Unsigned 256-bit Integers

A new `github.com/sebitt27/dcrd/math/uint256` module is now available that provides
highly optimized allocation free fixed precision unsigned 256-bit integer
arithmetic.

The package has a strong focus on performance and correctness and features
arithmetic, boolean comparison, bitwise logic, bitwise shifts, conversion
to/from relevant types, and full formatting support - all served with an
ergonomic API, full test coverage, and benchmarks.

Every operation is faster than the standard library `big.Int` equivalent and the
primary math operations provide reductions of over 90% in the calculation time.
Most other operations are also significantly faster.

See the
[uint256 module documentation](https://pkg.go.dev/github.com/sebitt27/dcrd/math/uint256)
for full details on usage, including a categorized summary, and performance
benchmarks.

### Standard Scripts

A new `github.com/sebitt27/dcrd/txscript/v4/stdscript` package is now available
that provides facilities for identifying and extracting data from transaction
scripts that are considered standard by the default policy of most nodes.

The package is part of the `github.com/sebitt27/dcrd/txscript/v4` module.

See the
[stdscript package documentation](https://pkg.go.dev/github.com/sebitt27/dcrd/txscript/v4/stdscript)
for full details on usage and a list of the recognized standard scripts.

### Standard Addresses

A new `github.com/sebitt27/dcrd/txscript/v4/stdaddr` package is now available that
provides facilities for working with human-readable Decred payment addresses.

The package is part of the `github.com/sebitt27/dcrd/txscript/v4` module.

See the
[stdaddr package documentation](https://pkg.go.dev/github.com/sebitt27/dcrd/txscript/v4/stdaddr)
for full details on usage and a list of the supported addresses.

## Changelog

This release consists of 877 commits from 16 contributors which total to 492
files changed, 77937 additional lines of code, and 30961 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/sebitt27/dcrd/compare/release-v1.6.0...release-v1.7.0).

### Protocol and network:

- chaincfg: Add extra seeders ([sebitt27/dcrd#2532](https://github.com/sebitt27/dcrd/pull/2532))
- server: Stop serving v1 cfilters over p2p ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- blockchain: Decouple processing and download logic ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- blockchain: Improve current detection ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- netsync: Rework inventory announcement handling ([sebitt27/dcrd#2548](https://github.com/sebitt27/dcrd/pull/2548))
- peer: Add inv type summary to debug message ([sebitt27/dcrd#2556](https://github.com/sebitt27/dcrd/pull/2556))
- netsync: Remove unused submit block flags param ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- netsync: Remove submit/processblock orphan flag ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- netsync: Remove orphan block handling ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- netsync: Rework sync model to use hdr annoucements ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- progresslog: Add support for header sync progress ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- netsync: Add header sync progress log ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- multi: Add chain verify progress percentage ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- peer: Remove getheaders response deadline ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- chaincfg: Update seed URL ([sebitt27/dcrd#2564](https://github.com/sebitt27/dcrd/pull/2564))
- upnp: Don't return loopback IPs in getOurIP ([sebitt27/dcrd#2566](https://github.com/sebitt27/dcrd/pull/2566))
- server: Prevent duplicate pending conns ([sebitt27/dcrd#2563](https://github.com/sebitt27/dcrd/pull/2563))
- multi: Use an APBF for recently confirmed txns ([sebitt27/dcrd#2580](https://github.com/sebitt27/dcrd/pull/2580))
- multi: Use an APBF for per peer known addrs ([sebitt27/dcrd#2583](https://github.com/sebitt27/dcrd/pull/2583))
- peer: Stop sending and logging reject messages ([sebitt27/dcrd#2586](https://github.com/sebitt27/dcrd/pull/2586))
- netsync: Stop sending reject messages ([sebitt27/dcrd#2586](https://github.com/sebitt27/dcrd/pull/2586))
- server: Stop sending reject messages ([sebitt27/dcrd#2586](https://github.com/sebitt27/dcrd/pull/2586))
- peer: Remove deprecated onversion reject return ([sebitt27/dcrd#2586](https://github.com/sebitt27/dcrd/pull/2586))
- peer: Remove unneeded PushRejectMsg ([sebitt27/dcrd#2586](https://github.com/sebitt27/dcrd/pull/2586))
- wire: Deprecate reject message ([sebitt27/dcrd#2586](https://github.com/sebitt27/dcrd/pull/2586))
- server: Respond to getheaders when same chain tip ([sebitt27/dcrd#2587](https://github.com/sebitt27/dcrd/pull/2587))
- netsync: Use an APBF for recently rejected txns ([sebitt27/dcrd#2590](https://github.com/sebitt27/dcrd/pull/2590))
- server: Only send fast block anns to full nodes ([sebitt27/dcrd#2606](https://github.com/sebitt27/dcrd/pull/2606))
- upnp: More accurate getOurIP ([sebitt27/dcrd#2571](https://github.com/sebitt27/dcrd/pull/2571))
- server: Correct tx not found ban reason ([sebitt27/dcrd#2677](https://github.com/sebitt27/dcrd/pull/2677))
- chaincfg: Add DCP0007 deployment ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- chaincfg: Introduce explicit ver upgrades agenda ([sebitt27/dcrd#2713](https://github.com/sebitt27/dcrd/pull/2713))
- blockchain: Implement reject new tx vers vote ([sebitt27/dcrd#2716](https://github.com/sebitt27/dcrd/pull/2716))
- blockchain: Implement reject new script vers vote ([sebitt27/dcrd#2716](https://github.com/sebitt27/dcrd/pull/2716))
- chaincfg: Add agenda for auto ticket revocations ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- multi: DCP0009 Auto revocations consensus change ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- chaincfg: Use single latest checkpoint ([sebitt27/dcrd#2762](https://github.com/sebitt27/dcrd/pull/2762))
- peer: Offset ping interval from idle timeout ([sebitt27/dcrd#2796](https://github.com/sebitt27/dcrd/pull/2796))
- chaincfg: Update checkpoint for upcoming release ([sebitt27/dcrd#2794](https://github.com/sebitt27/dcrd/pull/2794))
- chaincfg: Update min known chain work for release ([sebitt27/dcrd#2795](https://github.com/sebitt27/dcrd/pull/2795))
- netsync: Request init state immediately upon sync ([sebitt27/dcrd#2812](https://github.com/sebitt27/dcrd/pull/2812))
- blockchain: Reject old block vers for HFV ([sebitt27/dcrd#2752](https://github.com/sebitt27/dcrd/pull/2752))
- netsync: Rework next block download logic ([sebitt27/dcrd#2828](https://github.com/sebitt27/dcrd/pull/2828))
- chaincfg: Add AssumeValid param ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- chaincfg: Introduce subsidy split change agenda ([sebitt27/dcrd#2847](https://github.com/sebitt27/dcrd/pull/2847))
- multi: Implement DCP0010 subsidy consensus vote ([sebitt27/dcrd#2848](https://github.com/sebitt27/dcrd/pull/2848))
- server: Force PoW upgrade to v9 ([sebitt27/dcrd#2875](https://github.com/sebitt27/dcrd/pull/2875))

### Transaction relay (memory pool):

- mempool: Limit ancestor tracking in mempool ([sebitt27/dcrd#2458](https://github.com/sebitt27/dcrd/pull/2458))
- mempool: Remove old fix sequence lock rejection ([sebitt27/dcrd#2496](https://github.com/sebitt27/dcrd/pull/2496))
- mempool: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- mempool: Enforce explicit versions ([sebitt27/dcrd#2716](https://github.com/sebitt27/dcrd/pull/2716))
- mempool: Remove unneeded max tx ver std checks ([sebitt27/dcrd#2716](https://github.com/sebitt27/dcrd/pull/2716))
- mempool: Update fraud proof data ([sebitt27/dcrd#2804](https://github.com/sebitt27/dcrd/pull/2804))
- mempool: CheckTransactionInputs check fraud proof ([sebitt27/dcrd#2804](https://github.com/sebitt27/dcrd/pull/2804))

### Mining:

- mining: Move txPriorityQueue to a separate file ([sebitt27/dcrd#2431](https://github.com/sebitt27/dcrd/pull/2431))
- mining: Move interfaces to mining/interface.go ([sebitt27/dcrd#2431](https://github.com/sebitt27/dcrd/pull/2431))
- mining: Add method comments to blockManagerFacade ([sebitt27/dcrd#2431](https://github.com/sebitt27/dcrd/pull/2431))
- mining: Move BgBlkTmplGenerator to separate file ([sebitt27/dcrd#2431](https://github.com/sebitt27/dcrd/pull/2431))
- mining: Prevent panic in child prio item handling ([sebitt27/dcrd#2434](https://github.com/sebitt27/dcrd/pull/2434))
- mining: Add Config struct to house mining params ([sebitt27/dcrd#2436](https://github.com/sebitt27/dcrd/pull/2436))
- mining: Move block chain functions to Config ([sebitt27/dcrd#2436](https://github.com/sebitt27/dcrd/pull/2436))
- mining: Move txMiningView from mempool package ([sebitt27/dcrd#2467](https://github.com/sebitt27/dcrd/pull/2467))
- mining: Switch to custom waitGroup impl ([sebitt27/dcrd#2477](https://github.com/sebitt27/dcrd/pull/2477))
- mining: Remove leftover block manager facade iface ([sebitt27/dcrd#2510](https://github.com/sebitt27/dcrd/pull/2510))
- mining: No error log on expected head reorg errors ([sebitt27/dcrd#2560](https://github.com/sebitt27/dcrd/pull/2560))
- mining: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- mining: Add error kinds for auto revocations ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- mining: Add auto revocation priority to tx queue ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- mining: Add HeaderByHash to Config ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- mining: Prevent unnecessary reorg with equal votes ([sebitt27/dcrd#2840](https://github.com/sebitt27/dcrd/pull/2840))
- mining: Update to latest block vers for HFV ([sebitt27/dcrd#2753](https://github.com/sebitt27/dcrd/pull/2753))

### RPC:

- rpcserver: Upgrade is deprecated; switch to Upgrader ([sebitt27/dcrd#2409](https://github.com/sebitt27/dcrd/pull/2409))
- multi: Add TAdd support to getrawmempool ([sebitt27/dcrd#2448](https://github.com/sebitt27/dcrd/pull/2448))
- rpcserver: Update getrawmempool txtype help ([sebitt27/dcrd#2452](https://github.com/sebitt27/dcrd/pull/2452))
- rpcserver: Hash auth using random-keyed MAC ([sebitt27/dcrd#2486](https://github.com/sebitt27/dcrd/pull/2486))
- rpcserver: Use next stake diff from snapshot ([sebitt27/dcrd#2493](https://github.com/sebitt27/dcrd/pull/2493))
- rpcserver: Make authenticate match header auth ([sebitt27/dcrd#2502](https://github.com/sebitt27/dcrd/pull/2502))
- rpcserver: Check unauthorized access in const time ([sebitt27/dcrd#2509](https://github.com/sebitt27/dcrd/pull/2509))
- multi: Subscribe for work ntfns in rpcserver ([sebitt27/dcrd#2501](https://github.com/sebitt27/dcrd/pull/2501))
- rpcserver: Prune block templates in websocket path ([sebitt27/dcrd#2503](https://github.com/sebitt27/dcrd/pull/2503))
- rpcserver: Remove version from gettxout result ([sebitt27/dcrd#2517](https://github.com/sebitt27/dcrd/pull/2517))
- rpcserver: Add tree param to gettxout ([sebitt27/dcrd#2517](https://github.com/sebitt27/dcrd/pull/2517))
- rpcserver/netsync: Remove notifystakedifficulty ([sebitt27/dcrd#2519](https://github.com/sebitt27/dcrd/pull/2519))
- rpcserver: Remove v1 getcfilter{,header} ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- rpcserver: Remove unused Filterer interface ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- rpcserver: Update getblockchaininfo best header ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- rpcserver: Remove unused LocateBlocks iface method ([sebitt27/dcrd#2538](https://github.com/sebitt27/dcrd/pull/2538))
- rpcserver: Allow TLS client cert authentication ([sebitt27/dcrd#2482](https://github.com/sebitt27/dcrd/pull/2482))
- rpcserver: Add invalidate/reconsiderblock support ([sebitt27/dcrd#2536](https://github.com/sebitt27/dcrd/pull/2536))
- rpcserver: Support getblockchaininfo genesis block ([sebitt27/dcrd#2550](https://github.com/sebitt27/dcrd/pull/2550))
- rpcserver: Calc verify progress based on best hdr ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- rpcserver: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- rpcserver: Allow gettreasurybalance empty blk str ([sebitt27/dcrd#2640](https://github.com/sebitt27/dcrd/pull/2640))
- rpcserver: Add median time to verbose results ([sebitt27/dcrd#2638](https://github.com/sebitt27/dcrd/pull/2638))
- rpcserver: Allow interface names for dial addresses ([sebitt27/dcrd#2623](https://github.com/sebitt27/dcrd/pull/2623))
- rpcserver: Add script version to gettxout ([sebitt27/dcrd#2650](https://github.com/sebitt27/dcrd/pull/2650))
- rpcserver: Remove unused help entry ([sebitt27/dcrd#2648](https://github.com/sebitt27/dcrd/pull/2648))
- rpcserver: Set script version in raw tx results ([sebitt27/dcrd#2663](https://github.com/sebitt27/dcrd/pull/2663))
- rpcserver: Impose additional read limits ([sebitt27/dcrd#2675](https://github.com/sebitt27/dcrd/pull/2675))
- rpcserver: Add more strict request origin check ([sebitt27/dcrd#2676](https://github.com/sebitt27/dcrd/pull/2676))
- rpcserver: Use duplicate tx error for recently mined transactions ([sebitt27/dcrd#2705](https://github.com/sebitt27/dcrd/pull/2705))
- rpcserver: Wait for sync on rpc request ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- rpcserver: Update websocket ping timeout handling ([sebitt27/dcrd#2866](https://github.com/sebitt27/dcrd/pull/2866))

### dcrd command-line flags and configuration:

- multi: Rename BMGR subsystem to SYNC ([sebitt27/dcrd#2500](https://github.com/sebitt27/dcrd/pull/2500))
- server/indexers: Remove v1 cfilter indexing support ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- config: Add utxocachemaxsize ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- main: Update slog for LOGFLAGS=nodatetime support ([sebitt27/dcrd#2608](https://github.com/sebitt27/dcrd/pull/2608))
- config: Allow interface names for listener addresses ([sebitt27/dcrd#2623](https://github.com/sebitt27/dcrd/pull/2623))
- config: Correct dir create failure error message ([sebitt27/dcrd#2682](https://github.com/sebitt27/dcrd/pull/2682))
- config: Add logsize config option ([sebitt27/dcrd#2711](https://github.com/sebitt27/dcrd/pull/2711))
- config: conditionally generate rpc credentials ([sebitt27/dcrd#2779](https://github.com/sebitt27/dcrd/pull/2779))
- multi: Add assumevalid config option ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))

### gencerts utility changes:

- gencerts: Add certificate authority capabilities ([sebitt27/dcrd#2478](https://github.com/sebitt27/dcrd/pull/2478))
- gencerts: Add RSA support (4096 bit keys only) ([sebitt27/dcrd#2551](https://github.com/sebitt27/dcrd/pull/2551))

### addblock utility changes:

- cmd/addblock: update block importer ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- addblock: Run index subscriber as a goroutine ([sebitt27/dcrd#2760](https://github.com/sebitt27/dcrd/pull/2760))
- addblock: Fix blockchain initialization ([sebitt27/dcrd#2760](https://github.com/sebitt27/dcrd/pull/2760))
- addblock: Use chain bulk import mode ([sebitt27/dcrd#2782](https://github.com/sebitt27/dcrd/pull/2782))

### findcheckpoint utility changes:

- findcheckpoint: Fix blockchain initialization ([sebitt27/dcrd#2759](https://github.com/sebitt27/dcrd/pull/2759))

### Documentation:

- docs: Fix JSON-RPC API gettxoutsetinfo description ([sebitt27/dcrd#2443](https://github.com/sebitt27/dcrd/pull/2443))
- docs: Add JSON-RPC API getpeerinfo missing fields ([sebitt27/dcrd#2443](https://github.com/sebitt27/dcrd/pull/2443))
- docs: Fix JSON-RPC API gettreasurybalance fmt ([sebitt27/dcrd#2443](https://github.com/sebitt27/dcrd/pull/2443))
- docs: Fix JSON-RPC API gettreasuryspendvotes fmt ([sebitt27/dcrd#2443](https://github.com/sebitt27/dcrd/pull/2443))
- docs: Add JSON-RPC API searchrawtxns req limit ([sebitt27/dcrd#2443](https://github.com/sebitt27/dcrd/pull/2443))
- docs: Update JSON-RPC API getrawmempool ([sebitt27/dcrd#2453](https://github.com/sebitt27/dcrd/pull/2453))
- progresslog: Add package documentation ([sebitt27/dcrd#2499](https://github.com/sebitt27/dcrd/pull/2499))
- netsync: Add package documentation ([sebitt27/dcrd#2500](https://github.com/sebitt27/dcrd/pull/2500))
- multi: update error code related documentation ([sebitt27/dcrd#2515](https://github.com/sebitt27/dcrd/pull/2515))
- docs: Update JSON-RPC API getwork to match reality ([sebitt27/dcrd#2526](https://github.com/sebitt27/dcrd/pull/2526))
- docs: Remove notifystakedifficulty JSON-RPC API ([sebitt27/dcrd#2519](https://github.com/sebitt27/dcrd/pull/2519))
- docs: Remove v1 getcfilter{,header} JSON-RPC API ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- chaincfg: Update doc.go ([sebitt27/dcrd#2528](https://github.com/sebitt27/dcrd/pull/2528))
- blockchain: Update README.md and doc.go ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- docs: Add invalidate/reconsiderblock JSON-RPC API ([sebitt27/dcrd#2536](https://github.com/sebitt27/dcrd/pull/2536))
- docs: Add release notes for v1.6.0 ([sebitt27/dcrd#2451](https://github.com/sebitt27/dcrd/pull/2451))
- multi: Update README.md files for go modules ([sebitt27/dcrd#2559](https://github.com/sebitt27/dcrd/pull/2559))
- apbf: Add README.md ([sebitt27/dcrd#2579](https://github.com/sebitt27/dcrd/pull/2579))
- docs: Add release notes for v1.6.1 ([sebitt27/dcrd#2601](https://github.com/sebitt27/dcrd/pull/2601))
- docs: Update min recommended specs in README.md ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- stdaddr: Add README.md ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add serialized pubkey info to README.md ([sebitt27/dcrd#2619](https://github.com/sebitt27/dcrd/pull/2619))
- docs: Add release notes for v1.6.2 ([sebitt27/dcrd#2630](https://github.com/sebitt27/dcrd/pull/2630))
- docs: Add scriptpubkey json returns ([sebitt27/dcrd#2650](https://github.com/sebitt27/dcrd/pull/2650))
- stdscript: Add README.md ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stake: Comment on max SSGen outputs with treasury ([sebitt27/dcrd#2664](https://github.com/sebitt27/dcrd/pull/2664))
- docs: Update JSON-RPC API for script version ([sebitt27/dcrd#2663](https://github.com/sebitt27/dcrd/pull/2663))
- docs: Update JSON-RPC API for max request limits ([sebitt27/dcrd#2675](https://github.com/sebitt27/dcrd/pull/2675))
- docs: Add SECURITY.md file ([sebitt27/dcrd#2717](https://github.com/sebitt27/dcrd/pull/2717))
- sampleconfig: Add missing log options ([sebitt27/dcrd#2723](https://github.com/sebitt27/dcrd/pull/2723))
- docs: Update go versions in README.md ([sebitt27/dcrd#2722](https://github.com/sebitt27/dcrd/pull/2722))
- docs: Correct generate description ([sebitt27/dcrd#2724](https://github.com/sebitt27/dcrd/pull/2724))
- database: Correct README rpcclient link ([sebitt27/dcrd#2725](https://github.com/sebitt27/dcrd/pull/2725))
- docs: Add accuracy and reliability to README.md ([sebitt27/dcrd#2726](https://github.com/sebitt27/dcrd/pull/2726))
- sampleconfig: Update for deprecated nodnsseed ([sebitt27/dcrd#2728](https://github.com/sebitt27/dcrd/pull/2728))
- docs: Update for secp256k1 v4 module ([sebitt27/dcrd#2732](https://github.com/sebitt27/dcrd/pull/2732))
- docs: Update for new modules ([sebitt27/dcrd#2744](https://github.com/sebitt27/dcrd/pull/2744))
- sampleconfig: update rpc credentials documentation ([sebitt27/dcrd#2779](https://github.com/sebitt27/dcrd/pull/2779))
- docs: Update for addrmgr v2 module ([sebitt27/dcrd#2797](https://github.com/sebitt27/dcrd/pull/2797))
- docs: Update for rpc/jsonrpc/types v3 module ([sebitt27/dcrd#2801](https://github.com/sebitt27/dcrd/pull/2801))
- stdscript: Update README.md for provably pruneable ([sebitt27/dcrd#2803](https://github.com/sebitt27/dcrd/pull/2803))
- docs: Update for txscript v3 module ([sebitt27/dcrd#2815](https://github.com/sebitt27/dcrd/pull/2815))
- docs: Update for dcrutil v4 module ([sebitt27/dcrd#2818](https://github.com/sebitt27/dcrd/pull/2818))
- uint256: Add README.md ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- docs: Update for peer v3 module ([sebitt27/dcrd#2820](https://github.com/sebitt27/dcrd/pull/2820))
- docs: Update for database v3 module ([sebitt27/dcrd#2822](https://github.com/sebitt27/dcrd/pull/2822))
- docs: Update for blockchain/stake v4 module ([sebitt27/dcrd#2824](https://github.com/sebitt27/dcrd/pull/2824))
- docs: Update for gcs v3 module ([sebitt27/dcrd#2830](https://github.com/sebitt27/dcrd/pull/2830))
- docs: Fix typos and trailing whitespace ([sebitt27/dcrd#2843](https://github.com/sebitt27/dcrd/pull/2843))
- docs: Add max line length and wrapping guidelines ([sebitt27/dcrd#2843](https://github.com/sebitt27/dcrd/pull/2843))
- docs: Update for math/uint256 module ([sebitt27/dcrd#2842](https://github.com/sebitt27/dcrd/pull/2842))
- docs: Update simnet env docs for subsidy split ([sebitt27/dcrd#2848](https://github.com/sebitt27/dcrd/pull/2848))
- docs: Update for blockchain v4 module ([sebitt27/dcrd#2831](https://github.com/sebitt27/dcrd/pull/2831))
- docs: Update for rpcclient v7 module ([sebitt27/dcrd#2851](https://github.com/sebitt27/dcrd/pull/2851))
- primitives: Add skeleton README.md ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))

### Contrib changes:

- contrib: Update OpenBSD rc script for 6.9 features ([sebitt27/dcrd#2646](https://github.com/sebitt27/dcrd/pull/2646))
- contrib: Bump Dockerfile.alpine to alpine:3.14.0 ([sebitt27/dcrd#2681](https://github.com/sebitt27/dcrd/pull/2681))
- build: Use go 1.17 in Dockerfiles ([sebitt27/dcrd#2722](https://github.com/sebitt27/dcrd/pull/2722))
- build: Pin docker images with SHA instead of tag ([sebitt27/dcrd#2735](https://github.com/sebitt27/dcrd/pull/2735))
- build/contrib: Improve docker support ([sebitt27/dcrd#2740](https://github.com/sebitt27/dcrd/pull/2740))

### Developer-related package and module changes:

- dcrjson: Reject dup method type registrations ([sebitt27/dcrd#2417](https://github.com/sebitt27/dcrd/pull/2417))
- peer: various cleanups ([sebitt27/dcrd#2396](https://github.com/sebitt27/dcrd/pull/2396))
- blockchain: Create treasury buckets during upgrade ([sebitt27/dcrd#2441](https://github.com/sebitt27/dcrd/pull/2441))
- blockchain: Fix stxosToScriptSource ([sebitt27/dcrd#2444](https://github.com/sebitt27/dcrd/pull/2444))
- rpcserver: add NtfnManager interface ([sebitt27/dcrd#2410](https://github.com/sebitt27/dcrd/pull/2410))
- lru: Fix lookup race on small caches ([sebitt27/dcrd#2464](https://github.com/sebitt27/dcrd/pull/2464))
- gcs: update error types ([sebitt27/dcrd#2262](https://github.com/sebitt27/dcrd/pull/2262))
- main: Switch windows service dependency ([sebitt27/dcrd#2479](https://github.com/sebitt27/dcrd/pull/2479))
- blockchain: Simplify upgrade single run stage code ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- blockchain: Simplify upgrade batching logic ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- blockchain: Use new batching logic for filter init ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- blockchain: Use new batch logic for blkidx upgrade ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- blockchain: Use new batch logic for utxos upgrade ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- blockchain: Use new batch logic for spends upgrade ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- blockchain: Use new batch logic for clr failed ([sebitt27/dcrd#2457](https://github.com/sebitt27/dcrd/pull/2457))
- windows: Switch to os.Executable ([sebitt27/dcrd#2485](https://github.com/sebitt27/dcrd/pull/2485))
- blockchain: Revert fast add reversal ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Less order dependent full blocks tests ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Move context free tx sanity checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Move context free block sanity checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Rework contextual tx checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Move {coin,trsy}base contextual checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Move staketx-related contextual checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Move sigop-related contextual checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Make CheckBlockSanity context free ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Context free CheckTransactionSanity ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- blockchain: Move contextual treasury spend checks ([sebitt27/dcrd#2481](https://github.com/sebitt27/dcrd/pull/2481))
- mempool: Comment and stylistic updates ([sebitt27/dcrd#2480](https://github.com/sebitt27/dcrd/pull/2480))
- mining: Rename TxMiningView Remove method ([sebitt27/dcrd#2490](https://github.com/sebitt27/dcrd/pull/2490))
- mining: Unexport TxMiningView methods ([sebitt27/dcrd#2490](https://github.com/sebitt27/dcrd/pull/2490))
- mining: Update mergeUtxoView comment ([sebitt27/dcrd#2490](https://github.com/sebitt27/dcrd/pull/2490))
- blockchain: Consolidate deployment errors ([sebitt27/dcrd#2487](https://github.com/sebitt27/dcrd/pull/2487))
- blockchain: Consolidate unknown block errors ([sebitt27/dcrd#2487](https://github.com/sebitt27/dcrd/pull/2487))
- blockchain: Consolidate no filter errors ([sebitt27/dcrd#2487](https://github.com/sebitt27/dcrd/pull/2487))
- blockchain: Consolidate no treasury bal errors ([sebitt27/dcrd#2487](https://github.com/sebitt27/dcrd/pull/2487))
- blockchain: Convert to LRU block cache ([sebitt27/dcrd#2488](https://github.com/sebitt27/dcrd/pull/2488))
- blockchain: Remove unused error returns ([sebitt27/dcrd#2489](https://github.com/sebitt27/dcrd/pull/2489))
- blockmanager: Remove unused stakediff infra ([sebitt27/dcrd#2493](https://github.com/sebitt27/dcrd/pull/2493))
- server: Use next stake diff from snapshot ([sebitt27/dcrd#2493](https://github.com/sebitt27/dcrd/pull/2493))
- blockchain: Explicit hash in next stake diff calcs ([sebitt27/dcrd#2494](https://github.com/sebitt27/dcrd/pull/2494))
- blockchain: Explicit hash in LN agenda active func ([sebitt27/dcrd#2495](https://github.com/sebitt27/dcrd/pull/2495))
- blockmanager: Remove unused config field ([sebitt27/dcrd#2497](https://github.com/sebitt27/dcrd/pull/2497))
- blockmanager: Decouple block database code ([sebitt27/dcrd#2497](https://github.com/sebitt27/dcrd/pull/2497))
- blockmanager: Decouple from global config var ([sebitt27/dcrd#2497](https://github.com/sebitt27/dcrd/pull/2497))
- blockchain: Explicit hash in max block size func ([sebitt27/dcrd#2507](https://github.com/sebitt27/dcrd/pull/2507))
- progresslog: Make block progress log internal ([sebitt27/dcrd#2499](https://github.com/sebitt27/dcrd/pull/2499))
- server: Do not use unexported block manager cfg ([sebitt27/dcrd#2498](https://github.com/sebitt27/dcrd/pull/2498))
- blockmanager: Rework chain current logic ([sebitt27/dcrd#2498](https://github.com/sebitt27/dcrd/pull/2498))
- multi: Handle chain ntfn callback in server ([sebitt27/dcrd#2498](https://github.com/sebitt27/dcrd/pull/2498))
- server: Rename blockManager field to syncManager ([sebitt27/dcrd#2500](https://github.com/sebitt27/dcrd/pull/2500))
- server: Add temp sync manager interface ([sebitt27/dcrd#2500](https://github.com/sebitt27/dcrd/pull/2500))
- netsync: Split blockmanager into separate package ([sebitt27/dcrd#2500](https://github.com/sebitt27/dcrd/pull/2500))
- netsync: Rename blockManager to SyncManager ([sebitt27/dcrd#2500](https://github.com/sebitt27/dcrd/pull/2500))
- internal/ticketdb: update error types ([sebitt27/dcrd#2279](https://github.com/sebitt27/dcrd/pull/2279))
- secp256k1/ecdsa: update error types ([sebitt27/dcrd#2281](https://github.com/sebitt27/dcrd/pull/2281))
- secp256k1/schnorr: update error types ([sebitt27/dcrd#2282](https://github.com/sebitt27/dcrd/pull/2282))
- dcrjson: update error types ([sebitt27/dcrd#2271](https://github.com/sebitt27/dcrd/pull/2271))
- dcrec/secp256k1: update error types ([sebitt27/dcrd#2265](https://github.com/sebitt27/dcrd/pull/2265))
- blockchain/stake: update error types ([sebitt27/dcrd#2264](https://github.com/sebitt27/dcrd/pull/2264))
- multi: update database error types ([sebitt27/dcrd#2261](https://github.com/sebitt27/dcrd/pull/2261))
- blockchain: Remove unused treasury active func ([sebitt27/dcrd#2514](https://github.com/sebitt27/dcrd/pull/2514))
- stake: update ticket lottery errors ([sebitt27/dcrd#2433](https://github.com/sebitt27/dcrd/pull/2433))
- netsync: Improve is current detection ([sebitt27/dcrd#2513](https://github.com/sebitt27/dcrd/pull/2513))
- internal/mining: update mining error types ([sebitt27/dcrd#2515](https://github.com/sebitt27/dcrd/pull/2515))
- multi: sprinkle on more errors.As/Is ([sebitt27/dcrd#2522](https://github.com/sebitt27/dcrd/pull/2522))
- mining: Correct fee calculations during reorgs ([sebitt27/dcrd#2530](https://github.com/sebitt27/dcrd/pull/2530))
- fees: Remove deprecated DisableLog ([sebitt27/dcrd#2529](https://github.com/sebitt27/dcrd/pull/2529))
- rpcclient: Remove deprecated DisableLog ([sebitt27/dcrd#2527](https://github.com/sebitt27/dcrd/pull/2527))
- rpcclient: Remove notifystakedifficulty ([sebitt27/dcrd#2519](https://github.com/sebitt27/dcrd/pull/2519))
- rpc/jsonrpc/types: Remove notifystakedifficulty ([sebitt27/dcrd#2519](https://github.com/sebitt27/dcrd/pull/2519))
- netsync: Remove unneeded ForceReorganization ([sebitt27/dcrd#2520](https://github.com/sebitt27/dcrd/pull/2520))
- mining: Remove duplicate method ([sebitt27/dcrd#2520](https://github.com/sebitt27/dcrd/pull/2520))
- multi: use EstimateSmartFeeResult ([sebitt27/dcrd#2283](https://github.com/sebitt27/dcrd/pull/2283))
- rpcclient: Remove v1 getcfilter{,header} ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- rpc/jsonrpc/types: Remove v1 getcfilter{,header} ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- gcs: Remove unused v1 blockcf package ([sebitt27/dcrd#2525](https://github.com/sebitt27/dcrd/pull/2525))
- blockchain: Remove legacy sequence lock view ([sebitt27/dcrd#2534](https://github.com/sebitt27/dcrd/pull/2534))
- blockchain: Remove IsFixSeqLocksAgendaActive ([sebitt27/dcrd#2534](https://github.com/sebitt27/dcrd/pull/2534))
- blockchain: Explicit hash in estimate stake diff ([sebitt27/dcrd#2524](https://github.com/sebitt27/dcrd/pull/2524))
- netsync: Remove unneeded TipGeneration ([sebitt27/dcrd#2537](https://github.com/sebitt27/dcrd/pull/2537))
- netsync: Remove unused TicketPoolValue ([sebitt27/dcrd#2544](https://github.com/sebitt27/dcrd/pull/2544))
- netsync: Embed peers vs separate peer states ([sebitt27/dcrd#2541](https://github.com/sebitt27/dcrd/pull/2541))
- netsync/server: Update peer heights directly ([sebitt27/dcrd#2542](https://github.com/sebitt27/dcrd/pull/2542))
- netsync: Move proactive sigcache evict to server ([sebitt27/dcrd#2543](https://github.com/sebitt27/dcrd/pull/2543))
- blockchain: Add invalidate/reconsider infrastruct ([sebitt27/dcrd#2536](https://github.com/sebitt27/dcrd/pull/2536))
- rpc/jsonrpc/types: Add invalidate/reconsiderblock ([sebitt27/dcrd#2536](https://github.com/sebitt27/dcrd/pull/2536))
- netsync: Convert lifecycle to context ([sebitt27/dcrd#2545](https://github.com/sebitt27/dcrd/pull/2545))
- multi: Rework utxoset/view to use outpoints ([sebitt27/dcrd#2540](https://github.com/sebitt27/dcrd/pull/2540))
- blockchain: Remove compression version param ([sebitt27/dcrd#2547](https://github.com/sebitt27/dcrd/pull/2547))
- blockchain: Remove error from LatestBlockLocator ([sebitt27/dcrd#2548](https://github.com/sebitt27/dcrd/pull/2548))
- blockchain: Fix incorrect decompressScript calls ([sebitt27/dcrd#2552](https://github.com/sebitt27/dcrd/pull/2552))
- blockchain: Fix V3 spend journal migration ([sebitt27/dcrd#2552](https://github.com/sebitt27/dcrd/pull/2552))
- multi: Remove blockChain field from UtxoViewpoint ([sebitt27/dcrd#2553](https://github.com/sebitt27/dcrd/pull/2553))
- blockchain: Move UtxoEntry to a separate file ([sebitt27/dcrd#2553](https://github.com/sebitt27/dcrd/pull/2553))
- blockchain: Update UtxoEntry Clone method comment ([sebitt27/dcrd#2553](https://github.com/sebitt27/dcrd/pull/2553))
- progresslog: Make logger more generic ([sebitt27/dcrd#2555](https://github.com/sebitt27/dcrd/pull/2555))
- server: Remove several unused funcs ([sebitt27/dcrd#2561](https://github.com/sebitt27/dcrd/pull/2561))
- mempool: Store staged transactions as TxDesc ([sebitt27/dcrd#2319](https://github.com/sebitt27/dcrd/pull/2319))
- connmgr: Add func to iterate conn reqs ([sebitt27/dcrd#2562](https://github.com/sebitt27/dcrd/pull/2562))
- netsync: Correct check for needTx ([sebitt27/dcrd#2568](https://github.com/sebitt27/dcrd/pull/2568))
- rpcclient: Update EstimateSmartFee return type ([sebitt27/dcrd#2255](https://github.com/sebitt27/dcrd/pull/2255))
- server: Notify sync mgr later and track ntfn ([sebitt27/dcrd#2582](https://github.com/sebitt27/dcrd/pull/2582))
- apbf: Introduce Age-Partitioned Bloom Filters ([sebitt27/dcrd#2579](https://github.com/sebitt27/dcrd/pull/2579))
- apbf: Add basic usage example ([sebitt27/dcrd#2579](https://github.com/sebitt27/dcrd/pull/2579))
- apbf: Add support to go generate a KL table ([sebitt27/dcrd#2579](https://github.com/sebitt27/dcrd/pull/2579))
- apbf: Switch to fast reduce method ([sebitt27/dcrd#2584](https://github.com/sebitt27/dcrd/pull/2584))
- server: Remove unneeded child context ([sebitt27/dcrd#2593](https://github.com/sebitt27/dcrd/pull/2593))
- blockchain: Separate utxo state from tx flags ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Add utxoStateFresh to UtxoEntry ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Add size method to UtxoEntry ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Deep copy view entry script from tx ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Add utxoSetState to the database ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- multi: Add UtxoCache ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Make InitUtxoCache a UtxoCache method ([sebitt27/dcrd#2599](https://github.com/sebitt27/dcrd/pull/2599))
- blockchain: Add UtxoCacher interface ([sebitt27/dcrd#2599](https://github.com/sebitt27/dcrd/pull/2599))
- dcrutil: Correct ed25519 address constructor ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Introduce package infra for std addrs ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add infrastructure for v0 decoding ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2pk-ecdsa-secp256k1 support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2pk-ed25519 support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2pk-schnorr-secp256k1 support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2pkh-ecdsa-secp256k1 support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2pkh-ed25519 support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2pkh-schnorr-secp256k1 support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add v0 p2sh support ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- stdaddr: Add decode address example ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- txscript: Rename script bldr add data to unchecked ([sebitt27/dcrd#2611](https://github.com/sebitt27/dcrd/pull/2611))
- txscript: Add script bldr unchecked op add ([sebitt27/dcrd#2611](https://github.com/sebitt27/dcrd/pull/2611))
- rpcserver: Remove uncompressed pubkeys fast path ([sebitt27/dcrd#2617](https://github.com/sebitt27/dcrd/pull/2617))
- blockchain: Allow alternate tips for current check ([sebitt27/dcrd#2612](https://github.com/sebitt27/dcrd/pull/2612))
- txscript: Accept raw public keys in MultiSigScript ([sebitt27/dcrd#2615](https://github.com/sebitt27/dcrd/pull/2615))
- cpuminer: Remove unused MiningAddrs from Config ([sebitt27/dcrd#2616](https://github.com/sebitt27/dcrd/pull/2616))
- stdaddr: Add ability to obtain raw public key ([sebitt27/dcrd#2619](https://github.com/sebitt27/dcrd/pull/2619))
- stdaddr: Move from internal/staging to txscript ([sebitt27/dcrd#2620](https://github.com/sebitt27/dcrd/pull/2620))
- stdaddr: Accept vote and revoke limits separately ([sebitt27/dcrd#2624](https://github.com/sebitt27/dcrd/pull/2624))
- stake: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- indexers: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- blockchain: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- rpcclient: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- hdkeychain: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSStx ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSStxChange ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSSGen ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSSGenSHDirect ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSSGenPKHDirect ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSSRtx ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSSRtxPKHDirect ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToSSRtxSHDirect ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToAddrScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused PayToScriptHashScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToSchnorrPubKeyScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToEdwardsPubKeyScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToPubKeyScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToScriptHashScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToPubKeyHashSchnorrScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToPubKeyHashEdwardsScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused payToPubKeyHashScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused GenerateSStxAddrPush ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Remove unused ErrUnsupportedAddress ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- txscript: Break dcrutil dependency ([sebitt27/dcrd#2626](https://github.com/sebitt27/dcrd/pull/2626))
- stdaddr: Replace Address method with String ([sebitt27/dcrd#2633](https://github.com/sebitt27/dcrd/pull/2633))
- dcrutil: Convert to use new stdaddr package ([sebitt27/dcrd#2628](https://github.com/sebitt27/dcrd/pull/2628))
- dcrutil: Remove all code related to Address ([sebitt27/dcrd#2628](https://github.com/sebitt27/dcrd/pull/2628))
- blockchain: Trsy always inactive for genesis blk ([sebitt27/dcrd#2636](https://github.com/sebitt27/dcrd/pull/2636))
- blockchain: Use agenda flags for tx check context ([sebitt27/dcrd#2639](https://github.com/sebitt27/dcrd/pull/2639))
- blockchain: Move UTXO DB methods to separate file ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- blockchain: Move UTXO DB tests to separate file ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- ipc: Fix lifetimeEvent comments ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- blockchain: Add utxoDatabaseInfo ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- multi: Introduce UTXO database ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- blockchain: Decouple stxo and utxo migrations ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- multi: Migrate to UTXO database ([sebitt27/dcrd#2632](https://github.com/sebitt27/dcrd/pull/2632))
- main: Handle SIGHUP with clean shutdown ([sebitt27/dcrd#2645](https://github.com/sebitt27/dcrd/pull/2645))
- txscript: Split signing code to sign subpackage ([sebitt27/dcrd#2642](https://github.com/sebitt27/dcrd/pull/2642))
- database: Add Flush to DB interface ([sebitt27/dcrd#2649](https://github.com/sebitt27/dcrd/pull/2649))
- multi: Flush block DB before UTXO DB ([sebitt27/dcrd#2649](https://github.com/sebitt27/dcrd/pull/2649))
- blockchain: Flush UTXO DB after init utxoSetState ([sebitt27/dcrd#2649](https://github.com/sebitt27/dcrd/pull/2649))
- blockchain: Force flush in separateUtxoDatabase ([sebitt27/dcrd#2649](https://github.com/sebitt27/dcrd/pull/2649))
- version: Rework to support single version override ([sebitt27/dcrd#2651](https://github.com/sebitt27/dcrd/pull/2651))
- blockchain: Remove UtxoCacher DB Tx dependency ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add UtxoBackend interface ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Export UtxoSetState ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add FetchEntry to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add PutUtxos to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add FetchState to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add FetchStats to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add FetchInfo to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Move LoadUtxoDB to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Add Upgrade to UtxoBackend ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- multi: Remove UTXO db in BlockChain and UtxoCache ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- blockchain: Export ViewFilteredSet ([sebitt27/dcrd#2652](https://github.com/sebitt27/dcrd/pull/2652))
- stake: Return StakeAddress from cmtmt conversion ([sebitt27/dcrd#2655](https://github.com/sebitt27/dcrd/pull/2655))
- stdscript: Introduce pkg infra for std scripts ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pk-ecdsa-secp256k1 support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pk-ed25519 support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pk-schnorr-secp256k1 support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pkh-ecdsa-secp256k1 support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pkh-ed25519 support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pkh-schnorr-secp256k1 support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2sh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 ecdsa multisig support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 ecdsa multisig redeem support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 nulldata support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake sub p2pkh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake sub p2sh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake gen p2pkh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake gen p2sh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake revoke p2pkh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake revoke p2sh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake change p2pkh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake change p2sh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 treasury add support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 treasury gen p2pkh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 treasury gen p2sh support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add ecdsa multisig creation script ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 atomic swap redeem support ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add example for determining script type ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add example for p2pkh extract ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add example of script hash extract ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- blockchain: Use scripts in tickets address query ([sebitt27/dcrd#2657](https://github.com/sebitt27/dcrd/pull/2657))
- stake: Do not use standardness code in consensus ([sebitt27/dcrd#2658](https://github.com/sebitt27/dcrd/pull/2658))
- blockchain: Remove unneeded OP_TADD maturity check ([sebitt27/dcrd#2659](https://github.com/sebitt27/dcrd/pull/2659))
- stake: Add is treasury gen script ([sebitt27/dcrd#2660](https://github.com/sebitt27/dcrd/pull/2660))
- blockchain: No standardness code in consensus ([sebitt27/dcrd#2661](https://github.com/sebitt27/dcrd/pull/2661))
- gcs: No standardness code in consensus ([sebitt27/dcrd#2662](https://github.com/sebitt27/dcrd/pull/2662))
- stake: Remove stale TODOs from CheckSSGenVotes ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- stake: Remove stale TODOs from CheckSSRtx ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- txscript: Move contains stake opcode to consensus ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Move stake blockref script to consensus ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Move stake votebits script to consensus ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Remove unused IsPubKeyHashScript ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Remove unused IsStakeChangeScript ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Remove unused PushedData ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- blockchain: Flush UtxoCache when latch to current ([sebitt27/dcrd#2671](https://github.com/sebitt27/dcrd/pull/2671))
- dcrjson: Minor jsonerr.go update ([sebitt27/dcrd#2672](https://github.com/sebitt27/dcrd/pull/2672))
- rpcclient: Cancel client context on shutdown ([sebitt27/dcrd#2678](https://github.com/sebitt27/dcrd/pull/2678))
- blockchain: Remove serializeUtxoEntry error ([sebitt27/dcrd#2683](https://github.com/sebitt27/dcrd/pull/2683))
- blockchain: Add IsTreasuryEnabled to AgendaFlags ([sebitt27/dcrd#2686](https://github.com/sebitt27/dcrd/pull/2686))
- multi: Update block ntfns to contain AgendaFlags ([sebitt27/dcrd#2686](https://github.com/sebitt27/dcrd/pull/2686))
- multi: Update ProcessOrphans to use AgendaFlags ([sebitt27/dcrd#2686](https://github.com/sebitt27/dcrd/pull/2686))
- mempool: Add maybeAcceptTransaction AgendaFlags ([sebitt27/dcrd#2686](https://github.com/sebitt27/dcrd/pull/2686))
- secp256k1: Allow code generation to compile again ([sebitt27/dcrd#2687](https://github.com/sebitt27/dcrd/pull/2687))
- jsonrpc/types: Add missing Method type to vars ([sebitt27/dcrd#2688](https://github.com/sebitt27/dcrd/pull/2688))
- blockchain: Add UTXO backend error kinds ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Add helper to convert leveldb errors ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Add UtxoBackendIterator interface ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Add UtxoBackendTx interface ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Add levelDbUtxoBackendTx type ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- multi: Update UtxoBackend to use leveldb directly ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- multi: Move UTXO database ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Unexport levelDbUtxoBackend ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Always use node lookup methods ([sebitt27/dcrd#2685](https://github.com/sebitt27/dcrd/pull/2685))
- blockchain: Use short keys for block index ([sebitt27/dcrd#2685](https://github.com/sebitt27/dcrd/pull/2685))
- rpcclient: Shutdown breaks reconnect sleep ([sebitt27/dcrd#2696](https://github.com/sebitt27/dcrd/pull/2696))
- secp256k1: No deps on adaptor code for precomps ([sebitt27/dcrd#2690](https://github.com/sebitt27/dcrd/pull/2690))
- secp256k1: Always initialize adaptor instance ([sebitt27/dcrd#2690](https://github.com/sebitt27/dcrd/pull/2690))
- secp256k1: Optimize precomp values to use affine ([sebitt27/dcrd#2690](https://github.com/sebitt27/dcrd/pull/2690))
- rpcserver: Handle getwork nil err during reorg ([sebitt27/dcrd#2700](https://github.com/sebitt27/dcrd/pull/2700))
- secp256k1: Use blake256 directly in examples ([sebitt27/dcrd#2697](https://github.com/sebitt27/dcrd/pull/2697))
- secp256k1: Improve scalar mult readability ([sebitt27/dcrd#2695](https://github.com/sebitt27/dcrd/pull/2695))
- secp256k1: Optimize NAF conversion ([sebitt27/dcrd#2695](https://github.com/sebitt27/dcrd/pull/2695))
- blockchain: Verify state of DCP0007 voting ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- blockchain: Rename max expenditure funcs ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- stake: Add ExpiringNextBlock method to Node ([sebitt27/dcrd#2701](https://github.com/sebitt27/dcrd/pull/2701))
- rpcclient: Add GetNetworkInfo call ([sebitt27/dcrd#2703](https://github.com/sebitt27/dcrd/pull/2703))
- stake: Pre-allocate lottery ticket index slice ([sebitt27/dcrd#2710](https://github.com/sebitt27/dcrd/pull/2710))
- blockchain: Switch to treasuryValueType.IsDebit ([sebitt27/dcrd#2680](https://github.com/sebitt27/dcrd/pull/2680))
- blockchain: Sum amounts added to treasury ([sebitt27/dcrd#2680](https://github.com/sebitt27/dcrd/pull/2680))
- blockchain: Add maxTreasuryExpenditureDCP0007 ([sebitt27/dcrd#2680](https://github.com/sebitt27/dcrd/pull/2680))
- blockchain: Use new expenditure policy if activated ([sebitt27/dcrd#2680](https://github.com/sebitt27/dcrd/pull/2680))
- blockchain: Add checkTicketRedeemers ([sebitt27/dcrd#2702](https://github.com/sebitt27/dcrd/pull/2702))
- blockchain: Add NextExpiringTickets to BestState ([sebitt27/dcrd#2708](https://github.com/sebitt27/dcrd/pull/2708))
- multi: Add FetchUtxoEntry to mining Config ([sebitt27/dcrd#2709](https://github.com/sebitt27/dcrd/pull/2709))
- stake: Add func to create revocation from ticket ([sebitt27/dcrd#2707](https://github.com/sebitt27/dcrd/pull/2707))
- rpcserver: Use CreateRevocationFromTicket ([sebitt27/dcrd#2707](https://github.com/sebitt27/dcrd/pull/2707))
- multi: Don't use deprecated ioutil package ([sebitt27/dcrd#2722](https://github.com/sebitt27/dcrd/pull/2722))
- blockchain: Consolidate tx check flag construction ([sebitt27/dcrd#2716](https://github.com/sebitt27/dcrd/pull/2716))
- stake: Export Hash256PRNG UniformRandom ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- blockchain: Check auto revocations agenda state ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- multi: Add mempool IsAutoRevocationsAgendaActive ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- multi: Add auto revocations to agenda flags ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- multi: Check tx inputs auto revocations flag ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- blockchain: Add auto revocation error kinds ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- stake: Add auto revocation error kinds ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- blockchain: Move revocation checks block context ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- multi: Add isAutoRevocationsEnabled to CheckSSRtx ([sebitt27/dcrd#2719](https://github.com/sebitt27/dcrd/pull/2719))
- addrmgr: Remove deprecated code ([sebitt27/dcrd#2729](https://github.com/sebitt27/dcrd/pull/2729))
- peer: Remove deprecated DisableLog ([sebitt27/dcrd#2730](https://github.com/sebitt27/dcrd/pull/2730))
- database: Remove deprecated DisableLog ([sebitt27/dcrd#2731](https://github.com/sebitt27/dcrd/pull/2731))
- addrmgr: Decouple IP network checks from wire ([sebitt27/dcrd#2596](https://github.com/sebitt27/dcrd/pull/2596))
- addrmgr: Rename network address type ([sebitt27/dcrd#2596](https://github.com/sebitt27/dcrd/pull/2596))
- addrmgr: Decouple addrmgr from wire NetAddress ([sebitt27/dcrd#2596](https://github.com/sebitt27/dcrd/pull/2596))
- multi: add spend pruner ([sebitt27/dcrd#2641](https://github.com/sebitt27/dcrd/pull/2641))
- multi: synchronize spend prunes and notifications ([sebitt27/dcrd#2641](https://github.com/sebitt27/dcrd/pull/2641))
- blockchain: workSorterLess -> betterCandidate ([sebitt27/dcrd#2747](https://github.com/sebitt27/dcrd/pull/2747))
- mempool: Add HeaderByHash to Config ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- rpctest: Remove unused BlockVersion const ([sebitt27/dcrd#2754](https://github.com/sebitt27/dcrd/pull/2754))
- blockchain: Handle genesis auto revocation agenda ([sebitt27/dcrd#2755](https://github.com/sebitt27/dcrd/pull/2755))
- indexers: remove index manager ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- indexers: add index subscriber ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- indexers: refactor interfaces ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- indexers: async transaction index ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- indexers: update address index ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- indexers: async exists address index ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- multi: integrate index subscriber ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- multi: avoid using subscriber lifecycle in catchup ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- multi: remove spend deps on index disc. notifs ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- multi: copy snapshot pkScript ([sebitt27/dcrd#2219](https://github.com/sebitt27/dcrd/pull/2219))
- blockchain: Conditionally log difficulty retarget ([sebitt27/dcrd#2761](https://github.com/sebitt27/dcrd/pull/2761))
- multi: Use single latest checkpoint ([sebitt27/dcrd#2763](https://github.com/sebitt27/dcrd/pull/2763))
- blockchain: Move diff retarget log to connect ([sebitt27/dcrd#2765](https://github.com/sebitt27/dcrd/pull/2765))
- multi: source index notif. from block notif ([sebitt27/dcrd#2256](https://github.com/sebitt27/dcrd/pull/2256))
- server: fix wireToAddrmgrNetAddress data race ([sebitt27/dcrd#2758](https://github.com/sebitt27/dcrd/pull/2758))
- multi: Flush cache before fetching UTXO stats ([sebitt27/dcrd#2767](https://github.com/sebitt27/dcrd/pull/2767))
- blockchain: Don't use deprecated ioutil package ([sebitt27/dcrd#2769](https://github.com/sebitt27/dcrd/pull/2769))
- blockchain: Fix ticket db disconnect revocations ([sebitt27/dcrd#2768](https://github.com/sebitt27/dcrd/pull/2768))
- blockchain: Add convenience ancestor of func ([sebitt27/dcrd#2771](https://github.com/sebitt27/dcrd/pull/2771))
- blockchain: Use new ancestor of convenience func ([sebitt27/dcrd#2771](https://github.com/sebitt27/dcrd/pull/2771))
- blockchain: Remove unused latest blk locator func ([sebitt27/dcrd#2772](https://github.com/sebitt27/dcrd/pull/2772))
- blockchain: Remove unused next lottery data func ([sebitt27/dcrd#2773](https://github.com/sebitt27/dcrd/pull/2773))
- secp256k1: Correct 96-bit accum double overflow ([sebitt27/dcrd#2778](https://github.com/sebitt27/dcrd/pull/2778))
- blockchain: Further decouple upgrade code ([sebitt27/dcrd#2776](https://github.com/sebitt27/dcrd/pull/2776))
- blockchain: Add bulk import mode ([sebitt27/dcrd#2782](https://github.com/sebitt27/dcrd/pull/2782))
- multi: Remove flags from SyncManager ProcessBlock ([sebitt27/dcrd#2783](https://github.com/sebitt27/dcrd/pull/2783))
- netsync: Remove flags from processBlockMsg ([sebitt27/dcrd#2783](https://github.com/sebitt27/dcrd/pull/2783))
- multi: Remove flags from blockchain ProcessBlock ([sebitt27/dcrd#2783](https://github.com/sebitt27/dcrd/pull/2783))
- multi: Remove flags from ProcessBlockHeader ([sebitt27/dcrd#2785](https://github.com/sebitt27/dcrd/pull/2785))
- blockchain: Remove flags maybeAcceptBlockHeader ([sebitt27/dcrd#2785](https://github.com/sebitt27/dcrd/pull/2785))
- version: Use uint32 for major/minor/patch ([sebitt27/dcrd#2789](https://github.com/sebitt27/dcrd/pull/2789))
- wire: Write message header directly ([sebitt27/dcrd#2790](https://github.com/sebitt27/dcrd/pull/2790))
- stake: Correct treasury enabled vote discovery ([sebitt27/dcrd#2780](https://github.com/sebitt27/dcrd/pull/2780))
- blockchain: Correct treasury spend vote data ([sebitt27/dcrd#2780](https://github.com/sebitt27/dcrd/pull/2780))
- blockchain: UTXO database migration fix ([sebitt27/dcrd#2798](https://github.com/sebitt27/dcrd/pull/2798))
- blockchain: Handle zero-length UTXO backend state ([sebitt27/dcrd#2798](https://github.com/sebitt27/dcrd/pull/2798))
- mining: Remove unnecessary tx copy ([sebitt27/dcrd#2792](https://github.com/sebitt27/dcrd/pull/2792))
- multi: Use dcrutil Tx in NewTxDeepTxIns ([sebitt27/dcrd#2802](https://github.com/sebitt27/dcrd/pull/2802))
- indexers: synchronize index subscriber ntfn sends/receives ([sebitt27/dcrd#2806](https://github.com/sebitt27/dcrd/pull/2806))
- stdscript: Add exported MaxDataCarrierSizeV0 ([sebitt27/dcrd#2803](https://github.com/sebitt27/dcrd/pull/2803))
- stdscript: Add ProvablyPruneableScriptV0 ([sebitt27/dcrd#2803](https://github.com/sebitt27/dcrd/pull/2803))
- stdscript: Add num required sigs support ([sebitt27/dcrd#2805](https://github.com/sebitt27/dcrd/pull/2805))
- netsync: Remove unused RpcServer ([sebitt27/dcrd#2811](https://github.com/sebitt27/dcrd/pull/2811))
- netsync: Consolidate initial sync handling ([sebitt27/dcrd#2812](https://github.com/sebitt27/dcrd/pull/2812))
- stdscript: Add v0 p2pk-ed25519 extract ([sebitt27/dcrd#2807](https://github.com/sebitt27/dcrd/pull/2807))
- stdscript: Add v0 p2pk-schnorr-secp256k1 extract ([sebitt27/dcrd#2807](https://github.com/sebitt27/dcrd/pull/2807))
- stdscript: Add v0 p2pkh-ed25519 extract ([sebitt27/dcrd#2807](https://github.com/sebitt27/dcrd/pull/2807))
- stdscript: Add v0 p2pkh-schnorr-secp256k1 extract ([sebitt27/dcrd#2807](https://github.com/sebitt27/dcrd/pull/2807))
- stdscript: Add script to address conversion ([sebitt27/dcrd#2807](https://github.com/sebitt27/dcrd/pull/2807))
- stdscript: Move from internal/staging to txscript ([sebitt27/dcrd#2810](https://github.com/sebitt27/dcrd/pull/2810))
- mining: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- mempool: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- chaingen: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- blockchain: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- indexers: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- indexers: Remove unused trsy enabled params ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript/sign: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript/sign: Remove unused trsy enabled params ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- rpcserver: Convert to use stdscript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove deprecated ExtractAtomicSwapDataPushes ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused GenerateProvablyPruneableOut ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused MultiSigScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused MultisigRedeemScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused CalcMultiSigStats ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused IsMultisigScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused IsMultisigSigScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused ExtractPkScriptAltSigType ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused GetScriptClass ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused GetStakeOutSubclass ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused typeOfScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isTreasurySpendScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isMultisigScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused ExtractPkScriptAddrs ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused scriptHashToAddrs ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused pubKeyHashToAddrs ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isTreasuryAddScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractMultisigScriptDetails ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isStakeChangeScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isPubKeyHashScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isStakeRevocationScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isStakeGenScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isStakeSubmissionScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractStakeScriptHash ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractStakePubKeyHash ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isNullDataScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractPubKeyHash ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isPubKeyAltScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractPubKeyAltDetails ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isPubKeyScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractPubKey ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractUncompressedPubKey ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractCompressedPubKey ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isPubKeyHashAltScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused extractPubKeyHashAltDetails ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused isStandardAltSignatureType ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused MaxDataCarrierSize ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused ScriptClass ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused ErrNotMultisigScript ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused ErrTooManyRequiredSigs ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- txscript: Remove unused ErrTooMuchNullData ([sebitt27/dcrd#2808](https://github.com/sebitt27/dcrd/pull/2808))
- stdaddr: Use txscript for opcode definitions ([sebitt27/dcrd#2809](https://github.com/sebitt27/dcrd/pull/2809))
- stdscript: Add v0 stake-tagged p2pkh extract ([sebitt27/dcrd#2816](https://github.com/sebitt27/dcrd/pull/2816))
- stdscript: Add v0 stake-tagged p2sh extract ([sebitt27/dcrd#2816](https://github.com/sebitt27/dcrd/pull/2816))
- server: sync rebroadcast inv sends/receives ([sebitt27/dcrd#2814](https://github.com/sebitt27/dcrd/pull/2814))
- multi: Move last ann block from peer to netsync ([sebitt27/dcrd#2821](https://github.com/sebitt27/dcrd/pull/2821))
- uint256: Introduce package infrastructure ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add set from big endian bytes ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add set from little endian bytes ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add get big endian bytes ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add get little endian bytes ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add zero support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add uint32 casting support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add uint64 casting support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add equality comparison support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add less than comparison support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add less or equals comparison support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add greater than comparison support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add greater or equals comparison support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add general comparison support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add addition support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add subtraction support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add multiplication support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add squaring support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add division support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add negation support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add is odd support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise left shift support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise right shift support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise not support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise or support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise and support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise xor support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bit length support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add text formatting support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add conversion to stdlib big int support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add conversion from stdlib big int support ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add basic usage example ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- stake: Rename func to identify stake cmtmnt output ([sebitt27/dcrd#2824](https://github.com/sebitt27/dcrd/pull/2824))
- progresslog: Make header logging concurrent safe ([sebitt27/dcrd#2833](https://github.com/sebitt27/dcrd/pull/2833))
- netsync: Contiguous hashes for initial state reqs ([sebitt27/dcrd#2825](https://github.com/sebitt27/dcrd/pull/2825))
- multi: Allow discrete mining with invalidated tip ([sebitt27/dcrd#2838](https://github.com/sebitt27/dcrd/pull/2838))
- primitives: Add difficulty bits <-> uint256 ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add work calc from diff bits ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add hash to uint256 conversion ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add check proof of work ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add core merkle tree root calcs ([sebitt27/dcrd#2826](https://github.com/sebitt27/dcrd/pull/2826))
- primitives: Add inclusion proof funcs ([sebitt27/dcrd#2827](https://github.com/sebitt27/dcrd/pull/2827))
- indexers: update indexer error types ([sebitt27/dcrd#2770](https://github.com/sebitt27/dcrd/pull/2770))
- rpcserver: Submit transactions directly ([sebitt27/dcrd#2835](https://github.com/sebitt27/dcrd/pull/2835))
- netsync: Remove unused tx submission processing ([sebitt27/dcrd#2835](https://github.com/sebitt27/dcrd/pull/2835))
- internal/staging: add ban manager ([sebitt27/dcrd#2554](https://github.com/sebitt27/dcrd/pull/2554))
- uint256: Correct base 10 output formatting ([sebitt27/dcrd#2844](https://github.com/sebitt27/dcrd/pull/2844))
- multi: Add assumeValid to BlockChain ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- blockchain: Track assumed valid node ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- blockchain: Set BFFastAdd based on assume valid ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- blockchain: Assume valid skip script validation ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- blockchain: Bulk import skip script validation ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- hdkeychain: Add a strict BIP32 child derivation method ([sebitt27/dcrd#2845](https://github.com/sebitt27/dcrd/pull/2845))
- mempool: Consolidate tx check flag construction ([sebitt27/dcrd#2846](https://github.com/sebitt27/dcrd/pull/2846))
- standalone: Add modified subsidy split calcs ([sebitt27/dcrd#2848](https://github.com/sebitt27/dcrd/pull/2848))

### Developer-related module management:

- rpcclient: Prepare v6.0.1 ([sebitt27/dcrd#2455](https://github.com/sebitt27/dcrd/pull/2455))
- multi: Start blockchain v4 module dev cycle ([sebitt27/dcrd#2463](https://github.com/sebitt27/dcrd/pull/2463))
- multi: Start rpcclient v7 module dev cycle ([sebitt27/dcrd#2463](https://github.com/sebitt27/dcrd/pull/2463))
- multi: Start gcs v3 module dev cycle ([sebitt27/dcrd#2463](https://github.com/sebitt27/dcrd/pull/2463))
- multi: Start blockchain/stake v4 module dev cycle ([sebitt27/dcrd#2511](https://github.com/sebitt27/dcrd/pull/2511))
- multi: Start txscript v4 module dev cycle ([sebitt27/dcrd#2511](https://github.com/sebitt27/dcrd/pull/2511))
- multi: Start dcrutil v4 module dev cycle ([sebitt27/dcrd#2511](https://github.com/sebitt27/dcrd/pull/2511))
- multi: Start dcrec/secp256k1 v4 module dev cycle ([sebitt27/dcrd#2511](https://github.com/sebitt27/dcrd/pull/2511))
- rpc/jsonrpc/types: Start v3 module dev cycle ([sebitt27/dcrd#2517](https://github.com/sebitt27/dcrd/pull/2517))
- multi: Round 1 prerel module release ver updates ([sebitt27/dcrd#2569](https://github.com/sebitt27/dcrd/pull/2569))
- multi: Round 2 prerel module release ver updates ([sebitt27/dcrd#2570](https://github.com/sebitt27/dcrd/pull/2570))
- multi: Round 3 prerel module release ver updates ([sebitt27/dcrd#2572](https://github.com/sebitt27/dcrd/pull/2572))
- multi: Round 4 prerel module release ver updates ([sebitt27/dcrd#2573](https://github.com/sebitt27/dcrd/pull/2573))
- multi: Round 5 prerel module release ver updates ([sebitt27/dcrd#2574](https://github.com/sebitt27/dcrd/pull/2574))
- multi: Round 6 prerel module release ver updates ([sebitt27/dcrd#2575](https://github.com/sebitt27/dcrd/pull/2575))
- multi: Update to siphash v1.2.2 ([sebitt27/dcrd#2577](https://github.com/sebitt27/dcrd/pull/2577))
- peer: Start v3 module dev cycle ([sebitt27/dcrd#2585](https://github.com/sebitt27/dcrd/pull/2585))
- addrmgr: Start v2 module dev cycle ([sebitt27/dcrd#2592](https://github.com/sebitt27/dcrd/pull/2592))
- blockchain: Prerel module release ver updates ([sebitt27/dcrd#2634](https://github.com/sebitt27/dcrd/pull/2634))
- blockchain: Bump database module minor version ([sebitt27/dcrd#2654](https://github.com/sebitt27/dcrd/pull/2654))
- multi: Require last database/v2.0.3-x version ([sebitt27/dcrd#2689](https://github.com/sebitt27/dcrd/pull/2689))
- multi: Introduce database/v3 module ([sebitt27/dcrd#2689](https://github.com/sebitt27/dcrd/pull/2689))
- multi: Use database/v3 module ([sebitt27/dcrd#2693](https://github.com/sebitt27/dcrd/pull/2693))
- main: Use pseudo-versions in bumped mods ([sebitt27/dcrd#2698](https://github.com/sebitt27/dcrd/pull/2698))
- blockchain: Add replace to chaincfg dependency ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- dcrjson: Introduce v4 module ([sebitt27/dcrd#2733](https://github.com/sebitt27/dcrd/pull/2733))
- secp256k1: Prepare v4.0.0 ([sebitt27/dcrd#2732](https://github.com/sebitt27/dcrd/pull/2732))
- docs: Update for dcrjson v4 module ([sebitt27/dcrd#2734](https://github.com/sebitt27/dcrd/pull/2734))
- dcrjson: Prepare v4.0.0 ([sebitt27/dcrd#2734](https://github.com/sebitt27/dcrd/pull/2734))
- blockchain: Prerel module release ver updates ([sebitt27/dcrd#2748](https://github.com/sebitt27/dcrd/pull/2748))
- gcs: Prerel module release ver updates ([sebitt27/dcrd#2749](https://github.com/sebitt27/dcrd/pull/2749))
- multi: Update gcs prerel version ([sebitt27/dcrd#2750](https://github.com/sebitt27/dcrd/pull/2750))
- multi: update build tags to pref. go1.17 syntax ([sebitt27/dcrd#2764](https://github.com/sebitt27/dcrd/pull/2764))
- chaincfg: Prepare v3.1.0 ([sebitt27/dcrd#2799](https://github.com/sebitt27/dcrd/pull/2799))
- addrmgr: Prepare v2.0.0 ([sebitt27/dcrd#2797](https://github.com/sebitt27/dcrd/pull/2797))
- rpc/jsonrpc/types: Prepare v3.0.0 ([sebitt27/dcrd#2801](https://github.com/sebitt27/dcrd/pull/2801))
- txscript: Prepare v4.0.0 ([sebitt27/dcrd#2815](https://github.com/sebitt27/dcrd/pull/2815))
- hdkeychain: Prepare v3.0.1 ([sebitt27/dcrd#2817](https://github.com/sebitt27/dcrd/pull/2817))
- dcrutil: Prepare v4.0.0 ([sebitt27/dcrd#2818](https://github.com/sebitt27/dcrd/pull/2818))
- connmgr: Prepare v3.1.0 ([sebitt27/dcrd#2819](https://github.com/sebitt27/dcrd/pull/2819))
- peer: Prepare v3.0.0 ([sebitt27/dcrd#2820](https://github.com/sebitt27/dcrd/pull/2820))
- database: Prepare v3.0.0 ([sebitt27/dcrd#2822](https://github.com/sebitt27/dcrd/pull/2822))
- blockchain/stake: Prepare v4.0.0 ([sebitt27/dcrd#2824](https://github.com/sebitt27/dcrd/pull/2824))
- gcs: Prepare v3.0.0 ([sebitt27/dcrd#2830](https://github.com/sebitt27/dcrd/pull/2830))
- math/uint256: Prepare v1.0.0 ([sebitt27/dcrd#2842](https://github.com/sebitt27/dcrd/pull/2842))
- blockchain: Prepare v4.0.0 ([sebitt27/dcrd#2831](https://github.com/sebitt27/dcrd/pull/2831))
- rpcclient: Prepare v7.0.0 ([sebitt27/dcrd#2851](https://github.com/sebitt27/dcrd/pull/2851))
- version: Include VCS build info in version string ([sebitt27/dcrd#2841](https://github.com/sebitt27/dcrd/pull/2841))
- main: Update to use all new module versions ([sebitt27/dcrd#2853](https://github.com/sebitt27/dcrd/pull/2853))
- main: Remove module replacements ([sebitt27/dcrd#2855](https://github.com/sebitt27/dcrd/pull/2855))

### Testing and Quality Assurance:

- rpcserver: Add handleGetTreasuryBalance tests ([sebitt27/dcrd#2390](https://github.com/sebitt27/dcrd/pull/2390))
- rpcserver: Add handleGet{Generate,HashesPerSec} tests ([sebitt27/dcrd#2365](https://github.com/sebitt27/dcrd/pull/2365))
- mining: Cleanup txPriorityQueue tests ([sebitt27/dcrd#2431](https://github.com/sebitt27/dcrd/pull/2431))
- blockchain: fix errorlint warnings ([sebitt27/dcrd#2411](https://github.com/sebitt27/dcrd/pull/2411))
- rpcserver: Add handleGetHeaders test ([sebitt27/dcrd#2366](https://github.com/sebitt27/dcrd/pull/2366))
- rpcserver: add ticketsforaddress tests ([sebitt27/dcrd#2405](https://github.com/sebitt27/dcrd/pull/2405))
- rpcserver: add ticketvwap tests ([sebitt27/dcrd#2406](https://github.com/sebitt27/dcrd/pull/2406))
- rpcserver: add handleTxFeeInfo tests ([sebitt27/dcrd#2407](https://github.com/sebitt27/dcrd/pull/2407))
- rpcserver: add handleTicketFeeInfo tests ([sebitt27/dcrd#2408](https://github.com/sebitt27/dcrd/pull/2408))
- rpcserver: add handleVerifyMessage tests ([sebitt27/dcrd#2413](https://github.com/sebitt27/dcrd/pull/2413))
- rpcserver: add handleSendRawTransaction tests ([sebitt27/dcrd#2410](https://github.com/sebitt27/dcrd/pull/2410))
- rpcserver: add handleGetVoteInfo tests ([sebitt27/dcrd#2432](https://github.com/sebitt27/dcrd/pull/2432))
- database: Fix errorlint warnings ([sebitt27/dcrd#2484](https://github.com/sebitt27/dcrd/pull/2484))
- mining: Add mining test harness ([sebitt27/dcrd#2480](https://github.com/sebitt27/dcrd/pull/2480))
- mining: Add NewBlockTemplate tests ([sebitt27/dcrd#2480](https://github.com/sebitt27/dcrd/pull/2480))
- mining: Move TxMiningView tests to mining ([sebitt27/dcrd#2480](https://github.com/sebitt27/dcrd/pull/2480))
- rpcserver: add handleGetRawTransaction tests ([sebitt27/dcrd#2483](https://github.com/sebitt27/dcrd/pull/2483))
- blockchain: Improve synthetic treasury vote tests ([sebitt27/dcrd#2488](https://github.com/sebitt27/dcrd/pull/2488))
- rpcserver: Add handleGetMempoolInfo test ([sebitt27/dcrd#2492](https://github.com/sebitt27/dcrd/pull/2492))
- connmgr: Increase test timeouts ([sebitt27/dcrd#2505](https://github.com/sebitt27/dcrd/pull/2505))
- run_tests.sh: Avoid command substitution ([sebitt27/dcrd#2506](https://github.com/sebitt27/dcrd/pull/2506))
- mempool: Make sequence lock tests more consistent ([sebitt27/dcrd#2496](https://github.com/sebitt27/dcrd/pull/2496))
- mempool: Rework sequence lock acceptance tests ([sebitt27/dcrd#2496](https://github.com/sebitt27/dcrd/pull/2496))
- rpcserver: Add handleGetTxOut tests ([sebitt27/dcrd#2516](https://github.com/sebitt27/dcrd/pull/2516))
- rpcserver: Add handleGetNetworkHashPS test ([sebitt27/dcrd#2512](https://github.com/sebitt27/dcrd/pull/2512))
- rpcserver: Add handleGetMiningInfo test ([sebitt27/dcrd#2512](https://github.com/sebitt27/dcrd/pull/2512))
- blockchain: Simplify TestFixedSequenceLocks ([sebitt27/dcrd#2534](https://github.com/sebitt27/dcrd/pull/2534))
- chaingen: Support querying block test name by hash ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- blockchain: Improve test harness logging ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- blockchain: Support separate test block generation ([sebitt27/dcrd#2518](https://github.com/sebitt27/dcrd/pull/2518))
- rpcserver: add handleVersion, handleHelp rpc tests ([sebitt27/dcrd#2549](https://github.com/sebitt27/dcrd/pull/2549))
- blockchain: Use ReplaceVoteBits in utxoview tests ([sebitt27/dcrd#2553](https://github.com/sebitt27/dcrd/pull/2553))
- blockchain: Add unit test coverage for UtxoEntry ([sebitt27/dcrd#2553](https://github.com/sebitt27/dcrd/pull/2553))
- rpctest: Don't use installed node ([sebitt27/dcrd#2523](https://github.com/sebitt27/dcrd/pull/2523))
- apbf: Add comprehensive tests ([sebitt27/dcrd#2579](https://github.com/sebitt27/dcrd/pull/2579))
- apbf: Add benchmarks ([sebitt27/dcrd#2579](https://github.com/sebitt27/dcrd/pull/2579))
- rpcserver: Add handleGetRawMempool test ([sebitt27/dcrd#2589](https://github.com/sebitt27/dcrd/pull/2589))
- build: Test against go 1.16 ([sebitt27/dcrd#2598](https://github.com/sebitt27/dcrd/pull/2598))
- blockchain: Add test name to TestUtxoEntry errors ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Add UtxoCache test coverage ([sebitt27/dcrd#2591](https://github.com/sebitt27/dcrd/pull/2591))
- blockchain: Use new style for chainio test errors ([sebitt27/dcrd#2595](https://github.com/sebitt27/dcrd/pull/2595))
- rpcserver: Add handleInvalidateBlock test ([sebitt27/dcrd#2604](https://github.com/sebitt27/dcrd/pull/2604))
- blockchain: Mock time.Now for utxo cache tests ([sebitt27/dcrd#2605](https://github.com/sebitt27/dcrd/pull/2605))
- blockchain: Add UtxoCache Initialize tests ([sebitt27/dcrd#2599](https://github.com/sebitt27/dcrd/pull/2599))
- blockchain: Add TestShutdownUtxoCache tests ([sebitt27/dcrd#2599](https://github.com/sebitt27/dcrd/pull/2599))
- rpcserver: Add handleReconsiderBlock test ([sebitt27/dcrd#2613](https://github.com/sebitt27/dcrd/pull/2613))
- stdaddr: Add benchmarks ([sebitt27/dcrd#2610](https://github.com/sebitt27/dcrd/pull/2610))
- rpctest: Make tests work properly with latest code ([sebitt27/dcrd#2614](https://github.com/sebitt27/dcrd/pull/2614))
- mempool: Remove unused field from test struct ([sebitt27/dcrd#2618](https://github.com/sebitt27/dcrd/pull/2618))
- mempool: Remove unused func from tests ([sebitt27/dcrd#2621](https://github.com/sebitt27/dcrd/pull/2621))
- rpctest: Don't use Fatalf in goroutines ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- chaingen: Remove unused PurchaseCommitmentScript ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- rpctest: Convert to use new stdaddr package ([sebitt27/dcrd#2625](https://github.com/sebitt27/dcrd/pull/2625))
- dcrutil: Move address params iface and mock impls ([sebitt27/dcrd#2628](https://github.com/sebitt27/dcrd/pull/2628))
- stdscript: Add v0 p2pk-ecdsa-secp256k1 benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pk-ed25519 benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pk-schnorr-secp256k1 benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pkh-ecdsa-secp256k1 benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pkh-ed25519 benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2pkh-schnorr-secp256k1 benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 p2sh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 ecdsa multisig benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 ecdsa multisig redeem benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add extract v0 multisig redeem benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 nulldata benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake sub p2pkh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake sub p2sh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake gen p2pkh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake gen p2sh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake revoke p2pkh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake revoke p2sh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake change p2pkh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 stake change p2sh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 treasury add benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 treasury gen p2pkh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 treasury gen p2sh benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add determine script type benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- stdscript: Add v0 atomic swap redeem benchmark ([sebitt27/dcrd#2656](https://github.com/sebitt27/dcrd/pull/2656))
- txscript: Separate short form script parsing ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Explicit consensus p2sh tests ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- txscript: Explicit consensus any kind p2sh tests ([sebitt27/dcrd#2666](https://github.com/sebitt27/dcrd/pull/2666))
- stake: No standardness code in tests ([sebitt27/dcrd#2667](https://github.com/sebitt27/dcrd/pull/2667))
- blockchain: Add outpointKey tests ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- blockchain: Add block index key collision tests ([sebitt27/dcrd#2685](https://github.com/sebitt27/dcrd/pull/2685))
- secp256k1: Rework NAF tests ([sebitt27/dcrd#2695](https://github.com/sebitt27/dcrd/pull/2695))
- secp256k1: Cleanup NAF benchmark ([sebitt27/dcrd#2695](https://github.com/sebitt27/dcrd/pull/2695))
- rpctest: Add P2PAddress() function ([sebitt27/dcrd#2704](https://github.com/sebitt27/dcrd/pull/2704))
- tests: Remove hardcoded CC=gcc from run_tests.sh ([sebitt27/dcrd#2706](https://github.com/sebitt27/dcrd/pull/2706))
- build: Test against Go 1.17 ([sebitt27/dcrd#2712](https://github.com/sebitt27/dcrd/pull/2712))
- blockchain: Support voting multiple agendas in test ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- blockchain: Single out treasury policy test ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- blockchain: Correct test harness err msg ([sebitt27/dcrd#2714](https://github.com/sebitt27/dcrd/pull/2714))
- blockchain: Test new max expenditure policy ([sebitt27/dcrd#2680](https://github.com/sebitt27/dcrd/pull/2680))
- chaingen: Add spendable coinbase out snapshots ([sebitt27/dcrd#2715](https://github.com/sebitt27/dcrd/pull/2715))
- mempool: Accept test mungers for create tickets ([sebitt27/dcrd#2721](https://github.com/sebitt27/dcrd/pull/2721))
- build: Don't set GO111MODULE unnecessarily ([sebitt27/dcrd#2722](https://github.com/sebitt27/dcrd/pull/2722))
- build: Don't manually test changing go.{mod,sum} ([sebitt27/dcrd#2722](https://github.com/sebitt27/dcrd/pull/2722))
- stake: Add CalculateRewards tests ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- stake: Add CheckSSRtx tests ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- blockchain: Test auto revocations deployment ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- chaingen: Add revocation mungers ([sebitt27/dcrd#2718](https://github.com/sebitt27/dcrd/pull/2718))
- addrmgr: Improve test coverage ([sebitt27/dcrd#2596](https://github.com/sebitt27/dcrd/pull/2596))
- addrmgr: Remove unnecessary test cases ([sebitt27/dcrd#2596](https://github.com/sebitt27/dcrd/pull/2596))
- rpcserver: Tune large tspend test amount ([sebitt27/dcrd#2679](https://github.com/sebitt27/dcrd/pull/2679))
- build: Pin GitHub Actions to SHA ([sebitt27/dcrd#2736](https://github.com/sebitt27/dcrd/pull/2736))
- blockchain: Add calcTicketReturnAmounts tests ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- blockchain: Add checkTicketRedeemers tests ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- blockchain: Add auto revocation validation tests ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- mining: Add auto revocation block template tests ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- mempool: Add tests with auto revocations enabled ([sebitt27/dcrd#2720](https://github.com/sebitt27/dcrd/pull/2720))
- txscript: Add versioned short form parsing ([sebitt27/dcrd#2756](https://github.com/sebitt27/dcrd/pull/2756))
- txscript: Test consistency and cleanup ([sebitt27/dcrd#2757](https://github.com/sebitt27/dcrd/pull/2757))
- mempool: Add blockHeight to AddFakeUTXO for tests ([sebitt27/dcrd#2804](https://github.com/sebitt27/dcrd/pull/2804))
- mempool: Test fraud proof handling ([sebitt27/dcrd#2804](https://github.com/sebitt27/dcrd/pull/2804))
- stdscript: Add extract v0 stake-tagged p2pkh bench ([sebitt27/dcrd#2816](https://github.com/sebitt27/dcrd/pull/2816))
- stdscript: Add extract v0 stake-tagged p2sh bench ([sebitt27/dcrd#2816](https://github.com/sebitt27/dcrd/pull/2816))
- mempool: Update test to check hash value ([sebitt27/dcrd#2804](https://github.com/sebitt27/dcrd/pull/2804))
- stdscript: Add num required sigs benchmark ([sebitt27/dcrd#2805](https://github.com/sebitt27/dcrd/pull/2805))
- uint256: Add big endian set benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add little endian set benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add big endian get benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add little endian get benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add zero benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add equality comparison benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add less than comparison benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add greater than comparison benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add general comparison benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add addition benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add subtraction benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add multiplication benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add squaring benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add division benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add negation benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add is odd benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise left shift benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise right shift benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise not benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise or benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise and benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bitwise xor benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add bit length benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add text formatting benchmarks ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add conversion to stdlib big int benchmark ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- uint256: Add conversion from stdlib big int benchmark ([sebitt27/dcrd#2787](https://github.com/sebitt27/dcrd/pull/2787))
- primitives: Add diff bits conversion benchmarks ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add work calc benchmark ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add hash to uint256 benchmark ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add check proof of work benchmark ([sebitt27/dcrd#2788](https://github.com/sebitt27/dcrd/pull/2788))
- primitives: Add merkle root benchmarks ([sebitt27/dcrd#2826](https://github.com/sebitt27/dcrd/pull/2826))
- primitives: Add inclusion proof benchmarks ([sebitt27/dcrd#2827](https://github.com/sebitt27/dcrd/pull/2827))
- blockchain: Add AssumeValid tests ([sebitt27/dcrd#2839](https://github.com/sebitt27/dcrd/pull/2839))
- chaingen: Add vote subsidy munger ([sebitt27/dcrd#2848](https://github.com/sebitt27/dcrd/pull/2848))

### Misc:

- release: Bump for 1.7 release cycle ([sebitt27/dcrd#2429](https://github.com/sebitt27/dcrd/pull/2429))
- secp256k1: Correct const name for doc comment ([sebitt27/dcrd#2445](https://github.com/sebitt27/dcrd/pull/2445))
- multi: Fix various typos ([sebitt27/dcrd#2607](https://github.com/sebitt27/dcrd/pull/2607))
- rpcserver: Fix createrawssrtx comments ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- blockchain: Fix comment formatting in generator ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- stake: Fix MaxOutputsPerSSRtx comment ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- stake: Fix CheckSSGenVotes function comment ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- stake: Fix CheckSSRtx function comment ([sebitt27/dcrd#2665](https://github.com/sebitt27/dcrd/pull/2665))
- database: Add comment on os.MkdirAll behavior ([sebitt27/dcrd#2670](https://github.com/sebitt27/dcrd/pull/2670))
- multi: Address some linter complaints ([sebitt27/dcrd#2684](https://github.com/sebitt27/dcrd/pull/2684))
- txscript: Fix a couple of a comment typos ([sebitt27/dcrd#2692](https://github.com/sebitt27/dcrd/pull/2692))
- blockchain: Remove inapplicable comment ([sebitt27/dcrd#2742](https://github.com/sebitt27/dcrd/pull/2742))
- mining: Fix error in comment ([sebitt27/dcrd#2743](https://github.com/sebitt27/dcrd/pull/2743))
- blockchain: Fix several typos ([sebitt27/dcrd#2745](https://github.com/sebitt27/dcrd/pull/2745))
- blockchain: Update a few BFFastAdd comments ([sebitt27/dcrd#2781](https://github.com/sebitt27/dcrd/pull/2781))
- multi: Address some linter complaints ([sebitt27/dcrd#2791](https://github.com/sebitt27/dcrd/pull/2791))
- netsync: Correct typo ([sebitt27/dcrd#2813](https://github.com/sebitt27/dcrd/pull/2813))
- netsync: Fix misc typos ([sebitt27/dcrd#2834](https://github.com/sebitt27/dcrd/pull/2834))
- mining: Fix typo ([sebitt27/dcrd#2834](https://github.com/sebitt27/dcrd/pull/2834))
- blockchain: Correct comment typos for find fork ([sebitt27/dcrd#2828](https://github.com/sebitt27/dcrd/pull/2828))
- rpcserver: Rename var to make linter happy ([sebitt27/dcrd#2835](https://github.com/sebitt27/dcrd/pull/2835))
- blockchain: Wrap at max line length ([sebitt27/dcrd#2843](https://github.com/sebitt27/dcrd/pull/2843))
- release: Bump for 1.7.0 ([sebitt27/dcrd#2856](https://github.com/sebitt27/dcrd/pull/2856))

### Code Contributors (alphabetical order):

- briancolecoinmetrics
- Dave Collins
- David Hill
- degeri
- Donald Adu-Poku
- J Fixby
- Jamie Holdstock
- Joe Gruffins
- Jonathan Chappelow
- Josh Rickmar
- lolandhold
- Matheus Degiovani
- Naveen
- Ryan Staudt
- Youssef Boukenken
- Wisdom Arerosuoghene
