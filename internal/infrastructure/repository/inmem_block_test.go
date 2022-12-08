package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"blockchain-parser/internal/constant"
	"blockchain-parser/internal/entity"
	errorpkg "blockchain-parser/internal/error"
)

func TestInMemBlock_GetLastParsedBlock(t *testing.T) {
	type fields struct {
		parsedBlocks      map[int]entity.Block
		parsedBlockNumber int
	}
	type args struct {
		ctx context.Context
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Block
		wantErr error
	}{
		{
			name: "block NOT found",
			fields: fields{
				parsedBlocks:      map[int]entity.Block{},
				parsedBlockNumber: 1,
			},
			args: args{
				ctx: ctx,
			},
			want:    entity.Block{},
			wantErr: errorpkg.BlockNotFound,
		},
		{
			name: "block found",
			fields: fields{
				parsedBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusParsed,
					},
				},
				parsedBlockNumber: 1,
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Block{
				Number: 1,
				Status: constant.BlockStatusParsed,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InMemBlock{
				parsedBlocks:      tt.fields.parsedBlocks,
				parsedBlockNumber: tt.fields.parsedBlockNumber,
			}
			got, err := r.GetLastParsedBlock(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetLastParsedBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastParsedBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemBlock_GetLastBlock(t *testing.T) {
	type fields struct {
		processingBlocks      map[int]entity.Block
		parsedBlocks          map[int]entity.Block
		parsedBlockNumber     int
		processingBlockNumber int
	}
	type args struct {
		ctx context.Context
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Block
		wantErr error
	}{
		{
			name: "block NOT found",
			fields: fields{
				parsedBlocks:          map[int]entity.Block{},
				processingBlocks:      map[int]entity.Block{},
				parsedBlockNumber:     1,
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
			},
			want:    entity.Block{},
			wantErr: errorpkg.BlockNotFound,
		},
		{
			name: "return processing block",
			fields: fields{
				parsedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusProcessing,
					},
				},
				parsedBlockNumber:     1,
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Block{
				Number: 2,
				Status: constant.BlockStatusProcessing,
			},
			wantErr: nil,
		},
		{
			name: "return parsed block",
			fields: fields{
				parsedBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusParsed,
					},
				},
				processingBlocks:      map[int]entity.Block{},
				parsedBlockNumber:     1,
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Block{
				Number: 1,
				Status: constant.BlockStatusParsed,
			},
			wantErr: nil,
		},
		{
			name: "return last parsed block, there is processing block though",
			fields: fields{
				parsedBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusParsed,
					},
				},
				processingBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusProcessing,
					},
				},
				parsedBlockNumber:     2,
				processingBlockNumber: 1,
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Block{
				Number: 2,
				Status: constant.BlockStatusParsed,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemBlock{
				processingBlocks:      tt.fields.processingBlocks,
				parsedBlocks:          tt.fields.parsedBlocks,
				parsedBlockNumber:     tt.fields.parsedBlockNumber,
				processingBlockNumber: tt.fields.processingBlockNumber,
			}
			got, err := r.GetLastBlock(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetLastBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemBlock_GetFailedBlock(t *testing.T) {
	type fields struct {
		failedBlocks     map[int]entity.Block
		processingBlocks map[int]entity.Block
	}

	ctx := context.Background()
	expiredUpdatedAt := time.Now().Add(-time.Hour)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Block
		wantErr error
	}{
		{
			name: "block NOT found (processing block is NOT expired)",
			fields: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					2: {
						Number:    2,
						Status:    constant.BlockStatusProcessing,
						UpdatedAt: time.Now(),
					},
				},
			},
			args: args{
				ctx: ctx,
			},
			want:    entity.Block{},
			wantErr: errorpkg.BlockNotFound,
		},
		{
			name: "return processing block",
			fields: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					2: {
						Number:    2,
						Status:    constant.BlockStatusProcessing,
						UpdatedAt: expiredUpdatedAt,
					},
				},
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Block{
				Number:    2,
				Status:    constant.BlockStatusProcessing,
				UpdatedAt: expiredUpdatedAt,
			},
			wantErr: nil,
		},
		{
			name: "return failed block",
			fields: fields{
				failedBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusFailed,
					},
				},
				processingBlocks: map[int]entity.Block{},
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Block{
				Number: 2,
				Status: constant.BlockStatusFailed,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemBlock{
				failedBlocks:     tt.fields.failedBlocks,
				processingBlocks: tt.fields.processingBlocks,
			}
			got, err := r.GetFailedBlock(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetFailedBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFailedBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemBlock_Upsert(t *testing.T) {
	type fields struct {
		failedBlocks          map[int]entity.Block
		processingBlocks      map[int]entity.Block
		parsedBlocks          map[int]entity.Block
		parsedBlockNumber     int
		processingBlockNumber int
	}

	ctx := context.Background()

	type args struct {
		ctx   context.Context
		block entity.Block
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr error
	}{
		{
			name: "unknown block status",
			fields: fields{
				failedBlocks:          map[int]entity.Block{},
				processingBlocks:      map[int]entity.Block{},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 1,
					Status: "unknown",
				},
			},
			want: fields{
				failedBlocks:          map[int]entity.Block{},
				processingBlocks:      map[int]entity.Block{},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			wantErr: errorpkg.UnknownBlockStatus,
		},
		{
			name: "block became processing, NOT update processing block number",
			fields: fields{
				failedBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusFailed,
					},
				},
				processingBlocks:      map[int]entity.Block{},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 1,
					Status: constant.BlockStatusProcessing,
				},
			},
			want: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusProcessing,
					},
				},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			wantErr: nil,
		},
		{
			name: "block became processing, update processing block number",
			fields: fields{
				failedBlocks:          map[int]entity.Block{},
				processingBlocks:      map[int]entity.Block{},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 1,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 2,
					Status: constant.BlockStatusProcessing,
				},
			},
			want: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusProcessing,
					},
				},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			wantErr: nil,
		},
		{
			name: "block became parsed, NOT update parsed block number",
			fields: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusFailed,
					},
				},
				parsedBlocks:      map[int]entity.Block{},
				parsedBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 1,
					Status: constant.BlockStatusParsed,
				},
			},
			want: fields{
				failedBlocks:     map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{},
				parsedBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusParsed,
					},
				},
				parsedBlockNumber: 2,
			},
			wantErr: nil,
		},
		{
			name: "block became parsed, update parsed block number",
			fields: fields{
				failedBlocks:      map[int]entity.Block{},
				processingBlocks:  map[int]entity.Block{},
				parsedBlocks:      map[int]entity.Block{},
				parsedBlockNumber: 1,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 2,
					Status: constant.BlockStatusParsed,
				},
			},
			want: fields{
				failedBlocks:     map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{},
				parsedBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusParsed,
					},
				},
				parsedBlockNumber: 2,
			},
			wantErr: nil,
		},
		{
			name: "block became failed, NOT update processing block number",
			fields: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusProcessing,
					},
				},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 1,
					Status: constant.BlockStatusFailed,
				},
			},
			want: fields{
				failedBlocks: map[int]entity.Block{
					1: {
						Number: 1,
						Status: constant.BlockStatusFailed,
					},
				},
				processingBlocks:      map[int]entity.Block{},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			wantErr: nil,
		},
		{
			name: "block became failed, update processing block number",
			fields: fields{
				failedBlocks: map[int]entity.Block{},
				processingBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusProcessing,
					},
				},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 2,
			},
			args: args{
				ctx: ctx,
				block: entity.Block{
					Number: 2,
					Status: constant.BlockStatusFailed,
				},
			},
			want: fields{
				failedBlocks: map[int]entity.Block{
					2: {
						Number: 2,
						Status: constant.BlockStatusFailed,
					},
				},
				processingBlocks:      map[int]entity.Block{},
				parsedBlocks:          map[int]entity.Block{},
				processingBlockNumber: 1,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemBlock{
				failedBlocks:          tt.fields.failedBlocks,
				processingBlocks:      tt.fields.processingBlocks,
				parsedBlocks:          tt.fields.parsedBlocks,
				parsedBlockNumber:     tt.fields.parsedBlockNumber,
				processingBlockNumber: tt.fields.processingBlockNumber,
			}
			if err := r.Upsert(tt.args.ctx, tt.args.block); !errors.Is(err, tt.wantErr) {
				t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(r.failedBlocks, tt.want.failedBlocks) {
				t.Errorf("failedBlocks got = %v, want %v", r.failedBlocks, tt.want.failedBlocks)
			}
			if !reflect.DeepEqual(r.processingBlocks, tt.want.processingBlocks) {
				t.Errorf("processingBlocks got = %v, want %v", r.processingBlocks, tt.want.processingBlocks)
			}
			if !reflect.DeepEqual(r.parsedBlocks, tt.want.parsedBlocks) {
				t.Errorf("parsedBlocks got = %v, want %v", r.parsedBlocks, tt.want.parsedBlocks)
			}
			if !reflect.DeepEqual(r.parsedBlockNumber, tt.want.parsedBlockNumber) {
				t.Errorf("parsedBlockNumber got = %v, want %v", r.parsedBlockNumber, tt.want.parsedBlockNumber)
			}
			if !reflect.DeepEqual(r.processingBlockNumber, tt.want.processingBlockNumber) {
				t.Errorf("processingBlockNumber got = %v, want %v", r.processingBlockNumber, tt.want.processingBlockNumber)
			}
		})
	}
}
