package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"

	"blockchain-parser/internal/constant"
	"blockchain-parser/internal/entity"
	errorpkg "blockchain-parser/internal/error"
	"blockchain-parser/internal/service/mocks"
)

func TestParserWorker_getProcessingBlock(t *testing.T) {

	t.Run("get failed block", func(tt *testing.T) {
		ctx := context.Background()
		now := time.Now()
		failedBlock := entity.Block{
			Number:    1,
			Status:    constant.BlockStatusFailed,
			UpdatedAt: now,
		}
		processingBlock := entity.Block{
			Number:    1,
			Status:    constant.BlockStatusProcessing,
			UpdatedAt: now,
		}

		ctrl := gomock.NewController(nil)
		blockRepoMock := mocks.NewMockBlockRepository(ctrl)
		blockRepoMock.EXPECT().GetFailedBlock(ctx).Return(failedBlock, nil).Times(1)
		blockRepoMock.EXPECT().Upsert(ctx, processingBlock).Return(nil).Times(1)

		lockerMock := mocks.NewMockLocker(ctrl)
		lockerMock.EXPECT().Lock().Times(1)
		lockerMock.EXPECT().Unlock().Times(1)

		monkey.Patch(time.Now, func() time.Time {
			return now
		})

		w := NewParserWorker(
			nil,
			nil,
			blockRepoMock,
			nil,
			lockerMock,
		)

		block, err := w.getProcessingBlock(ctx, 0)
		if !errors.Is(err, nil) {
			t.Errorf("get failed block error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(block, processingBlock) {
			t.Errorf("get failed block got = %v, want %v", block, processingBlock)
		}
	})

	t.Run("get last block", func(tt *testing.T) {
		ctx := context.Background()
		now := time.Now()
		processingBlockFromBD := entity.Block{
			Number:    1,
			Status:    constant.BlockStatusProcessing,
			UpdatedAt: time.Now().Add(-time.Minute),
		}
		newProcessingBlock := entity.Block{
			Number:    2,
			Status:    constant.BlockStatusProcessing,
			UpdatedAt: now,
		}

		ctrl := gomock.NewController(nil)
		blockRepoMock := mocks.NewMockBlockRepository(ctrl)
		blockRepoMock.EXPECT().GetFailedBlock(ctx).Return(entity.Block{}, errorpkg.BlockNotFound).Times(1)
		blockRepoMock.EXPECT().GetLastBlock(ctx).Return(processingBlockFromBD, nil).Times(1)
		blockRepoMock.EXPECT().Upsert(ctx, newProcessingBlock).Return(nil).Times(1)

		lockerMock := mocks.NewMockLocker(ctrl)
		lockerMock.EXPECT().Lock().Times(1)
		lockerMock.EXPECT().Unlock().Times(1)

		monkey.Patch(time.Now, func() time.Time {
			return now
		})

		w := NewParserWorker(
			nil,
			nil,
			blockRepoMock,
			nil,
			lockerMock,
		)

		block, err := w.getProcessingBlock(ctx, 10)
		if !errors.Is(err, nil) {
			t.Errorf("get failed block error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(block, newProcessingBlock) {
			t.Errorf("get failed block got = %v, want %v", block, newProcessingBlock)
		}
	})

	t.Run("no block for parsing", func(tt *testing.T) {
		ctx := context.Background()
		now := time.Now()
		processingBlockFromBD := entity.Block{
			Number:    1,
			Status:    constant.BlockStatusProcessing,
			UpdatedAt: time.Now().Add(-time.Minute),
		}
		newProcessingBlock := entity.Block{
			Number:    2,
			Status:    constant.BlockStatusProcessing,
			UpdatedAt: now,
		}

		ctrl := gomock.NewController(nil)
		blockRepoMock := mocks.NewMockBlockRepository(ctrl)
		blockRepoMock.EXPECT().GetFailedBlock(ctx).Return(entity.Block{}, errorpkg.BlockNotFound).Times(1)
		blockRepoMock.EXPECT().GetLastBlock(ctx).Return(processingBlockFromBD, nil).Times(1)

		lockerMock := mocks.NewMockLocker(ctrl)
		lockerMock.EXPECT().Lock().Times(1)
		lockerMock.EXPECT().Unlock().Times(1)

		monkey.Patch(time.Now, func() time.Time {
			return now
		})

		w := NewParserWorker(
			nil,
			nil,
			blockRepoMock,
			nil,
			lockerMock,
		)

		block, err := w.getProcessingBlock(ctx, 1)
		if !errors.Is(err, errorpkg.NoBlockForParsing) {
			t.Errorf("get failed block error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(block, entity.Block{}) {
			t.Errorf("get failed block got = %v, want %v", block, newProcessingBlock)
		}
	})
}

func TestParserWorker_processBlock(t *testing.T) {

	t.Run("getting txns failed", func(tt *testing.T) {
		ctx := context.Background()
		block := entity.Block{
			Number: 1,
			Status: constant.BlockStatusProcessing,
		}
		gettingTxnsError := errors.New("error")

		ctrl := gomock.NewController(nil)
		blockChainClientMock := mocks.NewMockBlockChainClient(ctrl)
		blockChainClientMock.EXPECT().GetTxnsByBlockByNumber(ctx, block.Number).Return(nil, gettingTxnsError).Times(1)

		w := NewParserWorker(
			nil,
			nil,
			nil,
			blockChainClientMock,
			nil,
		)

		err := w.processBlock(ctx, block)
		if !errors.Is(err, gettingTxnsError) {
			t.Errorf("get process block error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("check subscription input address failed", func(tt *testing.T) {
		ctx := context.Background()
		block := entity.Block{
			Number: 34534,
			Status: constant.BlockStatusProcessing,
		}
		txns := []entity.Transaction{
			{
				From:             "0xd7def8de6bff40e7fa3a19b6749aca84bd5ba0ae",
				To:               "0x00000000006c3852cbef3e08e8df289169ede581",
				Value:            "0xb1a2bc2ec50000",
				BlockNumber:      34534,
				TransactionIndex: 0,
			},
			{
				From:             "0x069acf904f610cbf8ef1540349092852e46b4e95",
				To:               "0x8e9f0cd8f96e8e7b6531d01617e883d67f9dd150",
				Value:            "0x5f3bcfe512dcc00",
				BlockNumber:      34534,
				TransactionIndex: 1,
			},
			{
				From:             "0x4d7f1790644af787933c9ff0e2cff9a9b4299abb",
				To:               "0x417651e5e427b77fb1c258a9fbdbf4632fe348f9",
				Value:            "0x2386f26fc100000",
				BlockNumber:      34534,
				TransactionIndex: 2,
			},
		}
		failCheckSubscriptionErr := errors.New("error")

		ctrl := gomock.NewController(nil)
		blockChainClientMock := mocks.NewMockBlockChainClient(ctrl)
		blockChainClientMock.EXPECT().GetTxnsByBlockByNumber(ctx, block.Number).Return(txns, nil).Times(1)

		subscriptionRepoMock := mocks.NewMockSubscriberRepository(ctrl)
		subscriptionRepoMock.EXPECT().Get(ctx, "0x00000000006c3852cbef3e08e8df289169ede581").Return(entity.Subscriber{}, failCheckSubscriptionErr).Times(1)

		w := NewParserWorker(
			nil,
			subscriptionRepoMock,
			nil,
			blockChainClientMock,
			nil,
		)

		err := w.processBlock(ctx, block)
		if !errors.Is(err, failCheckSubscriptionErr) {
			t.Errorf("get process block error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("check subscription output address failed", func(tt *testing.T) {
		ctx := context.Background()
		block := entity.Block{
			Number: 34534,
			Status: constant.BlockStatusProcessing,
		}
		txns := []entity.Transaction{
			{
				From:             "0xd7def8de6bff40e7fa3a19b6749aca84bd5ba0ae",
				To:               "0x00000000006c3852cbef3e08e8df289169ede581",
				Value:            "0xb1a2bc2ec50000",
				BlockNumber:      34534,
				TransactionIndex: 0,
			},
			{
				From:             "0x069acf904f610cbf8ef1540349092852e46b4e95",
				To:               "0x8e9f0cd8f96e8e7b6531d01617e883d67f9dd150",
				Value:            "0x5f3bcfe512dcc00",
				BlockNumber:      34534,
				TransactionIndex: 1,
			},
			{
				From:             "0x4d7f1790644af787933c9ff0e2cff9a9b4299abb",
				To:               "0x417651e5e427b77fb1c258a9fbdbf4632fe348f9",
				Value:            "0x2386f26fc100000",
				BlockNumber:      34534,
				TransactionIndex: 2,
			},
		}
		failCheckSubscriptionErr := errors.New("error")

		ctrl := gomock.NewController(nil)
		blockChainClientMock := mocks.NewMockBlockChainClient(ctrl)
		blockChainClientMock.EXPECT().GetTxnsByBlockByNumber(ctx, block.Number).Return(txns, nil).Times(1)

		subscriptionRepoMock := mocks.NewMockSubscriberRepository(ctrl)
		subscriptionRepoMock.EXPECT().Get(ctx, "0x00000000006c3852cbef3e08e8df289169ede581").Return(entity.Subscriber{Address: "0x00000000006c3852cbef3e08e8df289169ede581"}, nil).Times(1)
		subscriptionRepoMock.EXPECT().Get(ctx, "0xd7def8de6bff40e7fa3a19b6749aca84bd5ba0ae").Return(entity.Subscriber{}, failCheckSubscriptionErr).Times(1)

		w := NewParserWorker(
			nil,
			subscriptionRepoMock,
			nil,
			blockChainClientMock,
			nil,
		)

		err := w.processBlock(ctx, block)
		if !errors.Is(err, failCheckSubscriptionErr) {
			t.Errorf("get process block error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("check subscription input and output", func(tt *testing.T) {
		ctx := context.Background()
		block := entity.Block{
			Number: 34534,
			Status: constant.BlockStatusProcessing,
		}

		txn1 := entity.Transaction{
			From:             "0xd7def8de6bff40e7fa3a19b6749aca84bd5ba0ae",
			To:               "0x00000000006c3852cbef3e08e8df289169ede581",
			Value:            "0xb1a2bc2ec50000",
			BlockNumber:      34534,
			TransactionIndex: 0,
		}
		txn2 := entity.Transaction{
			From:             "0x069acf904f610cbf8ef1540349092852e46b4e95",
			To:               "0x8e9f0cd8f96e8e7b6531d01617e883d67f9dd150",
			Value:            "0x5f3bcfe512dcc00",
			BlockNumber:      34534,
			TransactionIndex: 1,
		}
		txn3 := entity.Transaction{
			From:             "0x4d7f1790644af787933c9ff0e2cff9a9b4299abb",
			To:               "0x417651e5e427b77fb1c258a9fbdbf4632fe348f9",
			Value:            "0x2386f26fc100000",
			BlockNumber:      34534,
			TransactionIndex: 2,
		}

		txns := []entity.Transaction{txn1, txn2, txn3}

		ctrl := gomock.NewController(nil)
		blockChainClientMock := mocks.NewMockBlockChainClient(ctrl)
		blockChainClientMock.EXPECT().GetTxnsByBlockByNumber(ctx, block.Number).Return(txns, nil).Times(1)

		subscriptionRepoMock := mocks.NewMockSubscriberRepository(ctrl)
		subscriptionRepoMock.EXPECT().Get(ctx, txn1.To).Return(entity.Subscriber{Address: txn1.To}, nil).Times(1)
		subscriptionRepoMock.EXPECT().Get(ctx, txn1.From).Return(entity.Subscriber{}, errorpkg.SubscriberNotFound).Times(1)

		subscriptionRepoMock.EXPECT().Get(ctx, txn2.To).Return(entity.Subscriber{}, errorpkg.SubscriberNotFound).Times(1)
		subscriptionRepoMock.EXPECT().Get(ctx, txn2.From).Return(entity.Subscriber{Address: txn2.From}, nil).Times(1)

		subscriptionRepoMock.EXPECT().Get(ctx, txn3.To).Return(entity.Subscriber{}, errorpkg.SubscriberNotFound).Times(1)
		subscriptionRepoMock.EXPECT().Get(ctx, txn3.From).Return(entity.Subscriber{}, errorpkg.SubscriberNotFound).Times(1)

		txnRepoMock := mocks.NewMockTransactionRepository(ctrl)
		txnRepoMock.EXPECT().Save(ctx, txn1).Return(nil).Times(1)
		txnRepoMock.EXPECT().Save(ctx, txn2).Return(nil).Times(1)

		w := NewParserWorker(
			txnRepoMock,
			subscriptionRepoMock,
			nil,
			blockChainClientMock,
			nil,
		)

		err := w.processBlock(ctx, block)
		if !errors.Is(err, nil) {
			t.Errorf("get process block error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("saving transaction failed", func(tt *testing.T) {
		ctx := context.Background()
		block := entity.Block{
			Number: 34534,
			Status: constant.BlockStatusProcessing,
		}

		txn1 := entity.Transaction{
			From:             "0xd7def8de6bff40e7fa3a19b6749aca84bd5ba0ae",
			To:               "0x00000000006c3852cbef3e08e8df289169ede581",
			Value:            "0xb1a2bc2ec50000",
			BlockNumber:      34534,
			TransactionIndex: 0,
		}

		txns := []entity.Transaction{txn1}
		savingTxnErr := errors.New("error")

		ctrl := gomock.NewController(nil)
		blockChainClientMock := mocks.NewMockBlockChainClient(ctrl)
		blockChainClientMock.EXPECT().GetTxnsByBlockByNumber(ctx, block.Number).Return(txns, nil).Times(1)

		subscriptionRepoMock := mocks.NewMockSubscriberRepository(ctrl)
		subscriptionRepoMock.EXPECT().Get(ctx, txn1.To).Return(entity.Subscriber{Address: txn1.To}, nil).Times(1)
		subscriptionRepoMock.EXPECT().Get(ctx, txn1.From).Return(entity.Subscriber{}, errorpkg.SubscriberNotFound).Times(1)

		txnRepoMock := mocks.NewMockTransactionRepository(ctrl)
		txnRepoMock.EXPECT().Save(ctx, txn1).Return(savingTxnErr).Times(1)

		w := NewParserWorker(
			txnRepoMock,
			subscriptionRepoMock,
			nil,
			blockChainClientMock,
			nil,
		)

		err := w.processBlock(ctx, block)
		if !errors.Is(err, savingTxnErr) {
			t.Errorf("get process block error = %v, wantErr %v", err, nil)
			return
		}
	})
}
