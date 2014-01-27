package main

import (
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/image"
	"game/engine"
	scenes "game/scenes/menu"
	"time"
)

const (
	FPS int = 60
)

func main() {
	var (
		display    *allegro.Display
		eventQueue *allegro.EventQueue
		fpsTimer   *allegro.Timer
		err        error
	)

	// Event Queue
	if eventQueue, err = allegro.CreateEventQueue(); err == nil {
		defer eventQueue.Destroy()
	} else {
		panic(err)
	}

	// Display
	allegro.SetNewDisplayFlags(allegro.WINDOWED)
	if display, err = allegro.CreateDisplay(640, 480); err == nil {
		defer display.Destroy()
		display.SetWindowTitle("My Game")
		eventQueue.RegisterEventSource(display.EventSource())
	} else {
		panic(err)
	}

	// Mouse
	if err = allegro.InstallMouse(); err == nil {
		var mouseEventSource *allegro.EventSource
		if mouseEventSource, err = allegro.MouseEventSource(); err == nil {
			eventQueue.RegisterEventSource(mouseEventSource)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

	// FPS Timer
	secondsPerFrame := float64(1) / float64(FPS)
	if fpsTimer, err = allegro.CreateTimer(secondsPerFrame); err == nil {
		defer fpsTimer.Destroy()
		eventQueue.RegisterEventSource(fpsTimer.EventSource())
	} else {
		panic(err)
	}

	// Submodules
	font.Init()
	image.Init()

	engine.Init()
	engine.GoTo(eventQueue, &scenes.MenuScene{Name: "main"})

	fpsTimer.Start()

	var (
		running    bool          = true
		lastUpdate time.Time     = time.Now()
		step       time.Duration = time.Duration(secondsPerFrame * float64(time.Second))
		lag        time.Duration
	)
	for running {
		event := eventQueue.WaitForEvent(false)
		if !engine.HandleEvent(event) {
			continue
		}

		switch event.Type {
		case allegro.EVENT_TIMER:
			// We don't care what the count is at, only that it ticked.
			// Reset it to 0 to prevent a possible eventual overflow.
			fpsTimer.SetCount(0)
			now := time.Now()
			elapsed := now.Sub(lastUpdate)
			lag += elapsed
			for lag >= step {
				engine.Update()
				lag -= step
			}
			engine.Render(float64(lag) / float64(step))

		case allegro.EVENT_DISPLAY_CLOSE:
			running = false
		}
	}
}
