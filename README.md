# Mongo Logs to Kibana

This tool helps to visualise any Mongo logs into Kibana. It can very helpful to perform a deep investigation.  
It is responsible to parse any Mongo logs and populate ElasticSearch. 
This tool can handle Mongo logs of Gbs.   
Once populated, you can visualise the logs and perform some queries in Kibana.  

This tool can populate ElasticSearch running locally or from your server.  

## Technical Overview
- Golang
- ElasticSearch / Kibana
- Docker / Docker-compose

## Requirements
- Docker
- Make

## Commands
> These commands have been tested only on MacOS.  

##### How to run ElasticSearch / Kibana locally
```
cd elk
docker-compose up
```
> The url for ElasticSearch will be http://localhost:9200  
> The url for Kibana will be http://localhost:5601  

##### How to build the tool
```
make build
```

##### How to run the tool

```
ELASTICSEARCH_INDEX=your_elastic_search_index \
ELASTICSEARCH_URL=elastic_search_url \
LOG_FILE=your_log_file_path \
make run
```


### Visualise the Mongo Logs in Kibana 

1. Open Kibana (if local, open http://localhost:5601)
2. Create an index pattern in Kibana
3. Visualise and query