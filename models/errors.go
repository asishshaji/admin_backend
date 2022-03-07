package models

import "fmt"

var ErrStudentExists = fmt.Errorf("student already exists")
var ErrInvalidCredentials = fmt.Errorf("invalid credentials")
var ErrNoStudentExists = fmt.Errorf("no student with given studentname")
var ErrNoStudentWithIdExists = fmt.Errorf("no Student with id exists")
var ErrParsingStudent = fmt.Errorf("error parsing student data from database")

var ErrNoAdminWithUsername = fmt.Errorf("no admin with username exists")

var ErrMentorExists = fmt.Errorf("mentor already exists")

var ErrNoValidRecordFound = fmt.Errorf("no valid document found")
var ErrTaskSubmissionExists = fmt.Errorf("task submission already exists")
