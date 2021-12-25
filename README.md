# eatigo
## How to use: 

Build docker image
```
sh build.sh
```

Running up:
```
docker-compose up
```

Call the API
```
curl --location --request GET 'localhost:8080/eatigo/v1/restaurants?place=Sukhumvit%20Road%20Bangkok'
```

Call the API again with next page
```
curl --location --request GET 'localhost:8080/eatigo/v1/restaurants?place=Sukhumvit%20Road%20Bangkok&cursor=Aap_uED8Z-WzL1tRU2FbDZfwFBpDGTYMbanXKUUWR_MHmYqqg058ws9HAxaPLETdiHCN2hxZ-WxMU8GszdPQfF7seOSPN2o2ceqdchNDsaT0IS3zUw-7K2E_dvethMpZ0oAXsxD0O1tppYIYsEJ4qSr26PjSeUl7DWZcgelFo5XRHN1CBPxIpaiGGgbazAhBmxqm5zL5HNO1d9P9irrzQgAmgcVEqIeuE6_bIE3SuJ6aEjPZnSh7gG1Llu4xPtDYFSZZpDq1MfALnZkkHFk1DMOkhRY7GkYzuTIIm7_VFcvxtgyCLV5Go8IYOoJzfAf3HFpszXC_Q5nKvcGeVAvTW4Pv5ZIZto712JEgLNF7ff_TaW1eJRH7hPrqs8MF94ki1watAw6ZdXSt4byxHDpaUR-A1_-5KDqjoYAbaNWb3fteknIk6c0zglh5tido6O_mf4c'
```