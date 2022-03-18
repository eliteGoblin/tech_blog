FROM node:16

RUN npm install hexo-cli -g

ARG USER_NAME=frank.sun
ARG USER_ID=1001
ARG GROUP_ID=1001
ARG USER_HOME=/home/${USER_NAME}/

RUN groupadd -g ${GROUP_ID} ${USER_NAME} &&\
    useradd -l -u ${USER_ID} -g  ${USER_NAME} ${USER_NAME} &&\
    install -d -m 0755 -o ${USER_NAME} -g ${USER_NAME} ${USER_HOME}

USER ${USER_NAME}
RUN mkdir ${USER_HOME}/hexo_workground
WORKDIR ${USER_HOME}/hexo_workground
# CMD sh

ENV ACTION_DEPLOY_KEY=

ADD --chown=${USER_ID}:${GROUP_ID} ./ .

ENTRYPOINT [ "sh", "/home/frank.sun/hexo_workground/entrypoint.sh" ]