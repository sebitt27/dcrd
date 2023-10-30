// Copyright (c) 2019-2022 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"

	"github.com/sebitt27/dcrd/blockchain/standalone/v2"
	"github.com/sebitt27/dcrd/chaincfg/chainhash"
	"github.com/sebitt27/dcrd/database/v3"
	"github.com/sebitt27/dcrd/dcrutil/v4"
	"github.com/sebitt27/dcrd/gcs/v4"
	"github.com/sebitt27/dcrd/wire"
)

const (
	// HeaderCmtFilterIndex is the proof index for the filter header commitment.
	HeaderCmtFilterIndex = 0
)

// headerCommitmentData houses information the block header commits to via the
// commitment root.
type headerCommitmentData struct {
	filter     *gcs.FilterV2
	filterHash chainhash.Hash
}

// v1Leaves returns the individual commitment hashes that comprise the leaves of
// the merkle tree for a v1 header commitment.
func (c *headerCommitmentData) v1Leaves() []chainhash.Hash {
	return []chainhash.Hash{c.filterHash}
}

// CalcCommitmentRootV1 calculates and returns the required v1 block commitment
// root from the filter hash it commits to.
//
// This function is safe for concurrent access.
func CalcCommitmentRootV1(filterHash chainhash.Hash) chainhash.Hash {
	// NOTE: The commitment root is actually the merkle root of a merkle tree
	// whose leaves are each of the individual commitments.  However, since
	// there is only currently a single commitment, the merkle root will simply
	// be the hash of the sole item, so there is no point in doing extra work.
	//
	// The callers could certainly simply avoid calling this function and do the
	// same thing directly, however, providing versioned functions helps make
	// it clear exactly what each header commitment version commits to and makes
	// the code more consistent with multiple versions.
	return filterHash
}

// FetchUtxoViewParentTemplate loads utxo details from the point of view of just
// having connected the given block, which must be a block template that
// connects to the parent of the tip of the main chain.  In other words, the
// given block must be a sibling of the current tip of the main chain.
//
// This should typically only be used by mining code when it is unable to
// generate a template that extends the current tip due to being unable to
// acquire the minimum required number of votes to extend it.
//
// This function is safe for concurrent access however the returned view is NOT.
func (b *BlockChain) FetchUtxoViewParentTemplate(block *wire.MsgBlock) (*UtxoViewpoint, error) {
	b.chainLock.Lock()
	defer b.chainLock.Unlock()

	// The block template must build off the parent of the current tip of the
	// main chain.
	tip := b.bestChain.Tip()
	if tip.parent == nil {
		str := fmt.Sprintf("unable to fetch utxos for non-existent parent of "+
			"the current tip %s", tip.hash)
		return nil, ruleError(ErrInvalidTemplateParent, str)
	}
	parentHash := block.Header.PrevBlock
	if parentHash != tip.parent.hash {
		str := fmt.Sprintf("previous block must be the parent of the current "+
			"chain tip %s, but got %s", tip.parent.hash, parentHash)
		return nil, ruleError(ErrInvalidTemplateParent, str)
	}

	// Since the block template is building on the parent of the current tip,
	// undo the transactions and spend information for the tip block to reach
	// the point of view of the block template.
	view := NewUtxoViewpoint(b.utxoCache)
	view.SetBestHash(&tip.hash)
	tipBlock, err := b.fetchMainChainBlockByNode(tip)
	if err != nil {
		return nil, err
	}
	parent, err := b.fetchMainChainBlockByNode(tip.parent)
	if err != nil {
		return nil, err
	}

	// Determine if treasury agenda is active.
	isTreasuryEnabled, err := b.isTreasuryAgendaActive(tip.parent)
	if err != nil {
		return nil, err
	}

	// Load all of the spent txos for the tip block from the spend journal.
	var stxos []spentTxOut
	err = b.db.View(func(dbTx database.Tx) error {
		stxos, err = dbFetchSpendJournalEntry(dbTx, tipBlock, isTreasuryEnabled)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Update the view to unspend all of the spent txos and remove the utxos
	// created by the tip block.  Also, if the block votes against its parent,
	// reconnect all of the regular transactions.
	err = view.disconnectBlock(tipBlock, parent, stxos, isTreasuryEnabled)
	if err != nil {
		return nil, err
	}

	// The view is now from the point of view of the parent of the current tip
	// block.  However, calculating the commitment root requires the view to
	// include outputs created in the candidate block, so update the view to
	// mark all utxos referenced by the block as spent and add all transactions
	// being created by the block to it.  In the case the block votes against
	// the parent, also disconnect all of the regular transactions in the parent
	// block.
	utilBlock := dcrutil.NewBlock(block)
	err = view.connectBlock(b.db, utilBlock, parent, nil, isTreasuryEnabled)
	if err != nil {
		return nil, err
	}

	return view, nil
}

// HeaderProof houses a merkle tree inclusion proof and associated proof index
// for a header commitment.  This information allows clients to efficiently
// prove whether or not the commitment root of a block header commits to
// specific data at the given index.
type HeaderProof struct {
	ProofIndex  uint32
	ProofHashes []chainhash.Hash
}

// FilterByBlockHash returns the version 2 GCS filter for the given block hash
// along with a header commitment inclusion proof when they exist.  This
// function returns the filters regardless of whether or not their associated
// block is part of the main chain.
//
// An error that wraps ErrNoFilter will be returned when the filter for the
// given block hash does not exist.
//
// This function is safe for concurrent access.
func (b *BlockChain) FilterByBlockHash(hash *chainhash.Hash) (*gcs.FilterV2, *HeaderProof, error) {
	// Avoid a database lookup when there is no way the filter data for the
	// requested block is available.
	node := b.index.LookupNode(hash)
	if node == nil || !b.index.NodeStatus(node).HaveData() {
		str := fmt.Sprintf("no filter available for block %s", hash)
		return nil, nil, contextError(ErrNoFilter, str)
	}

	// Attempt to load the filter and associated header commitments from the
	// database.
	var filter *gcs.FilterV2
	var leaves []chainhash.Hash
	err := b.db.View(func(dbTx database.Tx) error {
		var err error
		filter, err = dbFetchGCSFilter(dbTx, hash)
		if err != nil {
			return err
		}
		if filter == nil {
			str := fmt.Sprintf("no filter available for block %s", hash)
			return contextError(ErrNoFilter, str)
		}

		leaves, err = dbFetchHeaderCommitments(dbTx, hash)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	// Generate the header commitment inclusion proof for the filter.
	const proofIndex = HeaderCmtFilterIndex
	proof := standalone.GenerateInclusionProof(leaves, proofIndex)
	headerProof := &HeaderProof{
		ProofIndex:  proofIndex,
		ProofHashes: proof,
	}
	return filter, headerProof, nil
}
