## 1. ZooKeeper 介绍

> 大家可以了解一下[Paxos的小岛(Island)](https://www.douban.com/note/208430424/)，以便更好的理解Zookeeper的概念

### 1.1 什么是Zookeeper

`ZooKeeper` 是一个开源的**分布式协调服务框架**，为分布式系统提供一致性服务。

那么什么是分布式？什么是协调程序？和集群又有什么区别？

举一个例子来说明，现在有一个网上商城购物系统，并发量太大单机系统承受不住，那我们可以多加几台服务器支持大并发量的访问需求，这个就是所谓的**`Cluster` 集群** 。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210308193345.png)

如果我们将这个网上商城购物系统拆分成多个子系统，比如订单系统、积分系统、购物车系统等等，**然后将这些子系统部署在不同的服务器上** ，这个时候就是 **`Distributed` 分布式** 。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210308193736.png)

对于集群来说，多加几台服务器就行（当然还得解决session共享，负载均衡等问题），而对于分布式来说，你首先需要将业务进行拆分，然后再加服务器，同时还要去解决分布式带来的一系列问题。比如各个分布式组件如何**协调**起来，如何减少各个系统之间的耦合度，如何处理分布式事务，如何去配置整个分布式系统，如何解决各分布式子系统的数据不一致问题等等。`ZooKeeper` 主要就是解决这些问题的。

### 1.2 使用ZooKeeper的开源项目

许多著名的开源项目用到了 ZooKeeper，比如：

1. **Kafka** : ZooKeeper 主要为 Kafka 提供 Broker 和 Topic 的注册以及多个 Partition 的负载均衡等功能。
2. **Hbase** : ZooKeeper 为 Hbase 提供确保整个集群只有一个 Master 以及保存和提供 regionserver 状态信息（是否在线）等功能。
3. **Hadoop** : ZooKeeper 为 Namenode 提供高可用支持。
4. **Dubbo**：阿里巴巴集团开源的分布式服务框架，它使用 ZooKeeper 来作为其命名服务，维护全局的服务地址列表。

### 1.3 ZooKeeper的三种运行模式

ZooKeeper 有三种运行模式：单机模式、伪集群模式和集群模式。

- 单机模式：这种模式一般适用于开发测试环境，一方面我们没有那么多机器资源，另外就是平时的开发调试并不需要极好的稳定性。
- 集群模式：一个 ZooKeeper 集群通常由一组机器组成，一般 3 台以上就可以组成一个可用的 ZooKeeper 集群了。组成 ZooKeeper 集群的每台机器都会在内存中维护当前的服务器状态，并且每台机器之间都会互相保持通信。
- 伪集群模式：这是一种特殊的集群模式，即集群的所有服务器都部署在一台机器上。当你手头上有一台比较好的机器，如果作为单机模式进行部署，就会浪费资源，这种情况下，ZooKeeper 允许你在一台机器上通过启动不同的端口来启动多个 ZooKeeper 服务实例，从而以集群的特性来对外服务。

## 2. CAP和BASE理论

一个分布式系统必然会存在一个问题：**因为分区容忍性（partition tolerance）的存在，就必定要求我们需要在系统可用性（availability）和数据一致性（consistency）中做出权衡** 。这就是著名的 `CAP` 定理。

举个例子来说明，假如班级代表整个分布式系统，而学生是整个分布式系统中一个个独立的子系统。这个时候班里的小红小明偷偷谈恋爱被班里的小花发现了，小花欣喜若狂告诉了周围的人，然后小红小明谈恋爱的消息在班级里传播起来了。当在消息的传播（散布）过程中，你问班里一个同学的情况，如果他回答你不知道，那么说明整个班级系统出现了数据不一致的问题（因为小花已经知道这个消息了）。而如果他直接不回答你，因为现在消息还在班级里传播（为了保证一致性，需要所有人都知道才可提供服务），这个时候就出现了系统的可用性问题。

这个例子中前者就是 `Eureka` 的处理方式，它保证了AP（可用性），后者就 `ZooKeeper` 的处理方式，它保证了CP（数据一致性）。

CAP理论中，`P`（分区容忍性）是必然要满足的，因为毕竟是分布式，不能把所有的应用全放到一个服务器里面，这样服务器是吃不消的。所以，只能从AP（可用性）和CP（一致性）中找平衡。

怎么个平衡法呢？在这种环境下出现了**BASE理论**：即使无法做到强一致性，但分布式系统可以根据自己的业务特点，采用适当的方式来使系统达到最终的一致性。BASE理论由：`Basically Avaliable` 基本可用、`Soft state` 软状态、`Eventually consistent` 最终一致性组成。

- **基本可用(Basically Available)**：基本可用是指分布式系统在出现故障的时候，允许损失部分可用性，即保证核心可用。例如，电商大促时，为了应对访问量激增，部分用户可能会被引导到降级页面，服务层在该页面只提供降级服务。
- **软状态(Soft State)**： 软状态是指允许系统存在中间状态，而该中间状态不会影响系统整体可用性。分布式存储中一般一份数据至少会有多个副本，允许不同节点间副本同步的延时就是软状态的体现。
- **最终一致性(Eventual Consistency)**： 最终一致性是指系统中的所有数据副本经过一定时间后，最终能够达到一致的状态。弱一致性和强一致性相反，最终一致性是弱一致性的一种特殊情况。

一句话概括就是：平时系统要求是基本可用，运行有可容忍的延迟状态，但是，无论如何经过一段时间的延迟后系统最终必须达成数据是一致的。

> ACID 是传统数据库常用的设计理念，追求强一致性模型。BASE 支持的是大型分布式系统，通过牺牲强一致性获得高可用性。

其实可能发现不管是CAP理论，还是BASE理论，他们都是理论，这些理论是需要算法来实现的，这些算法有2PC、3PC、Paxos、Raft、ZAB，它们所解决的问题全部都是：**在分布式环境下，怎么让系统尽可能的高可用，而且数据能最终能达到一致**。

## 3. Zookeeper的特点

