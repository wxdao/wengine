package opengl

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	. "github.com/wxdao/wengine"
	"math"
)

type deferredShading struct {
	renderer *renderer

	gBuffer   uint32
	gPosition uint32
	gNormal   uint32
	gDiffuse  uint32
	gDepth    uint32

	sBuffer   uint32
	sDirMap   uint32
	sPointMap uint32
	sDepth    uint32

	quad uint32
}

func (r *deferredShading) init() error {
	if err := r.initGBuffer(); err != nil {
		return err
	}

	if err := r.initSBuffer(); err != nil {
		return err
	}

	if err := r.initQuad(); err != nil {
		return err
	}

	return nil
}

func (r *deferredShading) initGBuffer() error {
	scrWidth, scrHeight := r.renderer.context.ScreenSize()

	gl.GenFramebuffers(1, &r.gBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.gBuffer)

	gl.GenTextures(1, &r.gPosition)
	gl.BindTexture(gl.TEXTURE_2D, r.gPosition)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB32F, int32(scrWidth), int32(scrHeight), 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, r.gPosition, 0)

	gl.GenTextures(1, &r.gNormal)
	gl.BindTexture(gl.TEXTURE_2D, r.gNormal)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB32F, int32(scrWidth), int32(scrHeight), 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, r.gNormal, 0)

	gl.GenTextures(1, &r.gDiffuse)
	gl.BindTexture(gl.TEXTURE_2D, r.gDiffuse)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(scrWidth), int32(scrHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT2, gl.TEXTURE_2D, r.gDiffuse, 0)

	gl.DrawBuffers(3, &[]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}[0])

	gl.GenRenderbuffers(1, &r.gDepth)
	gl.BindRenderbuffer(gl.RENDERBUFFER, r.gDepth)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(scrWidth), int32(scrHeight))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, r.gDepth)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return errors.New("framebuffer failed")
	}

	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return nil
}

func (r *deferredShading) initSBuffer() error {
	gl.GenFramebuffers(1, &r.sBuffer)

	gl.GenTextures(1, &r.sDirMap)
	gl.BindTexture(gl.TEXTURE_2D, r.sDirMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, int32(r.renderer.shadowMapResolution), int32(r.renderer.shadowMapResolution), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &[]float32{1, 1, 1, 1}[0])

	gl.GenRenderbuffers(1, &r.sDepth)
	gl.BindRenderbuffer(gl.RENDERBUFFER, r.sDepth)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(r.renderer.shadowMapResolution), int32(r.renderer.shadowMapResolution))

	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

	return nil
}

func (r *deferredShading) initQuad() error {
	var vbo [2]uint32

	vertices := []mgl32.Vec3{
		{1, -1, 0},
		{1, 1, 0},
		{-1, 1, 0},
		{1, -1, 0},
		{-1, 1, 0},
		{-1, -1, 0},
	}

	uvs := []mgl32.Vec2{
		{1, 0},
		{1, 1},
		{0, 1},
		{1, 0},
		{0, 1},
		{0, 0},
	}

	gl.GenVertexArrays(1, &r.quad)
	gl.BindVertexArray(r.quad)

	gl.GenBuffers(2, &vbo[0])

	// vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*3*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// uvs
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(uvs)*2*4, gl.Ptr(uvs), gl.STATIC_DRAW)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

func (r *deferredShading) render(targetFBO uint32, lights []*LightComponent, meshes []*MeshComponent, scene *Scene, cameras []*CameraComponent) error {
	for _, camera := range cameras {
		err := r.geometryPass(lights, meshes, scene, camera)
		if err != nil {
			return err
		}

		err = r.scenePass(targetFBO, lights, camera)
		if err != nil {
			return err
		}

		err = r.shadowPass(targetFBO, lights, meshes, camera)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *deferredShading) geometryPass(lights []*LightComponent, meshes []*MeshComponent, scene *Scene, camera *CameraComponent) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.gBuffer)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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
		if _, exists := r.renderer.materials[mesh.Material]; !exists {
			err := r.renderer.helpLoad(mesh.Material)
			if err != nil {
				return err
			}
			continue
		}
		shader, err := r.selectShader(mesh)
		if err != nil {
			return err
		}
		if shader == nil {
			continue
		}

		uniform := deferredMeshUniform{}
		uniform.model = model
		uniform.view = view
		uniform.projection = projection

		r.applyShaderToMesh(shader, uniform, mesh)

		if err := r.renderer.drawMesh(rMesh); err != nil {
			return err
		}

		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.UseProgram(0)
	}

	return nil
}

