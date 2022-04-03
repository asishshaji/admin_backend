package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminJWTClaims struct {
	AdminId primitive.ObjectID
	IsAdmin bool
	jwt.StandardClaims
}

type Response struct {
	Message interface{} `json:"message"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `json:"username"`
	Password string             `json:"-"`
}

type Student struct {
	ID               primitive.ObjectID   `bson:"_id" json:"_id"`
	Email            string               `json:"email" validate:"required"`
	FirstName        string               `json:"first_name" validate:"required"`
	Domains          []string             `json:"domains,omitempty"`
	LastName         string               `json:"last_name" validate:"required"`
	MiddleName       string               `json:"middle_name"`
	CreatedAt        primitive.DateTime   `json:"-"`
	UpdatedAt        primitive.DateTime   `json:"-"`
	Password         string               `json:"password" validate:"required,min=4"`
	DOB              string               `json:"dob" validate:"required"`
	Gender           Gender               `json:"gender" validate:"required"`
	PhoneNumber      string               `json:"phone_number" validate:"required"`
	PhoneNumberAlt   string               `json:"phone_number_alt"`
	College          string               `json:"college" validate:"required"`
	Course           string               `json:"course" validate:"required"`
	Specialization   string               `json:"specialization" validate:"required"`
	HasArrears       bool                 `json:"has_arrears" validate:"required"`
	Place            string               `json:"place" validate:"required"`
	Semester         string               `json:"semester" validate:"required"`
	District         string               `json:"district" validate:"required"`
	State            string               `json:"state" validate:"required"`
	Country          string               `json:"country" validate:"required"`
	DateOfJoining    string               `json:"date_of_joining"`
	CourseEndingDate string               `json:"course_ending_date"`
	Mentors          []primitive.ObjectID `json:"-"`
}

type Students []Student

func (students Students) ToStudentResponse() []StudentResponse {
	studentReponse := []StudentResponse{}

	for _, stu := range students {
		studentReponse = append(studentReponse, StudentResponse{
			ID:               stu.ID,
			Email:            stu.Email,
			FirstName:        stu.FirstName,
			Domains:          stu.Domains,
			LastName:         stu.LastName,
			MiddleName:       stu.MiddleName,
			CreatedAt:        stu.CreatedAt,
			UpdatedAt:        stu.UpdatedAt,
			DOB:              stu.DOB,
			Gender:           Gender(stu.Gender).String(),
			PhoneNumber:      stu.PhoneNumber,
			PhoneNumberAlt:   stu.PhoneNumberAlt,
			College:          stu.College,
			Course:           stu.Course,
			Specialization:   stu.Specialization,
			HasArrears:       stu.HasArrears,
			Place:            stu.Place,
			Semester:         stu.Semester,
			District:         stu.District,
			State:            stu.State,
			Country:          stu.Country,
			DateOfJoining:    stu.DateOfJoining,
			CourseEndingDate: stu.CourseEndingDate,
		})
	}

	return studentReponse

}

type Mentor struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" validate:"required"`
	Title        string             `json:"title" validate:"required"`
	Organization string             `json:"organization" validate:"required"`
	Domain       string             `json:"domain" validate:"required"`
	Image        string             `json:"image"`
	Videos       []Videos           `bson:"videos"`
	CreatedAt    primitive.DateTime `bson:",omitempty"`
	UpdatedAt    primitive.DateTime `bson:",omitempty"`
}

func (mentor *Mentor) Validate() error {
	validate := validator.New()

	return validate.Struct(mentor)
}

func (dto *Mentor) ToResponse() *MentorResponse {

	return &MentorResponse{
		ID:           dto.ID,
		Name:         dto.Name,
		Title:        dto.Title,
		Organization: dto.Organization,
		Domain:       dto.Domain,
		CreatedAt:    dto.CreatedAt,
		Image:        dto.Image,
		Videos:       dto.Videos,
	}
}

type Task struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Semester  string             `json:"semester"`
	Domain    string             `json:"domain"` // TYPE CAN BE RETAIL, ED-Tech
	Title     string             `json:"title"`  // title of task
	Detail    string             `json:"detail"`
	CreatedAt primitive.DateTime `json:"created_at" bson:",omitempty"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:",omitempty"`
	CreatorID primitive.ObjectID `json:"creator_id"`
}

type TaskSubmission struct {
	UserId    primitive.ObjectID `json:"userid"`
	TaskId    primitive.ObjectID `json:"taskid"`
	Comment   string             `json:"comment"`
	FileURL   string             `json:"fileurl" bson:",omitempty"`
	Status    Status             `json:"status"`
	CreatedAt primitive.DateTime `bson:",omitempty"`
	UpdatedAt primitive.DateTime `bson:",omitempty"`
}

type StaticModel struct {
	Name      string
	CreatedOn primitive.DateTime `bson:"created_at"`
}

type TT struct {
	Mentors []primitive.ObjectID
}

type Token struct {
	UserId primitive.ObjectID `bson:"user_id"`
	Token  string             `bson:"token"`
}

type NotificationMessage struct {
	UserToken string
	Contents  map[string]string
	Heading   map[string]string
}

type NotificationEntity struct {
	Title     string
	Content   string
	Image     string
	CreatedAt primitive.DateTime `bson:"created_at"`
	UserId    primitive.ObjectID `bson:"user_id"`
}