> 该部分来源于[讲解 Zookeeper 的五个核心知识点](https://mp.weixin.qq.com/s?__biz=MzI4NjI1OTI4Nw==&mid=2247489891&idx=1&sn=eb7b6a4d4f2560df31eb41e10dc66264&chksm=ebdef85bdca9714d89dcd84894ce8c3af4c2ad5bab43760adfa30384fe7345c6d93fe55b47d5&mpshare=1&scene=24&srcid=0308ammHQgSRw3GGV7RWu4M3&sharer_sharetime=1615167268668&sharer_shareid=18383980e942ee6dfd94ea4b7b61fcbe&ascene=14&devicetype=android-29&version=2700153b&nettype=WIFI&abtest_cookie=AAACAA%3D%3D&lang=zh_CN&exportkey=AUeEnm4ZGyJa8Eg6QqeW7W8%3D&pass_ticket=He9OMc%2Bmhiuj111RUsXzzTJbt%2B9kiaQht3Dd7kCsxQpc8HgWiMvTeMy4aVZ1XSPB&wx_header=1)。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309102817.png)

1. **集群**：Zookeeper是一个领导者（Leader），多个跟随者（Follower）组成的集群。
2. **高可用性**：集群中只要有半数以上节点存活，Zookeeper集群就能正常服务。
3. **全局数据一致**：每个Server保存一份相同的数据副本，Client无论连接到哪个Server，数据都是一致的。
4. **更新请求顺序进行**：来自同一个Client的更新请求按其发送顺序依次执行。
5. **数据更新原子性**：一次数据更新要么成功，要么失败。
6. **实时性**：在一定时间范围内，Client能读到最新数据。
7. 从设计模式角度来看，zk是一个基于**观察者设计模式**的框架，它负责管理跟存储大家都关心的数据，然后接受观察者的注册，数据反生变化zk会通知在zk上注册的观察者做出反应。
8. Zookeeper是一个**分布式协调系统**，满足CP性，跟SpringCloud中的Eureka满足AP不一样。

## 4. 一致性协议之 ZAB

> 推荐大家先了解其他的一致性算法，如2PC、3PC、Paxos、Raft，可参考[大数据中的 2PC、3PC、Paxos、Raft、ZAB](https://mp.weixin.qq.com/s/b5mGEbn-FLb9vhOh1OpwIg)。

作为一个优秀高效且可靠的分布式协调框架，`ZooKeeper` 在解决分布式数据一致性问题时并没有直接使用 `Paxos` ，而是专门定制了一致性协议叫做 `ZAB(ZooKeeper Automic Broadcast)` 原子广播协议，该协议能够很好地支持 **崩溃恢复** 。

### 4.1 ZAB 中的三个角色

ZAB 中三个主要的角色，Leader 领导者、Follower跟随者、Observer观察者 。

- `Leader` ：集群中 **唯一的写请求处理者** ，能够发起投票（投票也是为了进行写请求）。
- `Follower`：能够接收客户端的请求，如果是读请求则可以自己处理，**如果是写请求则要转发给 `Leader`** 。在选举过程中会参与投票，**有选举权和被选举权** 。
- `Observer` ：就是没有选举权和被选举权的 `Follower` 。

在 `ZAB` 协议中对 `zkServer`(即上面我们说的三个角色的总称) 还有两种模式的定义，分别是 **消息广播** 和 **崩溃恢复** 。

### 4.2 ZXID和myid

**ZooKeeper** 采用全局递增的事务 id 来标识，所有 proposal(提议)在被提出的时候加上了**ZooKeeper Transaction Id** 。ZXID是64位的Long类型，**这是保证事务的顺序一致性的关键**。ZXID中高32位表示纪元**epoch**，低32位表示事务标识**xid**。你可以认为zxid越大说明存储数据越新，如下图所示：

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309202934.png)

1. 每个leader都会具有不同的**epoch**值，表示一个纪元/朝代，用来标识 leader周期。每个新的选举开启时都会生成一个新的epoch，从1开始，每次选出新的Leader，epoch递增1，并会将该值更新到所有的zkServer的zxid的epoch。
2. **xid**是一个依次递增的事务编号。数值越大说明数据越新，可以简单理解为递增的事务id。**每次epoch变化，都将低32位的序号重置**，这样保证了zxid的全局递增性。

每个ZooKeeper服务器，都需要在数据文件夹下创建一个名为myid的文件，该文件包含整个ZooKeeper集群唯一的id（整数）。例如，某ZooKeeper集群包含三台服务器，hostname分别为zoo1、zoo2和zoo3，其myid分别为1、2和3，则在配置文件中其id与hostname必须一一对应，如下所示。在该配置文件中，`server.`后面的数据即为myid

```
tex

server.1=zoo1:2888:3888
server.2=zoo2:2888:3888
server.3=zoo3:2888:3888
```

### 4.3 历史队列

每一个follower节点都会有一个**先进先出**（FIFO)的队列用来存放收到的事务请求，保证执行事务的顺序。所以：

- 可靠提交由ZAB的事务一致性协议保证
- 全局有序由TCP协议保证
- 因果有序由follower的历史队列(history queue)保证

### 4.4 消息广播模式

ZAB协议两种模式：消息广播模式和崩溃恢复模式。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309204307.png)

说白了就是 `ZAB` 协议是如何处理写请求的，上面我们不是说只有 `Leader` 能处理写请求嘛？那么我们的 `Follower` 和 `Observer` 是不是也需要 **同步更新数据** 呢？总不能数据只在 `Leader` 中更新了，其他角色都没有得到更新吧。

第一步肯定需要 `Leader` 将写请求 **广播** 出去呀，让 `Leader` 问问 `Followers` 是否同意更新，如果超过半数以上的同意那么就进行 `Follower` 和 `Observer` 的更新（和 `Paxos` 一样）。消息广播机制是通过如下图流程**保证事务的顺序一致性**的：

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309205602.png)

1. leader从客户端收到一个写请求
2. leader生成一个新的事务并为这个事务生成一个唯一的ZXID
3. leader将这个事务发送给所有的follows节点，将带有 zxid 的消息作为一个提案(proposal)分发给所有 follower。
4. follower节点将收到的事务请求加入到历史队列(history queue)中，当 follower 接收到 proposal，先将 proposal 写到硬盘，写硬盘成功后再向 leader 回一个 ACK
5. 当leader收到大多数follower（超过一半）的ack消息，leader会向follower发送commit请求（leader自身也要提交这个事务）
6. 当follower收到commit请求时，会判断该事务的ZXID是不是比历史队列中的任何事务的ZXID都小，如果是则提交事务，如果不是则等待比它更小的事务的commit(保证顺序性)
7. Leader将处理结果返回给客户端

**过半写成功策略**：Leader节点接收到写请求后，这个Leader会将写请求广播给各个Server，各个Server会将该写请求加入历史队列，并向Leader发送ACK信息，当Leader收到一半以上的ACK消息后，说明该写操作可以执行。Leader会向各个server发送commit消息，各个server收到消息后执行commit操作。

这里要注意以下几点：

- Leader并不需要得到Observer的ACK，即Observer无投票权
- Leader不需要得到所有Follower的ACK，只要收到过半的ACK即可，**同时Leader本身对自己有一个ACK**
- Observer虽然无投票权，但仍须同步Leader的数据从而在处理读请求时可以返回尽可能新的数据

另外，Follower/Observer也可以接受写请求，此时：

- Follower/Observer接受写请求以后，不能直接处理，而需要将写请求转发给Leader处理
- 除了多了一步请求转发，其它流程与直接写Leader无任何区别
- Leader处理写请求是通过上面的消息广播模式，实质上最后所有的zkServer都要执行写操作，这样数据才会一致

而对于读请求，Leader/Follower/Observer都可直接处理读请求，从本地内存中读取数据并返回给客户端即可。由于处理读请求不需要各个服务器之间的交互，因此Follower/Observer越多，整体可处理的读请求量越大，也即读性能越好。

### 4.5 崩溃恢复模式

恢复模式大致可以分为四个阶段：选举、发现、同步、广播。

