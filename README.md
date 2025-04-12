# GOBANK 

This project was made following Anthony GG 
(tutorials)[https://www.youtube.com/watch?v=pwZuNmAzaH8] as a starting 
point. 

After finishing that I continued working on the project adding different 
features: 

**Finished:**

- Changed the project architecture to have one more layer 
  +  API layer: here we check the received input, check if there are any 
  errors, call the Repository layer & return the HTTP response (& body 
  if it corresponds) 
  + Repository layer: Where the buissnes logic is located 
  + DB layer: SQL queries.

- Integrated Swagger in order to make the documentation of the API 
 easier.
- Added Air to have hot reload

**Pending:**

- Create the transfer logic
- Add an invoice generation endpoint 
- Create accounts by uploading a CSV 
