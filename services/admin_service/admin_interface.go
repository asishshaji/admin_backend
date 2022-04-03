package admin_service

import (
	"context"

	"github.com/asishshaji/admin-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdminService interface {
	Login(ctx context.Context, username, password string) (string, error)
	AddTask(ctx context.Context, task models.TaskDTO, creatorID primitive.ObjectID) error
	UpdateTask(ctx context.Context, task models.TaskDTO) error
	DeleteTask(c context.Context, taskId primitive.ObjectID) error
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetUsers(ctx context.Context) ([]models.StudentResponse, error)
	GetTaskSubmissions(c context.Context) ([]models.TaskSubmissionsAdminResponse, error)
	EditTaskSubmission(ctx context.Context, uid primitive.ObjectID, taskId primitive.ObjectID, status models.Status) error
	GetTaskSubmissionsForUser(ctx context.Context, userId primitive.ObjectID) ([]models.TaskSubmissionsAdminResponse, error)

	CreateMentor(ctx context.Context, mentor models.MentorDTO) error
	UpdateMentor(ctx context.Context, mentor models.MentorDTO) error
	GetMentors(ctx context.Context) ([]models.MentorResponse, error)

	CreateDomain(ctx context.Context, domainString string) error
	CreateCollege(ctx context.Context, college string) error
	CreateCourse(ctx context.Context, course string) error

	GetData(ctx context.Context) (models.Data, error)
}
