# 公主连结公会战管理（后端）

之前的太拉了，重写

已经开发完成

#### 功能清单

- [x] 实时更新boss状态
- [x] 出刀、尾刀、我进了、我出了、挂树了
- [x] 调整boss状态
- [x] 最新出刀显示

#### 部署

下载正确的可执行文件，需要在文件的同目录下有config，log，db三个目录文件。config下需要有config.toml文件。

需要结合[客户端](https://github.com/NaitoSama/pcrcli)使用。

#### config.toml文件结构

```
[General]
HttpPort = "8081"  # 服务端的端口号
RegisterCode = "peko"  # 用户注册码

[DB]
AdminName = "admin"  # 管理员名字
AdminPasswd = "admin"  # 管理员密码  （仅在生成db文件前可用，如果需要修改，请数据库中修改，或者把数据库删了再启动一次软件）

[Boss]
StageOne = [1,2,3,4,5]  # boss第一阶段血量，顺序为第一个boss、第二个boss...
StageTwo = [10,20,30,40,50]  # boss第二阶段血量，下面以此类推
StageThree = [100,200,300,400,500]
StageFour = [1000,2000,3000,4000,5000]
StageFive = [10000,20000,30000,40000,50000]
StageSix = [100000,200000,300000,400000,500000]
StageSwitchRound = [5,10,14,20,24]  # boss切换阶段的周目数，第一个5代表了boss在第5周目切换为二阶段（5周目也为二阶段）。总共设置了六个阶段，用不上的可以设置为100或更大（但不能太大）。

[ClanBattle]
CanBeUndoRecordsUP = 10 # 对应撤回功能，定义了可撤回的最大的历史记录，设置为10代表只能撤回最新10条记录内的数据。
```

