
#### Blog建立

*  hexo init
    -  切换源到淘宝　[切换淘宝npm镜像源](https://blog.skyx.in/archives/206/)
*  安装插件
    ```
    npm install hexo-generator-index --save
    npm install hexo-generator-archive --save
    npm install hexo-generator-category --save
    npm install hexo-generator-tag --save
    npm install hexo-server --save
    npm install hexo-deployer-git --save
    npm install hexo-deployer-heroku --save
    npm install hexo-deployer-rsync --save
    npm install hexo-deployer-openshift --save
    npm install hexo-renderer-marked@0.2 --save
    npm install hexo-renderer-stylus@0.2 --save
    npm install hexo-generator-feed@1 --save
    npm install hexo-generator-sitemap@1 --save
    ```
*  下载next: git clone https://github.com/iissnan/hexo-theme-next themes/next
,并设置为主题: ./blog/_config.yml theme: next
* 配置next的语言等
    ```
    language: zh-Hans
    ```
*  [SEO](http://blog.mobing.net/content/hexo/hexo-next-seo.html)

#### Reference

[hexo-reference](https://www.npmjs.com/package/hexo-reference)

#### Blog上传图片管理


#### SEO 



### 新环境

已经checkout blog source, 新电脑

```
sudo apt-get install nodejs
sudo apt-get install npm
npm install
// 没有将next作为子module，没有保存hexo config file
git clone https://github.com/theme-next/hexo-theme-next.git themes/next
```

