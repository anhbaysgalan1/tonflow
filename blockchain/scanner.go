package blockchain

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"tonflow/pkg"
)

func Scan(c *Client, txChan chan<- *tlb.Transaction, errCh chan<- error) {
	ctx := c.liteClient.StickyContext(context.Background())

	// get the latest block of master chain
	master, err := c.tonClient.GetMasterchainInfo(ctx)
	if err != nil {
		errCh <- fmt.Errorf("failed to get masterchain info: %w", err)
		return
	}
	log.Debugf("masterchain info: %s", pkg.PrintAny(master))

	// getting information about other work-chains and its shards of first master block
	// to init storage of last seen shard seq numbers
	firstShards, err := c.tonClient.GetBlockShardsInfo(ctx, master)
	if err != nil {
		errCh <- fmt.Errorf("failed to get work-chains and shards of first master block: %w", err)
		return
	}
	log.Debugf("all workchains and its shards 1: %s", pkg.PrintAny(master))

	// storage for last seen shard seqno
	// TODO: load from DB. So needs to save somewhere in code
	shardLastSeqno := map[string]uint32{}

	// save shard workchain | shard and seqno
	for _, shard := range firstShards {
		shardLastSeqno[getShardID(shard)] = shard.SeqNo
	}

	log.Debugf("shardLastSeqno 1: %s", pkg.PrintAny(shardLastSeqno))

	for {
		log.Debugf("scanning %d master block ...\n", master.SeqNo)

		// getting information about other work-chains and shards of master block
		currentShards, err := c.tonClient.GetBlockShardsInfo(ctx, master)
		if err != nil {
			errCh <- fmt.Errorf("failed to get other work-chains and shards of master block: %w", err)
			return
		}
		log.Debugf("all workchains and its shards 2: %s", pkg.PrintAny(master))

		// shards in master block may have holes, e.g. shard seqno 2756461, then 2756463, and no 2756462 in master chain
		// thus we need to scan a bit back in case of discovering a hole, till last seen, to fill the misses.
		var newShards []*tlb.BlockInfo
		for _, shard := range currentShards {
			notSeen, err := getNotSeenShards(ctx, c.tonClient, shard, shardLastSeqno)
			if err != nil {
				errCh <- fmt.Errorf("failed to get not seen shards: %w", err)
				return
			}
			shardLastSeqno[getShardID(shard)] = shard.SeqNo
			log.Debugf("shardLastSeqno 2: %s", pkg.PrintAny(shardLastSeqno))

			newShards = append(newShards, notSeen...)
		}

		var txList []*tlb.Transaction

		// for each shard block getting transactions
		for _, shard := range newShards {
			log.Debugf("scanning block %d of shard %d ...", shard.SeqNo, shard.Shard)

			var fetchedIDs []*tlb.TransactionID
			var after *tlb.TransactionID
			var more = true

			// load all transactions in batches with 100 transactions in each while exists
			for more {
				fetchedIDs, more, err = c.tonClient.GetBlockTransactions(ctx, shard, 100, after)
				if err != nil {
					errCh <- fmt.Errorf("failed to get tx ids: %w", err)
					return
				}

				if more {
					// set load offset for next query (pagination)
					after = fetchedIDs[len(fetchedIDs)-1]
				}

				for _, id := range fetchedIDs {
					// get full transaction by id
					tx, err := c.tonClient.GetTransaction(ctx, shard, address.NewAddress(0, 0, id.AccountID), id.LT)
					if err != nil {
						errCh <- fmt.Errorf("failed to get tx data: %w", err)
						return
					}
					txList = append(txList, tx)
				}
			}
		}

		for _, transaction := range txList {
			// log.Debug("TRANSACTION:", i, transaction.String())
			txChan <- transaction
		}

		if len(txList) == 0 {
			log.Debugf("no transactions in %d block", master.SeqNo)
		}

		master, err = c.tonClient.WaitNextMasterBlock(ctx, master)
		if err != nil {
			errCh <- fmt.Errorf("failed to wait next master: %w", err)
			return
		}
	}
}

// func to get storage map key
func getShardID(shard *tlb.BlockInfo) string {
	return fmt.Sprintf("%d|%d", shard.Workchain, shard.Shard)
}

func getNotSeenShards(ctx context.Context, api *ton.APIClient, shard *tlb.BlockInfo, shardLastSeqno map[string]uint32) (ret []*tlb.BlockInfo, err error) {
	if no, ok := shardLastSeqno[getShardID(shard)]; ok && no == shard.SeqNo {
		return nil, nil
	}

	b, err := api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, fmt.Errorf("failed to get block data: %w", err)
	}

	parents, err := b.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, fmt.Errorf("failed to get parent blocks (%d:%x:%d): %w", shard.Workchain, uint64(shard.Shard), shard.Shard, err)
	}

	for _, parent := range parents {
		ext, err := getNotSeenShards(ctx, api, parent, shardLastSeqno)
		if err != nil {
			return nil, fmt.Errorf("failed to get not seen shards into func getNotSeenShards: %w", err)
		}
		ret = append(ret, ext...)
	}

	ret = append(ret, shard)
	return ret, nil
}
