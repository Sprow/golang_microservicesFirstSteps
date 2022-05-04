To run all services in docker use:
`docker-compose up`


content_task microservice wait for task(url) from user and sand it to rabbitmq.

content_scraper get task(url) in rabbitmq and scrap all site content.
Save site content in mongodb and sent new task (oid) to parser_direct Queue

content_parser get oid in RabbitMQ --> make api call to content_scraper
to get scraped content --> parse --> save result in mongo