Simple logger that takes result of REST api calll and writes it into end of file.

Another possible consideration will be use of bash scrip (i do believe its doable)
but i suggest purpose of exercise is to see ability to access REST api from Go.

Logger is configured through CMD line, takes host, port where to connect and file where 
write to. Chosen approach was to make it as simple as possible.
Logger logs time of execution in UTF, errors if any happens and results in simple table.

With more input it will be clearer what are expectations from logger.


 