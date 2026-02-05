it seems like bin.ts is the entry point of the program;

It starts the server based on what file is being sent into args()?

so at least with npm run dev, it starts with fixtures/db.json as the file, and I suppose it gets port, host and static from .env ? or the default set in the bin.ts file's try block?

If the bin file doesn't find the file it was supposed to use, it just exits the process.

otherwise - Check if the file is empty and set it up to be a json file if it is.

I presume adapter is the interface between the scripts and the database itself, which tells the computer how to handle and modify the data?

I also presume the observer is sort of like a logger - It observes and logs out what the adapter is doing to the database, which uses the lowdb npm package which is a json DB.
It awaits a read to populate the variable with the data from the db.

Once the db variable is setup, start the server throug the createapp function which takes in the db and some options.
here the function connects the db up to the server.ts script, creates a new express server/app?
Set ups middleware -> using sirv to serve static files (public folder), CORS and json body parser.
then the function sets up the different endpoints;

- get- / - uses eta to serve up the index.html file with the data from the db.
- get- /:name - Name is empty by default, but should be populated by the req.params - the endpoint itself creates a new query object, and populates the object based on req.query keys -> converting the keys into ints and checks if the number is not NaN and sends back the data based on the find function from service.ts, using the name params and the query.
- get- /:name/:id - Finds a specific entry in the database based on the name and id passed through req.params, sending the data fetched by the findById function from service.ts
- post- /:name - Checks if the body provided is an item suitable for the DB; If it is create an entry in the DB through the create function in service.ts
- put- /:name - Same as above, but Updates the entry entry based on the body provided with the put method.
- put- /:name/:id - Updates an entry with the put method based on the req.body provided.
- patch- /:name - Partial update of an entry based on name
- patch- /:name/:id - Partial update of an entry based on name and ID
- delete- /:name/:id - Deletes a specific entry based on name and ID.

The app then listens to the given port, printing out some information in the console.
at the end there seems to be some development stuff;
using the observer to flag whether we're writing to the DB. It does seem to be mainly logging out what is going on with the database?
