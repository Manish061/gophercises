package db

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/etcd-io/bbolt"
)

var taskBucket = []byte("TaskBucket")
var db *bbolt.DB
var taskLength = 0

// Task template for user tasks
type Task struct {
	Key       int
	Value     string
	Completed string
}

type taskDetail struct {
	Value     string
	Completed string
}

//Init initialises a db at path 'dbPath'
func Init(dbPath string) error {
	var err error
	db, err = bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

//AddTask add a task to the db
func AddTask(task string) (int, error) {
	var id = 0
	errUpdating := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		id64, _ := bucket.NextSequence()
		id = int(id64)
		key := itob(id)
		t := &taskDetail{Value: task, Completed: ""}
		out, errMarshalling := json.Marshal(t)
		if errMarshalling != nil {
			return errMarshalling
		}
		return bucket.Put(key, []byte(string(out)))
	})
	if errUpdating != nil {
		return -1, errUpdating
	}
	return id, nil
}

//AllTasks list all the user tasks
func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		t := taskDetail{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			errUnmarshalling := json.Unmarshal(v, &t)
			if errUnmarshalling != nil {
				return errUnmarshalling
			}
			if t.Completed != "" {
				timeC, errParsing := time.Parse("2006-01-02 15:04:05 -0700 MST", t.Completed)
				if errParsing != nil {
					return errParsing
				}
				tasks = append(tasks, Task{
					Key:       btoi(k),
					Value:     string(t.Value),
					Completed: timeC.String(),
				})
			} else {
				tasks = append(tasks, Task{
					Key:       btoi(k),
					Value:     string(t.Value),
					Completed: "",
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

//GetCompletedTasksForToday returns the list of tasks completed today
func GetCompletedTasksForToday() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		t := taskDetail{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			errUnmarshalling := json.Unmarshal(v, &t)
			if errUnmarshalling != nil {
				return errUnmarshalling
			}
			if t.Completed == "" {
				continue
			}
			timeC, errParsing := time.Parse("2006-01-02 15:04:05 -0700 MST", t.Completed)
			if errParsing != nil {
				return errParsing
			}
			today := time.Now().UTC()
			if timeC.Day() == today.Day() && timeC.Month() == today.Month() && timeC.Year() == today.Year() {
				tasks = append(tasks, Task{
					Key:   btoi(k),
					Value: string(t.Value),
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

//DeleteTask deletes a particular user task
func DeleteTask(key int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

// CompleteTask updates the time when the task is completed
func CompleteTask(key int, prevValue string, updateTime time.Time) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		t := &taskDetail{Value: prevValue, Completed: updateTime.String()}
		out, errMarshalling := json.Marshal(t)
		if errMarshalling != nil {
			return errMarshalling
		}
		return b.Put(itob(key), []byte(string(out)))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
