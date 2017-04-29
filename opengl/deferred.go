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

	sDirBuffer   uint32
	sDirMap      uint32
	sPointBuffer uint32
	sPointMap    uint32
	sPointDepth  uint32
	sSpotBuffer  uint32
	sSpotMap     uint32
	sSpotDepth   uint32

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
	// directional light
	gl.GenFramebuffers(1, &r.sDirBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.sDirBuffer)

	gl.GenTextures(1, &r.sDirMap)
	gl.BindTexture(gl.TEXTURE_2D, r.sDirMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, int32(r.renderer.dirLightShadowMapResolution), int32(r.renderer.dirLightShadowMapResolution), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &[]float32{1, 1, 1, 1}[0])

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, r.sDirMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return errors.New("framebuffer failed")
	}

	// spot light
	gl.GenFramebuffers(1, &r.sSpotBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.sSpotBuffer)

	gl.GenTextures(1, &r.sSpotMap)
	gl.BindTexture(gl.TEXTURE_2D, r.sSpotMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32F, int32(r.renderer.spotLightShadowMapResolution), int32(r.renderer.spotLightShadowMapResolution), 0, gl.RED, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &[]float32{1, 1, 1, 1}[0])

	gl.GenRenderbuffers(1, &r.sSpotDepth)
	gl.BindRenderbuffer(gl.RENDERBUFFER, r.sSpotDepth)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(r.renderer.spotLightShadowMapResolution), int32(r.renderer.spotLightShadowMapResolution))

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, r.sSpotMap, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, r.sSpotDepth)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return errors.New("framebuffer failed")
	}

	// point light
	gl.GenFramebuffers(1, &r.sPointBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.sPointBuffer)

	gl.GenTextures(1, &r.sPointMap)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.sPointMap)
	for i := 0; i < 6; i++ {
		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, gl.R16F, int32(r.renderer.pointLightShadowMapResolution), int32(r.renderer.pointLightShadowMapResolution), 0, gl.RED, gl.FLOAT, nil)
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.GenTextures(1, &r.sPointDepth)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.sPointDepth)
	for i := 0; i < 6; i++ {
		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, gl.DEPTH_COMPONENT, int32(r.renderer.pointLightShadowMapResolution), int32(r.renderer.pointLightShadowMapResolution), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, r.sPointMap, 0)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, r.sPointDepth, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println(gl.CheckFramebufferStatus(gl.FRAMEBUFFER))
		return errors.New("framebuffer failed")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
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

		err = r.blendAmbient(targetFBO, camera)
		if err != nil {
			return err
		}

		err = r.lightsPass(targetFBO, lights, meshes, camera)
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

func (r *deferredShading) blendAmbient(targetFBO uint32, camera *CameraComponent) error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, targetFBO)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)

	shader := defaultShaders["deferred_blend_ambient"]
	gl.UseProgram(shader.program)

	gl.Uniform3fv(shader.getLocation("ambient"), 1, &camera.Ambient[0])
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.gDiffuse)
	gl.Uniform1i(shader.getLocation("gDiffuse"), 0)

	gl.BindVertexArray(r.quad)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

	gl.UseProgram(0)
	gl.Disable(gl.BLEND)

	return nil
}

