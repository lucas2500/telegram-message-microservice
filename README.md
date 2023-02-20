
# Telegram Message Microservice

Um microserviço para envio de mensagens ao Telegram escrito em Golang (Fiber Framework) e RabbitMQ.

# Configuração do .env
**API_HTTP_PORT**: Porta HTTP da API\
**RABBITMQ_DEFAULT_USER**: Usuário do RabbitMQ\
**RABBITMQ_DEFAULT_PASS**: Senha do RabbitMQ\
**RABBITMQ_DEFAULT_HOST**: Endereço do host do RabbitMQ\
**RABBITMQ_DEFAULT_VHOST**: Virtual host do RabbitMQ\
**RABBITMQ_DEFAULT_PORT**: Porta do RabbitMQ\
**RABBITMQ_QUEUE_NAME**: Nome da fila\
**RABBITMQ_EXCHANGE_NAME**: Nome da exchange\
**RABBITMQ_QUEUE_ROUTING_KEY**: Nome da routing key\
**TELEGRAM_BASE_URL**: Base URL da API do Telegram\
**TELEGRAM_ROUTE**: Rota da API do Telegram para envio de mensagens\
**WORKERS_NUMBER**: Total de workers que serão instanciados para processar as mensagens

#

## Fluxo do serviço

![App Screenshot](https://raw.githubusercontent.com/lucas2500/telegram-message-microservice/master/Fluxo%20do%20servi%C3%A7o.png)

