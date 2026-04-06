package tables

import (
	"fmt"
	"html/template"
)

func printDualListBoxJS(sourceField, targetField, url string, params ...map[string]any) template.HTML {
	var q string
	if len(params) > 0 {
		for _, param := range params {
			for k, v := range param {
				q += fmt.Sprintf("%s=%s&", k, v)
			}
		}
	}
	if q != "" {
		q = q[:len(q)-1]
	}
	js := fmt.Sprintf(`
	// Trigger once on load so the first selection is pre-populated
    // if there is already a value selected in the dropdown
	// Define the source and target SelectBox element
	var $sourceField = $('[name="%s"]')
	var $targetField = $('[name="%s"]')

    $sourceField.off('change').on('change', function() {
		var value = $(this).val();
		var id = $('[name="id"]').val() || 0;
		if (!value) {
			$targetField.empty().bootstrapDualListbox('refresh');
			return;
		}

		$.ajax({
			url: '%s?id=' + id + '&value=' + value + '&%s',
			type: 'POST',
			success: function(response) {
				if (response.code === 200) {
					var options = response.data;

					// Clear existing options from both panels
					$targetField.find('option').remove();
					// Step 2: Clear ALL options from the underlying hidden select
					$targetField.empty();

					// Add new options
					if (response.data && Array.isArray(response.data)) {
						$.each(options, function(i, opt) {
							$targetField.append(
								$('<option>', { value: opt.value, text: opt.text })
									.prop('selected', !!opt.selected)
							);
						});
					}

					// Refresh the Bootstrap Dual Listbox
					$targetField.bootstrapDualListbox('refresh', true);
				}
			},
			error: function(xhr) {
				console.error('timeslot fetch failed:', xhr.status, xhr.responseText);
			}
		});
	});


	// Check if we've already initialized this specific field
	// ensure that this onchange event won't be fired multiple times
	if (!$sourceField.data('init-dual-box')) {

		// Mark as initialized
		$sourceField.data('init-dual-box', true);

		// Trigger initial load
		var initialValue = $sourceField.val();
		if (initialValue) {
			$sourceField.trigger('change');
		}
	}
`, sourceField, targetField, url, q)

	return template.HTML(js)
}

func printSampleReqJSON() string {
	return `
[
	{
		"subject": {
			"id": 1000,
			"name": "math",
			"is_heavy": true,
			"requires_lab": false
		},
		"teacher": {
			"id": 1000,
			"first_name": "Jone",
			"last_name": "Doe"
		},
		"class": {
			"id":1047,
			"grade": 1,
			"class": "1"
		},
		"weekly_sessions": 5,
		"min_day_gap": 0,
		"preferred_days": "1,2,3,4"
	},
	{
		"subject": {
			"id": 1001,
			"name": "chinese",
			"is_heavy": true,
			"requires_lab": false
		},
		"teacher": {
			"id": 1000,
			"first_name": "Frank",
			"last_name": "Joe"
		},
		"class": {
			"id":1048,
			"grade": 1,
			"class": "2"
		},
		"weekly_sessions": 5,
		"min_day_gap": 0,
		"preferred_days": "1,2,3,4"
	}
]
`
}

func printSampleSubjectJSON() string {
	return `
[
	{
		"name": "mathematics",
		"code": "math",
		"description": "this is math subject example descripton",
		"is_heavy": true,
		"requires_lab": false
	},
	{
		"name": "science",
		"code": "sci",
		"description": "this is math subject example descripton....",
		"is_heavy": true,
		"requires_lab": true
	},
	{
		"name": "music",
		"code": "mus",
		"description": "this is math subject example descripton....",
		"is_heavy": false,
		"requires_lab": false
	}
]
`
}

func printSampleTimeslotsJSON() string {
	return `
[
    {"day": 1, "start_time": "09:00", "end_time": "09:45"},
    {"day": 1, "start_time": "10:00", "end_time": "10:45"},
    {"day": 1, "start_time": "11:00", "end_time": "11:45"},
    {"day": 1, "start_time": "13:00", "end_time": "13:45"},
    {"day": 1, "start_time": "14:00", "end_time": "14:45"},
    {"day": 1, "start_time": "15:00", "end_time": "15:45"},

    {"day": 2, "start_time": "09:00", "end_time": "09:45"},
    {"day": 2, "start_time": "10:00", "end_time": "10:45"},
    {"day": 2, "start_time": "11:00", "end_time": "11:45"},
    {"day": 2, "start_time": "13:00", "end_time": "13:45"},
    {"day": 2, "start_time": "14:00", "end_time": "14:45"},
    {"day": 2, "start_time": "15:00", "end_time": "15:45"},

    {"day": 3, "start_time": "09:00", "end_time": "09:45"},
    {"day": 3, "start_time": "10:00", "end_time": "10:45"},
    {"day": 3, "start_time": "11:00", "end_time": "11:45"},
    {"day": 3, "start_time": "13:00", "end_time": "13:45"},
    {"day": 3, "start_time": "14:00", "end_time": "14:45"},
    {"day": 3, "start_time": "15:00", "end_time": "15:45"},

    {"day": 4, "start_time": "09:00", "end_time": "09:45"},
    {"day": 4, "start_time": "10:00", "end_time": "10:45"},
    {"day": 4, "start_time": "11:00", "end_time": "11:45"},
    {"day": 4, "start_time": "13:00", "end_time": "13:45"},
    {"day": 4, "start_time": "14:00", "end_time": "14:45"},
    {"day": 4, "start_time": "15:00", "end_time": "15:45"},

    {"day": 5, "start_time": "09:00", "end_time": "09:45"},
    {"day": 5, "start_time": "10:00", "end_time": "10:45"},
    {"day": 5, "start_time": "11:00", "end_time": "11:45"},
    {"day": 5, "start_time": "13:00", "end_time": "13:45"},
    {"day": 5, "start_time": "14:00", "end_time": "14:45"},
    {"day": 5, "start_time": "15:00", "end_time": "15:45"}
]
`
}
