package wengine

import (
	"github.com/go-gl/mathgl/mgl32"
)

type ComponentMap map[int]Component

type Object struct {
	enabled  bool
	parent   *Object
	children []*Object

	translation, rotation, scale, model mgl32.Mat4

	TransformInfo

	behavior Behavior

	components ComponentMap
}

func NewObject() *Object {
	return &Object{translation: mgl32.Ident4(), rotation: mgl32.Ident4(), scale: mgl32.Ident4(), components: make(ComponentMap)}
}

func (o *Object) Enabled() bool {
	return o.enabled
}

func (o *Object) SetEnabled(enabled bool) {
	o.enabled = enabled
}

func (o *Object) SetParent(parent *Object) {
	o.parent = parent
	for _, child := range parent.children {
		if child == o {
			return
		}
	}
	parent.children = append(parent.children, o)
}

func (o *Object) TranslationMatrix() mgl32.Mat4 {
	return o.translation
}

func (o *Object) RotationMatrix() mgl32.Mat4 {
	return o.rotation
}

func (o *Object) ScaleMatrix() mgl32.Mat4 {
	return o.scale
}

func (o *Object) ModelMatrix() mgl32.Mat4 {
	return o.model
}

func (o *Object) Translate(v mgl32.Vec3) {
	o.translation = mgl32.Translate3D(v.Elem()).Mul4(o.translation)
}

func (o *Object) ResetTranslation() {
	o.translation = mgl32.Ident4()
}

func (o *Object) Rotate(angle float32, axis mgl32.Vec3) {
	o.rotation = mgl32.HomogRotate3D(angle, axis).Mul4(o.rotation)
}

func (o *Object) ResetRotation() {
	o.rotation = mgl32.Ident4()
}

func (o *Object) Scale(v mgl32.Vec3) {
	o.scale = mgl32.Scale3D(v.Elem()).Mul4(o.scale)
}

func (o *Object) ResetScale() {
	o.scale = mgl32.Ident4()
}

func (o *Object) SetBehavior(behavior Behavior) {
	o.behavior = behavior
}

func (o *Object) AttachComponent(component Component) {
	component.setObject(o)
	o.components[component.Type()] = component
}

func (o *Object) Components() ComponentMap {
	return o.components
}

type TransformInfo struct {
	position, up, right, forward mgl32.Vec3
}

func (t *TransformInfo) Position() mgl32.Vec3 {
	return t.position
}

func (t *TransformInfo) Up() mgl32.Vec3 {
	return t.up
}

func (t *TransformInfo) Right() mgl32.Vec3 {
	return t.right
}

func (t *TransformInfo) Forward() mgl32.Vec3 {
	return t.forward
}

type BehaviorContext struct {
	Context   *Context
	Object    *Object
	DeltaTime float64
	Time      float64
}

type Behavior interface {
	Start(bctx *BehaviorContext)
	Update(bctx *BehaviorContext)
}
