package delivery

import (
	pb "api/app/protobuf"
	"context"
)

func (s *server) AddMarketDeck(ctx context.Context, request *pb.AddMarketDeckRequest) (*pb.AddMarketDeckReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) UpdateMarketDeck(ctx context.Context, request *pb.UpdateMarketDeckRequest) (
	*pb.UpdateMarketDeckReply, error,
) {
	//TODO implement me
	panic("implement me")
}

func (s *server) DeleteMarketDeck(ctx context.Context, request *pb.DeleteMarketDeckRequest) (
	*pb.DeleteMarketDeckReply, error,
) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetMarketDecks(ctx context.Context, request *pb.GetMarketDeckRequest) (*pb.GetMarketDeckReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetPopularMarketDecks(
	ctx context.Context, request *pb.GetPopularMarketDecksRequest,
) (*pb.GetPopularMarketDecksReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetMarketDeck(ctx context.Context, request *pb.GetMarketDeckRequest) (*pb.GetMarketDeckReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetMyMarketDecks(ctx context.Context, request *pb.GetMyMarketDecksRequest) (
	*pb.GetMyMarketDecksReply, error,
) {
	//TODO implement me
	panic("implement me")
}
