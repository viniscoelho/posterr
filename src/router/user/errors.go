package user

type InvalidToggleError struct{}

func (e InvalidToggleError) Error() string {
	return "invalid toggle selected"
}
