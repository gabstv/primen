package imgui

import (
	"errors"
	"sort"

	"github.com/inkyblackness/imgui-go/v2"
)

type uiMemory struct {
	variables []*UIVariable
}

func newUIMemory() *uiMemory {
	return &uiMemory{
		variables: make([]*UIVariable, 0, 64),
	}
}

func (m *uiMemory) find(name string) int {
	i := sort.Search(len(m.variables), func(i int) bool { return m.variables[i].Name >= name })
	if i < len(m.variables) && m.variables[i].Name == name {
		return i
	}
	return -1
}

func (m *uiMemory) Add(name string, value interface{}) error {
	if m.find(name) != -1 {
		return errors.New("variable already exists")
	}
	m.variables = append(m.variables, &UIVariable{
		Name:  name,
		Value: value,
	})
	sort.Sort(UIVariables(m.variables))
	return nil
}

func (m *uiMemory) MustBool(name string, defaultv bool) bool {
	x := m.find(name)
	if x != -1 {
		return m.variables[x].Value.(bool)
	}
	m.Add(name, defaultv)
	return m.MustBool(name, defaultv)
}

func (m *uiMemory) MustInt(name string, defaultv int) int {
	x := m.find(name)
	if x != -1 {
		return m.variables[x].Value.(int)
	}
	m.Add(name, defaultv)
	return m.MustInt(name, defaultv)
}

func (m *uiMemory) Set(name string, v interface{}) {
	x := m.find(name)
	if x != -1 {
		m.variables[x].Value = v
	}
	m.Add(name, v)
}

func (m *uiMemory) ImGuiVec4(name string) (imgui.Vec4, bool) {
	x := m.find(name)
	if x != -1 {
		return m.variables[x].Value.(imgui.Vec4), true
	}
	return imgui.Vec4{}, false
}

func (m *uiMemory) ImGuiVec2(name string) (imgui.Vec2, bool) {
	x := m.find(name)
	if x != -1 {
		return m.variables[x].Value.(imgui.Vec2), true
	}
	return imgui.Vec2{}, false
}

func (m *uiMemory) ImGuiFloat(name string) (float32, bool) {
	x := m.find(name)
	if x != -1 {
		return m.variables[x].Value.(float32), true
	}
	return 0, false
}

type UIVariable struct {
	Name  string
	Value interface{}
}

type UIVariables []*UIVariable

func (a UIVariables) Len() int           { return len(a) }
func (a UIVariables) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UIVariables) Less(i, j int) bool { return a[i].Name < a[j].Name }