func (r *deferredShading) lightsPass(targetFBO uint32, lights []*LightComponent, meshes []*MeshComponent, camera *CameraComponent) error {
	for _, light := range lights {
		var shadowMapShader, shader *glShaderProgram
		switch light.LightSource {
		case LIGHT_SOURCE_DIRECTIONAL:
			switch light.ShadowType {
			case LIGHT_SHADOW_TYPE_NONE:
				shader = defaultShaders["deferred_dirLight_noshadow"]
				gl.UseProgram(shader.program)
			default:
				shadowMapShader = defaultShaders["shadow_map_dirLight"]
				lightMatrix, err := r.generateDirLightShadowMap(shadowMapShader, light, meshes)
				if err != nil {
					return err
				}
				shader = defaultShaders["deferred_dirLight"]
				gl.UseProgram(shader.program)

				gl.UniformMatrix4fv(shader.getLocation("lightMatrix"), 1, false, &lightMatrix[0])
				gl.ActiveTexture(gl.TEXTURE3)
				gl.BindTexture(gl.TEXTURE_2D, r.sDirMap)
				gl.Uniform1i(shader.getLocation("sDirMap"), 3)
			}

			uniform := dirLightUniform{
				position:  light.Object().Position(),
				direction: light.Object().Forward(),
				diffuse:   light.Diffuse,
				specular:  light.Specular,
			}
			gl.Uniform3fv(
				shader.getLocation("dirLight.position"),
				1,
				&uniform.position[0],
			)
			gl.Uniform3fv(
				shader.getLocation("dirLight.direction"),
				1,
				&uniform.direction[0],
			)
			gl.Uniform3fv(
				shader.getLocation("dirLight.diffuse"),
				1,
				&uniform.diffuse[0],
			)
			gl.Uniform3fv(
				shader.getLocation("dirLight.specular"),
				1,
				&uniform.specular[0],
			)
		case LIGHT_SOURCE_POINT:
			switch light.ShadowType {
			case LIGHT_SHADOW_TYPE_NONE:
				shader = defaultShaders["deferred_pointLight_noshadow"]
				gl.UseProgram(shader.program)
			default:
				shadowMapShader = defaultShaders["shadow_map_pointLight"]
				err := r.generatePointLightShadowMap(shadowMapShader, light, meshes)
				if err != nil {
					return err
				}
				shader = defaultShaders["deferred_pointLight"]
				gl.UseProgram(shader.program)

				gl.ActiveTexture(gl.TEXTURE3)
				gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.sPointMap)
				gl.Uniform1i(shader.getLocation("sPointMap"), 3)
			}

			uniform := pointLightUniform{
				position:   light.Object().Position(),
				lightRange: light.Range,
				diffuse:    light.Diffuse,
				specular:   light.Specular,
			}
			gl.Uniform3fv(
				shader.getLocation("pointLight.position"),
				1,
				&uniform.position[0],
			)
			gl.Uniform1f(
				shader.getLocation("pointLight.range"),
				uniform.lightRange,
			)
			gl.Uniform3fv(
				shader.getLocation("pointLight.diffuse"),
				1,
				&uniform.diffuse[0],
			)
			gl.Uniform3fv(
				shader.getLocation("pointLight.specular"),
				1,
				&uniform.specular[0],
			)
		case LIGHT_SOURCE_SPOT:
			switch light.ShadowType {
			case LIGHT_SHADOW_TYPE_NONE:
				shader = defaultShaders["deferred_spotLight_noshadow"]
				gl.UseProgram(shader.program)
			default:
				shadowMapShader = defaultShaders["shadow_map_spotLight"]
				lightMatrix, err := r.generateSpotLightShadowMap(shadowMapShader, light, meshes)
				if err != nil {
					return err
				}
				shader = defaultShaders["deferred_spotLight"]
				gl.UseProgram(shader.program)

				gl.UniformMatrix4fv(shader.getLocation("lightMatrix"), 1, false, &lightMatrix[0])
				gl.ActiveTexture(gl.TEXTURE3)
				gl.BindTexture(gl.TEXTURE_2D, r.sSpotMap)
				gl.Uniform1i(shader.getLocation("sSpotMap"), 3)
			}

			uniform := spotLightUniform{
				position:   light.Object().Position(),
				direction:  light.Object().Forward(),
				cosAngle:   float32(math.Cos(float64(light.Angle / 2))),
				lightRange: light.Range,
				diffuse:    light.Diffuse,
				specular:   light.Specular,
			}
			gl.Uniform3fv(
				shader.getLocation("spotLight.position"),
				1,
				&uniform.position[0],
			)
			gl.Uniform3fv(
				shader.getLocation("spotLight.direction"),
				1,
				&uniform.direction[0],
			)
			gl.Uniform1f(
				shader.getLocation("spotLight.cosAngle"),
				uniform.cosAngle,
			)
			gl.Uniform1f(
				shader.getLocation("spotLight.range"),
				uniform.lightRange,
			)
			gl.Uniform3fv(
				shader.getLocation("spotLight.diffuse"),
				1,
				&uniform.diffuse[0],
			)
			gl.Uniform3fv(
				shader.getLocation("spotLight.specular"),
				1,
				&uniform.specular[0],
			)
		}
		if shader == nil {
			return errors.New("no available shader")
		}

		cameraPos := camera.Object().Position()
		gl.UniformMatrix3fv(
			shader.getLocation("cameraPosition"),
			1,
			false,
			&cameraPos[0],
		)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, r.gPosition)
		gl.Uniform1i(shader.getLocation("gPosition"), 0)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, r.gNormal)
		gl.Uniform1i(shader.getLocation("gNormal"), 1)
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, r.gDiffuse)
		gl.Uniform1i(shader.getLocation("gDiffuse"), 2)

		gl.BindFramebuffer(gl.FRAMEBUFFER, targetFBO)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE)

		gl.BindVertexArray(r.quad)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		gl.Disable(gl.BLEND)

		gl.UseProgram(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	}

	return nil
}

