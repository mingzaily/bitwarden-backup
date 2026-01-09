package provider

import (
	"fmt"
	"sync"
)

// Registry Provider 注册表
type Registry struct {
	providers map[string]DestinationProvider
	mu        sync.RWMutex
}

var (
	defaultRegistry *Registry
	once            sync.Once
)

// GetRegistry 获取默认注册表
func GetRegistry() *Registry {
	once.Do(func() {
		defaultRegistry = &Registry{
			providers: make(map[string]DestinationProvider),
		}
		// 注册内置 Provider
		defaultRegistry.Register(NewLocalProvider())
		defaultRegistry.Register(NewWebDAVProvider())
		defaultRegistry.Register(NewServerProvider())
		defaultRegistry.Register(NewS3Provider())
	})
	return defaultRegistry
}

// Register 注册 Provider
func (r *Registry) Register(p DestinationProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Type()] = p
}

// Get 获取指定类型的 Provider
func (r *Registry) Get(providerType string) (DestinationProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
	return p, nil
}

// Types 返回所有已注册的 Provider 类型
func (r *Registry) Types() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.providers))
	for t := range r.providers {
		types = append(types, t)
	}
	return types
}
