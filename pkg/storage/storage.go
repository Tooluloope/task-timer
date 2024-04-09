package storage

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/icza/gox/timex"
	"github.com/teris-io/shortid"
	"github.com/tooluloope/task-timer/pkg/config"
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Storage interface {
	SaveTask(task Task) (string, error)
	GetTask(name string) (Task, error)
	GetAllTasks() ([]Task, error)
	DeleteTask(name string) error
	UpdateTask(task Task) error
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

func NewCSVStorage(filepath string, sid *shortid.Shortid) *CSVStorage {
	return &CSVStorage{
		filepath: filepath,
		sid:      sid,
	}
}

func (t *Task) TotalTime() (totalTime time.Duration, err error) {

	totalTime = 0
	layout := "2006-01-02T15:04:05"

	for _, interval := range t.TimeIntervals {
		start, err := time.Parse(layout, interval.StartTime.Format(layout))
		if err != nil {
			return 0, err
		}
		end, err := time.Parse(layout, interval.EndTime.Format(layout))
		if err != nil {
			return 0, err
		}
		totalTime += end.Sub(start)
	}
	return totalTime, nil
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
		header := []string{"ID", "Name", "Total Time", "Tags", "Created At", "Updated At"}
		if err = writer.Write(header); err != nil {
			return
		}
	}

	totalTime, err := task.TotalTime()
	if err != nil {
		return
	}
	id, err = c.sid.Generate()
	if err != nil {
		return
	}

	record := []string{id, task.Name, timex.ShortDuration(totalTime), strings.Join(task.Tags, ","), task.CreatedAt.String(), task.UpdatedAt.String()}

	if err = writer.Write(record); err != nil {
		return
	}

	writer.Flush()
	return id, err
}

func (c *CSVStorage) GetTask(name string) (Task, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Task{}, nil
}

func (c *CSVStorage) GetAllTasks() ([]Task, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return []Task{}, nil
}

func (c *CSVStorage) DeleteTask(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return nil
}

func (c *CSVStorage) UpdateTask(task Task) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return nil
}
