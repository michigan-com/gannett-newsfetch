package mongoqueue

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrTimeout                      = errors.New("timeout")
	ErrCannotCompleteWhenNotRunning = errors.New("cannot complete a job that's not running")
)

type Params struct {
	Logger Logger

	JobTimeout              time.Duration
	JobRetryMinDelay        time.Duration
	JobRetryMaxDelay        time.Duration
	JobRetryDelayMultiplier float64
	FailedJobExpiration     time.Duration

	JobRemovalAge time.Duration
}

type Queue struct {
	Collection *mgo.Collection
	Params
}

type WorkerFunc func(op string, args map[string]interface{}) error

type Worker struct {
	Func    WorkerFunc
	Timeout time.Duration
}

type JobState string

type PermanentError interface {
	IsPermanent() bool
}

const (
	Pending       JobState = "p"
	FailedPending          = "f"
	Running                = "r"
	Failed                 = "F"
	Succeeded              = "S"
	Canceled               = "C"
	Unknown                = "?"
)

func (s JobState) IsFinal() bool {
	return s == Succeeded || s == Canceled || s == Failed
}

type Job struct {
	// Name is a user-defined unique identifier of the job. An attempt to add a duplicate job with
	// the same name will be ignored; this makes enqueuing idempotent and safe to retry when
	// communication fails. bson.NewObjectId().Hex() is a good way to generate these.
	Name string `json:"name" bson:"_id"`

	// Op determines the specific operation to perform. A single queue can support multiple
	// different kinds of operations.
	Op string `json:"op" bson:"op"`

	// Args is a dictionary of arbitrary additional data. This can be used by the code
	// of the operation, or by someone inspecting the job.
	Args map[string]interface{} `json:"args" bson:"args"`

	// Tags []string `json:"tags" bson:"tags"`

	// State indicates whether the job is pending, running, completed or waiting
	// for a retry.
	State JobState `json:"state" bson:"state"`

	// Attempts indicates the number of times the job has been started.
	// It's zero while waiting for the initial execution, one for the jobs
	// executed once, and more than one for failed jobs that have been retried.
	Attempts int `json:"attempts" bson:"attempts"`

	FirstExecutionTime time.Time `json:"firsttime" bson:"firsttime"`
	LastExecutionTime  time.Time `json:"starttime" bson:"starttime"`
	LastCompletionTime time.Time `json:"endtime" bson:"endtime"`

	// EligibilityTime is when the job becomes eligible for the (next) execution.
	EligibilityTime time.Time `json:"eligtime" bson:"eligtime"`

	// ScheduledTime is the earliest time the job can be executed for the first time. It's basically the first EligibilityTime, and is used as a basis for FailedJobExpiration countdown.
	ScheduledTime time.Time `json:"schedtime" bson:"schedtime"`

	// FirstAddTime is the time the job has been added.
	FirstAddTime time.Time `json:"addtime" bson:"addtime"`

	// LastAddTime is the time the the last attempt to add the job (duplicate or not).
	LastAddTime time.Time `json:"lastaddtime" bson:"lastaddtime"`

	Error string `json:"error,omitempty",bson:"error,omitempty"`
}

type Request struct {
	// Name is a user-defined unique identifier of the job. An attempt to add a duplicate job with
	// the same name will be ignored; this makes enqueuing idempotent and safe to retry when
	// communication fails. bson.NewObjectId().Hex() is a good way to generate these.
	Name string `json:"name" bson:"name"`

	// Op determines the specific operation to perform. A single queue can support multiple
	// different kinds of operations.
	Op string

	// Args is a dictionary of arbitrary additional data. This can be used by the code
	// of the operation, or by someone inspecting the job.
	Args map[string]interface{}

	// Tags []string

	// CancelsTags []string
}

type RunParams struct {
	PollInterval time.Duration

	// Shutdown is a channel that, when closed, will cause the runner to shut down.
	Shutdown chan struct{}
}

type Stats struct {
	PendingJobs int

	JobCountsByState   map[JobState]int
	JobCountsByOpState map[string]map[JobState]int
}

