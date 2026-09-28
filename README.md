# go-ddd — 领域驱动设计（DDD）四层架构 · Go 实现

参考 Java 版订单中台 [3thinkthendoit/think-oms-ddd](https://github.com/3thinkthendoit/think-oms-ddd)
落地的一套 **DDD 四层架构**（接口层 / 应用层 / 领域层 / 基础设施层）Go 工程。

业务场景对齐参考项目：**多渠道订单接入 → 订单风控审核 → 金额拆分 → 分仓拆单 → 推 WMS 履约 → 发货回传上游平台**。

> 本项目是一次完整的重构：旧的「菱形架构」分包（`acl` / `ohs` / `application/{assembler,gateway,local}` /
> `domain/station` / `infrastructure/component`）已全部废弃，替换为下面的四层分包。

---

## 一、目录结构

<!-- DIRSTRUCTURE_START_MARKER -->
<pre>
go-ddd/
├─ interface/ .............................. 接口层（北向网关，驱动适配器）
│  ├─ ohs/ ................................. 远程网关
│  │  ├─ controller/ ....................... HTTP 接口
│  │  ├─ dto/ .............................. 北向网关入参 / 出参（req / resp）
│  │  ├─ job/ .............................. 定时任务（拉取淘宝 / 拼多多订单）
│  │  ├─ listener/ ......................... 本地事件监听
│  │  ├─ mq/ ............................... 远程事件监听（WMS 发货消息）
│  │  └─ rpc/ .............................. 微服务接口
│  ├─ local/ ............................... 本地网关（进程内统一入口）
│  └─ httpx/ ............................... 北向 HTTP 适配器公共辅助
├─ application/ ............................ 应用层
│  └─ service/ ............................. 应用服务（用例编排，CQRS）
├─ domain/ ................................. 领域层（核心，不依赖任何外层与框架）
│  ├─ common/assert/ ....................... 领域内共用的断言工具
│  ├─ model/ ............................... 领域模型
│  │  ├─ aggregate/ ........................ 聚合 / 聚合根
│  │  │  ├─ create/ ........................ 订单创建聚合
│  │  │  ├─ orderfulfill/ .................. 订单履约聚合（分仓 / 拆单）
│  │  │  └─ shippingcallback/ .............. 发货回传聚合
│  │  ├─ constant/ ......................... 业务枚举
│  │  ├─ dp/ ............................... DP 模型（Domain Primitive）
│  │  └─ valueobject/ ...................... 值对象
│  ├─ pl/ .................................. 转换层（PL：Payload Layer）
│  │  ├─ command/ .......................... 创建 / 变更聚合根的 command
│  │  ├─ event/ ............................ 领域事件
│  │  ├─ query/ ............................ 查询对象
│  │  ├─ request/ .......................... 南向网关 request
│  │  └─ response/ ......................... 南向网关 response
│  ├─ port/ ................................ 领域层端口（南向网关）
│  │  ├─ gateway/ .......................... 其他域的网关接口
│  │  ├─ publisher/ ........................ 领域事件发布接口
│  │  └─ repository/ ....................... 资源库接口
│  └─ service/ ............................. 领域服务（与其他域协作的领域行为）
├─ infrastructure/ ......................... 基础设施层
│  ├─ acl/ ................................. 南向网关实现（防腐层）
│  │  ├─ api/ .............................. 外部平台协议实现
│  │  │  ├─ douyin/ ........................ 抖音
│  │  │  ├─ mijia/ ......................... 小米 / 米家
│  │  │  ├─ pdd/ ........................... 拼多多
│  │  │  └─ taobao/ ........................ 淘宝天猫
│  │  ├─ gateway/ .......................... port-gateway 实现
│  │  ├─ pl/ ............................... 领域对象 ↔ PO ↔ DTO 转换
│  │  ├─ publisher/ ........................ port-publisher 实现
│  │  └─ repository/ ....................... port-repository 实现
│  ├─ common/ .............................. 通用工具
│  │  ├─ demo/ ............................. 本地演示数据
│  │  └─ util/ ............................. 工具类
│  ├─ config/ .............................. 配置文件
│  └─ core/ ................................ 技术组件实现
│     ├─ event/ ............................ 进程内事件总线
│     ├─ http/ ............................. HTTP 客户端
│     ├─ mq/ ............................... 消息队列客户端
│     ├─ orm/ .............................. 持久化（含 po）
│     ├─ redis/ ............................ 缓存
│     ├─ scheduler/ ........................ 定时任务调度
│     └─ server/ ........................... HTTP 服务端
├─ doc/example/ ............................ 示例请求报文
├─ starter/ ................................ 启动类
│  ├─ main.go .............................. 入口 + 各层 bean 扫描
│  └─ wire.go .............................. 领域层装配
├─ go.mod
└─ README.md
</pre>
<!-- DIRSTRUCTURE_END_MARKER -->

---

## 二、依赖规则

```
        ┌────────────────────────────────────────────┐
        │  interface（北向网关）                       │
        │  controller / job / listener / mq / rpc      │
        └───────────────┬────────────────────────────┘
                        │ 只能向下依赖
        ┌───────────────▼────────────────────────────┐
        │  application（应用层）                       │
        │  OrderAppService：用例编排，无业务规则        │
        └───────────────┬────────────────────────────┘
                        │
        ┌───────────────▼────────────────────────────┐
        │  domain（领域层，最内层）                     │
        │  model / pl / port / service                 │
        │  ✅ 不依赖 interface、application、infrastructure
        │  ✅ 不依赖任何框架（连 DI 框架都不 import）    │
        └───────────────▲────────────────────────────┘
                        │ 实现 port（依赖倒置）
        ┌───────────────┴────────────────────────────┐
        │  infrastructure（基础设施层）                 │
        │  acl / core / common / config                │
        └────────────────────────────────────────────┘
```

几条硬约束：

1. **领域层零框架依赖。** `domain/**` 不 import `gs`、不 import 任何基础设施包。领域模型只表达业务规则。
2. **外层依赖内层，内层通过 port 反转依赖。** 领域层定义 `port/gateway`、`port/repository`、`port/publisher`，
   基础设施层在 `acl/*` 里实现，运行时由容器注入。
3. **聚合根状态只能通过领域方法变更。** 所有字段私有，只暴露只读访问器 + 行为方法
   （`Check` / `PriceCalculate` / `PriorityProcessing` / `Dispatch` / `Split` / `Hangup`）。
4. **PL 传输对象用导出字段。** `domain/pl/**` 是纯数据载体，遵循 Go 惯例用导出字段；
   只有 `domain/model/**`（聚合、实体、值对象、DP）才做封装与不变式保护。

---

## 三、与 Java 参考项目的对应关系

| Java（think-oms-ddd） | 本项目 | 说明 |
| --- | --- | --- |
| `think-oms-interface` | `interface/` | 接口层（北向网关） |
| `com.think.oms.ohs.*` | `interface/ohs/` | controller / dto / job / listener / mq / rpc |
| `com.think.oms.local` | `interface/local/` | 本地网关 |
| `think-oms-application` | `application/` | 应用层 |
| `com.think.oms.app.service.OrderAppService` | `application/service/order_app_service.go` | 用例编排 |
| `think-oms-domain` | `domain/` | 领域层 |
| `domain/model/aggregate/{create,orderfulfill,shippingcallback}` | 同名同结构 | 三个聚合 |
| `domain/model/{constant,dp,valueobject}` | 同名同结构 | 枚举 / DP / 值对象 |
| `domain/pl/{command,event,query,request,response}` | 同名同结构 | 转换层 |
| `domain/port/{gateway,publisher,repository}` | 同名同结构 | 领域端口 |
| `domain/service/*` | `domain/service/*` | 领域服务 |
| `think-oms-infrastructure` | `infrastructure/` | 基础设施层 |
| `infrastructure/acl/{api,gateway,pl,publisher,repository}` | 同名同结构 | 南向网关实现 |
| `infrastructure/core/{http,mybatis,redis,rockermq}` | `infrastructure/core/{http,orm,redis,mq}` | 技术组件 |
| `ThinkOmsApplication` | `starter/main.go` + `starter/wire.go` | 启动类 |
| `@Service` / `@Component` / `@RestController` | 各包 `init()` 里的 `gs.Object(...)` | bean 注册 |
| `@Autowired` | `autowire:""` 结构体标签 | 依赖注入 |
| `@EventListener` + `@Async` | `Bus.SubscribeAsync(...)` | 事件订阅 |
| `@Scheduled(fixedRate=...)` | `Scheduler.Every(name, interval, fn)` | 定时任务 |
| `@RequestMapping` / `@PostMapping` | `gs.PostMapping(path, handler)` | 路由注册 |
| `@Value("${key:default}")` | `value:"${key:=default}"` 结构体标签 | 配置绑定 |
| Spring `ApplicationEventPublisher` | `infrastructure/core/event.Bus` | 进程内事件总线 |
| MyBatis + MySQL | `infrastructure/core/orm.Store`（内存实现） | 持久化 |
| RocketMQ | `infrastructure/core/mq.RocketMqClient`（进程内实现） | 消息队列 |

**容器**：沿用项目原有的 [go-spring](https://github.com/go-spring/spring-core)（v1.1.1），
它是 Go 生态里最接近 Spring 的 IoC 容器，能 1:1 承载上表的注解语义。

---

## 四、一次下单的完整链路

```
POST /api/order/mijia/create
  └─ MiJiaController            解析米家协议 → OrderCreateCommand
      └─ OrderLocalService      本地网关
          └─ OrderAppService.CreateOrder   ── 应用层编排
              ├─ OrderCreateAggregate.Create      组装聚合根 + 命令校验
              ├─ OrderCreateDomainService.IsExist 幂等：查订单是否已接入
              ├─ OrderCreateDomainService.InitBaseInfo
              │    └─ SkuInfoQueryGateway → 商品中心（外部 sku → 内部 sku）
              ├─ OrderCreateDomainService.Audit
              │    ├─ aggregate.Check()           领域不变式校验
              │    └─ RiskCheckGateway → 风控；命中则 aggregate.Hangup()
              ├─ aggregate.PriceCalculate()       优惠 / 运费按权重分摊到 sku item
              ├─ aggregate.PriorityProcessing()   生鲜优先级 3 / 虚拟商品 0
              ├─ OrderCreateRepository.Save()     主表 + sku + sku item 落库
              └─ OrderEventPublisher.Publish(OrderCreatedEvent)
                    │
                    ▼  (进程内事件总线，异步)
              OrderCreatedListener
                └─ OrderAppService.DispatchOrder
                    ├─ OrderFulfillRepository.OfByOrderNo   还原履约聚合根
                    ├─ OrderFulfillDomainService.InitBaseInfo
                    │    └─ WarehouseQueryGateway → 仓库库存 / 指定仓映射
                    ├─ aggregate.Check() / Dispatch() / Split()
                    │    分仓三级策略：店铺指定仓 > sku 指定仓 > 最优仓（能力-距离 + 库存贪心）
                    ├─ OrderFulfillRepository.Save()        拆单结果落库
                    └─ OrderEventPublisher.Publish(OrderFulfillEvent)
                          │
                          ▼
                    OrderFulfilledListener
                      └─ OrderAppService.FulfillOrder
                          ├─ OrderFulfillRepository.QueryFulfillOrderInfos
                          ├─ OrderFulfillGateway.Fulfill()   推 WMS
                          └─ OrderFulfillRepository.UpdateOrderFulfill()

POST /api/wms/callback  （真实环境由 WMS 直接投递 MQ）
  └─ WmsCallbackController → RocketMqClient.Send(topic)
      └─ WmsConsumer             远程事件监听
          └─ OrderAppService.ShippingCallback
              ├─ ShippingCallbackAggregate.Create / InitBaseInfo
              │    └─ OrderInfoGateway 查订单，合并 WMS 发货数量
              ├─ aggregate.Check()                 判断是否需要回传上游平台
              ├─ OrderShippingDomainService.ShippingCallback
              │    └─ ShippingCallbackGateway → 淘宝 / 拼多多 / 抖音 / 米家协议
              ├─ SkuShippingRepository.Save()       发货数量 + 回传记录落库
              └─ OrderCreateRepository.Update()     全部发完 → 主订单置为「已发货」
```

---

## 五、快速开始

### 1. 启动

```bash
# 国内网络建议先设置代理
export GOPROXY=https://goproxy.cn,direct

go run ./starter
# 日志：refresh 41 beans / http server started on :8000
```

端口与开关在 `infrastructure/config/application.properties`，例如：

```properties
web.server.port=8000
oms.scheduler.enabled=false      # 定时拉单任务开关
oms.sku-center.fallback=true     # 商品中心不可用时使用内置商品目录
```

### 2. 接入一笔订单

```bash
curl -X POST http://127.0.0.1:8000/api/order/mijia/create \
  -H "Content-Type: application/json" \
  -d @doc/example/mijia-create-order.json
```

示例订单包含 4 个 sku，一次请求即可观察到全部领域行为：

| sku | 类型 | 预期行为 |
| --- | --- | --- |
| SKU1001 | 生鲜单品 ×2 | 优先级 3；最优仓选中上海仓 |
| SKU1002 | 百货单品 ×1 | 优先级 1；同上海仓，与 SKU1001 合成 1 张发货单 |
| SKU1003 | 虚拟商品 ×1 | 优先级 0，**不参与分仓发货** |
| SKU1004 | 组合商品 ×1 | 被指定从广州仓发货，单独拆出 1 张发货单 |

### 3. 查询订单

```bash
curl -X POST http://127.0.0.1:8000/orderInfoIface/query \
  -H "Content-Type: application/json" \
  -d '{"externalOrderNo":"MJ20260917001","orderSource":6}'
```

### 4. 发货回传

把上面 `data[0].orderNo` 和日志里拆出的发货单号填进去：

```bash
# 发货单 1（上海仓）
curl -X POST http://127.0.0.1:8000/api/wms/callback \
  -H "Content-Type: application/json" \
  -d '{"orderNo":"<OMS订单号>","omsOrderNo":"<发货单号1>",
       "shippingInfos":[{"skuCode":"SKU1001","skuAmount":2,"expressCode":"SF","expressNo":"SF1001"},
                        {"skuCode":"SKU1002","skuAmount":1,"expressCode":"SF","expressNo":"SF1002"}]}'

# 发货单 2（广州仓）
curl -X POST http://127.0.0.1:8000/api/wms/callback \
  -H "Content-Type: application/json" \
  -d '{"orderNo":"<OMS订单号>","omsOrderNo":"<发货单号2>",
       "shippingInfos":[{"skuCode":"SKU1004","skuAmount":1,"expressCode":"YTO","expressNo":"YTO2001"}]}'
```

两张发货单都回传后，订单状态自动变为 **已发货（3）**。

> **部分发货**：WMS 回传的 `skuAmount` 是**本次发货数量（增量）**，不是累计值。
> 同一 sku 分多次发出（例如先发 1 件、再发 1 件）时数量会累加，全部发满后主订单才流转为「已发货」。

### 5. 单元测试

```bash
go test ./domain/...      # 领域层回归测试（发货回传累加 / 留痕序列化 / 回传判断）
go build ./... && go vet ./... && gofmt -l .
```

### 6. 其它可验证的场景

```bash
# 幂等：重复接入同一笔订单 → 订单 externalOrderNo=xxx 已经存在
# 领域校验：orderStatus 改为 -1（未支付）或 2（取消）→ 只允许已支付订单接入
# 领域校验：orderStatus 改为 0（未定义状态）→ 同样拒绝
# 风控挂起：地址里带上「黑名单」→ 订单状态变为「挂起(4)」，不触发履约
```

> 若本机设置了 `HTTP_PROXY`，访问 `127.0.0.1` 需加 `--noproxy '*'`，否则请求会被代理拦成 502。

---

## 六、相对 Java 参考实现做的修正

参考项目是一份教学 DEMO，部分代码是骨架 / 存在笔误。移植时按业务意图修正并在代码注释中标明：

**领域规则**

| 位置 | 原实现 | 修正后 |
| --- | --- | --- |
| `OrderCreateAggregate.check` | `Assert.isTrue(PAYED != status)` —— 断言反了，已支付反而抛异常 | 要求 `status == PAYED`，未支付不允许接入 |
| `OrderCreateAggregate.check` | `Assert.isNull(invoiceInfo)` 等三处 —— 断言反了 | 改为「必须非空」校验 |
| `OrderCreateAggregate.priceCalculate` | 金额分摊逻辑被注释掉，`skuItems` 无金额 | 按 sku 行小计权重真实分摊运费 / 优惠，余数由末位 item 承担 |
| `OrderFulfillAggregate.split` | `orderNo = splitResults.get(0).getParentOrderNo()` 赋错字段 | 同仓库复用同一父单号，生成正确的拆单结果 |
| `FulfillSkuItem` | `dispatchInfo` 未初始化（NPE） | 构造时初始化，`Dispatch` 累加 |
| `FulfillSkuItem` | 分仓数量用 `shippingAmount`（已发货量，初始 0） | 改用待发货量 `skuAmount - shippingAmount` |
| `OrderFulfillAggregate.dispatch` | 虚拟商品（优先级 0）也会参与分仓 | 优先级 0 直接跳过，不产生发货单 |
| `ShippingCallbackAggregate.check` | `orderId.getOrderSource()` 为 null（NPE） | 领域服务查询订单后回填订单来源与外部订单号 |
| `ShippingCallbackAggregate` | 回传只带 `wmsOrderNo`，无法定位 OMS 订单 | 命令同时携带 `orderNo` + `wmsOrderNo` |
| 发货回传落库 | 回传记录生成后没有落库 | 发货数量与回传记录一起落库 |
| 部分 sku 回传 | 会把其它 sku 的已发货数量冲成 0 | 补全时带入订单上已有的发货数量 |
| 同一 sku 多次发货 | `ModifyShippingInfo` 直接赋值，把累计已发量冲成「本次发货量」 | 改为 `AddShippingInfo` **累加**；WMS 回传的是本次发货增量，不累加则部分发货永远凑不满、主订单无法流转为「已发货」 |
| 回传记录明细 | 直接 `json.Marshal` 聚合内的 `ShippingSkuItem`，字段全未导出 → 留痕是 `{"SKU1001":{}}` | 转成带 json tag 的快照结构并按 skuCode 排序，留痕内容真实且可比对 |
| `initSkuShippingInfo` | 把 OMS 订单号当作 `omsOrderNos` 查 WMS | 改用拆单后的发货单号 `WmsOrderNo()` |
| 米家 `orderStatus` | `int` 接收 + `== 0` 即改写为「已支付」 | 改用 `*int` 区分「字段缺省」与「显式传 0」，避免零值把显式未支付悄悄改写成已支付、架空领域守卫 |

**枚举与常量**

- `PayType`：原 `TRAN_FEE(1)` 与 `UNION_PAY(1)` 编码重复 → 修正为 `BALANCE(1)/ALIPAY(2)/WX_PAY(3)/UNION_PAY(4)`
- `Currency`：原 `INSTALL_FEE/TRAN_FEE` 是复制粘贴残留 → 修正为 `CNY(1)/USD(2)`
- `SkuType.SINGEL` / `SkuCategory.FRUSH_FOOD` 拼写 → `SINGLE` / `FRESH_FOOD`
- `OrderSource` 补充 `MI_JIA(6)`，否则 `MiJiaController` 没有对应的来源编码

**Go 化调整**

- 断言 / 校验统一返回 `error` 而非抛异常，由应用层决定中断用例还是走业务补偿
- `pl` 传输对象使用导出字段；只有 `domain/model` 做封装
- 领域层不 import DI 框架，领域服务的 bean 注册集中在 `starter/wire.go`
- 补齐参考项目缺失的 `WarehouseQueryGateway`（分仓需要仓库库存与指定仓映射）

---

## 七、接入真实中间件

各技术组件都是**可替换的进程内实现**，接口形态与真实组件一致，替换时上层代码不用动：

| 组件 | 位置 | 替换方式 |
| --- | --- | --- |
| 持久化 | `infrastructure/core/orm` | 保持 `Store` 方法签名，换成 GORM / sqlx + MySQL |
| 缓存 | `infrastructure/core/redis` | 换成 `go-redis` |
| 消息队列 | `infrastructure/core/mq` | 换成 RocketMQ / Kafka 客户端 |
| 商品中心 | `acl/gateway/sku_info_query_gateway_impl.go` | 实现 `parseSkuCenterResponse`，关闭 `oms.sku-center.fallback` |
| WMS | `acl/gateway/order_fulfill_gateway_impl.go` | 关闭 `oms.wms.mock`，调用真实接口 |
| 风控 | `acl/gateway/risk_check_gateway_impl.go` | 替换内置规则为风控系统调用 |
| 仓库域 | `acl/gateway/warehouse_query_gateway_impl.go` | 替换内置仓库数据为仓储系统调用 |
| 各平台协议 | `acl/api/{taobao,pdd,douyin,mijia}` | 关闭 `oms.platform-mock.*`，接入开放平台 SDK |

内置演示数据集中在这几处，接入真实环境时删除即可：`infrastructure/common/demo`、
`acl/gateway/sku_info_query_gateway_impl.go` 的 `builtinSkuCatalog`、
`acl/gateway/warehouse_query_gateway_impl.go` 的 `builtinWarehouses`。

---

## 八、License

见 [LICENSE](LICENSE)。
