package tests

import (
	"fmt"
	dbtest "github.com/yourusername/basic-a/business/data/tests"
	"github.com/yourusername/basic-a/foundation/docker"
	"testing"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}
