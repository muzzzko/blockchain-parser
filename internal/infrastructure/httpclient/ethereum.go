package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"blockchain-parser/config"
	"blockchain-parser/internal/entity"
	errorpkg "blockchain-parser/internal/error"
)

const (
	ethJSONRPCVersion = "2.0"

	ethGetBlockNumberMethod = "eth_blockNumber"
	ethGetBlockByNumber     = "eth_getBlockByNumber"
)

type Ethereum struct {
	clnt http.Client
	cfg  config.EthereumHttpClient
}

func NewEthereum(cfg config.EthereumHttpClient) *Ethereum {
	clnt := http.Client{
		Timeout: cfg.Timeout,
	}

	return &Ethereum{
		clnt: clnt,
		cfg:  cfg,
	}
}

func (c *Ethereum) GetBlockNumber(ctx context.Context) (int, error) {
	id := rand.Int31()
	body := ethereumRequestBody{
		Version: ethJSONRPCVersion,
		Method:  ethGetBlockNumberMethod,
		Params:  []interface{}{},
		ID:      id,
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return 0, fmt.Errorf("fail marshal body in GetBlock: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.cfg.Host, &buf)
	if err != nil {
		return 0, fmt.Errorf("fail create request in GetBlock: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.clnt.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if os.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("fail get block number in GetBlock: %w", errorpkg.TimeoutErr)
		}

		return 0, fmt.Errorf("fail get block number in GetBlock: %w", err)
	}

	ethGetBlockNumberResponse := ethereumGetBlockNumberResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&ethGetBlockNumberResponse); err != nil {
		return 0, fmt.Errorf("fail unmarshal response in GetBlock: %w", err)
	}

	if ethGetBlockNumberResponse.Error != nil {
		return 0, fmt.Errorf("error code (%d), message (%s) in GetBlock: %w", ethGetBlockNumberResponse.Error.Code, ethGetBlockNumberResponse.Error.Message, errorpkg.HttpErr)
	}

	if ethGetBlockNumberResponse.Result == nil {
		return 0, errors.New("result is nil in GetBlock response")
	}

	if ethGetBlockNumberResponse.ID != id {
		return 0, errors.New("mismatch request and response IDs in GetBlock")
	}

	blockNumber, err := strconv.ParseInt(*ethGetBlockNumberResponse.Result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("parse block number in GetBlock response: %w", err)
	}

	return int(blockNumber), nil
}

func (c *Ethereum) GetTxnsByBlockByNumber(ctx context.Context, blockNumber int) ([]entity.Transaction, error) {
	id := rand.Int31()
	body := ethereumRequestBody{
		Version: ethJSONRPCVersion,
		Method:  ethGetBlockByNumber,
		Params: []interface{}{
			fmt.Sprintf("0x%x", blockNumber),
			true,
		},
		ID: id,
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, fmt.Errorf("fail marshal body in GetTxnsByBlockByNumber: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.cfg.Host, &buf)
	if err != nil {
		return nil, fmt.Errorf("fail create request in GetTxnsByBlockByNumber: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.clnt.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if os.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("fail get transaction in GetTxnsByBlockByNumber: %w", errorpkg.TimeoutErr)
		}

		return nil, fmt.Errorf("fail get transaction GetTxnsByBlockByNumber: %w", err)
	}

	ethGetBlockByNumberResponse := ethereumGetBlockByNumberResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&ethGetBlockByNumberResponse); err != nil {
		return nil, fmt.Errorf("fail unmarshal response in GetTxnsByBlockByNumber: %w", err)
	}

	if ethGetBlockByNumberResponse.Error != nil {
		return nil, fmt.Errorf("error code (%d), message (%s) in GetTxnsByBlockByNumber: %w", ethGetBlockByNumberResponse.Error.Code, ethGetBlockByNumberResponse.Error.Message, errorpkg.HttpErr)
	}

	if ethGetBlockByNumberResponse.ID != id {
		return nil, errors.New("mismatch request and response IDs in GetTxnsByBlockByNumber")
	}

	return mapResponseToTxns(ethGetBlockByNumberResponse), nil
}
