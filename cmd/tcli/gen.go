// Code generated by sql2var <https://github.com/gabstv/sql2var>. DO NOT EDIT.
// source: main.go

package main

const newSystemTpl = "package {{.Tags.Package}}\n\nimport (\n    \"github.com/gabstv/primen\"\n    \"github.com/gabstv/ebiten\"\n)\n\n// {{.Tags.Component}} is the data of a {{tolower .Tags.Component}} component.\ntype {{.Tags.Component}} struct {\n    // public and private struct fields\n}\n\n// {{.Tags.Component}}Component will get the registered {{tolower .Tags.Component}} component of the world.\n// If a component is not present, it will create a new component\n// using world.NewComponent\nfunc {{.Tags.Component}}Component(w primen.WorldDicter) *primen.Component {\n\tc := w.Component(\"{{tolower .Tags.Package}}.{{.Tags.Component}}Component\")\n\tif c == nil {\n\t\tvar err error\n\t\tc, err = w.NewComponent(primen.NewComponentInput{\n\t\t\tName: \"{{tolower .Tags.Package}}.{{.Tags.Component}}Component\",\n\t\t\tValidateDataFn: func(data interface{}) bool {\n                if data == nil {\n                    return false\n                }\n\t\t\t\t_, ok := data.(*{{.Tags.Component}})\n                return ok\n\t\t\t},\n\t\t\tDestructorFn: func(_ primen.WorldDicter, entity primen.Entity, data interface{}) {\n\t\t\t\t//TODO: fill\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\tpanic(err)\n\t\t}\n\t}\n\treturn c\n}\n\n// {{.Tags.Component}}System creates the {{tolower .Tags.Component}} system\nfunc {{.Tags.Component}}System(w *primen.World) *primen.System {\n\tif sys := w.System(\"{{.Tags.Package}}.{{.Tags.Component}}System\"); sys != nil {\n\t\treturn sys\n\t}\n\tsys := w.NewSystem(\"{{.Tags.Package}}.{{.Tags.Component}}System\", {{.Tags.Priority}}, {{.Tags.Component}}SystemExec, {{.Tags.Component}}Component(w))\n\t//sys.AddTag(primen.WorldTagDraw)\n\tsys.AddTag(primen.WorldTagUpdate)\n\treturn sys\n}\n\n// {{.Tags.Component}}SystemExec is the main function of the {{.Tags.Component}}System\nfunc {{.Tags.Component}}SystemExec(ctx primen.Context, screen *ebiten.Image) {\n\tv := ctx.System().View()\n\tworld := v.World()\n\tmatches := v.Matches()\n\t{{tolower .Tags.Component}}comp := {{.Tags.Component}}Component(world)\n\tfor _, m := range matches {\n\t\t_ = m.Components[{{tolower .Tags.Component}}comp].(*{{.Tags.Component}})\n\t}\n}\n\n// {{.Tags.Component}}CS ensures that all the required components and systems are added to the world.\nfunc {{.Tags.Component}}CS(w *primen.World) {\n\t{{.Tags.Component}}Component(w)\n\t{{.Tags.Component}}System(w)\n\t//TODO: add all additional required components and systems\n} \n\nfunc init() {\n\tprimen.DefaultComp(func(e *primen.Engine, w *primen.World) {\n\t\t{{.Tags.Component}}Component(w)\n\t\t//TODO: add all additional required components\n\t})\n\tprimen.DefaultSys(func(e *primen.Engine, w *primen.World) {\n\t\t{{.Tags.Component}}System(w)\n\t\t//TODO: add all additional required systems\n\t})\n}\n\n"
