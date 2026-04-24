package timer

import (
	"context"
	"errors"
	"time"

	gen "github.com/zimlewis/tomato/gen/proto"
	errs "github.com/zimlewis/tomato/internal/errors"
	"github.com/zimlewis/tomato/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)


var waitTime = []int64{25, 15, 30}

type Service struct {
	gen.UnimplementedTimerServer
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetClock(ctx context.Context, _ *emptypb.Empty) (*gen.GetClockResponse, error) {
	clock, err := s.repo.GetClock(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot get code: %s", err)
	}

	return &gen.GetClockResponse{
		Clock: int32(clock),
	}, nil
}

func (s *Service) Current(ctx context.Context, _ *emptypb.Empty) (*gen.CurrentTimer, error) {
	var result gen.CurrentTimer
	
	clock, err := s.repo.GetClock(ctx)
	if err != nil {
		return nil, err
	}

	startTime, err := s.repo.GetStartTime(ctx)
	if errors.Is(err, errs.ErrDidNotStart) {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
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
	return nil, err
}

func (s *Service) Stop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	// Delete the start time
	err := s.repo.DeleteStartTime(ctx)
	return nil, err
}

func (s *Service) Switch(ctx context.Context, dir *gen.SwitchRequest) (*emptypb.Empty, error) {
	// Get the clock type and switch it arcodingly
	clock, err := s.repo.GetClock(ctx)
	if err != nil {
		return nil, err
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
		return nil, status.Errorf(codes.InvalidArgument, "Wrong direction format: %s\n", dir.String())
	}

	// Set the clock type 
	err = s.repo.SetClock(ctx, int(valueToSwitch))
	if err != nil {
		return nil, err
	}

	// Delete the start time
	err = s.repo.DeleteStartTime(ctx)
	return nil, err
}
