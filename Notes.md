# 多个文件SQL语句交错执行

## 问题

多个SQL文件中的语句交错执行，并且要求所属同一个文件的SQL的执行顺序和文件中SQL语句顺序一致。

### 问题的数学模型

- 假设有n个SQL文件，记为S1、S2，S3 ... Sn, 每个SQL文件中的语句数量分别为：M1, M2, M3 ...Mn, 所有SQL文件中总的SQL数量记为T, T=(M1+M2+M3..._+Mn)
- n个sql文件，对应n个数据库客户端

- "n个SQL文件中语句交错执行，并且保持单个文件中SQL语句顺序一致"，问题可以简化为: 
    1. 长度为T的序列，先把其中M1个位置标识为S1, 再M2个位置标识为S2, ... 最后Mn个位置标识为Mn.
    2. 从序列T开始遍历，遇到标识为Sx的结点，就从Sx文件中顺序取出一条SQL进行执行，后续遇到同样Sx的结点，又从Sx的文件中顺序取出下一条SQL语句。通过该方法保持同属一个文件的SQL语句顺序一致。
    3. 执行文件Sx的SQL语句时，通过文件对应的第x数据库客户端执行。
    
 - 所以：问题的第一步是怎样完成长度T的序列，M1个位置标识为S1, M2个位置标识成S2, ... Mn个位置标识为Mn， 数学模型可以理解为：
    1. 从长度为T的序列，先选择M1个结点，再在剩余的结点中选择M2个，再在剩余的结点中选择M3个, ... 剩余的结点Mn个选择Mn个。
    2. 所以，枚举数列的排列组合，排列枚举总数量是：C(T,M1)*C(T-M1, M2)*C(T-M1-M2, M3)*C(Mn, Mn), 
    3. 比如：有3个文件，每个文件中SQL语句分别为：1, 2, 3, 则总的排列数量为：C(6,1)*C(5,2)*C(3,3) = 6 * 10 * 1 = 60 

### 实现T长度序列，选择M1个结点标识成S1 ... Mn个结点标识成Mn 

- 因为实现排列组合算法 C(A,B) 采用的是递归实现，如果进行一边串C(A,B)*C(A-B,D)*C(A-B-D, E),则递归深度会非常大，可能会造成调用栈溢出。
- 实现则通过goroutine和channel, 单个goroutine实现单个C(X, Y)的算法，结合channel形成pipeline，串联实现"(A,B)*C(A-B,D)*C(A-B-D, E)"的算法, 避免调用栈过深的问题。

## 实现

- 为了加快整体的SQL执行速度，当通过"Pipeline"输出一条排列组合时，则通过新的goroutine执行该排列组合对应的SQL语句。
- 为提升程序的扩展性、测试性，通过uber的dig库进行动态注册，并且struct的依赖均采用interface.
- 如："sql_disorder_executor --conf ./config/config.json mock" 的子命令，dig注册Mock的数据库验证程序执行的正确性。 而"sql_disorder_executor --conf ./config/config.json"会连接数据库进行执行。

### 待优化点

- 计算过程中的"组合序列"采用数组的形式，后续可以采用BitMap的形式进行优化。
- 程序的优雅退出，现在程序执行过得中生成"组合序列"并创建新的goroutine执行"组合序列"对应的SQL语句，主进程等待所有的goroutine执行完毕后，才能退出，并没有超时退出机制和通知goroutine退出的机制。
- 完善的单测。

## 运行

- 启动TIDB 
```
git clone https://github.com/pingcap/tidb-docker-compose.git
cd tidb-docker-compose && docker-compose pull
docker-compose up -d
```

-. Create Table

```
using test; 
CREATE TABLE `x` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `a` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4;
```

-. Run the app with MockClient

```
go build
./sql_disorder_executor --conf ./config/config.json mock > runlog.log 
```

-. Run the app with DB

```
go build
./sql_disorder_executor --conf ./config/config.json > runlog.log 
```

-. Verify the result

```
grep "loop_1," runlog.log --color
```