## dcrd v1.0.7

This release of dcrd primarily contains improvements to the infrastructure and
other quality assurance changes that are bringing us closer to providing full
support for Lightning Network.

A lot of work required for Lightning Network support went into getting the
required code merged into the upstream project, btcd, which now fully supports
it.  These changes also must be synced and integrated with dcrd as well and
therefore many of the changes in this release are related to that process.

## Notable Changes

### Dust check removed from stake transactions

The standard policy for regular transactions is to reject any transactions that
have outputs so small that they cost more to the network than their value.  This
behavior is desirable for regular transactions, however it was also being
applied to vote and revocation transactions which could lead to a situation
where stake pools with low fees could result in votes and revocations having
difficulty being mined.

This check has been changed to only apply to regular transactions now in order
to prevent any issues.  Stake transactions have several other checks that make
this one unnecessary for them.

### New `feefilter` peer-to-peer message

A new optional peer-to-peer message named `feefilter` has been added that allows
peers to inform others about the minimum transaction fee rate they are willing
to accept.  This will enable peers to avoid notifying others about transactions
they will not accept anyways and therefore can result in a significant bandwidth
savings.

### Bloom filter service bit enforcement

Peers that are configured to disable bloom filter support will now disconnect
remote peers that send bloom filter related commands rather than simply ignoring
them.  This allows any light clients that do not observe the service bit to
potentially find another peer that provides the service.  Additionally, remote
peers that have negotiated a high enough protocol version to observe the service
bit and still send bloom filter related commands anyways will now be banned.


## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/sebitt27/dcrd/compare/v1.0.5...v1.0.7).