// New creates a new job queue backed by the given MongoDB collection.
func New(collection *mgo.Collection, params Params) *Queue {
	if params.JobTimeout == 0 {
		params.JobTimeout = 60 * time.Second
	}
	if params.JobRetryMinDelay == 0 {
		params.JobRetryMinDelay = 1 * time.Second
	}
	if params.JobRetryMaxDelay == 0 {
		params.JobRetryMaxDelay = 1 * time.Hour
	}
	if params.JobRetryDelayMultiplier < 0.9 {
		params.JobRetryDelayMultiplier = 2
	}
	if params.FailedJobExpiration == 0 {
		params.FailedJobExpiration = 2 * 24 * time.Hour
	}
	if params.JobRemovalAge == 0 {
		params.JobRemovalAge = 14 * 24 * time.Hour
	}

	return &Queue{collection, params}
}

func (q *Queue) Migrate() error {
	err := q.Collection.EnsureIndex(mgo.Index{
		Key: []string{"state", "op", "eligtime"},
	})
	if err != nil {
		return errors.Wrap(err, "mongoqueue EnsureIndex failed")
	}

	_, err = q.Collection.UpdateAll(bson.M{
		"state":    Pending,
		"attempts": bson.M{"$gte": 1},
	}, bson.M{
		"$set": bson.M{
			"state": FailedPending,
		},
	})
	if err != nil {
		return errors.Wrap(err, "mongoqueue FailedPending migration failed")
	}

	return nil
}

// NewJobID returns a unique string suitable for use as a job name. This is NOT the best way
// to generate job names; the best way is to use a unique descriptive string that would prevent
// identical jobs from being enqueued.
//
// E.g., for a job that runs once per day, you can use "opname-YYYYmmdd" as the name.
func NewJobName() string {
	return bson.NewObjectId().Hex()
}

func (q *Queue) Add(request Request) error {
	now := bson.Now()

	_, err := q.Collection.Upsert(bson.M{
		"_id": request.Name,
	}, bson.M{
		"$setOnInsert": bson.M{
			"op":   request.Op,
			"args": request.Args,
			// "tags": request.Tags,

			"addtime":   now,
			"eligtime":  now,
			"schedtime": now,

			"state": Pending,
		},
		"$set": bson.M{
			"lastaddtime": now,
		},
		// "$inc": bson.M{},
		// "$max": bson.M{},
	})

	if err != nil {
		return fmt.Errorf("mongoqueue failed to add job (%v/%v): %v", request.Op, request.Name, err)
	}

	if err == nil {
		if q.Logger != nil {
			q.Logger.Printf("mongoqueue job added: %v/%v at %v (took %v)", request.Op, request.Name, now, bson.Now().Sub(now))
		}
	} else {
	}
	return err
}

type statsItem struct {
	op    string `bson:"_id"`
	count int    `bson:"count"`
}

func (q *Queue) Stats() (*Stats, error) {
	stats := Stats{
		JobCountsByOpState: make(map[string]map[JobState]int),
		JobCountsByState:   make(map[JobState]int),
	}

	var items []bson.M
	err := q.Collection.Pipe([]bson.M{
		bson.M{"$group": bson.M{"_id": bson.M{"op": "$op", "state": "$state"}, "count": bson.M{"$sum": 1}}},
	}).All(&items)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		key := item["_id"].(bson.M)
		op := key["op"].(string)
		state := JobState(key["state"].(string))
		count := item["count"].(int)

		if stats.JobCountsByOpState[op] == nil {
			stats.JobCountsByOpState[op] = make(map[JobState]int)
		}
		stats.JobCountsByOpState[op][state] = stats.JobCountsByOpState[op][state] + count
		stats.JobCountsByState[state] = stats.JobCountsByState[state] + count
		if state == Pending || state == FailedPending || state == Running {
			stats.PendingJobs += count
		}
	}
	// fmt.Printf("items = %#v\n", items)
	// fmt.Printf("stats = %#v\n", stats)

	return &stats, nil
}

func (q *Queue) Cleanup() error {
	now := bson.Now()
	threshold := now.Add(-1 * q.JobRemovalAge)

	info, err := q.Collection.RemoveAll(bson.M{
		"state":    bson.M{"$in": []JobState{Succeeded, Failed, Canceled}},
		"eligtime": bson.M{"$lte": threshold},
	})
	if err != nil {
		return err
	}
	if info.Removed > 0 {
		if q.Logger != nil {
			q.Logger.Printf("mongoqueue removed %v old jobs", info.Removed)
		}
	}

	return nil
}

