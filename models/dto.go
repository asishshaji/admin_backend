package models

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MentorDTO struct {
	Id           string `json:"_id"`
	Name         string `validate:"required"`
	Title        string `validate:"required"`
	Organization string `validate:"required"`
	Image        string `validate:"required"`
	Domain       string `validate:"required"`
}

func (mentor MentorDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(mentor)
}

func (dto MentorDTO) ToMentor() Mentor {
	id, _ := primitive.ObjectIDFromHex(dto.Id)
	return Mentor{
		ID:           id,
		Name:         dto.Name,
		Title:        dto.Title,
		Organization: dto.Organization,
		Image:        dto.Image,
		Domain:       dto.Domain,
	}
}

type TaskSubmissionDTO struct {
	TaskId  string `json:"task_id"`
	Comment string `json:"comment"`
	FileURL string `json:"file_url"`
}

type TaskDTO struct {
	ID       string `validate:"required"`
	Semester string `json:"semester" validate:"required"`
	Type     string `json:"type" validate:"required"`  // TYPE CAN BE RETAIL, ED-Tech
	Title    string `json:"title" validate:"required"` // title of task
	Detail   string `json:"detail" validate:"required"`
}

func (tD TaskDTO) ToTask() Task {
	return Task{
		Semester: tD.Semester,
		Type:     tD.Type,
		Title:    tD.Title,
		Detail:   tD.Detail,
	}
}

func (task *TaskDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(task)
}

type StudentDTO struct {
	Username         string       `json:"username" validate:"required"`
	FirstName        string       `json:"first_name" validate:"required"`
	PreferedType     PreferedType `json:"type"`
	LastName         string       `json:"last_name" validate:"required"`
	MiddleName       string       `json:"middle_name"`
	Password         string       `json:"password" validate:"required,min=4"`
	DOB              string       `json:"dob" validate:"required"`
	Gender           Gender       `json:"gender" validate:"required"`
	PhoneNumber      string       `json:"phone_number" validate:"required"`
	PhoneNumberAlt   string       `json:"phone_number_alt"`
	College          string       `json:"college" validate:"required"`
	Course           string       `json:"course" validate:"required"`
	Specialization   string       `json:"specialization" validate:"required"`
	HasArrears       bool         `json:"has_arrears" validate:"required"`
	Place            string       `json:"place" validate:"required"`
	Semester         string       `json:"semester" validate:"required"`
	District         string       `json:"district" validate:"required"`
	State            string       `json:"state" validate:"required"`
	Country          string       `json:"country" validate:"required"`
	DateOfJoining    string       `json:"date_of_joining"`
	CourseEndingDate string       `json:"course_ending_date"`
}

func (Student *StudentDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(Student)
}

func (stu StudentDTO) ToStudent() Student {
	return Student{
		Username:         stu.Username,
		FirstName:        stu.FirstName,
		PreferedType:     stu.PreferedType,
		LastName:         stu.LastName,
		MiddleName:       stu.MiddleName,
		DOB:              stu.DOB,
		Gender:           Gender(stu.Gender),
		PhoneNumber:      stu.PhoneNumber,
		PhoneNumberAlt:   stu.PhoneNumberAlt,
		College:          stu.College,
		Course:           stu.Course,
		Specialization:   stu.Specialization,
		HasArrears:       stu.HasArrears,
		Place:            stu.Place,
		Semester:         stu.Semester,
		Password:         stu.Password,
		District:         stu.District,
		State:            stu.State,
		Country:          stu.Country,
		DateOfJoining:    stu.DateOfJoining,
		CourseEndingDate: stu.CourseEndingDate,
	}
}
