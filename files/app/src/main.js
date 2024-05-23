import * as AWS from '@aws-sdk/client-s3'
import App from './App.svelte'
import 'unfetch'
AWS.config.region = 'eu-west-1' // Region
AWS.config.credentials = new AWS.CognitoIdentityCredentials({
  IdentityPoolId: 'eu-west-1:7b929dbf-6831-4a5c-977c-e4ca4ddc7d8b',
})

const app = new App({
  target: document.getElementById('app'),
})

export default app
