package services

import "time"

// Overlaps retorna true quando dois intervalos [aStart, aEnd) e [bStart, bEnd)
// se sobrepõem. Intervalos apenas contíguos (fim == início) não conflitam.
func Overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
