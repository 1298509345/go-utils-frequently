package optional

import (
	"fmt"
	"testing"
	"time"
)

// 示例使用（用户服务配置）
type UserConfig struct {
	Name    string
	Age     int
	Timeout time.Duration
}

// 具体Option实现
func WithName(name string) Op[UserConfig] {
	return func(c *UserConfig) {
		c.Name = name
	}
}

func WithAge(age int) Op[UserConfig] {
	return func(c *UserConfig) {
		c.Age = age
	}
}

func (u *UserConfig) Validate() error {
	if u.Name == "" || u.Age <= 0 {
		return fmt.Errorf("invalid user config")
	}
	return nil
}

func TestEg(t *testing.T) {
	config := New(
		&UserConfig{
			Timeout: 30 * time.Second,
		},
		WithName("Alice"),
		WithAge(25),
	)
	t.Log(config)

	config, err := NewWithErr(
		&UserConfig{
			Timeout: 30 * time.Second,
		},
		WithName("Alice"),
		WithAge(0),
	)
	t.Log(config, err)
}
