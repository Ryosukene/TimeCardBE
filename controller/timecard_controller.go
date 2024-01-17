package controller

import (
	"fmt"
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type IAttendanceRecordController interface {
	ClockIn(c echo.Context) error
	ClockOut(c echo.Context) error
	GetAllRecords(c echo.Context) error
	GetRecordById(c echo.Context) error
	GetRecordByDate(c echo.Context) error
	GetRecordsByDate(c echo.Context) error
	GetRecordsByDepartment(c echo.Context) error
	GetRecordsByDateDepartment(c echo.Context) error
	GetAllUsers(c echo.Context) error
	CreateRecord(c echo.Context) error
	UpdateRecord(c echo.Context) error
	DeleteRecord(c echo.Context) error
}

type attendanceRecordController struct {
	aru usecase.IAttendanceRecordUsecase
}

func NewAttendanceRecordController(aru usecase.IAttendanceRecordUsecase) IAttendanceRecordController {
	return &attendanceRecordController{aru}
}
func (arc *attendanceRecordController) ClockIn(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	floatUserId, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "User ID is not a float64")
	}
	userId := uint(floatUserId)

	type ClockInRequest struct {
		ClockInTime time.Time `json:"clock_in_time"`
	}

	var req ClockInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	record := model.AttendanceRecord{
		UserID:      userId,
		ClockInTime: req.ClockInTime,
	}

	recordRes, err := arc.aru.CreateRecord(record)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, recordRes)
}
func (arc *attendanceRecordController) ClockOut(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	floatUserId, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "User ID is not a float64")
	}
	userId := uint(floatUserId)

	type ClockOutRequest struct {
		ClockOutTime time.Time `json:"clock_out_time"`
	}

	var req ClockOutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var recordRes model.AttendanceRecordResponse
	recordRes, err := arc.aru.GetRecordByDate(userId, req.ClockOutTime)
	var record model.AttendanceRecord = ConvertResponseToRecord(recordRes)

	record.ClockOutTime = req.ClockOutTime

	recordRes, err = arc.aru.UpdateRecord(model.AttendanceRecord(record), userId, uint(record.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recordRes)

}
func ConvertResponseToRecord(response model.AttendanceRecordResponse) model.AttendanceRecord {
	return model.AttendanceRecord{
		ID:           response.ID,
		UserID:       response.UserID,
		ClockInTime:  response.ClockInTime,
		ClockOutTime: response.ClockOutTime,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
	}
}

func (arc *attendanceRecordController) GetAllRecords(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// userId := claims["user_id"].(uint)
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	floatUserId, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "User ID is not a float64")
	}
	userId := uint(floatUserId)

	recordsRes, err := arc.aru.GetAllRecords(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recordsRes)
}

func (arc *attendanceRecordController) GetRecordsByDepartment(c echo.Context) error {
	// URL パラメータから部署を取得
	department := c.QueryParam("department")

	records, err := arc.aru.GetRecordsByDepartment(department)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, records)
}
func (arc *attendanceRecordController) GetRecordsByDateDepartment(c echo.Context) error {
	// URL パラメータから日付を取得
	dateParam := c.QueryParam("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid date format")
	}

	// URL パラメータから部署を取得
	department := c.QueryParam("department")

	records, err := arc.aru.GetRecordsByDateDepartment(date, department)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, records)
}

func (arc *attendanceRecordController) GetRecordById(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	floatUserId, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "User ID is not a float64")
	}
	userId := uint(floatUserId)

	id := c.Param("recordId")
	recordId, _ := strconv.Atoi(id)

	recordRes, err := arc.aru.GetRecordById(userId, uint(recordId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recordRes)
}

func (arc *attendanceRecordController) GetRecordByDate(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	floatUserId, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "User ID is not a float64")
	}
	userId := uint(floatUserId)

	date := c.Param("date")

	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(dateTime)
	recordRes, err := arc.aru.GetRecordByDate(userId, dateTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recordRes)
}
func (arc *attendanceRecordController) GetRecordsByDate(c echo.Context) error {
	// URL パラメータから日付を取得
	dateParam := c.QueryParam("date")
	date, err := time.Parse("2006-01-02", dateParam)
	fmt.Println(date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid date format")
	}

	records, err := arc.aru.GetRecordsByDate(date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, records)
}
func (arc *attendanceRecordController) GetAllUsers(c echo.Context) error {
	users, err := arc.aru.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
func (arc *attendanceRecordController) CreateRecord(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	floatUserId, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "User ID is not a float64")
	}
	userId := uint(floatUserId)

	record := model.AttendanceRecord{}
	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	record.UserID = userId

	recordRes, err := arc.aru.CreateRecord(record)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, recordRes)
}

func (arc *attendanceRecordController) UpdateRecord(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(uint)
	id := c.Param("recordId")
	recordId, _ := strconv.Atoi(id)

	record := model.AttendanceRecord{}
	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	recordRes, err := arc.aru.UpdateRecord(record, userId, uint(recordId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recordRes)
}

func (arc *attendanceRecordController) DeleteRecord(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(uint)
	id := c.Param("recordId")
	recordId, _ := strconv.Atoi(id)

	err := arc.aru.DeleteRecord(userId, uint(recordId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
