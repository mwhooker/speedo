// Open the seed studio grove kit box and find the little green bag labeled "Grove - LCD", open the
// bag and plug the grove connector cable into the grove slot labeled "I2C".
package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	board := edison.NewEdisonAdaptor("edison")
	lidar := NewLIDARLiteDriver(board, "lidar")
	screen := i2c.NewGroveLcdDriver(board, "screen")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			distance, err := lidar.Distance()
			if err != nil {
				fmt.Println("error: ", err)
			}
			dist := float32(distance) / 2.54
			fmt.Println("Distance (in)", dist)
			screen.Clear()
			screen.Home()
			screen.Write(fmt.Sprintf("%f\"", dist))
		})
	}

	robot := gobot.NewRobot("screenBot",
		[]gobot.Connection{board},
		[]gobot.Device{lidar, screen},
		work,
	)

	gbot.AddRobot(robot)

	/*
		if errs := screen.Start(); len(errs) > 0 {
			for _, err := range errs {
				fmt.Println(err)
			}
		}
	*/

	gbot.Start()
}
