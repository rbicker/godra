godra
=====
* godra is a simple login / logout provider for ory hydra, written in go.
* Consent is automatically given.
* It is provided "as is".

# basic configuration
The application is configured using environment variables (default value in bracket):
* **HYDRA_PRIVATE_URL**: hydra's private url (http://localhost:4445)
* **PORT**: godra server's http port (5000)
* **MONGO_URL**: mongodb server url (mongodb://localhost:27017)
* **MONGO_DB**: name of the mongodb database (db)
* **MONGO_COLLECTION**: name of the mongodb collection (users)

# customize login page
It is possible to provide a custom html-header or -footer by providing the path to html files as env variables:
  * **CUSTOM_HEADER_PATH**
  * **CUSTOM_FOOTER_PATH**

* If you need to serve an additional directory containing static files (your logo for example), you can do so by setting **CUSTOM_STATIC_PATH**. The folder will be served under **/static**.
* If you want to serve an additional stylesheet, you can do so by setting **CUSTOM_STYLESHEET_PATH**.