func (q *Queue) Recover() error {
	now := bson.Now()
	threshold := now.Add(-1 * q.JobTimeout)
	iter := q.Collection.Find(bson.M{
		"state":     Running,
		"starttime": bson.M{"$lte": threshold},
	}).Iter()

	var job *Job
	count := 0
	for iter.Next(&job) {
		if q.Logger != nil {
			q.Logger.Printf("NOTICE mongoqueue recovering job: %v/%v age=%v state=%v", job.Op, job.Name, now.Sub(job.LastExecutionTime), job.State)
		}
		err := q.failed(job, ErrTimeout.Error(), false)
		if err != nil {
			iter.Close()
			return fmt.Errorf("mongoqueue job recovery failed (%v/%v): %v", job.Op, job.Name, err)
		}
		count++
	}

	err := iter.Close()
	if err != nil {
		err = fmt.Errorf("mongoqueue recovery query failed: %v", err)
		return err
	}

	return nil
}

func (q *Queue) Run(workers map[string]Worker, params RunParams) {
	found := true

Loop:
	for {
		if found {
			select {
			case <-params.Shutdown:
				break Loop
			default:
			}
		} else {
			err := q.Recover()
			if err != nil {
				if q.Logger != nil {
					q.Logger.Printf("ERROR: %v", err)
				}
			}
			err = q.Cleanup()
			if err != nil {
				if q.Logger != nil {
					q.Logger.Printf("ERROR: %v", err)
				}
			}
			select {
			case <-params.Shutdown:
				break Loop
			case <-time.After(params.PollInterval):
			}
		}
		var err error
		found, err = q.RunOnce(workers)
		if err != nil {
		} else {
		}
	}
}

func (q *Queue) RunOnce(workers map[string]Worker) (bool, error) {
	ops := make([]string, 0, len(workers))
	for k := range workers {
		ops = append(ops, k)
	}

	job, err := q.Dequeue(ops)
	if err != nil {
		return false, err
	}
	if job == nil {
		return false, nil
	}

	worker := workers[job.Op]
	workErr := worker.Func(job.Op, job.Args)

	if workErr == nil {
		err = q.Succeeded(job.Name)
	} else {
		permanent := false
		if perr, ok := workErr.(PermanentError); ok {
			permanent = perr.IsPermanent()
		}
		err = q.Failed(job.Name, workErr.Error(), permanent)
	}

	return true, nil
}

func (q *Queue) Dequeue(ops []string) (*Job, error) {
	now := bson.Now()

	var job *Job
	_, err := q.Collection.Find(bson.M{
		"op":       bson.M{"$in": ops},
		"state":    bson.M{"$in": []JobState{Pending, FailedPending}},
		"eligtime": bson.M{"$lte": now},
	}).Apply(mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"state":     Running,
				"starttime": now,
			},
			"$inc": bson.M{
				"attempts": 1,
			},
		},
		ReturnNew: true,
	}, &job)

	if err != nil {
		if err == mgo.ErrNotFound {
			// err = q.Collection.Find(bson.M{
			// 	"op":    bson.M{"$in": ops},
			// 	"state": Pending,
			// }).One(&job)
			// if err == nil {
			// 	if q.Logger != nil {
			// 		q.Logger.Printf("=> %v (op %v) not eligible until %v, now = %v", job.Name, job.Op, job.EligibilityTime, now)
			// 	}
			// }

			// if q.Logger != nil {
			// 	q.Logger.Printf("mongoqueue polling returning no jobs for ops: %v", ops)
			// }
			return nil, nil
		} else {
			return nil, err
		}
	}

	err = q.Collection.Update(bson.M{
		"_id":       job.Name,
		"firsttime": bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"firsttime": now,
		},
	})
	if err != nil && err != mgo.ErrNotFound {
		err = fmt.Errorf("mongoqueue failed update of firsttime when dequeuing (%v/%v): %v", job.Op, job.Name, err)
		if q.Logger != nil {
			q.Logger.Printf("ERROR %v", err)
		}
		job.FirstExecutionTime = now
		q.failed(job, err.Error(), false)
		return nil, err
	}

	// println("test")

	if q.Logger != nil {
		q.Logger.Printf("mongoqueue job dequeued: %v/%v %v", job.Op, job.Name, job.Args)
	}
	return job, nil
}

