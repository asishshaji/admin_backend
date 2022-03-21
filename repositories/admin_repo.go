package admin_repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/asishshaji/admin-api/models"
	"github.com/asishshaji/admin-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminRepository struct {
	l                        *log.Logger
	adminCollection          *mongo.Collection
	taskCollection           *mongo.Collection
	typeCollection           *mongo.Collection
	studentCollection        *mongo.Collection
	taskSubmissionCollection *mongo.Collection
	mentorCollection         *mongo.Collection
	domainCollection         *mongo.Collection
	collegeCollection        *mongo.Collection
	tokenCollection          *mongo.Collection

	courseCollection *mongo.Collection
}

func NewAdminRepository(l *log.Logger, db *mongo.Database) IAdminRepository {

	return AdminRepository{
		l:                        l,
		adminCollection:          db.Collection("admin"),
		mentorCollection:         db.Collection("mentor"),
		taskCollection:           db.Collection("tasks"),
		typeCollection:           db.Collection("types"),
		studentCollection:        db.Collection("students"),
		tokenCollection:          db.Collection("tokens"),
		taskSubmissionCollection: db.Collection("task_submission"),
		domainCollection:         db.Collection("domains"),
		collegeCollection:        db.Collection("colleges"),
		courseCollection:         db.Collection("courses"),
	}
}
func (aR AdminRepository) GenerateAdminCredentials(ctx context.Context, username, password string) error {
	opts := options.Update().SetUpsert(true)

	res, err := aR.adminCollection.UpdateOne(ctx, bson.M{
		"username": username,
	}, bson.M{
		"$set": bson.M{
			"username": username,
			"password": password,
		},
	}, opts)

	fmt.Println(res.MatchedCount)

	if err != nil {
		return err
	}

	return nil
}

func (adminRepo AdminRepository) GetAdmin(ctx context.Context, username string) (*models.Admin, error) {
	admin := new(models.Admin)

	res := adminRepo.adminCollection.FindOne(ctx, bson.M{"username": username})

	if res.Err() == mongo.ErrNoDocuments {
		adminRepo.l.Println("No admin with username", username, "exists")
		return nil, res.Err()
	}

	err := res.Decode(admin)

	if err != nil {
		adminRepo.l.Println(err)
		return nil, err
	}

	return admin, nil
}

func (aR AdminRepository) AddTask(ctx context.Context, task models.Task) error {

	res, err := aR.taskCollection.InsertOne(ctx, task)

	if err != nil {
		aR.l.Println(err)
		return err
	}

	aR.l.Println("Inserted task with ID", res.InsertedID)

	return nil
}

func (aR AdminRepository) UpdateTask(ctx context.Context, task models.Task) error {
	opts := options.Update().SetUpsert(true)

	up, err := utils.ToDoc(task)

	if err != nil {
		return err
	}

	doc := bson.M{"$set": up}

	_, err = aR.taskCollection.UpdateByID(ctx, task.Id, doc, opts)
	if err != nil {
		aR.l.Println(err)
		return err
	}
	return nil
}

func (aR AdminRepository) GetTasks(ctx context.Context) ([]models.Task, error) {
	tasks := []models.Task{}

	cursor, err := aR.taskCollection.Find(ctx, bson.M{})
	if err != nil {
		aR.l.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &tasks); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return tasks, nil
}

func (aR AdminRepository) CreateDomain(c context.Context, domain models.StaticModel) error {
	return insertStaticModelData(c, aR.domainCollection, domain)
}

func insertStaticModelData(c context.Context, collection *mongo.Collection, data models.StaticModel) error {
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(c, bson.M{
		"name": data.Name,
	}, bson.M{
		"$set": data,
	}, opts)

	if err != nil {
		return err
	}
	return nil

}

func (aR AdminRepository) CreateCollege(c context.Context, college models.StaticModel) error {
	return insertStaticModelData(c, aR.collegeCollection, college)
}
func (aR AdminRepository) CreateCourse(c context.Context, course models.StaticModel) error {
	return insertStaticModelData(c, aR.courseCollection, course)
}

func (aR AdminRepository) GetUsers(ctx context.Context) (models.Students, error) {
	students := new([]models.Student)
	cursor, err := aR.studentCollection.Find(ctx, bson.M{})
	if err != nil {
		aR.l.Println(err)
		return nil, err
	}
	if err = cursor.All(ctx, students); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return *students, nil
}
func (aR AdminRepository) DeleteTask(ctx context.Context, taskId primitive.ObjectID) error {
	res, err := aR.taskCollection.DeleteOne(ctx, bson.M{
		"_id": taskId,
	})
	if res.DeletedCount == 0 {
		aR.l.Println("No task found to be deleted")
		return errors.New("no task found with given id")

	}
	if err != nil {
		aR.l.Println(err)
		return err
	}
	return nil
}

