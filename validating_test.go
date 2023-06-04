package validate

/*
import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestIssue(t *testing.T) {
	m := map[string]interface{}{
		"names": []string{"John", "Jane", "abc"},
		"cpt": []map[string]any{
			{
				"code":              "001",
				"encounter_uid":     "1",
				"billing_provider":  "Test provider",
				"resident_provider": "Test Resident Provider",
			},
			{
				"code":              "OBS01",
				"encounter_uid":     "1",
				"billing_provider":  "Test provider",
				"resident_provider": "Test Resident Provider",
			},
			{
				"code":              "SU002",
				"billing_provider":  "Test provider",
				"resident_provider": "Test Resident Provider",
			},
		},
	}

	v := Map(m)
	v.StopOnError = false
	v.StringRule("cpt.*.encounter_uid", "required")
	assert.False(t, v.Validate())
}

func TestIssue_124(t *testing.T) {
	m := map[string]interface{}{
		"names": []string{"John", "Jane", "abc"},
		"address": []map[string]any{
			{"number": "1b", "country": "en"},
			{"number": "1", "country": "cz"},
		},
	}

	v := Map(m)
	v.StopOnError = false
	v.StringRule("address.*.abc", "required")

	assert.False(t, v.Validate())
	assert.Error(t, v.Errors.ErrOrNil())
	// dump.Println(v.Errors)

	// TODO how to use on struct.
	// type user struct {
	// 	Tags []string `json:"tags" validate:"required|slice"`
	// }
}
*/
