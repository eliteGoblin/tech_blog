

# Docker 

*  本地和CI公用一个docker image: pin version as possible
*  docker image 存在Dockerhub上, 没有npm更新时，用同一个hexo container
*  dockerhub login用`dhl`: 
*  更新了package.json呢? 
    -  docker image 按date tag
    -  更新了，手动push image到docker hub
    -  update本地脚本和github workflow, 用最新image.

## user, group

为了让docker create的markdown能被外部vscode编辑(否则在host上, owner是root), 需要让docker container内部运行和host一样的user和group.

```sh
RUN groupadd -g ${GROUP_ID} ${USER_NAME} &&\
    useradd -l -u ${USER_ID} -g  ${USER_NAME} ${USER_NAME} &&\
    install -d -m 0755 -o ${USER_NAME} -g ${USER_NAME} ${USER_HOME}

USER ${USER_NAME}
RUN mkdir ${USER_HOME}/hexo_workground
WORKDIR ${USER_HOME}/hexo_workground
```

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
docker run -ti --rm --name="hexo_deploy" -p 4000:4000 -e "ACTION_DEPLOY_KEY=$(cat ~/.ssh/github-deploy-key)" -v "$(pwd)"/../hexo_workground/source/:/home/frank.sun/hexo_workground/source  elitegoblin/hexo_workground
```

## Workflow

