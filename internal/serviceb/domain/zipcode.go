package domain

const zipcodeLength = 8

type Zipcode string

func NewZipcode(value string) (Zipcode, error) {
	if len(value) != zipcodeLength {
		return "", ErrInvalidZipcode
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return "", ErrInvalidZipcode
		}
	}
	return Zipcode(value), nil
}

func (z Zipcode) String() string {
	return string(z)
}
