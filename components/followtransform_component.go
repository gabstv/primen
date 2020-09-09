// Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.

package components

import (
    "sort"
    

    "github.com/gabstv/ecs/v2"
)








const uuidFollowTransformComponent = "1245B961-718E-4580-AEB7-893877FD948C"
const capFollowTransformComponent = 256

type drawerFollowTransformComponent struct {
    Entity ecs.Entity
    Data   FollowTransform
}

// WatchFollowTransform is a helper struct to access a valid pointer of FollowTransform
type WatchFollowTransform interface {
    Entity() ecs.Entity
    Data() *FollowTransform
}

type slcdrawerFollowTransformComponent []drawerFollowTransformComponent
func (a slcdrawerFollowTransformComponent) Len() int           { return len(a) }
func (a slcdrawerFollowTransformComponent) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a slcdrawerFollowTransformComponent) Less(i, j int) bool { return a[i].Entity < a[j].Entity }


type mWatchFollowTransform struct {
    c *FollowTransformComponent
    entity ecs.Entity
}

func (w *mWatchFollowTransform) Entity() ecs.Entity {
    return w.entity
}

func (w *mWatchFollowTransform) Data() *FollowTransform {
    
    
    id := w.c.indexof(w.entity)
    if id == -1 {
        return nil
    }
    return &w.c.data[id].Data
}

// FollowTransformComponent implements ecs.BaseComponent
type FollowTransformComponent struct {
    initialized bool
    flag        ecs.Flag
    world       ecs.BaseWorld
    wkey        [4]byte
    data        []drawerFollowTransformComponent
    
}

// GetFollowTransformComponent returns the instance of the component in a World
func GetFollowTransformComponent(w ecs.BaseWorld) *FollowTransformComponent {
    return w.C(uuidFollowTransformComponent).(*FollowTransformComponent)
}

// SetFollowTransformComponentData updates/adds a FollowTransform to Entity e
func SetFollowTransformComponentData(w ecs.BaseWorld, e ecs.Entity, data FollowTransform) {
    GetFollowTransformComponent(w).Upsert(e, data)
}

// GetFollowTransformComponentData gets the *FollowTransform of Entity e
func GetFollowTransformComponentData(w ecs.BaseWorld, e ecs.Entity) *FollowTransform {
    return GetFollowTransformComponent(w).Data(e)
}

// WatchFollowTransformComponentData gets a pointer getter of an entity's FollowTransform.
//
// The pointer must not be stored because it may become invalid overtime.
func WatchFollowTransformComponentData(w ecs.BaseWorld, e ecs.Entity) WatchFollowTransform {
    return &mWatchFollowTransform{
        c: GetFollowTransformComponent(w),
        entity: e,
    }
}

// UUID implements ecs.BaseComponent
func (FollowTransformComponent) UUID() string {
    return "1245B961-718E-4580-AEB7-893877FD948C"
}

// Name implements ecs.BaseComponent
func (FollowTransformComponent) Name() string {
    return "FollowTransformComponent"
}

func (c *FollowTransformComponent) indexof(e ecs.Entity) int {
    i := sort.Search(len(c.data), func(i int) bool { return c.data[i].Entity >= e })
    if i < len(c.data) && c.data[i].Entity == e {
        return i
    }
    return -1
}

// Upsert creates or updates a component data of an entity.
// Not recommended to be used directly. Use SetFollowTransformComponentData to change component
// data outside of a system loop.
func (c *FollowTransformComponent) Upsert(e ecs.Entity, data interface{}) {
    v, ok := data.(FollowTransform)
    if !ok {
        panic("data must be FollowTransform")
    }
    
    id := c.indexof(e)
    
    if id > -1 {
        
        dwr := &c.data[id]
        dwr.Data = v
        
        return
    }
    
    rsz := false
    if cap(c.data) == len(c.data) {
        rsz = true
        c.world.CWillResize(c, c.wkey)
        
    }
    newindex := len(c.data)
    c.data = append(c.data, drawerFollowTransformComponent{
        Entity: e,
        Data:   v,
    })
    if len(c.data) > 1 {
        if c.data[newindex].Entity < c.data[newindex-1].Entity {
            c.world.CWillResize(c, c.wkey)
            
            sort.Sort(slcdrawerFollowTransformComponent(c.data))
            rsz = true
        }
    }
    
    if rsz {
        
        c.world.CResized(c, c.wkey)
        c.world.Dispatch(ecs.Event{
            Type: ecs.EvtComponentsResized,
            ComponentName: "FollowTransformComponent",
            ComponentID: "1245B961-718E-4580-AEB7-893877FD948C",
        })
    }
    
    c.world.CAdded(e, c, c.wkey)
    c.world.Dispatch(ecs.Event{
        Type: ecs.EvtComponentAdded,
        ComponentName: "FollowTransformComponent",
        ComponentID: "1245B961-718E-4580-AEB7-893877FD948C",
        Entity: e,
    })
}

// Remove a FollowTransform data from entity e
//
// Warning: DO NOT call remove inside the system entities loop
func (c *FollowTransformComponent) Remove(e ecs.Entity) {
    
    
    i := c.indexof(e)
    if i == -1 {
        return
    }
    
    //c.data = append(c.data[:i], c.data[i+1:]...)
    c.data = c.data[:i+copy(c.data[i:], c.data[i+1:])]
    c.world.CRemoved(e, c, c.wkey)
    
    c.world.Dispatch(ecs.Event{
        Type: ecs.EvtComponentRemoved,
        ComponentName: "FollowTransformComponent",
        ComponentID: "1245B961-718E-4580-AEB7-893877FD948C",
        Entity: e,
    })
}

func (c *FollowTransformComponent) Data(e ecs.Entity) *FollowTransform {
    
    
    index := c.indexof(e)
    if index > -1 {
        return &c.data[index].Data
    }
    return nil
}

// Flag returns the 
func (c *FollowTransformComponent) Flag() ecs.Flag {
    return c.flag
}

// Setup is called by ecs.BaseWorld
//
// Do not call this directly
func (c *FollowTransformComponent) Setup(w ecs.BaseWorld, f ecs.Flag, key [4]byte) {
    if c.initialized {
        panic("FollowTransformComponent called Setup() more than once")
    }
    c.flag = f
    c.world = w
    c.wkey = key
    c.data = make([]drawerFollowTransformComponent, 0, 256)
    c.initialized = true
    
}


func init() {
    ecs.RegisterComponent(func() ecs.BaseComponent {
        return &FollowTransformComponent{}
    })
}
