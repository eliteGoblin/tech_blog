
#### ENOSPC Error


```
npm dedupe
# If not work, TRY
# This will increase the limit for the number of files you can watch
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
```


#### Clean && Delete

*  delete blog file in source/_post/xxx.md xxx.html xxx/_
*  delete db.json file directly
*  hexo clean

