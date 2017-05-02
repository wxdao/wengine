package opengl

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	. "github.com/wxdao/wengine"
	"image"
	"sort"
	"strings"
)

func init() {
	RegisterRenderer("opengl", newRenderer())
}

type renderer struct {
	context *Context

	meshes    map[string]*glMesh
	materials map[string]*glMaterial
	programs  map[string]*glShaderProgram

	dirLightShadowMapResolution   int
	pointLightShadowMapResolution int
	spotLightShadowMapResolution  int

	pc renderPath

	assetsToInstall []string

	lastScene *Scene

	versionStr string
}

func newRenderer() *renderer {
	r := &renderer{
		meshes:                        map[string]*glMesh{},
		materials:                     map[string]*glMaterial{},
		programs:                      map[string]*glShaderProgram{},
		dirLightShadowMapResolution:   3072,
		pointLightShadowMapResolution: 512,
		spotLightShadowMapResolution:  1024,
		assetsToInstall:               []string{},
	}
	//r.pc = &forwardShading{renderer: r}
	r.pc = &deferredShading{renderer: r}
	return r
}

func (r *renderer) Init(context *Context) error {
	r.context = context

	if err := gl.Init(); err != nil {
		return errors.New(fmt.Sprint("unable to init opengl:", err))
	}
	r.versionStr = gl.GoStr(gl.GetString(gl.VERSION))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)

	if err := r.loadDefaultShaders(); err != nil {
		return err
	}

	if err := r.pc.init(); err != nil {
		return err
	}

	return nil
}

func (r *renderer) loadDefaultShaders() error {
	for name, shader := range defaultShaders {
		if err := shader.install(); err != nil {
			return err
		}
		r.programs["___"+name] = shader
		println("installed shader: ___" + name)
	}

	return nil
}

func (r *renderer) Version() string {
	return r.versionStr
}

func (r *renderer) Render(scene *Scene) error {
	r.installAll()
	// find all cameras
	cameras := []*CameraComponent{}
	meshes := []*MeshComponent{}
	lights := []*LightComponent{}
	for _, obj := range scene.Objects() {
		if !obj.Enabled() {
			continue
		}
		for _, compo := range obj.Components() {
			switch compo.Type() {
			case COMPO_CAMERA:
				cameras = append(cameras, compo.(*CameraComponent))
			case COMPO_MESH:
				meshes = append(meshes, compo.(*MeshComponent))
			case COMPO_LIGHT:
				lights = append(lights, compo.(*LightComponent))
			}
		}
	}
	// sort by depth, decreasing
	sort.Slice(cameras, func(i, j int) bool {
		return !(cameras[i].Depth < cameras[j].Depth)
	})
	// hand over to renderPath
	if err := r.pc.render(0, lights, meshes, scene, cameras); err != nil {
		return err
	}
	return nil
}

func (r *renderer) NotifyInstall(assets []string) error {
	r.assetsToInstall = append(r.assetsToInstall, assets...)
	return nil
}

func (r *renderer) drawMesh(mesh *glMesh) error {
	if !mesh.installed() {
		return errors.New("uninstalled")
	}
	gl.BindVertexArray(mesh.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(mesh.num))
	gl.BindVertexArray(0)
	return nil
}

func (r *renderer) helpLoad(asset string) error {
	if _, exists := r.context.Assets()[asset]; exists {
		err := r.context.LoadAssets([]string{asset})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) installAll() error {
	for _, name := range r.assetsToInstall {
		asset, exists := r.context.Assets()[name]
		if !exists {
			return errors.New("no internal asset")
		}
		switch a := asset.(type) {
		case *MeshAsset:
			if _, exists := r.meshes[name]; exists {
				continue
			}
			r.meshes[name] = &glMesh{MeshAsset: a}
			if err := r.meshes[name].install(); err != nil {
				return err
			}
			println("installed mesh: " + name)
		case *MaterialAsset:
			if _, exists := r.materials[name]; exists {
				continue
			}
			r.materials[name] = &glMaterial{MaterialAsset: a}
			if err := r.materials[name].install(); err != nil {
				return err
			}
			println("installed material: " + name)
		}
	}
	r.assetsToInstall = []string{}
	return nil
}

// -----------------------------------------------------------

type glMesh struct {
	*MeshAsset

	num int
	vao uint32
}

func (m *glMesh) installed() bool {
	if m.vao == 0 {
		return false
	}
	return true
}

func (m *glMesh) install() error {
	if m.installed() {
		return nil
	}

	vertices, uvs, norms := m.InternalData()
	m.num = len(vertices)

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	vbo := [3]uint32{}
	gl.GenBuffers(3, &vbo[0])

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
	// norms
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(norms)*3*4, gl.Ptr(norms), gl.STATIC_DRAW)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

// -----------------------------------------------------------

type glMaterial struct {
	*MaterialAsset

	diffuseMap uint32
}

func (m *glMaterial) installed() bool {
	if (m.DiffuseMapPath != "" || m.DiffuseMapBuffer != nil) && m.diffuseMap == 0 {
		return false
	}
	return true
}

func (m *glMaterial) install() error {
	diffuseImage := m.InternalData()
	if m.diffuseMap == 0 && diffuseImage != nil {
		// invert y
		flipped := image.NewRGBA(diffuseImage.Bounds())
		xLen := diffuseImage.Rect.Size().X
		yLen := diffuseImage.Rect.Size().Y
		for y := 0; y < yLen; y++ {
			for x := 0; x < xLen; x++ {
				end := yLen - y
				o := diffuseImage.At(x, end)
				flipped.Set(x, y, o)
			}
		}

		gl.GenTextures(1, &m.diffuseMap)
		gl.BindTexture(gl.TEXTURE_2D, m.diffuseMap)
		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(xLen),
			int32(yLen),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(flipped.Pix),
		)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
	return nil
}

// -----------------------------------------------------------

type glShaderProgram struct {
	vertexSource, geometrySource, fragmentSource string

	modelLoc, viewLoc, projectionLoc uint32

	program uint32

	locations map[string]int32
}

func (p *glShaderProgram) installed() bool {
	if p.program == 0 {
		return false
	}
	return true
}

func (p *glShaderProgram) install() error {
	if p.installed() {
		return nil
	}

	program := gl.CreateProgram()

	vertexShader, err := compileShader(p.vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return errors.New("vertex: " + err.Error())
	}
	defer gl.DeleteShader(vertexShader)
	gl.AttachShader(program, vertexShader)

	if p.geometrySource != "" {
		geometryShader, err := compileShader(p.geometrySource, gl.GEOMETRY_SHADER)
		if err != nil {
			return errors.New("geometry: " + err.Error())
		}
		gl.AttachShader(program, geometryShader)
		defer gl.DeleteShader(geometryShader)
	}

	fragmentShader, err := compileShader(p.fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return errors.New("fragment: " + err.Error())
	}
	gl.AttachShader(program, fragmentShader)
	defer gl.DeleteShader(fragmentShader)

	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return errors.New("program: " + log)
	}
	p.program = program
	p.locations = map[string]int32{}
	return nil
}

func (p *glShaderProgram) getLocation(name string) int32 {
	location, exists := p.locations[name]
	if !exists {
		location = gl.GetUniformLocation(p.program, gl.Str(name+"\x00"))
		p.locations[name] = location
	}
	return location
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	cSource, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(shader, 1, cSource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, errors.New(log)
	}
	return shader, nil
}
