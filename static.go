// Code generated by "esc -o static.go -ignore \.map$ app"; DO NOT EDIT.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/app/eventemitter3.js": {
		local:   "app/eventemitter3.js",
		size:    3518,
		modtime: 1517990254,
		compressed: `
H4sIAAAAAAAC/5RXW2/bOBb+K7IeVBJmDGnb3QcphFG02adus8C0T4IQKPRRwhmZ1PCSNLA1v31AiZQl
N0kzLxYv37l+h4f0qrGCGS4FAnzgDYrl7e/ATEypeepANhH86KQyOkliK3bQcAG7eBU293JnW8DjZ+Oh
FBAuoNUQOX1B/0njqCVJxu+m3u/wOERlRWAUPaDnzD1ysZOP2/GTP4e4a+Vt3W7Hz7MIDW2zdT+5ueca
b64eQJirPTcG1Hvnet+jKSf4oMBYJaKwEgEyRBCFD9OKRJo0Q+5WotTVODLD6KFWEaPPpEDBn5YrSBI/
KJxMkyQMe3sMabJKsVvnYY37NafVUgGP0ZVSUqH4Uy2ENFHDxc4zEr2L13odv4txYe6VfIzshskd0Ph/
15+/f7m6+Xr97ea/19+/fo6J7Z2+mjrf6cFTmB/6vnAxlGm1YXXbojqwS+YF40QFHYBZVUJVeFclEscj
4J7U5CQIZExd70HOYtjsG6mQ08bfki6iaVroS7VpQdyZ+0Kv11gi5XI+edCjQ5aXJ2edcXyIrYZIG8WZ
iYdEKno9FPymU9JIZ29zX+vrR/F/JTtQ5olIGv8VFxPdHOFDP8008lERPh6fyfmp4AQeKXCEfXvqwJP2
7R6ilmsDAlS0t9pEtxDVU6HFI8964Pk8Cle5m0ZQIMOISWHgh6FmnErBgIrjcZX1SBB1PALhmDRUbuXa
5CZkCDY34Gpfl021nU82jVjMaTmfEV3lC3Rn9T3SOEcLGU2m6SdphVmvMYFT3hoXCj6klF5cLIEn00Po
PN9BCwZm7prqpIchn42FDJkvDVpp2nuemYLaQJIgfqKcLvaQsG2LCRpU4c3NzYC7uTkekaSrDGPCZtUy
GPla70HTWdtw3A1k0bJyhzillP7kVDjYonDVbyIuIligsBoPn8tVkogx1XJrNrrlDFCGczMVvA/hDsys
eH972t/KVm+FqxFWG/QaCgHGuegX0YUKnQXnj70Z6glyIGLhc2mGgFfCRzfGLzbNtODGVREOvKIp4VT4
o0xGBj8qVT8hjgt1yQu1XmNdqoqKUrniDBHr510d+f4n7gYats6zbZYHZ/J0aQH2fK44HHyiRxNNMDHE
vzDSVD74VTYca0YsqekZhLS0Vnd276ahs/EG1S51B/3IDbtH9XC6k2QQVbCXD/DFx42AOCh5kHwXpe6m
IC0+sFpDlOU+RAcI/dx3DUxWaTGg/vUKisxw71/FEXFCfvgFkqgT9t+/xLouFtD/eQOaaIcfbhZLM8Jm
ldVeZLiwl21h12vMSnuRVafcl7YqBrV117VPM70M98PbxDHYkQdaB5ZGE2lhLx8GjYGs0lav8+UAL1Hm
N31wbhIYK24V1H8E0l7EkSXy/WtIIhbYD69jiQroHTS1bU3uKp5hl4fuuVR3l23RDanulqnuqiJY8tme
m2K494+F1dlRlIKeX4qhKyCXa//YWGX4XI7B2yTTM8kleT+1gbEF8JdbAA8tIHLL435YaUbLHA8Xl7/4
z6SdhLv0sftZUWqOR5Ukq/G6Px5FkkwPgRV11/+kdHxQh37LaEosLStS06nLsMu6YOs1Rrpk1VK9WzhZ
cLOZEZwk1t//JatwYb3C7ZnvNKOUTpu2TKvc5pN//Swvz+X8Y9t+efkamp4zW3Tq8WcdPkm8MYNxjt72
Xhi5OKuepqEv18QCWu92U60sC3CAQcN/wI5KwhZ/QCjzvcAxkCTITP+pGO7Joa/cDymzCqMM97j4OwAA
//8y9qzVvg0AAA==
`,
	},

	"/app/howler.js": {
		local:   "app/howler.js",
		size:    30738,
		modtime: 1521559205,
		compressed: `
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
ZnJsWZgw0r22bQLy/xsAAP//L4Ds4hJ4AAA=
`,
	},

	"/app/room.js": {
		local:   "app/room.js",
		size:    6766,
		modtime: 1521665938,
		compressed: `
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
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/app": {
		isDir: true,
		local: "app",
	},
}
