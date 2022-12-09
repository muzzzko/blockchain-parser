package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BlockChainParser struct {
	parser Parser
}

func NewBlockChainParser(parser Parser) *BlockChainParser {
	return &BlockChainParser{
		parser: parser,
	}
}

func (h *BlockChainParser) GetCurrentBlock(w http.ResponseWriter, _ *http.Request) {
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

func (h *BlockChainParser) Subscribe(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("address %s was subscribed", blockChainParserSubscribe.Address)

	w.WriteHeader(http.StatusNoContent)
}

func (h *BlockChainParser) GetTransactions(w http.ResponseWriter, r *http.Request) {
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
