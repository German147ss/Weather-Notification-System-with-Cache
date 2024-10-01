1. Diseño de la Arquitectura
La solución se basará en microservicios con la capacidad de escalar y ser resiliente a fallos. Los componentes clave son:

API del Clima (CPTEC): Servicio que consumirá la API externa para obtener la previsión del tiempo y olas.
Gestión de Usuarios: Servicio para administrar las preferencias de los usuarios (ciudad seleccionada, agendamiento de notificaciones, opt-out).
Servicio de Notificaciones: Servicio que enviará las notificaciones en el horario agendado a los usuarios.
Cola de Mensajes (Opcional): Se podría utilizar una cola de mensajes como RabbitMQ o Kafka para manejar las notificaciones en un sistema más robusto.
Almacenamiento: Base de datos para persistir la información de usuarios, preferencias, y registros de notificaciones.