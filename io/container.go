package io

import (
	"bytes"
	"context"
	"errors"
	"image"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Container interface {
	Len() int64
	FS() Filesystem
	Load(name string) chan struct{}
	Unload(name string) (bool, error)
	LoadAll(names []string) (progress chan float64, done chan struct{})
	UnloadAll()
	Get(name string) ([]byte, error)
	GetImage(name string) (image.Image, error)
}

type container struct {
	ctx          context.Context
	fs           Filesystem
	m            sync.RWMutex
	loadedfiles  map[string][]byte
	loadingfiles map[string]chan struct{}
	loadedlen    int64 // atomic
	loadinglen   int64 // atomic
	loadingn     int32 // atomic
}

func (c *container) Len() int64 {
	return atomic.LoadInt64(&c.loadedlen)
}

func (c *container) FS() Filesystem {
	return c.fs
}

func (c *container) Load(name string) chan struct{} {
	c.m.RLock()
	xch, ok := c.loadingfiles[name]
	c.m.RUnlock()
	if ok {
		// already loading/loaded
		return xch
	}
	atomic.AddInt32(&c.loadingn, 1)
	ldch := make(chan struct{}, 1)
	go func() {
		defer close(ldch)
		c.m.Lock()
		c.loadingfiles[name] = ldch
		c.m.Unlock()
		defer atomic.AddInt32(&c.loadingn, -1)
		f, err := c.fs.Open(name)
		if err != nil {
			log.Println("container io error: " + err.Error())
			return
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			log.Println("container stat error: " + err.Error())
			return
		}
		atomic.AddInt64(&c.loadinglen, fi.Size())
		defer atomic.AddInt64(&c.loadinglen, -fi.Size())
		buf := new(bytes.Buffer)
		buf.Grow(int(fi.Size()))
		if _, err := buf.ReadFrom(f); err != nil {
			log.Println("container io error: " + err.Error())
			return
		}
		c.m.Lock()
		c.loadedfiles[name] = buf.Bytes()
		c.m.Unlock()
		atomic.AddInt64(&c.loadedlen, int64(buf.Len()))
	}()
	return ldch
}

func (c *container) Unload(name string) (bool, error) {
	c.m.RLock()
	loadingch, ok := c.loadingfiles[name]
	c.m.RUnlock()
	if !ok {
		return false, nil
	}
	select {
	case <-time.After(time.Microsecond * 500):
		return false, errors.New("resource is loading")
	case <-loadingch:
	}
	c.m.Lock()
	delete(c.loadingfiles, name)
	x := len(c.loadedfiles[name])
	delete(c.loadedfiles, name)
	c.m.Unlock()
	atomic.AddInt64(&c.loadedlen, int64(x))
	return true, nil
}

func (c *container) LoadAll(names []string) (progress chan float64, done chan struct{}) {
	progress = make(chan float64, 64)
	done = make(chan struct{})
	if len(names) < 1 {
		close(done)
		return
	}
	go func() {
		defer close(done)
		defer func() {
			select {
			case <-c.ctx.Done():
			case <-time.After(time.Second):
			}
			close(progress)
		}()
		var step float64 = 1 / float64(len(names))
		var current float64
		for _, name := range names {
			done1 := c.Load(name)
			select {
			case <-c.ctx.Done():
				return
			case <-done1:
			}
			current += step
			progress <- current
		}
	}()
	return
}

func (c *container) UnloadAll() {
	c.m.RLock()
	names := make([]string, 0, len(c.loadingfiles))
	for name, _ := range c.loadingfiles {
		names = append(names, name)
	}
	c.m.RUnlock()
	for _, name := range names {
		_, _ = c.Unload(name)
	}
}

func (c *container) Get(name string) ([]byte, error) {
	c.m.RLock()
	fd, ok := c.loadedfiles[name]
	c.m.RUnlock()
	if ok {
		return fd, nil
	}
	c.m.RLock()
	fch, ok := c.loadingfiles[name]
	c.m.RUnlock()
	if ok {
		<-fch
		c.m.RLock()
		fd = c.loadedfiles[name]
		c.m.RUnlock()
		if fd == nil {
			return nil, errors.New("x not found")
		}
		return fd, nil
	}
	// check if exists
	if _, err := c.fs.Stat(name); err != nil {
		return nil, err
	}
	ch := c.Load(name)
	<-ch
	c.m.RLock()
	fd = c.loadedfiles[name]
	c.m.RUnlock()
	if fd == nil {
		return nil, errors.New("x not found")
	}
	return fd, nil
}

func (c *container) GetImage(name string) (image.Image, error) {
	b, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	return img, err
}

func NewContainer(ctx context.Context, fs Filesystem) Container {
	c := &container{
		ctx:          ctx,
		fs:           fs,
		loadedfiles:  make(map[string][]byte),
		loadingfiles: make(map[string]chan struct{}),
	}
	return c
}
