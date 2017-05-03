package wengine

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strconv"
	"strings"
)

type Asset interface {
	load() error
	Loaded() bool
}

type MeshAsset struct {
	Path   string
	Buffer []byte

	vertices []mgl32.Vec3
	uvs      []mgl32.Vec2
	norms    []mgl32.Vec3
}

func DefaultCubeAsset() *MeshAsset {
	return &MeshAsset{Buffer: []byte(defaultCubeObj)}
}

func DefaultPlaneAsset() *MeshAsset {
	return &MeshAsset{Buffer: []byte(defaultPlaneObj)}
}

func (mesh *MeshAsset) InternalData() (vertices []mgl32.Vec3, uvs []mgl32.Vec2, norms []mgl32.Vec3) {
	vertices = mesh.vertices
	uvs = mesh.uvs
	norms = mesh.norms
	return
}

func (mesh *MeshAsset) Loaded() bool {
	return mesh.vertices != nil && mesh.norms != nil && mesh.uvs != nil
}

func (mesh *MeshAsset) load() error {
	var objReader io.Reader
	if mesh.Buffer != nil {
		objReader = bytes.NewReader(mesh.Buffer)
	} else if mesh.Path != "" {
		file, err := os.Open(mesh.Path)
		defer file.Close()
		if err != nil {
			return err
		}
		objReader = file
	}
	vertices := []mgl32.Vec3{}
	uvs := []mgl32.Vec2{}
	norms := []mgl32.Vec3{}
	faces := [][3][3]int{}
	scanner := bufio.NewScanner(objReader)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), " ")
		if len(cols) == 0 {
			continue
		}
		switch cols[0] {
		case "v":
			if len(cols) < 4 {
				return errors.New("corrupted obj file")
			}
			vec := mgl32.Vec3{}
			for i := 1; i <= 3; i++ {
				f, err := strconv.ParseFloat(cols[i], 32)
				if err != nil {
					return errors.New("corrupted value in obj file")
				}
				vec[i-1] = float32(f)
			}
			vertices = append(vertices, vec)
		case "vn":
			if len(cols) < 4 {
				return errors.New("corrupted obj file")
			}
			vec := mgl32.Vec3{}
			for i := 1; i <= 3; i++ {
				f, err := strconv.ParseFloat(cols[i], 32)
				if err != nil {
					return errors.New("corrupted value in obj file")
				}
				vec[i-1] = float32(f)
			}
			norms = append(norms, vec)
		case "vt":
			if len(cols) < 3 {
				return errors.New("corrupted obj file")
			}
			vec := mgl32.Vec2{}
			for i := 1; i <= 2; i++ {
				f, err := strconv.ParseFloat(cols[i], 32)
				if err != nil {
					return errors.New("corrupted value in obj file")
				}
				vec[i-1] = float32(f)
			}
			uvs = append(uvs, vec)
		case "f":
			if len(cols) < 4 {
				return errors.New("corrupted obj file")
			}
			face := [3][3]int{}
			for i := 1; i <= 3; i++ {
				faceVals := strings.Split(cols[i], "/")
				if len(faceVals) < 3 {
					return errors.New("corrupted obj file")
				}
				vec := [3]int{}
				for j := 0; j < 3; j++ {
					f, err := strconv.Atoi(faceVals[j])
					if err != nil {
						if j != 0 {
							f = 1
						} else {
							return errors.New("corrupted value in obj file")
						}
					}
					vec[j] = f
				}
				face[i-1] = vec
			}
			faces = append(faces, face)
		}
	}
	if len(uvs) == 0 {
		uvs = make([]mgl32.Vec2, 1)
	}
	if len(norms) == 0 {
		norms = make([]mgl32.Vec3, 1)
	}
	for _, face := range faces {
		for _, v := range face {
			mesh.vertices = append(mesh.vertices, vertices[v[0]-1])
			mesh.uvs = append(mesh.uvs, uvs[v[1]-1])
			mesh.norms = append(mesh.norms, norms[v[2]-1])
		}
	}
	return nil
}

// -----------------------------------------------------------

type MeshMaterialAsset struct {
	DiffuseColor     mgl32.Vec4
	DiffuseMapPath   string
	DiffuseMapBuffer []byte

	colorImage *image.RGBA
}

func (m *MeshMaterialAsset) InternalData() (colorImage *image.RGBA) {
	colorImage = m.colorImage
	return
}

func (m *MeshMaterialAsset) Loaded() bool {
	return !(m.DiffuseMapPath != "" || m.DiffuseMapBuffer != nil) ||
		((m.DiffuseMapPath != "" || m.DiffuseMapBuffer != nil) && m.colorImage != nil)
}

