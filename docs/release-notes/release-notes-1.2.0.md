# dcrd v1.2.0

This release of dcrd contains significant performance enhancements,
infrastructure improvements, improved access to chain-related information for
providing better SPV (Simplified Payment Verification) support, and other
quality assurance changes.

A significant amount of infrastructure work has also been done this release
cycle towards being able to support several planned scalability optimizations.

## Downgrade Warning

The database format in v1.2.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Significantly Faster Startup

The startup time has been improved by roughly 17x on slower hard disk drives
(HDDs) and 8x on solid state drives (SSDs).

In order to achieve these speedups, there is a one time database migration, as
previously mentioned, that will likely take a while to complete (typically
around 5 to 6 minutes on HDDs and 2 to 3 minutes on SSDs).

### Support For DNS Seed Filtering

In order to better support the forthcoming SPV wallets, support for finding
other peers based upon their enabled services has been added.  This is useful
for both SPV wallets and full nodes since SPV wallets will require access to
full nodes in order to retrieve the necessary proofs and full nodes are
generally not interested in making outgoing connections to SPV wallets.

### Committed Filters

With the intention of supporting light clients, such as SPV wallets, in a
privacy-preserving way while still minimizing the amount of data that needs to
be downloaded, this release adds support for committed filters.  A committed
filter is a combination of a probalistic data structure that is used to test
whether an element is a member of a set with a predetermined collision
probability along with a commitment by consensus-validating full nodes to that
data.

A committed filter is created for every block which allows light clients to
download the filters and match against them locally rather than uploading
personal data to other nodes.

A new service flag is also provided to allow clients to discover nodes that
provide access to filters.

There is a one time database update to build and store the filters for all
existing historical blocks which will likely take a while to complete (typically
around 2 to 3 minutes on HDDs and 1 to 1.5 minutes on SSDs).

### Updated Atomic Swap Contracts

The standard checks for atomic swap contracts have been updated to ensure the
contracts enforce the secret size for safer support between chains with
disparate script rules.

### RPC Server Changes

#### New `getchaintips` RPC

