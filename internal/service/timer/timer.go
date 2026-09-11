package timer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	gen "github.com/zimlewis/tomato/gen/proto"
	"github.com/zimlewis/tomato/internal/tomatoerrs"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)


var waitTime = []int64{25, 5, 30}

type repository interface {
	DeleteStartTime(ctx context.Context) error
	GetClock(ctx context.Context) (uint16, error)
	GetStartTime(ctx context.Context) (int64, error)
	SetClock(ctx context.Context, clockIndex int) error
	SetStartTime(ctx context.Context, time int64) error
}

type Service struct {
	gen.UnimplementedTimerServer
	repo   repository
	logger *slog.Logger
}


func New(repo repository, logger *slog.Logger) *Service {
	return &Service{ repo: repo, logger: logger }
}



func (s *Service) SetClock(ctx context.Context, req *gen.SetClockRequest) (*emptypb.Empty, error) {
	valueToSwitch := req.Clock
	if valueToSwitch < 0 || valueToSwitch > 2 {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			fmt.Errorf("Invalid clock value(value must be in range from 0 to 2): %d", valueToSwitch),
			codes.InvalidArgument,
			"invalid clock value",
		)
	}

	// Set the clock type 
	err := s.repo.SetClock(ctx, int(valueToSwitch))
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger, 
			err,
			codes.Internal,
			"unable to set clock type",
		)
	}

	// Delete the start time
	err = s.repo.DeleteStartTime(ctx)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger, 
			err, 
			codes.Internal, 
			"unable to reset start time",
		)
	}
	return nil, nil
}

func (s *Service) GetClock(ctx context.Context, _ *emptypb.Empty) (*gen.GetClockResponse, error) {
	clock, err := s.repo.GetClock(ctx)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			err,
			codes.Internal,
			"cannot get the code",
		)

	}

	return &gen.GetClockResponse{
		Clock: int32(clock),
	}, nil
}

func (s *Service) Current(ctx context.Context, _ *emptypb.Empty) (*gen.CurrentTimer, error) {
	var result gen.CurrentTimer
	
	clock, err := s.repo.GetClock(ctx)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			err,
			codes.Internal,
			"cannot get current clock",
		)
	}

	startTime, err := s.repo.GetStartTime(ctx)
	if errors.Is(err, tomatoerrs.ErrDidNotStart) {
		return nil, tomatoerrs.GRPCError(
			s.logger, 
			err, 
			codes.NotFound,
			"the session did not start",
		)
	}
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			err,
			codes.Internal,
			"cannot get start the session",
		)
	}

	currentTime := time.Now().Unix()
	elapsed := currentTime - startTime

	timeLeft := waitTime[clock] * 60 - elapsed 

	result.Clock = int32(clock)
	result.TimeLeft = timeLeft

	return &result, nil
}

func (s *Service) Start(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	// Set the start time to current time
	currentTime := time.Now().Unix()
	err := s.repo.SetStartTime(ctx, currentTime)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			err,
			codes.Internal,
			"cannot start session",
		)
	}
	return nil, nil
}

func (s *Service) Stop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	// Delete the start time
	err := s.repo.DeleteStartTime(ctx)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			err,
			codes.Internal,
			"cannot stop the clock",
		)
	}
	return nil, nil
}

func (s *Service) Switch(ctx context.Context, dir *gen.SwitchRequest) (*emptypb.Empty, error) {
	// Get the clock type and switch it arcodingly
	clock, err := s.repo.GetClock(ctx)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger, 
			err, 
			codes.Internal, 
			"cannot get the clock",
		)
	}

	var valueToSwitch int
	switch dir.Dir.String() {
	case "UP":
		valueToSwitch = int(clock) + 1
		if valueToSwitch > 2 {
			valueToSwitch = 0
		}
	case "DOWN":
		valueToSwitch = int(clock) - 1
		if valueToSwitch < 0 {
			valueToSwitch = 2
		}
	default:
		return nil, tomatoerrs.GRPCError(
			s.logger, 
			fmt.Errorf("wrong direction format: %s", dir.String()), 
			codes.InvalidArgument, 
			"wrong direction format",
		)
	}

	// Set the clock type 
	err = s.repo.SetClock(ctx, int(valueToSwitch))
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger, 
			err, 
			codes.Internal, 
			"cannot set clock",
		)
	}

	// Delete the start time
	err = s.repo.DeleteStartTime(ctx)
	if err != nil {
		return nil, tomatoerrs.GRPCError(
			s.logger,
			err,
			codes.Internal,
			"cannot delete start time",
		)
	}
	return nil, nil
}
