FROM node:16

ENV ACTION_DEPLOY_KEY=
VOLUME /hexo_workground

WORKDIR /hexo_workground

ADD ./ .

RUN npm install hexo-cli -g

ENTRYPOINT [ "sh", "/hexo_workground/entrypoint.sh" ]