### Protocol and network:
- Allow reorg of block one [sebitt27/dcrd#745](https://github.com/sebitt27/dcrd/pull/745)
- blockchain: use the time source [sebitt27/dcrd#747](https://github.com/sebitt27/dcrd/pull/747)
- peer: Strictly enforce bloom filter service bit [sebitt27/dcrd#768](https://github.com/sebitt27/dcrd/pull/768)
- wire/peer: Implement feefilter p2p message [sebitt27/dcrd#779](https://github.com/sebitt27/dcrd/pull/779)
- chaincfg: update checkpoints for 1.0.7 release  [sebitt27/dcrd#816](https://github.com/sebitt27/dcrd/pull/816)

### Transaction relay (memory pool):
- mempool: Break dependency on chain instance [sebitt27/dcrd#754](https://github.com/sebitt27/dcrd/pull/754)
- mempool: unexport the mutex [sebitt27/dcrd#755](https://github.com/sebitt27/dcrd/pull/755)
- mempool: Add basic test harness infrastructure [sebitt27/dcrd#756](https://github.com/sebitt27/dcrd/pull/756)
- mempool: Improve tx input standard checks [sebitt27/dcrd#758](https://github.com/sebitt27/dcrd/pull/758)
- mempool: Update comments for dust calcs [sebitt27/dcrd#764](https://github.com/sebitt27/dcrd/pull/764)
- mempool: Only perform standard dust checks on regular transactions  [sebitt27/dcrd#806](https://github.com/sebitt27/dcrd/pull/806)

### RPC:
- Fix gettxout includemempool handling [sebitt27/dcrd#738](https://github.com/sebitt27/dcrd/pull/738)
- Improve help text for getmininginfo [sebitt27/dcrd#748](https://github.com/sebitt27/dcrd/pull/748)
- rpcserverhelp: update TicketFeeInfo help [sebitt27/dcrd#801](https://github.com/sebitt27/dcrd/pull/801)
- blockchain: Improve getstakeversions efficiency [sebitt27/dcrd#81](https://github.com/sebitt27/dcrd/pull/813)

### dcrd command-line flags:
- config: introduce new flags to accept/reject non-std transactions [sebitt27/dcrd#757](https://github.com/sebitt27/dcrd/pull/757)
- config: Add --whitelist option [sebitt27/dcrd#352](https://github.com/sebitt27/dcrd/pull/352)
- config: Improve config file handling [sebitt27/dcrd#802](https://github.com/sebitt27/dcrd/pull/802)
- config: Improve blockmaxsize check [sebitt27/dcrd#810](https://github.com/sebitt27/dcrd/pull/810)

### dcrctl:
- Add --walletrpcserver option [sebitt27/dcrd#736](https://github.com/sebitt27/dcrd/pull/736)

### Documentation
- docs: add commit prefix notes  [sebitt27/dcrd#788](https://github.com/sebitt27/dcrd/pull/788)

### Developer-related package changes:
- blockchain: check errors and remove ineffectual assignments [sebitt27/dcrd#689](https://github.com/sebitt27/dcrd/pull/689)
- stake: less casting [sebitt27/dcrd#705](https://github.com/sebitt27/dcrd/pull/705)
- blockchain: chainstate only needs prev block hash [sebitt27/dcrd#706](https://github.com/sebitt27/dcrd/pull/706)
- remove dead code [sebitt27/dcrd#715](https://github.com/sebitt27/dcrd/pull/715)
- Use btclog for determining valid log levels [sebitt27/dcrd#738](https://github.com/sebitt27/dcrd/pull/738)
- indexers: Minimize differences with upstream code [sebitt27/dcrd#742](https://github.com/sebitt27/dcrd/pull/742)
- blockchain: Add median time to state snapshot [sebitt27/dcrd#753](https://github.com/sebitt27/dcrd/pull/753)
- blockmanager: remove unused GetBlockFromHash function [sebitt27/dcrd#761](https://github.com/sebitt27/dcrd/pull/761)
- mining: call CheckConnectBlock directly [sebitt27/dcrd#762](https://github.com/sebitt27/dcrd/pull/762)
- blockchain: add missing error code entries [sebitt27/dcrd#763](https://github.com/sebitt27/dcrd/pull/763)
- blockchain: Sync main chain flag on ProcessBlock [sebitt27/dcrd#767](https://github.com/sebitt27/dcrd/pull/767)
- blockchain: Remove exported CalcPastTimeMedian func [sebitt27/dcrd#770](https://github.com/sebitt27/dcrd/pull/770)
- blockchain: check for error [sebitt27/dcrd#772](https://github.com/sebitt27/dcrd/pull/772)
- multi: Optimize by removing defers [sebitt27/dcrd#782](https://github.com/sebitt27/dcrd/pull/782)
- blockmanager: remove unused logBlockHeight [sebitt27/dcrd#787](https://github.com/sebitt27/dcrd/pull/787)
- dcrutil: Replace DecodeNetworkAddress with DecodeAddress [sebitt27/dcrd#746](https://github.com/sebitt27/dcrd/pull/746)
- txscript: Force extracted addrs to compressed [sebitt27/dcrd#775](https://github.com/sebitt27/dcrd/pull/775)
- wire: Remove legacy transaction decoding [sebitt27/dcrd#794](https://github.com/sebitt27/dcrd/pull/794)
- wire: Remove dead legacy tx decoding code [sebitt27/dcrd#796](https://github.com/sebitt27/dcrd/pull/796)
- mempool/wire: Don't make policy decisions in wire [sebitt27/dcrd#797](https://github.com/sebitt27/dcrd/pull/797)
- dcrjson: Remove unused cmds & types [sebitt27/dcrd#795](https://github.com/sebitt27/dcrd/pull/795)
- dcrjson: move cmd types [sebitt27/dcrd#799](https://github.com/sebitt27/dcrd/pull/799)
- multi: Separate tx serialization type from version [sebitt27/dcrd#798](https://github.com/sebitt27/dcrd/pull/798)
- dcrjson: add Unconfirmed field to dcrjson.GetAccountBalanceResult [sebitt27/dcrd#812](https://github.com/sebitt27/dcrd/pull/812)
- multi: Error descriptions should be lowercase [sebitt27/dcrd#771](https://github.com/sebitt27/dcrd/pull/771)
- blockchain: cast to int64  [sebitt27/dcrd#817](https://github.com/sebitt27/dcrd/pull/817)

### Testing and Quality Assurance:
- rpcserver: Upstream sync to add basic RPC tests [sebitt27/dcrd#750](https://github.com/sebitt27/dcrd/pull/750)
- rpctest: Correct several issues tests and joins [sebitt27/dcrd#751](https://github.com/sebitt27/dcrd/pull/751)
- rpctest: prevent process leak due to panics [sebitt27/dcrd#752](https://github.com/sebitt27/dcrd/pull/752)
- rpctest: Cleanup resources on failed setup [sebitt27/dcrd#759](https://github.com/sebitt27/dcrd/pull/759)
- rpctest: Use ports based on the process id [sebitt27/dcrd#760](https://github.com/sebitt27/dcrd/pull/760)
- rpctest/deps: Update dependencies and API [sebitt27/dcrd#765](https://github.com/sebitt27/dcrd/pull/765)
- rpctest: Gate rpctest-based behind a build tag [sebitt27/dcrd#766](https://github.com/sebitt27/dcrd/pull/766)
- mempool: Add test for max orphan entry eviction [sebitt27/dcrd#769](https://github.com/sebitt27/dcrd/pull/769)
- fullblocktests: Add more consensus tests [sebitt27/dcrd#77](https://github.com/sebitt27/dcrd/pull/773)
- fullblocktests: Sync upstream block validation [sebitt27/dcrd#774](https://github.com/sebitt27/dcrd/pull/774)
- rpctest: fix a harness range bug in syncMempools [sebitt27/dcrd#778](https://github.com/sebitt27/dcrd/pull/778)
- secp256k1: Add regression tests for field.go [sebitt27/dcrd#781](https://github.com/sebitt27/dcrd/pull/781)
- secp256k1: Sync upstream test consolidation [sebitt27/dcrd#783](https://github.com/sebitt27/dcrd/pull/783)
- txscript: Correct p2sh hashes in json test data  [sebitt27/dcrd#785](https://github.com/sebitt27/dcrd/pull/785)
- txscript: Replace CODESEPARATOR json test data [sebitt27/dcrd#786](https://github.com/sebitt27/dcrd/pull/786)
- txscript: Remove multisigdummy from json test data [sebitt27/dcrd#789](https://github.com/sebitt27/dcrd/pull/789)
- txscript: Remove max money from json test data [sebitt27/dcrd#790](https://github.com/sebitt27/dcrd/pull/790)
- txscript: Update signatures in json test data [sebitt27/dcrd#791](https://github.com/sebitt27/dcrd/pull/791)
- txscript: Use native encoding in json test data [sebitt27/dcrd#792](https://github.com/sebitt27/dcrd/pull/792)
- rpctest: Store logs and data in same path [sebitt27/dcrd#780](https://github.com/sebitt27/dcrd/pull/780)
- txscript: Cleanup reference test code  [sebitt27/dcrd#793](https://github.com/sebitt27/dcrd/pull/793)

### Misc:
- Update deps to pull in additional logging changes [sebitt27/dcrd#734](https://github.com/sebitt27/dcrd/pull/734)
- Update markdown files for GFM changes [sebitt27/dcrd#744](https://github.com/sebitt27/dcrd/pull/744)
- blocklogger: Show votes, tickets, & revocations [sebitt27/dcrd#784](https://github.com/sebitt27/dcrd/pull/784)
- blocklogger: Remove STransactions from transactions calculation [sebitt27/dcrd#811](https://github.com/sebitt27/dcrd/pull/811)

### Contributors (alphabetical order):

- Alex Yocomm-Piatt
- Atri Viss
- Chris Martin
- Dave Collins
- David Hill
- Donald Adu-Poku
- Jimmy Song
- John C. Vernaleo
- Jolan Luff
- Josh Rickmar
- Olaoluwa Osuntokun
- Marco Peereboom
