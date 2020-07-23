// Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.

package components

import (
    
    "sort"

    "github.com/gabstv/ecs/v2"
    
)









const uuidCameraSystem = "E19B710D-139B-47BD-AF0C-340414BC7226"

type viewCameraSystem struct {
    entities []VICameraSystem
    world ecs.BaseWorld
    
}

type VICameraSystem struct {
    Entity ecs.Entity
    
    Transform *Transform 
    
    Camera *Camera 
    
}

type sortedVICameraSystems []VICameraSystem
func (a sortedVICameraSystems) Len() int           { return len(a) }
func (a sortedVICameraSystems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedVICameraSystems) Less(i, j int) bool { return a[i].Entity < a[j].Entity }

func newviewCameraSystem(w ecs.BaseWorld) *viewCameraSystem {
    return &viewCameraSystem{
        entities: make([]VICameraSystem, 0),
        world: w,
    }
}

func (v *viewCameraSystem) Matches() []VICameraSystem {
    
    return v.entities
    
}

func (v *viewCameraSystem) indexof(e ecs.Entity) int {
    i := sort.Search(len(v.entities), func(i int) bool { return v.entities[i].Entity >= e })
    if i < len(v.entities) && v.entities[i].Entity == e {
        return i
    }
    return -1
}

// Fetch a specific entity
func (v *viewCameraSystem) Fetch(e ecs.Entity) (data VICameraSystem, ok bool) {
    
    i := v.indexof(e)
    if i == -1 {
        return VICameraSystem{}, false
    }
    return v.entities[i], true
}

func (v *viewCameraSystem) Add(e ecs.Entity) bool {
    
    
    // MUST NOT add an Entity twice:
    if i := v.indexof(e); i > -1 {
        return false
    }
    v.entities = append(v.entities, VICameraSystem{
        Entity: e,
        Transform: GetTransformComponent(v.world).Data(e),
Camera: GetCameraComponent(v.world).Data(e),

    })
    if len(v.entities) > 1 {
        if v.entities[len(v.entities)-1].Entity < v.entities[len(v.entities)-2].Entity {
            sort.Sort(sortedVICameraSystems(v.entities))
        }
    }
    return true
}

func (v *viewCameraSystem) Remove(e ecs.Entity) bool {
    
    
    if i := v.indexof(e); i != -1 {

        v.entities = append(v.entities[:i], v.entities[i+1:]...)
        return true
    }
    return false
}

func (v *viewCameraSystem) clearpointers() {
    
    
    for i := range v.entities {
        e := v.entities[i].Entity
        
        v.entities[i].Transform = nil
        
        v.entities[i].Camera = nil
        
        _ = e
    }
}

func (v *viewCameraSystem) rescan() {
    
    
    for i := range v.entities {
        e := v.entities[i].Entity
        
        v.entities[i].Transform = GetTransformComponent(v.world).Data(e)
        
        v.entities[i].Camera = GetCameraComponent(v.world).Data(e)
        
        _ = e
        
    }
}

// CameraSystem implements ecs.BaseSystem
type CameraSystem struct {
    initialized bool
    world       ecs.BaseWorld
    view        *viewCameraSystem
    enabled     bool
    
}

// GetCameraSystem returns the instance of the system in a World
func GetCameraSystem(w ecs.BaseWorld) *CameraSystem {
    return w.S(uuidCameraSystem).(*CameraSystem)
}

// Enable system
func (s *CameraSystem) Enable() {
    s.enabled = true
}

// Disable system
func (s *CameraSystem) Disable() {
    s.enabled = false
}

// Enabled checks if enabled
func (s *CameraSystem) Enabled() bool {
    return s.enabled
}

// UUID implements ecs.BaseSystem
func (CameraSystem) UUID() string {
    return "E19B710D-139B-47BD-AF0C-340414BC7226"
}

func (CameraSystem) Name() string {
    return "CameraSystem"
}

// ensure matchfn
var _ ecs.MatchFn = matchCameraSystem

// ensure resizematchfn
var _ ecs.MatchFn = resizematchCameraSystem

func (s *CameraSystem) match(eflag ecs.Flag) bool {
    return matchCameraSystem(eflag, s.world)
}

func (s *CameraSystem) resizematch(eflag ecs.Flag) bool {
    return resizematchCameraSystem(eflag, s.world)
}

func (s *CameraSystem) ComponentAdded(e ecs.Entity, eflag ecs.Flag) {
    if s.match(eflag) {
        if s.view.Add(e) {
            // TODO: dispatch event that this entity was added to this system
            
        }
    } else {
        if s.view.Remove(e) {
            // TODO: dispatch event that this entity was removed from this system
            
        }
    }
}

func (s *CameraSystem) ComponentRemoved(e ecs.Entity, eflag ecs.Flag) {
    if s.match(eflag) {
        if s.view.Add(e) {
            // TODO: dispatch event that this entity was added to this system
            
        }
    } else {
        if s.view.Remove(e) {
            // TODO: dispatch event that this entity was removed from this system
            
        }
    }
}

func (s *CameraSystem) ComponentResized(cflag ecs.Flag) {
    if s.resizematch(cflag) {
        s.view.rescan()
        
    }
}

func (s *CameraSystem) ComponentWillResize(cflag ecs.Flag) {
    if s.resizematch(cflag) {
        
        s.view.clearpointers()
    }
}

func (s *CameraSystem) V() *viewCameraSystem {
    return s.view
}

func (*CameraSystem) Priority() int64 {
    return 50
}

func (s *CameraSystem) Setup(w ecs.BaseWorld) {
    if s.initialized {
        panic("CameraSystem called Setup() more than once")
    }
    s.view = newviewCameraSystem(w)
    s.world = w
    s.enabled = true
    s.initialized = true
    
}


func init() {
    ecs.RegisterSystem(func() ecs.BaseSystem {
        return &CameraSystem{}
    })
}