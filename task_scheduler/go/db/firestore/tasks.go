package firestore

import (
	"context"
	"fmt"
	"time"

	fs "cloud.google.com/go/firestore"
	"go.skia.org/infra/go/firestore"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
	"go.skia.org/infra/task_scheduler/go/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	COLLECTION_TASKS = "tasks"
)

// Fix all timestamps for the given task.
func fixTaskTimestamps(task *db.Task) {
	task.Created = fixTimestamp(task.Created)
	task.DbModified = fixTimestamp(task.DbModified)
	task.Finished = fixTimestamp(task.Finished)
	task.Started = fixTimestamp(task.Started)
}

// tasks returns a reference to the tasks collection.
func (d *firestoreDB) tasks() *fs.CollectionRef {
	return d.client.Collection(COLLECTION_TASKS)
}

// See documentation for db.TaskReader interface.
func (d *firestoreDB) GetTaskById(id string) (*db.Task, error) {
	doc, err := firestore.Get(d.tasks().Doc(id), DEFAULT_ATTEMPTS, GET_SINGLE_TIMEOUT)
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var rv db.Task
	if err := doc.DataTo(&rv); err != nil {
		return nil, err
	}
	return &rv, nil
}

// See documentation for db.TaskReader interface.
func (d *firestoreDB) GetTasksFromDateRange(start, end time.Time, repo string) ([]*db.Task, error) {
	// Adjust start and end times for Firestore resolution.
	start = fixTimestamp(start)
	end = fixTimestamp(end)

	// TODO(borenet): We can make this part of the query,
	// but it would required building a composite index.
	// It's possible that we should require a repo when
	// searching by timestamp; indexing a timestamp causes
	// the whole collection to cap out at a maximum of
	// 500 writes per second.
	q := d.tasks().Where("Created", ">=", start).Where("Created", "<", end).OrderBy("Created", fs.Asc)
	rv := []*db.Task{}
	if err := firestore.IterDocs(q, DEFAULT_ATTEMPTS, GET_MULTI_TIMEOUT, func(doc *fs.DocumentSnapshot) error {
		var task db.Task
		if err := doc.DataTo(&task); err != nil {
			return err
		}
		if repo == "" || task.Repo == repo {
			rv = append(rv, &task)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return rv, nil
}

// See documentation for db.TaskDB interface.
func (d *firestoreDB) AssignId(task *db.Task) error {
	task.Id = d.tasks().NewDoc().ID
	return nil
}

// putTasks sets the contents of the given tasks in Firestore, as part of the
// given transaction. It is used by PutTask and PutTasks.
func (d *firestoreDB) putTasks(tasks []*db.Task, tx *fs.Transaction) (rvErr error) {
	// Set the modification time of the tasks.
	now := fixTimestamp(time.Now())
	isNew := make([]bool, len(tasks))
	prevModified := make([]time.Time, len(tasks))
	for idx, task := range tasks {
		if util.TimeIsZero(task.Created) {
			return fmt.Errorf("Created not set. Task %s created time is %s. %v", task.Id, task.Created, task)
		}
		isNew[idx] = util.TimeIsZero(task.DbModified)
		prevModified[idx] = task.DbModified
		task.DbModified = now
	}
	defer func() {
		if rvErr != nil {
			for idx, task := range tasks {
				task.DbModified = prevModified[idx]
			}
		}
	}()

	// Find the previous versions of the tasks. Ensure that they weren't
	// updated concurrently.
	refs := make([]*fs.DocumentRef, 0, len(tasks))
	for _, task := range tasks {
		refs = append(refs, d.tasks().Doc(task.Id))
	}
	docs, err := tx.GetAll(refs)
	if err != nil {
		return err
	}
	for idx, doc := range docs {
		if !doc.Exists() {
			// This is expected for new tasks.
			if !isNew[idx] {
				sklog.Errorf("Task is not new but wasn't found in the DB.")
				// If the task is supposed to exist but does not, then
				// we have a problem.
				return db.ErrConcurrentUpdate
			}
		} else if isNew[idx] {
			// If the task is not supposed to exist but does, then
			// we have a problem.
			return db.ErrConcurrentUpdate
		}
		// If the task already exists, check the DbModified timestamp
		// to ensure that someone else didn't update it.
		if !isNew[idx] {
			var old db.Task
			if err := doc.DataTo(&old); err != nil {
				return err
			}
			if old.DbModified != prevModified[idx] {
				return db.ErrConcurrentUpdate
			}
		}
	}

	// Set the new contents of the tasks.
	for _, task := range tasks {
		ref := d.tasks().Doc(task.Id)
		if err := tx.Set(ref, task); err != nil {
			return err
		}
	}
	return nil
}

// See documentation for db.TaskDB interface.
func (d *firestoreDB) PutTask(task *db.Task) error {
	return d.PutTasks([]*db.Task{task})
}

// See documentation for db.TaskDB interface.
func (d *firestoreDB) PutTasks(tasks []*db.Task) error {
	if len(tasks) > firestore.MAX_TRANSACTION_DOCS/2 {
		sklog.Errorf("Inserting %d tasks; Firestore maximum per transaction is %d", len(tasks), firestore.MAX_TRANSACTION_DOCS)
	}
	for _, task := range tasks {
		if task.Id == "" {
			if err := d.AssignId(task); err != nil {
				return err
			}
		}
		fixTaskTimestamps(task)
	}

	if err := firestore.RunTransaction(d.client, DEFAULT_ATTEMPTS, PUT_MULTI_TIMEOUT, func(ctx context.Context, tx *fs.Transaction) error {
		return d.putTasks(tasks, tx)
	}); err != nil {
		return err
	}
	for _, task := range tasks {
		d.TrackModifiedTask(task)
	}
	return nil
}
