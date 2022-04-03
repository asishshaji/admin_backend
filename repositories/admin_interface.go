package admin_repository

import (
	"context"

	"github.com/asishshaji/admin-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdminRepository interface {
	GenerateAdminCredentials(ctx context.Context, username, password string) error

	GetAdmin(ctx context.Context, username string) (*models.Admin, error)
	AddTask(ctx context.Context, task models.Task) error
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskId primitive.ObjectID) error
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetUsers(ctx context.Context) (models.Students, error)
	GetTaskSubmissions(c context.Context) ([]models.TaskSubmissionsAdminResponse, error)
	GetTaskSubmissionsForUser(c context.Context, userid primitive.ObjectID) ([]models.TaskSubmissionsAdminResponse, error)
	EditTaskSubmissionStatus(c context.Context, status models.Status, taskid primitive.ObjectID) error

	CreateMentor(c context.Context, mentor models.Mentor) error
	UpdateMentor(c context.Context, mentor models.Mentor) error
	GetMentors(c context.Context) ([]models.Mentor, error)

	CreateDomain(c context.Context, domain models.StaticModel) error
	CreateCollege(c context.Context, college models.StaticModel) error
	CreateCourse(c context.Context, course models.StaticModel) error
	GetToken(c context.Context, uid primitive.ObjectID) (models.Token, error)

	GetDomains(ctx context.Context) ([]models.StaticModel, error)
	GetColleges(ctx context.Context) ([]models.StaticModel, error)
	GetCourses(ctx context.Context) ([]models.StaticModel, error)

	CreateNotification(ctx context.Context, notification models.NotificationEntity) error
}
