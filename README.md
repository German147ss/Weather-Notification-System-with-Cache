# meli-challenge

## Descripción
Tengo una API de usuarios que permite guardar los datos del usuario basándose en la solicitud del cliente. Además, al momento del registro, devuelve la información meteorológica de la ubicación del usuario por primera vez. La API también permite programar notificaciones diarias que se enviarán a través de una cola de mensajes.

La API de clima incorpora un manejo de caché para reducir la carga en la API externa y evitar solicitudes duplicadas, optimizando así el tiempo de respuesta y el uso de recursos.

Finalmente, tenemos un servicio de notificaciones que ejecuta un cron cada minuto. Este servicio verifica en la tabla de preferencias de usuario si hay notificaciones programadas. Si encuentra alguna, obtiene la información meteorológica actual y publica la notificación en el canal de la cola de mensajes.

