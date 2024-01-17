package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"time"
)

type IAttendanceRecordUsecase interface {
	GetRecordByDate(uint, time.Time) (model.AttendanceRecordResponse, error)
	GetRecordsByDate(date time.Time) ([]model.AttendanceRecordResponse, error)
	GetRecordsByDepartment(department string) ([]model.AttendanceRecordResponse, error)
	GetRecordsByDateDepartment(date time.Time, department string) ([]model.AttendanceRecordResponse, error)
	GetAllRecords(userId uint) ([]model.AttendanceRecordResponse, error)
	GetRecordById(userId uint, recordId uint) (model.AttendanceRecordResponse, error)
	GetAllUsers() ([]model.UserResponse, error)
	CreateRecord(record model.AttendanceRecord) (model.AttendanceRecordResponse, error)
	UpdateRecord(record model.AttendanceRecord, userId uint, recordId uint) (model.AttendanceRecordResponse, error)
	DeleteRecord(userId uint, recordId uint) error
}

type attendanceRecordUsecase struct {
	ar repository.IAttendanceRecordRepository
	av validator.IAttendanceRecordValidator
}

func NewAttendanceRecordUsecase(ar repository.IAttendanceRecordRepository, av validator.IAttendanceRecordValidator) IAttendanceRecordUsecase {
	return &attendanceRecordUsecase{ar, av}
}

func (aru *attendanceRecordUsecase) GetRecordByDate(userId uint, date time.Time) (model.AttendanceRecordResponse, error) {
	record := model.AttendanceRecord{}
	if err := aru.ar.GetRecordByDate(&record, userId, date); err != nil {
		return model.AttendanceRecordResponse{}, err
	}
	return model.AttendanceRecordResponse{
		ID:           record.ID,
		UserID:       record.UserID,
		ClockInTime:  record.ClockInTime,
		ClockOutTime: record.ClockOutTime,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}, nil
}
func (aru *attendanceRecordUsecase) GetRecordsByDate(date time.Time) ([]model.AttendanceRecordResponse, error) {
	var records []model.AttendanceRecord

	if err := aru.ar.GetRecordsByDate(&records, date); err != nil {
		return nil, err
	}

	var responses []model.AttendanceRecordResponse
	for _, record := range records {
		responses = append(responses, model.AttendanceRecordResponse{
			ID:           record.ID,
			UserID:       record.UserID,
			ClockInTime:  record.ClockInTime,
			ClockOutTime: record.ClockOutTime,
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		})
	}
	println("usecase GetRecordsByDate")
	return responses, nil
}

func (aru *attendanceRecordUsecase) GetRecordsByDepartment(department string) ([]model.AttendanceRecordResponse, error) {
	var records []model.AttendanceRecord

	if err := aru.ar.GetRecordsByDepartment(&records, department); err != nil {
		return nil, err
	}

	var responses []model.AttendanceRecordResponse
	for _, record := range records {
		responses = append(responses, model.AttendanceRecordResponse{
			ID:           record.ID,
			UserID:       record.UserID,
			ClockInTime:  record.ClockInTime,
			ClockOutTime: record.ClockOutTime,
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		})
	}
	return responses, nil
}

func (aru *attendanceRecordUsecase) GetRecordsByDateDepartment(date time.Time, department string) ([]model.AttendanceRecordResponse, error) {
	var records []model.AttendanceRecord

	if err := aru.ar.GetRecordsByDateDepartment(&records, date, department); err != nil {
		return nil, err
	}

	var responses []model.AttendanceRecordResponse
	for _, record := range records {
		userResponse := model.UserResponse{
			ID:         record.User.ID,
			Email:      record.User.Email,
			Department: record.User.Department,
			Name:       record.User.Name,
		}

		responses = append(responses, model.AttendanceRecordResponse{
			ID:           record.ID,
			UserID:       record.UserID,
			ClockInTime:  record.ClockInTime,
			ClockOutTime: record.ClockOutTime,
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
			User:         userResponse,
		})
	}
	return responses, nil
}
func (aru *attendanceRecordUsecase) GetAllRecords(userId uint) ([]model.AttendanceRecordResponse, error) {
	records := []model.AttendanceRecord{}
	if err := aru.ar.GetAllRecords(&records, userId); err != nil {
		return nil, err
	}
	resRecords := make([]model.AttendanceRecordResponse, len(records))
	for i, v := range records {
		resRecords[i] = model.AttendanceRecordResponse{
			ID:           v.ID,
			UserID:       v.UserID,
			ClockInTime:  v.ClockInTime,
			ClockOutTime: v.ClockOutTime,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		}
	}
	return resRecords, nil
}

func (aru *attendanceRecordUsecase) GetAllUsers() ([]model.UserResponse, error) {
	users := []model.User{}
	if err := aru.ar.GetAllUsers(&users); err != nil {
		return nil, err
	}
	resUsers := make([]model.UserResponse, len(users))
	for i, v := range users {
		resUsers[i] = model.UserResponse{
			ID:         v.ID,
			Email:      v.Email,
			Department: v.Department,
			Name:       v.Name,
		}
	}
	return resUsers, nil
}
func (aru *attendanceRecordUsecase) GetRecordById(userId uint, recordId uint) (model.AttendanceRecordResponse, error) {
	record := model.AttendanceRecord{}
	if err := aru.ar.GetRecordById(&record, userId, recordId); err != nil {
		return model.AttendanceRecordResponse{}, err
	}
	return model.AttendanceRecordResponse{
		ID:           record.ID,
		UserID:       record.UserID,
		ClockInTime:  record.ClockInTime,
		ClockOutTime: record.ClockOutTime,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}, nil
}

func (aru *attendanceRecordUsecase) CreateRecord(record model.AttendanceRecord) (model.AttendanceRecordResponse, error) {

	if err := aru.av.ValidateClockIn(record); err != nil {
		return model.AttendanceRecordResponse{}, err
	}
	if err := aru.ar.CreateRecord(&record); err != nil {
		return model.AttendanceRecordResponse{}, err
	}
	return model.AttendanceRecordResponse{
		ID:           record.ID,
		UserID:       record.UserID,
		ClockInTime:  record.ClockInTime,
		ClockOutTime: record.ClockOutTime,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}, nil
}

func (aru *attendanceRecordUsecase) UpdateRecord(record model.AttendanceRecord, userId uint, recordId uint) (model.AttendanceRecordResponse, error) {

	if err := aru.av.ValidateClockOut(record); err != nil {
		return model.AttendanceRecordResponse{}, err
	}
	if err := aru.ar.UpdateRecord(&record, userId, recordId); err != nil {
		return model.AttendanceRecordResponse{}, err
	}
	return model.AttendanceRecordResponse{
		ID:           record.ID,
		UserID:       record.UserID,
		ClockInTime:  record.ClockInTime,
		ClockOutTime: record.ClockOutTime,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}, nil
}

func (aru *attendanceRecordUsecase) DeleteRecord(userId uint, recordId uint) error {
	return aru.ar.DeleteRecord(userId, recordId)
}
