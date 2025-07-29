package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"

	"github.com/google/uuid"
)

type (
	IScheduleService interface {
		CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (dto.ScheduleResponse, error)
		GetAllSchedule(ctx context.Context) ([]dto.ScheduleResponse, error)
		GetScheduleByID(ctx context.Context, scheduleID string) (dto.ScheduleFullResponse, error)
		UpdateSchedule(ctx context.Context, req dto.UpdateScheduleRequest) (dto.ScheduleResponse, error)
		DeleteScheduleByID(ctx context.Context, req dto.DeleteScheduleRequest) (dto.ScheduleResponse, error)

		GetSchedulesByFieldID(ctx context.Context, fieldID string) ([]dto.ScheduleResponse, error)
		GetScheduleByFieldIDAndDay(ctx context.Context, fieldID string, day int) (dto.ScheduleResponse, error)
	}

	ScheduleService struct {
		scheduleRepo repository.IScheduleRepository
		fieldRepo    repository.IFieldRepository
	}
)

func NewScheduleService(scheduleRepo repository.IScheduleRepository, fieldRepo repository.IFieldRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		fieldRepo:    fieldRepo,
	}
}

func (ss *ScheduleService) CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (dto.ScheduleResponse, error) {
	loc := helpers.GetAppLocation()

	fieldUUID, err := uuid.Parse(req.FieldID)
	if err != nil {
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	_, _, err = ss.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		return dto.ScheduleResponse{}, constants.ErrFieldNotFound
	}

	schedule := model.Schedule{
		ScheduleID: uuid.New(),
		FieldID:    fieldUUID,
		DayOfWeek:  req.DayOfWeek,
		OpenTime:   req.OpenTime.In(loc),
		CloseTime:  req.CloseTime.In(loc),
	}

	if err := ss.scheduleRepo.CreateSchedule(ctx, nil, schedule); err != nil {
		return dto.ScheduleResponse{}, constants.ErrCreateSchedule
	}

	res := dto.ScheduleResponse{
		ScheduleID: schedule.ScheduleID,
		FieldID:    schedule.FieldID,
		DayOfWeek:  schedule.DayOfWeek,
		DayName:    helpers.DayIntToName(schedule.DayOfWeek),
		OpenTime:   schedule.OpenTime.In(loc).Format("15:04"),
		CloseTime:  schedule.CloseTime.In(loc).Format("15:04"),
	}

	return res, nil
}

func (ss *ScheduleService) GetAllSchedule(ctx context.Context) ([]dto.ScheduleResponse, error) {
	loc := helpers.GetAppLocation()

	schedules, err := ss.scheduleRepo.GetAllSchedule(ctx, nil)
	if err != nil {
		return nil, constants.ErrGetAllSchedule
	}

	var res []dto.ScheduleResponse
	for _, s := range schedules {
		res = append(res, dto.ScheduleResponse{
			ScheduleID: s.ScheduleID,
			FieldID:    s.FieldID,
			DayOfWeek:  s.DayOfWeek,
			DayName:    helpers.DayIntToName(s.DayOfWeek),
			OpenTime:   s.OpenTime.In(loc).Format("15:04"),
			CloseTime:  s.CloseTime.In(loc).Format("15:04"),
		})
	}
	return res, nil
}

func (ss *ScheduleService) GetScheduleByID(ctx context.Context, scheduleID string) (dto.ScheduleFullResponse, error) {
	loc := helpers.GetAppLocation()

	if _, err := uuid.Parse(scheduleID); err != nil {
		return dto.ScheduleFullResponse{}, constants.ErrInvalidUUID
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, scheduleID)
	if err != nil {
		return dto.ScheduleFullResponse{}, constants.ErrScheduleNotFound
	}

	field := schedule.Field
	category := field.Category

	fieldDTO := dto.FieldFullResponse{
		FieldID:      field.FieldID,
		FieldName:    field.FieldName,
		FieldAddress: field.FieldAddress,
		FieldPrice:   field.FieldPrice,
		FieldImage:   field.FieldImage,
		Category: dto.CategoryCompactResponse{
			CategoryID:  category.CategoryID,
			Name:        category.Name,
			Description: category.Description,
		},
	}

	res := dto.ScheduleFullResponse{
		ScheduleID: schedule.ScheduleID,
		DayOfWeek:  schedule.DayOfWeek,
		DayName:    helpers.DayIntToName(schedule.DayOfWeek),
		OpenTime:   schedule.OpenTime.In(loc).Format("15:04"),
		CloseTime:  schedule.CloseTime.In(loc).Format("15:04"),
		Field:      dto.FieldCompactResponse(fieldDTO),
	}

	return res, nil
}

