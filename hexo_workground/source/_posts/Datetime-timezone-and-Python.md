---
layout: post
title: 'Datetime, timezone and Python'
date: 2021-06-16 06:32:37
tags:
---

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210618084911.png" alt="20210618084911" style="width:500px"/> 

## Preface

写server程序，datetime是绕不开的话题，且经常需要处理多时区转换，时间的术语和格式让人眼花缭乱, 如

*  UNIX time, epoch
*  datetime string: `1990-12-31T23:59:60Z`, `1990-12-31T15:59:60-08:00`
*  UTC, GMT, Zulu time
*  DST: daylight saving time
*  Timezone name 如`AEST`, `AEDT`
*  Timezone name又如 "Australia/Sydney"
*  时区与longitude的关系
*  ISO 8601, RFC 3339

需要系统梳理时间标准，术语; 这就需要对背后的地理知识稍作了解. 

另外给出如何用Python的`Arrow` package实现常用的时间/时区操作. 

<!-- more -->

## Background

生活总我们提到的时间多为当地时间: 如写本文时，身在悉尼，当地时间为: `Thu 17 Jun 2021 07:58:48`; 但在同一时刻，北京时间为`Thu 17 Jun 2021 05:58:48`, 即北京时间落后于悉尼时间2h. 

日常电子设备，操作系统的时间一般设置为当地时间，方便使用。可见，除了具体的Year, Month, Day, Hour, Minute; 我们还需要location信息，才能唯一的标定一个时间。

## 无歧义的时间: UTC, GMT, UNIX time

为了解决同样的datetime string在不同location, 代表的绝对时刻不一致的问题; 同时为了方便我们在不同location之间转换，我们需要有一个"标准"的location作为参照，即UTC时间(后续阐述其由来)。

只要知道身处何地，当地时间都能转换为UTC时间; 这样有了一致的标准，不同location的当地时间就能互相转换了。

例如: 北京时间比UTC时间提前8hr, 悉尼时间比UTC时间提前10hr(当前为冬令时, 后续会讲到), 因此悉尼时间比北京提前2hr.

Note: 

*  时间提前/落后: 这里表示同一时刻，北京时间的显示比悉尼时间小/落后，北京时间落后于悉尼时间2h, 针对同一时刻. 即从悉尼打电话到北京，发现北京时间"慢"2小时
*  以上又意味着北京将比悉尼晚2小时见到太阳

UTC时间取自Greenwich时间，位于伦敦; 因此Greenwich Mean Time与UTC相等; 军用的Zulu time也与UTC相等, 即: 

> UTC == GMT == Zulu time

编程时常用的UNIX Time也是基于UTC的: 

> Unix time (also known as Epoch time, POSIX time, seconds since the Epoch or UNIX Epoch time)

> It is the number of seconds that have elapsed since the Unix epoch, minus leap seconds; the Unix epoch is 00:00:00 UTC on 1 January 1970

可见，UNIX time可以无歧义的表示某一时刻，因为它的标准化属性，计算机系统内部存储常用它. 

## Timezones相关地理知识

如果想进一步理解UTC的由来，标准时区的划分，需要review一些地理知识: 

首先地球在不停的自传, 考虑自转时太阳可视为基本不动

<img src="https://upload.wikimedia.org/wikipedia/commons/3/32/Earth_rotation.gif" alt="20210617075335" style="width:500px"/>

自转时时而面对太阳，时而背对太阳; 因此地球上的每个location都会依次经历 日出->日落->日出->..., 循环往复. 

时区的划分与自转紧密相关。还记得地理课的"口诀"吗? 地球自西向东转动; 东西是怎么来的呢？ 东西方向感觉是相对的，为什么一个地方的东西南北方向是确定的，没有歧义呢? 

首先东西方向随地球的自转: **东代表日出，西代表日落**, 即太阳从东方升起，到西方落下，这种约定在任何地方都是一样的;  

东西方向确定后，南北自然确定，同样遵循约定: : 

> east being in the clockwise direction of rotation from north and west being directly opposite east

