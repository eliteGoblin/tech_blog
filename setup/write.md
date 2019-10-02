
#### 命令

*  hexo new "TCP connection reuse when HTTP in golang"
*  hexo s
*  hexo g // generate html file
*  hexo d // push html to gitpage

#### 注意

*  header 不要跳级，会让索引混乱

#### tricks

#####  插入pdf

npm install --save hexo-pdf
{% pdf ./bash_freshman.pdf %}


#### 添加about页面

hexo new page "about"

菜单显示 about 链接，在主题的 _configy.yml 设置中将 menu 中 about 前面的注释去掉即可。

menu:
  home: /
  archives: /archives
  tags: /tags
  about: /about

[hexo wiki](https://github.com/iissnan/hexo-theme-next/wiki/)