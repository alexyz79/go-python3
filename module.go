/*
Unless explicitly stated otherwise all files in this repository are licensed
under the $license_for_repo License.
This product includes software developed at Datadog (https://www.datadoghq.com/).
Copyright 2018 Datadog, Inc.
*/

package python3

/*
#include "Python.h"
#include "macro.h"
*/
import "C"
import (
	"unsafe"
)

// PyCFunction ...
// type PyCFunction func(*PyObject, *PyObject) *PyObject
type PyCFunction unsafe.Pointer

// PyMethodDef ...
type PyMethodDef struct {
	Name  string // name of the method
	Meth  PyCFunction
	Flags MethodDefFlags
	Doc   string
}

// PyModuleDef ...
type PyModuleDef struct {
	Name    string // name of the method
	Doc     string
	Size    int
	Methods []PyMethodDef
}

// MethodDefFlags ...
type MethodDefFlags int32

const (
	// MethVarArgs ...
	MethVarArgs MethodDefFlags = C.METH_VARARGS
	// MethKeyWords ...
	MethKeyWords = C.METH_KEYWORDS
	// MethNoArgs ...
	MethNoArgs = C.METH_NOARGS
	// MethO ...
	MethO = C.METH_O
	// MethClass ...
	MethClass = C.METH_CLASS
	// MethStatic ...
	MethStatic = C.METH_STATIC
	// MethCoexist ...
	MethCoexist = C.METH_COEXIST
)

//Module : https://docs.python.org/3/c-api/module.html#c.PyModule_Type
var Module = togo((*C.PyObject)(unsafe.Pointer(&C.PyModule_Type)))

//PyModule_Check : https://docs.python.org/3/c-api/module.html#c.PyModule_Check
func PyModule_Check(p *PyObject) bool {
	return C._go_PyModule_Check(toc(p)) != 0
}

//PyModule_CheckExact : https://docs.python.org/3/c-api/module.html#c.PyModule_CheckExact
func PyModule_CheckExact(p *PyObject) bool {
	return C._go_PyModule_CheckExact(toc(p)) != 0
}

//PyModule_NewObject : https://docs.python.org/3/c-api/module.html#c.PyModule_NewObject
func PyModule_NewObject(name *PyObject) *PyObject {
	return togo(C.PyModule_NewObject(toc(name)))
}

//PyModule_New : https://docs.python.org/3/c-api/module.html#c.PyModule_New
func PyModule_New(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return togo(C.PyModule_New(cname))
}

//PyModule_GetDict : https://docs.python.org/3/c-api/module.html#c.PyModule_GetDict
func PyModule_GetDict(module *PyObject) *PyObject {
	return togo(C.PyModule_GetDict(toc(module)))
}

//PyModule_GetNameObject : https://docs.python.org/3/c-api/module.html#c.PyModule_GetNameObject
func PyModule_GetNameObject(module *PyObject) *PyObject {
	return togo(C.PyModule_GetNameObject(toc(module)))
}

//PyModule_GetName : https://docs.python.org/3/c-api/module.html#c.PyModule_GetName
func PyModule_GetName(module *PyObject) string {
	cname := C.PyModule_GetName(toc(module))
	return C.GoString(cname)
}

//PyModule_GetState : https://docs.python.org/3/c-api/module.html#c.PyModule_GetState
func PyModule_GetState(module *PyObject) unsafe.Pointer {
	return unsafe.Pointer(C.PyModule_GetNameObject(toc(module)))
}

//PyModule_GetFilenameObject : https://docs.python.org/3/c-api/module.html#c.PyModule_GetFilenameObject
func PyModule_GetFilenameObject(module *PyObject) *PyObject {
	return togo(C.PyModule_GetFilenameObject(toc(module)))
}

// PyModule_Create https://docs.python.org/3/c-api/module.html#c.PyModule_Create
func PyModule_Create(module PyModuleDef) *PyObject {

	// ml_meth:  C._go_PyCFunction((C.intptr_t)(reflect.ValueOf(methodDef.Meth).Pointer())),

	n := C.size_t(len(module.Methods) + 1)
	pyMethodDefs := C._go_malloc_PyMethodDefArray(n)
	for i, methodDef := range module.Methods {
		if i < len(module.Methods) {
			C._go_set_PyMethodDef(pyMethodDefs, C.int(i), &C.PyMethodDef{
				ml_name:  C.CString(methodDef.Name),
				ml_meth:  C.PyCFunction(methodDef.Meth),
				ml_flags: C.int(methodDef.Flags),
				ml_doc:   C.CString(methodDef.Doc),
			})
		}
	}

	C._go_set_PyMethodDef(pyMethodDefs, C.int(len(module.Methods)), &C.PyMethodDef{
		ml_name:  nil,
		ml_meth:  nil,
		ml_flags: 0,
		ml_doc:   nil,
	})

	pyObjectHEADINIT := C.PyObject{
		ob_refcnt: 1,
		ob_type:   nil,
	}

	pyModuleDefHEADINIT := C.PyModuleDef_Base{
		ob_base: pyObjectHEADINIT,
		m_init:  nil,
		m_index: 0,
		m_copy:  nil,
	}

	pyModuleDef := C.PyModuleDef{
		m_base:     pyModuleDefHEADINIT,
		m_name:     C.CString(module.Name),
		m_doc:      C.CString(module.Doc),
		m_size:     C.Py_ssize_t(module.Size),
		m_methods:  pyMethodDefs,
		m_slots:    nil,
		m_traverse: nil,
		m_clear:    nil,
		m_free:     nil,
	}

	return togo(C._go_PyModule_Create((*C.PyModuleDef)(&pyModuleDef)))
}
