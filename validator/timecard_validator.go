package validator

import (
	"go-rest-api/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IAttendanceRecordValidator interface {
	Validate(record model.AttendanceRecord) error
	ValidateClockIn(record model.AttendanceRecord) error
	ValidateClockOut(record model.AttendanceRecord) error
}

type attendanceRecordValidator struct{}

func NewAttendanceRecordValidator() IAttendanceRecordValidator {
	return &attendanceRecordValidator{}
}

func (arv *attendanceRecordValidator) Validate(record model.AttendanceRecord) error {
	return validation.ValidateStruct(&record,
		validation.Field(
			&record.ClockInTime,
			validation.Required.Error("clock-in time is required"),
			validation.Max(time.Now()).Error("clock-in time cannot be in the future"),
		),
		validation.Field(
			&record.ClockOutTime,
			validation.Required.Error("clock-out time is required"),
			validation.Max(time.Now()).Error("clock-out time cannot be in the future"),
			validation.By(func(value interface{}) error {
				if value.(time.Time).Before(record.ClockInTime) {
					return validation.NewError("validation", "clock-out time cannot be before clock-in time")
				}
				return nil
			}),
		),
	)
}
func (arv *attendanceRecordValidator) ValidateClockIn(record model.AttendanceRecord) error {
	return validation.ValidateStruct(&record,
		validation.Field(
			&record.ClockInTime,
			validation.Required.Error("clock-in time is required"),
			validation.Max(time.Now()).Error("clock-in time cannot be in the future"),
		),
	)
}
func (arv *attendanceRecordValidator) ValidateClockOut(record model.AttendanceRecord) error {
	return validation.ValidateStruct(&record,
		validation.Field(
			&record.ClockOutTime,
			validation.Required.Error("clock-out time is required"),
			validation.Max(time.Now()).Error("clock-out time cannot be in the future"),
			validation.By(func(value interface{}) error {
				if value.(time.Time).Before(record.ClockInTime) {
					return validation.NewError("validation", "clock-out time cannot be before clock-in time")
				}
				return nil
			}),
		),
	)
}