func (aR AdminRepository) GetTaskSubmissions(c context.Context) ([]models.TaskSubmissionsAdminResponse, error) {

	lookupStage1 := bson.D{{
		"$lookup", bson.D{{
			"from", "students",
		}, {
			"localField", "userid",
		}, {
			"foreignField", "_id",
		}, {
			"as", "student",
		}},
	}}

	unwindStage1 := bson.D{{
		"$unwind", bson.D{{
			"path", "$student",
		}},
	}}

	projectStage1 := bson.D{{
		"$project", bson.D{
			{
				"student.username", 1,
			},
			{
				"student._id", 1,
			},
			{
				"_id", 1,
			},
			{
				"taskid", 1,
			},
			{
				"comment", 1,
			},
			{
				"fileurl", 1,
			},
			{
				"status", 1,
			},
			{
				"updatedat", 1,
			},
		},
	}}

	lookupStage2 := bson.D{{
		"$lookup", bson.D{{
			"from", "tasks",
		},
			{
				"localField", "taskid",
			},
			{
				"foreignField", "_id",
			},
			{
				"as", "task",
			},
		},
	}}

	unwindStage2 := bson.D{{
		"$unwind", bson.D{{
			"path", "$task",
		}},
	}}

	projectStage2 := bson.D{{
		"$project", bson.D{{
			"taskid", 0,
		}},
	}}

	cursor, err := aR.taskSubmissionCollection.Aggregate(c, mongo.Pipeline{lookupStage1, unwindStage1, projectStage1, lookupStage2, unwindStage2, projectStage2})
	if err != nil {
		aR.l.Println(err)
		return nil, err
	}
	var responseData []models.TaskSubmissionsAdminResponse
	if err = cursor.All(c, &responseData); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return responseData, nil
}

func (aR AdminRepository) EditTaskSubmissionStatus(c context.Context, status models.Status, taskid primitive.ObjectID) error {
	res, err := aR.taskSubmissionCollection.UpdateByID(c, taskid, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})

	if res.MatchedCount == 0 {
		aR.l.Println(models.ErrNoValidRecordFound)
		return models.ErrNoValidRecordFound
	}

	if err != nil {
		aR.l.Println(err)
		return err
	}
	return nil
}

func (aR AdminRepository) GetTaskSubmissionsForUser(c context.Context, userid primitive.ObjectID) ([]models.TaskSubmissionsAdminResponse, error) {
	filter := bson.D{{
		"$match", bson.D{{
			"userid", userid,
		}},
	}}
	lookupStage1 := bson.D{{
		"$lookup", bson.D{{
			"from", "students",
		}, {
			"localField", "userid",
		}, {
			"foreignField", "_id",
		}, {
			"as", "student",
		}},
	}}

	unwindStage1 := bson.D{{
		"$unwind", bson.D{{
			"path", "$student",
		}},
	}}

	projectStage1 := bson.D{{
		"$project", bson.D{
			{
				"student.username", 1,
			},
			{
				"student._id", 1,
			},
			{
				"_id", 1,
			},
			{
				"taskid", 1,
			},
			{
				"comment", 1,
			},
			{
				"fileurl", 1,
			},
			{
				"status", 1,
			},
		},
	}}

	lookupStage2 := bson.D{{
		"$lookup", bson.D{{
			"from", "tasks",
		},
			{
				"localField", "taskid",
			},
			{
				"foreignField", "_id",
			},
			{
				"as", "task",
			},
		},
	}}

	unwindStage2 := bson.D{{
		"$unwind", bson.D{{
			"path", "$task",
		}},
	}}

	projectStage2 := bson.D{{
		"$project", bson.D{{
			"taskid", 0,
		}},
	}}

	cursor, err := aR.taskSubmissionCollection.Aggregate(c, mongo.Pipeline{filter, lookupStage1, unwindStage1, projectStage1, lookupStage2, unwindStage2, projectStage2})
	if err != nil {
		aR.l.Println(err)
		return nil, err
	}
	var responseData []models.TaskSubmissionsAdminResponse
	if err = cursor.All(c, &responseData); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return responseData, nil

}

func (aR AdminRepository) CreateMentor(c context.Context, mentor models.Mentor) error {
	res, err := aR.mentorCollection.InsertOne(c, mentor)

	if mongo.IsDuplicateKeyError(err) {
		aR.l.Println("mentor already exists")
		return models.ErrMentorExists
	}

	if err != nil {
		aR.l.Println("Failed to create mentor")
		return err
	}

	aR.l.Println("Inserted mentor with ID", res.InsertedID)

	return nil
}

func (aR AdminRepository) UpdateMentor(c context.Context, mentor models.Mentor) error {
	opts := options.Update().SetUpsert(true)

	up, err := utils.ToDoc(mentor)

	if err != nil {
		return err
	}

	doc := bson.M{"$set": up}

	res, err := aR.mentorCollection.UpdateByID(c, mentor.ID, doc, opts)
	if err != nil {
		aR.l.Println(err)
		return err
	}
	aR.l.Println(res.MatchedCount)
	return nil
}

func (aR AdminRepository) GetMentors(c context.Context) ([]models.Mentor, error) {
	mentors := []models.Mentor{}

	cursor, err := aR.mentorCollection.Find(c, bson.M{})

	if err != nil {
		aR.l.Println(err)

		return nil, err
	}

	if err = cursor.All(c, &mentors); err != nil {
		aR.l.Println(err)
		return nil, err
	}

	return mentors, nil
}

func (aR AdminRepository) GetToken(ctx context.Context, uid primitive.ObjectID) (models.Token, error) {
	token := models.Token{}

	res := aR.tokenCollection.FindOne(ctx, bson.M{
		"user_id": uid,
	})

	if res.Err() != nil {
		return token, res.Err()
	}

	err := res.Decode(&token)
	if err != nil {
		aR.l.Println(err)
		return token, err
	}

	return token, nil

}
