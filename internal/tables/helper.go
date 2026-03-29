package tables

import (
	"fmt"
	"html/template"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// Global instances for performance (thread-safe)
var (
	decoder  = form.NewDecoder()
	validate = validator.New()
)

// MapAndValidate decodes GoAdmin values and runs struct-tag validation
func MapAndValidate[T any](values map[string][]string) (*T, error) {
	// fmt.Printf("DEBUG: values map content: %#v\n", values)
	var result T

	// 1. Decode map[string][]string into the struct
	if err := decoder.Decode(&result, values); err != nil {
		return nil, err
	}

	// 2. Validate the struct based on 'validate' tags
	if err := validate.Struct(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

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
