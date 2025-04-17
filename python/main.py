project='core-eso'

from google.auth import jwt
from google.cloud import pubsub_v1
import google.auth
from google.oauth2 import service_account

# audience = 'https://pubsub.googleapis.com/google.pubsub.v1.Publisher'
# credentials = jwt.Credentials.from_service_account_file('../certs/jwt-access-svc-account.json',audience=audience)

credentials =  service_account.Credentials.from_service_account_file('../certs/jwt-access-svc-account.json').with_always_use_jwt_access(True)

publisher = pubsub_v1.PublisherClient(credentials=credentials)
project_path = f"projects/{project}"
for topic in publisher.list_topics(request={"project": project_path}):
  print(topic.name)

