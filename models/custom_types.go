package models

type Gender int32

func (g Gender) String() string {
	switch g {
	case MALE:
		return "male"
	case FEMALE:
		return "female"
	}
	return ""
}

const (
	MALE   Gender = 1
	FEMALE Gender = 2
)

type Status string

const (
	ACTIVE    Status = "active"    // submitted
	COMPLETED Status = "completed" // verified by admin
	INACTIVE  Status = "inactive"  // not started
	REJECTED  Status = "rejected"
)

func (s Status) String() string {
	switch s {
	case ACTIVE:
		return "active"
	case COMPLETED:
		return "completed"
	case INACTIVE:
		return "inactive"
	case REJECTED:
		return "rejected"
	}
	return ""
}
