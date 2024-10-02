# meli-challenge
- Se da de alta el usuario.
Se carga su ciudad y se obtiene su prevision del tiempo -> la busca en el servicio de clima, este lo devuelve
y se muestra.
Se dispara un job para enviar notificaciones.

- Se da de baja el usuario para las notificaciones.


se da de alta en el job que busca cada hora para notificarle a los usuarios.


Por otro lado el servicio de clima estara pendiente de cada solicitud para devolver el clima segun la ciudad, donde primero se solicita si esta cacheada
si no esta cacheada se busca cptec, se devuelve y se cachea.
El cache dura una hora, cada hora se realiza la solicutd al cptec.


