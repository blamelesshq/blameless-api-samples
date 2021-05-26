## Nodejs example code to ingest SLI Data into Blameless API

To run this code example, first install all dependencies:
```javascript
$ yarn install

// or with NPM
$ npm install
```

Then run:
```javascript
$ yarn start

// or with NPM
$ npm start
```


In order to get the auth token to make queries to the Blameless API, you will need to set some variables in your `.env` file:

```
AUTHZERO_CLIENT_ID=<AUTHZERO_CLIENT_ID>
AUTHZERO_CLIENT_SECRET=<AUTHZERO_CLIENT_SECRET>
AUTHZERO_API_AUDIENCE=<AUTH0_API_AUDIENCE>
BLAMELESS_HOST=<BLA
MELESS_HOST (i.e: blameless.blameless.io)>
```