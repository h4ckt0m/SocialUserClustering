from linkedin_api import Linkedin
import json

# Authenticate using any Linkedin account credentials
api = Linkedin('user@gmail.com', 'enter_password')
enterprise = api.get_company('companyname')
json_object = json.dumps(enterprise, indent= 4)
print(json_object)


# GET a profile
profile = api.get_profile('user_profile')

# GET a profiles contact info
contact_info = api.get_profile_contact_info('user_profile')

# GET 1st degree connections of a given profile
connections = api.get_profile_connections('1234asc12304')

# GET profile views
search_companies = api.search_companies('')

# GET company profile connections
prof_conn = api.get_profile_connections('')
