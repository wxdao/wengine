package wengine

type Renderer interface {
	Init(context *Context) error
	Version() string
	Render(scene *Scene) error
	NotifyInstall(assets []string) error
}

var (
	registeredRenderers = map[string]Renderer{}
)

func RegisterRenderer(name string, renderer Renderer) {
	registeredRenderers[name] = renderer
}

type RendererSetting struct {
}
