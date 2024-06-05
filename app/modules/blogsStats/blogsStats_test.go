package blogsStats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_findNonDuplicates(t *testing.T) {
	arr1 := []string{ "topic1" }
	arr2 := []string{ "topic1", "topic2" }
	res := findNonDuplicates(arr1, arr2)
	assert.Equal(t, []string{"topic2"}, res)
}

func Test_toArray(t *testing.T) {
	res := toArray("res1,res2,res3")
	assert.Equal(t, []string{"res1", "res2", "res3"}, res)
}

func Test_monthFormat(t *testing.T) {
	specificTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC) // Set hour, minute, etc. to 0 for midnight
	assert.Equal(t, "2020-01", specificTime.Format(monthFormat))
}

func Test_dayFormat(t *testing.T) {
	specificTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC) // Set hour, minute, etc. to 0 for midnight
	assert.Equal(t, "2020-01-02", specificTime.Format(dayFormat))
}