func (r *deferredShading) generateDirLightShadowMap(shader *glShaderProgram, light *LightComponent, meshes []*MeshComponent) (*mgl32.Mat4, error) {
	gl.Viewport(0, 0, int32(r.renderer.dirLightShadowMapResolution), int32(r.renderer.dirLightShadowMapResolution))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.sDirBuffer)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

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
	gl.UniformMatrix4fv(shader.getLocation("lightMatrix"), 1, false, &lightMatrix[0])

	for _, mesh := range meshes {
		if !mesh.CastShadow {
			continue
		}
		rMesh, exists := r.renderer.meshes[mesh.Mesh]
		if !exists {
			err := r.renderer.helpLoad(mesh.Mesh)
			if err != nil {
				return nil, err
			}
			continue
		}
		model := mesh.Object().ModelMatrix()
		gl.UniformMatrix4fv(shader.getLocation("model"), 1, false, &model[0])
		if err := r.renderer.drawMesh(rMesh); err != nil {
			return nil, err
		}
	}

	scrWidth, scrHeight := r.renderer.context.ScreenSize()
	gl.Viewport(0, 0, int32(scrWidth), int32(scrHeight))
	return &lightMatrix, nil
}

func (r *deferredShading) generatePointLightShadowMap(shader *glShaderProgram, light *LightComponent, meshes []*MeshComponent) error {
	gl.Viewport(0, 0, int32(r.renderer.pointLightShadowMapResolution), int32(r.renderer.pointLightShadowMapResolution))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.sPointBuffer)
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(shader.program)

	lightPosition := light.Object().Position()
	lightRange := light.Range

	perspectiveMatrix := mgl32.Perspective(mgl32.DegToRad(90), 1, 0.1, light.Range)
	lightMatrices := []mgl32.Mat4{
		perspectiveMatrix.Mul4(mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{1, 0, 0}), mgl32.Vec3{0, -1, 0})),
		perspectiveMatrix.Mul4(mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{-1, 0, 0}), mgl32.Vec3{0, -1, 0})),
		perspectiveMatrix.Mul4(mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0, 1, 0}), mgl32.Vec3{0, 0, 1})),
		perspectiveMatrix.Mul4(mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0, -1, 0}), mgl32.Vec3{0, 0, -1})),
		perspectiveMatrix.Mul4(mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0, 0, 1}), mgl32.Vec3{0, -1, 0})),
		perspectiveMatrix.Mul4(mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0, 0, -1}), mgl32.Vec3{0, -1, 0})),
	}
	for i := 0; i < 6; i++ {
		gl.UniformMatrix4fv(
			shader.getLocation(fmt.Sprintf("lightMatrices[%d]", i)),
			1,
			false,
			&lightMatrices[i][0],
		)
	}

	gl.Uniform3fv(shader.getLocation("lightPosition"), 1, &lightPosition[0])
	gl.Uniform1f(shader.getLocation("lightRange"), lightRange)

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
		gl.UniformMatrix4fv(shader.getLocation("model"), 1, false, &model[0])
		if err := r.renderer.drawMesh(rMesh); err != nil {
			return err
		}
	}

	scrWidth, scrHeight := r.renderer.context.ScreenSize()
	gl.Viewport(0, 0, int32(scrWidth), int32(scrHeight))
	return nil
}

