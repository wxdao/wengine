package wengine

import (
	"errors"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"time"
)

const (
	WINDOW_MODE_WINDOWED = iota
	WINDOW_MODE_FULL_SCREEN
)

type Config struct {
	Context       *Context
	Width, Height int
	WindowMode    int
	WindowTitle   string
	FrameLimit    int
}

type App struct {
	width, height int
	winMode       int
	title         string
	frameLimit    int

	lastScene *Scene

	currentTime float64
	lastTime    float64

	context *Context
}

func NewApp(config *Config) (*App, error) {
	if config.Context == nil {
		config.Context = NewContext()
	}
	return &App{
		width:      config.Width,
		height:     config.Height,
		title:      config.WindowTitle,
		winMode:    config.WindowMode,
		frameLimit: config.FrameLimit,
		context:    config.Context,
	}, nil
}

func (a *App) Context() *Context {
	return a.context
}

func (a *App) Run() error {
	switch a.context.rendererName {
	case "opengl":
		return a.runOpenGL()
	default:
		return errors.New("no available renderer runner")
	}
}

func (a *App) runOpenGL() error {
	var err error

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err = glfw.Init(); err != nil {
		return errors.New(fmt.Sprint("unable to init glfw:", err))
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var window *glfw.Window
	switch a.winMode {
	case WINDOW_MODE_WINDOWED:
		window, err = glfw.CreateWindow(a.width, a.height, a.title, nil, nil)
		break
	case WINDOW_MODE_FULL_SCREEN:
		window, err = glfw.CreateWindow(a.width, a.height, a.title, glfw.GetPrimaryMonitor(), nil)
		break
	default:
		return errors.New("window uninitialized")
	}
	if err != nil {
		return errors.New(fmt.Sprint("unable to init window:", err))
	}
	window.MakeContextCurrent()

	scrWidth, scrHeight := window.GetFramebufferSize()
	a.context.SetScreenSize(scrWidth, scrHeight)

	if err := a.context.renderer.Init(a.context); err != nil {
		return err
	}

	fmt.Println("OpenGL version", a.context.renderer.Version())

	glfw.SwapInterval(1)

	fps := 0
	fpsDisplayLastTime := a.currentTime

	a.currentTime = glfw.GetTime()

	for !window.ShouldClose() {
		a.lastTime = a.currentTime
		a.currentTime = glfw.GetTime()

		currentScene := a.context.currentScene
		if currentScene != a.lastScene {
			result, err := a.context.asyncLoadScene(currentScene)
			if err != nil {
				return err
			}
			err = <-result
			if err != nil {
				return err
			}
			a.executeBehaviors(true)
		}
		a.lastScene = currentScene

		currentScene.updateTransforms()
		err := a.context.renderer.Render(currentScene)
		if err != nil {
			return err
		}
		window.SwapBuffers()
		glfw.PollEvents()
		a.executeBehaviors(false)

		if a.frameLimit != 0 {
			sleepTime := 1.0/float64(a.frameLimit) - (glfw.GetTime() - a.currentTime)
			if sleepTime > 0 {
				time.Sleep(time.Duration(float64(time.Second) * sleepTime))
			}
		}

		fps += 1
		if glfw.GetTime()-fpsDisplayLastTime >= 1 {
			fpsDisplayLastTime = glfw.GetTime()
			fmt.Println("fps:", fps)
			fps = 0
		}
	}
	return nil
}

func (a *App) executeBehaviors(runStart bool) {
	for _, obj := range a.context.currentScene.objects {
		if obj.enabled && obj.behavior != nil {
			bctx := BehaviorContext{Context: a.context, Object: obj, Time: a.currentTime, DeltaTime: a.currentTime - a.lastTime}
			if runStart {
				obj.behavior.Start(&bctx)
			} else {
				obj.behavior.Update(&bctx)
			}
		}
	}
}
