package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/snabb/isoweek"
	"github.com/timchurchard/kobo-readstat/pkg"
)

func Test_readsInWeekToTableLine(t *testing.T) {
	type args struct {
		stats     pkg.Stats
		weekStart time.Time
		hour      int
	}

	exampleStats := pkg.Stats{
		Years: map[int]pkg.YearStats{
			2000: {
				Months: nil,
				Weeks: map[int]pkg.MonthStats{
					3: {
						FinishedBooks:    nil,
						FinishedArticles: nil,
						Books: map[string]*pkg.StatsBook{
							"A": {
								Reads: []pkg.StatsRead{
									{Time: "2000-01-17T15:04:05.000", Duration: 555},
									{Time: "2000-01-18T15:14:05.000", Duration: 666},
									{Time: "2000-01-19T15:24:05.000", Duration: 1234},
									{Time: "2000-01-20T15:34:05.000", Duration: 1234},
									{Time: "2000-01-21T15:44:05.000", Duration: 1234},
									{Time: "2000-01-22T15:54:05.000", Duration: 1234},
									{Time: "2000-01-23T15:04:05.000", Duration: 6666},
								},
							},
						},
						Articles: nil,
					},
				},
			},
		},
		Content: nil,
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "2000 week 3 should mark ",
			args: args{
				stats:     exampleStats,
				weekStart: isoweek.StartTime(2000, 3, time.UTC),
				hour:      15,
			},
			want: "|1500|-     | --   |  --- |   ---|    --|     -|------|2h42m37s|",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := readsInWeekToTableLine(tt.args.stats, tt.args.weekStart, tt.args.hour)
			fmt.Printf("%s\n%s\n", got, tt.want)

			if got != tt.want {
				t.Errorf("readsInWeekToTableLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
