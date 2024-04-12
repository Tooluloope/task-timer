package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/icza/gox/timex"
	"github.com/teris-io/shortid"
	"github.com/tooluloope/task-timer/pkg/config"
)

type TaskStatus int
type CSVColumn int

const (
	ID CSVColumn = iota
	Name
	TimeIntervals
	TotalTime
	Tags
	Status
	CreatedAt
	UpdatedAt
)

const (
	Created TaskStatus = iota
	Running
	Stopped
	Paused
	Completed
)

type TimeInterval struct {
	StartTime time.Time
	EndTime   time.Time
}

type Task struct {
	ID            string
	Name          string
	TimeIntervals []TimeInterval
	Tags          []string
	Status        string
	TotalTime     time.Duration
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Storage interface {
	SaveTask(task Task) (string, error)
	GetTasksByName(name string) ([]Task, error)
	GetTaskByID(id string) (Task, error)
	GetAllTasks() ([]Task, error)
	DeleteTask(id string) error
	StartTask(task Task) error
	StopTask(task Task) error
}

type CSVStorage struct {
	filepath string
	mu       sync.Mutex
	sid      *shortid.Shortid
}

var Data *CSVStorage

func init() {
	filepath := config.EnvConfigs.DataPath
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		log.Fatal(err)
	}

	Data = NewCSVStorage(filepath, sid)
}

func (ts TaskStatus) String() string {
	return [...]string{"Created", "Running", "Stopped", "Paused", "Completed"}[ts]
}

func NewCSVStorage(filepath string, sid *shortid.Shortid) *CSVStorage {
	return &CSVStorage{
		filepath: filepath,
		sid:      sid,
	}
}

func (t *Task) getTotalTime() (totalTime time.Duration, err error) {

	totalTime = t.TotalTime

	for _, interval := range t.TimeIntervals {
		if interval.EndTime.IsZero() {
			totalTime += time.Since(interval.StartTime)
		} else {
			totalTime += interval.EndTime.Sub(interval.StartTime)
		}
	}
	return totalTime, nil
}

func (t *Task) serializeIntervals() (timeIntervals []string, err error) {
	for _, interval := range t.TimeIntervals {
		timeIntervals = append(timeIntervals, fmt.Sprintf("%s,%s",
			interval.StartTime.Format(time.RFC3339),
			interval.EndTime.Format(time.RFC3339)))
	}

	return
}

func deserializeIntervals(intervals string) (timeIntervals []TimeInterval, err error) {

	if intervals == "" {
		return
	}

	intervalStrings := strings.Split(intervals, ";")
	var (
		startTime time.Time
		endTime   time.Time
	)

	for _, intervalString := range intervalStrings {
		times := strings.Split(intervalString, ",")
		startTime, err = time.Parse(time.RFC3339, times[0])
		if err != nil {
			return
		}
		endTime, err = time.Parse(time.RFC3339, times[1])
		if err != nil {
			return
		}
		timeIntervals = append(timeIntervals, TimeInterval{StartTime: startTime, EndTime: endTime})
	}

	return
}

func (c *CSVStorage) SaveTask(task Task) (id string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fileInfo, err := os.Stat(c.filepath)
	writeHeader := os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0)

	file, err := os.OpenFile(c.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if writeHeader {
		header := []string{"ID", "Name", "Time Intervals", "Total Time", "Tags", "Status", "Created At", "Updated At"}
		if err = writer.Write(header); err != nil {
			return
		}
	}
	id, err = c.sid.Generate()
	if err != nil {
		return
	}
	task.ID = id

	record, err := taskToRecord(task)
	if err != nil {
		return
	}

	if err = writer.Write(record); err != nil {
		return
	}

	if err = writer.Error(); err != nil {
		return
	}

	writer.Flush()
	return id, err
}

func (c *CSVStorage) GetTaskByID(id string) (task Task, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	records, err := c.readAllRecords()
	if err != nil {
		return Task{}, err
	}

	for _, record := range records {

		if record[ID] == id {
			return RecordToTask(record)
		}
	}

	return
}

func (c *CSVStorage) GetTasksByName(name string) (tasks []Task, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	records, err := c.readAllRecords()
	if err != nil {
		return []Task{}, err
	}

	for _, record := range records {
		if record[Name] == name {
			task, err := RecordToTask(record)
			if err != nil {
				return []Task{}, err
			}
			tasks = append(tasks, task)
		}
	}
	return
}

func (c *CSVStorage) GetAllTasks() ([]Task, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return []Task{}, nil
}

func (c *CSVStorage) DeleteTask(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return nil
}

func (c *CSVStorage) StartTask(task Task) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if task.Status == Running.String() {
		return fmt.Errorf("Task already started")
	}

	task.Status = Running.String()
	totalTime, err := task.getTotalTime()
	if err != nil {
		return err
	}
	task.TotalTime = totalTime
	task.TimeIntervals = []TimeInterval{{
		StartTime: time.Now(),
	}}
	task.UpdatedAt = time.Now()
	return updateTask(task, c)
}

func (c *CSVStorage) StopTask(task Task) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if task.Status == Stopped.String() {
		return fmt.Errorf("Task already ended")
	}

	task.Status = Stopped.String()
	lastRecordedTime := task.TimeIntervals[len(task.TimeIntervals)-1]
	task.TimeIntervals[len(task.TimeIntervals)-1] = TimeInterval{
		StartTime: lastRecordedTime.StartTime,
		EndTime:   time.Now(),
	}
	task.UpdatedAt = time.Now()
	return updateTask(task, c)
}

func updateTask(task Task, ctx *CSVStorage) error {

	records, err := ctx.readAllRecords()
	if err != nil {
		return err
	}

	updated := false
	for i, record := range records {

		if record[ID] == task.ID {

			updatedRecord, err := taskToRecord(task)
			if err != nil {
				return err
			}

			records[i] = updatedRecord
			updated = true

			break
		}
	}

	if !updated {
		return fmt.Errorf("no task found with id %s", task.ID)
	}
	return ctx.writeAllRecords(records)
}

func (c *CSVStorage) readAllRecords() ([][]string, error) {
	file, err := os.Open(c.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func (c *CSVStorage) writeAllRecords(records [][]string) error {
	file, err := os.Create(c.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		return err
	}
	return nil
}

func taskToRecord(task Task) ([]string, error) {
	timeIntervals, err := task.serializeIntervals()
	if err != nil {
		return nil, err
	}

	totalTime, err := task.getTotalTime()
	if err != nil {
		return nil, err
	}

	record := []string{task.ID, task.Name, strings.Join(timeIntervals, ";"), timex.ShortDuration(totalTime), strings.Join(task.Tags, ";"), task.Status, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339)}
	return record, nil
}

func RecordToTask(record []string) (task Task, err error) {

	var (
		createdAt time.Time
		updatedAt time.Time
		intervals []TimeInterval
	)

	intervals, err = deserializeIntervals(record[TimeIntervals])
	if err != nil {
		return Task{}, err
	}
	createdAt, err = time.Parse(time.RFC3339, record[CreatedAt])
	if err != nil {
		return Task{}, err
	}
	updatedAt, err = time.Parse(time.RFC3339, record[UpdatedAt])
	if err != nil {
		return Task{}, err
	}
	totalTime, err := time.ParseDuration(record[TotalTime])
	if err != nil {
		return Task{}, err
	}

	task = Task{
		ID:            record[ID],
		Name:          record[Name],
		TimeIntervals: intervals,
		Tags:          strings.Split(record[Tags], ";"),
		Status:        record[Status],
		TotalTime:     totalTime,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return
}
