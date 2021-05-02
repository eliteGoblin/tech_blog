## 自动上传插件

picgo: 可在extension中找到插件，看说明

ctrl + alt + e: browse电脑文件，上传
ctrl + alt + u: clipboard上传

[第四章 vscode作为Markdown编辑器](https://www.jianshu.com/p/cb8d2194d5ef)

选择github作为图床, vscode中配置. (默认smms api deprecated, 上传失败)

在Markdown中用快捷键, 自动生成链接

## 图片简单proessing

```
sudo apt-get install imagemagick
```

convert
```
convert howtogeek.png howtogeek.jpg
# quality for jpeg
convert howtogeek.png -quality 95 howtogeek.jpg
```

resize

```
# try best to keep ratio
convert example.png -resize 200x100 example.png
# 强制按size
convert example.png -resize 200x100! example.png
# 仅指定width
convert example.png -resize 200 example.png
# 仅指定height
convert example.png -resize x100 example.png
# rotate
convert howtogeek.jpg -rotate 90 howtogeek-rotated.jpg
```

目前是600: 

```
convert Ginkgo_leaves_1280x720.jpg -resize 600 ginkgo.jpg
```
refer: https://www.howtogeek.com/109369/how-to-quickly-resize-convert-modify-images-from-the-linux-terminal/