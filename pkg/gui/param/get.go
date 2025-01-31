package param

import "math"

// Get single string value based on group and label
func GetString(group, label string) (string, error) {
	if sParams[group] == nil {
		return "", ErrParameterNotFound
	}

	if p := sParams[group].GetParam(label); p != nil {
		return p.Get()
	}

	return "", ErrParameterNotFound
}

// Get all string values based on group
func GetStrings(group string) ([]string, error) {
	if sParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return sParams[group].GetValues()
}

// Get float minimum for the parameter specified by group/label
func GetFloatMinimum(group string, label string) (float64, error) {
	if fParams[group] == nil {
		return 0, ErrParameterNotFound
	}

	if p := fParams[group].GetParam(label); p != nil {
		if len(p.GetRelatives()) > 0 {
			return p.GetRelatives()[0].Get()
		} else {
			return -math.MaxFloat64, nil
		}

	}

	return 0, ErrParameterNotFound
}

// Get float maximum for the parameter specified by group/label
func GetFloatMaximum(group string, label string) (float64, error) {
	if fParams[group] == nil {
		return 0, ErrParameterNotFound
	}

	if p := fParams[group].GetParam(label); p != nil {
		if len(p.GetRelatives()) > 1 {
			return p.GetRelatives()[1].Get()
		} else {
			return math.MaxFloat64, nil
		}
	}

	return 0, ErrParameterNotFound
}

// Get single float value based on group and label
func GetFloat(group, label string) (float64, error) {
	if fParams[group] == nil {
		return 0, ErrParameterNotFound
	}

	if p := fParams[group].GetParam(label); p != nil {
		return p.Get()
	}

	return 0, ErrParameterNotFound
}

// Get checkbox value
func GetFloatCheck(group string, label string) (bool, error) {
	if fParams[group] == nil {
		return false, ErrParameterNotFound
	}

	if p := fParams[group].GetParam(label); p != nil {
		return p.IsChecked(), nil
	}

	return false, ErrParameterNotFound
}

// Get all checkbox values based on group
func GetFloatChecked(group string) ([]bool, error) {
	if fParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return fParams[group].GetChecked(), nil
}

// Get all float minima based on group
func GetFloatMinima(group string) ([]float64, error) {
	if fParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return fParams[group].GetMinima()
}

// Get all float maxima based on group
func GetFloatMaximas(group string) ([]float64, error) {
	if fParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return fParams[group].GetMaxima()
}

// Get all float values based on group
func GetFloats(group string) ([]float64, error) {
	if fParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return fParams[group].GetValues()
}

// Get single int value based on group and label
func GetInt(group, label string) (int, error) {
	if iParams[group] == nil {
		return 0, ErrParameterNotFound
	}

	if p := iParams[group].GetParam(label); p != nil {
		return p.Get()
	}

	return 0, ErrParameterNotFound
}

// Get all int values based on group
func GetInts(group string) ([]int, error) {
	if iParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return iParams[group].GetValues()
}

// Get checkbox value
func GetIntCheck(group string, label string) (bool, error) {
	if iParams[group] == nil {
		return false, ErrParameterNotFound
	}

	if p := iParams[group].GetParam(label); p != nil {
		return p.IsChecked(), nil
	}

	return false, ErrParameterNotFound
}

// Get all checkbox values based on group
func GetIntChecked(group string) ([]bool, error) {
	if iParams[group] == nil {
		return nil, ErrParameterNotFound
	}

	return iParams[group].GetChecked(), nil
}
