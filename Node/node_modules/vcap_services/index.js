'use strict';

/**
 * if VCAP_SERVICES exists or the instance name exists in the
 * environemnt, then it returns the credentials
 * for the last service that starts with 'name' or {} otherwise
 * If plan is specified it will return the credentials for
 * the service instance that match that plan or {} otherwise
 * @param  String name, service name
 * @param  String plan, service plan
 * @param  String iname, instance name
 * @return {Object} the service credentials or {} if
 * name is not found in VCAP_SERVICES or instance name
 * is set as an environmental variable. Env var must be
 * upper case.
 */

const getCredentials = function(name, plan, iname) {
  if (process.env.VCAP_SERVICES) {
    var services = JSON.parse(process.env.VCAP_SERVICES);
    for (var service_name in services) {
      if (service_name.indexOf(name) === 0) {
        for (var i = services[service_name].length - 1; i >= 0; i--) {
          var instance = services[service_name][i];
          if ((!plan || plan === instance.plan) && (!iname || iname === instance.name))
            return instance.credentials || {};
        }
      }
    }
  }
  //Check if env vars were set directly
  var env = process.env;
  var instance = {};

  if(iname) {
    iname = iname.toUpperCase().replace(/[\s&-]/g, '_');
    if(env[iname]) {
      try {
        instance = JSON.parse(env[iname]);
      } catch(e) {
        console.warn('Error parsing JSON from process.env.' + iname );
        console.warn(e);
      }
    }
  }
  return instance;
};

/**
 * Returns the credentials that match the service label
 * pass credentials in the following format:
 * {
 *  "watson_conversation_username": "username",
 *  "watson_conversation_password": "password",
 * }
 * @param {string} serviceLabel The service label
 * @param {object} credentials The credentials from starterkit
 */
const getCredentialsFromLocalConfig = function(serviceLabel, credentials) {
  const creds = {};
  const key = `watson_${serviceLabel}_`;
  if(credentials) {
    Object.keys(credentials)
    .filter(c => c.indexOf(key) === 0)
    .forEach(k => {
      if (k.substr(key.length) === 'apikey') {
        creds[`iam_${k.substr(key.length)}`] = credentials[k]
      }
      else {
        creds[k.substr(key.length)] = credentials[k]
      }
    });
  }
  return creds;
}

/**
 * Returns all the credentials that match the service label from env variables
 *
 * @param {string} serviceLabel The service label
 * @param {object} credentialsFromFile OPTIONAL: The credentials for starterkit from local file
 */
const getCredentialsForStarter = function(serviceLabel, credsFromFile) {
  let creds = {};
  if (credsFromFile) {
    creds = getCredentialsFromLocalConfig(serviceLabel, credsFromFile);
  }
  else if (process.env.VCAP_SERVICES) {
    creds = getCredentials(serviceLabel);
  }
  else if (process.env[`service_watson_${serviceLabel}`]){
    creds = JSON.parse(process.env[`service_watson_${serviceLabel}`]);
  }
  return creds;
}

module.exports = {
  getCredentials,
  getCredentialsFromLocalConfig,
  getCredentialsForStarter,
};
