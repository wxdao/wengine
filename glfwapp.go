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

var keyMap = map[glfw.Key]string{
	glfw.KeyUnknown:      "unknown",
	glfw.KeySpace:        "space",
	glfw.KeyApostrophe:   "",
	glfw.KeyComma:        ",",
	glfw.KeyMinus:        "-",
	glfw.KeyPeriod:       ".",
	glfw.KeySlash:        "-",
	glfw.Key0:            "0",
	glfw.Key1:            "1",
	glfw.Key2:            "2",
	glfw.Key3:            "3",
	glfw.Key4:            "4",
	glfw.Key5:            "5",
	glfw.Key6:            "6",
	glfw.Key7:            "7",
	glfw.Key8:            "8",
	glfw.Key9:            "9",
	glfw.KeySemicolon:    ";",
	glfw.KeyEqual:        "=",
	glfw.KeyA:            "a",
	glfw.KeyB:            "b",
	glfw.KeyC:            "c",
	glfw.KeyD:            "d",
	glfw.KeyE:            "e",
	glfw.KeyF:            "f",
	glfw.KeyG:            "g",
	glfw.KeyH:            "h",
	glfw.KeyI:            "i",
	glfw.KeyJ:            "j",
	glfw.KeyK:            "k",
	glfw.KeyL:            "l",
	glfw.KeyM:            "m",
	glfw.KeyN:            "n",
	glfw.KeyO:            "o",
	glfw.KeyP:            "p",
	glfw.KeyQ:            "q",
	glfw.KeyR:            "r",
	glfw.KeyS:            "s",
	glfw.KeyT:            "t",
	glfw.KeyU:            "u",
	glfw.KeyV:            "v",
	glfw.KeyW:            "w",
	glfw.KeyX:            "x",
	glfw.KeyY:            "y",
	glfw.KeyZ:            "z",
	glfw.KeyLeftBracket:  "(",
	glfw.KeyBackslash:    "\\",
	glfw.KeyRightBracket: ")",
	glfw.KeyGraveAccent:  "`",
	glfw.KeyWorld1:       "",
	glfw.KeyWorld2:       "",
	glfw.KeyEscape:       "esc",
	glfw.KeyEnter:        "enter",
	glfw.KeyTab:          "tab",
	glfw.KeyBackspace:    "backspace",
	glfw.KeyInsert:       "insert",
	glfw.KeyDelete:       "delete",
	glfw.KeyRight:        "right",
	glfw.KeyLeft:         "left",
	glfw.KeyDown:         "down",
	glfw.KeyUp:           "up",
	glfw.KeyPageUp:       "page up",
	glfw.KeyPageDown:     "page down",
	glfw.KeyHome:         "home",
	glfw.KeyEnd:          "end",
	glfw.KeyCapsLock:     "caps lock",
	glfw.KeyScrollLock:   "scroll lock",
	glfw.KeyNumLock:      "num lock",
	glfw.KeyPrintScreen:  "print screen",
	glfw.KeyPause:        "pause",
	glfw.KeyF1:           "f1",
	glfw.KeyF2:           "f2",
	glfw.KeyF3:           "f3",
	glfw.KeyF4:           "f4",
	glfw.KeyF5:           "f5",
	glfw.KeyF6:           "f6",
	glfw.KeyF7:           "f7",
	glfw.KeyF8:           "f8",
	glfw.KeyF9:           "f9",
	glfw.KeyF10:          "f10",
	glfw.KeyF11:          "f11",
	glfw.KeyF12:          "f12",
	glfw.KeyF13:          "f13",
	glfw.KeyF14:          "f14",
	glfw.KeyF15:          "f15",
	glfw.KeyF16:          "f16",
	glfw.KeyF17:          "f17",
	glfw.KeyF18:          "f18",
	glfw.KeyF19:          "f19",
	glfw.KeyF20:          "f20",
	glfw.KeyF21:          "f21",
	glfw.KeyF22:          "f22",
	glfw.KeyF23:          "f23",
	glfw.KeyF24:          "f24",
	glfw.KeyF25:          "f25",
	glfw.KeyKP0:          "keypad 0",
	glfw.KeyKP1:          "keypad 1",
	glfw.KeyKP2:          "keypad 2",
	glfw.KeyKP3:          "keypad 3",
	glfw.KeyKP4:          "keypad 4",
	glfw.KeyKP5:          "keypad 5",
	glfw.KeyKP6:          "keypad 6",
	glfw.KeyKP7:          "keypad 7",
	glfw.KeyKP8:          "keypad 8",
	glfw.KeyKP9:          "keypad 9",
	glfw.KeyKPDecimal:    "keypad .",
	glfw.KeyKPDivide:     "keypad /",
	glfw.KeyKPMultiply:   "keypad *",
	glfw.KeyKPSubtract:   "keypad -",
	glfw.KeyKPAdd:        "keypad +",
	glfw.KeyKPEnter:      "keypad enter",
	glfw.KeyKPEqual:      "keypad =",
	glfw.KeyLeftShift:    "left shift",
	glfw.KeyLeftControl:  "left control",
	glfw.KeyLeftAlt:      "left alt",
	glfw.KeyLeftSuper:    "left super",
	glfw.KeyRightShift:   "right shift",
	glfw.KeyRightControl: "right control",
	glfw.KeyRightAlt:     "right alt",
	glfw.KeyRightSuper:   "right super",
	glfw.KeyMenu:         "menu",
}

