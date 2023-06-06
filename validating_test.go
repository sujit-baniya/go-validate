package validate

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestIssue1(t *testing.T) {
	m := map[string]interface{}{
		"names": []string{"John", "Jane", "abc"},
		"coding": []map[string]any{
			{
				"details": map[string]any{
					"em": map[string]any{
						"code":              "001",
						"encounter_uid":     "1",
						"billing_provider":  "Test provider",
						"resident_provider": "Test Resident Provider",
					},
					"cpt": []map[string]any{
						{
							"code":              "001",
							"encounter_uid":     "1",
							"work_item_uid":     "1",
							"billing_provider":  "Test provider",
							"resident_provider": "Test Resident Provider",
						},
						{
							"code":              "OBS01",
							"encounter_uid":     "1",
							"work_item_uid":     "1",
							"billing_provider":  "Test provider",
							"resident_provider": "Test Resident Provider",
						},
						{
							"code":              "SU002",
							"encounter_uid":     "1",
							"work_item_uid":     "1",
							"billing_provider":  "Test provider",
							"resident_provider": "Test Resident Provider",
						},
					},
				},
			},
		},
	}

	v := Map(m)
	v.StopOnError = false
	v.StringRule("coding.*.details", "required")
	v.StringRule("coding.*.details.em", "required")
	v.StringRule("coding.*.details.cpt.*.encounter_uid", "required")
	v.StringRule("coding.*.details.cpt.*.work_item_uid", "required")
	assert.True(t, v.Validate())
}

func TestIssue5(t *testing.T) {
	m := map[string]interface{}{
		"names": []string{"John", "Jane", "abc"},
		"cpt": []map[string]any{
			{
				"code":              "001",
				"encounter_uid":     "1",
				"work_item_uid":     "1",
				"billing_provider":  "Test provider",
				"resident_provider": "Test Resident Provider",
			},
			{
				"code":              "OBS01",
				"encounter_uid":     "1",
				"work_item_uid":     "1",
				"billing_provider":  "Test provider",
				"resident_provider": "Test Resident Provider",
			},
			{
				"code":              "SU002",
				"encounter_uid":     "1",
				"billing_provider":  "Test provider",
				"resident_provider": "Test Resident Provider",
			},
		},
	}

	v := Map(m)
	v.StopOnError = false
	v.StringRule("cpt.*.encounter_uid", "required")
	v.StringRule("cpt.*.work_item_uid", "required")
	assert.False(t, v.Validate())
}

func TestIssue3(t *testing.T) {
	m := map[string]interface{}{
		"names": []string{"John", "Jane", "abc"},
		"coding": []map[string]any{
			{
				"details": map[string]any{
					"em": map[string]any{
						"code":              "001",
						"encounter_uid":     "1",
						"billing_provider":  "Test provider",
						"resident_provider": "Test Resident Provider",
					},
				},
			},
		},
	}

	v := Map(m)
	v.StopOnError = false
	v.StringRule("coding.*.details", "required")
	v.StringRule("coding.*.details.em", "required")
	v.StringRule("coding.*.details.cpt.*.encounter_uid", "required")
	v.StringRule("coding.*.details.cpt.*.work_item_uid", "required")
	assert.True(t, v.Validate())
}

func TestIssue2(t *testing.T) {
	m := map[string]interface{}{
		"names": []string{"John", "Jane", "abc"},
		"coding": []map[string]any{
			{
				"details": map[string]any{
					"em": map[string]any{
						"code":              "001",
						"encounter_uid":     "1",
						"billing_provider":  "Test provider",
						"resident_provider": "Test Resident Provider",
					},
					"cpt": []map[string]any{
						{
							"code":              "001",
							"work_item_uid":     "1",
							"billing_provider":  "Test provider",
							"resident_provider": "Test Resident Provider",
						},
						{
							"code":              "OBS01",
							"encounter_uid":     "1",
							"work_item_uid":     "1",
							"billing_provider":  "Test provider",
							"resident_provider": "Test Resident Provider",
						},
						{
							"code":              "SU002",
							"encounter_uid":     "1",
							"work_item_uid":     "1",
							"billing_provider":  "Test provider",
							"resident_provider": "Test Resident Provider",
						},
					},
				},
			},
		},
	}

	v := Map(m)
	v.StopOnError = false
	v.StringRule("coding.*.details", "required")
	v.StringRule("coding.*.details.em", "required")
	v.StringRule("coding.*.details.cpt.*.encounter_uid", "required")
	v.StringRule("coding.*.details.cpt.*.work_item_uid", "required")
	assert.False(t, v.Validate())
}
