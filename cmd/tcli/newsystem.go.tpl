-- define: newSystemTpl
package {{.Tags.Package}}

import (
    "github.com/gabstv/troupe/pkg/troupe"
    "github.com/hajimehoshi/ebiten"
)

// {{.Tags.Component}} is the data of a {{tolower .Tags.Component}} component.
type {{.Tags.Component}} struct {
    // public and private struct fields
}

// {{.Tags.Component}}Component will get the registered {{tolower .Tags.Component}} component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func {{.Tags.Component}}Component(w troupe.WorldDicter) *troupe.Component {
	c := w.Component("{{tolower .Tags.Package}}.{{.Tags.Component}}Component")
	if c == nil {
		var err error
		c, err = w.NewComponent(troupe.NewComponentInput{
			Name: "{{tolower .Tags.Package}}.{{.Tags.Component}}Component",
			ValidateDataFn: func(data interface{}) bool {
                if data == nil {
                    return false
                }
				_, ok := data.(*{{.Tags.Component}})
                return ok
			},
			DestructorFn: func(_ troupe.WorldDicter, entity troupe.Entity, data interface{}) {
				//TODO: fill
			},
		})
		if err != nil {
			panic(err)
		}
	}
	return c
}

// {{.Tags.Component}}System creates the {{tolower .Tags.Component}} system
func {{.Tags.Component}}System(w *troupe.World) *troupe.System {
	if sys := w.System("{{.Tags.Package}}.{{.Tags.Component}}System"); sys != nil {
		return sys
	}
	sys := w.NewSystem("{{.Tags.Package}}.{{.Tags.Component}}System", {{.Tags.Priority}}, {{.Tags.Component}}SystemExec, {{.Tags.Component}}Component(w))
	//sys.AddTag(troupe.WorldTagDraw)
	sys.AddTag(troupe.WorldTagUpdate)
	return sys
}

// {{.Tags.Component}}SystemExec is the main function of the {{.Tags.Component}}System
func {{.Tags.Component}}SystemExec(ctx troupe.Context, screen *ebiten.Image) {
	v := ctx.System().View()
	world := v.World()
	matches := v.Matches()
	{{tolower .Tags.Component}}comp := {{.Tags.Component}}Component(world)
	for _, m := range matches {
		_ = m.Components[{{tolower .Tags.Component}}comp].(*{{.Tags.Component}})
	}
}

// {{.Tags.Component}}CS ensures that all the required components and systems are added to the world.
func {{.Tags.Component}}CS(w *troupe.World) {
	{{.Tags.Component}}Component(w)
	{{.Tags.Component}}System(w)
	//TODO: add all additional required components and systems
} 

func init() {
	troupe.DefaultComp(func(e *troupe.Engine, w *troupe.World) {
		{{.Tags.Component}}Component(w)
		//TODO: add all additional required components
	})
	troupe.DefaultSys(func(e *troupe.Engine, w *troupe.World) {
		{{.Tags.Component}}System(w)
		//TODO: add all additional required systems
	})
}

-- end