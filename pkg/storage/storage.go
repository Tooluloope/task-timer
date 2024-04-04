package storage

import (
	"encoding/csv"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/icza/gox/timex"
)

type TimeInterval struct {
	StartTime string
	EndTime string
}

type Task struct {
	Name string
	TimeIntervals []TimeInterval
	Tags []string
	CreatedAt string
	UpdatedAt string
}

type Storage interface {
	SaveTask(task Task) error
	GetTask(name string) (Task, error)
	GetAllTasks() ([]Task, error)
	DeleteTask(name string) error
	UpdateTask(task Task) error
}

type CSVStorage struct {
	filepath string
	mu sync.Mutex
}

var Data = NewCSVStorage("/Users/tolulopeadetula/Documents/GitHub/task-timer/data/data.csv")


func NewCSVStorage(filepath string) *CSVStorage {
	return &CSVStorage{
		filepath: filepath,
	}
}

func (t Task) TotalTime() (totalTime time.Duration, err error) {

	totalTime = 0
    layout := "2006-01-02T15:04:05"

    for _, interval := range t.TimeIntervals {
        start, err := time.Parse(layout, interval.StartTime)
        if err != nil {
            return 0, err
        }
        end, err := time.Parse(layout, interval.EndTime)
        if err != nil {
            return 0, err
        }
        totalTime += end.Sub(start)
    }
    return totalTime, nil
}

func (c *CSVStorage) SaveTask(task Task) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.OpenFile(c.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()



	totalTime, err := task.TotalTime()
	if err != nil {
		return err
	}

	record := []string{task.Name, timex.ShortDuration(totalTime), strings.Join(task.Tags, ","), task.CreatedAt, task.UpdatedAt}

	return writer.Write(record)
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
