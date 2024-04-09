package storage

import (
	"encoding/csv"
	"fmt"
	"io"
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
	GetTaskByName(name string) (Task, error)
	GetTaskByID(id string) (Task, error)
	GetAllTasks() ([]Task, error)
	DeleteTask(id string) error
	UpdateTask(id string) error
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

	for _, interval := range t.TimeIntervals {
		start, err := time.Parse(time.RFC3339, interval.StartTime.Format(time.RFC3339))
		if err != nil {
			return 0, err
		}
		end, err := time.Parse(time.RFC3339, interval.EndTime.Format(time.RFC3339))
		if err != nil {
			return 0, err
		}
		totalTime += end.Sub(start)
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
		header := []string{"ID", "Name", "Time Intervals", "Total Time", "Tags", "Created At", "Updated At"}
		if err = writer.Write(header); err != nil {
			return
		}
	}

	timeIntervals, err := task.serializeIntervals()
	if err != nil {
		return
	}

	totalTime, err := task.TotalTime()
	if err != nil {
		return
	}

	id, err = c.sid.Generate()
	if err != nil {
		return
	}

	record := []string{id, task.Name, strings.Join(timeIntervals, ";"), timex.ShortDuration(totalTime), strings.Join(task.Tags, ";"), task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339)}

	if err = writer.Write(record); err != nil {
		return
	}

	writer.Flush()
	return id, err
}

func (c *CSVStorage) GetTaskByID(id string) (task Task, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.filepath)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return Task{}, err
		}

		if record[0] == id {

			var (
				createdAt time.Time
				updatedAt time.Time
				intervals []TimeInterval
			)

			intervals, err = deserializeIntervals(record[2])
			if err != nil {
				return Task{}, err
			}
			createdAt, err = time.Parse(time.RFC3339, record[5])
			if err != nil {
				return Task{}, err
			}
			updatedAt, err = time.Parse(time.RFC3339, record[6])
			if err != nil {
				return Task{}, err
			}
			return Task{

				ID:            record[0],
				Name:          record[1],
				TimeIntervals: intervals,
				Tags:          strings.Split(record[4], ";"),
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			}, nil
		}
	}

	return Task{}, nil
}

func (c *CSVStorage) GetTaskByName(name string) (Task, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Task{}, nil
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

func (c *CSVStorage) UpdateTask(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return nil
}
