package admin_service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/asishshaji/admin-api/models"
	admin_repository "github.com/asishshaji/admin-api/repositories"
	file_service "github.com/asishshaji/admin-api/services/file"
	"github.com/asishshaji/admin-api/services/notification_service"
	"github.com/asishshaji/admin-api/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	l                   *log.Logger
	adminRepo           admin_repository.IAdminRepository
	rClient             *redis.Client
	imageService        file_service.IFileService
	notificationService notification_service.INotificationService
}

func NewAdminService(l *log.Logger, adminRepo admin_repository.IAdminRepository, rClient *redis.Client, fileService file_service.IFileService, notification notification_service.INotificationService) IAdminService {
	return AdminService{
		l:                   l,
		adminRepo:           adminRepo,
		rClient:             rClient,
		imageService:        fileService,
		notificationService: notification,
	}
}

func (aS AdminService) Login(ctx context.Context, username, password string) (string, error) {

	var token string = "16000112-a6ab-11ec-abda-ee142f97fd44"

	msg := models.NotificationMessage{
		UserToken: token,
		Heading:   map[string]string{"en": "Hi Welcome"},
		Contents:  map[string]string{"en": "Helloooo. Notification Body"},
	}

	aS.notificationService.SendNotification(ctx, msg)

	admin, err := aS.adminRepo.GetAdmin(ctx, username)

	if err != nil {
		return "", models.ErrNoAdminWithUsername
	}
	authenticate := utils.CheckpasswordHash(password, admin.Password)

	if !authenticate {
		return "", models.ErrInvalidCredentials
	}

	adminClaims := &models.AdminJWTClaims{
		admin.ID,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	tokenMethod := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)
	t, err := tokenMethod.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		aS.l.Println(err)
		return "", err
	}

	return t, nil
}

func (aS AdminService) AddTask(ctx context.Context, task models.TaskDTO, creatorID primitive.ObjectID) error {

	t := task.ToTask()
	t.CreatorID = creatorID
	t.Id = primitive.NewObjectIDFromTimestamp(time.Now())
	t.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	t.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := aS.adminRepo.AddTask(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (aS AdminService) UpdateTask(ctx context.Context, task models.TaskDTO) error {

	tId, _ := primitive.ObjectIDFromHex(task.ID)

	t := task.ToTask()
	t.Id = tId
	t.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := aS.adminRepo.UpdateTask(ctx, t)
	if err != nil {
		return err
	}
	return nil
}

func (aS AdminService) GetTasks(ctx context.Context) ([]models.Task, error) {
	return aS.adminRepo.GetTasks(ctx)
}

func (aS AdminService) GetUsers(ctx context.Context) ([]models.StudentResponse, error) {
	studentModels, err := aS.adminRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	studentResponse := studentModels.ToStudentResponse()

	return studentResponse, nil
}

func (aS AdminService) DeleteTask(c context.Context, taskId primitive.ObjectID) error {
	return aS.adminRepo.DeleteTask(c, taskId)
}

func (aS AdminService) GetTaskSubmissions(c context.Context) ([]models.TaskSubmissionsAdminResponse, error) {
	return aS.adminRepo.GetTaskSubmissions(c)
}
func (aS AdminService) EditTaskSubmission(ctx context.Context, uid primitive.ObjectID, taskId primitive.ObjectID, status models.Status) error {
	err := aS.adminRepo.EditTaskSubmissionStatus(ctx, status, taskId)
	if err != nil {
		return err
	}

	tK, err := aS.adminRepo.GetToken(ctx, uid)
	if err != nil {
		return err
	}

	msg := models.NotificationMessage{
		UserToken: tK.Token,
		Heading:   map[string]string{"en": "Your task is " + status.String()},
	}

	err = aS.notificationService.SendNotification(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

func (aS AdminService) GetTaskSubmissionsForUser(ctx context.Context, userId primitive.ObjectID) ([]models.TaskSubmissionsAdminResponse, error) {

	return aS.adminRepo.GetTaskSubmissionsForUser(ctx, userId)

}

func (aS AdminService) CreateMentor(ctx context.Context, mentor models.MentorDTO) error {

	m := mentor.ToMentor()
	m.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	m.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	m.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return aS.adminRepo.CreateMentor(ctx, m)
}

func (aS AdminService) UpdateMentor(ctx context.Context, mentor models.MentorDTO) error {

	m := mentor.ToMentor()
	m.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return aS.adminRepo.UpdateMentor(ctx, m)
}

func (aS AdminService) GetMentors(ctx context.Context) ([]models.MentorResponse, error) {

	mentors, err := aS.adminRepo.GetMentors(ctx)
	if err != nil {
		return nil, err
	}

	mentorResponses := []models.MentorResponse{}

	for _, dto := range mentors {
		mentorResponses = append(mentorResponses, *dto.ToResponse())
	}

	return mentorResponses, nil
}
func (aS AdminService) CreateDomain(ctx context.Context, domainString string) error {

	domain := models.StaticModel{
		Name:      domainString,
		CreatedOn: primitive.NewDateTimeFromTime(time.Now()),
	}

	return aS.adminRepo.CreateDomain(ctx, domain)
}

func (aS AdminService) CreateCollege(ctx context.Context, college string) error {
	c := models.StaticModel{
		Name:      college,
		CreatedOn: primitive.NewDateTimeFromTime(time.Now()),
	}

	return aS.adminRepo.CreateCollege(ctx, c)
}

func (aS AdminService) CreateCourse(ctx context.Context, course string) error {
	c := models.StaticModel{
		Name:      course,
		CreatedOn: primitive.NewDateTimeFromTime(time.Now()),
	}

	return aS.adminRepo.CreateCourse(ctx, c)
}
