[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/vedicsociety/brucheion-pro/tree/master)

<img src="files/static/img/logo-flat.png" alt="" width="500">

# Brucheion, the Virtual Research Environment

Brucheion is a Virtual Research Environment (VRE) to create Linked Open Data (LOD) for historical languages and the research of historical objects in multiuser environment. This project is a combination of Go (Golang) and Svelte, two powerful technologies for building high-performance, scalable web applications. Golang handles the backend, offering its characteristic speed and efficiency, while Svelte is used on the frontend for its ability to compile your code to efficient, imperative code that directly manipulates the DOM. All transcript data stored in Postgres DB.

# Overview
## Project Structure
### Backend
* The backend is written in Golang
* Using an extended [version](https://github.com/vedicsociety/platform) of the framework built in "Pro Go The Complete Guide to Programming Reliable and Efficient Software Using Golang (2022) Adam Freeman" 
* the entry point of the code is from main.go
* config.json is the source for configuration
* data is stored in Postgres and accessed via sql in ./sql, migrations are managed from ./sql/migrations

### Frontend
* files/app hosts the frontend npm project 
* Svelte is used to create the UI
* Libraries like bootstrap, opensea dragon etc are used as needed

### Sources
* Refer to "Pro Go The Complete Guide to Programming Reliable and Efficient Software Using Golang (2022) Adam Freeman" for more details
  
## File Structure
```
|-- auth        //backend code powering the Authorization
|-- files       //this entire directory is served statically by the backend, i.e. any file in this can be accessed via http
|   |-- app         //the frontend project using Svelte
|   |   |-- dist      //the packaged javascript is distributed from here
|   |   |-- src       //the source of the javascript files
|   |   |   |-- components    //Svelte kit components
|   |   |   |-- lib           //js Libraries like opensea dragon
|   |   |   |-- routes        //All routes exposed by frontend
|   |   |   `-- transitions   //helper js for transitions
|   |   `-- test        //frontend tests
|   |-- css           //css files for static distribution
|   |   `-- images    //image files used by CSS
|   |-- fonts         //fonts
|   |-- img           //images
|   |-- js            //static js libraries
|   `-- static        //static (TODO: merge with above)
|-- gocite          //go package for the cite format      
|-- handlers        //go package for different endpoints exposed
|   |-- admin       // /admin endpoints
|   |-- api         // /api endpoints
|   |-- root        // / endpoints
|   `-- newauth     // /auth endpoints for Auth0 authenticate handles
|-- models          // Data structs
|   `-- repo        // go package exposing DB interaction
|-- scripts         // scripts used for deployment
|-- sql             // all the db queries in .sql files
|   `-- migrations    //db migrations 
|-- templates       //go templates used by backend
|   `-- email       //templates used for email
`-- utils           //go util package

```



## Features

   * user management
   * import collection from .cex file
   * manage collections (create, delete, share, etc)
   * edit image references
   * manage image collection
   * support three types of images: IIIf, static, DZI
  
# Setup & Installation
### Local
* Verify the config.json file to see external dependencies like DB are correctly configured
* `make dev` will compile the frontend and run the backend


### Deploy to heroku
* Use the deploy to heroku on the top of the README to deploy to heroku
* the config in config.json can be overriden using environment vairbales
  * if the config.json has
```
"system" : {
        "prefix": "BR_",
},
"logging" : {
  "level": "debug"
},
```
  an environment variable `BR_logging_level=info` will override the config.json
* Using the above override any needed configuration


* all parameters load from two places: 
  ** config.json file  
  ** os environment with prefix, defined in config.json (system:prefix) 
* paremeters, defined in os environment redefine values in config.json
* new env can be added to os.env to leaf level only 
* for redefine a variable, you must specify all path. for example:

# Development
## Handlers
* By agreement, handlers are opened  for unauthenticated users via api/v1 URL, handlers via api/v2 URL are opened for autenticated users only 
* Permissions for handler's access defined in main.go

## DB
* All interactions to DB are done via sql scripts in sql directory
* These are exposed as methods in models/repository.go. To do so
  * add the new sql file 
  * add new entry in model/repo/sql_repo.go's SqlCommands
  * update the config.json to expose the sql script via sql.commands 
  * expose a new method in models/repository.go's Repository interface
  * implement the newly introduced method in model/repo via a new or exisitng go file
* All migrations are in sql/migrations

## Frontend
* Some parts of frontend are wrotten using Svelte framework
* For compile Svelte modules use command: cd files/app/ && npm run dev (or npm run build)

# Contributing

We welcome contributions from the community. If you'd like to contribute, please fork the repository and make changes as you'd like. Pull requests are warmly welcome.
# License

This project is licensed under the MIT License - see the LICENSE file for details.

