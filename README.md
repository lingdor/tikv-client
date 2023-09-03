# tikv-client
tikv-client is simple command line program


#install

download in page:https://github.com/lingdor/tikv-client/releases

#demo


```shell
./tikv --pd 127.0.0.1:2379
```



```shell
put xx "good\"123"
```
output:
done!



```shell
get key1 xx
```
output:

key1:
xx: good\"123

```shell
set names utf8
set names gbk
```

```shell
rawget key1
```

```shell
cat aa.bmp | rawset aa.jpg
```

