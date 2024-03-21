# Generate Doc Page for Gofakeit Functions

## Step 1: Retrieve JSON of All Functions

The `gofakeit` library includes a command line tool called `gofakeitserver`. It is an HTTP server for the tool and has an endpoint which can return a JSON of all the available functions.

First, we need to install the `gofakeitserver` command line tool. **Make sure you're in the root directory of the project.**

```shell
export GOBIN=$(pwd)/bin
mkdir -p $GOBIN
go install github.com/brianvoe/gofakeit/v7/cmd/gofakeitserver@latest
```

Then, we need to start the server.

```shell
${GOBIN}/gofakeitserver &
```

You should see the following output:

```text
Running on port 8080
```

Now, we can retrieve the JSON of all the available functions.

```shell
pushd docs/gofakeit-functions
curl -o gofakeit-functions.json http://localhost:8080/list
```

It will create a file called `gofakeit-functions.json` in the `docs/gofakeit-functions` directory.

## Step 2: Convert JSON to Markdown

We can run the Python script to generate a markdown file from the JSON. **You should have Python 3 installed (no dependencies required; tested on Python 3.11).**

```shell
python convert-json-to-md.py
```

It will create a file called `gofakeit-functions.md` in the `docs/gofakeit-functions` directory.

## Step 3: Clean Up

```shell
popd
kill $(lsof -t -i :8080)
```