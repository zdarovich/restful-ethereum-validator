package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/ethclient"
	"restful-ethereum-validator/beaconclient"
	"strconv"
)

type EthereumService interface {
	GetCurrentSlot(ctx context.Context) (uint64, error)
	GetBlockReward(ctx context.Context, slot uint64) (string, string, error)
	GetSyncDuties(ctx context.Context, slot uint64) ([]string, error)
}

type Reward struct {
	Status string `json:"status"`
	Reward string `json:"reward"`
}

type ethereumService struct {
	client       *ethclient.Client
	beaconClient beaconclient.Client
	rewardCache  *lru.Cache[uint64, Reward]
	dutiesCache  *lru.Cache[uint64, []string]
}

func NewEthereumService(beaconClient beaconclient.Client, client *ethclient.Client) EthereumService {
	return &ethereumService{
		client:       client,
		beaconClient: beaconClient,
		rewardCache:  lru.NewCache[uint64, Reward](100),
		dutiesCache:  lru.NewCache[uint64, []string](100),
	}
}
func (s *ethereumService) GetCurrentSlot(ctx context.Context) (uint64, error) {
	block, err := s.beaconClient.GetBeaconBlockInfo("head")
	if err != nil {
		return 0, err
	}
	if block == nil {
		return 0, err
	}
	slot, err := strconv.ParseUint(block.Data.Message.Slot, 10, 64)
	if err != nil {
		return 0, err
	}
	return slot, nil
}

func (s *ethereumService) GetBlockReward(ctx context.Context, slot uint64) (string, string, error) {
	if reward, ok := s.rewardCache.Get(slot); ok {
		return reward.Status, reward.Reward, nil
	}
	block, err := s.beaconClient.GetBeaconBlockInfo(strconv.Itoa(int(slot)))
	if err != nil {
		return "", "", err
	}
	if block == nil {
		return "", "", fmt.Errorf("block not found for slot %d", slot)
	}
	blockRewardInfo, err := s.beaconClient.GetBlockRewardInfo(slot)
	if err != nil {
		return "", "", err
	}
	feeRecipient := block.Data.Message.Body.ExecutionPayload.FeeRecipient
	code, err := s.client.CodeAt(ctx, common.HexToAddress(feeRecipient), nil)
	if err != nil {
		return "", "", err
	}
	rewardStatus := "Vanilla"
	// if address is not contract then it is MEV
	if len(code) == 0 {
		rewardStatus = "MEV"
	}
	s.rewardCache.Add(slot, Reward{Status: rewardStatus, Reward: blockRewardInfo.Data.Total})

	return rewardStatus, blockRewardInfo.Data.Total, nil
}

func (s *ethereumService) GetSyncDuties(ctx context.Context, slot uint64) ([]string, error) {
	if duties, ok := s.dutiesCache.Get(slot); ok {
		return duties, nil
	}
	block, err := s.beaconClient.GetBeaconBlockInfo(strconv.Itoa(int(slot)))
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("block not found for slot %d", slot)
	}
	syncComiteeRewards, err := s.beaconClient.GetSyncCommitteeRewards(slot)
	if err != nil {
		return nil, err
	}
	validatorIndexes := make([]string, len(syncComiteeRewards))
	for i, reward := range syncComiteeRewards {
		validatorIndexes[i] = reward.ValidatorIndex
	}

	vs, err := s.beaconClient.GetValidatorsByStatus(slot, validatorIndexes, []string{"active_ongoing"})
	if err != nil {
		return nil, err
	}
	validatorsPubKeys := make([]string, len(vs))
	for i, v := range vs {
		validatorsPubKeys[i] = v.ValidatorInformation.PubKey
	}
	s.dutiesCache.Add(slot, validatorsPubKeys)

	return validatorsPubKeys, nil
}