1. **选举阶段**（Leader election）：当leader崩溃后，集群进入选举阶段（下面会将如何选举Leader），开始选举出潜在的准 leader，然后进入下一个阶段。
2. **发现阶段**（Discovery）：用于在从节点中发现最新的ZXID和事务日志。准Leader接收所有Follower发来各自的最新epoch值。Leader从中选出最大的epoch，基于此值加1，生成新的epoch分发给各个Follower。各个Follower收到全新的epoch后，返回ACK给Leader，带上各自最大的ZXID和历史提议日志。Leader选出最大的ZXID，并更新自身历史日志，此时Leader就用拥有了最新的提议历史。（注意：每次epoch变化时，ZXID的第32位从0开始计数）。
3. **同步阶段**（Synchronization）：主要是利用 leader 前一阶段获得的最新提议历史，同步给集群中所有的Follower。只有当超过半数Follower同步成功，这个准Leader才能成为正式的Leader。这之后，follower 只会接收 zxid 比自己的 lastZxid 大的提议。
4. 广播阶段（Broadcast）：集群恢复到广播模式，开始接受客户端的写请求。

> 在发现阶段，或许有人会问：既然Leader被选为主节点，已经是集群里数据最新的了，为什么还要从节点中寻找最新事务呢？这是为了防止某些意外情况。所以这一阶段，Leader集思广益，接收所有Follower发来各自的最新epoch值。

这里有两点要注意：

（1）**确保已经被Leader提交的提案最终能够被所有的Follower提交**

假设 `Leader (server2)` 发送 `commit` 请求（忘了请看上面的消息广播模式），他发送给了 `server3`，然后要发给 `server1` 的时候突然挂了。这个时候重新选举的时候我们如果把 `server1` 作为 `Leader` 的话，那么肯定会产生数据不一致性，因为 `server3` 肯定会提交刚刚 `server2` 发送的 `commit` 请求的提案，而 `server1` 根本没收到所以会丢弃。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309214242.png)

那怎么解决呢？

**这个时候 `server1` 已经不可能成为 `Leader` 了，因为 `server1` 和 `server3` 进行投票选举的时候会比较 `ZXID` ，而此时 `server3` 的 `ZXID` 肯定比 `server1` 的大了**（后面讲到选举机制时就明白了）。同理，只能由server3当Leader，server3当上Leader之后，在同步阶段，会将最新提议历史同步给集群中所有的Follower，这就保证数据一致性了。如果server2在某个时刻又重新恢复了，它作为`Follower` 的身份进入集群中，再向Leader同步当前最新提议和Zxid即可。

（2）**确保跳过那些已经被丢弃的提案**

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309215027.png)

假设 `Leader (server2)` 此时同意了提案N1，自身提交了这个事务并且要发送给所有 `Follower` 要 `commit` 的请求，却在这个时候挂了，此时肯定要重新进行 `Leader` 的选举，假如此时选 `server1` 为 `Leader` （这无所谓，server1和server2都可以当选）。但是过了一会，这个 **挂掉的 `Leader` 又重新恢复了** ，此时它肯定会作为 `Follower` 的身份进入集群中，需要注意的是刚刚 `server2` 已经同意提交了提案N1，但其他 `server` 并没有收到它的 `commit` 信息，所以其他 `server` 不可能再提交这个提案N1了，这样就会出现数据不一致性问题了，所以 **该提案N1最终需要被抛弃掉** 。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309215338.png)

### 4.6 脑裂问题

脑裂问题：所谓的“脑裂”即“大脑分裂”，也就是本来一个“大脑”被拆分了两个或多个“大脑”。通俗的说，就是比如当你的 cluster 里面有两个节点，它们都知道在这个 cluster 里需要选举出一个 master。那么当它们两之间的通信完全没有问题的时候，就会达成共识，选出其中一个作为 master。但是如果它们之间的通信出了问题，那么两个结点都会觉得现在没有 master，所以每个都把自己选举成 master，于是 cluster 里面就会有两个 master。

ZAB为解决脑裂问题，要求集群内的节点数量为2N+1, 当网络分裂后，始终有一个集群的节点数量过半数，而另一个集群节点数量小于N+1（即小于半数）, 因为选主需要过半数节点同意，所以任何情况下集群中都不可能出现大于一个leader的情况。

因此，有了过半机制，对于一个Zookeeper集群，要么没有Leader，要没只有1个Leader，这样就避免了脑裂问题。

## 5. Zookeeper选举机制

`Leader` 选举可以分为两个不同的阶段，第一个是我们提到的 `Leader` 宕机需要重新选举，第二则是当 `Zookeeper` 启动时需要进行系统的 `Leader` 初始化选举。下面是zkserver的几种状态：

- **LOOKING** 不确定Leader状态。该状态下的服务器认为当前集群中没有Leader，会发起Leader选举。
- **FOLLOWING** 跟随者状态。表明当前服务器角色是Follower，并且它知道Leader是谁。
- **LEADING** 领导者状态。表明当前服务器角色是Leader，它会维护与Follower间的心跳。
- **OBSERVING** 观察者状态。表明当前服务器角色是Observer，与Folower唯一的不同在于不参与选举，也不参与集群写操作时的投票。

### 5.1 初始化Leader选举

假设我们集群中有3台机器，那也就意味着我们需要2台同意（超过半数）。这里假设服务器1~3的myid分别为1,2,3，初始化Leader选举过程如下：

1. 服务器 1 启动，发起一次选举。它会首先 **投票给自己** ，投票内容为`(myid, ZXID)`，因为初始化所以 `ZXID` 都为0，此时 `server1` 发出的投票为`(1, 0)`，即`myid`为1， `ZXID`为0。此时服务器 1 票数一票，不够半数以上，选举无法完成，服务器 1 状态保持为 LOOKING。
2. 服务器 2 启动，再发起一次选举。服务器2首先也会将投票选给自己`(2, 0)`，并将投票信息广播出去（`server1`也会，只是它那时没有其他的服务器了），`server1` 在收到 `server2` 的投票信息后会将投票信息与自己的作比较。**首先它会比较 `ZXID` ，`ZXID` 大的优先为 `Leader`，如果相同则比较 `myid`，`myid` 大的优先作为 `Leader`**。所以，**此时`server1` 发现 `server2` 更适合做 `Leader`，它就会将自己的投票信息更改为`(2, 0)`然后再广播出去**，之后`server2` 收到之后发现和自己的一样无需做更改。此时，服务器1票数0票，服务器2票数2票，**投票已经超过半数**，确定 `server2` 为 `Leader`。服务器 1更改状态为 FOLLOWING，服务器 2 更改状态为 LEADING。
3. 服务器 3 启动，发起一次选举。此时服务器 1，2已经不是 LOOKING 状态，它会直接以 `FOLLOWING` 的身份加入集群。

### 5.2 运行时Leader选举

运行时候如果Leader节点崩溃了会走崩溃恢复模式，新Leader选出前会暂停对外服务，大致可以分为四个阶段：选举、发现、同步、广播（见4.5节），此时Leader选举流程如下：

1. Leader挂掉，剩下的两个 `Follower` 会将自己的状态 **从 `Following` 变为 `Looking` 状态** ，每个Server会发出一个投票，第一次都是投自己，其中投票内容为`(myid, ZXID)`，注意这里的 `zxid` 可能不是0了
2. 收集来自各个服务器的投票
3. 处理投票，处理逻辑：**优先比较ZXID，然后比较myid**
4. 统计投票，只要超过半数的机器接收到同样的投票信息，就可以确定leader
5. 改变服务器状态Looking变为Following或Leading
6. 然后依次进入发现、同步、广播阶段

