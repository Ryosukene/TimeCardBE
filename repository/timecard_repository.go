package repository

import (
	"fmt"
	"go-rest-api/model"
	"log"
	"time"

	"gorm.io/gorm"
)

type IAttendanceRecordRepository interface {
	GetRecordByDate(record *model.AttendanceRecord, userId uint, date time.Time) error
	GetAllRecords(records *[]model.AttendanceRecord, userId uint) error
	GetRecordById(record *model.AttendanceRecord, userId uint, recordId uint) error
	GetRecordsByDate(records *[]model.AttendanceRecord, date time.Time) error
	GetRecordsByDepartment(records *[]model.AttendanceRecord, department string) error
	GetRecordsByDateDepartment(records *[]model.AttendanceRecord, date time.Time, department string) error
	GetAllUsers(records *[]model.User) error
	CreateRecord(record *model.AttendanceRecord) error
	UpdateRecord(record *model.AttendanceRecord, userId uint, recordId uint) error
	DeleteRecord(userId uint, recordId uint) error
}

type attendanceRecordRepository struct {
	db *gorm.DB
}

func NewAttendanceRecordRepository(db *gorm.DB) IAttendanceRecordRepository {
	return &attendanceRecordRepository{db}
}

func (ar *attendanceRecordRepository) GetRecordByDate(record *model.AttendanceRecord, userId uint, date time.Time) error {
	// 日付のみを抽出（時間部分は無視）

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	// 指定された日付に一致するレコードを検索
	err := ar.db.Where("user_id = ? AND clock_in_time >= ? AND clock_in_time < ?", userId, dayStart, dayEnd).First(record).Error

	if err != nil {
		return err
	}
	return nil
}
func (ar *attendanceRecordRepository) GetRecordsByDate(records *[]model.AttendanceRecord, date time.Time) error {
	// 日付の開始時刻と終了時刻を計算
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	fmt.Println(dayStart)
	// 指定された日付に一致するレコードを検索
	err := ar.db.Where("clock_in_time >= ? AND clock_in_time < ?", dayStart, dayEnd).Find(records).Error
	if err != nil {
		return err
	}

	return nil
}
func (ar *attendanceRecordRepository) GetRecordsByDepartment(records *[]model.AttendanceRecord, department string) error {

	err := ar.db.Joins("join users on users.id = attendance_records.user_id").Where("users.department = ?", department).Find(records).Error
	if err != nil {
		return err
	}

	return nil
}

func (ar *attendanceRecordRepository) GetRecordsByDateDepartment(records *[]model.AttendanceRecord, date time.Time, department string) error {
	// 日付の開始時刻と終了時刻を計算
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	err := ar.db.Preload("User").Joins("join users on users.id = attendance_records.user_id").
		Where("attendance_records.clock_in_time >= ? AND attendance_records.clock_in_time < ? AND users.department = ?", dayStart, dayEnd, department).
		Find(records).Error

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	return nil
}

func (ar *attendanceRecordRepository) GetAllUsers(records *[]model.User) error {
	if err := ar.db.Order("id asc").Find(records).Error; err != nil {
		return err
	}
	return nil
}

func (ar *attendanceRecordRepository) GetAllRecords(records *[]model.AttendanceRecord, userId uint) error {
	if err := ar.db.Joins("User").Where("user_id = ?", userId).Order("created_at").Find(records).Error; err != nil {
		return err
	}
	return nil
}

func (ar *attendanceRecordRepository) GetRecordById(record *model.AttendanceRecord, userId uint, recordId uint) error {
	if err := ar.db.Joins("User").Where("attendance_records.user_id = ? AND attendance_records.id = ?", userId, recordId).First(record).Error; err != nil {
		return err
	}
	return nil
}

func (ar *attendanceRecordRepository) CreateRecord(record *model.AttendanceRecord) error {
	if err := ar.db.Create(record).Error; err != nil {
		return err
	}
	return nil
}

func (ar *attendanceRecordRepository) UpdateRecord(record *model.AttendanceRecord, userId uint, recordId uint) error {
	result := ar.db.Model(record).Where("id = ? AND user_id = ?", recordId, userId).Updates(record)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (ar *attendanceRecordRepository) DeleteRecord(userId uint, recordId uint) error {
	result := ar.db.Where("id = ? AND user_id = ?", recordId, userId).Delete(&model.AttendanceRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
