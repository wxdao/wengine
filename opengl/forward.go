package opengl

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	. "github.com/wxdao/wengine"
	"math"
)

type forwardShading struct {
	renderer *renderer
}

func (r *forwardShading) init() error {
	return nil
}

func (r *forwardShading) render(targetFBO uint32, lights []*LightComponent, meshes []*MeshComponent, scene *Scene, cameras []*CameraComponent) error {
	for _, camera := range cameras {
		err := r.scenePass(targetFBO, lights, meshes, scene, camera)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *forwardShading) scenePass(targetFBO uint32, lights []*LightComponent, meshes []*MeshComponent, scene *Scene, camera *CameraComponent) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, targetFBO)
	gl.ClearColor(
		camera.Ambient.X(),
		camera.Ambient.Y(),
		camera.Ambient.Z(),
		1,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// build view & projection matrices
	cameraObj := camera.Object()
	view := mgl32.LookAtV(cameraObj.Position(), cameraObj.Position().Add(cameraObj.Forward()), cameraObj.Up())
	scrWidth, scrHeight := r.renderer.context.ScreenSize()
	var projection mgl32.Mat4
	switch camera.Mode {
	case CAMERA_MODE_PERSPECTIVE:
		projection = mgl32.Perspective(
			camera.FOV,
			float32(scrWidth)/float32(scrHeight),
			camera.NearPlane,
			camera.FarPlane,
		)
	case CAMERA_MODE_ORTHOGRAPHIC:
		projection = mgl32.Ortho(
			-camera.Width/2,
			camera.Width/2,
			-camera.Width*float32(scrHeight)/float32(scrWidth)/2,
			camera.Width*float32(scrHeight)/float32(scrWidth)/2,
			camera.NearPlane,
			camera.FarPlane,
		)
	}
	// render meshes
	for _, mesh := range meshes {
		model := mesh.Object().ModelMatrix()
		rMesh, exists := r.renderer.meshes[mesh.Mesh]
		if !exists {
			err := r.renderer.helpLoad(mesh.Mesh)
			if err != nil {
				return err
			}
			continue
		}
		if mesh.Material == "" {
			return errors.New("mesh with no material")
		}
		if _, exists := r.renderer.meshMaterials[mesh.Material]; !exists {
			err := r.renderer.helpLoad(mesh.Material)
			if err != nil {
				return err
			}
			continue
		}
		shader, err := r.selectShader(mesh, lights)
		if err != nil {
			return err
		}
		if shader == nil {
			continue
		}

		uniform := forwardMeshUniform{}
		uniform.model = model
		uniform.view = view
		uniform.projection = projection
		uniform.cameraPosition = cameraObj.Position()
		if len(lights) > 0 {
			dirLights := []dirLightUniform{}
			pointLights := []pointLightUniform{}
			for _, light := range lights {
				switch light.LightSource {
				case LIGHT_SOURCE_DIRECTIONAL:
					dirLights = append(dirLights, dirLightUniform{
						position:  light.Object().Position(),
						direction: light.Object().Forward(),
						diffuse:   light.Diffuse,
						specular:  light.Specular,
					})
				case LIGHT_SOURCE_POINT:
					pointLights = append(pointLights, pointLightUniform{
						position:   light.Object().Position(),
						lightRange: light.Range,
						diffuse:    light.Diffuse,
						specular:   light.Specular,
					})
				}
			}
			uniform.dirLights = dirLights
			uniform.pointLights = pointLights
		}

		r.applyShaderToMesh(shader, uniform, mesh, camera)

		if err := r.renderer.drawMesh(rMesh); err != nil {
			return err
		}

		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.UseProgram(0)
	}
	return nil
}

func (r *forwardShading) selectShader(mesh *MeshComponent, lights []*LightComponent) (*glShaderProgram, error) {
	if mesh.Shader != "" {
		if shader, exists := r.renderer.programs[mesh.Shader]; exists {
			return shader, nil
		} else {
			err := r.renderer.helpLoad(mesh.Shader)
			return nil, err
		}
	}
	hasLights := len(lights) > 0
	material := r.renderer.meshMaterials[mesh.Material]
	if material.diffuseMap != 0 {
		if hasLights {
			return defaultShaders["mesh_texture"], nil
		}
		return defaultShaders["mesh_texture_nolight"], nil
	} else {
		if hasLights {
			return defaultShaders["mesh_color"], nil
		}
		return defaultShaders["mesh_color_nolight"], nil
	}
}

type dirLightUniform struct {
	position  mgl32.Vec3
	direction mgl32.Vec3
	diffuse   mgl32.Vec3
	specular  mgl32.Vec3
}

type pointLightUniform struct {
	position   mgl32.Vec3
	lightRange float32
	diffuse    mgl32.Vec3
	specular   mgl32.Vec3
}

type spotLightUniform struct {
	position   mgl32.Vec3
	direction  mgl32.Vec3
	cosAngle   float32
	lightRange float32
	diffuse    mgl32.Vec3
	specular   mgl32.Vec3
}

type forwardMeshUniform struct {
	projection, model, view mgl32.Mat4

	cameraPosition mgl32.Vec3

	dirLights   []dirLightUniform
	pointLights []pointLightUniform
}

func (r *forwardShading) applyShaderToMesh(shader *glShaderProgram, uniform forwardMeshUniform, mesh *MeshComponent, camera *CameraComponent) error {
	modelLoc := gl.GetUniformLocation(shader.program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(shader.program, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(shader.program, gl.Str("projection\x00"))
	cameraPositionLoc := gl.GetUniformLocation(shader.program, gl.Str("cameraPosition\x00"))
	material := r.renderer.meshMaterials[mesh.Material]

	gl.UseProgram(shader.program)

	gl.UniformMatrix4fv(modelLoc, 1, false, &uniform.model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &uniform.view[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &uniform.projection[0])
	gl.UniformMatrix3fv(cameraPositionLoc, 1, false, &uniform.cameraPosition[0])

	if material.diffuseMap != 0 {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, material.diffuseMap)
		gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("diffuseMap\x00")), 0)
	} else {
		gl.Uniform4fv(gl.GetUniformLocation(shader.program, gl.Str("color\x00")), 1, &material.DiffuseColor[0])
	}

	gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str("ambient\x00")), 1, &camera.Ambient[0])

	if len(uniform.dirLights) > 0 {
		num_dirLight := int(math.Min(10, float64(len(uniform.dirLights))))
		for i := 0; i < num_dirLight; i++ {
			light := uniform.dirLights[i]
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].position\x00", i))), 1, &light.position[0])
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].direction\x00", i))), 1, &light.direction[0])
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].diffuse\x00", i))), 1, &light.diffuse[0])
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].specular\x00", i))), 1, &light.specular[0])
		}
		gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("num_dirLight\x00")), int32(num_dirLight))
	}

	if len(uniform.pointLights) > 0 {
		num_pointLight := int(math.Min(10, float64(len(uniform.pointLights))))
		for i := 0; i < num_pointLight; i++ {
			light := uniform.pointLights[i]
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].position\x00", i))), 1, &light.position[0])
			gl.Uniform1f(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].range\x00", i))), light.lightRange)
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].diffuse\x00", i))), 1, &light.diffuse[0])
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].specular\x00", i))), 1, &light.specular[0])
		}
		gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("num_pointLight\x00")), int32(num_pointLight))
	}

	return nil
}
