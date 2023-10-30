# dcrd v1.3.0

This release of dcrd contains significant performance enhancements for startup
speed, validation, and network operations that directly benefit lightweight
clients, such as SPV (Simplified Payment Verification) wallets, a policy change
to reduce the default minimum transaction fee rate, a new public test network
version, removal of bloom filter support, infrastructure improvements, and other
quality assurance changes.

**It is highly recommended that everyone upgrade to this latest release as it
contains many important scalability improvements and is required to be able to
use the new public test network.**

## Downgrade Warning

The database format in v1.3.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Reduction of Default Minimum Transaction Fee Rate Policy

The default setting for the policy which specifies the minimum transaction fee
rate that will be accepted and relayed to the rest of the network has been
reduced to 0.0001 DCR/kB (10,000 atoms/kB) from the previous value of 0.001
DCR/kB (100,000 atoms/kB).

Transactions should not attempt to use the reduced fee rate until the majority
of the network has upgraded to this release as otherwise the transactions will
likely have issues relaying through the network since old nodes that have not
updated their policy will reject them due to not paying a high enough fee.

### Several Speed Optimizations

This release contains several enhancements to improve speed for startup,
the initial sync process, validation, and network operations.

In order to achieve these speedups, there is a one time database migration, as
previously mentioned, that typically only takes a few seconds to complete on
most hardware.

#### Further Improved Startup Speed

The startup time has been improved by roughly 2x on both slower hard disk drives
(HDDs) and solid state drives (SSDs) as compared to v1.2.0.

#### Significantly Faster Network Operations

The ability to serve information to other peers on the network has received
several optimizations which, in addition to generally improving the overall
scalability and throughput of the network, also directly benefits SPV
(Simplified Payment Verification) clients by delivering the block headers they
require roughly 3x to 4x faster.

#### Signature Hash Calculation Optimization

Part of validating that transactions are only spending coins that the owner has
authorized involves ensuring the validity of cryptographic signatures.  This
release provides a speedup of about 75% to a key portion of that validation
which results in a roughly 20% faster initial sync process.

### Bloom Filters Removal

Bloom filters were deprecated as of the last release in favor of the more recent
privacy-preserving GCS committed filters.  Consequently, this release removes
support for bloom filters completely.  There are no known clients which use
bloom filters, however, if there are any unknown clients which use them, those
clients will need to be updated to use the GCS committed filters accordingly.

### Public Test Network Version 3

