# Go With Grafana And Prometheus



## Prometheus & Grafana 简单认识与理解

### 1. 什么是 Prometheus？

#### 1.1 通俗理解

👉 它是个`「监控数据管家」`，专门收集和存储各种监控指标的数字（比如CPU温度、网站访问量、内存用量等）并**保存在它自带的数据库里**。

---

#### 1.2 核心功能

|特点	|说明|
|:------|:------|
|📈 时间序列数据库	|按时间顺序记录指标变化（如：CPU使用率每分钟的值）|
|🔍 主动拉取模式	|定期去各个服务"抄表"（HTTP访问 /metrics 接口）|
|🔔 告警系统	|发现异常时能自动发送报警（需要额外配置Alertmanager）|
|🧩 多维数据模型	|用标签区分数据（如区分不同服务器的CPU指标）|

---


### 2. 什么是 Grafana？

#### 2.1 ​通俗理解​：

👉 它是个`「数据可视化画家」`，能把枯燥的数字变成漂亮的**图表和仪表盘**。它**不存储数据，只负责展示**。

---


#### 2.2 核心能力​

|特点	|说明|
|:------|:------|
|🎨 可视化专家	|支持折线图/柱状图/仪表盘/热力图等30+图表类型|
|🔌 万能连接器	|能对接Prometheus、MySQL、InfluxDB等50+数据源|
|📱 看板定制	|自由组合监控面板（类似组装汽车仪表盘）|
|🚨 智能告警	|可基于图表阈值设置报警规则|

----

## 安装部署（Ubuntu版）

### 1. 安装 Prometheus

```bash
# 下载安装包（最新稳定版）
wget https://github.com/prometheus/prometheus/releases/download/v2.48.0/prometheus-2.48.0.linux-amd64.tar.gz

# 解压到系统目录
tar xvfz prometheus-*.tar.gz
sudo mv prometheus-2.48.0.linux-amd64 /opt/prometheus
```

---

### 2. Prometheus 配置与启动

```bash
# 创建配置文件（保留默认配置即可）
sudo nano /opt/prometheus/prometheus.yml

# 启动服务（前台运行方便调试）
cd /opt/prometheus
./prometheus
```

✅ ​验证安装​：浏览器打开 http://localhost:9090 【ip可根据实际需求修改】

![promethuespage](./image/promethuespage.png)


---

### 3. 安装 Grafana

```bash
# 下载安装包
wget https://dl.grafana.com/oss/release/grafana-10.2.1.linux-amd64.tar.gz

# 解压到系统目录
tar xvfz grafana-*.tar.gz
sudo mv grafana-10.2.1 /opt/grafana
```

---

### 4. Grafana 启动与验证

```bash
# 启动服务
cd /opt/grafana/bin
./grafana-server
```

✅ ​验证安装​：浏览器打开 http://localhost:3000 【ip可根据实际需求修改】

🔑 默认账号: admin(账号)/admin(密码) 
> (首次登录要求改密码)

![grafanapage](./image/grafanapage.png)

---

