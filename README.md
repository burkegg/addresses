# addresses

This project serves some basic property info to the browser.  Type in part of
a street address (not the city) to see results.


## Usage

### Requirements to run:
1. Internet connection.
1. Docker / Docker Compose

### Instructions to run:
1.  Clone down the repo with `git clone git@github.com:burkegg/addresses.git` or `git clone https://github.com/burkegg/addresses.git`
depending on whether you're using ssh or not.
1.  Run the command: `docker-compose up`.
1.  Navigate to `localhost:8080`
1.  Start typing an address.  Note - the search only works for the street address, not the city/state.

## Tools / Thought process

The CLI is Cobra/Viper which lets us pass flags in easily to deploy in different
environments.  To get things going quicker, I used `dbt`, Dynamic Binary Toolkit
to set up the boilerplate.  In the words of the dbt author: "You could
set up all that boilerplate manually, but WHY?"

To see which flags are available, go to:  `/addresses/cmd/run.go`.  The flags can be passed
in through environment variables or on startup.  Viper also lets you use a config file.  The idea
here is to be able to pass off an app to the infra folks who will do some magic with k8s or `waves hands` *hey look over there!*

The server logic is all found at `addresses/pkg/addresses/addresses.go`. The app is written in Go using Gorm to handle postgres interactions and Gin to handle
routing requests.  Gorm is an ORM that is nice for a small app like this.  You don't even have
to build the tables yourself - GORM will do it for you based on your types.  Gin is a an http library that
I'm used to - that's the only reasoning there.

The server has one endpoint `/addresses` which I set up as a `POST` request.
It's really a `GET`, but I had some other ideas of things I wanted to
include, and sometimes it's easier to just pass in a JSON object than to
append a bunch of stuff to the request.  I didn't end up doing that other stuff,
and didn't get back to change it to a GET.

On startup, the server reads in the csv file found at `addresses/pkg/assets/addressdata.csv`.  This is a csv of data from Redfin with some basic info about houses for sale in San Francisco.

The logic of searching for houses is handled by Postgres with `ILIKE`.  For example
`where ILIKE %searchterm%`, which matches any number of characters at the beginning or end of the search.  Thanks Postgres.  :)

The web client is found at `addresses/pkg/addresses/assets/js/App.jsx`.  This
is an embedded file in the go application.  I don't know how to use webpack or Create-React-App
in this manner - but with a simple UI like this, I feel like this is an okay approach.
We only need a couple components, and the logic is all pretty simple.  You only NEED one here
but I put in one to demonstrate passing props to a functional component.  The downside is that testing the
client probably isn't going to happen here - but if you're building a UI like this it's probably going to stay pretty simple.

The styling is all very simple and done in-line.  It's not optimal for rendering speed apparently,
but I haven't worked on a website fast enough for that amount of time to matter.  Plus here everything
is pretty quick and simple.
