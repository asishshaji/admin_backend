package admin_controller

import "github.com/labstack/echo/v4"

type IAdminController interface {

	// users
	Login(c echo.Context) error
	GetUsers(c echo.Context) error

	CreateType(c echo.Context) error

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
}
