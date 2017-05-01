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

var keyMap = map[int]string{
	int(glfw.KeyUnknown):      "unknown",
	int(glfw.KeySpace):        "space",
	int(glfw.KeyApostrophe):   "",
	int(glfw.KeyComma):        ",",
	int(glfw.KeyMinus):        "-",
	int(glfw.KeyPeriod):       ".",
	int(glfw.KeySlash):        "-",
	int(glfw.Key0):            "0",
	int(glfw.Key1):            "1",
	int(glfw.Key2):            "2",
	int(glfw.Key3):            "3",
	int(glfw.Key4):            "4",
	int(glfw.Key5):            "5",
	int(glfw.Key6):            "6",
	int(glfw.Key7):            "7",
	int(glfw.Key8):            "8",
	int(glfw.Key9):            "9",
	int(glfw.KeySemicolon):    ";",
	int(glfw.KeyEqual):        "=",
	int(glfw.KeyA):            "a",
	int(glfw.KeyB):            "b",
	int(glfw.KeyC):            "c",
	int(glfw.KeyD):            "d",
	int(glfw.KeyE):            "e",
	int(glfw.KeyF):            "f",
	int(glfw.KeyG):            "g",
	int(glfw.KeyH):            "h",
	int(glfw.KeyI):            "i",
	int(glfw.KeyJ):            "j",
	int(glfw.KeyK):            "k",
	int(glfw.KeyL):            "l",
	int(glfw.KeyM):            "m",
	int(glfw.KeyN):            "n",
	int(glfw.KeyO):            "o",
	int(glfw.KeyP):            "p",
	int(glfw.KeyQ):            "q",
	int(glfw.KeyR):            "r",
	int(glfw.KeyS):            "s",
	int(glfw.KeyT):            "t",
	int(glfw.KeyU):            "u",
	int(glfw.KeyV):            "v",
	int(glfw.KeyW):            "w",
	int(glfw.KeyX):            "x",
	int(glfw.KeyY):            "y",
	int(glfw.KeyZ):            "z",
	int(glfw.KeyLeftBracket):  "(",
	int(glfw.KeyBackslash):    "\\",
	int(glfw.KeyRightBracket): ")",
	int(glfw.KeyGraveAccent):  "`",
	int(glfw.KeyWorld1):       "",
	int(glfw.KeyWorld2):       "",
	int(glfw.KeyEscape):       "esc",
	int(glfw.KeyEnter):        "enter",
	int(glfw.KeyTab):          "tab",
	int(glfw.KeyBackspace):    "backspace",
	int(glfw.KeyInsert):       "insert",
	int(glfw.KeyDelete):       "delete",
	int(glfw.KeyRight):        "right",
	int(glfw.KeyLeft):         "left",
	int(glfw.KeyDown):         "down",
	int(glfw.KeyUp):           "up",
	int(glfw.KeyPageUp):       "page up",
	int(glfw.KeyPageDown):     "page down",
	int(glfw.KeyHome):         "home",
	int(glfw.KeyEnd):          "end",
	int(glfw.KeyCapsLock):     "caps lock",
	int(glfw.KeyScrollLock):   "scroll lock",
	int(glfw.KeyNumLock):      "num lock",
	int(glfw.KeyPrintScreen):  "print screen",
	int(glfw.KeyPause):        "pause",
	int(glfw.KeyF1):           "f1",
	int(glfw.KeyF2):           "f2",
	int(glfw.KeyF3):           "f3",
	int(glfw.KeyF4):           "f4",
	int(glfw.KeyF5):           "f5",
	int(glfw.KeyF6):           "f6",
	int(glfw.KeyF7):           "f7",
	int(glfw.KeyF8):           "f8",
	int(glfw.KeyF9):           "f9",
	int(glfw.KeyF10):          "f10",
	int(glfw.KeyF11):          "f11",
	int(glfw.KeyF12):          "f12",
	int(glfw.KeyF13):          "f13",
	int(glfw.KeyF14):          "f14",
	int(glfw.KeyF15):          "f15",
	int(glfw.KeyF16):          "f16",
	int(glfw.KeyF17):          "f17",
	int(glfw.KeyF18):          "f18",
	int(glfw.KeyF19):          "f19",
	int(glfw.KeyF20):          "f20",
	int(glfw.KeyF21):          "f21",
	int(glfw.KeyF22):          "f22",
	int(glfw.KeyF23):          "f23",
	int(glfw.KeyF24):          "f24",
	int(glfw.KeyF25):          "f25",
	int(glfw.KeyKP0):          "keypad 0",
	int(glfw.KeyKP1):          "keypad 1",
	int(glfw.KeyKP2):          "keypad 2",
	int(glfw.KeyKP3):          "keypad 3",
	int(glfw.KeyKP4):          "keypad 4",
	int(glfw.KeyKP5):          "keypad 5",
	int(glfw.KeyKP6):          "keypad 6",
	int(glfw.KeyKP7):          "keypad 7",
	int(glfw.KeyKP8):          "keypad 8",
	int(glfw.KeyKP9):          "keypad 9",
	int(glfw.KeyKPDecimal):    "keypad .",
	int(glfw.KeyKPDivide):     "keypad /",
	int(glfw.KeyKPMultiply):   "keypad *",
	int(glfw.KeyKPSubtract):   "keypad -",
	int(glfw.KeyKPAdd):        "keypad +",
	int(glfw.KeyKPEnter):      "keypad enter",
	int(glfw.KeyKPEqual):      "keypad =",
	int(glfw.KeyLeftShift):    "left shift",
	int(glfw.KeyLeftControl):  "left control",
	int(glfw.KeyLeftAlt):      "left alt",
	int(glfw.KeyLeftSuper):    "left super",
	int(glfw.KeyRightShift):   "right shift",
	int(glfw.KeyRightControl): "right control",
	int(glfw.KeyRightAlt):     "right alt",
	int(glfw.KeyRightSuper):   "right super",
	int(glfw.KeyMenu):         "menu",

	int(glfw.MouseButton1): "mouse 1",
	int(glfw.MouseButton2): "mouse 2",
	int(glfw.MouseButton3): "mouse 3",
	int(glfw.MouseButton4): "mouse 4",
	int(glfw.MouseButton5): "mouse 5",
	int(glfw.MouseButton6): "mouse 6",
	int(glfw.MouseButton7): "mouse 7",
	int(glfw.MouseButton8): "mouse 8",
}

