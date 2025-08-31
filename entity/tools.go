package entity

import "github.com/davfer/archit/helpers/sli"

func Contains[K Entity](entities []K, e K) (ok bool) {
	if e.GetID().IsEmpty() {
		return
	}

	_, ok = sli.Find(entities, func(i K) bool {
		return !i.GetID().IsEmpty() && i.GetID().Equals(e.GetID())
	})

	return
}