func (q *Queue) Succeeded(name string) error {
	var job *Job
	err := q.Collection.FindId(name).One(&job)
	if err != nil {
		return err
	}

	if job.State.IsFinal() {
		return ErrCannotCompleteWhenNotRunning
	}

	if q.Logger != nil {
		q.Logger.Printf("mongoqueue job succeeded: %v/%v", job.Op, job.Name)
	}

	err = q.Collection.UpdateId(name, bson.M{
		"$set": bson.M{
			"state": Succeeded,
		},
	})
	return err
}

func (q *Queue) Failed(name string, errmsg string, permanent bool) error {
	var job *Job
	err := q.Collection.FindId(name).One(&job)
	if err != nil {
		return err
	}

	if job.State != Running {
		return ErrCannotCompleteWhenNotRunning
	}

	if q.Logger != nil {
		q.Logger.Printf("mongoqueue job failed (%v/%v): %v", job.Op, job.Name, errmsg)
	}

	return q.failed(job, errmsg, permanent)
}

func (q *Queue) failed(job *Job, errmsg string, permanent bool) error {
	now := bson.Now()

	// if q.Logger != nil {
	// 	q.Logger.Printf("mongoqueue failed: job.ScheduledTime = %v, q.FailedJobExpiration = %v, Add() = %v, now = %v, Before = %v", job.ScheduledTime, q.FailedJobExpiration, job.ScheduledTime.Add(q.FailedJobExpiration), now, job.ScheduledTime.Add(q.FailedJobExpiration).Before(now))
	// }

	if permanent {
		if q.Logger != nil {
			q.Logger.Printf("WARNING mongoqueue permanent failure of job %v/%v: %v", job.Op, job.Name, errmsg)
		}
		return q.Collection.UpdateId(job.Name, bson.M{
			"$set": bson.M{
				"state":   Failed,
				"error":   errmsg,
				"endtime": now,
			},
		})
	}

	if job.ScheduledTime.Add(q.FailedJobExpiration).Before(now) {
		if q.Logger != nil {
			q.Logger.Printf("WARNING mongoqueue failing job expired %v/%v: %v", job.Op, job.Name, errmsg)
		}
		return q.Collection.UpdateId(job.Name, bson.M{
			"$set": bson.M{
				"state":   Failed,
				"error":   errmsg,
				"endtime": now,
			},
		})
	}

	delay := time.Duration(float64(q.JobRetryMinDelay) * math.Pow(q.JobRetryDelayMultiplier, float64(job.Attempts)))
	if delay > q.JobRetryMaxDelay {
		delay = q.JobRetryMaxDelay
	}

	if q.Logger != nil {
		q.Logger.Printf("NOTICE mongoqueue will retry job %v/%v after %v: %v", job.Op, job.Name, delay, errmsg)
	}

	return q.Collection.UpdateId(job.Name, bson.M{
		"$set": bson.M{
			"state":    FailedPending,
			"error":    errmsg,
			"endtime":  now,
			"eligtime": now.Add(delay),
		},
	})
}

func (q *Queue) RerunFailedJobsNow() error {
	now := bson.Now()

	_, err := q.Collection.UpdateAll(bson.M{
		"state": FailedPending,
	}, bson.M{
		"$min": bson.M{
			"eligtime": now,
		},
	})
	if err != nil {
		return errors.Wrap(err, "mongoqueue update (for pending jobs) failed")
	}

	_, err = q.Collection.UpdateAll(bson.M{
		"state": Failed,
	}, bson.M{
		"$set": bson.M{
			"state":     FailedPending,
			"schedtime": now,
			"attempts":  0,
			"eligtime":  now,
		},
	})
	if err != nil {
		return errors.Wrap(err, "mongoqueue update (for failed jobs) failed")
	}

	return nil
}

