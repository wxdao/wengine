package wengine

import "github.com/go-gl/mathgl/mgl32"

const (
	COMPO_CAMERA = iota
	COMPO_LIGHT
	COMPO_MESH
)

type Component interface {
	Type() int
	Object() *Object

	setObject(object *Object)
}

type componentBase struct {
	parentObject *Object
}

func (c *componentBase) setObject(object *Object) {
	c.parentObject = object
}

func (c *componentBase) Object() *Object {
	return c.parentObject
}

const (
	CAMERA_MODE_PERSPECTIVE = iota
	CAMERA_MODE_ORTHOGRAPHIC
)

type CameraComponent struct {
	Depth     int
	Mode      int
	NearPlane float32
	FarPlane  float32
	Ambient   mgl32.Vec3

	// perspective only
	FOV float32

	// orthographic only
	Width float32

	componentBase
}

func (CameraComponent) Type() int {
	return COMPO_CAMERA
}

type MeshComponent struct {
	Mesh     string
	Material string
	Shader   string

	CastShadow    bool
	ReceiveShader bool

	componentBase
}

func (MeshComponent) Type() int {
	return COMPO_MESH
}

const (
	LIGHT_SOURCE_DIRECTIONAL = iota
	LIGHT_SOURCE_POINT
)

const (
	LIGHT_SHADOW_TYPE_NONE = iota
	LIGHT_SHADOW_TYPE_SOFT
	LIGHT_SHADOW_TYPE_HARD
)

type LightComponent struct {
	LightSource int
	ShadowType  int

	// common values
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	// point light only
	Range float32

	componentBase
}

func (LightComponent) Type() int {
	return COMPO_LIGHT
}
