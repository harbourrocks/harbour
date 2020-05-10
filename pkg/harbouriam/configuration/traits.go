package configuration

// RedisTrait returns redis config
type IamTrait interface {
	GetIamConfig() Options
	SetIamConfig(Options)
}

// RequestModel holds the request
type IamModel struct {
	redisOptions Options
}

func (m IamModel) GetIamConfig() Options {
	return m.redisOptions
}

func (m *IamModel) SetIamConfig(s Options) {
	m.redisOptions = s
}

func AddIamConfig(trait IamTrait, s Options) {
	trait.SetIamConfig(s)
}
