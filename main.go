package main

import (
	"fmt"

	"github.com/tooluloope/task-timer/cmd"
	"github.com/tooluloope/task-timer/pkg/config"
)

func init() {}

func main() {
	fmt.Print(config.EnvConfigs.DataPath)

	cmd.Execute()
}