func (q *Queue) RerunJobNow(name string) (bool, error) {
	now := bson.Now()

	err := q.Collection.UpdateId(bson.M{
		"_id":   name,
		"state": []JobState{Pending, FailedPending},
	}, bson.M{
		"$min": bson.M{
			"eligtime": now,
		},
	})
	if err == nil {
		return true, nil
	}
	if err != mgo.ErrNotFound {
		return false, errors.Wrap(err, "mongoqueue update (for pending jobs) failed")
	}

	err = q.Collection.UpdateId(bson.M{
		"_id":   name,
		"state": Failed,
	}, bson.M{
		"$set": bson.M{
			"state":     FailedPending,
			"schedtime": now,
			"attempts":  0,
			"eligtime":  now,
		},
	})
	if err == nil {
		return true, nil
	}
	if err != mgo.ErrNotFound {
		return false, errors.Wrap(err, "mongoqueue update (for failed jobs) failed")
	}

	return false, nil
}

type WaitParams struct {
	Timeout  time.Duration
	Deadline time.Time

	PollInterval time.Duration
}

func (params *WaitParams) StartDeadline() {
	if params.Deadline.IsZero() {
		params.Deadline = time.Now().Add(params.Timeout)
	}
}

func (params WaitParams) IsAfterDeadline() bool {
	return params.Deadline.IsZero() || time.Now().After(params.Deadline)
}

func (q *Queue) Wait(names []string, params WaitParams) (bool, error) {
	if params.PollInterval == 0 {
		params.PollInterval = 1 * time.Second
	}
	if params.Timeout < params.PollInterval {
		params.PollInterval = params.Timeout
	}

	params.StartDeadline()
	for {
		n, err := q.Collection.Find(bson.M{"_id": bson.M{"$in": names}, "state": bson.M{"$in": []JobState{Pending, FailedPending, Running}}}).Count()
		if err != nil {
			return false, fmt.Errorf("mongoqueue waiting poll failed: %v", err)
		}

		if n == 0 {
			return true, nil
		}

		if params.IsAfterDeadline() {
			return false, nil
		}

		time.Sleep(params.PollInterval)
	}
}

func (q *Queue) GetStatus(jobId bson.ObjectId) (JobState, error) {
	var job *Job
	err := q.Collection.FindId(jobId).One(&job)

	if err != nil {
		return Unknown, err
	}
	return job.State, nil
}

func (q *Queue) GetPendingCount() (int, error) {
	pending, err := q.Collection.Find(bson.M{"state": bson.M{"$in": []JobState{Pending, FailedPending, Running}}}).Count()
	if err != nil {
		return -1, err
	}
	return pending, nil
}

func (q *Queue) GetPendingCountForOp(op string) (int, error) {
	pending, err := q.Collection.Find(bson.M{"state": bson.M{"$in": []JobState{Pending, FailedPending, Running}}}).Count()
	if err != nil {
		return -1, err
	}
	return pending, nil
}

func (q *Queue) GetPendingJobs() ([]*Job, error) {
	var jobs []*Job
	err := q.Collection.Find(bson.M{"state": bson.M{"$in": []JobState{Pending, FailedPending, Running}}}).All(&jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (q *Queue) WaitUntilEmpty(pollInterval time.Duration, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return ErrTimeout
		}
		cnt, err := q.GetPendingCount()
		if err != nil {
			return err
		}
		if cnt == 0 {
			return nil
		}
		time.Sleep(pollInterval)
	}
}

// func RecoverSumm() error {
// 	now := bson.Now()
// 	threshold := now.Add(-60 * time.Second)

// 	info, err := collection.UpdateAll(bson.M{
// 		"summ.running": true,
// 		"summ.runtime": bson.M{"$lt": threshold},
// 	}, bson.M{
// 		"$set": bson.M{
// 			"summ.running": false,
// 			"summ.pending": true,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// fmt.Fprintf(os.Stderr, "Recovered %d stuck resources.\n", info.Updated)
// 	if info.Updated > 0 {
// 		fmt.Fprintf(os.Stderr, "Recovered %d stuck resources.\n", info.Updated)
// 	}

// 	return nil
// }
