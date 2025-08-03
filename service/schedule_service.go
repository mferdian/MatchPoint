package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"
	"fieldreserve/utils"

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
	utils.Log.Infof("Creating schedule for field ID: %s", req.FieldID)

	loc := helpers.GetAppLocation()

	fieldUUID, err := uuid.Parse(req.FieldID)
	if err != nil {
		utils.Log.Errorf("Invalid field UUID: %v", err)
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	_, _, err = ss.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		utils.Log.Errorf("Field not found: %v", err)
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
		utils.Log.Errorf("Failed to create schedule: %v", err)
		return dto.ScheduleResponse{}, constants.ErrCreateSchedule
	}

	utils.Log.Infof("Schedule created successfully: %s", schedule.ScheduleID)

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

func (ss *ScheduleService) GetAllSchedule(ctx context.Context) ([]dto.ScheduleResponse, error) {
	utils.Log.Info("Fetching all schedules")

	loc := helpers.GetAppLocation()

	schedules, err := ss.scheduleRepo.GetAllSchedule(ctx, nil)
	if err != nil {
		utils.Log.Errorf("Failed to fetch schedules: %v", err)
		return nil, constants.ErrGetAllSchedule
	}

	utils.Log.Infof("Fetched %d schedules", len(schedules))

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
	utils.Log.Infof("Fetching schedule by ID: %s", scheduleID)

	loc := helpers.GetAppLocation()

	if _, err := uuid.Parse(scheduleID); err != nil {
		utils.Log.Errorf("Invalid schedule UUID: %v", err)
		return dto.ScheduleFullResponse{}, constants.ErrInvalidUUID
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, scheduleID)
	if err != nil {
		utils.Log.Errorf("Schedule not found: %v", err)
		return dto.ScheduleFullResponse{}, constants.ErrScheduleNotFound
	}

	utils.Log.Infof("Schedule found: %s", scheduleID)

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
	utils.Log.Infof("Updating schedule: %s", req.ScheduleID)

	loc := helpers.GetAppLocation()

	if _, err := uuid.Parse(req.ScheduleID); err != nil {
		utils.Log.Errorf("Invalid schedule UUID: %v", err)
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, req.ScheduleID)
	if err != nil {
		utils.Log.Errorf("Schedule not found: %v", err)
		return dto.ScheduleResponse{}, constants.ErrGetScheduleByID
	}

	if req.FieldID != nil {
		parsedUUID, err := uuid.Parse(*req.FieldID)
		if err != nil {
			utils.Log.Errorf("Invalid field UUID: %v", err)
			return dto.ScheduleResponse{}, constants.ErrInvalidUUID
		}
		schedule.FieldID = parsedUUID
	}

	if req.DayOfWeek != nil {
		if *req.DayOfWeek < 0 || *req.DayOfWeek > 6 {
			utils.Log.Warn("Invalid day of week")
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
		utils.Log.Warn("Close time must be after open time")
		return dto.ScheduleResponse{}, constants.ErrCloseTimeMustAfterOpen
	}

	if err := ss.scheduleRepo.UpdateSchedule(ctx, nil, schedule); err != nil {
		utils.Log.Errorf("Failed to update schedule: %v", err)
		return dto.ScheduleResponse{}, constants.ErrUpdateSchedule
	}

	utils.Log.Infof("Schedule updated successfully: %s", req.ScheduleID)

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
	utils.Log.Infof("Deleting schedule: %s", req.ScheduleID)

	loc := helpers.GetAppLocation()

	if _, err := uuid.Parse(req.ScheduleID); err != nil {
		utils.Log.Errorf("Invalid schedule UUID: %v", err)
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, req.ScheduleID)
	if err != nil {
		utils.Log.Errorf("Schedule not found: %v", err)
		return dto.ScheduleResponse{}, constants.ErrScheduleNotFound
	}

	if err := ss.scheduleRepo.DeleteScheduleByID(ctx, nil, req.ScheduleID); err != nil {
		utils.Log.Errorf("Failed to delete schedule: %v", err)
		return dto.ScheduleResponse{}, constants.ErrDeleteSchedule
	}

	utils.Log.Infof("Schedule deleted successfully: %s", req.ScheduleID)

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
	utils.Log.Infof("Fetching schedules by field ID: %s", fieldID)

	loc := helpers.GetAppLocation()

	if _, err := uuid.Parse(fieldID); err != nil {
		utils.Log.Errorf("Invalid field UUID: %v", err)
		return nil, constants.ErrInvalidUUID
	}

	schedules, err := ss.scheduleRepo.GetSchedulesByFieldID(ctx, nil, fieldID)
	if err != nil {
		utils.Log.Errorf("Failed to get schedules by field ID: %v", err)
		return nil, constants.ErrGetAllSchedule
	}

	utils.Log.Infof("Found %d schedules for field ID: %s", len(schedules), fieldID)

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
	utils.Log.Infof("Fetching schedule for field ID: %s and day: %d", fieldID, day)

	loc := helpers.GetAppLocation()

	if _, err := uuid.Parse(fieldID); err != nil {
		utils.Log.Errorf("Invalid field UUID: %v", err)
		return dto.ScheduleResponse{}, constants.ErrInvalidUUID
	}

	s, err := ss.scheduleRepo.GetScheduleByFieldIDAndDay(ctx, nil, fieldID, day)
	if err != nil {
		utils.Log.Errorf("Schedule not found for field ID %s and day %d: %v", fieldID, day, err)
		return dto.ScheduleResponse{}, constants.ErrScheduleNotFound
	}

	utils.Log.Infof("Schedule found: %s", s.ScheduleID)

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