type Config struct {
	Context       *Context
	Width, Height int
	WindowMode    int
	WindowTitle   string
	FrameLimit    int
	VSync         bool
}

type App struct {
	window *glfw.Window

	width, height int
	winMode       int
	title         string
	frameLimit    int
	vSync         bool

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
		vSync:      config.VSync,
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
	a.window.SetMouseButtonCallback(a.mouseCallBack)

	scrWidth, scrHeight := a.window.GetFramebufferSize()
	a.context.SetScreenSize(scrWidth, scrHeight)

	if err := a.context.renderer.Init(a.context); err != nil {
		return err
	}

	fmt.Println("OpenGL version", a.context.renderer.Version())

	if a.vSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}

	fps := 0
	fpsDisplayLastTime := a.currentTime

	a.currentTime = glfw.GetTime()

	// initialize input
	a.context.input.curMouseX, a.context.input.curMouseY = a.window.GetCursorPos()
	a.context.input.frameStart(a.currentTime)
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
		a.context.input.curMouseX, a.context.input.curMouseY = a.window.GetCursorPos()
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
		a.context.input.curKeyState[keyMap[int(key)]] = KEY_STATE_DOWN
	}
	if action == glfw.Release {
		a.context.input.curKeyState[keyMap[int(key)]] = KEY_STATE_UP
	}
}

func (a *App) mouseCallBack(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		a.context.input.curKeyState[keyMap[int(key)]] = KEY_STATE_DOWN
	}
	if action == glfw.Release {
		a.context.input.curKeyState[keyMap[int(key)]] = KEY_STATE_UP
	}
}