举个例子来说明，假设集群有三台服务器，`Leader (server2)`挂掉了，只剩下server1和server3。 `server1` 给自己投票为(1,99)，然后广播给其他 `server`，`server3` 首先也会给自己投票(3,95)，然后也广播给其他 `server`。`server1` 和 `server3` 此时会收到彼此的投票信息，和一开始选举一样，他们也会比较自己的投票和收到的投票（`zxid` 大的优先，如果相同那么就 `myid` 大的优先）。这个时候 `server1` 收到了 `server3` 的投票发现没自己的合适故不变，`server3` 收到 `server1` 的投票结果后发现比自己的合适于是更改投票为(1,99)然后广播出去，最后 `server1` 收到了发现自己的投票已经超过半数就把自己设为 `Leader`，`server3` 也随之变为 `Follower`。

## 6. Zookeeper数据模型

ZooKeeper 数据模型（Data model）采用层次化的多叉树形结构，每个节点上都可以存储数据，这些数据可以是数字、字符串或者是二级制序列。并且，每个节点还可以拥有 N 个子节点，最上层是根节点以`/`来代表。

每个数据节点在 ZooKeeper 中被称为 **znode**，它是 ZooKeeper 中数据的最小单元。并且，每个 znode 都一个唯一的路径标识。由于**ZooKeeper 主要是用来协调服务的，而不是用来存储业务数据的**，这种特性使得 Zookeeper 不能用于存放大量的数据，每个节点的存放数据上限为**1M**。

和文件系统一样，我们能够自由的增加、删除**znode**，在一个**znode**下增加、删除子**znode**，唯一的不同在于**znode**是可以存储数据的。默认有四种类型的**znode**：

1. **持久化目录节点 PERSISTENT**：客户端与zookeeper断开连接后，该节点依旧存在。
2. **持久化顺序编号目录节点 PERSISTENT_SEQUENTIAL**：客户端与zookeeper断开连接后，该节点依旧存在，只是Zookeeper给该节点名称进行顺序编号。
3. **临时目录节点 EPHEMERAL**：客户端与zookeeper断开连接后，该节点被删除。
4. **临时顺序编号目录节点 EPHEMERAL_SEQUENTIAL**：客户端与zookeeper断开连接后，该节点被删除，只是Zookeeper给该节点名称进行顺序编号。

在zookeeper客户端使用`get`命令可以查看znode的内容和状态信息：

```
bash

[zk: localhost:2181(CONNECTED) 2] get /zk01
updateed02
cZxid = 0x600000023
ctime = Mon Mar 01 21:20:26 CST 2021
mZxid = 0xb0000000d
mtime = Fri Mar 05 17:15:53 CST 2021
pZxid = 0xb00000018
cversion = 5
dataVersion = 7
aclVersion = 0
ephemeralOwner = 0x0
dataLength = 10
numChildren = 3
```

下面我们来看一下每个 znode 状态信息究竟代表的是什么吧

| znode 状态信息 | 解释                                                         |
| -------------- | ------------------------------------------------------------ |
| cZxid          | create ZXID，即该数据节点被创建时的事务 id                   |
| ctime          | create time，znode 被创建的毫秒数(从1970 年开始)             |
| mZxid          | modified ZXID，znode 最后更新的事务 id                       |
| mtime          | modified time，znode 最后修改的毫秒数(从1970 年开始)         |
| pZxid          | znode 最后更新子节点列表的事务 id，只有子节点列表变更才会更新 pZxid，子节点内容变更不会更新 |
| cversion       | znode 子节点变化号，znode 子节点修改次数，子节点每次变化时值增加 1 |
| dataVersion    | znode 数据变化号，节点创建时为 0，每更新一次节点内容(不管内容有无变化)该版本号的值增加 1 |
| aclVersion     | znode 访问控制列表(ACL )版本号，表示该节点 ACL 信息变更次数  |
| ephemeralOwner | 如果是临时节点，这个是 znode 拥有者的 sessionid。如果不是临时节，则 ephemeralOwner=0 |
| dataLength     | znode 的数据长度                                             |
| numChildren    | znode 子节点数量                                             |

## 7. Zookeeper监听通知机制

**Watcher** 监听机制是 **Zookeeper** 中非常重要的特性，我们基于 Zookeeper上创建的节点，可以对这些节点绑定**监听**事件，比如可以监听节点数据变更、节点删除、子节点状态变更等事件，通过这个事件机制，可以基于 **Zookeeper** 实现分布式锁、集群管理等多种功能，它有点类似于订阅的方式，即客户端向服务端 **注册** 指定的 `watcher` ，当服务端符合了 `watcher` 的某些事件或要求则会 **向客户端发送事件通知** ，客户端收到通知后找到自己定义的 `Watcher` 然后 **执行相应的回调方法** 。

当客户端在Zookeeper上某个节点绑定监听事件后，如果该事件被触发，Zookeeper会通过回调函数的方式通知客户端，但是客户端只会收到一次通知。如果后续这个节点再次发生变化，那么之前设置 **Watcher** 的客户端不会再次收到消息（Watcher是一次性的操作），可以通过循环监听去达到永久监听效果。

ZooKeeper 的 Watcher 机制，总的来说可以分为三个过程：

1. 客户端注册 Watcher，注册 watcher 有 3 种方式，getData、exists、getChildren。
2. 服务器处理 Watcher 。
3. 客户端回调 Watcher 客户端。

监听通知机制的流程如下：

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309140525.png)

1. 首先要有一个main()线程
2. 在main线程中创建zkClient，这时就会创建两个线程，一个负责网络连接通信（connet），一个负责监听（listener）。
3. 通过connect线程将注册的监听事件发送给Zookeeper。
4. 在Zookeeper的注册监听器列表中将注册的监听事件添加到列表中。
5. Zookeeper监听到有数据或路径变化，就会将这个消息发送给listener线程。
6. listener线程内部调用了process()方法。

## 8. Zookeeper会话（Session）

Session 可以看作是 ZooKeeper 服务器与客户端的之间的一个 TCP 长连接，客户端与服务端之间的任何交互操作都和Session 息息相关，其中包含zookeeper的临时节点的生命周期、客户端请求执行以及Watcher通知机制等。

> client 端连接 server 端默认的 2181 端口，也就是 session 会话。

接下来，我们从**全局的会话状态变化**到**创建会话**再到**会话管理**三个方面来看看Zookeeper是如何处理会话相关的操作。

### 8.1 会话状态

session会话状态有：

- **connecting**：连接中，session 一旦建立，状态就是 connecting 状态，时间很短。
- **connected**：已连接，连接成功之后的状态。
- **closed**：已关闭，发生在 session 过期，一般由于网络故障客户端重连失败，服务器宕机或者客户端主动断开。

