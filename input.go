package wengine

import "math"

const (
	KEY_STATE_UP = iota
	KEY_STATE_DOWN
)

const (
	AXIS_SOURCE_KEY = iota
	AXIS_SOURCE_MOUSE
)

const (
	AXIS_FROM_X = iota
	AXIS_FROM_Y
)

type AxisMeta struct {
	Source int

	// key
	PositiveKey, NegativeKey string

	// mouse & joystick
	From int

	Gravity, Dead, Sensitivity float64
	Invert                     bool
}

const (
	CURSOR_MODE_NORMAL = iota
	CURSOR_MODE_DISABLED
)

type Input struct {
	preKeyState map[string]int
	curKeyState map[string]int

	preMouseX, curMouseX float64
	preMouseY, curMouseY float64

	axes      map[string][]*AxisMeta
	axisValue map[*AxisMeta]float64

	cursorMode int

	currentTime, lastTime float64
}

func newInput() *Input {
	return &Input{
		preKeyState: map[string]int{},
		curKeyState: map[string]int{},
		axes:        map[string][]*AxisMeta{},
		axisValue:   map[*AxisMeta]float64{},
	}
}

func (i *Input) frameStart(currentTime float64) {
	i.currentTime = currentTime
	if i.lastTime == 0 {
		i.lastTime = i.currentTime
	}

	for meta, value := range i.axisValue {
		switch meta.Source {
		case AXIS_SOURCE_KEY:
			var deltaHold float64
			if meta.PositiveKey != "" {
				if i.GetKey(meta.PositiveKey) && !i.GetKeyDown(meta.PositiveKey) {
					deltaHold += i.currentTime - i.lastTime
				}
			}
			if meta.NegativeKey != "" {
				if i.GetKey(meta.NegativeKey) && !i.GetKeyDown(meta.NegativeKey) {
					deltaHold -= i.currentTime - i.lastTime
				}
			}
			if deltaHold == 0 {
				if value > 0 {
					value -= math.Min(meta.Gravity*(i.currentTime-i.lastTime), value)
				} else {
					value += math.Min(meta.Gravity*(i.currentTime-i.lastTime), -value)
				}
			} else {
				if value*deltaHold < 0 {
					value = 0
				}
				value += deltaHold * meta.Sensitivity
			}
			if value > 1 {
				value = 1
			}
			if value < -1 {
				value = -1
			}
			i.axisValue[meta] = value
		}
	}
}

func (i *Input) frameEnd() {
	for k, v := range i.curKeyState {
		i.preKeyState[k] = v
	}
	i.preMouseX, i.preMouseY = i.curMouseX, i.curMouseY
	i.lastTime = i.currentTime
}

func (i *Input) SetCursorMode(mode int) {
	i.cursorMode = mode
}

func (i *Input) GetKey(key string) bool {
	if i.curKeyState[key] == KEY_STATE_DOWN {
		return true
	}
	return false
}

func (i *Input) GetKeyDown(key string) bool {
	if i.curKeyState[key] == KEY_STATE_DOWN && i.preKeyState[key] == KEY_STATE_UP {
		return true
	}
	return false
}

func (i *Input) GetKeyUp(key string) bool {
	if i.curKeyState[key] == KEY_STATE_UP && i.preKeyState[key] == KEY_STATE_DOWN {
		return true
	}
	return false
}

func (i *Input) BindAxis(axis string, meta AxisMeta) {
	i.axes[axis] = append(i.axes[axis], &meta)
	i.axisValue[&meta] = 0
}

func (i *Input) ResetAxis(axis string) {
	for _, meta := range i.axes[axis] {
		delete(i.axisValue, meta)
	}
	delete(i.axes, axis)
}

func (i *Input) GetAxis(axis string) (final float64) {
	metas, exists := i.axes[axis]
	if !exists {
		return
	}
	for _, meta := range metas {
		if math.Abs(i.axisValue[meta]) >= math.Abs(final) {
			final = i.axisValue[meta]
		}
	}

	return
}
