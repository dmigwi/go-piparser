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
- For unix based o.s.

```bash
    mkdir -p ~/playground
```
- For windows
```bash
    md %USERPROFILE%\playground
```

## Start Server

```bash
    cd go-piparser/testutil

    $ go build . && ./testutil
    2019/03/05 11:44:09 Setting up the environment. Please Wait...
    2019/03/05 11:44:19 Serving on http://127.0.0.1:8080

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

- The tool may take a couple of seconds. Please wait ...


![Screenshot from 2019-03-13 19-19-14](https://user-images.githubusercontent.com/22055953/54296360-d3745080-45c5-11e9-89d4-d903d5acc0fc.png)
