package main

import (
	"flag"
	"game/engine"
	scenes "game/scenes/menu"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/image"
	"os"
	"runtime/pprof"
	"time"
)

const FPS int = 60

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	var (
		display    *allegro.Display
		eventQueue *allegro.EventQueue
		fpsTimer   *allegro.Timer
		err        error
	)

	flag.Parse()

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

	// Keyboard
	if err = allegro.InstallKeyboard(); err == nil {
		var keyboardEventSource *allegro.EventSource
		if keyboardEventSource, err = allegro.KeyboardEventSource(); err == nil {
			eventQueue.RegisterEventSource(keyboardEventSource)
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

	engine.Init(display)
	engine.GoTo(eventQueue, &scenes.MenuScene{Name: "main"})

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var (
		running    bool          = true
		lastUpdate time.Time     = time.Now()
		step       time.Duration = time.Duration(secondsPerFrame * float64(time.Second))
		lag        time.Duration = 0
	)

	fpsTimer.Start()

	for running {
		event := eventQueue.WaitForEvent(false)
		if !engine.HandleEvent(event) {
			continue
		}

		switch event.Type {
		case allegro.EVENT_TIMER:
			fpsTimer.SetCount(0)
			now := time.Now()
			elapsed := now.Sub(lastUpdate)
			lastUpdate = now
			lag += elapsed
			for lag >= step {
				engine.Update()
				lag -= step
			}
			engine.Render(float32(lag / step))

		case allegro.EVENT_DISPLAY_CLOSE:
			running = false
		}
	}
}
