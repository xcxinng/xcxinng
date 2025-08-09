# 安装
基于docker-compose安装, 参考[官方安装文档](https://docs.nebula-graph.com.cn/3.8.0/4.deployment-and-installation/2.compile-and-install-nebula-graph/3.deploy-nebula-graph-with-docker-compose/)

# 初次运行
```bash
# 进入nebula-docker-compose目录
cd /Users/xianchaoxing/go/src/example.com/nebula-docker-compose
docker compose up -d
```
如果需要手工停止容器（不移除容器/卷/网络）
```bash
docker compose stop
```

如果之前通过docker-compose stop停掉过容器，通过start命令重新运行容器:
```bash
docker compose start
```

# 连接
通过nebula-console连接:
```bash
# nebula-console 是一个二进制文件，放置在 ~/Downloads 目录下
cd ~/Downloads
# 默认用户名: root, 密码: nebula
./nebula-console -u root -p nebula -P 9669
```
文档参考: [nebula-console](https://docs.nebula-graph.com.cn/3.8.0/2.quick-start/3.quick-start-on-premise/3.connect-to-nebula-graph/)

# 关闭
```bash
# 临时停掉容器，不会删除
docker compose stop
# 删除容器(Stops containers and removes containers, networks, volumes, and images created by up.)
docker compose down
```

# studio安装
参考 [studio安装文档](https://docs.nebula-graph.com.cn/3.8.0/nebula-studio/deploy-connect/st-ug-deploy/#docker_studio)

默认docker-compose的配置会把studio和nebula-graph部署在不同网络下，需要通过下面命令把容器加入到graph的网络中：
```bash
docker network connect {nebula-net} {containerID}
```
studio本地目录(docker-compose.yml配置文件)： /Users/xianchaoxing/Downloads/nebula-graph-studio-3.10.0

# 访问studio

浏览器访问 http://localhost:7007

graphIPAddr: 通过docker inspect {containerID} 获取任意一个graph节点的IP地址
username： root
password： nebula
address: 172.20.0.9:9669
