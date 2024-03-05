package entity

import "github.com/davfer/archit/helpers/sli"

func Contains[K Entity](entities []K, e K) (ok bool) {
	if e.GetId().IsEmpty() {
		return
	}

	_, ok = sli.Find(entities, func(i K) bool {
		return !i.GetId().IsEmpty() && i.GetId().Equals(e.GetId())
	})

	return
}