The public test network has been reset and bumped to version 3.  All of the new
consensus rules voted in by version 2 of the public test network have been
retained and are therefore active on the new version 3 test network without
having to vote them in again.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/sebitt27/dcrd/compare/release-v1.2.0...release-v1.3.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.3.0 release ([sebitt27/dcrd#1385](https://github.com/sebitt27/dcrd/pull/1385))
- multi: Remove everything to do about bloom filters ([sebitt27/dcrd#1162](https://github.com/sebitt27/dcrd/pull/1162))
- wire: Remove TxSerializeWitnessSigning ([sebitt27/dcrd#1180](https://github.com/sebitt27/dcrd/pull/1180))
- addrmgr: Skip low quality addresses for getaddr ([sebitt27/dcrd#1135](https://github.com/sebitt27/dcrd/pull/1135))
- addrmgr: Fix race in save peers ([sebitt27/dcrd#1259](https://github.com/sebitt27/dcrd/pull/1259))
- server: Only respond to getaddr once per conn ([sebitt27/dcrd#1257](https://github.com/sebitt27/dcrd/pull/1257))
- peer: Rework version negotiation ([sebitt27/dcrd#1250](https://github.com/sebitt27/dcrd/pull/1250))
- peer: Allow OnVersion callback to reject peer ([sebitt27/dcrd#1251](https://github.com/sebitt27/dcrd/pull/1251))
- server: Reject outbound conns to non-full nodes ([sebitt27/dcrd#1252](https://github.com/sebitt27/dcrd/pull/1252))
- peer: Improve net address service adverts ([sebitt27/dcrd#1253](https://github.com/sebitt27/dcrd/pull/1253))
- addrmgr: Expose method to update services ([sebitt27/dcrd#1254](https://github.com/sebitt27/dcrd/pull/1254))
- server: Update addrmgr services on outbound conns ([sebitt27/dcrd#1254](https://github.com/sebitt27/dcrd/pull/1254))
- server: Use local inbound var in version handler ([sebitt27/dcrd#1255](https://github.com/sebitt27/dcrd/pull/1255))
- server: Only advertise local addr when current ([sebitt27/dcrd#1256](https://github.com/sebitt27/dcrd/pull/1256))
- server: Use local addr var in version handler ([sebitt27/dcrd#1258](https://github.com/sebitt27/dcrd/pull/1258))
- chaincfg: split params into per-network files ([sebitt27/dcrd#1265](https://github.com/sebitt27/dcrd/pull/1265))
- server: Always reply to getheaders with headers ([sebitt27/dcrd#1295](https://github.com/sebitt27/dcrd/pull/1295))
- addrmgr: skip never-successful addresses ([sebitt27/dcrd#1313](https://github.com/sebitt27/dcrd/pull/1313))
- multi: Introduce default coin type for SLIP0044 ([sebitt27/dcrd#1293](https://github.com/sebitt27/dcrd/pull/1293))
- blockchain: Modify diff redux logic for testnet ([sebitt27/dcrd#1387](https://github.com/sebitt27/dcrd/pull/1387))
- multi: Reset testnet and bump to version 3 ([sebitt27/dcrd#1387](https://github.com/sebitt27/dcrd/pull/1387))
- multi: Remove testnet version 2 defs and refs ([sebitt27/dcrd#1387](https://github.com/sebitt27/dcrd/pull/1387))

### Transaction relay (memory pool):

- policy: Lower default relay fee to 0.0001/kB ([sebitt27/dcrd#1202](https://github.com/sebitt27/dcrd/pull/1202))
- mempool: Use blockchain for tx expiry check ([sebitt27/dcrd#1199](https://github.com/sebitt27/dcrd/pull/1199))
- mempool: use secp256k1 functions directly ([sebitt27/dcrd#1213](https://github.com/sebitt27/dcrd/pull/1213))
- mempool: Make expiry pruning self contained ([sebitt27/dcrd#1378](https://github.com/sebitt27/dcrd/pull/1378))
- mempool: Stricter orphan evaluation and eviction ([sebitt27/dcrd#1207](https://github.com/sebitt27/dcrd/pull/1207))
- mempool: use secp256k1 functions directly ([sebitt27/dcrd#1213](https://github.com/sebitt27/dcrd/pull/1213))
- multi: add specialized rebroadcast handling for stake txs ([sebitt27/dcrd#979](https://github.com/sebitt27/dcrd/pull/979))
- mempool: Make expiry pruning self contained ([sebitt27/dcrd#1378](https://github.com/sebitt27/dcrd/pull/1378))

### RPC:

- rpcserver: Improve JSON-RPC compatibility ([sebitt27/dcrd#1150](https://github.com/sebitt27/dcrd/pull/1150))
- rpcserver: Correct rebroadcastwinners handler ([sebitt27/dcrd#1234](https://github.com/sebitt27/dcrd/pull/1234))
- dcrjson: Add Expiry field to CreateRawTransactionCmd ([sebitt27/dcrd#1149](https://github.com/sebitt27/dcrd/pull/1149))
- dcrjson: add estimatesmartfee ([sebitt27/dcrd#1201](https://github.com/sebitt27/dcrd/pull/1201))
- rpc: Use upstream gorilla/websocket ([sebitt27/dcrd#1218](https://github.com/sebitt27/dcrd/pull/1218))
- dcrjson: add createvotingaccount and dropvotingaccount rpc methods ([sebitt27/dcrd#1217](https://github.com/sebitt27/dcrd/pull/1217))
- multi: Change NoSplitTransaction param to SplitTx ([sebitt27/dcrd#1231](https://github.com/sebitt27/dcrd/pull/1231))
- rpcclient: pass default value for NewPurchaseTicketCmd's comment param ([sebitt27/dcrd#1232](https://github.com/sebitt27/dcrd/pull/1232))
- multi: No winning ticket ntfns for big reorg depth ([sebitt27/dcrd#1235](https://github.com/sebitt27/dcrd/pull/1235))
- multi: modify PurchaseTicketCmd ([sebitt27/dcrd#1241](https://github.com/sebitt27/dcrd/pull/1241))
- multi: move extension commands into associated normal command files ([sebitt27/dcrd#1238](https://github.com/sebitt27/dcrd/pull/1238))
- dcrjson: Fix NewCreateRawTransactionCmd comment ([sebitt27/dcrd#1262](https://github.com/sebitt27/dcrd/pull/1262))
- multi: revert TicketChange addition to PurchaseTicketCmd ([sebitt27/dcrd#1278](https://github.com/sebitt27/dcrd/pull/1278))
- rpcclient: Implement fmt.Stringer for Client ([sebitt27/dcrd#1298](https://github.com/sebitt27/dcrd/pull/1298))
- multi: add amount field to TransactionInput ([sebitt27/dcrd#1297](https://github.com/sebitt27/dcrd/pull/1297))
- dcrjson: Ready GetStakeInfoResult for SPV wallets ([sebitt27/dcrd#1333](https://github.com/sebitt27/dcrd/pull/1333))
- dcrjson: add fundrawtransaction command ([sebitt27/dcrd#1316](https://github.com/sebitt27/dcrd/pull/1316))
- dcrjson: Make linter happy by renaming Id to ID ([sebitt27/dcrd#1374](https://github.com/sebitt27/dcrd/pull/1374))
- dcrjson: Remove unused vote bit concat codec funcs ([sebitt27/dcrd#1384](https://github.com/sebitt27/dcrd/pull/1384))
- rpcserver: Cleanup cfilter handling ([sebitt27/dcrd#1398](https://github.com/sebitt27/dcrd/pull/1398))

### dcrd command-line flags and configuration:

- multi: Correct clean and expand path handling ([sebitt27/dcrd#1186](https://github.com/sebitt27/dcrd/pull/1186))

### dcrctl utility changes:

- dcrctl: Fix --skipverify failing if rpc.cert not found ([sebitt27/dcrd#1163](https://github.com/sebitt27/dcrd/pull/1163))

### Documentation:

- hdkeychain: Correct hash algorithm in comment ([sebitt27/dcrd#1171](https://github.com/sebitt27/dcrd/pull/1171))
- Fix typo in blockchain ([sebitt27/dcrd#1185](https://github.com/sebitt27/dcrd/pull/1185))
- docs: Update node.js example for v8.11.1 LTS ([sebitt27/dcrd#1209](https://github.com/sebitt27/dcrd/pull/1209))
- docs: Update txaccepted method in json_rpc_api.md ([sebitt27/dcrd#1242](https://github.com/sebitt27/dcrd/pull/1242))
- docs: Correct blockmaxsize and blockprioritysize ([sebitt27/dcrd#1339](https://github.com/sebitt27/dcrd/pull/1339))
- server: Correct comment in getblocks handler ([sebitt27/dcrd#1269](https://github.com/sebitt27/dcrd/pull/1269))
- config: Fix typo ([sebitt27/dcrd#1274](https://github.com/sebitt27/dcrd/pull/1274))
- multi: Fix badges in README ([sebitt27/dcrd#1277](https://github.com/sebitt27/dcrd/pull/1277))
- stake: Correct comment in connectNode ([sebitt27/dcrd#1325](https://github.com/sebitt27/dcrd/pull/1325))
- txscript: Update comments for removal of flags ([sebitt27/dcrd#1336](https://github.com/sebitt27/dcrd/pull/1336))
- docs: Update docs for versioned modules ([sebitt27/dcrd#1391](https://github.com/sebitt27/dcrd/pull/1391))
- mempool: Correct min relay tx fee comment to DCR ([sebitt27/dcrd#1396](https://github.com/sebitt27/dcrd/pull/1396))

### Developer-related package and module changes:

- blockchain: CheckConnectBlockTemplate with tests ([sebitt27/dcrd#1086](https://github.com/sebitt27/dcrd/pull/1086))
- addrmgr: Simplify package API ([sebitt27/dcrd#1136](https://github.com/sebitt27/dcrd/pull/1136))
- txscript: Remove unused strict multisig flag ([sebitt27/dcrd#1203](https://github.com/sebitt27/dcrd/pull/1203))
- txscript: Move sig hash logic to separate file ([sebitt27/dcrd#1174](https://github.com/sebitt27/dcrd/pull/1174))
- txscript: Remove SigHashAllValue ([sebitt27/dcrd#1175](https://github.com/sebitt27/dcrd/pull/1175))
- txscript: Decouple and optimize sighash calc ([sebitt27/dcrd#1179](https://github.com/sebitt27/dcrd/pull/1179))
- wire: Remove TxSerializeWitnessValueSigning ([sebitt27/dcrd#1176](https://github.com/sebitt27/dcrd/pull/1176))
- hdkeychain: Satisfy fmt.Stringer interface ([sebitt27/dcrd#1168](https://github.com/sebitt27/dcrd/pull/1168))
- blockchain: Validate tx expiry in block context ([sebitt27/dcrd#1187](https://github.com/sebitt27/dcrd/pull/1187))
- blockchain: rename ErrRegTxSpendStakeOut to ErrRegTxCreateStakeOut ([sebitt27/dcrd#1195](https://github.com/sebitt27/dcrd/pull/1195))
- multi: Break coinbase dep on standardness rules ([sebitt27/dcrd#1196](https://github.com/sebitt27/dcrd/pull/1196))
- txscript: Cleanup code for the substr opcode ([sebitt27/dcrd#1206](https://github.com/sebitt27/dcrd/pull/1206))
- multi: use secp256k1 types and fields directly ([sebitt27/dcrd#1211](https://github.com/sebitt27/dcrd/pull/1211))
- dcrec: add Pubkey func to secp256k1 and edwards elliptic curves ([sebitt27/dcrd#1214](https://github.com/sebitt27/dcrd/pull/1214))
- blockchain: use secp256k1 functions directly ([sebitt27/dcrd#1212](https://github.com/sebitt27/dcrd/pull/1212))
- multi: Replace btclog with slog ([sebitt27/dcrd#1216](https://github.com/sebitt27/dcrd/pull/1216))
- multi: Define vgo modules ([sebitt27/dcrd#1223](https://github.com/sebitt27/dcrd/pull/1223))
- chainhash: Define vgo module ([sebitt27/dcrd#1224](https://github.com/sebitt27/dcrd/pull/1224))
- wire: Refine vgo deps ([sebitt27/dcrd#1225](https://github.com/sebitt27/dcrd/pull/1225))
- addrmrg: Refine vgo deps ([sebitt27/dcrd#1226](https://github.com/sebitt27/dcrd/pull/1226))
- chaincfg: Refine vgo deps ([sebitt27/dcrd#1227](https://github.com/sebitt27/dcrd/pull/1227))
- multi: Return fork len from ProcessBlock ([sebitt27/dcrd#1233](https://github.com/sebitt27/dcrd/pull/1233))
- blockchain: Panic on fatal assertions ([sebitt27/dcrd#1243](https://github.com/sebitt27/dcrd/pull/1243))
- blockchain: Convert to full block index in mem ([sebitt27/dcrd#1229](https://github.com/sebitt27/dcrd/pull/1229))
- blockchain: Optimize checkpoint handling ([sebitt27/dcrd#1230](https://github.com/sebitt27/dcrd/pull/1230))
- blockchain: Optimize block locator generation ([sebitt27/dcrd#1237](https://github.com/sebitt27/dcrd/pull/1237))
- multi: Refactor and optimize inv discovery ([sebitt27/dcrd#1239](https://github.com/sebitt27/dcrd/pull/1239))
- peer: Minor function definition order cleanup ([sebitt27/dcrd#1247](https://github.com/sebitt27/dcrd/pull/1247))
- peer: Remove superfluous dup version check ([sebitt27/dcrd#1248](https://github.com/sebitt27/dcrd/pull/1248))
- txscript: export canonicalDataSize ([sebitt27/dcrd#1266](https://github.com/sebitt27/dcrd/pull/1266))
- blockchain: Add BuildMerkleTreeStore alternative for MsgTx ([sebitt27/dcrd#1268](https://github.com/sebitt27/dcrd/pull/1268))
- blockchain: Optimize exported header access ([sebitt27/dcrd#1273](https://github.com/sebitt27/dcrd/pull/1273))
- txscript: Cleanup P2SH and stake opcode handling ([sebitt27/dcrd#1318](https://github.com/sebitt27/dcrd/pull/1318))
- txscript: Significantly improve errors ([sebitt27/dcrd#1319](https://github.com/sebitt27/dcrd/pull/1319))
- txscript: Remove pay-to-script-hash flag ([sebitt27/dcrd#1321](https://github.com/sebitt27/dcrd/pull/1321))
- txscript: Remove DER signature verification flag ([sebitt27/dcrd#1323](https://github.com/sebitt27/dcrd/pull/1323))
- txscript: Remove verify minimal data flag ([sebitt27/dcrd#1326](https://github.com/sebitt27/dcrd/pull/1326))
- txscript: Remove script num require minimal flag ([sebitt27/dcrd#1328](https://github.com/sebitt27/dcrd/pull/1328))
- txscript: Make PeekInt consistent with PopInt ([sebitt27/dcrd#1329](https://github.com/sebitt27/dcrd/pull/1329))
- build: Add experimental support for vgo ([sebitt27/dcrd#1215](https://github.com/sebitt27/dcrd/pull/1215))
- build: Update some vgo dependencies to use tags ([sebitt27/dcrd#1219](https://github.com/sebitt27/dcrd/pull/1219))
- stake: add ExpiredByBlock to stake.Node ([sebitt27/dcrd#1221](https://github.com/sebitt27/dcrd/pull/1221))
- server: Minor function definition order cleanup ([sebitt27/dcrd#1271](https://github.com/sebitt27/dcrd/pull/1271))
- server: Convert CF code to use new inv discovery ([sebitt27/dcrd#1272](https://github.com/sebitt27/dcrd/pull/1272))
- multi: add valueIn parameter to wire.NewTxIn ([sebitt27/dcrd#1287](https://github.com/sebitt27/dcrd/pull/1287))
- txscript: Remove low S verification flag ([sebitt27/dcrd#1308](https://github.com/sebitt27/dcrd/pull/1308))
- txscript: Remove unused old sig hash type ([sebitt27/dcrd#1309](https://github.com/sebitt27/dcrd/pull/1309))
- txscript: Remove strict encoding verification flag ([sebitt27/dcrd#1310](https://github.com/sebitt27/dcrd/pull/1310))
- blockchain: Combine block by hash functions ([sebitt27/dcrd#1330](https://github.com/sebitt27/dcrd/pull/1330))
- multi: Continue conversion from chainec to dcrec ([sebitt27/dcrd#1304](https://github.com/sebitt27/dcrd/pull/1304))
- multi: Remove unused secp256k1 sig parse parameter ([sebitt27/dcrd#1335](https://github.com/sebitt27/dcrd/pull/1335))
- blockchain: Refactor db main chain idx to blk idx ([sebitt27/dcrd#1332](https://github.com/sebitt27/dcrd/pull/1332))
- blockchain: Remove main chain index from db ([sebitt27/dcrd#1334](https://github.com/sebitt27/dcrd/pull/1334))
- blockchain: Implement new chain view ([sebitt27/dcrd#1337](https://github.com/sebitt27/dcrd/pull/1337))
- blockmanager: remove unused Pause() API ([sebitt27/dcrd#1340](https://github.com/sebitt27/dcrd/pull/1340))
- chainhash: Remove dup code from hash funcs ([sebitt27/dcrd#1342](https://github.com/sebitt27/dcrd/pull/1342))
- connmgr: Fix the ConnReq print out causing panic ([sebitt27/dcrd#1345](https://github.com/sebitt27/dcrd/pull/1345))
- gcs: Pool MatchAny data allocations ([sebitt27/dcrd#1348](https://github.com/sebitt27/dcrd/pull/1348))
- blockchain: Faster chain view block locator ([sebitt27/dcrd#1338](https://github.com/sebitt27/dcrd/pull/1338))
- blockchain: Refactor to use new chain view ([sebitt27/dcrd#1344](https://github.com/sebitt27/dcrd/pull/1344))
- blockchain: Remove unnecessary genesis block check ([sebitt27/dcrd#1368](https://github.com/sebitt27/dcrd/pull/1368))
- chainhash: Update go build module support ([sebitt27/dcrd#1358](https://github.com/sebitt27/dcrd/pull/1358))
- wire: Update go build module support ([sebitt27/dcrd#1359](https://github.com/sebitt27/dcrd/pull/1359))
- addrmgr: Update go build module support ([sebitt27/dcrd#1360](https://github.com/sebitt27/dcrd/pull/1360))
- chaincfg: Update go build module support ([sebitt27/dcrd#1361](https://github.com/sebitt27/dcrd/pull/1361))
- connmgr: Refine go build module support ([sebitt27/dcrd#1363](https://github.com/sebitt27/dcrd/pull/1363))
- secp256k1: Refine go build module support ([sebitt27/dcrd#1362](https://github.com/sebitt27/dcrd/pull/1362))
- dcrec: Refine go build module support ([sebitt27/dcrd#1364](https://github.com/sebitt27/dcrd/pull/1364))
- certgen: Update go build module support ([sebitt27/dcrd#1365](https://github.com/sebitt27/dcrd/pull/1365))
- dcrutil: Refine go build module support ([sebitt27/dcrd#1366](https://github.com/sebitt27/dcrd/pull/1366))
- hdkeychain: Refine go build module support ([sebitt27/dcrd#1369](https://github.com/sebitt27/dcrd/pull/1369))
- txscript: Refine go build module support ([sebitt27/dcrd#1370](https://github.com/sebitt27/dcrd/pull/1370))
- multi: Remove go modules that do not build ([sebitt27/dcrd#1371](https://github.com/sebitt27/dcrd/pull/1371))
- database: Refine go build module support ([sebitt27/dcrd#1372](https://github.com/sebitt27/dcrd/pull/1372))
- build: Refine build module support ([sebitt27/dcrd#1384](https://github.com/sebitt27/dcrd/pull/1384))
- blockmanager: make pruning transactions consistent ([sebitt27/dcrd#1376](https://github.com/sebitt27/dcrd/pull/1376))
- blockchain: Optimize reorg to use known status ([sebitt27/dcrd#1367](https://github.com/sebitt27/dcrd/pull/1367))
- blockchain: Make block index flushable ([sebitt27/dcrd#1375](https://github.com/sebitt27/dcrd/pull/1375))
- blockchain: Mark fastadd block valid ([sebitt27/dcrd#1392](https://github.com/sebitt27/dcrd/pull/1392))
- release: Bump module versions and deps ([sebitt27/dcrd#1390](https://github.com/sebitt27/dcrd/pull/1390))
- blockchain: Mark fastadd block valid ([sebitt27/dcrd#1392](https://github.com/sebitt27/dcrd/pull/1392))
- gcs: use dchest/siphash ([sebitt27/dcrd#1395](https://github.com/sebitt27/dcrd/pull/1395))
- dcrec: Make function defs more consistent ([sebitt27/dcrd#1432](https://github.com/sebitt27/dcrd/pull/1432))

### Testing and Quality Assurance:

- addrmgr: Simplify tests for KnownAddress ([sebitt27/dcrd#1133](https://github.com/sebitt27/dcrd/pull/1133))
- blockchain: move block validation rule tests into fullblocktests ([sebitt27/dcrd#1141](https://github.com/sebitt27/dcrd/pull/1141))
- addrmgr: Test timestamp update during AddAddress ([sebitt27/dcrd#1137](https://github.com/sebitt27/dcrd/pull/1137))
- txscript: Consolidate tests into txscript package ([sebitt27/dcrd#1177](https://github.com/sebitt27/dcrd/pull/1177))
- txscript: Add JSON-based signature hash tests ([sebitt27/dcrd#1178](https://github.com/sebitt27/dcrd/pull/1178))
- txscript: Correct JSON-based signature hash tests ([sebitt27/dcrd#1181](https://github.com/sebitt27/dcrd/pull/1181))
- txscript: Add benchmark for sighash calculation ([sebitt27/dcrd#1179](https://github.com/sebitt27/dcrd/pull/1179))
- mempool: Refactor pool membership test logic ([sebitt27/dcrd#1188](https://github.com/sebitt27/dcrd/pull/1188))
- blockchain: utilize CalcNextReqStakeDifficulty in fullblocktests ([sebitt27/dcrd#1189](https://github.com/sebitt27/dcrd/pull/1189))
- fullblocktests: add additional premine and malformed tests ([sebitt27/dcrd#1190](https://github.com/sebitt27/dcrd/pull/1190))
- txscript: Improve substr opcode test coverage ([sebitt27/dcrd#1205](https://github.com/sebitt27/dcrd/pull/1205))
- txscript: Convert reference tests to new format ([sebitt27/dcrd#1320](https://github.com/sebitt27/dcrd/pull/1320))
- txscript: Remove P2SH flag from test data ([sebitt27/dcrd#1322](https://github.com/sebitt27/dcrd/pull/1322))
- txscript: Remove DERSIG flag from test data ([sebitt27/dcrd#1324](https://github.com/sebitt27/dcrd/pull/1324))
- txscript: Remove MINIMALDATA flag from test data ([sebitt27/dcrd#1327](https://github.com/sebitt27/dcrd/pull/1327))
- fullblocktests: Add expired stake tx test ([sebitt27/dcrd#1184](https://github.com/sebitt27/dcrd/pull/1184))
- travis: simplify Docker files ([sebitt27/dcrd#1275](https://github.com/sebitt27/dcrd/pull/1275))
- docker: Add dockerfiles for running dcrd nodes ([sebitt27/dcrd#1317](https://github.com/sebitt27/dcrd/pull/1317))
- blockchain: Improve spend journal tests ([sebitt27/dcrd#1246](https://github.com/sebitt27/dcrd/pull/1246))
- txscript: Cleanup and add tests for left opcode ([sebitt27/dcrd#1281](https://github.com/sebitt27/dcrd/pull/1281))
- txscript: Cleanup and add tests for right opcode ([sebitt27/dcrd#1282](https://github.com/sebitt27/dcrd/pull/1282))
- txscript: Cleanup and add tests for the cat opcode ([sebitt27/dcrd#1283](https://github.com/sebitt27/dcrd/pull/1283))
- txscript: Cleanup and add tests for rotr opcode ([sebitt27/dcrd#1285](https://github.com/sebitt27/dcrd/pull/1285))
- txscript: Cleanup and add tests for rotl opcode ([sebitt27/dcrd#1286](https://github.com/sebitt27/dcrd/pull/1286))
- txscript: Cleanup and add tests for lshift opcode ([sebitt27/dcrd#1288](https://github.com/sebitt27/dcrd/pull/1288))
- txscript: Cleanup and add tests for rshift opcode ([sebitt27/dcrd#1289](https://github.com/sebitt27/dcrd/pull/1289))
- txscript: Cleanup and add tests for div opcode ([sebitt27/dcrd#1290](https://github.com/sebitt27/dcrd/pull/1290))
- txscript: Cleanup and add tests for mod opcode ([sebitt27/dcrd#1291](https://github.com/sebitt27/dcrd/pull/1291))
- txscript: Update CSV to match tests in DCP0003 ([sebitt27/dcrd#1292](https://github.com/sebitt27/dcrd/pull/1292))
- txscript: Introduce repeated syntax to test data ([sebitt27/dcrd#1299](https://github.com/sebitt27/dcrd/pull/1299))
- txscript: Allow multi opcode test data repeat ([sebitt27/dcrd#1300](https://github.com/sebitt27/dcrd/pull/1300))
- txscript: Improve and correct some script tests ([sebitt27/dcrd#1303](https://github.com/sebitt27/dcrd/pull/1303))
- main: verify network pow limits ([sebitt27/dcrd#1302](https://github.com/sebitt27/dcrd/pull/1302))
- txscript: Remove STRICTENC flag from test data ([sebitt27/dcrd#1311](https://github.com/sebitt27/dcrd/pull/1311))
- txscript: Cleanup plus tests for checksig opcodes ([sebitt27/dcrd#1315](https://github.com/sebitt27/dcrd/pull/1315))
- blockchain: Add negative tests for forced reorg ([sebitt27/dcrd#1341](https://github.com/sebitt27/dcrd/pull/1341))
- dcrjson: Consolidate tests into dcrjson package ([sebitt27/dcrd#1373](https://github.com/sebitt27/dcrd/pull/1373))
- txscript: add additional data push op code tests ([sebitt27/dcrd#1346](https://github.com/sebitt27/dcrd/pull/1346))
- txscript: add/group control op code tests ([sebitt27/dcrd#1349](https://github.com/sebitt27/dcrd/pull/1349))
- txscript: add/group stack op code tests ([sebitt27/dcrd#1350](https://github.com/sebitt27/dcrd/pull/1350))
- txscript: group splice opcode tests ([sebitt27/dcrd#1351](https://github.com/sebitt27/dcrd/pull/1351))
- txscript: add/group bitwise logic, comparison & rotation op code tests ([sebitt27/dcrd#1352](https://github.com/sebitt27/dcrd/pull/1352))
- txscript: add/group numeric related opcode tests ([sebitt27/dcrd#1353](https://github.com/sebitt27/dcrd/pull/1353))
- txscript: group reserved op code tests ([sebitt27/dcrd#1355](https://github.com/sebitt27/dcrd/pull/1355))
- txscript: add/group crypto related op code tests ([sebitt27/dcrd#1354](https://github.com/sebitt27/dcrd/pull/1354))
- multi: Reduce testnet2 refs in unit tests ([sebitt27/dcrd#1387](https://github.com/sebitt27/dcrd/pull/1387))
- blockchain: Avoid deployment expiration in tests ([sebitt27/dcrd#1450](https://github.com/sebitt27/dcrd/pull/1450))

### Misc:

- release: Bump for v1.3.0 ([sebitt27/dcrd#1388](https://github.com/sebitt27/dcrd/pull/1388))
- multi: Correct typos found by misspell ([sebitt27/dcrd#1197](https://github.com/sebitt27/dcrd/pull/1197))
- main: Correct mem profile error message ([sebitt27/dcrd#1183](https://github.com/sebitt27/dcrd/pull/1183))
- multi: Use saner permissions saving certs ([sebitt27/dcrd#1263](https://github.com/sebitt27/dcrd/pull/1263))
- server: only call time.Now() once ([sebitt27/dcrd#1313](https://github.com/sebitt27/dcrd/pull/1313))
- multi: linter cleanup ([sebitt27/dcrd#1305](https://github.com/sebitt27/dcrd/pull/1305))
- multi: Remove unnecessary network name funcs ([sebitt27/dcrd#1387](https://github.com/sebitt27/dcrd/pull/1387))
- config: Warn if testnet2 database exists ([sebitt27/dcrd#1389](https://github.com/sebitt27/dcrd/pull/1389))

### Code Contributors (alphabetical order):

- Dave Collins
- David Hill
- Dmitry Fedorov
- Donald Adu-Poku
- harzo
- hypernoob
- J Fixby
- Jonathan Chappelow
- Josh Rickmar
- Markus Richter
- matadormel
- Matheus Degiovani
- Michael Eze
- Orthomind
- Shuai Qi
- Tibor Bősze
- Victor Oliveira
