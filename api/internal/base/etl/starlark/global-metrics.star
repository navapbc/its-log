SUMMARIES = [
    {'operation': 'app_fhir_v3_patient_call_synthetic_count', 'tags': '2.local postman', 'value': '6'},
    {'operation': 'app_fhir_v2_patient_call_synthetic_count', 'tags': '2.local postman', 'value': '1'},
    {'operation': 'app_fhir_v1_patient_call_synthetic_count', 'tags': '2.local postman', 'value': '1'},
    {'operation': 'app_fhir_v1_coverage_call_synthetic_count', 'tags': '2.local postman', 'value': '1'},
    {'operation': 'app_fhir_v3_coverage_call_synthetic_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_fhir_v2_coverage_call_synthetic_count', 'tags': '2.local postman', 'value': '3'},
    {'operation': 'app_fhir_v1_eob_call_synthetic_count', 'tags': '2.local postman', 'value': '6'},
    {'operation': 'app_fhir_v1_call_synthetic_count', 'tags': '2.local postman', 'value': '8'},
    {'operation': 'app_fhir_v2_eob_call_synthetic_count', 'tags': '2.local postman', 'value': '9'},
    {'operation': 'app_fhir_v2_call_synthetic_count', 'tags': '2.local postman', 'value': '13'},
    {'operation': 'app_fhir_v3_eob_call_synthetic_count', 'tags': '2.local postman', 'value': '2'},
    {'operation': 'app_fhir_v3_call_synthetic_count', 'tags': '2.local postman', 'value': '12'},
    {'operation': 'app_token_refresh_response_2xx_count', 'tags': '2.local postman', 'value': '2'},
    {'operation': 'app_token_authorization_code_2xx_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_token_authorization_code_for_synthetic_bene_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_approval_view_post_ok_synthetic_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_auth_demoscope_required_choice_sharing_synthetic_bene_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_auth_ok_synthetic_bene_distinct_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_auth_ok_synthetic_bene_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_approval_view_get_ok_synthetic_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_sls_callback_ok_synthetic_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_authentication_matched_new_bene_synthetic_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_authentication_start_ok_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_medicare_login_redirect_ok_count', 'tags': '2.local postman', 'value': '4'},
    {'operation': 'app_authorize_initial_count', 'tags': '2.local postman', 'value': '12'},
    {'operation': 'app_fhir_v3_coverage_call_synthetic_count', 'tags': '1.TestApp', 'value': '1'},
    {'operation': 'app_fhir_v3_patient_call_synthetic_count', 'tags': '1.TestApp', 'value': '1'},
    {'operation': 'app_fhir_v3_eob_call_synthetic_count',	'tags': '1.TestApp', 'value': '1'},
    {'operation': 'app_fhir_v3_call_synthetic_count', 'tags': '1.TestApp', 'value': '3'},
    {'operation': 'app_fhir_v2_coverage_call_synthetic_count', 'tags': '1.TestApp', 'value': '1'},
    {'operation': 'app_fhir_v2_patient_call_synthetic_count', 'tags': '1.TestApp', 'value': '1'},
    {'operation': 'app_fhir_v2_eob_call_synthetic_count', 'tags': '1.TestApp', 'value': '1'},
    {'operation': 'app_fhir_v2_call_synthetic_count', 'tags': '1.TestApp', 'value': 	'3'},
    {'operation': 'app_token_authorization_code_for_synthetic_bene_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_approval_view_post_ok_synthetic_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_auth_demoscope_required_choice_sharing_synthetic_bene_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_auth_ok_synthetic_bene_distinct_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_auth_ok_synthetic_bene_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_approval_view_get_ok_synthetic_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_sls_callback_ok_synthetic_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_authentication_matched_new_bene_synthetic_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_authentication_start_ok_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_medicare_login_redirect_ok_count', 'tags': '1.TestApp', 'value': '2'},
    {'operation': 'app_authorize_initial_count', 'tags': '1.TestApp', 'value': '6'}
]

