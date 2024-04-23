package beaconclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client interface {
	GetSyncCommitteeRewards(slot uint64) ([]SyncCommitteeReward, error)
	GetValidatorsByStatus(slot uint64, ids []string, statuses []string) ([]Validator, error)
	GetBlockRewardInfo(slot uint64) (*BlockRewardResponse, error)
	GetBeaconBlockInfo(slot string) (*BeaconBlockResponse, error)
}

type beaconClient struct {
	baseURL string
	client  *http.Client
}

func NewBeaconClient(baseURL string) Client {
	return &beaconClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (vc *beaconClient) GetSyncCommitteeRewards(slot uint64) ([]SyncCommitteeReward, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/rewards/sync_committee/%d", vc.baseURL, slot)

	reqBody := []byte(`[]`)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := vc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	var rewardsResponse RewardsResponse
	err = json.NewDecoder(resp.Body).Decode(&rewardsResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return rewardsResponse.Data, nil
}

func (vc *beaconClient) GetValidatorsByStatus(slot uint64, ids []string, statuses []string) ([]Validator, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/states/%d/validators", vc.baseURL, slot)

	payload := map[string][]string{
		"ids":      ids,
		"statuses": statuses,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := vc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	var validatorResponse ValidatorResponse
	err = json.NewDecoder(resp.Body).Decode(&validatorResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return validatorResponse.Data, nil
}

func (vc *beaconClient) GetBlockRewardInfo(slot uint64) (*BlockRewardResponse, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/rewards/blocks/%d", vc.baseURL, slot)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	resp, err := vc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	var blockRewardResponse BlockRewardResponse
	err = json.NewDecoder(resp.Body).Decode(&blockRewardResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &blockRewardResponse, nil
}

func (vc *beaconClient) GetBeaconBlockInfo(slot string) (*BeaconBlockResponse, error) {
	url := fmt.Sprintf("%s/eth/v2/beacon/blocks/%s", vc.baseURL, slot)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	resp, err := vc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	var beaconBlockResponse BeaconBlockResponse
	err = json.NewDecoder(resp.Body).Decode(&beaconBlockResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &beaconBlockResponse, nil
}

type BlockRewardResponse struct {
	ExecutionOptimistic bool            `json:"execution_optimistic"`
	Finalized           bool            `json:"finalized"`
	Data                BlockRewardData `json:"data"`
}

type BlockRewardData struct {
	ProposerIndex     string `json:"proposer_index"`
	Total             string `json:"total"`
	Attestations      string `json:"attestations"`
	SyncAggregate     string `json:"sync_aggregate"`
	ProposerSlashings string `json:"proposer_slashings"`
	AttesterSlashings string `json:"attester_slashings"`
}

type RewardsResponse struct {
	ExecutionOptimistic bool                  `json:"execution_optimistic"`
	Finalized           bool                  `json:"finalized"`
	Data                []SyncCommitteeReward `json:"data"`
}

type SyncCommitteeReward struct {
	ValidatorIndex string `json:"validator_index"`
	Reward         string `json:"reward"`
}

type ValidatorResponse struct {
	ExecutionOptimistic bool        `json:"execution_optimistic"`
	Finalized           bool        `json:"finalized"`
	Data                []Validator `json:"data"`
}

type Validator struct {
	Index                string               `json:"index"`
	Balance              string               `json:"balance"`
	Status               string               `json:"status"`
	ValidatorInformation ValidatorInformation `json:"validator"`
}

type ValidatorInformation struct {
	PubKey                     string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

type BeaconBlockResponse struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
	Finalized           bool   `json:"finalized"`
	Data                struct {
		Message struct {
			Slot          string `json:"slot"`
			ProposerIndex string `json:"proposer_index"`
			ParentRoot    string `json:"parent_root"`
			StateRoot     string `json:"state_root"`
			Body          struct {
				RandaoReveal string `json:"randao_reveal"`
				Eth1Data     struct {
					DepositRoot  string `json:"deposit_root"`
					DepositCount string `json:"deposit_count"`
					BlockHash    string `json:"block_hash"`
				} `json:"eth1_data"`
				Graffiti          string        `json:"graffiti"`
				ProposerSlashings []interface{} `json:"proposer_slashings"`
				AttesterSlashings []interface{} `json:"attester_slashings"`
				Attestations      []struct {
					AggregationBits string `json:"aggregation_bits"`
					Data            struct {
						Slot            string `json:"slot"`
						Index           string `json:"index"`
						BeaconBlockRoot string `json:"beacon_block_root"`
						Source          struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"source"`
						Target struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"target"`
					} `json:"data"`
					Signature string `json:"signature"`
				} `json:"attestations"`
				Deposits []struct {
					Proof []string `json:"proof"`
					Data  struct {
						Pubkey                string `json:"pubkey"`
						WithdrawalCredentials string `json:"withdrawal_credentials"`
						Amount                string `json:"amount"`
						Signature             string `json:"signature"`
					} `json:"data"`
				} `json:"deposits"`
				VoluntaryExits []interface{} `json:"voluntary_exits"`
				SyncAggregate  struct {
					SyncCommitteeBits      string `json:"sync_committee_bits"`
					SyncCommitteeSignature string `json:"sync_committee_signature"`
				} `json:"sync_aggregate"`
				ExecutionPayload struct {
					ParentHash    string   `json:"parent_hash"`
					FeeRecipient  string   `json:"fee_recipient"`
					StateRoot     string   `json:"state_root"`
					ReceiptsRoot  string   `json:"receipts_root"`
					LogsBloom     string   `json:"logs_bloom"`
					PrevRandao    string   `json:"prev_randao"`
					BlockNumber   string   `json:"block_number"`
					GasLimit      string   `json:"gas_limit"`
					GasUsed       string   `json:"gas_used"`
					Timestamp     string   `json:"timestamp"`
					ExtraData     string   `json:"extra_data"`
					BaseFeePerGas string   `json:"base_fee_per_gas"`
					BlockHash     string   `json:"block_hash"`
					Transactions  []string `json:"transactions"`
					Withdrawals   []struct {
						Index          string `json:"index"`
						ValidatorIndex string `json:"validator_index"`
						Address        string `json:"address"`
						Amount         string `json:"amount"`
					} `json:"withdrawals"`
					BlobGasUsed   string `json:"blob_gas_used"`
					ExcessBlobGas string `json:"excess_blob_gas"`
				} `json:"execution_payload"`
				BlsToExecutionChanges []interface{} `json:"bls_to_execution_changes"`
				BlobKzgCommitments    []interface{} `json:"blob_kzg_commitments"`
			} `json:"body"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"data"`
}
