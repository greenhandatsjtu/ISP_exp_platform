# ISP_exp_platform
《信息安全综合实践》在线实验平台

## 运行项目

### 搭建K8S环境

推荐使用[easzlab](https://github.com/easzlab)/[kubeasz](https://github.com/easzlab/kubeasz)脚本安装，也可自行安装。

> **注**：为K8S集群安装插件时会遇到`ImagePullBackOff`问题，这是因为国内的网络问题导致拉去镜像失败。可以在阿里云申请一个自己的加速器（ https://cr.console.aliyun.com/cn-hangzhou/instances/mirrors ），然后修改`/etc/docker/daemon.json`中的`registry-mirrors`，填入申请到的加速器地址，并执行以下命令重启docker即可：
>
> ```bash
> sudo systemctl daemon-reload
> sudo systemctl restart docker
> ```
>

K8S环境搭建后自行验证可用性，可参考：https://github.com/easzlab/kubeasz/tree/master/docs/practice/go_web_app

### 后端

#### 安装`MySQL`和`Redis`

> **注**：适用于`Ubuntu`系统，其他发行版自行查找资料安装

首先安装`MySQL`：

```bash
sudo apt-get update
sudo apt-get install mysql-server
```

初始化配置：

```bash
sudo mysql_secure_installation
```

为了遇到避免数据库不支持中文的问题，请先检查数据库编码是否为**utf8mb4**。否则需进行如下操作：

编辑`/etc/mysql/mysql.conf.d/mysqld.cnf`，添加：

```
character_set_server=utf8mb4
collation-server=utf8mb4_unicode_ci
```

然后安装`Redis`：

```bash
sudo apt-get install redis-server
```

设置mysql和redis自启动：

```bash
sudo systemctl enable mysql
sudo systemctl enable redis-server
```

```
Synchronizing state of mysql.service with SysV init with /lib/systemd/systemd-sysv-install...
Executing /lib/systemd/systemd-sysv-install enable mysql
```

#### 安装`Go`和后端依赖

`apt`安装的`Go`版本过低，因此采用直接下载二进制的方式：

```bash
wget https://studygolang.com/dl/golang/go1.14.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.14.6.linux-amd64.tar.gz
```

然后设置`env`启用`module`：

```bash
go env -w GOPROXY=https://goproxy.io,direct
go env -w GO111MODULE=on
```

接着安装后端的依赖，但在此之前，由于我们用到的client-go（`Go`和K8S交互的接口）的版本与Kubernetes版本有关，对于不同的kubernetes版本，需要下载对应的client-go版本才行。

例如，利用`kubectl version`查询得到本机K8S版本：

```
Client Version: version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.2", GitCommit:"52c56ce7a8272c798dbc29846288d7cd9fbae032", GitTreeState:"clean", BuildDate:"2020-04-16T11:56:40Z", GoVersion:"go1.13.9", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.2", GitCommit:"52c56ce7a8272c798dbc29846288d7cd9fbae032", GitTreeState:"clean", BuildDate:"2020-04-16T11:48:36Z", GoVersion:"go1.13.9", Compiler:"gc", Platform:"linux/amd64"}
```

根据[`client-go`安装文档](https://github.com/kubernetes/client-go/blob/master/INSTALL.md)，我们需要安装`0.18.2`版本的`client-go`，因此需要在`go.mod`中将`k8s.io/client-go v0.x.y`改为`k8s.io/client-go v0.18.2`。

然后在命令执行如下命令安装依赖即可：

```bash
go get -v
```

#### 初始化数据库

在开启后端前，需要先迁移数据库（即根据后端定义的`struct`自动生成数据库）和生成随机数据。

编辑`backend/main.go`，在`import`中添加`"backend/utils"`，并取消下面这两行的注释：

```go
// init database with fake data
utils.InitDatabase(database.Db)
return
```

然后`go run .`即可，初始化完成后，务必撤销对`main.go`的修改。

随机生成的用户中，`ID`1-2为系统管理员，3-5为教务老师，6-15为教师，16-120为学生（其中16-20为助教），密码都为`test`，用邮箱和密码登录。

#### 运行后端

```bash
go build .
./backend -d=true
```

后端运行在`18080`端口。

> 若后端部署在K8S集群外，需要将K8S节点的`.kube/config`（K8S配置文件）复制到本机的`~/.kube/config`下才能正常运行后端。

> 现在主机重启后，后端不会自动运行，因此还需要再手动执行上面的命令

### 前端

预先安装`nodejs`，然后进入`frontend`文件夹，执行：

```bash
npm install
```

若在开发中，则`npm run serve`即可；若需要部署服务，则`npm run build`，会生成`dist`文件夹。

### `Nginx`部署服务

为了避免跨域问题，我们需要用`Nginx`将前后端结合起来。

首先安装`Nginx`：

```bash
sudo apt install nginx
sudo systemctl enable nginx
```

创建两个文件夹：

```bash
mkdir /etc/nginx/sites-available
mkdir /etc/nginx/sites-enabled
```

在`sites-available`里创建`k8s.conf`，写入如下内容：

```nginx
server{
	listen 80;
	charset utf-8;
    root {项目路径}/ISP_exp_platform/frontend/dist;
	location / {
		try_files $uri $uri/ @router;
		index index.html index.htm;
	}
	location @router{
		rewrite ^.*$/index.html last;
	}
	location ^~/api/ {
		proxy_pass http://localhost:18080;
	}
}
```

接着编辑`/etc/nginx/nginx.conf`，在`http`块的末尾添加：

```
http {
    ...
    include sites-enabled/*;
}
```

接下来：

```bash
ln -s /etc/nginx/sites-available/k8s.conf /etc/nginx/sites-enabled/k8s.conf
```

然后重启`nginx`即可：

```bash
sudo systemctl restart nginx
```

至此服务部署完成，可以自行访问`localhost`验证。

#### 测试

之前使用过`dorowu/ubuntu-desktop-lxde-vnc`和`citizenstig/dvwa`来测试前后端和K8S功能（主要测试能否开启实验）。

自行编写实验的`deployment`和`service`的yaml，然后在前端页面上传yaml，并尝试开启实验，访问提示的端口，进行验证。

示例：

`deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vnc
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vnc
  template:
    metadata:
      labels:
        app: vnc
    spec:
      containers:
        - name: vnc
          image: dorowu/ubuntu-desktop-lxde-vnc:xenial
          ports:
            - containerPort: 80
```

`service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: vnc
spec:
  ports:
    - port: 80
      targetPort: 80
  selector:
    app: vnc
  type: NodePort
```

