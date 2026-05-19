# Sample of actual itslog_event records. Used in testing
# To use in testing, uncomment out the summarize() call at the bottom of the file,
# comment out the query where events is set.
# And in summarize, rather than for event in events, comment that line out and uncomment the
# for event in EVENTS line
EVENTS = [
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Patient/-279722155", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Patient/-279722155", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Patient/-279722155", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Patient/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Patient/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Patient/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v2/fhir/Patient/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v1/fhir/Patient/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Coverage/dual--279722155", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Coverage/dual--279722155", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v1/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v2/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v2/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v2/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token/", "request_method": "POST", "response_code": 200, "type": "request_response_middleware"},
    {"action": "authorized", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token/", "type": "AccessToken", "fhir_id_v2": "-10000010257297", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v2/o/authorize/cd8c4e8d-5c61-4ee0-86b1-53258e016d51/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"allow": "True", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "auth_status": "OK", "path": "/v2/o/authorize/cd8c4e8d-5c61-4ee0-86b1-53258e016d51/", "share_demographic_scopes": "True", "type": "Authorization", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/v2/o/authorize/cd8c4e8d-5c61-4ee0-86b1-53258e016d51/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "path": "/mymedicare/sls-callback", "request_method": "GET", "response_code": 302, "type": "request_response_middleware", "auth_path": "/v2/o/authorize/cd8c4e8d-5c61-4ee0-86b1-53258e016d51/"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "type": "mymedicare_cb:get_and_update_user_initial_auth", "app_id": "2", "app_name": "local postman"},
    {"auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "type": "Authentication:success", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257297", "fhir_id_v3": "-279722155", "type": "mymedicare_cb:create_beneficiary_record", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-279722155", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257297", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "sls_userinfo_status_code": 200, "type": "Authentication:start", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b464-cb92-40f2-8897-414c79fcd62f", "type": "SLSx_userinfo", "response_code": 403, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b464-cb92-40f2-8897-414c79fcd62f", "type": "SLSx_userinfo", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/sso/session", "type": "SLSx_token", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/mymedicare/login", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/v2/o/authorize/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/Coverage/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/Coverage/:resource_id", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/Patient/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token/", "request_method": "POST", "response_code": 200, "type": "request_response_middleware"},
    {"action": "authorized", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token/", "type": "AccessToken", "fhir_id_v2": "-10000010256655", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "path": "/v2/o/authorize/d476a6bd-d146-47f9-bf2b-147d69dd79f1/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"allow": "True", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "auth_status": "OK", "path": "/v2/o/authorize/d476a6bd-d146-47f9-bf2b-147d69dd79f1/", "share_demographic_scopes": "True", "type": "Authorization", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "path": "/v2/o/authorize/d476a6bd-d146-47f9-bf2b-147d69dd79f1/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "path": "/mymedicare/sls-callback", "request_method": "GET", "response_code": 302, "type": "request_response_middleware", "auth_path": "/v2/o/authorize/d476a6bd-d146-47f9-bf2b-147d69dd79f1/"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "type": "mymedicare_cb:get_and_update_user_initial_auth", "app_id": "2", "app_name": "local postman"},
    {"auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "type": "Authentication:success", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256655", "fhir_id_v3": "-555129999", "type": "mymedicare_cb:create_beneficiary_record", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-555129999", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256655", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "sls_userinfo_status_code": 200, "type": "Authentication:start", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b45f-4812-4789-a1bc-2ea791eed6ed", "type": "SLSx_userinfo", "response_code": 403, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b45f-4812-4789-a1bc-2ea791eed6ed", "type": "SLSx_userinfo", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/sso/session", "type": "SLSx_token", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/mymedicare/login", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/v2/o/authorize/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/jsi18n/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/bluebutton/myapplication/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/bluebutton/myapplication/2/change/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/jsi18n/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/bluebutton/myapplication/2/change/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/jsi18n/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/bluebutton/myapplication/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/login/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/login/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin", "request_method": "GET", "response_code": 301, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "path": "/v3/o/token/", "request_method": "POST", "response_code": 200, "type": "request_response_middleware", "auth_grant_type": "refresh_token"},
    {"fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "type": "mymedicare_cb:get_and_update_user_refresh"},
    {"fhir_id_v2": "-982745970", "type": "fhir.server.authentication.match_fhir_id"},
    {"path": "patient search", "type": "fhir_auth_pre_fetch"},
    {"fhir_id_v2": "-10000010257292", "type": "fhir.server.authentication.match_fhir_id"},
    {"path": "patient search", "type": "fhir_auth_pre_fetch"},
    {"action": "authorized", "path": "/v3/o/token/", "type": "AccessToken", "fhir_id_v2": "-10000010257292", "app_id": 2, "app_name": "local postman"},
    {"action": "revoked", "type": "AccessToken", "fhir_id_v2": "-10000010257292", "app_id": 2, "app_name": "local postman"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/Patient/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "path": "/v/fhir/Coverage/", "request_method": "GET", "response_code": 404, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token/", "request_method": "POST", "response_code": 200, "type": "request_response_middleware"},
    {"action": "authorized", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token/", "type": "AccessToken", "fhir_id_v2": "-10000010255002", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "path": "/v2/o/authorize/f7c4d629-74e6-4451-88b1-3846648a4f29/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"allow": "True", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "auth_status": "OK", "path": "/v2/o/authorize/f7c4d629-74e6-4451-88b1-3846648a4f29/", "share_demographic_scopes": "True", "type": "Authorization", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "path": "/v2/o/authorize/f7c4d629-74e6-4451-88b1-3846648a4f29/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "path": "/mymedicare/sls-callback", "request_method": "GET", "response_code": 302, "type": "request_response_middleware", "auth_path": "/v2/o/authorize/f7c4d629-74e6-4451-88b1-3846648a4f29/"},
    {"auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "type": "Authentication:success", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "type": "mymedicare_cb:get_and_update_user_initial_auth", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010255002", "fhir_id_v3": "-866682817", "type": "mymedicare_cb:create_beneficiary_record", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-866682817", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010255002", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "sls_userinfo_status_code": 200, "type": "Authentication:start", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b44f-1246-40d6-a723-e347a58fb0a1", "type": "SLSx_userinfo", "response_code": 403, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b44f-1246-40d6-a723-e347a58fb0a1", "type": "SLSx_userinfo", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/sso/session", "type": "SLSx_token", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/mymedicare/login", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/v2/o/authorize/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 403, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/ExplanationOfBenefit/carrier--10000930145217", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v1/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v2/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 2, "app_name": "local postman", "path": "/v3/o/token/", "request_method": "POST", "response_code": 200, "type": "request_response_middleware", "auth_grant_type": "refresh_token"},
    {"fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "type": "mymedicare_cb:get_and_update_user_refresh"},
    {"fhir_id_v2": "-982745970", "type": "fhir.server.authentication.match_fhir_id"},
    {"path": "patient search", "type": "fhir_auth_pre_fetch"},
    {"fhir_id_v2": "-10000010257292", "type": "fhir.server.authentication.match_fhir_id"},
    {"path": "patient search", "type": "fhir_auth_pre_fetch"},
    {"action": "authorized", "path": "/v3/o/token/", "type": "AccessToken", "fhir_id_v2": "-10000010257292", "app_id": 2, "app_name": "local postman"},
    {"action": "revoked", "type": "AccessToken", "fhir_id_v2": "-10000010257292", "app_id": 2, "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v3/o/token/", "request_method": "POST", "response_code": 200, "type": "request_response_middleware"},
    {"action": "authorized", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v3/o/token/", "type": "AccessToken", "fhir_id_v2": "-10000010257292", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/o/authorize/b4733b07-0a7d-41a7-99f2-2f6c1b3a942e/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"allow": "True", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "auth_status": "OK", "path": "/v3/o/authorize/b4733b07-0a7d-41a7-99f2-2f6c1b3a942e/", "share_demographic_scopes": "True", "type": "Authorization", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/v3/o/authorize/b4733b07-0a7d-41a7-99f2-2f6c1b3a942e/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "path": "/mymedicare/sls-callback", "request_method": "GET", "response_code": 302, "type": "request_response_middleware", "auth_path": "/v3/o/authorize/b4733b07-0a7d-41a7-99f2-2f6c1b3a942e/"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "type": "mymedicare_cb:get_and_update_user_initial_auth", "app_id": "2", "app_name": "local postman"},
    {"auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "path": "v3/mymedicare/sls-callback", "type": "Authentication:success", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257292", "fhir_id_v3": "-982745970", "type": "mymedicare_cb:create_beneficiary_record", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-982745970", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010257292", "type": "fhir.server.authentication.match_fhir_id", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "v3/mymedicare/sls-callback", "sls_userinfo_status_code": 200, "type": "Authentication:start", "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b464-30c3-44a9-8404-b605926be2af", "type": "SLSx_userinfo", "response_code": 403, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b464-30c3-44a9-8404-b605926be2af", "type": "SLSx_userinfo", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"auth_require_demographic_scopes": "True", "path": "/sso/session", "type": "SLSx_token", "response_code": 200, "app_id": "2", "app_name": "local postman"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/mymedicare/login", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "2", "app_name": "local postman", "auth_require_demographic_scopes": "True", "path": "/v3/o/authorize/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/logout/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/jsi18n/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/auth/user/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/jsi18n/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/axes/accessfailurelog/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/login/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/favicon.ico", "request_method": "GET", "response_code": 404, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/login/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/admin", "request_method": "GET", "response_code": 301, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/connect/userinfo", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/fhir/Patient/-35180292", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/connect/userinfo", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/v3/o/token", "request_method": "POST", "response_code": 200, "type": "request_response_middleware", "auth_grant_type": "authorization_code"},
    {"action": "authorized", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v3/o/token", "type": "AccessToken", "fhir_id_v2": "-10000010256951", "app_id": "1", "app_name": "TestApp"},
    {"app_id": "1", "app_name": "TestApp", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/o/authorize/3a491c4b-3e32-4d68-a406-a6332692336a/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"allow": "True", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "auth_status": "OK", "path": "/v3/o/authorize/3a491c4b-3e32-4d68-a406-a6332692336a/", "share_demographic_scopes": "True", "type": "Authorization", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "app_id": "1", "app_name": "TestApp"},
    {"app_id": "1", "app_name": "TestApp", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/v3/o/authorize/3a491c4b-3e32-4d68-a406-a6332692336a/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "1", "app_name": "TestApp", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "path": "/mymedicare/sls-callback", "request_method": "GET", "response_code": 302, "type": "request_response_middleware", "auth_path": "/v3/o/authorize/3a491c4b-3e32-4d68-a406-a6332692336a/"},
    {"auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "path": "v3/mymedicare/sls-callback", "type": "Authentication:success", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "type": "mymedicare_cb:get_and_update_user_initial_auth", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256951", "fhir_id_v3": "-35180292", "type": "mymedicare_cb:create_beneficiary_record", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-35180292", "type": "fhir.server.authentication.match_fhir_id", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256951", "type": "fhir.server.authentication.match_fhir_id", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "v3/mymedicare/sls-callback", "sls_userinfo_status_code": 200, "type": "Authentication:start", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b461-03d8-4911-ba71-1bff75735457", "type": "SLSx_userinfo", "response_code": 403, "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b461-03d8-4911-ba71-1bff75735457", "type": "SLSx_userinfo", "response_code": 200, "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "/sso/session", "type": "SLSx_token", "response_code": 200, "app_id": "1", "app_name": "TestApp"},
    {"app_id": "1", "app_name": "TestApp", "auth_require_demographic_scopes": "True", "path": "/mymedicare/login", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "1", "app_name": "TestApp", "auth_require_demographic_scopes": "True", "path": "/v3/o/authorize/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/connect/userinfo", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/fhir/Coverage/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/fhir/Patient/-10000010256645", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/fhir/ExplanationOfBenefit/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": 1, "app_name": "TestApp", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/connect/userinfo", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/v2/o/token", "request_method": "POST", "response_code": 200, "type": "request_response_middleware", "auth_grant_type": "authorization_code"},
    {"action": "authorized", "auth_crosswalk_action": "C", "auth_grant_type": "authorization_code", "auth_require_demographic_scopes": "True", "path": "/v2/o/token", "type": "AccessToken", "fhir_id_v2": "-10000010256645", "app_id": "1", "app_name": "TestApp"},
    {"app_id": "1", "app_name": "TestApp", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/o/authorize/3792648d-46e7-4705-9cb0-9ba18f0208fc/", "request_method": "POST", "response_code": 302, "type": "request_response_middleware"},
    {"allow": "True", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "auth_status": "OK", "path": "/v2/o/authorize/3792648d-46e7-4705-9cb0-9ba18f0208fc/", "share_demographic_scopes": "True", "type": "Authorization", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "app_id": "1", "app_name": "TestApp"},
    {"app_id": "1", "app_name": "TestApp", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/v2/o/authorize/3792648d-46e7-4705-9cb0-9ba18f0208fc/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"},
    {"app_id": "1", "app_name": "TestApp", "auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "path": "/mymedicare/sls-callback", "request_method": "GET", "response_code": 302, "type": "request_response_middleware", "auth_path": "/v2/o/authorize/3792648d-46e7-4705-9cb0-9ba18f0208fc/"},
    {"auth_crosswalk_action": "C", "auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "type": "Authentication:success", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "type": "mymedicare_cb:get_and_update_user_initial_auth", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256645", "fhir_id_v3": "-591611569", "type": "mymedicare_cb:create_beneficiary_record", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-591611569", "type": "fhir.server.authentication.match_fhir_id", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "fhir_id_v2": "-10000010256645", "type": "fhir.server.authentication.match_fhir_id", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "patient search", "type": "fhir_auth_pre_fetch", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "v2/mymedicare/sls-callback", "sls_userinfo_status_code": 200, "type": "Authentication:start", "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b45e-ddb1-4195-8347-57c8c872fb42", "type": "SLSx_userinfo", "response_code": 403, "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "/v1/users/0854b45e-ddb1-4195-8347-57c8c872fb42", "type": "SLSx_userinfo", "response_code": 200, "app_id": "1", "app_name": "TestApp"},
    {"auth_require_demographic_scopes": "True", "path": "/sso/session", "type": "SLSx_token", "response_code": 200, "app_id": "1", "app_name": "TestApp"},
    {"app_id": "1", "app_name": "TestApp", "auth_require_demographic_scopes": "True", "path": "/mymedicare/login", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "1", "app_name": "TestApp", "auth_require_demographic_scopes": "True", "path": "/v2/o/authorize/", "request_method": "GET", "response_code": 302, "type": "request_response_middleware"},
    {"app_id": "", "app_name": "", "path": "/", "request_method": "GET", "response_code": 200, "type": "request_response_middleware"}
]

# We do not want to track metrics for the following app names
APP_NAMES_TO_IGNORE = ['TestApp', 'BlueButton Client (Test - Internal Use Only)', 'MyMedicare PROD', 'new-relic']

REQUEST_RESPONSE_MIDDLEWARE_TYPE = 'request_response_middleware'

AUDIT_EVENT_TYPES = [
    "Authentication:start",
    "Authentication:success",
    "Authorization",
    "AccessToken"
]

REDIRECT_STATUS_CODE = 302
OK_STATUS_CODE = 200


# Various helper functions
def to_int(value):
    if value == None:
        return None
    return int(value)

def is_real(fhir_id_v2, fhir_id_v3):
    """
    True if either FHIR ID indicates a real beneficiary.
    Mirrors SQL: fhir_id_v2 > 0 OR COALESCE(fhir_id_v3, 0) > 0
    """
    v2 = to_int(fhir_id_v2)
    v3 = to_int(fhir_id_v3)
    v3_val = v3 if v3 != None else 0
    return (v2 != None and v2 > 0) or v3_val > 0

def response_code_equals(response_code, val):
    return to_int(response_code) == val

def response_code_not_equals(response_code, val):
    return to_int(response_code) != val

def response_code_in(response_code, vals):
    return to_int(response_code) in vals

def response_code_not_in(response_code, vals):
    return to_int(response_code) not in vals

def response_code_range(response_code, low, high):
    rc = to_int(response_code)
    if rc == None:
        return False
    return (low <= rc and rc < high)

def response_code_gte(response_code, low):
    rc = to_int(response_code)
    if rc == None:
        return False
    return rc >= low

def path_matches_versioned_token(path):
    for v in ["v1", "v2", "v3"]:
        if path.startswith("/" + v + "/o/token"):
            return True
    return False

def path_matches_authorize_prefix(path):
    for v in ["v1", "v2", "v3"]:
        if path.startswith("/" + v + "/o/authorize/"):
            return True
    return False

def path_matches_authorize_full(path):
    for v in ["v1", "v2", "v3"]:
        prefix = "/" + v + "/o/authorize/"
        if path.startswith(prefix) and path.endswith("/") and len(path) > len(prefix):
            return True
    return False

def evaluate_request_response_metrics(tags):
    # Specifically analyze logs with type = 'request_response_middleware'
    matched = []

    path = tags.get("path") or ""
    request_method = tags.get("request_method") or ""
    response_code = tags.get("response_code")
    fhir_id_v2 = tags.get("fhir_id_v2")
    fhir_id_v3 = tags.get("fhir_id_v3")
    lastupdated = tags.get("req_qparam_lastupdated") or ""
    auth_grant_type = tags.get("auth_grant_type") or ""
    sdk_header = tags.get("req_header_bluebutton_sdk") or ""

    # FHIR Resource call stats tracking
    # Top level conditional used in each real call, then check if real beneficiary or not, then check path
    if request_method == "GET" and response_code_equals(response_code, OK_STATUS_CODE):
        if is_real(fhir_id_v2, fhir_id_v3):
            # V1 Real Beneficiaries
            if path.startswith("/v1/fhir"):
                matched.append("app_fhir_v1_call_real_count")

            if path.startswith("/v1/fhir/ExplanationOfBenefit"):
                matched.append("app_fhir_v1_eob_call_real_count")

            if path.startswith("/v1/fhir/Coverage"):
                matched.append("app_fhir_v1_coverage_call_real_count")

            if path.startswith("/v1/fhir/Patient"):
                matched.append("app_fhir_v1_patient_call_real_count")

            if path.startswith("/v1/fhir/ExplanationOfBenefit") and lastupdated != "":
                matched.append("app_fhir_v1_eob_since_call_real_count")

            if path.startswith("/v1/fhir/Coverage") and lastupdated != "":
                matched.append("app_fhir_v1_coverage_since_call_real_count")

            # V2 Real Beneficiaries
            if path.startswith("/v2/fhir"):
                matched.append("app_fhir_v2_call_real_count")

            if path.startswith("/v2/fhir/ExplanationOfBenefit"):
                matched.append("app_fhir_v2_eob_call_real_count")

            if path.startswith("/v2/fhir/Coverage"):
                matched.append("app_fhir_v2_coverage_call_real_count")

            if path.startswith("/v2/fhir/Patient"):
                matched.append("app_fhir_v2_patient_call_real_count")

            if path.startswith("/v2/fhir/ExplanationOfBenefit") and lastupdated != "":
                matched.append("app_fhir_v2_eob_since_call_real_count")

            if path.startswith("/v2/fhir/Coverage") and lastupdated != "":
                matched.append("app_fhir_v2_coverage_since_call_real_count")

            # V3 Real Beneficiaries
            if path.startswith("/v3/fhir"):
                matched.append("app_fhir_v3_call_real_count")

            if path.startswith("/v3/fhir/ExplanationOfBenefit"):
                matched.append("app_fhir_v3_eob_call_real_count")

            if path.startswith("/v3/fhir/Coverage"):
                matched.append("app_fhir_v3_coverage_call_real_count")

            if path.startswith("/v3/fhir/Patient"):
                matched.append("app_fhir_v3_patient_call_real_count")

            if path.startswith("/v3/fhir/ExplanationOfBenefit") and lastupdated != "":
                matched.append("app_fhir_v3_eob_since_call_real_count")

            if path.startswith("/v3/fhir/Coverage") and lastupdated != "":
                matched.append("app_fhir_v3_coverage_since_call_real_count")

            if path.startswith("/v3/fhir/Patient/") and "insurance-card" in path:
                matched.append("app_fhir_v3_generate_insurance_card_call_real_count")

        else:
            # V1 Synthetic Beneficiaries
            if path.startswith("/v1/fhir"):
                matched.append("app_fhir_v1_call_synthetic_count")

            if path.startswith("/v1/fhir/ExplanationOfBenefit"):
                matched.append("app_fhir_v1_eob_call_synthetic_count")

            if path.startswith("/v1/fhir/Coverage"):
                matched.append("app_fhir_v1_coverage_call_synthetic_count")

            if path.startswith("/v1/fhir/Patient"):
                matched.append("app_fhir_v1_patient_call_synthetic_count")

            if path.startswith("/v1/fhir/ExplanationOfBenefit") and lastupdated != "":
                matched.append("app_fhir_v1_eob_since_call_synthetic_count")

            if path.startswith("/v1/fhir/Coverage") and lastupdated != "":
                matched.append("app_fhir_v1_coverage_since_call_synthetic_count")

            # V2 Synthetic Beneficiaries
            if path.startswith("/v2/fhir"):
                matched.append("app_fhir_v2_call_synthetic_count")

            if path.startswith("/v2/fhir/ExplanationOfBenefit"):
                matched.append("app_fhir_v2_eob_call_synthetic_count")

            if path.startswith("/v2/fhir/Coverage"):
                matched.append("app_fhir_v2_coverage_call_synthetic_count")

            if path.startswith("/v2/fhir/Patient"):
                matched.append("app_fhir_v2_patient_call_synthetic_count")

            if path.startswith("/v2/fhir/ExplanationOfBenefit") and lastupdated != "":
                matched.append("app_fhir_v2_eob_since_call_synthetic_count")

            if path.startswith("/v2/fhir/Coverage") and lastupdated != "":
                matched.append("app_fhir_v2_coverage_since_call_synthetic_count")

            # V3 Synthetic Beneficiaries
            if path.startswith("/v3/fhir"):
                matched.append("app_fhir_v3_call_synthetic_count")

            if path.startswith("/v3/fhir/ExplanationOfBenefit"):
                matched.append("app_fhir_v3_eob_call_synthetic_count")

            if path.startswith("/v3/fhir/Coverage"):
                matched.append("app_fhir_v3_coverage_call_synthetic_count")

            if path.startswith("/v3/fhir/Patient"):
                matched.append("app_fhir_v3_patient_call_synthetic_count")

            if path.startswith("/v3/fhir/ExplanationOfBenefit") and lastupdated != "":
                matched.append("app_fhir_v3_eob_since_call_synthetic_count")

            if path.startswith("/v3/fhir/Coverage") and lastupdated != "":
                matched.append("app_fhir_v3_coverage_since_call_synthetic_count")

            if path.startswith("/v3/fhir/Patient/") and "insurance-card" in path:
                matched.append("app_fhir_v3_generate_insurance_card_call_synthetic_count")

    if path.startswith("/v1/fhir/metadata") and request_method == "GET" and response_code_equals(response_code, OK_STATUS_CODE):
        matched.append("app_fhir_v1_metadata_call_count")

    if path.startswith("/v2/fhir/metadata") and request_method == "GET" and response_code_equals(response_code, OK_STATUS_CODE):
        matched.append("app_fhir_v2_metadata_call_count")

    if path.startswith("/v3/fhir/metadata") and request_method == "GET" and response_code_equals(response_code, OK_STATUS_CODE):
        matched.append("app_fhir_v3_metadata_call_count")

    # Token request stats
    if request_method == "POST" and path_matches_versioned_token(path) and auth_grant_type == "refresh_token":
        if response_code_range(response_code, OK_STATUS_CODE, 300):
            matched.append("app_token_refresh_response_2xx_count")
        elif response_code_range(response_code, 400, 500):
            matched.append("app_token_refresh_response_4xx_count")
        elif response_code_gte(response_code, 500):
            matched.append("app_token_refresh_response_5xx_count")

    if request_method == "POST" and path_matches_versioned_token(path) and auth_grant_type == "authorization_code":
        if response_code_range(response_code, OK_STATUS_CODE, 300):
            matched.append("app_token_authorization_code_2xx_count")
        elif response_code_range(response_code, 400, 500):
            matched.append("app_token_authorization_code_4xx_count")
        elif response_code_gte(response_code, 500):
            matched.append("app_token_authorization_code_5xx_count")

    # Auth flow stats
    if path_matches_authorize_prefix(path):
        matched.append("app_authorize_initial_count")

    if path == "/mymedicare/login" and response_code_equals(response_code, REDIRECT_STATUS_CODE):
        matched.append("app_medicare_login_redirect_ok_count")

    if path == "/mymedicare/login" and response_code_not_equals(response_code, REDIRECT_STATUS_CODE):
        matched.append("app_medicare_login_redirect_fail_count")

    if path == "/mymedicare/sls-callback" and response_code_equals(response_code, REDIRECT_STATUS_CODE):
        if is_real(fhir_id_v2, fhir_id_v3):
            matched.append("app_sls_callback_ok_real_count")
        else:
            matched.append("app_sls_callback_ok_synthetic_count")

    if path == "/mymedicare/sls-callback" and response_code_not_equals(response_code, REDIRECT_STATUS_CODE):
        matched.append("app_sls_callback_fail_count")

    if path_matches_authorize_full(path) and request_method == "GET" and response_code_in(response_code, [OK_STATUS_CODE, REDIRECT_STATUS_CODE]):
        if is_real(fhir_id_v2, fhir_id_v3):
            matched.append("app_approval_view_get_ok_real_count")
        else:
            matched.append("app_approval_view_get_ok_synthetic_count")

    if path_matches_authorize_full(path) and request_method == "GET" and response_code_not_in(response_code, [OK_STATUS_CODE, REDIRECT_STATUS_CODE]):
        matched.append("app_approval_view_get_fail_count")

    if path_matches_authorize_full(path) and request_method == "POST" and response_code_in(response_code, [OK_STATUS_CODE, REDIRECT_STATUS_CODE]):
        if is_real(fhir_id_v2, fhir_id_v3):
            matched.append("app_approval_view_post_ok_real_count")
        else:
            matched.append("app_approval_view_post_ok_synthetic_count")

    if path_matches_authorize_full(path) and request_method == "POST" and response_code_not_in(response_code, [OK_STATUS_CODE, REDIRECT_STATUS_CODE]):
        matched.append("app_approval_view_post_fail_count")

    if sdk_header == "python":
        matched.append("app_sdk_requests_python_count")

    if sdk_header == "node":
        matched.append("app_sdk_requests_node_count")

    return matched

def evaluate_audit_metrics(tags):
    # Specifically analyze log events with type equal to Authentication:start, Authentication:success, AccessToken or Authorization
    matched = []

    event_type = tags.get("type") or ""
    auth_status = tags.get("auth_status") or ""
    allow = tags.get("allow") or ""
    auth_grant_type = tags.get("auth_grant_type") or ""
    auth_action = tags.get("action") or ""
    auth_crosswalk_action = tags.get("auth_crosswalk_action") or ""
    auth_require_demographic_scopes = tags.get("auth_require_demographic_scopes") or ""
    share_demographic_scopes = tags.get("share_demographic_scopes") or ""
    sls_status = tags.get("sls_userinfo_status_code")
    path = tags.get("path") or ""

    # For Authorization/Authentication:success events
    crosswalk_fhir_id = tags.get("fhir_id_v2")
    crosswalk_fhir_id_v3 = tags.get("fhir_id_v3")

    if event_type == "Authorization":

        if auth_status == "OK" and allow == "True":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_ok_real_bene_count")
            else:
                matched.append("app_auth_ok_synthetic_bene_count")

        if auth_status == "FAIL":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_fail_or_deny_real_bene_count")
            else:
                matched.append("app_auth_fail_or_deny_synthetic_bene_count")


        if auth_status == "OK" and allow == "True" and auth_require_demographic_scopes == "True" and share_demographic_scopes == "True":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_demoscope_required_choice_sharing_real_bene_count")
            else:
                matched.append("app_auth_demoscope_required_choice_sharing_synthetic_bene_count")

        if auth_status == "OK" and allow == "True" and auth_require_demographic_scopes == "True" and share_demographic_scopes == "False":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_demoscope_required_choice_not_sharing_real_bene_count")
            else:
                matched.append("app_auth_demoscope_required_choice_not_sharing_synthetic_bene_count")

        if allow == "False" and auth_require_demographic_scopes == "True":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_demoscope_required_choice_deny_real_bene_count")
            else:
                matched.append("app_auth_demoscope_required_choice_deny_synthetic_bene_count")

        if auth_status == "OK" and allow == "True" and auth_require_demographic_scopes == "False":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_demoscope_not_required_not_sharing_real_bene_count")
            else:
                matched.append("app_auth_demoscope_not_required_not_sharing_synthetic_bene_count")

        if allow == "False" and auth_require_demographic_scopes == "False":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_auth_demoscope_not_required_deny_real_bene_count")
            else:
                matched.append("app_auth_demoscope_not_required_deny_synthetic_bene_count")

        if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3) and (path.startswith('/v1/o/authorize') or path.startswith('/v2/o/authorize')):
            matched.append("app_auth_v1_v2_user_clicks_connect_bene_count")

        if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3) and path.startswith('/v3/o/authorize'):
            matched.append("app_auth_v3_user_clicks_connect_bene_count")

    if event_type == "AccessToken":

        if auth_action == "authorized" and auth_grant_type == "refresh_token":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_token_refresh_for_real_bene_count")
            else:
                matched.append("app_token_refresh_for_synthetic_bene_count")

        if auth_grant_type == "authorization_code":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_token_authorization_code_for_real_bene_count")
            else:
                matched.append("app_token_authorization_code_for_synthetic_bene_count")

    if event_type == "Authentication:start":

        if response_code_equals(sls_status, OK_STATUS_CODE):
            matched.append("app_authentication_start_ok_count")
        else:
            matched.append("app_authentication_start_fail_count")

    if event_type == "Authentication:success":

        if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3) and (path == 'v1/mymedicare/sls-callback' or path == 'v2/mymedicare/sls-callback'):
            matched.append("app_auth_v1_v2_user_makes_it_to_permission_screen_bene_count")

        if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3) and path == 'v3/mymedicare/sls-callback':
            matched.append("app_auth_v3_user_makes_it_to_permission_screen_bene_count")

        if auth_crosswalk_action == "C":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_authentication_matched_new_bene_real_count")
            else:
                matched.append("app_authentication_matched_new_bene_synthetic_count")

        if auth_crosswalk_action == "R":
            if is_real(crosswalk_fhir_id, crosswalk_fhir_id_v3):
                matched.append("app_authentication_matched_returning_bene_real_count")
            else:
                matched.append("app_authentication_matched_returning_bene_synthetic_count")

    return matched

def parse_simple_json(s):
    """
    Parse a flat JSON object string into a Starlark dict.
    Starlark does not support while loops, so for loops are used instead
    """
    result = {}

    s = s.strip()
    if s.startswith("{"):
        s = s[1:]
    if s.endswith("}"):
        s = s[:-1]
    s = s.strip()

    if not s:
        return result

    n = len(s)
    pos = [0]  # Use a list so we can mutate it inside helpers

    def skip_whitespace():
        for _ in range(n):
            if pos[0] >= n or s[pos[0]] != " ":
                break
            pos[0] += 1

    def parse_string():
        pos[0] += 1  # skip opening quote
        val = ""
        for _ in range(n):
            if pos[0] >= n:
                break
            c = s[pos[0]]
            if c == '"':
                pos[0] += 1  # skip closing quote
                break
            if c == '\\' and pos[0] + 1 < n:
                pos[0] += 1
                val += s[pos[0]]
            else:
                val += c
            pos[0] += 1
        return val

    def parse_primitive():
        # Check for null, true, false
        if s[pos[0]:pos[0]+4] == "null":
            pos[0] += 4
            return None
        if s[pos[0]:pos[0]+4] == "true":
            pos[0] += 4
            return True
        if s[pos[0]:pos[0]+5] == "false":
            pos[0] += 5
            return False
        # Numeric — read until comma, space, or end
        num_str = ""
        for _ in range(n):
            if pos[0] >= n:
                break
            c = s[pos[0]]
            if c == "," or c == " " or c == "}":
                break
            num_str += c
            pos[0] += 1
        return num_str

    # Main parse loop — iterate up to n times (one per key-value pair max)
    for _ in range(n):
        if pos[0] >= n:
            break

        skip_whitespace()
        if pos[0] >= n:
            break

        # Expect a quoted key
        if s[pos[0]] != '"':
            break

        key = parse_string()

        # Skip whitespace and colon
        for _ in range(n):
            if pos[0] >= n:
                break
            c = s[pos[0]]
            if c != " " and c != ":":
                break
            pos[0] += 1

        if pos[0] >= n:
            break

        # Parse value
        if s[pos[0]] == '"':
            val = parse_string()
        else:
            val = parse_primitive()

        result[key] = val

        # Skip whitespace and comma
        for _ in range(n):
            if pos[0] >= n:
                break
            c = s[pos[0]]
            if c != " " and c != ",":
                break
            pos[0] += 1

    return result


def summarize():
    # Entry point for app-summary ETL

    # Dictionary to track event counts grouped by metric/app
    # For each event that applies to a specific metric for an app,
    # We increment the count for that metric/app key
    accumulator = {}

    # List that will contain the summary rows that will be returned/written to itslog_summary
    returning_summary_rows = []

    events = query('events', 'SELECT * FROM itslog_events')

    for event in events:
    # for event in EVENTS:
        # json_event = json.decode(event.value)
        event_dict = parse_simple_json(event.get('value'))

        app_name = event_dict.get('app_name')
        app_id = event_dict.get('app_id')

        if not app_name or not app_id:
            continue

        if event_dict.get('type') == REQUEST_RESPONSE_MIDDLEWARE_TYPE:
            for metric in evaluate_request_response_metrics(event_dict):
                key = (metric, app_name)
                accumulator[key] = accumulator.get(key, 0) + 1

        if event_dict.get('type') in AUDIT_EVENT_TYPES:
            for metric in evaluate_audit_metrics(event_dict):
                key = (metric, app_name)
                accumulator[key] = accumulator.get(key, 0) + 1

    for (metric_tag, tag), count in accumulator.items():
        returning_summary_rows.append({
            "operation": metric_tag,
            "tags": tag,
            "value": str(count),
            "count": 1,
        })
    for row in returning_summary_rows:
        print("ROW: ", row)

    return returning_summary_rows

summarize()
