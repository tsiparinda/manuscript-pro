[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/vedicsociety/brucheion-pro/tree/master)

<img src="files/static/img/logo-flat.png" alt="" width="500">

# Brucheion, the Virtual Research Environment

Brucheion is a Virtual Research Environment (VRE) to create Linked Open Data (LOD) for historical languages and the research of historical objects in multiuser environment. This project is a combination of Go (Golang) and Svelte, two powerful technologies for building high-performance, scalable web applications. Golang handles the backend, offering its characteristic speed and efficiency, while Svelte is used on the frontend for its ability to compile your code to efficient, imperative code that directly manipulates the DOM. All transcript data stored in Postgres DB.

# Overview
  

## Features

   * user management
   * import collection from .cex file
   * manage collections (create, delete, share, etc)
   * edit image references
   * manage image collection
   * support three types of images: IIIf, static, DZI
  
# Setup & Installation

To get started with this project, you can use button on the top of this page to deploy project to Heroku.



### Setup the config parameters

* all parameters load from two places: 
  ** config.json file  
  ** os environment with prefix, defined in config.json (system:prefix) 
* paremeters, defined in os environment redefine values in config.json
* new env can be added to os.env to leaf level only 
* for redefine a variable, you must specify all path. for example:

config.json variable sql:driverName = "postgres"

os.env should be: [prefix]_sql_driverName = "postgres"

where prefix is the config.json parameter system:prefix

* description of some impotant envs:
"templates":"reload"   reload templates be every rendering (true for developers)
"templates":"path"      path to templates
"files":"path"          path to static files
"sessions":"cyclekey"   reset session by restart application (true for developers)


# Contributing

We welcome contributions from the community. If you'd like to contribute, please fork the repository and make changes as you'd like. Pull requests are warmly welcome.
# License

This project is licensed under the MIT License - see the LICENSE file for details.

