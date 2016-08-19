package between_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/concourse/time-resource/between"
	"github.com/concourse/time-resource/models"
)

type testCase struct {
	start string
	stop  string

	timeToCompare string
	extraTime     time.Duration

	result bool
}

const exampleFormat = "3:04 PM -0700"

func (testCase testCase) Run() {
	startJSON, err := json.Marshal(testCase.start)
	Expect(err).NotTo(HaveOccurred())

	stopJSON, err := json.Marshal(testCase.stop)
	Expect(err).NotTo(HaveOccurred())

	var start, stop models.TimeOfDay

	err = json.Unmarshal(startJSON, &start)
	Expect(err).NotTo(HaveOccurred())

	err = json.Unmarshal(stopJSON, &stop)
	Expect(err).NotTo(HaveOccurred())

	timeOfDay, err := time.Parse(exampleFormat, testCase.timeToCompare)
	Expect(err).NotTo(HaveOccurred())

	// realTime := time.Now().In(timeOfDay.Location())

	result := between.Between(time.Duration(start), time.Duration(stop), timeOfDay.UTC())
	Expect(result).To(Equal(testCase.result))
}

var _ = DescribeTable("Between", (testCase).Run,
	Entry("between the start and stop time", testCase{
		start:         "2:00 AM +0000",
		stop:          "4:00 AM +0000",
		timeToCompare: "3:00 AM +0000",
		result:        true,
	}),
	Entry("between the start and stop time down to the minute", testCase{
		start:         "2:01 AM +0000",
		stop:          "2:03 AM +0000",
		timeToCompare: "2:02 AM +0000",
		result:        true,
	}),
	Entry("not between the start and stop time", testCase{
		start:         "2:00 AM +0000",
		stop:          "4:00 AM +0000",
		timeToCompare: "5:00 AM +0000",
		result:        false,
	}),
	Entry("after the stop time, down to the minute", testCase{
		start:         "2:00 AM +0000",
		stop:          "4:00 AM +0000",
		timeToCompare: "4:10 AM +0000",
		result:        false,
	}),
	Entry("before the start time, down to the minute", testCase{
		start:         "11:07 AM +0000",
		stop:          "11:10 AM +0000",
		timeToCompare: "11:05 AM +0000",
		result:        false,
	}),
	Entry("one nanosecond before the start time", testCase{
		start:         "3:04 AM +0000",
		stop:          "3:07 AM +0000",
		timeToCompare: "3:03 AM +0000",
		extraTime:     time.Minute - time.Nanosecond,
		result:        false,
	}),
	Entry("equal to the start time", testCase{
		start:         "3:04 AM +0000",
		stop:          "3:07 AM +0000",
		timeToCompare: "3:04 AM +0000",
		result:        true,
	}),
	Entry("one nanosecond before the stop time", testCase{
		start:         "3:04 AM +0000",
		stop:          "3:07 AM +0000",
		timeToCompare: "3:06 AM +0000",
		extraTime:     time.Minute - time.Nanosecond,
		result:        true,
	}),
	Entry("equal to the stop time", testCase{
		start:         "3:04 AM +0000",
		stop:          "3:07 AM +0000",
		timeToCompare: "3:07 AM +0000",
		result:        false,
	}),
	Entry("between the start and stop time but on a different day", testCase{
		start:         "2:00 AM +0000",
		stop:          "4:00 AM +0000",
		timeToCompare: "3:00 AM +0000",
		result:        true,
	}),

	// This is a regression test against how we used to determine the day offset,
	// which was by subtracting from midnight in UTC. Pretty funky.
	Entry("start and stop whose time of day plus timezone offset results in more than a day compared to midnight UTC", testCase{
		start:         "11:00 PM -0700",
		stop:          "11:30 PM -0700",
		timeToCompare: "11:15 PM -0700",
		result:        true,
	}),

	// Our date parsing library always returns the date as 1/1 since we only
	// give it a time. If the stop time is before the start time then assume
	// that the stop is in the next day.
	Entry("between the start and stop time but the stop time is before the start time", testCase{
		start:         "5:00 AM +0000",
		stop:          "1:00 AM +0000",
		timeToCompare: "6:00 AM +0000",
		result:        true,
	}),
	Entry("between the start and stop time but the stop time is before the start time (ignoring the date)", testCase{
		start:         "5:00 AM +0000",
		stop:          "1:00 AM +0000",
		timeToCompare: "6:00 AM +0000",
		result:        true,
	}),
	Entry("between the start and stop time but the stop time is before the start time (when the time to compare is in the early hours)", testCase{
		start:         "8:00 PM +0000",
		stop:          "8:00 AM +0000",
		timeToCompare: "1:00 AM +0000",
		extraTime:     24 * time.Hour, // real life doesn't use time parsing; it'll actually be a day in advance
		result:        true,
	}),
	Entry("between the start and stop time but the stop time is before the start time", testCase{
		start:         "5:00 AM +0000",
		stop:          "1:00 AM +0000",
		timeToCompare: "4:00 AM +0000",
		result:        false,
	}),

	Entry("between the start and stop time but the compare time is in a different timezone", testCase{
		start:         "2:00 AM -0600",
		stop:          "6:00 AM -0600",
		timeToCompare: "1:00 AM -0700",
		result:        true,
	}),
)
