---
layout: post
title: 'Python3 virtual environment: venv and conda'
date: 2022-01-25 10:07:21
tags:
---

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20220318120901.png" alt="20220318120901" style="width:500px"/>

## Preface

Python并不承诺backward compatibility: 我们可能需要在不同版本切换; 而且global package management的问题在于不区分特定库版本: 无法满足version A和B同时存在. 

我们需要隔离的Python环境: virtual env: 每个Python project之间隔离: 选定自己的Python版本和dependency.

<!-- more -->

## venv

Projects一般建立自己的virtual environment, 实现development isolation.

常用的命令: 

```s
# create a virtual env same with system's Python version
python3 -m venv venv_test
# activate venv
source venv_test/bin/activate
# deactivate a venv
deactivate
```

使用venv使得package安装隔离: 

```s
which python
# /usr/bin/python, a python2
python3 -m venv venv_test
source venv_test/bin/activate
which python
# /tmp/python_venv/venv_test/bin/python, python3
which pip # pip also with Python3
# pip 20.0.2 from /tmp/python_venv/venv_test/lib/python3.8/site-packages/pip (python 3.8)
echo $PATH # PATH also added venv folder
# /tmp/python_venv/venv_test/bin:...
```

`venv_test/bin/activate` 是shell script, 用来设置venv的环境变量, 指明当前Python及package的路径等;

新创建的venv只有pip和其依赖的package: 

```
pip list  
Package       Version
------------- -------
pip           20.0.2 
pkg-resources 0.0.0  
setuptools    44.0.0 
```

用 `pip freeze`来show和lock我们自己安装的package: 

```s
pip freeze
# nothing to show
pip install arrow
pip freeze
#arrow==1.2.2
#python-dateutil==2.8.2
#six==1.16.0
pip show arrow # 安装在venv folder, 而不是system
#Location: /tmp/python_venv/venv_test/lib/python3.8/site-packages
#Requires: python-dateutil
#Required-by: 
python -c "import site; import sys; print(sys.prefix); print(site.getsitepackages())" # 安装位置也能从这里看到
# /tmp/python_venv/venv_test 
# ['/tmp/python_venv/venv_test/lib/python3.8/site-packages',...]
```

## Conda

Conda多用于machinelearning, 也可用来manage全局的不同environment, 同时提供了package management, 约等于virtualenv+pip. Conda本身支持非常多的package, 尤其是ML; 没有的话再用pip安装. 

[Conda vs. pip vs. virtualenv commands](https://docs.conda.io/projects/conda/en/latest/commands.html#conda-vs-pip-vs-virtualenv-commands)  

Python环境称为environment, 相关cmd: 

```s
# 建立env能指定Python版本
conda create --name py35 python=3.5
# list所有environment
conda env list
# conda environments:
# base                     /home/frank.sun/local/anaconda3
# leetcode-cn           *  /home/frank.sun/local/anaconda3/envs/leetcode-cn
# delete a environment
conda env remove --name bio-env
conda activate leetcode-cn
conda deactivate
```

Package management: 
```s
conda install PACKAGENAME
# freeze version
conda list --explicit > bio-env.txt
# install from file
conda env create --file bio-env.txt
```

Note: 

*  `conda env list` 中的 `*` 并不一定代表当前shell用哪个, 使用前还是`conda activate xx`

# Reference

*  [Python Virtual Environments: A Primer](https://realpython.com/python-virtual-environments-a-primer/)  