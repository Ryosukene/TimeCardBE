package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
)

func main() {
	db := db.NewDB()
	userValidator := validator.NewUserValidator()
	authUserValidator := validator.NewAuthUserValidator() // AuthUser用のバリデーターを追加
	taskValidator := validator.NewTaskValidator()
	attendanceRecordValidator := validator.NewAttendanceRecordValidator()

	userRepository := repository.NewUserRepository(db)
	authUserRepository := repository.NewAuthUserRepository(db) // AuthUser用のリポジトリを追加
	taskRepository := repository.NewTaskRepository(db)
	attendanceRecordRepository := repository.NewAttendanceRecordRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	authUserUsecase := usecase.NewAuthUserUsecase(authUserRepository, authUserValidator) // AuthUser用のユースケースを追加
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	attendanceRecordUsecase := usecase.NewAttendanceRecordUsecase(attendanceRecordRepository, attendanceRecordValidator)

	userController := controller.NewUserController(userUsecase)
	authUserController := controller.NewAuthUserController(authUserUsecase) // AuthUser用のコントローラーを追加
	taskController := controller.NewTaskController(taskUsecase)
	attendanceRecordController := controller.NewAttendanceRecordController(attendanceRecordUsecase)

	e := router.NewRouter(userController, authUserController, taskController, attendanceRecordController) // ルーターにAuthUserコントローラーを追加
	e.Logger.Fatal(e.Start(":8080"))
}
