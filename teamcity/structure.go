package teamcity

import (
	"fmt"
	"time"
)

var daysOfWeek = map[string]time.Weekday{}

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		daysOfWeek[d.String()] = d
	}
}

func parseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}

func expandStringSlice(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func flattenStringSlice(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func getChangeExpandedStringList(oraw interface{}, nraw interface{}) (remove []string, add []string) {
	old := oraw.([]interface{})
	new := nraw.([]interface{})

	remove = make([]string, 0)
	add = make([]string, 0)

	for _, n := range new {
		if _, contains := sliceContainsString(old, n.(string)); !contains {
			add = append(add, n.(string))
		}
	}
	for _, o := range old {
		if _, contains := sliceContainsString(new, o.(string)); !contains {
			remove = append(remove, o.(string))
		}
	}

	return
}

func sliceContainsString(slice []interface{}, s string) (int, bool) {
	for idx, value := range slice {
		v := value.(string)
		if v == s {
			return idx, true
		}
	}
	return -1, false
}
