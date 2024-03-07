# 分布式对象存储项目

## ✨ 项目功能

✅🚀 分布式存储，自由扩展节点

✅⚡ 对象存储

✅📚 文件压缩传输

✅⚡ 大对象断点续传

✅💾 多版本存储

✅👁️‍🗨️ 数据校验和去重

✅🔗 数据冗余和数据恢复

✅📚 压缩存储数据

## 📃 TODO

- 👩‍👧‍👦 多用户功能，**MySQL**存储数据
- ☁️ 本机、阿里云 OSS、腾讯云 COS等作为存储端，使用微服务框架，**consul** 进行服务发现
- 🗃️ 使用 **Raft** 保证数据的一致性
- 🌈 使用 **ceph**
- 🚀 **kubernetes** 部署，基于 **gitlab+Jenkins+Harbor** 自动化部署

## ⚗️ 技术栈

- [Golang](https://golang.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [ES](https://www.elastic.co/cn/elasticsearch/)
- RS纠删码
- [Docker](https://www.docker.com/)

## 🛠️ 部署

通过 Docker 部署

## 📜 License

GPL V3

## 📗实现文档

### 单机版对象存储系统

对象存储是以对象的方式来管理数据的（对象数据+元数据+ID），通过REST网络服务来访问对象。提升了储存系统的扩展性

首先是处理路由，这一章主要是实现了两个REST网络接口：

- 一个是`PUT`请求，将数据保存到服务器中，通过 `io.Copy()`将想要储存的内容复制到文件即可
- 另一个是`GET`请求用来下载对象。在服务器中寻找对象，同样通过`io.Copy()`将服务器该对象的数据写入到`HTTP`响应体中

### 可扩展的分布式系统

分布式系统的好处在于可扩展性，只需要加入新的节点就可以自由扩展集群的性能。

本章主要是实现了分布式，从整体上来看，就是接口服务节点和数据服务节点

- 接口服务结点主要的作用是：转发`PUT`和`GET`请求到数据服务节点
- 数据服务节点才是真正处理`PUT`和`GET`请求

另外为了验证系统，还加入了`locate`接口服务，用来定位对象资源。具体实现如下：

首先接口服务收到"`/locate/文件名`"的请求，会通过Locate方法获取数据服务的监听地址，先将文件名通过交换机 dataServers 发送到消息队列，数据节点收到之后（Consume）通过Send方法将自身的监听地址发送给自己的消息队列，然后接口服务通过Consume方法收到这个监听地址。

数据服务节点和接口服务节点主要通过消息队列`RabbitMQ`进行通信，具体使用在于

- `object`包：接口服务节点转发对象

  - 首先是`GET`请求，先通过`locate`方法获得对象保存的数据节点位置（监听地址），然后通过`Http.Get(url)`获得响应即对象内容，此时的`URL`的服务器是数据服务节点
  - 其次是`PUT`请求，先随机选择一个数据服务节点（可以通过遍历哈希表找到随机的一个），然后通过`Http.NewRequest("PUT", "URL", reader)`转发该请求，再通过`Client.Do(request)`获取响应内容（状态码）

  数据服务节点的`object`包和上一章的一样

- `heartbeat`包：数据服务节点每隔五秒发送心跳消息（该节点的监听地址），通过`Publish`发送到`apiServers`（交换机）,交换机绑定(`Bind`)一个消息队列，接口服务节点则通过该消息队列接收数据服务节点的心跳消息，并将监听地址保存到哈希表，每10秒清除没有发送心跳消息的数据服务节点（可能出问题了）

- `locate`包：接口服务节点通过`Publish`将定位消息（对象）发送到`dataServers`（交换机），同样绑定一个消息队列，数据服务节点收到定位消息，会在服务器寻找该对象是否存在，发送反馈消息(`Send`)，存在就会发送该对象保存在哪一个数据节点（监听地址）。

另外是测试的环境

- 通过在一台虚拟机绑定多个`IP`地址实现多个节点。
- 安装`RabbitMQ-server`，创建两个交换机 

### 元数据服务

这一章主要是为了记录对象的不同版本。元数据指的是对象的描述信息，比如说名字，版本，大小以及散列值。

我们用`Elasticsearch`来实现元数据服务，它的索引相当于数据库，类型相当于表，属性相当于列。

与上一章相比，我们的代码主要新增了`es`包，`utils`包。

- `es`包：封装元数据服务的一些`API`，包括查看最新版本，查看一个对象的所有版本，增加一个版本，获取一个版本。因为`Elasticsearch`也相当于一个服务器，对于不同的索引保存的信息在服务器中，可以类比成`MySQL`保存数据
- `utils`包：只要两个函数，一个是获取请求头部的散列值，一个是获取请求头部的长度信息。

主要新增了`version`接口用来查看对象的版本，新增`delete`实现对对象的删除。

下面来说一下对象的`put`和`get`实现细节

- `put`请求：我们可以从请求头部获取它的散列值，长度，名字，这样我们就可以调用`es`包下的`API`来添加一个版本。另外我们发送`PUT`请求的时候需要带上它的散列值。
- `get`请求：与上一章不同的是，我们是将对象的散列值作为全局唯一变量。首先我们很容易获取请求的名字和版本号，根据这两个信息，他们调用`es`包下的`API`来查找它的元数据，这样我们就可以获取它的散列值，用这个散列值替代上一章的名字去访问

- `delete`请求：用于删除一个版本，首先我们通过名字获取它的最新版本，然后增加一个版本，只不过这个版本的长度为0，散列值为空，用来表示这是一个删除标记，实际上是没有删除，我们依然可以通过/version接口查看之前的版本。

另外数据服务存储的不再是对象的名字，而是散列值，这个散列值就是从请求头部获得的。

因为散列值是根据对象内容来计算出来的，所以不同的版本有不同的散列值，以不同的文件存放在服务器中。

想想还觉得非常不合理，哪个傻逼会存放东西的时候，还要根据内容计算出一个散列值，然后再去存呢？



### 数据校验和去重

这一章主要是进行数据校验和去重。这是非常重要的，如果一份数据保存几百份甚至更多会占用我们的空间，这就是去重，数据校验是为了防止恶意客户端错误的信息去请求，我们当然要拒绝，另外也是为了防止数据在传输中发生的错误。

首先是去重，其实实现起来不是很复杂，数据服务启动的时候先将磁盘扫描一遍，将所有的文件名（散列值）缓存下来，以防止定位服务多次访问磁盘，这样可以直接访问内存，提高了速度；之后客户端请求一个数据的时候我们就会去这个缓存中寻找，如果有的话就直接添加一个版本就好了。当然这样做也是有很大的局限性，`PUT`一个对象会等待`1`秒，这是去重导致的性能问题。

另外就是数据校验的实现，这一个相对来说复杂一点，当用户`PUT`一个对象的时候，我们需要根据请求对象的内容来计算出一个散列值，然后和传入的散列值进行比对，成功才会上传这个对象，否则就会删除。这样也会导致一个问题，需要校验的话我们要等待对象上传完之后才能计算出散列值，但是如果对象太大的话，可能超出接口服务的内存。

我们的解决方法是给数据服务增加缓存功能，这样一来接口服务的PUT方法就被新增的`POST`，`PATCH`，`PUT`，`DELETE`四种方法替代了。下面我们详细说明这个缓存的临时接口

- `/temp POST`方法：一开始执行`POST`请求增加一个对象，此时请求的是散列值，响应`uuid`；这里的`uuid`是服务器随机生成的，具体生成方式我有点疑惑。这个`uuid`是临时对象的标识，除此外还会新建两个文件，这两个文件都是以临时对象标识`uuid`进行命名的，一个文件记录数据的信息，包括散列值，`uuid`，`size`，另一个文件是数据文件，此时只是做好准备，没有数据，之后会通过`patch`方法打补丁，将数据写入。
- `/temp PATCH`方法：通过`PATCH`方法将请求的内容写入到数据文件，这时候请求的是通过`POST`方法获得的`uuid`（临时对象标识）。具体实现是，根据`uuid`读取信息文件，找到并打开数据文件，将请求体内容保存到数据文件，这是哈还需要比较数据文件的`size`和信息文件的`size`字段，我不太明白这是为什么？难道还会不相等吗？
- `/temp PUT`方法：用来更新一个对象，将之前的临时文件转正，首先会通过`uuid`读取信息文件，还是一样获取到的`size`比较？删除信息文件，更改数据文件的位置（转正），在缓存中增加自己的文件名。
- `/temp DELETE`方法：用来删除一个对象，通过`uuid`删除信息文件和数据文件。

另外接口服务的`GET`方法有所改变，获取一个对象的时候我们也需要进行数据校验，主要是为了防止存储系统的数据降解导致数据随着时间的流逝而逐渐损坏。具体是这样实现的：根据传入的哈希值找到文件，根据文件内容计算散列值，与传入的散列值校对，成功的话返回文件名，根据文件名将文件内容写入到响应体中。

以上就是数据服务的变化，下面是接口服务的变化

- `GET`方法没有变化，主要变化是`PUT`方法
- 去重和`PUT`一个对象进行的数据校验都是在数据节点完成。当没有找到文件的时候就需要新增一个对象，首先会发送`POST`请求，此时会让数据服务建立信息文件和数据文件以就绪，请求获得`uuid`，然后根据这个`uuid`需要发送`PATCH`请求，将请求体内容写入到数据服务的临时数据文件中，然后通过`io.TeeReader`将请求体和`uuid`，`server`写入`reader`，通过这些内容计算出散列值，然后进行校对，这里我也有点不明白，这样怎么可能散列值会相等？之后比对成功就会发送PUT请求将对象转正，否则删除文件。

以上就是这一章所有的去重和数据服务的实现，但是我们会有一个问题，所有的文件都只有一份，对于保护用户数据来说是十分危险的，所以下一章我们会通过数据冗余来解决这个问题！

### 数据冗余和修复

上一章我们留下了一个数据丢失如何解决的问题，这一章我们主要就是实现数据冗余与修复，抵御数据丢失，我们利用的事`RS`码的`4+2`的数据冗余策略，接下来是具体`PUT`一个对象和`GET`一个对象的细节实现。

先提前对一些知识进行说明，RS码是将一个对象的内容分成4个分片，每个分片都有一个id，另外2个是修复分片，以防数据节点出问题，用来修复分片数据。

首先需要`PUT`一个对象。请求传入的是对象的名字以及`hash`。

流程是这样的

宏观上的函数调用过程 `StoreObject`——>`putStream`——>`NewRSStream`——>`NewTempPutStream`

我们可以从一步一步看，`putStream`流主要是获得六个数据服务的节点，然后通过`NewRSStream`流发送六个请求，这个时候的`hash`是`hash.id`，最后是调用`NewTempPutStream`发送post请求到数据服务，这时候会建立6个信息文件和6个数据文件。此外`post`请求也会响应`servers和uuids`，这个会通过`writer`流写入`RSStream`进行编码。这样一来在`StoreObject`中有一个`io.TeeReader`会将请求体写入到`stream`中，而这个`stream`就是`RSStream`，也就是会调用`write`方法，通过循环每一次写入32000字节的文件到缓冲中，当缓冲满了之后会调用`flush`方法，这个方法会对数据内容进行拆分，并且实现奇偶校验，接着遍历将六个分片分别调用`tempstream`的`write`方法，发送`patch`请求，写入到数据文件中。最后验证散列值，成功的话就将其转正，本质上是发送`put`请求，这个时候数据服务会通过`uuid`找到信息文件，信息文件保存了每个分片的`name（hash.id）`和`size`，然后找到数据文件，根据内容可以计算出这个分片的`hash`值，最后将这个临时对象重新命名变成 `hash.id.<shard of hash>`

- 数据服务
  - 一个变化就是对接口服务的定位信息的反馈，之前只是反馈一个节点地址，但是现在需要反馈所有节点监听地址及相应的分片`id`，我们用一个结构体来将这些数据反馈给接口服务。
  - 每个分片存储的文件格式为`hash.id.<shard of hash>`, 所以我们缓存的节点监听地址也有变化，值变成了对应的分片`id`，还是一样从所有文件进行寻找，将分片id和监听地址一一对应。
  - 最后一个变化就是将对象转正的时候，因为我们存储的格式变化，我们需要根据每个分片内容计算它的散列值，然后进行转正重命名，将其加入到缓存中。
- 接口服务
  - `locate`方法发生了改变，需要接收到六个分片id以及对应的监听地址。
  - 之前是随机选择一个数据节点，现在变成了选择六个数据节点，因为每个节点都需要存放分片数据，通过`storeObject`函数转发`post`请求，这时候有六个数据节点，每个节点都需要返送`post`请求，响应返回`uuid`，`server`；这些数据都会存放在写入流中，函数返回的是`stream`，这其中包括这个写入流，通过`io.TeeReader(r, stream)` 将内容输出到`reader`中，计算散列值，进行比对，然后转正。

其次是`GET`一个对象。请求体是对象的名字和版本号。

流程是这样的

整体上来看函数的调用过程 get——>getStream——>RSGetStream——>NewGetStream(NewTempStream)

还是一步一步往后看，首先根据请求的名字和版本号通过元数据服务可以获得hash值，然后发送定位消息，这时候会出现两种情况，一种是返回了六个节点的地址，另一种情况是有一个或者两个数据节点的内容损坏了，这时候我们需要对其进行修复，通过随机获得数据节点那个函数，可以获得出错的数据节点地址，接着RSGetStream会发送get请求，如果数据节点没有问题的话，就可以通过reader读取到数据，但是如果有问题那么读取到的内容为空，这时候出现了两种情况，所以针对读取到内容为空的数据节点进行修复，通过NewTempStream重新上传数据，和之前的put一个对象一样，会有一个writer，这样整个流程走完了，执行到了get函数的io.Copy(w, stream)，此时会将stream流获得的数据写入到响应体中，这时候调用RSGetStream实现的read方法，这个方法具体是先看一下缓存中有没有数据，没有的话就会从之前获得数据的reader中取出来，因为有几个为空，我们可以吧id拿出来，然后通过writer把这些数据写入到服务器中，这样一来读取的所有数据就会到响应体中。

- 接口服务
  - 根据请求体获取元数据，从而获得`hash`；然后通过`hash`和`size`定位到每个分片的`id`和`serve`（来自数据服务的反馈信息）；如果我们获取到的`server`不足六个，你们说明有的数据节点数据出现了问题，所以我们需要进行修复，通过选择节点的那个函数修复，下面详细说明
  - `ChooserRandomDataServers`这个函数主要是传入需要的节点数量以及已经有的`server`，然后获得所有数据节点中不包括这些有了的`server`。
  - 以这种方式获取到的`server`数据一般都没有，所以需要进行修复，重新发送`post`请求，修复数据，这些数据在写入流中，而之前存在的数据发送`get`请求获得每个分片的数据文件内容，会放在`reader`中，所以数据要么在`writer`中要么在`reader`中
  - 调用`io.Copy()`方法，将读取到的内容写入到响应体中，内容实现`Read`方法。
- 数据服务：主要是主要`GET`请求中的“`hash.x`", 从所有文件中匹配这样一个文件，只有唯一一个，通过文件内容计算出它的散列值与实际保存的文件名上的散列值作比较，最终将内容写入到响应体中即可

这样实现是比较复杂的，但是提高了整个存储系统的负载均衡。

### 断点续传

这一章主要是解决了由于网络问题导致下载和上传阻断的问题，我们让用户可以从断点开始上传或者下载。

首先是断点上传，因为我们在接收一个对象需要对其进行`hash`验证，所以我们必须使用特定的接口上传。我们采用`post`接口，当用户知道自己上传的是很大的对象的时候，应该主动用这个接口进行上传

下面详细说明断点上传一个对象的流程

前面和`put`请求实现类似，判断`hash`是否存在，如果已经存在不进行后序操作，否则调用`NewRSResumablePutStream`方法，这个方法主要是调用`NewRSPutStream`方法，这其中和上一章一样，每个分片发送post请求，获得putStream流，接着通过`putStream.Writers[i].(*objects.TempPutStream).Uuid`接口回调获得所有的uuid数据，这样再保存到我们的`token`中，`token`的数据包括`name，size，hash，servers，uuids`，这样一来我们直接响应，只不过内容是`token`，当然这个`token`是将这些数据进行编码得到的，解码可以得到其中的数据流。

然后客户端自主改用`post`方法，`temp/token`接口，实现对对象的上传，首先对请求体的`token`进行解码得到数据流，因为我们实现的是断点上传，所以我们需要知道当前已经上传了多少对象了，用`current`标识，具体实现如下

- 通过发送 `head`请求可以获得一个分片的数据文件的大小，乘以4就可以得到已经上传对象的大小了
- 这个`head`请求数据服务实现的，主要是通过传入的`uuid`获取数据文件，然后就可以得到其中一个分片数据文件的大小响应头部给接口服务。

然后请求内容也需要提供一个断点，我们会比较这两个值，肯定是会相等的，除非客户端传错了。然后我们会循环每次写32000字节的文件到缓存中，然后调用`stream`的`write`方法将内容上传到服务器中，这时候还没有转正，因为还没有校验，当所有的内容上传完毕的时候，我们会通过调用`NewRSResumableGetStream`方法，这个方法之后通过发送`get`请求（这个get请求时temp下的）获得所有的数据，通过这个我们就可以计算出真正的`hash`值，与传入的进行比较，相同的话将文件转正即可。

下面是断点下载对象的流程。

断点下载还是比较简单的，只要跳过已经下载的字节就可以了，前面的基本是一样的，唯一不同的是需要从请求中获取偏移量`offset`，然后通过这个`offset`我们就去跳过，调用`Seek`函数， 这个函数是核心，以下是具体实现

- 我们会去循环的跳过，每一次跳过`32000`，通过`io.ReadFull(s, buf)`，这个是核心，`s`是`stream`流，这里作为`io.Reader`接口主要是获取`reader`中的内容，由上一章我们可以知道，这个`reader`就是数据文件的内容读取流，我们把他转移到`buf`中，知道偏移量为0，这时候`reader`中的内容就会减少`offset`个字节。

最后我们再调用`io.Copy(w, stream)`，此时同样是从`stream`中的`reader`读取流中写入到响应，由于之前已经跳过了，所以这时候开始的地方就是我们需要的数据。

### 数据压缩

为了提高传输速度和节省网络带宽，我们对数据进行压缩，使用的事gzip算法。

gzip主要有两个API，一个是NewWriter，创建一个writer对象，用来写入压缩的数据

```golang
w := gzip.NewWriter(f) //f是一个io.writer接口 可以是一个文件
io.Copy(w, oldf)  //将之前文件的内容写入到w中，主要w进行压缩再将压缩后的内容写入到f中

gzipStream, err := gizp.NewReader(f) //可以将压缩文件内容直接以流的形式 写入到响应体中
io.Copy(w, gzipStream)
```

具体有以下几个地方需要进行压缩。

- 首先是put一个对象转正的时候，之前我们是重命名，这时候我们可以对数据进行压缩，然后存储，可以节省存储空间
- 然后就是数据服务响应文件内容的时候，直接使用gzip.NewReader，将压缩文件写入到响应体中。
- 最后就是客户端GET一个对象的时候，可以设置头部Accept-Encoding为true，这样就可以接受压缩文件了，接受过程和第一点相似，通过gzip.NewWriter。

整体来说本章不算很复杂，这样一来我们就节省网络带宽，节省存储

### 数据维护

这一章主要是对数据进行维护。

- 删除过期的元数据，将保留过久的版本删除，保留最近的五个版本
- 删除没有元数据引用的对象数据
- 对象数据的检查和修复

![](https://raw.githubusercontent.com/RobKing9/Blog_Pic/master/Git/20220821204219.png)

# 🤔 常见问题

## 这个项目是如何设计与实现的？

- 一开始我们只是一个单体服务，客户端发送`put`请求存储数据，我们就将其存到磁盘中，客户端发送`get`请求下载数据，我们从磁盘中读取出来给客户端。但是这样我们存在一个问题，当客户端请求骤增，服务器磁盘`IO`负载过高时，都会导致性能下降，并且不好扩展。
- 针对上面的问题，我们将接口服务和数据服务解耦，接口服务只负责接受客户端请求，数据服务只用来请求磁盘，这样一来我们就可以轻松地往集群中扩展新的接口服务节点和数据服务节点，而接口服务和数据服务之间通过消息队列`RabbitMQ`进行信息的传递。但是我们也存在了一个问题，当客户端多次`put`同一个对象的时候，我们在数据服务节点都会存在很多同样的数据，这样非常浪费存储空间，对此我们需要解决数据去重这个问题；但是如果我们`put`同一对象，而每次数据都不一样，这时候我们可以保存对象的多个版本，对数据进项版本控制。
- 我使用的是`ElasticSearch`进行版本控制，它类似于数据库，索引相当于数据库，类型相当于表，每一个属性相当于列。利用ES客户端可以找到指定版本的数据，可以查询所有的版本，ES会保存对象的元数据，包括名字，大小，散列值，这个散列值是客户端通过Sha-256计算出来的，当客户端put一个对象的时候，首先还是会保存在数据服务中，然后会在ES服务器中添加一个版本的元数据，每次版本号加一，当get一个对象的时候，还是一样请求ES服务器返回对应的元数据信息，之后再去请求数据服务返回相应的数据。
- 另外一个就是数据去重的问题我们需要解决。同时因为客户端数据在传输中可能出现数据丢失问题，或者有一些恶意的客户端发送不一致的信息，这时候服务器不能将这些错误的信息保存下来，还有一种情况是服务器因为数据放久了出现数据降解的问题，这时候都需要对数据进行验证，保证接收和发送的数据完整性。解决去重问题，我们可以在接口服务节点转发请求之前，先发送定位信息，通过交换机发送对象数据散列值的信息，数据服务节点会搜索本地磁盘是否存在这个对象，存在的话就会反馈保存的数据节点的监听地址，否则什么也不返回。对数据进行校验我们就需要通过对象的内容计算出散列值，然后和客户端传进来的散列值进行比较，如果不同则拒绝服务，因为我们要接受完数据才能进行计算校验，但是如果文件内容比较大的话，很有可能会超出接口服务节点的内存，所以我们需要将数据转移到数据服务保存到一个临时的地方，当数据校验通过的话，将文件存储到正确的地方，还有就是客户端下载一个对象，我们也需要对取出来的对象进行数据检验。因为有数据降解问题的存在，我们又会想到一个新的问题，如果服务器上一个数据丢失了怎么办？客户端就拿不到数据了
- 对于这个问题，一种方法是保存多份，还有另一种方法就是将一个对象分成很多分片，然后每个分片保存在不同的数据服务节点，我使用的是RS纠删码来设计的，其中有我们将对象分成四个数据分片和两个校验分片，大小都是对象的25%，我们将每个分片保存在六个数据服务节点中，只有其中四个我们就可以还原完整对象，所以我们可以允许两个节点数据出错，这时候用户put一个对象的时候，我们就需要选择六个数据节点，每一个节点进行上面的post，patch，put操作。另外就是如果某个数据节点的数据出现了问题我们还需要对其进行数据修复，当客户端get一个对象的时候，接口服务可以通过心跳机制获得所有数据节点的监听地址，但是如果某个出问题了，我们就可能只能收到其中的五个，那么就需要对另一个数据进行修复，根据RS的原理我们可以很容易进行修复。

## 消息队列是如何设计的？接口节点和数据节点如何交换信息的？

首先数据服务得知道有哪些可用的数据服务节点可以请求，那么数据节点就要发送信息给数据节点，这个过程叫做心跳机制，接口服务节点会绑定一个接口交换机，数据服务每隔五秒给接口服务发送自己的监听地址，接口服务会将这些监听地址按照时间保存到内存中，并且会清除10s没有发送心跳信息的数据节点，因为可能出问题了。这样接口服务转发请求的时候就可以直接选择一个数据节点发送请求

另外接口服务节点还会给数据服务节点发送数据的定位消息，以确定数据在磁盘中的具体情况。这个时候数据服务节点也会绑定一个数据交换机，当数据服务节点收到定位信息的时候，就会反馈这个数据保存在的数据节点地址。

## 心跳机制是如何设计的

数据服务会启动一个协程每隔五秒 通过数据服务绑定的交换机给接口服务发送自己的监听地址，接口服务收到之后会保存收到的时间以及节点地址到哈希表中，每次保存的时候都需要上锁，保证数据的正确性；并且会清除10s仍没有发送心跳信息的数据节点，因为这些节点肯定是出现问题了，同样接口服务也是开辟的一个协程来处理心跳信息。

## 怎么进行版本控制的？

我是通过Elasticsearch实现的，当客户端发送put请求的时候，我们会根据元数据想ES服务器发送put请求保存这个版本，并且每次版本号加一，当客户端需要获取一个版本的元数据信息，直接发送get请求到ES服务器，接着ES服务器会响应对应的元数据信息，如果客户端想知道所有的版本信息，接口服务会请求ES服务器查询所有的版本，并且我们可以让它进行版本号降序排列，这样有时候客户端不指定版本号，我们就可以返回最新的版本信息。

## 怎么对数据进行校验和去重的？

首先是去重，去重意味着如果客户端上传了一个内容相同的对象我们就拒绝保存，客户端put一个对象的时候，接口服务通过消息队列向数据服务节点发送定位消息，数据服务节点会在磁盘中寻找，当然为了提高查找速度，数据服务节点在启动的时候会通过一个协程，将所有的文件名即对象的散列值保存到哈希表中，这样接口服务节点请求定位信息的时候，数据服务节点直接在哈希表中查找是否存在，通过消息队列反馈给接口服务。如果不存在的话才会继续去请求保存。

其次是数据校验，当客户端发送put请求的时候，我们需要通过对象的内容计算出散列值，然后和客户端传进来的散列值进行比较，如果不同则拒绝服务，因为我们要接受完数据才能进行计算校验，但是如果文件内容比较大的话，很有可能会超出接口服务节点的内存，所以我们需要将数据转移到数据服务保存到一个临时的地方，具体我们可以让接口服务节点发送post请求，让数据服务提前做好准备，新建一个信息文件包括对象的散列值和大小，以及一个用来存储数据的文件，当然这时候还是空的，返回给接口服务节点uuid，接着接口服务节点通过uuid发送patch请求，将客户端请求的数据保存到数据服务的临时数据文件中；当所有内容已经保存到临时数据文件的时候，这时候接口服务也通过内容将散列值结算出来了，如果和客户端传进来的散列值相同的话，意味着数据校验通过，紧接着接口服务发送put请求，数据服务节点将文件转正，存储到正确的地方，但是如果没有通过校验的话，接口服务会通过delete请求，数据服务将会删除临时的信息文件和数据文件。还有就是客户端下载一个对象，我们也需要对取出来的对象进行数据检验。

## 怎么对数据进行存储和修复的？

对于这个问题，一种方法是保存多份，还有另一种方法就是将一个对象分成很多分片，然后每个分片保存在不同的数据服务节点，我使用的是RS纠删码来设计的，其中有我们将对象分成四个数据分片和两个校验分片，大小都是对象的25%，我们将每个分片保存在六个数据服务节点中，只有其中四个我们就可以还原完整对象，所以我们可以允许两个节点数据出错，这时候用户put一个对象的时候，我们就需要选择六个数据节点，每一个节点进行上面的post，patch，put操作。另外就是如果某个数据节点的数据出现了问题我们还需要对其进行数据修复，当客户端get一个对象的时候，接口服务可以通过心跳机制获得所有数据节点的监听地址，但是如果某个出问题了，我们就可能只能收到其中的五个，那么就需要对另一个数据进行修复，根据RS的校验片我们可以很容易进行修复。