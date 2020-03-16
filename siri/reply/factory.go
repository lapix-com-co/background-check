package reply

// Factory returns the valid strategy for the given input from the strategies associated.
type Factory struct {
	strategies []Strategy
}

func NewFactory() *Factory {
	return &Factory{
		strategies: []Strategy{
			NewNoAccentsCapitals(),
			NewCapitals(),
			&Arithmetical{},
			&CharLength{},
			&DocumentCharts{},
			&Names{},
		},
	}
}

// Strategy returns the strategy that match the given value.
func (f Factory) Strategy(u User, i string) Strategy {
	for _, strategy := range f.strategies {
		if strategy.Is(u, i) {
			return strategy
		}
	}

	return nil
}