type Config struct {
	Context       *Context
	Width, Height int
	WindowMode    int
	WindowTitle   string
	FrameLimit    int
}

type App struct {
	window *glfw.Window

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

	switch a.winMode {
	case WINDOW_MODE_WINDOWED:
		a.window, err = glfw.CreateWindow(a.width, a.height, a.title, nil, nil)
		break
	case WINDOW_MODE_FULL_SCREEN:
		a.window, err = glfw.CreateWindow(a.width, a.height, a.title, glfw.GetPrimaryMonitor(), nil)
		break
	default:
		return errors.New("window uninitialized")
	}
	if err != nil {
		return errors.New(fmt.Sprint("unable to init window:", err))
	}
	a.window.MakeContextCurrent()

	a.window.SetKeyCallback(a.keyCallBack)
	a.window.SetCursorPosCallback(a.cursorPos)

	scrWidth, scrHeight := a.window.GetFramebufferSize()
	a.context.SetScreenSize(scrWidth, scrHeight)

	if err := a.context.renderer.Init(a.context); err != nil {
		return err
	}

	fmt.Println("OpenGL version", a.context.renderer.Version())

	glfw.SwapInterval(1)

	fps := 0
	fpsDisplayLastTime := a.currentTime

	a.currentTime = glfw.GetTime()

	// initialize input
	a.context.input.frameStart(a.currentTime)
	a.context.input.curMouseX, a.context.input.curMouseY = a.window.GetCursorPos()
	a.context.input.frameEnd()

	for !a.window.ShouldClose() {
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
		a.window.SwapBuffers()
		glfw.PollEvents()
		a.context.input.frameStart(a.currentTime)
		a.executeBehaviors(false)
		a.context.input.frameEnd()

		switch a.context.input.cursorMode {
		case CURSOR_MODE_NORMAL:
			a.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		case CURSOR_MODE_DISABLED:
			a.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		}

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

func (a *App) Stop() {
	a.window.SetShouldClose(true)
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

func (a *App) keyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		a.context.input.curKeyState[keyMap[key]] = KEY_STATE_DOWN
	}
	if action == glfw.Release {
		a.context.input.curKeyState[keyMap[key]] = KEY_STATE_UP
	}
}

func (a *App) cursorPos(w *glfw.Window, xpos, ypos float64) {
	fmt.Println(a, xpos, ypos)
}
