package core

import (
	"errors"
	"testing"
	"time"
)
import "github.com/stretchr/testify/require"

func TestIsCorrect(t *testing.T) {
	testCorrect(t, []SourceItem{
		sourceItem("2020-01-01T10:00:00Z"),
		sourceItem("2020-01-01T11:00:00Z"),
		sourceItem("2020-01-02T11:00:00Z"),
		sourceItem("2020-02-01T01:00:00Z"),
		sourceItem("2022-01-01T01:00:00Z"),
	}, true)

	testCorrect(t, []SourceItem{
		sourceItem("2020-01-01T11:00:00Z"),
		sourceItem("2020-01-01T10:00:00Z"),
		sourceItem("2020-01-02T11:00:00Z"),
		sourceItem("2020-02-01T01:00:00Z"),
		sourceItem("2022-01-01T01:00:00Z"),
	}, false)

	testCorrect(t, []SourceItem{
		sourceItem("2020-01-01T10:00:00Z"),
		sourceItem("2020-01-02T11:00:00Z"),
		sourceItem("2020-01-01T11:00:00Z"),
		sourceItem("2020-02-01T01:00:00Z"),
		sourceItem("2022-01-01T01:00:00Z"),
	}, false)

	testCorrect(t, []SourceItem{
		sourceItem("2020-01-01T10:00:00Z"),
		sourceItem("2020-02-01T01:00:00Z"),
		sourceItem("2020-01-01T11:00:00Z"),
		sourceItem("2020-01-02T11:00:00Z"),
		sourceItem("2022-01-01T01:00:00Z"),
	}, false)

	testCorrect(t, []SourceItem{
		sourceItem("2020-01-01T10:00:00Z"),
		sourceItem("2020-01-01T11:00:00Z"),
		sourceItem("2022-01-01T01:00:00Z"),
		sourceItem("2020-01-02T11:00:00Z"),
		sourceItem("2020-02-01T01:00:00Z"),
	}, false)
}

func testCorrect(t *testing.T, items []SourceItem, expected bool) {
	src := Source{Items: items}

	if expected {
		require.True(t, src.IsCorrect())
	} else {
		require.False(t, src.IsCorrect())
	}
}

func TestCompare(t *testing.T) {
	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
			},
		},
		IS_SAME)

	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "de"),
			},
		},
		HAS_CONFLICTS)

	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T13:00:00Z", "xyz"),
			},
		},
		HAS_CONFLICTS)

	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-02T13:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		HAS_CONFLICTS)

	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T10:30:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		HAS_CONFLICTS)

	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
			},
		},
		IS_NEWER)

	testCompare(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
				sourceItemRev("2020-01-01T13:00:00Z", "xyz"),
			},
		},
		IS_OLDER)
}

func TestChooseNewer(t *testing.T) {
	testChooseNewer(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T12:00:00", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T12:00:00", "xyz"),
			},
		},
		Source{},
		errors.New("sources has conflicts: can not choose newer version"))

	testChooseNewer(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T12:00:00", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T13:00:00", "xyz"),
			},
		},
		Source{},
		errors.New("sources has conflicts: can not choose newer version"))

	testChooseNewer(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T13:00:00", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T12:00:00", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T12:00:00", "xyz"),
			},
		},
		nil)

	testChooseNewer(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T12:00:00", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00", "abc"),
				sourceItemRev("2020-01-01T11:00:00", "def"),
				sourceItemRev("2020-01-01T13:00:00", "xyz"),
			},
		},
		Source{},
		errors.New("sources has conflicts: can not choose newer version"))

	testChooseNewer(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		nil)

	testChooseNewer(t,
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
				sourceItemRev("2020-01-01T13:00:00Z", "xyz"),
			},
		},
		Source{
			Items: []SourceItem{
				sourceItemRev("2020-01-01T10:00:00Z", "abc"),
				sourceItemRev("2020-01-01T11:00:00Z", "def"),
				sourceItemRev("2020-01-01T12:00:00Z", "xyz"),
				sourceItemRev("2020-01-01T13:00:00Z", "xyz"),
			},
		},
		nil)

}

func testCompare(t *testing.T, o1, o2 Source, result CompareResult) {
	require.Equal(t, result, Compare(o1, o2))
}

func testChooseNewer(t *testing.T, o1, o2, result Source, errExpected error) {
	newer, err := ChooseNewer(o1, o2)
	if errExpected != nil {
		require.Equal(t, err, errExpected)
	} else {
		require.Nil(t, err)
		require.Equal(t, result, newer)
	}
}

func sourceItemRev(val, rev string) SourceItem {
	return SourceItem{
		History: History{
			Time:     timeParse(val),
			Revision: rev,
		},
	}
}

func sourceItem(val string) SourceItem {
	return SourceItem{
		History: History{
			Time: timeParse(val),
		},
	}
}

func timeParse(val string) time.Time {
	result, err := time.Parse(time.RFC3339, val)
	if err != nil {
		panic(err)
	}

	return result
}
