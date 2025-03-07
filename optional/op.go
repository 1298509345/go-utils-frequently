package optional

type (
	Op[T any] func(*T)
	Validator interface {
		Validate() error
	}
)

func New[T any](def *T, ops ...Op[T]) *T {
	if def == nil {
		return nil
	}
	for _, op := range ops {
		op(def)
	}
	return def
}

func NewWithErr[T any](def *T, ops ...Op[T]) (*T, error) {
	if def == nil {
		return nil, nil
	}
	for _, op := range ops {
		op(def)
	}
	if v, ok := any(def).(Validator); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}
	return def, nil
}
