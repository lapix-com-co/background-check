package reply

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type alwaysFalseStrategy struct{}

func (a alwaysFalseStrategy) Is(u User, s string) bool { return false }

func (a alwaysFalseStrategy) Answer(u User, s string) (string, bool, error) {
	return "", false, nil
}

type alwaysTrueStrategy struct{}

func (a alwaysTrueStrategy) Is(u User, s string) bool { return true }

func (a alwaysTrueStrategy) Answer(u User, s string) (string, bool, error) {
	return "", true, nil
}

func TestReplyFactory(t *testing.T) {
	t.Run("will choose the valid strategy", func(t *testing.T) {
		question := ""
		factory := &Factory{
			strategies: []Strategy{
				&alwaysFalseStrategy{},
				&alwaysTrueStrategy{},
			},
		}

		strategy := factory.Strategy(User{}, question)

		if !reflect.DeepEqual(&alwaysTrueStrategy{}, strategy) {
			t.Errorf("expect alwaysTrueStrategy but got %v", strategy)
		}
	})

	t.Run("will not return a strategy if it does not match", func(t *testing.T) {
		question := ""
		factory := &Factory{
			strategies: []Strategy{
				&alwaysFalseStrategy{},
			},
		}

		strategy := factory.Strategy(User{}, question)

		require.Nil(t, strategy)
	})
}
