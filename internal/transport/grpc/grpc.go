package grpc

import (
	"account/internal/application/dtos"
	pb "account/proto"
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedAccountServer
	server *grpc.Server

	useCase UseCase
}

func New(useCase UseCase) *server {
	return &server{
		server:  grpc.NewServer(),
		useCase: useCase}
}

func (s *server) Run(socket string) error {
	listener, err := net.Listen("tcp", socket)
	if err != nil {
		return err
	}

	pb.RegisterAccountServer(s.server, s)

	return s.server.Serve(listener)
}

func (s *server) Stop() {
	s.server.GracefulStop()
}

func (s *server) Register(c context.Context, req *pb.RegisterReq) (*emptypb.Empty, error) {
	err := s.useCase.Register(c, &dtos.RegisterInput{
		Phone: req.Phone,
	})

	return &emptypb.Empty{}, err
}

func (s *server) ConfirmCode(c context.Context, req *pb.ConfirmCodeReq) (*emptypb.Empty, error) {
	err := s.useCase.ConfirmCode(c, &dtos.ConfirmCodeInput{
		Phone: req.Phone,
		Code:  req.Code,
	})
	return &emptypb.Empty{}, err
}

func (s *server) CompleteRegister(c context.Context, req *pb.CompleteRegisterReq) (*pb.AuthRes, error) {
	res, err := s.useCase.CompleteRegister(c, &dtos.CompleteRegisterInput{
		Phone:    req.Phone,
		Code:     req.Code,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return &pb.AuthRes{
			AccessToken: "",
		}, err
	}

	return &pb.AuthRes{
		AccessToken: res.AccessToken,
	}, nil
}

func (s *server) Login(c context.Context, req *pb.LoginReq) (*pb.AuthRes, error) {
	res, err := s.useCase.Login(c, &dtos.LoginInput{Phone: req.Phone, Password: req.Password})
	if err != nil {
		return &pb.AuthRes{AccessToken: ""}, err
	}

	return &pb.AuthRes{
		AccessToken: res.AccessToken,
	}, nil
}

func (s *server) Recovery(c context.Context, req *pb.RecoveryReq) (*emptypb.Empty, error) {
	err := s.useCase.Recovery(c, &dtos.RecoveryInput{Phone: req.Phone})
	return &emptypb.Empty{}, err
}

func (s *server) CompleteRecovery(c context.Context, req *pb.CompleteRecoveryReq) (*pb.AuthRes, error) {
	dto := &dtos.CompleteRecovery{
		Phone:    req.Phone,
		Code:     req.Code,
		Password: req.Password,
	}
	session, err := s.useCase.CompleteRecovery(c, dto)
	if err != nil {
		return &pb.AuthRes{}, err
	}

	return &pb.AuthRes{AccessToken: session.AccessToken}, nil
}

func (s *server) Ban(c context.Context, req *pb.BanReq) (*emptypb.Empty, error) {
	expiry, err := time.Parse(time.RFC3339, req.Expiry)

	if err != nil {
		return &emptypb.Empty{}, err
	}

	dto := &dtos.BanInput{
		CallerUserID: int(req.CallerUserId),
		UserID:       int(req.UserId),
		Reason:       req.Reason,
		Expiry:       expiry,
	}
	err = s.useCase.Ban(c, dto)

	return &emptypb.Empty{}, err
}

func (s *server) UnBan(c context.Context, req *pb.UnBanReq) (*emptypb.Empty, error) {
	dto := &dtos.UnBanInput{
		UserID: int(req.UserId),
		BanID:  int(req.BanId),
		Reason: req.Reason,
	}
	err := s.useCase.UnBan(c, dto)
	return &emptypb.Empty{}, err
}

func (s *server) Authenticate(c context.Context, req *pb.AuthenticateReq) (*pb.AuthenticateRes, error) {
	output, err := s.useCase.Authenticate(c, &dtos.AuthenticateInput{
		AccessToken: req.AccessToken,
	})

	if err != nil {
		return &pb.AuthenticateRes{}, err
	}

	return &pb.AuthenticateRes{
		UserId: int32(output.UserID),
		Role: &pb.Role{
			Id:          int32(output.Role.ID),
			Name:        output.Role.Name,
			Permissions: output.Role.Permissions,
		},
		IsBanned: output.IsBanned,
	}, nil
}

func (s *server) FindByID(c context.Context, req *pb.FindByIDReq) (*pb.UserRes, error) {
	user, err := s.useCase.FindByID(c, &dtos.FindByIDInput{ID: int(req.Id)})
	if err != nil {
		return &pb.UserRes{}, err
	}
	return &pb.UserRes{
		Id:          int32(user.ID),
		Username:    user.Username,
		Photo:       user.Photo,
		Description: user.Description,
		Role: &pb.Role{
			Id:          int32(user.Role.ID),
			Name:        user.Role.Name,
			Permissions: user.Role.Permissions,
		},
	}, nil
}

func (s *server) FindByUsername(c context.Context, req *pb.FindByUsernameReq) (*pb.UserRes, error) {
	user, err := s.useCase.FindByUsername(c, &dtos.FindByUsernameInput{Username: req.Username})
	if err != nil {
		return &pb.UserRes{}, err
	}
	return &pb.UserRes{
		Id:          int32(user.ID),
		Username:    user.Username,
		Photo:       user.Photo,
		Description: user.Description,
		Role: &pb.Role{
			Id:          int32(user.Role.ID),
			Name:        user.Role.Name,
			Permissions: user.Role.Permissions,
		},
	}, nil
}

func (s *server) IsPhoneExists(c context.Context, req *pb.IsPhoneExistsReq) (*emptypb.Empty, error) {
	err := s.useCase.IsPhoneExists(c, &dtos.IsPhoneExistsInput{Phone: req.Phone})
	return &emptypb.Empty{}, err
}

func (s *server) ChangePassword(c context.Context, req *pb.ChangePasswordReq) (*pb.AuthRes, error) {
	res, err := s.useCase.ChangePassword(c, &dtos.ChangePasswordInput{
		UserID:      int(req.UserId),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
		Logout:      req.Logout,
	})
	if err != nil {
		return &pb.AuthRes{}, err
	}

	return &pb.AuthRes{
		AccessToken: res.AccessToken,
	}, nil
}

func (s *server) UpdatePhoto(c context.Context, req *pb.UpdatePhotoReq) (*pb.UpdatePhotoRes, error) {
	photo, err := s.useCase.UploadImage(c, &dtos.UploadImageInput{
		UserID: int(req.UserId),
		File:   req.File,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePhotoRes{
		Photo: photo.Filename,
	}, nil
}

func (s *server) RemovePhoto(c context.Context, req *pb.RemovePhotoReq) (*emptypb.Empty, error) {
	err := s.useCase.RemoveImage(c, &dtos.RemovePhotoInput{UserID: int(req.UserId)})
	return &emptypb.Empty{}, err
}
