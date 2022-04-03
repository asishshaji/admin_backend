package admin_service

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/asishshaji/admin-api/models"
	admin_repository "github.com/asishshaji/admin-api/repositories"
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
	notificationService notification_service.INotificationService
}

func NewAdminService(l *log.Logger, adminRepo admin_repository.IAdminRepository, rClient *redis.Client, notification notification_service.INotificationService) IAdminService {
	return AdminService{
		l:         l,
		adminRepo: adminRepo,
		rClient:   rClient,

		notificationService: notification,
	}
}

func (aS AdminService) Login(ctx context.Context, username, password string) (string, error) {

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
		aS.l.Println(err)
	}

	title := "Your task is " + status.String()
	content := "Dolor aliquip labore incididunt id dolore culpa. Sunt velit cillum magna ad. Voluptate est cupidatat excepteur ut voluptate. Et anim incididunt cillum occaecat anim laboris ullamco dolor."

	msg := models.NotificationMessage{
		UserToken: tK.Token,
		Heading:   map[string]string{"en": title},
		Contents:  map[string]string{"en": content},
	}

	err = aS.notificationService.SendNotification(ctx, msg)
	if err != nil {
		aS.l.Println(err)
	}

	err = aS.adminRepo.CreateNotification(ctx, models.NotificationEntity{
		UserId:    uid,
		Content:   content,
		Title:     title,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Image:     "",
	})

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

func (aS AdminService) GetData(ctx context.Context) (models.Data, error) {
	data := models.Data{}

	res, err := aS.rClient.Get(ctx, "static_data").Result()
	if err != redis.Nil {
		aS.l.Println("Getting data from cache")

		err = json.Unmarshal([]byte(res), &data)
		if err != nil {
			return data, err
		}
		return data, nil
	}

	domainEntities, err := aS.adminRepo.GetDomains(ctx)
	if err != nil {
		return data, err
	}
	collegeEntities, err := aS.adminRepo.GetColleges(ctx)
	if err != nil {
		return data, err
	}
	courseEntities, err := aS.adminRepo.GetCourses(ctx)
	if err != nil {
		return data, err
	}

	for _, v := range domainEntities {
		data.Domains = append(data.Domains, v.Name)
	}

	for _, v := range collegeEntities {
		data.Colleges = append(data.Colleges, v.Name)
	}

	for _, v := range courseEntities {
		data.Courses = append(data.Courses, v.Name)

	}

	b, err := json.Marshal(&data)
	if err != nil {
		return data, err
	}

	r := aS.rClient.Set(ctx, "static_data", b, time.Hour*2)
	if r.Err() != nil {
		return data, r.Err()
	}

	return data, nil
}
