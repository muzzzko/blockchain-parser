package repository

import (
	errorpkg "blockchain-parser/internal/error"
	"context"
	"errors"
	"reflect"
	"testing"

	"blockchain-parser/internal/entity"
)

func TestInMemSubscriber_Save(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		subscriber entity.Subscriber
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]entity.Subscriber
		wantErr error
	}{
		{
			name: "save subscriber",
			args: args{
				ctx: ctx,
				subscriber: entity.Subscriber{
					Address: "0x41da31",
				},
			},
			want: map[string]entity.Subscriber{
				"0x41da31": {
					Address: "0x41da31",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewInMemSubscriber()
			if err := r.Save(tt.args.ctx, tt.args.subscriber); !errors.Is(err, tt.wantErr) {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(r.data, tt.want) {
				t.Errorf("data got = %v, want %v", r.data, tt.want)
			}
		})
	}
}

func TestInMemSubscriber_Get(t *testing.T) {
	type fields struct {
		data map[string]entity.Subscriber
	}

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Subscriber
		wantErr error
	}{
		{
			name: "subscriber not found",
			fields: fields{
				data: map[string]entity.Subscriber{},
			},
			args: args{
				ctx:     ctx,
				address: "0x41da31",
			},
			want:    entity.Subscriber{},
			wantErr: errorpkg.SubscriberNotFound,
		},
		{
			name: "return subscriber",
			fields: fields{
				data: map[string]entity.Subscriber{
					"0x41da31": {
						Address: "0x41da31",
					},
				},
			},
			args: args{
				ctx:     ctx,
				address: "0x41da31",
			},
			want:    entity.Subscriber{Address: "0x41da31"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemSubscriber{
				data: tt.fields.data,
			}
			got, err := r.Get(tt.args.ctx, tt.args.address)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