func (r *deferredShading) scenePass(targetFBO uint32, lights []*LightComponent, camera *CameraComponent) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, targetFBO)
	gl.ClearColor(
		camera.Ambient.X(),
		camera.Ambient.Y(),
		camera.Ambient.Z(),
		1,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	if len(lights) == 0 {
		shader := &defaultDeferredShader_NOLIGHT
		gl.UseProgram(shader.program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, r.gDiffuse)

		gl.BindVertexArray(r.quad)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		gl.UseProgram(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		return nil
	}

	shader := &defaultDeferredShader
	gl.UseProgram(shader.program)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.gPosition)
	gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("gPosition\x00")), 0)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, r.gNormal)
	gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("gNormal\x00")), 1)
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, r.gDiffuse)
	gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("gDiffuse\x00")), 2)

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

	cameraPos := camera.Object().Position()
	gl.UniformMatrix3fv(gl.GetUniformLocation(
		shader.program, gl.Str("cameraPosition\x00")),
		1,
		false,
		&cameraPos[0],
	)

	gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str("ambient\x00")), 1, &camera.Ambient[0])

	if len(dirLights) > 0 {
		num_dirLight := int(math.Min(10, float64(len(dirLights))))
		for i := 0; i < num_dirLight; i++ {
			light := dirLights[i]
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].position\x00", i))),
				1,
				&light.position[0],
			)
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].direction\x00", i))),
				1,
				&light.direction[0],
			)
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].diffuse\x00", i))),
				1,
				&light.diffuse[0],
			)
			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("dirLights[%d].specular\x00", i))),
				1,
				&light.specular[0],
			)
		}
		gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("num_dirLight\x00")), int32(num_dirLight))
	}

	if len(pointLights) > 0 {
		num_pointLight := int(math.Min(10, float64(len(pointLights))))
		for i := 0; i < num_pointLight; i++ {
			light := pointLights[i]
			gl.Uniform3fv(
				gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].position\x00", i))),
				1,
				&light.position[0],
			)
			gl.Uniform1f(
				gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].range\x00", i))),
				light.lightRange,
			)
			gl.Uniform3fv(
				gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].diffuse\x00", i))),
				1,
				&light.diffuse[0],
			)
			gl.Uniform3fv(
				gl.GetUniformLocation(shader.program, gl.Str(fmt.Sprintf("pointLights[%d].specular\x00", i))),
				1,
				&light.specular[0],
			)
		}
		gl.Uniform1i(
			gl.GetUniformLocation(shader.program, gl.Str("num_pointLight\x00")),
			int32(num_pointLight),
		)
	}

	gl.BindVertexArray(r.quad)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

	gl.UseProgram(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return nil
}