A new RPC named `getchaintips` has been added which allows callers to query
information about the status of known side chains and their branch lengths.
It currently only provides support for side chains that have been seen while the
current instance of the process is running.  This will be further improved in
future releases.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/sebitt27/dcrd/compare/v1.1.2...v1.2.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.2.0 release ([sebitt27/dcrd#1139](https://github.com/sebitt27/dcrd/pull/1139))
- chaincfg: Introduce new type DNSSeed ([sebitt27/dcrd#961](https://github.com/sebitt27/dcrd/pull/961))
- blockmanager: sync with the most updated peer ([sebitt27/dcrd#984](https://github.com/sebitt27/dcrd/pull/984))
- multi: remove MsgAlert ([sebitt27/dcrd#1161](https://github.com/sebitt27/dcrd/pull/1161))
- multi: Add initial committed filter (CF) support ([sebitt27/dcrd#1151](https://github.com/sebitt27/dcrd/pull/1151))

### Transaction relay (memory pool):

- txscript: Correct nulldata standardness check ([sebitt27/dcrd#935](https://github.com/sebitt27/dcrd/pull/935))
- mempool: Optimize orphan map limiting ([sebitt27/dcrd#1117](https://github.com/sebitt27/dcrd/pull/1117))
- mining: Fix duplicate txns in the prio heap ([sebitt27/dcrd#1108](https://github.com/sebitt27/dcrd/pull/1108))
- mining: Stop transactions losing their dependants ([sebitt27/dcrd#1109](https://github.com/sebitt27/dcrd/pull/1109))

### RPC:

- rpcserver: skip cert create when RPC is disabled ([sebitt27/dcrd#949](https://github.com/sebitt27/dcrd/pull/949))
- rpcserver: remove redundant checks in blockTemplateResult ([sebitt27/dcrd#826](https://github.com/sebitt27/dcrd/pull/826))
- rpcserver: assert network for validateaddress rpc ([sebitt27/dcrd#963](https://github.com/sebitt27/dcrd/pull/963))
- rpcserver: Do not rebroadcast stake transactions ([sebitt27/dcrd#973](https://github.com/sebitt27/dcrd/pull/973))
- dcrjson: add ticket fee field to PurchaseTicketCmd ([sebitt27/dcrd#902](https://github.com/sebitt27/dcrd/pull/902))
- dcrwalletextcmds: remove getseed ([sebitt27/dcrd#985](https://github.com/sebitt27/dcrd/pull/985))
- dcrjson: Add SweepAccountCmd & SweepAccountResult ([sebitt27/dcrd#1027](https://github.com/sebitt27/dcrd/pull/1027))
- rpcserver: add sweepaccount to the wallet list of commands ([sebitt27/dcrd#1028](https://github.com/sebitt27/dcrd/pull/1028))
- rpcserver: add batched request support (json 2.0) ([sebitt27/dcrd#841](https://github.com/sebitt27/dcrd/pull/841))
- dcrjson: include summary totals in GetBalanceResult ([sebitt27/dcrd#1062](https://github.com/sebitt27/dcrd/pull/1062))
- multi: Implement getchaintips JSON-RPC ([sebitt27/dcrd#1098](https://github.com/sebitt27/dcrd/pull/1098))
- rpcserver: Add dcrd version info to getversion RPC ([sebitt27/dcrd#1097](https://github.com/sebitt27/dcrd/pull/1097))
- rpcserver: Correct getblockheader result text ([sebitt27/dcrd#1104](https://github.com/sebitt27/dcrd/pull/1104))
- dcrjson: add StartAutoBuyerCmd & StopAutoBuyerCmd ([sebitt27/dcrd#903](https://github.com/sebitt27/dcrd/pull/903))
- dcrjson: fix typo for StartAutoBuyerCmd ([sebitt27/dcrd#1146](https://github.com/sebitt27/dcrd/pull/1146))
- dcrjson: require passphrase for StartAutoBuyerCmd ([sebitt27/dcrd#1147](https://github.com/sebitt27/dcrd/pull/1147))
- dcrjson: fix StopAutoBuyerCmd registration bug ([sebitt27/dcrd#1148](https://github.com/sebitt27/dcrd/pull/1148))
- blockchain: Support testnet stake diff estimation ([sebitt27/dcrd#1115](https://github.com/sebitt27/dcrd/pull/1115))
- rpcserver: fix jsonRPCRead data race ([sebitt27/dcrd#1157](https://github.com/sebitt27/dcrd/pull/1157))
- dcrjson: Add VerifySeedCmd ([sebitt27/dcrd#1160](https://github.com/sebitt27/dcrd/pull/1160))

### dcrd command-line flags and configuration:

- mempool: Rename RelayNonStd config option ([sebitt27/dcrd#1024](https://github.com/sebitt27/dcrd/pull/1024))
- sampleconfig: Update min relay fee ([sebitt27/dcrd#959](https://github.com/sebitt27/dcrd/pull/959))
- sampleconfig: Correct comment ([sebitt27/dcrd#1063](https://github.com/sebitt27/dcrd/pull/1063))
- multi: Expand ~ to correct home directory on all OSes ([sebitt27/dcrd#1041](https://github.com/sebitt27/dcrd/pull/1041))

### checkdevpremine utility changes:

- checkdevpremine: Remove --skipverify option ([sebitt27/dcrd#969](https://github.com/sebitt27/dcrd/pull/969))
- checkdevpremine: Implement --notls option ([sebitt27/dcrd#969](https://github.com/sebitt27/dcrd/pull/969))
- checkdevpremine: Make file naming consistent ([sebitt27/dcrd#969](https://github.com/sebitt27/dcrd/pull/969))
- checkdevpremine: Fix comment ([sebitt27/dcrd#969](https://github.com/sebitt27/dcrd/pull/969))
- checkdevpremine: Remove utility ([sebitt27/dcrd#1068](https://github.com/sebitt27/dcrd/pull/1068))

### Documentation:

- fullblocktests: Add missing doc.go file ([sebitt27/dcrd#956](https://github.com/sebitt27/dcrd/pull/956))
- docs: Add fullblocktests entry and make consistent ([sebitt27/dcrd#956](https://github.com/sebitt27/dcrd/pull/956))
- docs: Add mempool entry to developer tools section ([sebitt27/dcrd#1058](https://github.com/sebitt27/dcrd/pull/1058))
- mempool: Add docs.go and flesh out README.md ([sebitt27/dcrd#1058](https://github.com/sebitt27/dcrd/pull/1058))
- docs: document packages and fix typo  ([sebitt27/dcrd#965](https://github.com/sebitt27/dcrd/pull/965))
- docs: rpcclient is now part of the main dcrd repo ([sebitt27/dcrd#970](https://github.com/sebitt27/dcrd/pull/970))
- dcrjson: Update README.md ([sebitt27/dcrd#982](https://github.com/sebitt27/dcrd/pull/982))
- docs: Remove carriage return ([sebitt27/dcrd#1106](https://github.com/sebitt27/dcrd/pull/1106))
- Adjust README.md for new Go version ([sebitt27/dcrd#1105](https://github.com/sebitt27/dcrd/pull/1105))
- docs: document how to use go test -coverprofile ([sebitt27/dcrd#1107](https://github.com/sebitt27/dcrd/pull/1107))
- addrmgr: Improve documentation ([sebitt27/dcrd#1125](https://github.com/sebitt27/dcrd/pull/1125))
- docs: Fix links for internal packages ([sebitt27/dcrd#1144](https://github.com/sebitt27/dcrd/pull/1144))

### Developer-related package changes:

- chaingen: Add revocation generation infrastructure ([sebitt27/dcrd#1120](https://github.com/sebitt27/dcrd/pull/1120))
- txscript: Add null data script creator ([sebitt27/dcrd#943](https://github.com/sebitt27/dcrd/pull/943))
- txscript: Cleanup and improve NullDataScript tests ([sebitt27/dcrd#943](https://github.com/sebitt27/dcrd/pull/943))
- txscript: Allow external signature hash calc ([sebitt27/dcrd#951](https://github.com/sebitt27/dcrd/pull/951))
- secp256k1: update func signatures ([sebitt27/dcrd#934](https://github.com/sebitt27/dcrd/pull/934))
- txscript: enforce MaxDataCarrierSize for GenerateProvablyPruneableOut ([sebitt27/dcrd#953](https://github.com/sebitt27/dcrd/pull/953))
- txscript: Remove OP_SMALLDATA ([sebitt27/dcrd#954](https://github.com/sebitt27/dcrd/pull/954))
- blockchain: Accept header in CheckProofOfWork ([sebitt27/dcrd#977](https://github.com/sebitt27/dcrd/pull/977))
- blockchain: Make func definition style consistent ([sebitt27/dcrd#983](https://github.com/sebitt27/dcrd/pull/983))
- blockchain: only fetch the parent block in BFFastAdd ([sebitt27/dcrd#972](https://github.com/sebitt27/dcrd/pull/972))
- blockchain: Switch to FindSpentTicketsInBlock ([sebitt27/dcrd#915](https://github.com/sebitt27/dcrd/pull/915))
- stake: Add Hash256PRNG init vector support ([sebitt27/dcrd#986](https://github.com/sebitt27/dcrd/pull/986))
- blockchain/stake: Use Hash256PRNG init vector ([sebitt27/dcrd#987](https://github.com/sebitt27/dcrd/pull/987))
- blockchain: Don't store full header in block node ([sebitt27/dcrd#988](https://github.com/sebitt27/dcrd/pull/988))
- blockchain: Reconstruct headers from block nodes ([sebitt27/dcrd#989](https://github.com/sebitt27/dcrd/pull/989))
- stake/multi: Don't return errors for IsX functions ([sebitt27/dcrd#995](https://github.com/sebitt27/dcrd/pull/995))
- blockchain: Rename block index to main chain index ([sebitt27/dcrd#996](https://github.com/sebitt27/dcrd/pull/996))
- blockchain: Refactor main block index logic ([sebitt27/dcrd#990](https://github.com/sebitt27/dcrd/pull/990))
- blockchain: Use hash values in structs ([sebitt27/dcrd#992](https://github.com/sebitt27/dcrd/pull/992))
- blockchain: Remove unused dump function ([sebitt27/dcrd#1001](https://github.com/sebitt27/dcrd/pull/1001))
- blockchain: Generalize and optimize chain reorg ([sebitt27/dcrd#997](https://github.com/sebitt27/dcrd/pull/997))
- blockchain: Pass parent block in connection code ([sebitt27/dcrd#998](https://github.com/sebitt27/dcrd/pull/998))
- blockchain: Explicit block fetch semanticss ([sebitt27/dcrd#999](https://github.com/sebitt27/dcrd/pull/999))
- blockchain: Use next detach block in reorg chain ([sebitt27/dcrd#1002](https://github.com/sebitt27/dcrd/pull/1002))
- blockchain: Limit header sanity check to header ([sebitt27/dcrd#1003](https://github.com/sebitt27/dcrd/pull/1003))
- blockchain: Validate num votes in header sanity ([sebitt27/dcrd#1005](https://github.com/sebitt27/dcrd/pull/1005))
- blockchain: Validate max votes in header sanity ([sebitt27/dcrd#1006](https://github.com/sebitt27/dcrd/pull/1006))
- blockchain: Validate stake diff in header context ([sebitt27/dcrd#1004](https://github.com/sebitt27/dcrd/pull/1004))
- blockchain: No votes/revocations in header sanity ([sebitt27/dcrd#1007](https://github.com/sebitt27/dcrd/pull/1007))
- blockchain: Validate max purchases in header sanity ([sebitt27/dcrd#1008](https://github.com/sebitt27/dcrd/pull/1008))
- blockchain: Validate early votebits in header sanity ([sebitt27/dcrd#1009](https://github.com/sebitt27/dcrd/pull/1009))
- blockchain: Validate block height in header context ([sebitt27/dcrd#1010](https://github.com/sebitt27/dcrd/pull/1010))
- blockchain: Move check block context func ([sebitt27/dcrd#1011](https://github.com/sebitt27/dcrd/pull/1011))
- blockchain: Block sanity cleanup and consistency ([sebitt27/dcrd#1012](https://github.com/sebitt27/dcrd/pull/1012))
- blockchain: Remove dup ticket purchase value check ([sebitt27/dcrd#1013](https://github.com/sebitt27/dcrd/pull/1013))
- blockchain: Only tickets before SVH in block sanity ([sebitt27/dcrd#1014](https://github.com/sebitt27/dcrd/pull/1014))
- blockchain: Remove unused vote bits function ([sebitt27/dcrd#1015](https://github.com/sebitt27/dcrd/pull/1015))
- blockchain: Move upgrade-only code to upgrade.go ([sebitt27/dcrd#1016](https://github.com/sebitt27/dcrd/pull/1016))
- stake: Static assert of vote commitment ([sebitt27/dcrd#1020](https://github.com/sebitt27/dcrd/pull/1020))
- blockchain: Remove unused error code ([sebitt27/dcrd#1021](https://github.com/sebitt27/dcrd/pull/1021))
- blockchain: Improve readability of parent approval ([sebitt27/dcrd#1022](https://github.com/sebitt27/dcrd/pull/1022))
- peer: rename mruinvmap, mrunoncemap to lruinvmap, lrunoncemap ([sebitt27/dcrd#976](https://github.com/sebitt27/dcrd/pull/976))
- peer: rename noncemap to noncecache ([sebitt27/dcrd#976](https://github.com/sebitt27/dcrd/pull/976))
- peer: rename inventorymap to inventorycache ([sebitt27/dcrd#976](https://github.com/sebitt27/dcrd/pull/976))
- connmgr: convert state to atomic ([sebitt27/dcrd#1025](https://github.com/sebitt27/dcrd/pull/1025))
- blockchain/mining: Full checks in CCB ([sebitt27/dcrd#1017](https://github.com/sebitt27/dcrd/pull/1017))
- blockchain: Validate pool size in header context ([sebitt27/dcrd#1018](https://github.com/sebitt27/dcrd/pull/1018))
- blockchain: Vote commitments in block sanity ([sebitt27/dcrd#1023](https://github.com/sebitt27/dcrd/pull/1023))
- blockchain: Validate early final state is zero ([sebitt27/dcrd#1031](https://github.com/sebitt27/dcrd/pull/1031))
- blockchain: Validate final state in header context ([sebitt27/dcrd#1034](https://github.com/sebitt27/dcrd/pull/1033))
- blockchain: Max revocations in block sanity ([sebitt27/dcrd#1034](https://github.com/sebitt27/dcrd/pull/1034))
- blockchain: Allowed stake txns in block sanity ([sebitt27/dcrd#1035](https://github.com/sebitt27/dcrd/pull/1035))
- blockchain: Validate allowed votes in block context ([sebitt27/dcrd#1036](https://github.com/sebitt27/dcrd/pull/1036))
- blockchain: Validate allowed revokes in blk contxt ([sebitt27/dcrd#1037](https://github.com/sebitt27/dcrd/pull/1037))
- blockchain/stake: Rename tix spent to tix voted ([sebitt27/dcrd#1038](https://github.com/sebitt27/dcrd/pull/1038))
- txscript: Require atomic swap contracts to specify the secret size ([sebitt27/dcrd#1039](https://github.com/sebitt27/dcrd/pull/1039))
- blockchain: Remove unused struct ([sebitt27/dcrd#1043](https://github.com/sebitt27/dcrd/pull/1043))
- blockchain: Store side chain blocks in database ([sebitt27/dcrd#1000](https://github.com/sebitt27/dcrd/pull/1000))
- blockchain: Simplify initial chain state ([sebitt27/dcrd#1045](https://github.com/sebitt27/dcrd/pull/1045))
- blockchain: Rework database versioning ([sebitt27/dcrd#1047](https://github.com/sebitt27/dcrd/pull/1047))
- blockchain: Don't require chain for db upgrades ([sebitt27/dcrd#1051](https://github.com/sebitt27/dcrd/pull/1051))
- blockchain/indexers: Allow interrupts ([sebitt27/dcrd#1052](https://github.com/sebitt27/dcrd/pull/1052))
- blockchain: Remove old version information ([sebitt27/dcrd#1055](https://github.com/sebitt27/dcrd/pull/1055))
- stake: optimize FindSpentTicketsInBlock slightly ([sebitt27/dcrd#1049](https://github.com/sebitt27/dcrd/pull/1049))
- blockchain: Do not accept orphans/genesis block ([sebitt27/dcrd#1057](https://github.com/sebitt27/dcrd/pull/1057))
- blockchain: Separate node ticket info population ([sebitt27/dcrd#1056](https://github.com/sebitt27/dcrd/pull/1056))
- blockchain: Accept parent in blockNode constructor ([sebitt27/dcrd#1056](https://github.com/sebitt27/dcrd/pull/1056))
- blockchain: Combine ErrDoubleSpend & ErrMissingTx ([sebitt27/dcrd#1064](https://github.com/sebitt27/dcrd/pull/1064))
- blockchain: Calculate the lottery IV on demand ([sebitt27/dcrd#1065](https://github.com/sebitt27/dcrd/pull/1065))
- blockchain: Simplify add/remove node logic ([sebitt27/dcrd#1067](https://github.com/sebitt27/dcrd/pull/1067))
- blockchain: Infrastructure to manage block index ([sebitt27/dcrd#1044](https://github.com/sebitt27/dcrd/pull/1044))
- blockchain: Add block validation status to index ([sebitt27/dcrd#1044](https://github.com/sebitt27/dcrd/pull/1044))
- blockchain: Migrate to new block indexuse it ([sebitt27/dcrd#1044](https://github.com/sebitt27/dcrd/pull/1044))
- blockchain: Lookup child in force head reorg ([sebitt27/dcrd#1070](https://github.com/sebitt27/dcrd/pull/1070))
- blockchain: Refactor block idx entry serialization ([sebitt27/dcrd#1069](https://github.com/sebitt27/dcrd/pull/1069))
- blockchain: Limit GetStakeVersions count ([sebitt27/dcrd#1071](https://github.com/sebitt27/dcrd/pull/1071))
- blockchain: Remove dry run flag ([sebitt27/dcrd#1073](https://github.com/sebitt27/dcrd/pull/1073))
- blockchain: Remove redundant stake ver calc func ([sebitt27/dcrd#1087](https://github.com/sebitt27/dcrd/pull/1087))
- blockchain: Reduce GetGeneration to TipGeneration ([sebitt27/dcrd#1083](https://github.com/sebitt27/dcrd/pull/1083))
- blockchain: Add chain tip tracking ([sebitt27/dcrd#1084](https://github.com/sebitt27/dcrd/pull/1084))
- blockchain: Switch tip generation to chain tips ([sebitt27/dcrd#1085](https://github.com/sebitt27/dcrd/pull/1085))
- blockchain: Simplify voter version calculation ([sebitt27/dcrd#1088](https://github.com/sebitt27/dcrd/pull/1088))
- blockchain: Remove unused threshold serialization ([sebitt27/dcrd#1092](https://github.com/sebitt27/dcrd/pull/1092))
- blockchain: Simplify chain tip tracking ([sebitt27/dcrd#1092](https://github.com/sebitt27/dcrd/pull/1092))
- blockchain: Cache tip and parent at init ([sebitt27/dcrd#1100](https://github.com/sebitt27/dcrd/pull/1100))
- mining: Obtain block by hash instead of top block ([sebitt27/dcrd#1094](https://github.com/sebitt27/dcrd/pull/1094))
- blockchain: Remove unused GetTopBlock function ([sebitt27/dcrd#1094](https://github.com/sebitt27/dcrd/pull/1094))
- multi: Rename BIP0111Version to NodeBloomVersion ([sebitt27/dcrd#1112](https://github.com/sebitt27/dcrd/pull/1112))
- mining/mempool: Move priority code to mining pkg ([sebitt27/dcrd#1110](https://github.com/sebitt27/dcrd/pull/1110))
- mining: Use single uint64 coinbase extra nonce ([sebitt27/dcrd#1116](https://github.com/sebitt27/dcrd/pull/1116))
- mempool/mining: Clarify tree validity semantics ([sebitt27/dcrd#1118](https://github.com/sebitt27/dcrd/pull/1118))
- mempool/mining: TxSource separation ([sebitt27/dcrd#1119](https://github.com/sebitt27/dcrd/pull/1119))
- connmgr: Use same Dial func signature as net.Dial ([sebitt27/dcrd#1113](https://github.com/sebitt27/dcrd/pull/1113))
- addrmgr: Declutter package API ([sebitt27/dcrd#1124](https://github.com/sebitt27/dcrd/pull/1124))
- mining: Correct initial template generation ([sebitt27/dcrd#1122](https://github.com/sebitt27/dcrd/pull/1122))
- cpuminer: Use header for extra nonce ([sebitt27/dcrd#1123](https://github.com/sebitt27/dcrd/pull/1123))
- addrmgr: Make writing of peers file safer ([sebitt27/dcrd#1126](https://github.com/sebitt27/dcrd/pull/1126))
- addrmgr: Save peers file only if necessary ([sebitt27/dcrd#1127](https://github.com/sebitt27/dcrd/pull/1127))
- addrmgr: Factor out common code ([sebitt27/dcrd#1138](https://github.com/sebitt27/dcrd/pull/1138))
- addrmgr: Improve isBad() performance ([sebitt27/dcrd#1134](https://github.com/sebitt27/dcrd/pull/1134))
- dcrutil: Disallow creation of hybrid P2PK addrs ([sebitt27/dcrd#1154](https://github.com/sebitt27/dcrd/pull/1154))
- chainec/dcrec: Remove hybrid pubkey support ([sebitt27/dcrd#1155](https://github.com/sebitt27/dcrd/pull/1155))
- blockchain: Only fetch inputs once in connect txns ([sebitt27/dcrd#1152](https://github.com/sebitt27/dcrd/pull/1152))
- indexers: Provide interface for index removal ([sebitt27/dcrd#1158](https://github.com/sebitt27/dcrd/pull/1158))

### Testing and Quality Assurance:

- travis: set GOVERSION environment properly ([sebitt27/dcrd#958](https://github.com/sebitt27/dcrd/pull/958))
- stake: Override false positive vet error ([sebitt27/dcrd#960](https://github.com/sebitt27/dcrd/pull/960))
- docs: make example code compile ([sebitt27/dcrd#970](https://github.com/sebitt27/dcrd/pull/970))
- blockchain: Add median time tests ([sebitt27/dcrd#991](https://github.com/sebitt27/dcrd/pull/991))
- chaingen: Update vote commitments on hdr updates ([sebitt27/dcrd#1023](https://github.com/sebitt27/dcrd/pull/1023))
- fullblocktests: Add tests for early final state ([sebitt27/dcrd#1031](https://github.com/sebitt27/dcrd/pull/1031))
- travis: test in docker container ([sebitt27/dcrd#1053](https://github.com/sebitt27/dcrd/pull/1053))
- blockchain: Correct error stringer tests ([sebitt27/dcrd#1066](https://github.com/sebitt27/dcrd/pull/1066))
- blockchain: Remove superfluous reorg tests ([sebitt27/dcrd#1072](https://github.com/sebitt27/dcrd/pull/1072))
- blockchain: Use chaingen for forced reorg tests ([sebitt27/dcrd#1074](https://github.com/sebitt27/dcrd/pull/1074))
- blockchain: Remove superfluous test checks ([sebitt27/dcrd#1075](https://github.com/sebitt27/dcrd/pull/1075))
- blockchain: move block validation rule tests into fullblocktests ([sebitt27/dcrd#1060](https://github.com/sebitt27/dcrd/pull/1060))
- fullblocktests: Cleanup after refactor ([sebitt27/dcrd#1080](https://github.com/sebitt27/dcrd/pull/1080))
- chaingen: Prevent dup block names in NextBlock ([sebitt27/dcrd#1079](https://github.com/sebitt27/dcrd/pull/1079))
- blockchain: Remove duplicate val tests ([sebitt27/dcrd#1082](https://github.com/sebitt27/dcrd/pull/1082))
- chaingen: Break dependency on blockchain ([sebitt27/dcrd#1076](https://github.com/sebitt27/dcrd/pull/1076))
- blockchain: Consolidate tests into the main package ([sebitt27/dcrd#1077](https://github.com/sebitt27/dcrd/pull/1077))
- chaingen: Export vote commitment script function ([sebitt27/dcrd#1081](https://github.com/sebitt27/dcrd/pull/1081))
- fullblocktests: Improve vote on wrong block tests ([sebitt27/dcrd#1081](https://github.com/sebitt27/dcrd/pull/1081))
- chaingen: Export func to check if block is solved ([sebitt27/dcrd#1089](https://github.com/sebitt27/dcrd/pull/1089))
- fullblocktests: Use new exported IsSolved func ([sebitt27/dcrd#1089](https://github.com/sebitt27/dcrd/pull/1089))
- chaingen: Accept mungers for create premine block ([sebitt27/dcrd#1090](https://github.com/sebitt27/dcrd/pull/1090))
- blockchain: Add tests for chain tip tracking ([sebitt27/dcrd#1096](https://github.com/sebitt27/dcrd/pull/1096))
- blockchain: move block validation rule tests into fullblocktests (2/x) ([sebitt27/dcrd#1095](https://github.com/sebitt27/dcrd/pull/1095))
- addrmgr: Remove obsolete coverage script ([sebitt27/dcrd#1103](https://github.com/sebitt27/dcrd/pull/1103))
- chaingen: Track expected blk heights separately ([sebitt27/dcrd#1101](https://github.com/sebitt27/dcrd/pull/1101))
- addrmgr: Improve test coverage ([sebitt27/dcrd#1111](https://github.com/sebitt27/dcrd/pull/1111))
- chaingen: Add revocation generation infrastructure ([sebitt27/dcrd#1120](https://github.com/sebitt27/dcrd/pull/1120))
- fullblocktests: Add some basic revocation tests ([sebitt27/dcrd#1121](https://github.com/sebitt27/dcrd/pull/1121))
- addrmgr: Test removal of corrupt peers file ([sebitt27/dcrd#1129](https://github.com/sebitt27/dcrd/pull/1129))

### Misc:

- release: Bump for v1.2.0 ([sebitt27/dcrd#1140](https://github.com/sebitt27/dcrd/pull/1140))
- goimports -w . ([sebitt27/dcrd#968](https://github.com/sebitt27/dcrd/pull/968))
- dep: sync ([sebitt27/dcrd#980](https://github.com/sebitt27/dcrd/pull/980))
- multi: Simplify code per gosimple linter ([sebitt27/dcrd#993](https://github.com/sebitt27/dcrd/pull/993))
- multi: various cleanups ([sebitt27/dcrd#1019](https://github.com/sebitt27/dcrd/pull/1019))
- multi: release the mutex earlier ([sebitt27/dcrd#1026](https://github.com/sebitt27/dcrd/pull/1026))
- multi: fix some maligned linter warnings ([sebitt27/dcrd#1025](https://github.com/sebitt27/dcrd/pull/1025))
- blockchain: Correct a few log statements ([sebitt27/dcrd#1042](https://github.com/sebitt27/dcrd/pull/1042))
- mempool: cleaner ([sebitt27/dcrd#1050](https://github.com/sebitt27/dcrd/pull/1050))
- multi: fix misspell linter warnings ([sebitt27/dcrd#1054](https://github.com/sebitt27/dcrd/pull/1054))
- dep: sync ([sebitt27/dcrd#1091](https://github.com/sebitt27/dcrd/pull/1091))
- multi: Properly capitalize Decred ([sebitt27/dcrd#1102](https://github.com/sebitt27/dcrd/pull/1102))
- build: Correct semver build handling ([sebitt27/dcrd#1097](https://github.com/sebitt27/dcrd/pull/1097))
- main: Make func definition style consistent ([sebitt27/dcrd#1114](https://github.com/sebitt27/dcrd/pull/1114))
- main: Allow semver prerel via linker flags ([sebitt27/dcrd#1128](https://github.com/sebitt27/dcrd/pull/1128))

### Code Contributors (alphabetical order):

- Andrew Chiw
- Daniel Krawsiz
- Dave Collins
- David Hill
- Donald Adu-Poku
- Javed Khan
- Jolan Luff
- Jon Gillham
- Josh Rickmar
- Markus Richter
- Matheus Degiovani
- Ryan Vacek
