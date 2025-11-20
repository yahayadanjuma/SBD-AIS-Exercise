Software Architecture for Big Data - Exercise AIS WS 2025

During the next few exercises you are going learn how to build webservice Golang.

We're going to have a look at state-of-the-art server technologies, facilitate databases and communication protocols
such as gRPC, REST. This whole application will be containerized and can be deployed to a cluster.

This repository contains numbered folders, with every folder corresponding to one exercise (Exc_1-9).
Some exercises may provide a skeleton/ folder to get you started with your assignment.

Modus Operandi
1. Fork this repository
2. Complete the exercise
3. Put the completed exercise in a solution/ folder, i.e. Exc_2/solution
4. Push the solution to your GitHub repository

If you have any questions, feel free to contact me via my university e-mail, or start a GitHub discussion.

Local Run Notes (my quick notes)

These are my personal notes from getting the exercise running on my machine. I'm still learning, so
this is written like a short cheat sheet I used.

- Start everything from the Exc_5/skeleton folder with:

    docker compose up -d --build

- Open the frontend in your browser at http://localhost.
- The frontend will also respond at http://orders.localhost (I configured Traefik to make both work).
- The API lives at http://orders.localhost/api/* (example: http://orders.localhost/api/menu).

- If you open the API URL in a browser you might get redirected to the UI — I added that so
  it looks nicer when you open the API in a browser. API clients that request JSON still get JSON.

- If orders.localhost doesn't work on Windows, add this line to your hosts file (edit as Admin):

    127.0.0.1 orders.localhost

- Traefik dashboard (local only): http://localhost:8080/dashboard/ (note: this is insecure, only
  for local debugging).

- Database: the Postgres service uses a named volume order_pg_vol. If Postgres fails to start with
  data format errors (happens when changing major Postgres versions), remove the volume and recreate it:

    docker compose down -v
    docker compose up -d --build

- If something looks wrong: try a browser hard-refresh or docker compose down -v then docker compose up -d --build.

These notes helped me while I was testing — feel free to edit or move them into a separate docs/ file.
# Software Architecture for Big Data - Exercise AIS WS 2025

During the next few exercises you are going learn how to build webservice Golang.

We're going to have a look at state-of-the-art server technologies, facilitate databases and communication protocols
such as gRPC, REST. This whole application will be containerized and can be deployed to a cluster.

This repository contains numbered folders, with every folder corresponding to one exercise (Exc_1-9).
Some exercises may provide a `skeleton/` folder to get you started with your assignment. 

## Modus Operandi
1. Fork this repository
2. Complete the exercise
3. Put the completed exercise in a `solution/` folder, i.e. `Exc_2/solution`
4. Push the solution to your GitHub repository


If you have any questions, feel free to contact me via my university e-mail, or start a GitHub discussion.

## Local Run Notes (my quick notes)

These are my personal notes from getting the exercise running on my machine. I'm still learning, so
this is written like a short cheat sheet I used.

- Start everything from the Exc_5/skeleton folder with:


docker compose up -d --build


- Open the frontend in your browser at `http://localhost`.
- The frontend will also respond at `http://orders.localhost` (I configured Traefik to make both work).
- The API lives at http://orders.localhost/api/ (example: `http://orders.localhost/api/menu).

- If you open the API URL in a browser you might get redirected to the UI — I added that s
	it looks nicer when you open the API in a browser. API clients that request JSON still get JSON.

- If orders.localhost doesn't work on Windows, add this line to your hosts file (edit as Admin):


127.0.0.1 orders.localhost


- Traefik dashboard (local only): `http://localhost:8080/dashboard/` (note: this is insecure, only
	for local debugging).

- Database: the Postgres service uses a named volume `order_pg_vol`. If Postgres fails to start with
	data format errors (happens when changing major Postgres versions), remove the volume and recreate it:


docker compose down -v
docker compose up -d --build


- If something looks wrong: try a browser hard-refresh or `docker compose down -v` then `docker compose up -d --build`.

These notes helped me while I was testing — feel free to edit or move them into a separate `docs/` file.