func (r *deferredShading) generateSpotLightShadowMap(shader *glShaderProgram, light *LightComponent, meshes []*MeshComponent) (*mgl32.Mat4, error) {
	gl.Viewport(0, 0, int32(r.renderer.spotLightShadowMapResolution), int32(r.renderer.spotLightShadowMapResolution))
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.sSpotBuffer)
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(shader.program)

	lightPosition := light.Object().Position()
	lightRange := light.Range

	lightMatrix := mgl32.Perspective(light.Angle, 1, 0.1, light.Range).Mul4(mgl32.LookAtV(
		lightPosition,
		lightPosition.Add(light.Object().Forward()),
		light.Object().Up(),
	))
	gl.UniformMatrix4fv(shader.getLocation("lightMatrix"), 1, false, &lightMatrix[0])

	gl.Uniform3fv(shader.getLocation("lightPosition"), 1, &lightPosition[0])
	gl.Uniform1f(shader.getLocation("lightRange"), lightRange)

	for _, mesh := range meshes {
		if !mesh.CastShadow {
			continue
		}
		rMesh, exists := r.renderer.meshes[mesh.Mesh]
		if !exists {
			err := r.renderer.helpLoad(mesh.Mesh)
			if err != nil {
				return nil, err
			}
			continue
		}
		model := mesh.Object().ModelMatrix()
		gl.UniformMatrix4fv(shader.getLocation("model"), 1, false, &model[0])
		if err := r.renderer.drawMesh(rMesh); err != nil {
			return nil, err
		}
	}

	scrWidth, scrHeight := r.renderer.context.ScreenSize()
	gl.Viewport(0, 0, int32(scrWidth), int32(scrHeight))
	return &lightMatrix, nil
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
		return defaultShaders["mesh_texture_deferred"], nil
	} else {
		return defaultShaders["mesh_color_deferred"], nil
	}
}

type deferredMeshUniform struct {
	projection, model, view mgl32.Mat4
}

func (r *deferredShading) applyShaderToMesh(shader *glShaderProgram, uniform deferredMeshUniform, mesh *MeshComponent) error {
	modelLoc := shader.getLocation("model")
	TImodelLoc := shader.getLocation("TImodel")
	viewLoc := shader.getLocation("view")
	projectionLoc := shader.getLocation("projection")
	material := r.renderer.materials[mesh.Material]
	tiModel := uniform.model.Mat3().Inv().Transpose()

	gl.UseProgram(shader.program)

	gl.UniformMatrix4fv(modelLoc, 1, false, &uniform.model[0])
	gl.UniformMatrix3fv(TImodelLoc, 1, false, &tiModel[0])
	gl.UniformMatrix4fv(modelLoc, 1, false, &uniform.model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &uniform.view[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &uniform.projection[0])

	if mesh.ReceiveShader {
		gl.Uniform1f(shader.getLocation("recvShadow"), 1)
	} else {
		gl.Uniform1f(shader.getLocation("recvShadow"), 0)
	}

	if material.diffuseMap != 0 {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, material.diffuseMap)
		gl.Uniform1i(shader.getLocation("diffuseMap"), 0)
	} else {
		gl.Uniform4fv(shader.getLocation("color"), 1, &material.DiffuseColor[0])
	}

	return nil
}
