package admin_controller

import "github.com/labstack/echo/v4"

type IAdminController interface {

	// users
	Login(c echo.Context) error
	GetUsers(c echo.Context) error

	CreateDomain(c echo.Context) error // create and update
	GetDomains(c echo.Context) error
	CreateCollege(c echo.Context) error
	CreateCourse(c echo.Context) error

	// Tasks
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	GetTasks(c echo.Context) error
	DeleteTask(c echo.Context) error

	// submissions
	GetTaskSubmissions(c echo.Context) error
	GetTaskSubmissionForUser(c echo.Context) error
	EditTaskSubmissionStatus(c echo.Context) error

	// mentors
	CreateMentor(c echo.Context) error
	UpdateMentor(c echo.Context) error
	GetMentors(c echo.Context) error
	// DeleteMentor(c echo.Context) error

	GetData(c echo.Context) error
	UploadFile(c echo.Context) error
}
