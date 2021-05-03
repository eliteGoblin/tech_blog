

## Docker 

*  本地和CI公用一个docker image: pin version as possible
*  docker image 存在Dockerhub上, 没有npm更新时，用同一个hexo container
*  dockerhub login用`dhl`: 
*  更新了package.json呢? 
    -  docker image 按date tag
    -  更新了，手动push image到docker hub
    -  update本地脚本和github workflow, 用最新image.

build
```s
cd /home/frank.sun/git_repo/blog/docker
docker build -t elitegoblin/hexo_workground --no-cache . 
# export IMAGE_HASH=$(docker images elitegoblin/hexo_workground | grep "node" | awk '{print $3}')
# export IMAGE_TAG=node16_hexo5_$(date '+%Y%m%d')
# docker tag ${IMAGE_HASH} elitegoblin/hexo_workground:${IMAGE_TAG}
# docker push elitegoblin/hexo_workground:${IMAGE_TAG}
# echo "New Image: " elitegoblin/hexo_workground:${IMAGE_TAG}
```

run
```s
# use latest for local atm
cd /home/frank.sun/git_repo/blog/docker
docker run -ti --rm --name="hexo_deploy" -p 4000:4000 -e "ACTION_DEPLOY_KEY=$(cat ~/.ssh/github-deploy-key)" -v "$(pwd)"/../hexo_workground/source/:/hexo_workground/source  elitegoblin/hexo_workground
```

## Workflow

