package service

import (
	"context"

	pb "membership/api/membership"
)

type MembershipService struct {
	pb.UnimplementedMembershipServer
}

func NewMembershipService() *MembershipService {
	return &MembershipService{}
}

func (s *MembershipService) CreateMembership(ctx context.Context, req *pb.CreateMembershipRequest) (*pb.CreateMembershipReply, error) {
	return &pb.CreateMembershipReply{}, nil
}
func (s *MembershipService) UpdateMembership(ctx context.Context, req *pb.UpdateMembershipRequest) (*pb.UpdateMembershipReply, error) {
	return &pb.UpdateMembershipReply{}, nil
}
func (s *MembershipService) DeleteMembership(ctx context.Context, req *pb.DeleteMembershipRequest) (*pb.DeleteMembershipReply, error) {
	return &pb.DeleteMembershipReply{}, nil
}
func (s *MembershipService) GetMembership(ctx context.Context, req *pb.GetMembershipRequest) (*pb.GetMembershipReply, error) {
	return &pb.GetMembershipReply{}, nil
}
func (s *MembershipService) ListMembership(ctx context.Context, req *pb.ListMembershipRequest) (*pb.ListMembershipReply, error) {
	return &pb.ListMembershipReply{}, nil
}
