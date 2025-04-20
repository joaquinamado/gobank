# GOBANK 

This project was made following Anthony GG 
[tutorials](https://www.youtube.com/watch?v=pwZuNmAzaH8) as a starting 
point. 

After finishing that I continued working on the project adding different 
features: 


## Roadmap 

[X] Change the project architecture to have one more layer  
  +  API layer: here we check the received input, check if there are any 
  errors, call the Repository layer & return the HTTP response (& body 
  if it corresponds) 
  + Repository layer: Where the buissnes logic is located 
  + DB layer: SQL queries.
[X] Integrate Swagger in order to make the documentation of the API 
 easier.
[X] Add Air to have hot reload
[X] Create the transfer logic GET (paginated) & POST
[X] Add an invoice generation endpoint 
[ ] Create accounts by uploading a CSV (make it as parallel as possible)  
[ ] WebSocket for account balance? Real time ?
[ ] Message queue for transfer (like notifications) ? 
[ ] Frontend (HTMX?)
[ ] Improve auth, reset token, roles?, close sesion from one device?

