package blockchain

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"tonflow/pkg"
	"tonflow/storage"
)

func getShardID(shard *tlb.BlockInfo) string {
	return fmt.Sprintf("%d|%d", shard.Workchain, shard.Shard)
}

func Scan(c *Client, storage storage.Storage, txChan chan<- *tlb.Transaction, errCh chan<- error) {
	ctx := c.liteClient.StickyContext(context.Background())

	shardLastSeqno, err := storage.GetLastSeqno(ctx)
	if err != nil {
		errCh <- fmt.Errorf("failed to get shard last seqno from storage: %w", err)
		return
	}

	master, err := c.tonClient.GetMasterchainInfo(ctx)
	if err != nil {
		errCh <- fmt.Errorf("failed to get master chain latest block: %w", err)
		return
	}

	if len(shardLastSeqno) == 0 {
		workchainsShards, err := c.tonClient.GetBlockShardsInfo(ctx, master)
		if err != nil {
			errCh <- fmt.Errorf("failed to get workchains and its shards of master block: %w", err)
			return
		}

		for _, shard := range workchainsShards {
			shardLastSeqno[getShardID(shard)] = shard.SeqNo
		}
	}

	for {
		currentShards, err := c.tonClient.GetBlockShardsInfo(ctx, master)
		if err != nil {
			errCh <- fmt.Errorf("failed to get workchains and its shards of master block: %w", err)
			return
		}

		var newShards []*tlb.BlockInfo
		for _, shard := range currentShards {
			notSeen, err := getNotSeenShards(ctx, c.tonClient, shard, shardLastSeqno)
			if err != nil {
				errCh <- fmt.Errorf("failed to get not seen shards: %w", err)
				return
			}

			shardLastSeqno[getShardID(shard)] = shard.SeqNo
			newShards = append(newShards, notSeen...)
		}

		log.Debugf("shardLastSeqno:\n%s", pkg.PrintAny(shardLastSeqno))

		err = storage.SetLastSeqno(context.Background(), shardLastSeqno)
		if err != nil {
			log.Errorf("failed to set shard last seqno into storage: %s", err)
		}

		var txList []*tlb.Transaction

		for _, shard := range newShards {
			var fetchedIDs []*tlb.TransactionID
			var after *tlb.TransactionID
			var more = true

			for more {
				fetchedIDs, more, err = c.tonClient.GetBlockTransactions(ctx, shard, 100, after)
				if err != nil {
					errCh <- fmt.Errorf("failed to get tx ids: %w", err)
					return
				}

				if more {
					after = fetchedIDs[len(fetchedIDs)-1]
				}

				for _, id := range fetchedIDs {
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
			txChan <- transaction
		}

		if len(txList) == 0 {
			log.Debugf("no transactions in %d master block", master.SeqNo)
		}

		master, err = c.tonClient.WaitNextMasterBlock(ctx, master)
		if err != nil {
			errCh <- fmt.Errorf("failed to wait next master block: %w", err)
			return
		}
	}
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
