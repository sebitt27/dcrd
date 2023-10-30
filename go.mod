module github.com/sebitt27/dcrd

go 1.19

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/base58 v1.0.5
	github.com/sebitt27/dcrd/addrmgr/v2 v2.0.2
	github.com/sebitt27/dcrd/bech32 v1.1.3
	github.com/sebitt27/dcrd/blockchain/stake/v5 v5.0.0
	github.com/sebitt27/dcrd/blockchain/standalone/v2 v2.2.0
	github.com/sebitt27/dcrd/blockchain/v5 v5.0.0
	github.com/sebitt27/dcrd/certgen v1.1.2
	github.com/sebitt27/dcrd/chaincfg/chainhash v1.0.4
	github.com/sebitt27/dcrd/chaincfg/v3 v3.2.0
	github.com/sebitt27/dcrd/connmgr/v3 v3.1.1
	github.com/sebitt27/dcrd/container/apbf v1.0.1
	github.com/sebitt27/dcrd/crypto/blake256 v1.0.1
	github.com/sebitt27/dcrd/crypto/ripemd160 v1.0.2
	github.com/sebitt27/dcrd/database/v3 v3.0.1
	github.com/sebitt27/dcrd/dcrec v1.0.1
	github.com/sebitt27/dcrd/dcrec/secp256k1/v4 v4.2.0
	github.com/sebitt27/dcrd/dcrjson/v4 v4.0.1
	github.com/sebitt27/dcrd/dcrutil/v4 v4.0.1
	github.com/sebitt27/dcrd/gcs/v4 v4.0.0
	github.com/sebitt27/dcrd/lru v1.1.2
	github.com/sebitt27/dcrd/math/uint256 v1.0.1
	github.com/sebitt27/dcrd/peer/v3 v3.0.2
	github.com/sebitt27/dcrd/rpc/jsonrpc/types/v4 v4.0.0
	github.com/sebitt27/dcrd/rpcclient/v8 v8.0.0
	github.com/sebitt27/dcrd/txscript/v4 v4.1.0
	github.com/sebitt27/dcrd/wire v1.6.0
	github.com/decred/dcrtest/dcrdtest v1.0.0
	github.com/decred/go-socks v1.1.0
	github.com/decred/slog v1.2.0
	github.com/gorilla/websocket v1.5.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/jrick/bitset v1.0.0
	github.com/jrick/logrotate v1.0.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	golang.org/x/sys v0.8.0
	golang.org/x/term v0.5.0
	lukechampine.com/blake3 v1.2.1
)

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/sebitt27/dcrd/dcrec/edwards/v2 v2.0.3 // indirect
	github.com/sebitt27/dcrd/hdkeychain/v3 v3.1.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
)

replace (
	github.com/sebitt27/dcrd/addrmgr/v2 => ./addrmgr
	github.com/sebitt27/dcrd/bech32 => ./bech32
	github.com/sebitt27/dcrd/blockchain/stake/v5 => ./blockchain/stake
	github.com/sebitt27/dcrd/blockchain/standalone/v2 => ./blockchain/standalone
	github.com/sebitt27/dcrd/blockchain/v5 => ./blockchain
	github.com/sebitt27/dcrd/certgen => ./certgen
	github.com/sebitt27/dcrd/chaincfg/chainhash => ./chaincfg/chainhash
	github.com/sebitt27/dcrd/chaincfg/v3 => ./chaincfg
	github.com/sebitt27/dcrd/connmgr/v3 => ./connmgr
	github.com/sebitt27/dcrd/container/apbf => ./container/apbf
	github.com/sebitt27/dcrd/crypto/blake256 => ./crypto/blake256
	github.com/sebitt27/dcrd/crypto/ripemd160 => ./crypto/ripemd160
	github.com/sebitt27/dcrd/database/v3 => ./database
	github.com/sebitt27/dcrd/dcrec => ./dcrec
	github.com/sebitt27/dcrd/dcrec/secp256k1/v4 => ./dcrec/secp256k1
	github.com/sebitt27/dcrd/dcrjson/v4 => ./dcrjson
	github.com/sebitt27/dcrd/dcrutil/v4 => ./dcrutil
	github.com/sebitt27/dcrd/gcs/v4 => ./gcs
	github.com/sebitt27/dcrd/hdkeychain/v3 => ./hdkeychain
	github.com/sebitt27/dcrd/limits => ./limits
	github.com/sebitt27/dcrd/lru => ./lru
	github.com/sebitt27/dcrd/math/uint256 => ./math/uint256
	github.com/sebitt27/dcrd/peer/v3 => ./peer
	github.com/sebitt27/dcrd/rpc/jsonrpc/types/v4 => ./rpc/jsonrpc/types
	github.com/sebitt27/dcrd/rpcclient/v8 => ./rpcclient
	github.com/sebitt27/dcrd/txscript/v4 => ./txscript
	github.com/sebitt27/dcrd/wire => ./wire
)
