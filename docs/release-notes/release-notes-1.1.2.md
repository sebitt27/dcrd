# dcrd v1.1.2

This release of dcrd primarily contains performance enhancements, infrastructure
improvements, and other quality assurance changes.

While it is not visible in this release, significant infrastructure work has
also been done this release cycle towards porting the Lightning Network (LN)
daemon which will ultimately allow LN payments to be backed by Decred.

## Notable Changes

### Faster Block Validation

A significant portion of block validation involves handling the stake tickets
which form an integral part of Decred's hybrid proof-of-work and proof-of-stake
system.  The code which handles this portion of validation has been
significantly optimized in this release such that overall block validation is
up to approximately 3 times faster depending on the specific underlying hardware
configuration.  This also has a noticeable impact on the speed of the initial
block download process as well as how quickly votes for winning tickets are
submitted to the network.

### Data Carrier Transaction Standardness Policy

The standard policy for transaction relay of data carrier transaction outputs
has been modified to support canonically-encoded small data pushes.  These
outputs are also known as `OP_RETURN` or `nulldata` outputs.  In particular,
single byte small integers data pushes (0-16) are now supported.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/sebitt27/dcrd/compare/v1.1.0...v1.1.2).

### Protocol and network:
- chaincfg: update checkpoints for 1.1.2 release [sebitt27/dcrd#946](https://github.com/sebitt27/dcrd/pull/946)
- chaincfg: Rename one of the testnet seeders [sebitt27/dcrd#873](https://github.com/sebitt27/dcrd/pull/873)
- stake: treap index perf improvement [sebitt27/dcrd#853](https://github.com/sebitt27/dcrd/pull/853)
- stake: ticket expiry perf improvement [sebitt27/dcrd#853](https://github.com/sebitt27/dcrd/pull/853)

### Transaction relay (memory pool):

- txscript: Correct nulldata standardness check [sebitt27/dcrd#935](https://github.com/sebitt27/dcrd/pull/935)

### RPC:

- rpcserver: searchrawtransactions skip first input for vote tx [sebitt27/dcrd#859](https://github.com/sebitt27/dcrd/pull/859)
- multi: update stakebase tx vin[0] structure [sebitt27/dcrd#859](https://github.com/sebitt27/dcrd/pull/859)
- rpcserver: Fix empty ssgen verbose results [sebitt27/dcrd#871](https://github.com/sebitt27/dcrd/pull/871)
- rpcserver: check for error in getwork request [sebitt27/dcrd#898](https://github.com/sebitt27/dcrd/pull/898)
- multi: Add NoSplitTransaction to purchaseticket [sebitt27/dcrd#904](https://github.com/sebitt27/dcrd/pull/904)
- rpcserver: avoid nested decodescript p2sh addrs [sebitt27/dcrd#929](https://github.com/sebitt27/dcrd/pull/929)
- rpcserver: skip generating certs when nolisten set [sebitt27/dcrd#932](https://github.com/sebitt27/dcrd/pull/932)
- rpc: Add localaddr and relaytxes to getpeerinfo [sebitt27/dcrd#933](https://github.com/sebitt27/dcrd/pull/933)
- rpcserver: update handleSendRawTransaction error handling [sebitt27/dcrd#939](https://github.com/sebitt27/dcrd/pull/939)

### dcrd command-line flags:

- config: add --nofilelogging option [sebitt27/dcrd#872](https://github.com/sebitt27/dcrd/pull/872)

### Documentation:

- rpcclient: Remove docker info from README.md [sebitt27/dcrd#886](https://github.com/sebitt27/dcrd/pull/886)
- bloom: Fix link in README [sebitt27/dcrd#922](https://github.com/sebitt27/dcrd/pull/922)
- doc: tiny fix url [sebitt27/dcrd#928](https://github.com/sebitt27/dcrd/pull/928)
- doc: update go version for example test run in readme [sebitt27/dcrd#936](https://github.com/sebitt27/dcrd/pull/936)

### Developer-related package changes:

- multi: Drop glide, use dep [sebitt27/dcrd#818](https://github.com/sebitt27/dcrd/pull/818)
- txsort: Implement stable tx sorting package  [sebitt27/dcrd#940](https://github.com/sebitt27/dcrd/pull/940)
- coinset: Remove package [sebitt27/dcrd#888](https://github.com/sebitt27/dcrd/pull/888)
- base58: Use new github.com/decred/base58 package [sebitt27/dcrd#888](https://github.com/sebitt27/dcrd/pull/888)
- certgen: Move self signed certificate code into package [sebitt27/dcrd#879](https://github.com/sebitt27/dcrd/pull/879)
- certgen: Add doc.go and README.md [sebitt27/dcrd#883](https://github.com/sebitt27/dcrd/pull/883)
- rpcclient: Allow request-scoped cancellation during Connect [sebitt27/dcrd#880](https://github.com/sebitt27/dcrd/pull/880)
- rpcclient: Import dcrrpcclient repo into rpcclient directory [sebitt27/dcrd#880](https://github.com/sebitt27/dcrd/pull/880)
- rpcclient: json unmarshal into unexported embedded pointer  [sebitt27/dcrd#941](https://github.com/sebitt27/dcrd/pull/941)
- bloom: Copy github.com/decred/dcrutil/bloom to bloom package [sebitt27/dcrd#881](https://github.com/sebitt27/dcrd/pull/881)
- Improve gitignore [sebitt27/dcrd#887](https://github.com/sebitt27/dcrd/pull/887)
- dcrutil: Import dcrutil repo under dcrutil directory [sebitt27/dcrd#888](https://github.com/sebitt27/dcrd/pull/888)
- hdkeychain: Move to github.com/sebitt27/dcrd/hdkeychain [sebitt27/dcrd#892](https://github.com/sebitt27/dcrd/pull/892)
- stake: Add IsStakeSubmission [sebitt27/dcrd#907](https://github.com/sebitt27/dcrd/pull/907)
- txscript: Require SHA256 secret hashes for atomic swaps [sebitt27/dcrd#930](https://github.com/sebitt27/dcrd/pull/930)

### Testing and Quality Assurance:

- gometalinter: run on subpkgs too [sebitt27/dcrd#878](https://github.com/sebitt27/dcrd/pull/878)
- travis: test Gopkg.lock [sebitt27/dcrd#889](https://github.com/sebitt27/dcrd/pull/889)
- hdkeychain: Work around go vet issue with examples [sebitt27/dcrd#890](https://github.com/sebitt27/dcrd/pull/890)
- bloom: Add missing import to examples [sebitt27/dcrd#891](https://github.com/sebitt27/dcrd/pull/891)
- bloom: workaround go vet issue in example [sebitt27/dcrd#895](https://github.com/sebitt27/dcrd/pull/895)
- tests: make lockfile test work locally [sebitt27/dcrd#894](https://github.com/sebitt27/dcrd/pull/894)
- peer: Avoid goroutine leaking during handshake timeout [sebitt27/dcrd#909](https://github.com/sebitt27/dcrd/pull/909)
- travis: add gosimple linter [sebitt27/dcrd#897](https://github.com/sebitt27/dcrd/pull/897)
- multi: Handle detected data race conditions [sebitt27/dcrd#920](https://github.com/sebitt27/dcrd/pull/920)
- travis: add ineffassign linter [sebitt27/dcrd#896](https://github.com/sebitt27/dcrd/pull/896)
- rpctest: Choose flags based on provided params [sebitt27/dcrd#937](https://github.com/sebitt27/dcrd/pull/937)

### Misc:

- gofmt [sebitt27/dcrd#876](https://github.com/sebitt27/dcrd/pull/876)
- dep: sync third-party deps [sebitt27/dcrd#877](https://github.com/sebitt27/dcrd/pull/877)
- Bump for v1.1.2 [sebitt27/dcrd#916](https://github.com/sebitt27/dcrd/pull/916)
- dep: Use upstream jrick/bitset [sebitt27/dcrd#899](https://github.com/sebitt27/dcrd/pull/899)
- blockchain: removed unused funcs and vars [sebitt27/dcrd#900](https://github.com/sebitt27/dcrd/pull/900)
- blockchain: remove unused file [sebitt27/dcrd#900](https://github.com/sebitt27/dcrd/pull/900)
- rpcserver: nil pointer dereference when submit orphan block [sebitt27/dcrd#906](https://github.com/sebitt27/dcrd/pull/906)
- multi: remove unused funcs and vars [sebitt27/dcrd#901](https://github.com/sebitt27/dcrd/pull/901)

### Code Contributors (alphabetical order):

- Alex Yocom-Piatt
- Dave Collins
- David Hill
- detailyang
- Donald Adu-Poku
- Federico Gimenez
- Jason Zavaglia
- John C. Vernaleo
- Jonathan Chappelow
- Jolan Luff
- Josh Rickmar
- Maninder Lall
- Matheus Degiovani
- Nicola Larosa
- Samarth Hattangady
- Ugwueze Onyekachi Michael
