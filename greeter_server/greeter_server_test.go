package main

import (
	"context"
	"fmt"
	"github.com/zencore/helloworld/expect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/zencore/helloworld/proto/helloworld"
	"google.golang.org/protobuf/proto"
)

// Expectation is a embedding type for assert.Assertions
type Expectation struct {
	*assert.Assertions
}

// New makes a new Expectation object for the specified TestingT.
func New(t assert.TestingT) *Expectation {
	assert := assert.New(t)
	return &Expectation{assert}
}

// TODO(james): Equality depends on order of repeated fields. This may break some tests unnecessarily.
// ProtoEqual asserts that the specified protobuf messages are equal.
func (a *Expectation) ProtoEqual(expected, actual proto.Message) bool {
	return a.True(
		proto.Equal(expected, actual),
		fmt.Sprintf("These two protobuf messages are not equal:\nexpected: %v\nactual:  %v", expected, actual),
	)
}

// chat-service_test.go
func Test_greeterServiceServer_ListMessages(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		req *pb.ListMessagesRequest
	}

	server := &ServerImpl{}

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

			expect := expect.New(t)
			expect.ProtoEqual(tt.want, got)
		})
	}
}