客户端需要与服务端创建一个会话，这个时候客户端需要提供一个服务端地址列表，`host1 : port,host2: port ,host3:port` ，一般由地址管理器(HostProvider)管理，然后根据地址创建zookeeper对象。这个时候客户端的状态则变更为**CONNECTING**，同时客户端会根据上述的地址列表，按照顺序的方式获取IP来尝试建立网络连接，直到成功连接上服务器，这个时候客户端的状态就可以变更为**CONNECTED**。在Zookeeper服务端提供服务的过程中，有可能遇到网络波动等原因，导致客户端与服务端断开了连接，这个时候客户端会进行重新连接操作这个时候的状态为**CONNECTING**,当连接再次建立后，客户端的状态会再次更改为**CONNECTED**，也就是说只要在Zookeeper运行期间，客户端的状态总是能保持在**CONNECTING**或者是**CONNECTED**。当然在建立连接的过程中，如果出现了连接超时、权限检查失败或者是在建立连接的过程中，我们主动退出连接操作，这个时候客户端的状态都会变成**CLOSE**状态。

### 8.2 会话ID的生成

一个会话必须包含以下几个基本的属性：

- **SessionID** : 会话的ID，用来唯一标识一个会话，每一次客户端建立连接的时候，Zookeeper服务端都会给其分配一个全局唯一的**sessionID**。在Zookeeper中，无论是哪台服务器为客户端分配的 `sessionID`，都务必保证全局唯一。
- **Timeout**：一次会话的超时时间，客户端在构造Zookeeper实例的时候，会配置一个**sessionTimeOut**参数用于指定会话的超时的时间。Zookeeper服务端会按照连接的客户端发来的**TimeOut**参数来计算并确定超时的时间。当由于服务器压力太大、网络故障或是客户端主动断开连接等各种原因导致客户端连接断开时，只要在超时规定的时间内能够重新连接上集群中任意一台服务器，那么之前创建的会话仍然有效。
- ExpirationTime：TimeOut是一个相对时间，而ExpirationTime则是在时间轴上的一个绝对过期时间。可能你也会想到，一个比较通用的计算方法就是：`ExpirationTime = CurrentTime + Timeout`。 这样算出来的时间最准确，但ZK可不是这么算的，下面会讲具体计算方式及这样做的原因。
- **TickTime**：下一次会话超时的时间点，为了便于Zookeeper对会话实行分桶策略管理，同时也是为了高效低耗地实现会话的超时检查与清理，Zookeeper会为每个会话标记一个下次会话超时时间点。TickTime是一个13位的Long类型的数值，一般情况下这个值接近**TimeOut**，但是并不完全相等。
- **isCloseing**：用来标记当前会话是否已经处于被关闭的状态。如果服务端检测到当前会话的超时时间已经到了，就会将isCloseing属性标记为已经关闭，这样以后即使再有这个会话的请求访问也不会被处理。

SessionID作为一个全局唯一的标识，我们可以来探究下Zookeeper是如何保证Session会话在集群环境下依然能保证全局唯一性的。

在**sessionTracker**初始化的时候，会调用**initializeNextSession**来生成sessionid，算法大概如下:

```
java

public static long initializeNextSession(long id ) {
    long nextSid = 0;
    nextSid = (System.currentTimeMillis() << 24) >> 8;
    nextSid=nextSid|(id << 56);
    return nextSid;
}
```

从这段代码，我们可以看到session的创建大概分为以下几个步骤：

**1. 获取当前时间的毫秒表示**

我们假设当前System.currentTimeMills()获取的值是1380895182327，其64位二进制表示为:

```
tex

00000000 00000000 00000001 01000001 10000011 11000100 01001101 11110111
```

**2. 接下来左移24位**，我们可以得到结果：

```
tex

01000001 100000011 11000100 01001101 11110111 00000000 00000000 00000000
```

可以看到低位已经把高位补齐，剩下的低位都使用了0补齐。

**3. 右移8位**，结果变成了：

```
tex

00000000 01000001 100000011 11000100 01001101 11110111 00000000 00000000
```

**4. 计算机器码标识ID**：

在initializeNextSession方法中，传入了一个id变量，这个变量就是当前zkServer的myid中配置的值，一般是一个整数，假设此时的值为2，转为64位二进制表示：

```
tex

00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000010
```

此时发现高位几乎都是0，进行左移56位以后，得到值如下：

```
tex

00000010 00000000 00000000 00000000 00000000 00000000 00000000 00000000
```

**5. 将前面第三步和第四步得到的结果进行 | 操作**，可以得到结果为：

```tex
tex

00000010 01000001 10000011 11000100 01001101 11110111 00000000 00000000
```

这个时候我们可以得到一个集群中唯一的序列号ID，整个算法大概可以理解为，**先通过高8位确定zkServer所在的机器以后**，后面的56位按照当前毫秒进行随机，可以看出来当前的算法还是蛮严谨的，基本上看不出来什么明显的问题，但是其实也有问题的。我们可以看到，zk选择了当前机器时间内的毫秒作为基数，但是如果时间到了2022年4月8号以后， `System.currentTimeMillis ()`的值会是多少呢？

```java
java

Calendar calendar = Calendar.getInstance();
calendar.clear();
calendar.set(2022,5,1);
long millis = calendar.getTimeInMillis();
System.out.println(Long.toBinaryString(millis));
```

输出：`0000000000000000000000011000000100011010110110000110110000000000`

可以看到，输出结果前面有23个0，接着我们左移24位以后会发现，这个时候的值竟然是个负数。

> 在java中最高位为1时表示负数，为0表示正整数

为了保证不会出现负数的情况，可以将有符号移位换成无符号移位，解决方案如下:

```java
java

public static long initializeNextSession(long id ) {
    long nextSid = 0;
    nextSid = (System.currentTimeMillis() << 24) >>> 8;
    nextSid=nextSid|(id << 56);
    return nextSid;
}
```

上面`>>>`为无符号右移，当目标是负数时，**在移位时忽略符号位，空位都以0补齐，这样就保证了结果永远是正数**。

### 8.3 SessionTracker与ClientCnxn

**SessionTracker**是Zookeeper中的会话管理器，负责整个zk生命周期中会话的**创建**、**管理**和**清理**操作，而每一个会话在Sessiontracker内部都保留了如下三个数据结构，大体如下:

```java
java

protected final ConcurrentHashMap<Long, SessionImpl> sessionsById =
    new ConcurrentHashMap<Long, SessionImpl>();
private final ConcurrentMap<Long, Integer> sessionsWithTimeout;
```

1. sessionsWithTimeout这是一个ConcurrentHashMap类型的数据结构，用来管理会话的超时时间，这个参数会被持久化到快照文件中去
2. sessionsById是一个HashMap类型的数据结构，用于根据sessionId来管理session实体
3. sessionsSets同样也是一个HashMap类型的数据结构，用来会话超时的时候进行归档，便于进行会话恢复和管理

**ClientCnxn**是Zookeeper客户端的核心工作类，负责维护客户端与服务端之间的网络连接并进行一系列网络通信。

ClientCnxn内部又包含两个线程，**SendThread**是一个I/O线程，主要负责Zookeeper客户端和服务端之间的网络I/O通信，**EventThread**是一个事件线程，主要负责对服务端事件进行处理。

