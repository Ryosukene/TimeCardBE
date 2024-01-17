package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, auc controller.IAuthUserController, tc controller.ITaskController, arc controller.IAttendanceRecordController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteNoneMode,
		CookieSameSite: http.SameSiteDefaultMode,
		//CookieMaxAge:   60,
	}))
	e.POST("/signup", uc.SignUp)
	e.POST("/create-user", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.PUT("/update-user", uc.UpdateUser)
	e.DELETE("/delete-user", uc.DeleteUser)

	e.POST("/auth/signup", auc.SignUp)
	e.POST("/auth/login", auc.LogIn)
	e.POST("/auth/logout", auc.LogOut)

	e.GET("/csrf", uc.CsrfToken)
	t := e.Group("/tasks")
	t.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)

	ar := e.Group("/attendance-records")
	ar.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))

	ar.GET("", arc.GetAllRecords)
	ar.GET("/:recordId", arc.GetRecordById)
	ar.GET("/date/:date", arc.GetRecordByDate)

	ar.POST("", arc.CreateRecord)
	ar.POST("/clock-in", arc.ClockIn)
	ar.POST("/clock-out", arc.ClockOut)
	ar.PUT("/:recordId", arc.UpdateRecord)
	ar.DELETE("/:recordId", arc.DeleteRecord)

	ar2 := e.Group("/adminrecords")
	ar2.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	ar2.GET("/date", arc.GetRecordsByDate)
	ar2.GET("/department", arc.GetRecordsByDepartment)
	ar2.GET("/date-department", arc.GetRecordsByDateDepartment)
	ar2.GET("/users", arc.GetAllUsers)

	return e
}