APP_GLOBAL_PAIRS = {
    'fhir_v1_patient_call_synthetic_count': 'app_fhir_v1_patient_call_synthetic_count',
    'fhir_v1_coverage_call_synthetic_count': 'app_fhir_v1_coverage_call_synthetic_count',
    'fhir_v1_eob_call_synthetic_count': 'app_fhir_v1_eob_call_synthetic_count',
    'fhir_v1_call_synthetic_count': 'app_fhir_v1_call_synthetic_count',
    'fhir_v1_patient_call_real_count': 'app_fhir_v1_patient_call_real_count',
    'fhir_v1_coverage_call_real_count': 'app_fhir_v1_coverage_call_real_count',
    'fhir_v1_eob_call_real_count': 'app_fhir_v1_eob_call_real_count',
    'fhir_v1_call_real_count': 'app_fhir_v1_call_real_count',
    'fhir_v1_metadata_call_count': 'app_fhir_v1_metadata_call_count',
    'fhir_v1_coverage_since_call_synthetic_count': 'app_fhir_v1_coverage_since_call_synthetic_count',
    'fhir_v1_coverage_since_call_real_count': 'app_fhir_v1_coverage_since_call_real_count',
    'fhir_v1_eob_since_call_synthetic_count': 'app_fhir_v1_eob_since_call_synthetic_count',
    'fhir_v1_eob_since_call_real_count': 'app_fhir_v1_eob_since_call_real_count',

    'fhir_v2_patient_call_synthetic_count': 'app_fhir_v2_patient_call_synthetic_count',
    'fhir_v2_coverage_call_synthetic_count': 'app_fhir_v2_coverage_call_synthetic_count',
    'fhir_v2_eob_call_synthetic_count': 'app_fhir_v2_eob_call_synthetic_count',
    'fhir_v2_call_synthetic_count': 'app_fhir_v2_call_synthetic_count',
    'fhir_v2_patient_call_real_count': 'app_fhir_v2_patient_call_real_count',
    'fhir_v2_coverage_call_real_count': 'app_fhir_v2_coverage_call_real_count',
    'fhir_v2_eob_call_real_count': 'app_fhir_v2_eob_call_real_count',
    'fhir_v2_call_real_count': 'app_fhir_v2_call_real_count',
    'fhir_v2_metadata_call_count': 'app_fhir_v2_metadata_call_count',
    'fhir_v2_coverage_since_call_synthetic_count': 'app_fhir_v2_coverage_since_call_synthetic_count',
    'fhir_v2_coverage_since_call_real_count': 'app_fhir_v2_coverage_since_call_real_count',
    'fhir_v2_eob_since_call_synthetic_count': 'app_fhir_v2_eob_since_call_synthetic_count',
    'fhir_v2_eob_since_call_real_count': 'app_fhir_v2_eob_since_call_real_count',

    'fhir_v3_generate_insurance_card_call_synthetic_count': 'app_fhir_v3_generate_insurance_card_call_synthetic_count',
    'fhir_v3_patient_call_synthetic_count': 'app_fhir_v3_patient_call_synthetic_count',
    'fhir_v3_coverage_call_synthetic_count': 'app_fhir_v3_coverage_call_synthetic_count',
    'fhir_v3_eob_call_synthetic_count': 'app_fhir_v3_eob_call_synthetic_count',
    'fhir_v3_call_synthetic_count': 'app_fhir_v3_call_synthetic_count',
    'fhir_v3_generate_insurance_card_call_real_count': 'app_fhir_v3_generate_insurance_card_call_real_count',
    'fhir_v3_patient_call_real_count': 'app_fhir_v3_patient_call_real_count',
    'fhir_v3_coverage_call_real_count': 'app_fhir_v3_coverage_call_real_count',
    'fhir_v3_eob_call_real_count': 'app_fhir_v3_eob_call_real_count',
    'fhir_v3_call_real_count': 'app_fhir_v3_call_real_count',
    'fhir_v3_metadata_call_count': 'app_fhir_v3_metadata_call_count',
    'fhir_v3_coverage_since_call_synthetic_count': 'app_fhir_v3_coverage_since_call_synthetic_count',
    'fhir_v3_coverage_since_call_real_count': 'app_fhir_v3_coverage_since_call_real_count',
    'fhir_v3_eob_since_call_synthetic_count': 'app_fhir_v3_eob_since_call_synthetic_count',
    'fhir_v3_eob_since_call_real_count': 'app_fhir_v3_eob_since_call_real_count',

    'auth_v1_v2_user_makes_it_to_permission_screen_bene_count': 'app_auth_v1_v2_user_makes_it_to_permission_screen_bene_count',
    'auth_v1_v2_user_clicks_connect_bene_count': 'app_auth_v1_v2_user_clicks_connect_bene_count',
    'auth_v3_user_makes_it_to_permission_screen_bene_count': 'app_auth_v3_user_makes_it_to_permission_screen_bene_count',
    'auth_v3_user_clicks_connect_bene_count': 'app_auth_v3_user_clicks_connect_bene_count'
}

# needed this function as starlark does not seem to have a sum function
def sum_total(list_of_summary_values):
    total = 0
    for value in list_of_summary_values:
        total += int(value)
    return total

def summarize():
    summaries = []
    for key, value in APP_GLOBAL_PAIRS.items():
        # Grab all summary records for specific app-level value
        # Starlark does not support string formatting, hence concatenation
        # We only need the value back
        summary_query = 'SELECT * from itslog_summary WHERE operation = \'' + value + '\';'
        summary_rows = query('summary', summary_query)

        # Only grab the value from the retrieved summaries
        summary_values_for_specific_metric = [s['value'] for s in summary_rows]

        total = sum_total(summary_values_for_specific_metric)
        # Construct the global summary record and append it to the list that will be returned
        global_summary = {
            'operation': key,
            'value': str(total),
            'count': 1
        }
        summaries.append(global_summary)

    return summaries

summarize()

