# my_blog_src

Source of my technical blog

## Setup docker environment

Docker container 包含working environment: 将所有dependency fix, 防止某个package upgrade, environment broken.

build docker
```s
cd /home/frank.sun/git_repo/blog/docker
docker build -t elitegoblin/hexo_workground --no-cache . 
```

run docker:
```s
cd /home/frank.sun/git_repo/blog/docker
docker run -ti --rm --name="hexo_deploy" -p 4000:4000 -e "ACTION_DEPLOY_KEY=$(cat ~/.ssh/github-deploy-key)" -v "$(pwd)"/../hexo_workground/source/:/home/frank.sun/hexo_workground/source  elitegoblin/hexo_workground
```


Pls refer to [setup](./doc/setup.md)

## Write blog
Create a new post

```sh
hexo new post "Datetime, timezone and Python"
```

New article will appear in `hexo_workground/source/_posts`

Local hosting:
```sh
hexo s
```

Delete an article:

```sh
```