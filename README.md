# go-ddd

<!-- DIRSTRUCTURE_START_MARKER -->
<pre>
go-ddd/
├─ acl/ .............. 南向网关
│  ├─ assembler/ ..... 南向网关数据转领域对象
│  └─ impl/ .......... port/gateway实现类
│     ├─ gateway/ .... gateway实现类
│     └─ repository/ . repository实现类
├─ application/ ...... 应用层
│  ├─ assembler/ ..... 北向网关数据转领域对象
│  ├─ gateway/ ....... 直接调用南向网关
│  ├─ local/ ......... 本地service,防腐
│  └─ service/ ....... 业务编排
├─ domain/ ........... 领域层
│  ├─ model/ ......... 领域模型
│  │  ├─ aggregate/ .. 聚合根
│  │  ├─ command/ .... 创建聚合根命令
│  │  ├─ entity/ ..... 实体
│  │  └─ vo/ ......... 值对象,DP
│  ├─ port/ .......... port接口
│  │  ├─ gateway/ .... 其他领域接口
│  │  └─ repository/ . 资源库接口
│  └─ service/ ....... 领域服务
├─ infrastructure/ ... 基础设施层
│  ├─ common/ ........ 通用类
│  │  ├─ config/ ..... 配置
│  │  └─ util/ ....... 工具类
│  └─ component/ ..... 组件
│     ├─ http/ ....... http 实现
│     ├─ mysql/ ...... mysql 实现
│     ├─ redis/ ...... redis 实现
│     └─ rocketmq/ ... MQ实现
├─ osh/ .............. 北向网关
│  ├─ adapter/ ....... 适配器
│  │  ├─ http/ ....... http接口
│  │  ├─ mq/ ......... MQ消息/领域消息
│  │  └─ rpc/ ........ 微服务接口
│  ├─ cqe/ ........... 查询命令分离
│  │  ├─ cmd/ ........ 命令，增删查
│  │  ├─ event/ ...... 远程事件，领域事件
│  │  └─ query/ ...... 查询服务
│  └─ dto/ ........... 北向网关返回对象
└─ starter/ .......... 启动类
   └─ main.go ........ 
</pre>
<!-- DIRSTRUCTURE_END_MARKER -->
