// +build !windows

package winpath

func paths() ([]string, error) {
	return nil, ErrWrongPlatform
}

func add(string) (bool, error) {
	return false, ErrWrongPlatform
}

func remove(string) (bool, error) {
	return false, ErrWrongPlatform
}
