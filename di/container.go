package di

import "go.uber.org/dig"

// 依赖注入实例管理
type Container interface {
	Register(constructor interface{}, opts ...dig.ProvideOption)
	Call(function interface{}, opts ...dig.InvokeOption) error
}

// 依赖注入实例管理，通过封装uber的dig库，实现实例的注册、依赖反转控制
type AppContainer struct {
	*dig.Container
}

func (s *AppContainer) Register(constructor interface{}, opts ...dig.ProvideOption) {
	err := s.Provide(constructor, opts...)
	if err != nil {
		panic(err)
	}
}

func (s *AppContainer) Call(function interface{}, opts ...dig.InvokeOption) error {
	return s.Invoke(function, opts...)
}

var (
	appContainer *AppContainer = nil
)

func GetContainer() Container {
	if appContainer == nil {
		appContainer = &AppContainer{
			dig.New(),
		}
	}

	return appContainer
}
