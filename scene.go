package wengine

import "github.com/go-gl/mathgl/mgl32"

type ObjectMap map[string]*Object

type Scene struct {
	objects ObjectMap
}

func NewScene() *Scene {
	return &Scene{objects: make(ObjectMap)}
}

func (s *Scene) Objects() ObjectMap {
	return s.objects
}

func (s *Scene) RegisterObject(name string, object *Object) {
	s.objects[name] = object
}

func (s *Scene) updateTransforms() {
	for _, obj := range s.Objects() {
		if !obj.enabled {
			continue
		}
		model := s.buildTransform(obj)
		obj.model = model
		obj.position = model.Mul4x1(mgl32.Vec4{0, 0, 0, 1}).Vec3()
		obj.up = model.Mul4x1(mgl32.Vec4{0, 1, 0, 0}).Vec3()
		obj.right = model.Mul4x1(mgl32.Vec4{1, 0, 0, 0}).Vec3()
		obj.forward = obj.up.Cross(obj.right)
	}
}

func (s *Scene) buildTransform(object *Object) mgl32.Mat4 {
	model := mgl32.Ident4()
	for t := object; t != nil; t = t.parent {
		model = t.TranslationMatrix().
			Mul4(t.RotationMatrix()).
			Mul4(t.ScaleMatrix()).
			Mul4(model)
	}
	return model
}
