package core

import "errors"

type CompareResult int

const (
	IS_OLDER = CompareResult(iota)
	IS_SAME
	IS_NEWER
	HAS_CONFLICTS
)

func ChooseNewer(o1, o2 Source) (Source, error) {
	result := Compare(o1, o2)
	switch result {
	case IS_SAME:
	case IS_NEWER:
		return o1, nil
	case IS_OLDER:
		return o2, nil
	case HAS_CONFLICTS:
		return Source{}, errors.New("sources has conflicts: can not choose newer version")
	}

	return Source{}, errors.New("incorrect compare result: must be unreachable")
}

func Compare(o1, o2 Source) CompareResult {
	if !o1.IsCorrect() || !o2.IsCorrect() {
		return HAS_CONFLICTS
	}

	h1 := extractHistory(o1)
	h2 := extractHistory(o2)

	var min int
	var max int
	var shortests []History
	var longests []History

	if len(h1) <= len(h2) {
		min = len(h1)
		max = len(h2)
		shortests = h1
		longests = h2
	} else {
		min = len(h2)
		max = len(h1)
		shortests = h2
		longests = h1
	}

	if min == max {
		state := IS_SAME

		for i := 0; i < max; i++ {
			state = compareHistoryItemWithStatus(h1[i], h2[i], state)
			if state != IS_SAME {
				return HAS_CONFLICTS
			}
		}

		return state
	}

	state := IS_SAME

	for i := 0; i < min; i++ {
		state = compareHistoryItemWithStatus(shortests[i], longests[i], state)
		if state == HAS_CONFLICTS {
			return HAS_CONFLICTS
		}

		if state == IS_OLDER {
			return HAS_CONFLICTS
		}
	}

	return state
}

func extractHistory(o Source) []History {
	var result []History

	for _, oi := range o.Items {
		result = append(result, oi.History)
	}

	return result
}

func compareHistoryItemWithStatus(h1, h2 History, last CompareResult) CompareResult {
	cmd := compareHistoryItem(h1, h2)
	if cmd == HAS_CONFLICTS {
		return HAS_CONFLICTS
	}

	if last == IS_SAME {
		if cmd == IS_NEWER {
			last = IS_NEWER
		} else if cmd == IS_OLDER {
			last = IS_OLDER
		}
	} else if last == IS_NEWER && cmd != IS_NEWER {
		return HAS_CONFLICTS
	} else if last == IS_OLDER && cmd != IS_OLDER {
		return HAS_CONFLICTS
	}

	return last
}

func compareHistoryItem(hi1, hi2 History) CompareResult {
	if hi1.Time == hi2.Time && hi1.Revision == hi2.Revision {
		return IS_SAME
	}

	if hi1.Time.After(hi2.Time) {
		return IS_NEWER
	}

	if hi1.Time.Before(hi2.Time) {
		return IS_OLDER
	}

	return HAS_CONFLICTS
}

func (o Source) IsCorrect() bool {
	h := extractHistory(o)
	if len(h) < 2 {
		return true
	}

	var last History = h[0]
	var cur History

	for i := 1; i < len(h); i++ {
		cur = h[i]
		if compareHistoryItem(last, cur) != IS_OLDER {
			return false
		}
		last = cur
	}

	return true
}