- **SendThread**：SendThread维护了客户端与服务端之间的会话生命周期，其通过一定的周期频率向服务端发送一个PING包来实现心跳检测。此外，SendThread管理了客户端所有的请求发送和响应接受操作，其将上层客户端API操作转换成相应的请求协议并发送到服务端，并完成对同步调用的返回和异步调用的回调。同时，SendThread还负责将来自服务端的事件传递给EventThread去处理。
- **EventThread**：EventThread负责客户端的事件处理，并触发客户端注册的Watcher监听。EventThread中有一个waitingEvents队列，用于临时存放那些需要被触发的Object，包括那些客户端注册的Watcher和异步接口中注册的回调器AsyncCallback。EventThread会不断地从waitingEvents队列中取出Object，识别出具体的类型，并分别调用process（Watcher）和processResult（AsyncCallback）接口方法来实现对事件的触发和回调。

ClientCnxn中有两个核心队列outgoingQueue和pendingQueue，分别代表客户端的请求发送队列和服务端响应的等待队列。

- outgoing队列专门用于存储那些客户端需要发送到服务端的Packet集合
- pending队列存储那些已经从客户端发送到服务端的，但是需要等待服务端响应的Packet结合

clientCnxnSocket是底层Socket通信层，定义了Socket通信的接口，为了便于对底层Socket层进行扩展，例如使用Netty来实现和使用过NIO来实现。在Zookeeper中默认的实现是ClientCnxnSocketNIO，主要负责对请求的发送和响应的接收过程。

**发送请求**

在TCP连接正常情况下，从outgoingQueue队列中按照先进先出的顺序提取出一个可发送的Packet对象，同时生成一个客户端请求序号XID并将其设置到Packet请求头中，然后将其序列化后发送。

请求发送完毕后，会立即将该Packet保存到pendingQueue队列中，以便等待服务端响应返回后进行相应的处理，处理完毕后返回给客户端

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309224615.jpg)

**接收响应**

客户端获取到来自服务端的完整响应数据后，根据不同的客户端请求类型，进行不同的处理：

- 如果检测到当前客户端还未进行初始化，则说明当前客户端与服务端之间正在进行会话创建，那么就直接将接受到的ByteBuffer序列化成ConnectResponse对象
- 如果当前客户端已经处于正常的会话周期，并且接受到的服务端响应是一个事件，那么就将接受到的ByteBuffer序列化成WatcherEvent对象，并将该事件放入待处理队列中
- 如果是一个常规的请求响应（Create、GetData和Exist等操作请求），就会从pendingQueue队列中取出一个Packet，按照XID顺序将接受到的ByteBuffer序列化成响应的Response对象

### 8.4 会话创建

会话的创建的流程如下：

1. client随机选一个服务端地址列表提供的地址，委托给`ClientCnxnSocket`去创建与zk之间的TCP长连接。
2. SendThread会负责根据当前客户端的设置，构造出一个ConnectRequest请求，该请求代表了客户端试图与服务器创建一个会话。同时，Zookeeper客户端还会进一步将请求包装成网络IO的Packet对象，放入请求发送队列——outgoingQueue中去。
3. 当客户端请求准备完毕后，ClientCnxnSocket从outgoingQueue中取出Packet对象，将其序列化成ByteBuffer后，向服务器进行发送。
4. 服务端的SessionTracker为该会话分配一个sessionId，并发送响应。
5. Client收到响应后，会首先判断当前的客户端状态是否是已初始化，如果尚未完成初始化，**那么就认为该响应一定是会话创建请求的响应**，直接交由readConnectResult方法来处理该请求。
6. ClientCnxnSocket会对接受到的服务端响应进行反序列化，得到ConnectResponse对象，并从中获取到Zookeeper服务端分配的会话SessionId。
7. 连接成功后，一方面需要通知SendThread线程，进一步对客户端进行会话参数设置，包括readTimeout和connectTimeout等，并更新客户端状态；另一方面，需要通知地址管理器HostProvider当前成功连接的服务器地址。
8. 为了能够让上层应用感知到会话的成功创建，SendThread会生成一个事件SyncConnected-None，代表客户端与服务器会话创建成功，并将该事件传递给EventThread线程。
9. EventThread线程收到事件后，会从ClientWatchManager管理器中查询出对应的Watcher，针对SyncConnected-None事件，那么就直接找出存储的Watcher，然后将其放到EventThread的waitingEvents队列中。
10. EventThread不断地从waitingEvents队列中取出待处理的Watcher对象，然后直接调用该对象的process接口方法，以达到触发Watcher的目的。

至此，Zookeeper客户端完整的一次会话创建过程已经全部完成了。

### 8.5 会话超时管理

Session是由ZK服务端来进行管理的，一个服务端可以为多个客户端服务，也就是说，有多个Session，那这些Session是怎么样被管理的呢？而分桶机制可以说就是其管理的一个手段。ZK服务端会维护着一个个"桶",然后把Session们分配到一个个的桶里面。而这个区分的维度，就是**ExpirationTime**

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309163108.png)

为什么要如此区分呢？因为ZK的服务端会在运行期间定时地对会话进行超时检测，如果不对Session进行维护的话，那在检测的时候岂不是要遍历所有的Session？这显然不是一个好办法，所以才以超时时间为维度来存放Session，这样在检测的时候，只需要扫描对应的桶就可以了。

那这样的话，新的问题就来了：每个Session的超时时间是一个很分散的值，假设有1000个Session，很可能就会有1000个不同的超时时间，进而有1000个桶，这样有啥意义吗？因此zk的ExpirationTime 用了下面的计算方式：

```
tex

ExpirationTime = CurrentTime + SessionTimeout;
ExpirationTime = (ExpirationTime / ExpirationInterval + 1) * ExpirationInterval;
```

可以看到，最终得到的ExpirationTime是ExpirationInterval的**倍数**，而ExpirationInterval就是ZK服务端定时检查过期Session的频率，默认为2000毫秒。所以说，每个Session的ExpirationTime最后都是一个近似值，是ExpirationInterval的倍数，这样的话，ZK在进行扫描的时候，只需要扫描一个桶即可。

> 另外让过期时间是ExpirationInterval的倍数还有一个好处就是，让检查时间和每个Session的过期时间在一个时间节点上。否则的话就会出现一个问题：ZK检查完毕的1毫秒后，就有一个Session新过期了，这种情况肯定是不好。

为了便于理解，我们可以举几个例子，Zk默认的间隔时间是2000ms：

- 比如我们计算出来一个sessionA在3000ms后过期，，那么其会坐落在`(3000/2000+1)*2000=5000ms`，放在4000ms这个key里。
- 比如我们计算出来一个sessionB在1500ms后过期，那么其会坐落在`(1500/2000+1)*2000=3500ms`，放在2000ms这个key里。

| 0    | 2000ms   | 4000ms   | 6000ms | 8000ms |
| :--- | :------- | :------- | :----- | :----- |
|      | sessionB | sessionA |        |        |

这样线程就不用遍历所有的会话去逐一检查它们的过期时间了，有点妙。如果服务端检测到当前会话的超时时间已经到了，就会将**isCloseing**属性标记为已经关闭，这样以后即使再有这个会话的请求访问也不会被处理。