func (ss *ScheduleService) UpdateSchedule(ctx context.Context, req dto.UpdateScheduleRequest) (dto.ScheduleResponse, error) {
	loc := helpers.GetAppLocation()
	
	if _, err := uuid.Parse(req.ScheduleID); err != nil {
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, req.ScheduleID)
	if err != nil {
		return dto.ScheduleResponse{}, constants.ErrGetScheduleByID
	}

	if req.FieldID != nil {
		parsedUUID, err := uuid.Parse(*req.FieldID)
		if err != nil {
			return dto.ScheduleResponse{}, constants.ErrInvalidUUID
		}
		schedule.FieldID = parsedUUID
	}

	if req.DayOfWeek != nil {
		if *req.DayOfWeek < 0 || *req.DayOfWeek > 6 {
			return dto.ScheduleResponse{}, constants.ErrInvalidDayOfWeek
		}
		schedule.DayOfWeek = *req.DayOfWeek
	}

	if req.OpenTime != nil {
		schedule.OpenTime = req.OpenTime.In(loc)
	}
	if req.CloseTime != nil {
		schedule.CloseTime = req.CloseTime.In(loc)
	}

	if !schedule.CloseTime.After(schedule.OpenTime) {
		return dto.ScheduleResponse{}, constants.ErrCloseTimeMustAfterOpen
	}

	if err := ss.scheduleRepo.UpdateSchedule(ctx, nil, schedule); err != nil {
		return dto.ScheduleResponse{}, constants.ErrUpdateSchedule
	}

	res := dto.ScheduleResponse{
		ScheduleID: schedule.ScheduleID,
		FieldID:    schedule.FieldID,
		DayOfWeek:  schedule.DayOfWeek,
		DayName:    helpers.DayIntToName(schedule.DayOfWeek),
		OpenTime:   schedule.OpenTime.Format("15:04"),
		CloseTime:  schedule.CloseTime.Format("15:04"),
	}

	return res, nil
}

func (ss *ScheduleService) DeleteScheduleByID(ctx context.Context, req dto.DeleteScheduleRequest) (dto.ScheduleResponse, error) {
	loc := helpers.GetAppLocation()
	
	if _, err := uuid.Parse(req.ScheduleID); err != nil {
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, req.ScheduleID)
	if err != nil {
		return dto.ScheduleResponse{}, constants.ErrScheduleNotFound
	}

	if err := ss.scheduleRepo.DeleteScheduleByID(ctx, nil, req.ScheduleID); err != nil {
		return dto.ScheduleResponse{}, constants.ErrDeleteSchedule
	}

	res := dto.ScheduleResponse{
		ScheduleID: schedule.ScheduleID,
		FieldID:    schedule.FieldID,
		DayOfWeek:  schedule.DayOfWeek,
		DayName:    helpers.DayIntToName(schedule.DayOfWeek),
		OpenTime:   schedule.OpenTime.In(loc).Format("15:04"),
		CloseTime:  schedule.CloseTime.In(loc).Format("15:04"),
	}

	return res, nil
}

func (ss *ScheduleService) GetSchedulesByFieldID(ctx context.Context, fieldID string) ([]dto.ScheduleResponse, error) {
	loc := helpers.GetAppLocation()
	
	if _, err := uuid.Parse(fieldID); err != nil {
		return nil, constants.ErrInvalidUUID
	}

	schedules, err := ss.scheduleRepo.GetSchedulesByFieldID(ctx, nil, fieldID)
	if err != nil {
		return nil, constants.ErrGetAllSchedule
	}

	var res []dto.ScheduleResponse
	for _, s := range schedules {
		res = append(res, dto.ScheduleResponse{
			ScheduleID: s.ScheduleID,
			FieldID:    s.FieldID,
			DayOfWeek:  s.DayOfWeek,
			DayName:    helpers.DayIntToName(s.DayOfWeek),
			OpenTime:   s.OpenTime.In(loc).Format("15:04"),
			CloseTime:  s.CloseTime.In(loc).Format("15:04"),
		})
	}
	return res, nil
}

func (ss *ScheduleService) GetScheduleByFieldIDAndDay(ctx context.Context, fieldID string, day int) (dto.ScheduleResponse, error) {
	loc := helpers.GetAppLocation()
	
	if _, err := uuid.Parse(fieldID); err != nil {
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	s, err := ss.scheduleRepo.GetScheduleByFieldIDAndDay(ctx, nil, fieldID, day)
	if err != nil {
		return dto.ScheduleResponse{}, constants.ErrScheduleNotFound
	}

	res := dto.ScheduleResponse{
		ScheduleID: s.ScheduleID,
		FieldID:    s.FieldID,
		DayOfWeek:  s.DayOfWeek,
		DayName:    helpers.DayIntToName(s.DayOfWeek),
		OpenTime:   s.OpenTime.In(loc).Format("15:04"),
		CloseTime:  s.CloseTime.In(loc).Format("15:04"),
	}

	return res, nil
}
