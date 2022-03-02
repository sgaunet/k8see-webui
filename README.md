
# k8see (Kubernetes Events Exporter)

Kubernetes Events Exporter is a suit of three tools to export kubernertes events in an external database. The goal is to get events in an SQL DB to be able to analyze what happened.

The 3 tools are :

* [k8see-exporter](https://github.com/sgaunet/k8see-exporter) : Deployment inside the kubernetes cluster to export events in a redis stream
* [k8see-importer](https://github.com/sgaunet/k8see-importer) : Tool that read the redis stream to import events in a database (PostGreSQL)
* [k8see-webui](https://github.com/sgaunet/k8see-webui) : Web interface to query the database


**It's actually in development, not production ready**

# Install 

k8see-webui can be installed as a docker image only for now. [You can get kubernetes manifest here.](https://github.com/sgaunet/k8see-deploy)]

# Development

## CSS

The project is using tailwindcss. No need to update, very basic usage, but if you want to update this part, [you'll find more informations here.](docs/css/README.md)