func (r *deferredShading) shadowPass(targetFBO uint32, lights []*LightComponent, meshes []*MeshComponent, camera *CameraComponent) error {
	for _, light := range lights {
		if light.ShadowType == LIGHT_SHADOW_TYPE_NONE {
			continue
		}
		switch light.LightSource {
		case LIGHT_SOURCE_DIRECTIONAL:
			// generate shadow map
			gl.Viewport(0, 0, int32(r.renderer.shadowMapResolution), int32(r.renderer.shadowMapResolution))
			gl.BindFramebuffer(gl.FRAMEBUFFER, r.sBuffer)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, r.sDirMap, 0)
			gl.DrawBuffer(gl.NONE)
			gl.ReadBuffer(gl.NONE)
			if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
				return errors.New("framebuffer failed")
			}
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			shader := &defaultShadowMapShader_DIRLIGHT
			gl.UseProgram(shader.program)
			lightMatrix := mgl32.Ortho(
				-100,
				100,
				-100,
				100,
				0.1,
				200,
			).Mul4(mgl32.LookAtV(
				light.Object().Position(),
				light.Object().Position().Add(light.Object().Forward()),
				light.Object().Up(),
			))
			gl.UniformMatrix4fv(gl.GetUniformLocation(shader.program, gl.Str("lightMatrix\x00")), 1, false, &lightMatrix[0])

			for _, mesh := range meshes {
				if !mesh.CastShadow {
					continue
				}
				rMesh, exists := r.renderer.meshes[mesh.Mesh]
				if !exists {
					err := r.renderer.helpLoad(mesh.Mesh)
					if err != nil {
						return err
					}
					continue
				}
				model := mesh.Object().ModelMatrix()
				gl.UniformMatrix4fv(gl.GetUniformLocation(shader.program, gl.Str("model\x00")), 1, false, &model[0])
				if err := r.renderer.drawMesh(rMesh); err != nil {
					return err
				}
			}

			scrWidth, scrHeight := r.renderer.context.ScreenSize()
			gl.Viewport(0, 0, int32(scrWidth), int32(scrHeight))

			// blend shadow
			gl.BindFramebuffer(gl.FRAMEBUFFER, targetFBO)
			gl.Enable(gl.BLEND)
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

			shader = &defaultBlendShadowShader_DIRLIGHT_DEFERRED
			gl.UseProgram(shader.program)

			gl.Uniform3fv(gl.GetUniformLocation(shader.program, gl.Str("ambient\x00")), 1, &camera.Ambient[0])
			gl.UniformMatrix4fv(gl.GetUniformLocation(shader.program, gl.Str("lightMatrix\x00")), 1, false, &lightMatrix[0])

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, r.sDirMap)
			gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("sDirMap\x00")), 0)
			gl.ActiveTexture(gl.TEXTURE1)
			gl.BindTexture(gl.TEXTURE_2D, r.gPosition)
			gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("gPosition\x00")), 1)
			gl.ActiveTexture(gl.TEXTURE2)
			gl.BindTexture(gl.TEXTURE_2D, r.gDiffuse)
			gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("gDiffuse\x00")), 2)

			gl.BindVertexArray(r.quad)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
			gl.BindVertexArray(0)

			gl.UseProgram(0)
			gl.Disable(gl.BLEND)
		}
	}
	return nil
}

func (r *deferredShading) selectShader(mesh *MeshComponent) (*glShaderProgram, error) {
	if mesh.Shader != "" {
		if shader, exists := r.renderer.programs[mesh.Shader]; exists {
			return shader, nil
		} else {
			err := r.renderer.helpLoad(mesh.Shader)
			return nil, err
		}
	}
	material := r.renderer.materials[mesh.Material]
	if material.diffuseMap != 0 {
		return &defaultMeshShader_TEXTURE_DEFERRED, nil
	} else {
		return &defaultMeshShader_COLOR_DEFERRED, nil
	}
}

type deferredMeshUniform struct {
	projection, model, view mgl32.Mat4
}

func (r *deferredShading) applyShaderToMesh(shader *glShaderProgram, uniform deferredMeshUniform, mesh *MeshComponent) error {
	modelLoc := gl.GetUniformLocation(shader.program, gl.Str("model\x00"))
	TImodelLoc := gl.GetUniformLocation(shader.program, gl.Str("TImodel\x00"))
	viewLoc := gl.GetUniformLocation(shader.program, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(shader.program, gl.Str("projection\x00"))
	material := r.renderer.materials[mesh.Material]
	tiModel := uniform.model.Mat3().Inv().Transpose()

	gl.UseProgram(shader.program)

	gl.UniformMatrix4fv(modelLoc, 1, false, &uniform.model[0])
	gl.UniformMatrix3fv(TImodelLoc, 1, false, &tiModel[0])
	gl.UniformMatrix4fv(modelLoc, 1, false, &uniform.model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &uniform.view[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &uniform.projection[0])

	if mesh.ReceiveShader {
		gl.Uniform1f(gl.GetUniformLocation(shader.program, gl.Str("recvShadow\x00")), 1)
	} else {
		gl.Uniform1f(gl.GetUniformLocation(shader.program, gl.Str("recvShadow\x00")), 0)
	}

	if material.diffuseMap != 0 {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, material.diffuseMap)
		gl.Uniform1i(gl.GetUniformLocation(shader.program, gl.Str("diffuseMap\x00")), 0)
	} else {
		gl.Uniform4fv(gl.GetUniformLocation(shader.program, gl.Str("color\x00")), 1, &material.DiffuseColor[0])
	}

	return nil
}
