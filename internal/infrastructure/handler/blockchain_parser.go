package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BlockChanParser struct {
	parser Parser
}

func NewBlockChanParser(parser Parser) *BlockChanParser {
	return &BlockChanParser{
		parser: parser,
	}
}

func (h *BlockChanParser) GetCurrentBlock(w http.ResponseWriter, _ *http.Request) {
	blockNumber := h.parser.GetCurrentBlock()
	if blockNumber == 0 {
		resp := ErrorResponse{
			Message: "fail get block",
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(resp)

		return
	}

	resp := blockChainParserGetCurrentBlockResponse{
		Block: fmt.Sprintf("0x%x", blockNumber),
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *BlockChanParser) Subscribe(w http.ResponseWriter, r *http.Request) {
	blockChainParserSubscribe := BlockChainParserSubscribe{}
	if err := json.NewDecoder(r.Body).Decode(&blockChainParserSubscribe); err != nil {
		resp := ErrorResponse{
			Message: fmt.Sprintf("fail decode request: %s", err),
		}

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)

		return
	}

	if ok := h.parser.Subscribe(blockChainParserSubscribe.Address); !ok {
		resp := ErrorResponse{
			Message: "fail subscribe",
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(resp)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BlockChanParser) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		resp := ErrorResponse{
			Message: "address is required",
		}

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)

		return
	}

	txns := h.parser.GetTransactions(address)

	w.WriteHeader(http.StatusOK)
	resp := mapTransactionsToGetTransactionsResponse(txns)
	_ = json.NewEncoder(w).Encode(resp)
}
