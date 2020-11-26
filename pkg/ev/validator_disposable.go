package ev

type DisposableInterface interface {
	disposable(domain string) bool
}

type SetDisposable struct {
	set StringSet
}

func (s SetDisposable) disposable(domain string) bool {
	_, ok := s.set[domain]
	return ok
}

type DisposableValidatorInterface interface {
	ValidatorInterface
}

func NewDisposableValidator(d DisposableInterface) ValidatorInterface {
	return DisposableValidator{d}
}

type DisposableValidator struct {
	d DisposableInterface
}

func (d DisposableValidator) Validate(email EmailAddressInterface) ValidationResultInterface {
	return NewValidatorResult(d.d.disposable(email.GetDomain()), nil, nil)
}