### 8.5 会话激活

在客户端与服务端完成连接之后生成过期时间，这个值并不是一直不变的，而是会随着客户端与服务端的交互来更新。**过期时间的更新，当然就伴随着Session在桶上的迁移**。过期时间计算的过程则是使用上面的公式，计算完新的超时时间以后，就可以放在桶相应位置上。激活的方式有：

- 客户端每向服务端发送请求，包括读请求和写请求，都会触发一次激活，因为这预示着客户端处于活跃状态
- 而如果客户端一直没有读写请求，那么它在TimeOut的三分之一时间内没有发送过请求的话，那么客户端会发送一次PING，来触发Session的激活。当然，如果客户端直接断开连接的话，那么TimeOut结束后就会被服务端扫描到然后进行清楚了

除此之外，由于会话之间的激活是按照分桶策略进行保存的，因此我们可以利用此策略优化对于会话的超时检查，在Zookeeper中，会话超时检查也是由**SessionTracker**负责的，内部有一个线程专门进行会话的超时检查，只要依次的对每一个区块的会话进行检查。由于分桶是按照**ExpriationInterval** 的倍数来进行会话分布的，因此只要在这些时间点检查即可，这样可以减少检查的次数，并且批量清理会话，实现较高的效率。

### 8.6 会话清理

会话检查操作以后，当发现有超时的会话的时候，会进行会话清理操作，而Zookeeper中的会话清理操作，主要是以下几个步骤:

1. 由于会话清理过程需要一定的时间，为了保证在清理的过程中，该会话不会再去接受和处理发来的请求，因此，在会话检查完毕后，**SessionTracker**会先将其会话的**isClose**标记为true，这样在会话清理期间接收到客户端的新请求也无法继续处理了。
2. 发起关闭会话请求给`PrepRequestProcessor`，使其**在整个Zk集群里生效**。
3. 收集需要清理的临时节点 ——通过sessionsWithTimeout和分桶策略找到超时的会话
4. 当会话对应的临时节点列表找到后，Zookeeper会将列表中所有的节点变成删除节点的请求，并且丢给事物变更队列**OutStandingChanges**中，接着**FinalRequestProcessor**处理器会触发删除节点的操作，从内存数据库中删除。
5. 当会话对应的临时节点被删除以后，就需要将会话从**SessionTracker**中移除了，主要从**SessionById**，**sessionsWithTimeOut**以及**sessionsSets**中将会话移除掉，当一切操作完成后，清理会话操作完成，这个时候将会关闭最终的连接**NioServerCnxn**。

### 8.7 会话重连

在Zookeeper运行过程中，也可能会出现会话断开后重连的情况，这个时候客户端会从连接列表中按照顺序的方式重新建立连接，直到连接上其中一台机器为止。这个时候可能出现两种状态，一种是正常的连接**CONNECTED**，这种情况是Zookeeper客户端在超时时间内连接上了服务端，此时sessionid不变；而超时以后才连接上服务端的话，这个时候的客户端会话状态则为**EXPIRED**，被视为非法会话。

而在重连之前，可能因为其他原因导致的断开连接，即CONNECTION_LESS，会抛出异常**org.apache.zookeeper.KeeperException$ConnectionLossException**。此时，会话可能会出现两种情况：

（1）会话失效：SESSION_EXPIRED

会话失效一般发生在ConnectionLoss期间，客户端尝试开始重连，但是在超时时间以后，才与服务端建立连接的情况，这个时候服务端就会通知客户端当前会话已经失效，我们只能选择重新创建一个会话，进行数据的处理操作

（2）会话转移：SESSION_MOVED

会话转移也是在重连过程中常发生的一种情况，例如在断开连接之前，会话是在服务端A上，但是在断开连接重连以后，最终与服务端B重新恢复了会话，这种情况就称之为会话转移。而会话转移可能会带来一个新的问题，例如在断开连接之前，可能刚刚发送一个创建节点的请求，请求发送完毕后断开了，很短时间内再次重连上了另一台服务端，这个时候又发送了一个一样的创建节点请求，这个时候一样的事物请求可能会被执行了多次。因此在Zookeeper3.2版本开始，就有了会话转移的概念，并且封装了一个**SessionMovedExection**异常出来，在处理客户端请求之前，会检查一遍，请求的会话是不是当前服务端的，如果不存在当前服务端的会话，会直接抛出**SessionMovedExection**异常，当然这个时候客户端已经断开了连接，接受不到服务端的异常响应了。

## 9. Zookeeper分布式锁

> 本小节来自[漫画：如何用Zookeeper实现分布式锁？](https://mp.weixin.qq.com/s/u8QDlrDj3Rl1YjY4TyKMCA)

**分布式锁**是雅虎研究员设计Zookeeper的初衷。利用Zookeeper的临时顺序节点，可以轻松实现分布式锁。

### 9.1 获取锁

首先，在Zookeeper当中创建一个持久节点ParentLock。当第一个客户端想要获得锁时，需要在ParentLock这个节点下面创建一个**临时顺序节点** Lock1。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175457.png)

之后，Client1查找ParentLock下面所有的临时顺序节点并排序，判断自己所创建的节点Lock1是不是顺序最靠前的一个。如果是第一个节点，则成功获得锁。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175518.png)

这时候，如果再有一个客户端 Client2 前来获取锁，则在ParentLock下再创建一个临时顺序节点Lock2。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175541.png)

Client2查找ParentLock下面所有的临时顺序节点并排序，判断自己所创建的节点Lock2是不是顺序最靠前的一个，结果发现节点Lock2并不是最小的。

于是，Client2向排序仅比它靠前的节点Lock1注册**Watcher**，用于监听Lock1节点是否存在。这意味着Client2抢锁失败，进入了等待状态。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175609.png)

这时候，如果又有一个客户端Client3前来获取锁，则在ParentLock下载再创建一个临时顺序节点Lock3。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175635.png)

Client3查找ParentLock下面所有的临时顺序节点并排序，判断自己所创建的节点Lock3是不是顺序最靠前的一个，结果同样发现节点Lock3并不是最小的。

于是，Client3向排序仅比它靠前的节点**Lock2**注册Watcher，用于监听Lock2节点是否存在。这意味着Client3同样抢锁失败，进入了等待状态。

这样一来，Client1得到了锁，Client2监听了Lock1，Client3监听了Lock2。这恰恰形成了一个等待队列，很像是Java当中ReentrantLock所依赖的**AQS**（AbstractQueuedSynchronizer）。

### 9.2 释放锁

释放锁分为两种情况：

**1.任务完成，客户端显示释放**

当任务完成时，Client1会显示调用删除节点Lock1的指令。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175800.png)

**2.任务执行过程中，客户端崩溃**

获得锁的Client1在任务执行过程中，如果Duang的一声崩溃，则会断开与Zookeeper服务端的链接。根据临时节点的特性，相关联的节点Lock1会随之自动删除。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175823.png)

由于Client2一直监听着Lock1的存在状态，当Lock1节点被删除，Client2会立刻收到通知。这时候Client2会再次查询ParentLock下面的所有节点，确认自己创建的节点Lock2是不是目前最小的节点。如果是最小，则Client2顺理成章获得了锁。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175847.png)