func (m *MeshMaterialAsset) load() error {
	var imgReader io.Reader
	if m.DiffuseMapBuffer != nil {
		imgReader = bytes.NewReader(m.DiffuseMapBuffer)
	} else if m.DiffuseMapPath != "" {
		file, err := os.Open(m.DiffuseMapPath)
		defer file.Close()
		if err != nil {
			return err
		}
		imgReader = file
	}
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return err
	}
	m.colorImage = image.NewRGBA(img.Bounds())
	draw.Draw(m.colorImage, m.colorImage.Bounds(), img, image.Point{0, 0}, draw.Src)
	return nil
}

// -----------------------------------------------------------

type SpriteMaterialAsset struct {
	TexturePath   string
	TextureBuffer []byte

	textureImage *image.RGBA
}

func (m *SpriteMaterialAsset) InternalData() (colorImage *image.RGBA) {
	colorImage = m.textureImage
	return
}

func (m *SpriteMaterialAsset) Loaded() bool {
	return m.textureImage != nil
}

func (m *SpriteMaterialAsset) load() error {
	var imgReader io.Reader
	if m.TextureBuffer != nil {
		imgReader = bytes.NewReader(m.TextureBuffer)
	} else if m.TexturePath != "" {
		file, err := os.Open(m.TexturePath)
		defer file.Close()
		if err != nil {
			return err
		}
		imgReader = file
	}
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return err
	}
	m.textureImage = image.NewRGBA(img.Bounds())
	draw.Draw(m.textureImage, m.textureImage.Bounds(), img, image.Point{0, 0}, draw.Src)
	return nil
}

// -----------------------------------------------------------

type ShaderAsset struct {
	VertexSource, GeometrySource, FragmentSource string
}

func (s *ShaderAsset) Loaded() bool {
	return s.VertexSource != "" && s.FragmentSource != ""
}

func (s *ShaderAsset) load() error {
	return nil
}

// -----------------------------------------------------------

const defaultCubeObj = `
# Blender v2.78 (sub 0) OBJ File: ''
# www.blender.org
mtllib cube.mtl
o Cube
v 1.000000 -1.000000 -1.000000
v 1.000000 -1.000000 1.000000
v -1.000000 -1.000000 1.000000
v -1.000000 -1.000000 -1.000000
v 1.000000 1.000000 -0.999999
v 0.999999 1.000000 1.000001
v -1.000000 1.000000 1.000000
v -1.000000 1.000000 -1.000000
vt 0.0001 0.2500
vt 0.2500 0.5000
vt 0.0001 0.5000
vt 0.5000 0.5000
vt 0.7500 0.2500
vt 0.7500 0.5000
vt 0.9999 0.2500
vt 0.9999 0.5000
vt 0.5000 0.0001
vt 0.7500 0.0001
vt 0.2500 0.2500
vt 0.7500 0.7500
vt 0.5000 0.2500
vt 0.5000 0.7500
vn 0.0000 -1.0000 0.0000
vn 0.0000 1.0000 0.0000
vn 1.0000 -0.0000 0.0000
vn 0.0000 -0.0000 1.0000
vn -1.0000 -0.0000 -0.0000
vn 0.0000 0.0000 -1.0000
s off
f 2/1/1 4/2/1 1/3/1
f 8/4/2 6/5/2 5/6/2
f 5/6/3 2/7/3 1/8/3
f 6/5/4 3/9/4 2/10/4
f 3/11/5 8/4/5 4/2/5
f 1/12/6 8/4/6 5/6/6
f 2/1/1 3/11/1 4/2/1
f 8/4/2 7/13/2 6/5/2
f 5/6/3 6/5/3 2/7/3
f 6/5/4 7/13/4 3/9/4
f 3/11/5 7/13/5 8/4/5
f 1/12/6 4/14/6 8/4/6
`
const defaultPlaneObj = `
# Blender v2.78 (sub 0) OBJ File: ''
# www.blender.org
o Plane
v -1.000000 0.000000 1.000000
v 1.000000 0.000000 1.000000
v -1.000000 0.000000 -1.000000
v 1.000000 0.000000 -1.000000
vt 0.9999 0.0001
vt 0.0001 0.9999
vt 0.0001 0.0001
vt 0.9999 0.9999
vn 0.0000 1.0000 0.0000
s off
f 2/1/1 3/2/1 1/3/1
f 2/1/1 4/4/1 3/2/1
`

const defaultSpriteObj = `
# Blender v2.78 (sub 0) OBJ File: ''
# www.blender.org
o Plane
v -1.000000 -1.000000 -0.000000
v 1.000000 -1.000000 -0.000000
v -1.000000 1.000000 0.000000
v 1.000000 1.000000 0.000000
vt 0.9999 0.0001
vt 0.0001 0.9999
vt 0.0001 0.0001
vt 0.9999 0.9999
vn 0.0000 -0.0000 1.0000
s off
f 2/1/1 3/2/1 1/3/1
f 2/1/1 4/4/1 3/2/1
`
