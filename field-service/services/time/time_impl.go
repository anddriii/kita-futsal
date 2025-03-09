package time

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"github.com/anddriii/kita-futsal/field-service/repositories"
)

// TimeService is a service struct that provides time-related functionalities.
type TimeService struct {
	repository repositories.IRepoRegistry
}

// Create creates a new time entry in the database.
// It takes a context and a TimeRequest object as input.
// Returns a TimeResponse containing the created time record or an error if the operation fails.
func (t *TimeService) Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error) {
	// Create a TimeRequest object from the request
	time := dto.TimeRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Save the time entry into the database
	timeResult, err := t.repository.GetTime().Create(ctx, &models.Time{
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
	})

	if err != nil {
		return nil, err
	}

	// Construct the response object
	response := dto.TimeResponse{
		UUID:      timeResult.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: timeResult.CreatedAt,
		UpdateAt:  timeResult.UpdatedAt,
	}

	return &response, nil
}

// GetAll retrieves all time records from the database.
// It takes a context as input and returns a slice of TimeResponse or an error.
func (t *TimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {
	// Retrieve all time records from the repository
	times, err := t.repository.GetTime().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert database results to response objects
	timeResults := make([]dto.TimeResponse, 0, len(times))
	for _, time := range times {
		timeResults = append(timeResults, dto.TimeResponse{
			UUID:      time.UUID,
			StartTime: time.StartTime,
			EndTime:   time.EndTime,
			CreatedAt: time.CreatedAt,
			UpdateAt:  time.UpdatedAt,
		})
	}

	return timeResults, nil
}

// GetByUUID retrieves a specific time record by its UUID.
// It takes a context and UUID string as input.
// Returns a TimeResponse object if the record exists, or an error if not found.
func (t *TimeService) GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error) {
	// Fetch the time record from the database by UUID
	time, err := t.repository.GetTime().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Construct the response object
	timeResult := dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdateAt:  time.UpdatedAt,
	}

	return &timeResult, nil
}

// NewTimeService creates a new instance of TimeService.
// It takes a repository registry as a dependency and returns an instance of ITimeService.
func NewTimeService(repository repositories.IRepoRegistry) ITimeService {
	return &TimeService{repository: repository}
}