同理，如果Client2也因为任务完成或者节点崩溃而删除了节点Lock2，那么Client3就会接到通知。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175913.png)

最终，Client3成功得到了锁。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309175937.png)

### 9.3 Zk和Redis分布式锁的比较

下面的表格总结了Zookeeper和Redis分布式锁的优缺点：

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309180014.png)

有人说Zookeeper实现的分布式锁支持可重入，Redis实现的分布式锁不支持可重入，这是**错误的观点**。两者都可以在客户端实现可重入逻辑。

> 什么是 “可重入”，可重入就是说某个线程已经获得某个锁，可以再次获取锁而不会出现死锁

## 9. Zookeeper几个应用场景

### 9.1 数据发布/订阅

当某些数据由几个机器共享，且这些信息经常变化数据量还小的时候，这些数据就适合存储到ZK中。

- **数据存储**：将数据存储到 Zookeeper 上的一个数据节点。
- **数据获取**：应用在启动初始化节点从 Zookeeper 数据节点读取数据，并在该节点上注册一个数据变更 **Watcher**
- **数据变更**：当变更数据时会更新 Zookeeper 对应节点数据，Zookeeper会将数据变更**通知**发到各客户端，客户端接到通知后重新读取变更后的数据即可。

### 9.2 统一配置管理

本质上，统一配置管理和数据发布/订阅是一样的。

分布式环境下，配置文件的同步可以由Zookeeper来实现。

1. 将配置文件写入Zookeeper的一个ZNode
2. 各个客户端服务监听这个ZNode
3. 一旦ZNode发生改变，Zookeeper将通知各个客户端服务

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309172415.png)

### 9.3 统一集群管理

可能我们会有这样的需求，我们需要了解整个集群中有多少机器在工作，我们想对及群众的每台机器的运行时状态进行数据采集，对集群中机器进行上下线操作等等。

例如，集群机器监控：这通常用于那种对集群中机器状态，机器在线率有较高要求的场景，能够快速对集群中机器变化作出响应。这样的场景中，往往有一个监控系统，实时检测集群机器是否存活。过去的做法通常是：监控系统通过某种手段（比如ping）定时检测每个机器，或者每个机器自己定时向监控系统汇报“我还活着”。 这种做法可行，但是存在两个比较明显的问题：

1. 集群中机器有变动的时候，牵连修改的东西比较多。
2. 有一定的延时。

利用ZooKeeper有两个特性，就可以实时另一种集群机器存活性监控系统：

1. 客户端在某个节点上注册一个Watcher，那么如果该节点的子节点变化了，会通知该客户端。
2. 创建EPHEMERAL类型的节点，一旦客户端和服务器的会话结束或过期，那么该节点就会消失。

如下图所示，监控系统在`/manage`节点上注册一个Watcher，如果`/manage`子节点列表有变动，监控系统就能够实时知道集群中机器的增减情况，至于后续处理就是监控系统的业务了。

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309232621.png)

### 9.4 负载均衡

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309180848.png)

多个相同的jar包在不同的服务器上开启相同的服务，可以通过nginx在服务端进行负载均衡的配置。也可以通过ZooKeeper在客户端进行负载均衡配置。

1. 多个服务注册
2. 客户端获取中间件地址集合
3. 从集合中随机选一个服务执行任务

**ZooKeeper负载均衡和Nginx负载均衡区别**：

- **ZooKeeper**不存在单点问题，zab机制保证单点故障可重新选举一个leader只负责服务的注册与发现，不负责转发，减少一次数据交换（消费方与服务方直接通信），需要自己实现相应的负载均衡算法。
- **Nginx**存在单点问题，单点负载高数据量大,需要通过 **KeepAlived** + **LVS** 备机实现高可用。每次负载，都充当一次中间人转发角色，增加网络负载量（消费方与服务方间接通信），自带负载均衡算法。

### 9.5 命名服务

![img](https://fastly.jsdelivr.net/gh/wxler/cdnPicture/imgs/20210309181007.png)

命名服务是指通过指定的名字来获取资源或者服务的地址，利用 zk 创建一个全局唯一的路径，这个路径就可以作为一个名字，指向集群中某个具体的服务器，提供的服务的地址，或者一个远程的对象等等。

阿里巴巴集团开源的分布式服务框架 Dubbo 中使用 ZooKeeper 来作为其命名服务，维护全局的服务地址列表。在 Dubbo 的实现中：

- 服务提供者在启动的时候，向 ZooKeeper 上的指定节点`/dubbo/${serviceName}/providers` 目录下写入自己的 URL 地址，这个操作就完成了服务的发布。
- 服务消费者启动的时候，订阅`/dubbo/${serviceName} /consumers` 目录下写入自己的 URL 地址。

注意：所有向 ZooKeeper 上注册的地址都是临时节点，这样就能够保证服务提供者和消费者能够自动感应资源的变化。

另外，Dubbo 还有针对服务粒度的监控，方法是订阅`/dubbo/${serviceName}` 目录下所有提供者和消费者的信息。

另外，**分布式锁和选举也是Zookeeper的典型应用场景**。

【参考资料】

1. https://blog.csdn.net/weixin_44766402/article/details/92682593
2. https://blog.csdn.net/u013679744/article/details/79222103
3. https://blog.csdn.net/u013374645/article/details/93140148
4. https://dbaplus.cn/news-141-1875-1.html
5. https://mp.weixin.qq.com/s/b5mGEbn-FLb9vhOh1OpwIg
6. https://mp.weixin.qq.com/s/W6QgmFTpXQ8EL-dVvLWsyg
7. https://snailclimb.gitee.io/javaguide/#/docs/system-design/distributed-system/zookeeper/zookeeper-plus
8. https://snailclimb.gitee.io/javaguide/#/docs/system-design/distributed-system/zookeeper/zookeeper-intro
9. https://blog.csdn.net/u013679744/article/details/79222103
10. https://www.cnblogs.com/raphael5200/p/5285583.html
11. https://www.zhihu.com/question/20004877
12. https://mp.weixin.qq.com/s/tInsv8fRVT1a-lE-uS1QEw
13. https://zhuanlan.zhihu.com/p/158566353
14. https://www.runoob.com/w3cnote/zookeeper-session.html
15. https://segmentfault.com/a/1190000022193168
16. https://mp.weixin.qq.com/s/8yHbaEsyiY1EdKenVcVEYA
17. https://mp.weixin.qq.com/s/Ybt7M_uichWg5YX-9-tyhw
18. https://mp.weixin.qq.com/s/4kjMJ0IKXP9T-cRCrf3uYQ
19. https://mp.weixin.qq.com/s/HWjynP8-777EltcQtm3e3Q
20. https://mp.weixin.qq.com/s/u8QDlrDj3Rl1YjY4TyKMCA
21. https://mp.weixin.qq.com/s/Gs4rrF8wwRzF6EvyrF_o4A
22. https://blog.csdn.net/yangguosb/article/details/80254240
23. https://www.cnblogs.com/tommyli/p/3766189.html
24. https://zhuanlan.zhihu.com/p/67654401