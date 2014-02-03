package main

// TODO:
// 	1. Refactor "engine" out (everything should be under "game")
//  2. Define Entity and Object interfaces
//		- Entity is something update-able
//		- Object is something visible on screen
//  3. Add game.AddObject() method?
//  4. Rework widgets to be auto-updated (add game.AddWidget()?)
//  5. Implement double-pass updates
//  6. Come up with a consistent scheme for rounding

import (
	"flag"
	"game/engine"
	scenes "game/scenes/menu"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/image"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const FPS int = 60

var (
	cpuprofile = flag.String("cpuprofile", "", "output a cpuprofile to...")
	memprofile = flag.String("memprofile", "", "output a memprofile to...")
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	flag.Parse()

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
	secondsPerFrame := 1 / float64(FPS)
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
		event       allegro.Event
		running     bool          = true
		needsUpdate bool          = false
		lastUpdate  time.Time     = time.Now()
		step        time.Duration = time.Duration(secondsPerFrame * float64(time.Second))
		lag         time.Duration
		now         time.Time
		elapsed     time.Duration
	)

	fpsTimer.Start()

	for {
		eventQueue.WaitForEvent(&event)
		if !engine.HandleEvent(&event) {
			continue
		}

		switch event.Type {
		case allegro.EVENT_TIMER:
			if event.Source == fpsTimer.EventSource() {
				needsUpdate = true
				fpsTimer.SetCount(0)
			}

		case allegro.EVENT_DISPLAY_CLOSE:
			running = false
		}

		if !running {
			break
		}

		if needsUpdate && eventQueue.IsEmpty() {
			now = time.Now()
			elapsed = now.Sub(lastUpdate)
			lastUpdate = now
			lag += elapsed
			for lag >= step {
				engine.Update()
				lag -= step
			}
			engine.Render(float32(lag / step))
			needsUpdate = false
		}
	}
}
