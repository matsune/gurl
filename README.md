
# gurl
gurl is a simple CLI tool that can transfer http/https data.  
Since gurl aims to be able to use intuitively, there are less options provided than conventional tools, like [curl](https://github.com/curl/curl).  
gurl has an interactive mode and you can construct requests interactively on the CUI.

## Usage
### One-liner mode
```shell
$ gurl METHOD URL [OPTIONS]
```

```shell
$ gurl get localhost:8080

> Status
200 OK

> Header
Content-Type: application/json; charset=utf-8
Date: Mon, 24 Dec 2018 12:21:02 GMT
Content-Length: 15

> Body
{
  "status": "ok"
}
```

#### Request options

|short|long|description| example |
|---|---|---|---|
|-u| --user= | Basic auth <user[:password]> |`$ gurl get localhost:8080 -u user:password`|
|-H| --header= | Extra header \<key:value> |`$ gurl get localhost:8080 -H Custom-Key:value`|
|-j| --json= | JSON data| `$ gurl post localhost:8080/json -j '{"user":"gurl", "password":"pass"}'`|
|-x| --xml= | XML data|`$ gurl post localhost:8080/xml -x '<map><user>gurl</user><password>pass</password></map>'`|
|-f| --form= | Form URL Encoded data \<key:value> |`$ gurl post localhost:8080/form_post -f user:gurl -f password:pass`|



### Interactive mode
```shell
$ gurl
# or 
$ gurl -i
```
![interactive mode](https://user-images.githubusercontent.com/12775019/50400302-0da13380-07c9-11e9-9b71-878070083a2b.gif)

You can skip the input prompt by passing the value with the `-i` flag in interactive mode.
```shell
$ gurl localhost:8080 -i
# URL input prompt will be skipped
```

By passing the `-o` flag when running in interactive mode, you can output the constructed request in the one-liner format at the end

```shell
$ ./gurl -i -o
? Choose method: POST
? URL http://localhost:8080/json
? Use Basic Authorization? No
? Add custom header? No
? Body: JSON

> Status
200 OK

...

> one-liners
./gurl POST http://localhost:8080/json -j '{"user":"gurl","password":"pass"}'
```
