package entity

func Contains[K Entity](entities []K, e K) bool {
	if e.GetId().IsEmpty() {
		return false
	}

	for _, i := range entities {
		if !i.GetId().IsEmpty() && i.GetId().Equals(e.GetId()) {
			return true
		}
	}

	return false
}
