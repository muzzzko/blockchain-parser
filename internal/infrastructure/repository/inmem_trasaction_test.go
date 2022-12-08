package repository

import (
	"blockchain-parser/internal/entity"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestInMemTransaction_Save(t *testing.T) {
	type fields struct {
		data map[string]map[string]*entity.Transaction
	}

	ctx := context.Background()
	txn := entity.Transaction{
		From:             "0x42352",
		To:               "0x245212",
		Value:            "0x3453",
		BlockNumber:      1,
		TransactionIndex: 5,
	}

	type args struct {
		ctx         context.Context
		transaction entity.Transaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]map[string]*entity.Transaction
		wantErr error
	}{
		{
			name: "save transaction",
			fields: fields{
				data: map[string]map[string]*entity.Transaction{},
			},
			args: args{
				ctx:         ctx,
				transaction: txn,
			},
			want: map[string]map[string]*entity.Transaction{
				"0x42352": {
					"1_5": &txn,
				},
				"0x245212": {
					"1_5": &txn,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemTransaction{
				data: tt.fields.data,
			}
			if err := r.Save(tt.args.ctx, tt.args.transaction); !errors.Is(err, tt.wantErr) {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(r.data, tt.want) {
				t.Errorf("data got = %v, want %v", r.data, tt.want)
			}
		})
	}
}

func TestInMemTransaction_GetTxnsByAddress(t *testing.T) {
	type fields struct {
		data map[string]map[string]*entity.Transaction
	}

	ctx := context.Background()
	txn := entity.Transaction{
		From:             "0x42352",
		To:               "0x245212",
		Value:            "0x3453",
		BlockNumber:      1,
		TransactionIndex: 5,
	}

	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Transaction
		wantErr bool
	}{
		{
			name: "no transactions",
			fields: fields{
				data: map[string]map[string]*entity.Transaction{},
			},
			args: args{
				ctx:     ctx,
				address: "0x42352",
			},
			want:    []entity.Transaction{},
			wantErr: false,
		},
		{
			name: "return transactions",
			fields: fields{
				data: map[string]map[string]*entity.Transaction{
					"0x42352": {
						"1_5": &txn,
					},
					"0x245212": {
						"1_5": &txn,
					},
				},
			},
			args: args{
				ctx:     ctx,
				address: "0x245212",
			},
			want: []entity.Transaction{
				txn,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemTransaction{
				data: tt.fields.data,
			}
			got, err := r.GetTxnsByAddress(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxnsByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTxnsByAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
