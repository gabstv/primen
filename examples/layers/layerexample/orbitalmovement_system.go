// Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.

package layerexample

import (
    
    "sort"

    "github.com/gabstv/ecs/v2"
    
    "github.com/gabstv/primen/components"
    
    "github.com/gabstv/primen/components/graphics"
    
    "github.com/hajimehoshi/ebiten"
    
)









const uuidOrbitalMovementSystem = "826684C9-E190-4BF2-93D7-2FA61A5BCEEC"

type viewOrbitalMovementSystem struct {
    entities []VIOrbitalMovementSystem
    world ecs.BaseWorld
    
}

type VIOrbitalMovementSystem struct {
    Entity ecs.Entity
    
    OrbitalMovement *OrbitalMovement 
    
    Sprite *graphics.Sprite 
    
    DrawLayer *graphics.DrawLayer 
    
    Transform *components.Transform 
    
}

type sortedVIOrbitalMovementSystems []VIOrbitalMovementSystem
func (a sortedVIOrbitalMovementSystems) Len() int           { return len(a) }
func (a sortedVIOrbitalMovementSystems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedVIOrbitalMovementSystems) Less(i, j int) bool { return a[i].Entity < a[j].Entity }

func newviewOrbitalMovementSystem(w ecs.BaseWorld) *viewOrbitalMovementSystem {
    return &viewOrbitalMovementSystem{
        entities: make([]VIOrbitalMovementSystem, 0),
        world: w,
    }
}

func (v *viewOrbitalMovementSystem) Matches() []VIOrbitalMovementSystem {
    
    return v.entities
    
}

func (v *viewOrbitalMovementSystem) indexof(e ecs.Entity) int {
    i := sort.Search(len(v.entities), func(i int) bool { return v.entities[i].Entity >= e })
    if i < len(v.entities) && v.entities[i].Entity == e {
        return i
    }
    return -1
}

// Fetch a specific entity
func (v *viewOrbitalMovementSystem) Fetch(e ecs.Entity) (data VIOrbitalMovementSystem, ok bool) {
    
    i := v.indexof(e)
    if i == -1 {
        return VIOrbitalMovementSystem{}, false
    }
    return v.entities[i], true
}

func (v *viewOrbitalMovementSystem) Add(e ecs.Entity) bool {
    
    
    // MUST NOT add an Entity twice:
    if i := v.indexof(e); i > -1 {
        return false
    }
    v.entities = append(v.entities, VIOrbitalMovementSystem{
        Entity: e,
        OrbitalMovement: GetOrbitalMovementComponent(v.world).Data(e),
Sprite: graphics.GetSpriteComponent(v.world).Data(e),
DrawLayer: graphics.GetDrawLayerComponent(v.world).Data(e),
Transform: components.GetTransformComponent(v.world).Data(e),

    })
    if len(v.entities) > 1 {
        if v.entities[len(v.entities)-1].Entity < v.entities[len(v.entities)-2].Entity {
            sort.Sort(sortedVIOrbitalMovementSystems(v.entities))
        }
    }
    return true
}

func (v *viewOrbitalMovementSystem) Remove(e ecs.Entity) bool {
    
    
    if i := v.indexof(e); i != -1 {

        v.entities = append(v.entities[:i], v.entities[i+1:]...)
        return true
    }
    return false
}

func (v *viewOrbitalMovementSystem) clearpointers() {
    
    
    for i := range v.entities {
        e := v.entities[i].Entity
        
        v.entities[i].OrbitalMovement = nil
        
        v.entities[i].Sprite = nil
        
        v.entities[i].DrawLayer = nil
        
        v.entities[i].Transform = nil
        
        _ = e
    }
}

func (v *viewOrbitalMovementSystem) rescan() {
    
    
    for i := range v.entities {
        e := v.entities[i].Entity
        
        v.entities[i].OrbitalMovement = GetOrbitalMovementComponent(v.world).Data(e)
        
        v.entities[i].Sprite = graphics.GetSpriteComponent(v.world).Data(e)
        
        v.entities[i].DrawLayer = graphics.GetDrawLayerComponent(v.world).Data(e)
        
        v.entities[i].Transform = components.GetTransformComponent(v.world).Data(e)
        
        _ = e
        
    }
}

// OrbitalMovementSystem implements ecs.BaseSystem
type OrbitalMovementSystem struct {
    initialized bool
    world       ecs.BaseWorld
    view        *viewOrbitalMovementSystem
    enabled     bool
    
    paused bool 
    
    globalScale float64 
    
    radiusScale float64 
    
    xframes chan struct{} 
    
    wave1 float64 
    
    waver float64 
    
    fgs []*ebiten.Image 
    
    bgs []*ebiten.Image 
    
}

// GetOrbitalMovementSystem returns the instance of the system in a World
func GetOrbitalMovementSystem(w ecs.BaseWorld) *OrbitalMovementSystem {
    return w.S(uuidOrbitalMovementSystem).(*OrbitalMovementSystem)
}

// Enable system
func (s *OrbitalMovementSystem) Enable() {
    s.enabled = true
}

// Disable system
func (s *OrbitalMovementSystem) Disable() {
    s.enabled = false
}

// Enabled checks if enabled
func (s *OrbitalMovementSystem) Enabled() bool {
    return s.enabled
}

// UUID implements ecs.BaseSystem
func (OrbitalMovementSystem) UUID() string {
    return "826684C9-E190-4BF2-93D7-2FA61A5BCEEC"
}

func (OrbitalMovementSystem) Name() string {
    return "OrbitalMovementSystem"
}

// ensure matchfn
var _ ecs.MatchFn = matchOrbitalMovementSystem

// ensure resizematchfn
var _ ecs.MatchFn = resizematchOrbitalMovementSystem

func (s *OrbitalMovementSystem) match(eflag ecs.Flag) bool {
    return matchOrbitalMovementSystem(eflag, s.world)
}

func (s *OrbitalMovementSystem) resizematch(eflag ecs.Flag) bool {
    return resizematchOrbitalMovementSystem(eflag, s.world)
}

func (s *OrbitalMovementSystem) ComponentAdded(e ecs.Entity, eflag ecs.Flag) {
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

func (s *OrbitalMovementSystem) ComponentRemoved(e ecs.Entity, eflag ecs.Flag) {
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

func (s *OrbitalMovementSystem) ComponentResized(cflag ecs.Flag) {
    if s.resizematch(cflag) {
        s.view.rescan()
        
    }
}

func (s *OrbitalMovementSystem) ComponentWillResize(cflag ecs.Flag) {
    if s.resizematch(cflag) {
        
        s.view.clearpointers()
    }
}

func (s *OrbitalMovementSystem) V() *viewOrbitalMovementSystem {
    return s.view
}

func (*OrbitalMovementSystem) Priority() int64 {
    return 0
}

func (s *OrbitalMovementSystem) Setup(w ecs.BaseWorld) {
    if s.initialized {
        panic("OrbitalMovementSystem called Setup() more than once")
    }
    s.view = newviewOrbitalMovementSystem(w)
    s.world = w
    s.enabled = true
    s.initialized = true
    s.setupVars()
}


func init() {
    ecs.RegisterSystem(func() ecs.BaseSystem {
        return &OrbitalMovementSystem{}
    })
}
