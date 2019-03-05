# Test Client

This is a http server that parses and serves data of the proposal token set on the URL path.


## Check the git version to confirm its installation

```bash
    git --version
```

## Clone this repository

```bash
    git clone https://github.com/dmigwi/go-piparser.git
```

## Create the tool clone directory

```bash
    mkdir -p ~/playground
```

## Start Server

```bash
    cd go-piparser/testclient

    $ go build . && ./testclient
    2019/03/05 11:44:09 Please Wait... Setting up the environment
    2019/03/05 11:44:09 Serving on 127.0.0.1:8080

```

### Get a proposal token

    e.g `27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50`

now access the URL `http://localhost:8080/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50` on the browser.


### Logs on the Terminal/CMD interface

```bash
    2019/03/05 11:44:29 Retrieving details for 27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50 ...
    2019/03/05 11:45:11 Processing charts data for 27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50 ...
    2019/03/05 11:45:11 Found Yes: 13206 No: 7985 and Total 21191 
    2019/03/05 11:45:11 Done.
```

### Charts will appear

- The tool may take a couple of seconds to clone (if no repo existed before). Please wait ...

- Consecutive proposal token queries are faster since data updates are done hourly.


![screenshot from 2019-03-05 11-58-00](https://user-images.githubusercontent.com/22055953/53793018-11c29d80-3f3e-11e9-911d-819a3e526f62.png)