即East逆时针转90度就是North, 这就是[Cardinal direction](https://en.wikipedia.org/wiki/Cardinal_direction) , 这也是地图口诀"上北下南，左西右东"的由来. 东西定为水平线, X轴; 则南北自然与其垂直的Y轴。

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617084340.png" alt="20210617084340" style="width:500px"/>  

那为什么地球是**自西向东运动**呢? 下图给出解释:  

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617084457.png" alt="20210617084457" style="width:500px"/>  

即太阳不动，地球自转，不停的有location见到太阳(日出), Syndey比Beijing先见到，我们就说Sydney在Beijing的东边; 因此自转方向也是自西向东. 

上图来自一个很好的回答[Why do we say Earth rotates from west to east?](https://earthscience.stackexchange.com/a/13361)

另外理解了东西, 顺便看下Clockwise和counter clockwise; 地球是否顺时针转取决与观察角度，与"自西向东"并无直接关系: 从北极看上去，是逆时针转动; 钻到下面，从南极往"上"看，是顺时针转动: 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617085332.png" alt="20210617085332" style="width:500px"/>  

附一张地球转动图:

<img src="https://upload.wikimedia.org/wikipedia/commons/6/61/AxialTiltObliquity.png" style="width:500px"/>  


## Timezone的划分

有了以上的地理知识，下图开始make sense: 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617075335.png" alt="20210617075335" style="width:500px"/>  

地球经度360度，一天24小时; 为方便时间计算: 

*  设立24个时区，每个时区跨15度
*  每个时区对应1小时的时间跨度: 这样环绕地球一周，正好是一天24小时. 
*  地处同一时区内(15度经度带), 认为时间一致， 为时区标准时间, 即时区中部经度时间为准; 即时区开始经线的7.5度偏移处。

### 时区与经度

时间与太阳升起有关，因此只与经度: longitude有关. 

Prime meridian, 即本初子午线，代表0度经度; 同时将Prime meridian的时间定义为UTC标准时间：即**经度的0与时间的0重合**.

时区0因为中心是longitude 0, 因此纬度范围[-7.5 7.5]; 

北京纬度大约为`116.3`, 落在东八区: [112.5, 127.5]

提到经度就不得不提东西半球的划分: 有了prime meridian, 向东(即地图向右)180度，为东半球; 向左180度为西半球. 东西半球重合于两条经线: prime meridian和antimeridian(位于Pacific Ocean, 线上很少有人居住, 因此用作国际日期变更线); 

这大概就是"东方国家"，"西方国家"的由来吧; 虽然地球是个球体，没有绝对的谁先看到太阳，谁后看到; 但以标准0度经线来看，处于东方的中国先见到太阳，因此时间领先于处于西方的美国。而北京(东八区)和纽约(当前西四区)横跨12个时区，正好相差12小时, 昼夜颠倒, 当前:  

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617092303.png" alt="20210617092303" style="width:500px"/>

### 经度取值范围

有了Prime meridian, 定义经度范围[-180, +180]: 

> Positive longitude is east of the prime meridian, while negative longitude is west of the prime meridian

E.g: 北京处于东区, 经度约`116.3`, 为正; 纽约处于西半球, 约`-73.9`, 为负. 画成坐标，0度经线为Y轴: 0度经线将地球分为东西半球

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617093049.png" alt="20210617093049" style="width:500px"/>

另一种经度表示, 如 `73° 56' 6.8712''W`, 用`W`代表西半球, 数值永远为正: [0, 180]; 一般地图程序还是用[-180, 180]的多. 

维度取值为[-90, 90], 分不清数字对是经纬度还是维经度, 看哪个数字大于90即为经度(当然不是任意lat, lon pair都能分清);   

一般地图多用`lat/lon`, 如Google地图Beijing坐标`@39.9385466,116.1172815`. 

### 国际日期变更线

时区地图是二维平面地图，但地球是球形; 意味着除了Prime Meridian, 还有一处东西时区会"重合". 即Timezone地图的"东西十二区". 

它是最东的时区，比UTC时间早12小时; 同时它也认为是最西时区，比UTC晚12小时; 如何对待: 分半, 靠东的认为属于东半球的延伸， 领先UTC 12小时; 靠西的认为属于西半球, 落后UTC 12小时; 

即同样处于一个15度经度带的不同地方，日期相对标准的UTC可能相差一天! 即有一条虚拟的线存在，两边日期不同. 这正是国际日期变更线: `International Date Line`

为了使不同行政地区不至于相差一天，此线并不是严格划在经度带的中心，而是尽量绕开了不同的行政地区. 

> The International Date Line, established in 1884, passes through the mid-Pacific Ocean and roughly follows a 180 degrees longitude north-south line on the Earth. It is located halfway round the world from the prime meridian—the zero degrees longitude established in Greenwich, England, in 1852.  

上述来自[noaa](https://oceanservice.noaa.gov/facts/international-date-line.html)

是一条弯曲的线: 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210617162023.png" alt="20210617162023" style="width:500px"/>  

因此为了简化，东西半球划分更新为: 

> 	area of the Earth east of the prime meridian and west of the International Date Line.

即Prime meridian和国际日期变更线划分. 

## Time standards

有了上节讨论了地理基础, 我们知道地球分为24个时区, 每个时区的标准时间对应为中央经线的时间. 

有了标准的时区划分，对于时间日期的文字表示，还要将其标准化，有几个标准:  

*  ISO8601
*  RFC 2822
*  RFC 3339

> Standards for date and time representation covering the formatting of date and time, and timezone information

以上3个标准, 即标准化日期，时间和时区表示; 相对于有歧义，不规范的string: `2020-09-12 5:20`. 甚至`199202930201`; 我们需要定义标准的格式. 

[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339)是对前两者的完善, 接下来主要讨论它. 

## RFC 3339

直观的看几个符合标准的datetime: 

```sh
# UTC时间, 以Z结尾
2019-10-12T07:20:50.52Z
# AEST时间, +10:00 表示比UTC提前10小时
2021-06-17T16:40:02.058752+10:00
# UTC - 8, PST: Pacific Standard Time
1996-12-19T16:39:57-08:00 # 等价于悉尼时间1996-12-20T11:39:57+11:00, 注意此时间为+11:00, 即悉尼施行夏令时, 不是AEST而采用AEDT.
```

可见: 

*  日期以`-`连接, Year: 
   +  Year: 4 digits
   +  Month: 2 digits
   +  Day: 2 digits
   +  Hour: 2 digits, [00, 23], 不允许24
   +  Min: 2 digits
   +  Sec: 2 digits; 小数点后经度不限
*  日期时间必须加`T`
*  时间最后必须加时区信息: 
   +  `Z`: Zulu time, == UTC == GMT, 即`UTC+0`, 或者
   +  `+10:00`: UTC + 10hr == 当前时间戳; 即当前时间比UTC领先10hr.

ISO8601严格的`ABNF`在此: 

```
date-fullyear   = 4DIGIT
date-month      = 2DIGIT  ; 01-12
date-mday       = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on
                           ; month/year
time-hour       = 2DIGIT  ; 00-23
time-minute     = 2DIGIT  ; 00-59
time-second     = 2DIGIT  ; 00-58, 00-59, 00-60 based on leap second
                           ; rules
time-secfrac    = "." 1*DIGIT
time-numoffset  = ("+" / "-") time-hour ":" time-minute
time-offset     = "Z" / time-numoffset

partial-time    = time-hour ":" time-minute ":" time-second
                  [time-secfrac]
full-date       = date-fullyear "-" date-month "-" date-mday
full-time       = partial-time time-offset

date-time       = full-date "T" full-time
```

Note: 

*  `T`, `Z`可为小写
*  RFC3339, 不一定非要求`T`, 可以为空格` `; 这也是RFC3339的format与ISO8601的唯一区别. 

## Timezone name 与 Daylight saving

时区虽然有了数字编号, 如UTC+8, UTC-10; 

但人们更熟悉名字，因此IANA标准: [Time Zone Database](https://www.iana.org/time-zones) 是将时区命名标准化: 用大城市来代替所属的时区, 在东八区，即我们的北京时间, 北京城市并没有上榜, IANA对应的时区是: `Asia/Shanghai`, 即`CST`: `China Standard Time`. 但我们应该知道它和北京时间是一个时间. 

Linux下列出所有IANA timezone names: 

```sh
timedatectl list-timezones # 当前共计348个
```

可以看到悉尼时区为: `Australia/Sydney`, 这个和之前提到的`AEST`(Australian Eastern Standard Time)和`AEDT`(Australian Eastern Daylight Time)又有什么关系呢?

这就是Daylight saving. 是当地政府的行为, 一般有季节性. 

中国并没有采用Daylight saving机制, 比如一个人起床时间固定在6:00AM, 夏天天已经亮了，但冬天仍很黑. 冬天时钟往后拨一小时(将夏季的7:00AM当做冬季6:00AM), 能让冬天遵循同一时间routine的人多享受一会太阳. 

Daylight saving提供额外的好处，势必会带来额外的cost: 同一timezone: 如`Australia/Sydney`, 冬季夏季对于标准UTC的offset不一致, 冬季采用AEST(UTC +10), 夏季采用AEDT(UTC +11), 因为 `localTime = UTC + offset`, `AEDT`的`+11`, 比`AEDT`领先1小时, 即更早1小时. 

在悉尼，一年有两天变更时间, 家里有挂钟可能需要手动调整(手机软件一般会自动变更)

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210618071755.png" alt="20210618071755" style="width:500px"/>

因此只提到IANA的标准timezone, 如`Australia/Sydney`， 需要指定日期才能知道对UTC的offset, 但`AEST`和`AEDT`没有歧义，它们能得出UTC offset.

```sh
# 悉尼，采用AEST时间, 冬季, +10:00 表示比UTC提前10小时
2021-06-17T16:40:02.058752+10:00
# 悉尼, 采用AEDT时间, 夏季
1996-12-20T11:39:57+11:00
```

## Datetime in Python

个人推荐用[arrow](https://arrow.readthedocs.io/en/latest/)处理datetime, Python自带了一堆package, 用起来比较乱. 

Arrow只需要一个package, 且提供了更强大的功能. 

使用时，一般先创建`arrow.Arrow` object, 为无歧义的datetime(timezone awareness, 不允许不带有歧义的时间, timezone不指定默认UTC), 然后用此object进行datetime计算，且可以转为Python自带的datetime, tzinfo等object.

Demo代码见[这里](https://github.com/eliteGoblin/code_4_blog/blob/master/datetime_python/arrow_demo.py)  

### Construct a Arrow object

```py
# get now of local time, timezone from system, AEST in my case
dt = arrow.now() 
print(dt.isoformat())       # 2021-06-18T07:35:21.865380+10:00
# get now of UTC time
dt = arrow.utcnow()
print(dt.isoformat()) 
# get now of another timezone
dt = arrow.now('US/Pacific')
# get from ISO 8601 string
dt = arrow.get('2013-05-11T21:23:58.970460+07:00')
# get as datetime, UTC
dt = arrow.get(2013, 5, 5)  # 2013-05-05T00:00:00+00:00
# get from UNIX time, to UTC
dt = arrow.get(1623966051)
print(dt.isoformat())       # 2021-06-17T21:40:51+00:00, 默认为Timezone为UTC
# get from UNIX time, to local(AEST in my case)
dt = arrow.get(1623966051, tzinfo='local')
# get from UNIX time, to another timezone, support IANA and abbreviation
dt = arrow.get(1623966051, tzinfo='AEST')
print(dt.isoformat())       # 2021-06-18T07:40:51+10:00
dt = arrow.get(1623966051, tzinfo='Australia/Sydney')
print(dt.isoformat())       # 2021-06-18T07:40:51+10:00
```

### Convert between timezones

```py
dt = arrow.now() 
print(dt.to('UTC').isoformat())                 # 2021-06-17T21:42:46.497456+00:00
print(dt.to('Australia/Sydney').isoformat())    # 2021-06-18T07:42:46.497456+10:00
dt.to(tz.gettz('US/Pacific')) # support Python tz package
```

### Parse and format using custom format

```py
# parse from own format, default UTC
dt = arrow.get('2013-05-05 12:30:45', 'YYYY-MM-DD HH:mm:ss')
print(dt.isoformat())                           # 2013-05-05T12:30:45+00:00
print(dt.format('YYYY-MM-DD HH:mm:ss ZZ') )     # 2013-05-07 05:23:16 -00:00, ZZ代表offset
```

[内置一系列Token来帮助parse及format](https://arrow.readthedocs.io/en/latest/#supported-tokens), 与[Linux strptime](https://linux.die.net/man/3/strptime)的token不一致:  

同时系统另外一些常见的format: 

```
FORMAT_ATOM    FORMAT_RFC1123 FORMAT_RFC822  FORMAT_W3C    
FORMAT_COOKIE  FORMAT_RFC2822 FORMAT_RFC850       
FORMAT_RFC1036 FORMAT_RFC3339 FORMAT_RSS   
```

### Get datetime properties

```py
dt = arrow.utcnow()
# get Python's datetime.datetime
dt.datetime     # datetime.datetime(2013, 5, 7, 4, 38, 15, 447644, tzinfo=tzutc())
dt.tzinfo       # tzutc()
dt.year
dt.date()       # datetime.date(2013, 5, 7)
dt.time()       # datetime.time(4, 38, 15, 447644)
dt.timetz()     # datetime.time(22, 4, 50, 202752, tzinfo=tzutc())
```

### Replace some field and shift 

```py
# replace some fields, not change other attribute
arw = arrow.utcnow()                    # 2021-06-17T08:12:37.986011+00:00
arw = arw.replace(hour=4, minute=40)    # 2021-06-17T04:40:37.986011+00:00
arw = arw.shift(weeks=+3)               # 2021-07-08T04:40:37.986011+00:00
# only change timezone, not update other fields, different with convert
arw = arw.replace(tzinfo='US/Pacific')  # 2021-07-08T04:40:37.986011-07:0
```

### Span and range of times

```py
# 当前time的hour开始及结束, 当前时间: 2021-06-17T22:15:38.970255+00:00
arrow.utcnow().span('hour') # (<Arrow [2021-06-17T22:00:00+00:00]>, <Arrow [2021-06-17T22:59:59.999999+00:00]>)
# 另一种方法是floor和ceiling, 得出一样的span start和end
hour_start = arrow.utcnow().floor('hour')
hour_end = arrow.utcnow().ceil('hour')
# get of range of time spans
start = datetime(2013, 5, 5, 12, 30)
end = datetime(2013, 5, 5, 17, 15)
for r in arrow.Arrow.span_range('hour', start, end):
    print(r)
"""
(<Arrow [2013-05-05T12:00:00+00:00]>, <Arrow [2013-05-05T12:59:59.999999+00:00]>)
(<Arrow [2013-05-05T13:00:00+00:00]>, <Arrow [2013-05-05T13:59:59.999999+00:00]>)
(<Arrow [2013-05-05T14:00:00+00:00]>, <Arrow [2013-05-05T14:59:59.999999+00:00]>)
(<Arrow [2013-05-05T15:00:00+00:00]>, <Arrow [2013-05-05T15:59:59.999999+00:00]>)
(<Arrow [2013-05-05T16:00:00+00:00]>, <Arrow [2013-05-05T16:59:59.999999+00:00]>)
"""
# iterate over range of times
for r in arrow.Arrow.range('hour', start, end):
    print(repr(r))
"""
<Arrow [2013-05-05T12:30:00+00:00]>
<Arrow [2013-05-05T13:30:00+00:00]>
<Arrow [2013-05-05T14:30:00+00:00]>
<Arrow [2013-05-05T15:30:00+00:00]>
<Arrow [2013-05-05T16:30:00+00:00]>
"""
```

## Ref

[Arrow: Better dates & times for Python](https://arrow.readthedocs.io/en/latest/)  
[iso8601: Format a Datetime object: ISO 8601, RFC 2822 or RFC 3339](https://rdrr.io/cran/anytime/man/iso8601.html)   
[RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339)    