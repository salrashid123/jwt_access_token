var log4js = require("log4js");
var logger = log4js.getLogger();

const {GoogleAuth, JWTAccess} =  require('google-auth-library');
const {PubSub} = require('@google-cloud/pubsub');

jkey = require("../certs/jwt-access-svc-account.json");
const projectId = 'core-eso';

const client = new JWTAccess(
	jkey.client_email, 
	jkey.private_key,
	jkey.private_key_id
 );
client.useJWTAccessWithScope = true;

const pubsub = new PubSub({
	credentials: client,
	projectId: projectId
});
pubsub.getTopics((err, topic) => {
	if (err) {
		console.log(err);
		return;
	}
	topic.forEach(function(entry) {
    logger.info(entry.name);
	});
});


