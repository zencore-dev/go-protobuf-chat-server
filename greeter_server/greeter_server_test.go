package main

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/zencore/helloworld/proto/helloworld"
)

// chat-service_test.go
func Test_greeterServiceServer_ListMessages(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		req *pb.ListMessagesRequest
	}

	server := NewServer()

	// Preamble setup.
	msgs := []*pb.Message{
		{
			MessageUuid:       "123",
			ExternalHistoryId: "opaque-identifier",
			State:             pb.Message_STATE_OK,
		},
		{
			MessageUuid:       "456",
			ExternalHistoryId: "opaque-identifier",
			State:             pb.Message_STATE_OK,
		},
		{
			MessageUuid:       "789",
			ExternalHistoryId: "opaque-identifier",
			State:             pb.Message_STATE_OK,
		},
	}

	for _, msg := range msgs {
		_, err := server.CommitMessage(ctx, &pb.CommitMessageRequest{Message: msg})
		assert.Nil(t, err)
	}

	tests := []struct {
		name        string
		server      pb.ChatServiceServer
		args        args
		want        *pb.ListMessagesResponse
		wantErr     bool
		wantErrCode codes.Code
	}{
		/*
			{
				name:   "OK_EMPTY",
				server: server,
				args: args{
					ctx: ctx,
					req: &pb.ListMessagesRequest{
						ExternalHistoryId: "non-existent-identifier",
					},
				},
				want: &pb.ListMessagesResponse{
					Messages: []*pb.Message{},
				},
			},
		*/
		{
			name:   "OK",
			server: server,
			args: args{
				ctx: ctx,
				req: &pb.ListMessagesRequest{
					ExternalHistoryId: "opaque-identifier",
				},
			},
			want: &pb.ListMessagesResponse{
				Messages: msgs,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.ListMessages(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("returned non-status error")
				}

				if st.Code() != tt.wantErrCode {
					t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}

			if len(got.GetMessages()) != len(tt.want.GetMessages()) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}

			// Filter out timestamps.
			for i := range got.GetMessages() {
				got.GetMessages()[i].CreateTime = nil
				got.GetMessages()[i].LastUpdateTime = nil
			}

			assert.Empty(t, cmp.Diff(tt.want, got, protocmp.Transform(), cmpopts.IgnoreUnexported()))
		})
	}
}
