


一个独立管理的 D N S子树称为一个区域 ( z o n e )。

一旦一个区域的授权机构被委派后，由它负责向该区域提供多个名字服务器

当一个新系统加入到一个区域中时，该区域的 D N S管理者为该新系统申请一个域名和一个 I P地址，并
将它们加到名字服务器的数据库中

一个名字服务器负责一个或多个区域


问题:

具体解析的顺序: 本地没有命中，直接找到根还是逆流一级级往上找。

实践递归和非递归查找，用工具。


If you add the +trace variable, dig can also perform a recursive lookup of a DNS record: 

https://aws.amazon.com/premiumsupport/knowledge-center/partial-dns-failures/


You can also perform a query that returns only the name servers:
dig -t NS www.amazon.com

 dig -t NS www.nearmap.com 并没有返回NS


 whois: https://dig.whois.com.au/whois/nearmap.com


 With no server specified, dig will query the DNS server configured on the system where you are running the command.


 tools like dig, nslookup, whois and host can be used t determine the authoritative DNS servers

 dig +short NS digitalinternals.com


 can start tracing the DNS server records all the way from the TLD (top-level domain)

 *  the TLD is com. First, get the SOA record for com TLD
    ```
    dig +short SOA com
    ```

一个名字服务器负责一个或多个区域

并不是每个名字服务器都知道如何同其他名字服务器联系。相反，每个名字服务器必须知道如何同根的名字服务器联系

1 9 9 3年4月时有8个根名字服务
器，所有的主名字服务器都必须知道根服务器的 I P地址（这些 I P地址在主名字服务器的配置
文件中，主服务器必须知道根服务器的 I P地址，而不是它们的域名） 


详细的DNS解析过程， NS和SOA区别的分析: 

[What is the difference between SOA record and NS record? Does NS record help the resolver to identify the ipaddress of domain without requesting the root server?](https://www.quora.com/What-is-the-difference-between-SOA-record-and-NS-record-Does-NS-record-help-the-resolver-to-identify-the-ipaddress-of-domain-without-requesting-the-root-server)

cache: 第一次需要时间，recursive, 第二次0ms


```
dig +norecurse skymap.cn
dig skymap.cn
```

C N A M E 这表示“规范名字 (canonical name)”。它用来表示一个域名（标识符串） ，而
有规范名字的域名通常被称为别名 ( a l i a s )。

[权威 DNS 和递归 DNS](https://www.alibabacloud.com/help/zh/doc-detail/60303.htm)


递归查询例子很详细 例解DNS递归/迭代名称解析原理(https://blog.csdn.net/lycb_gz/article/details/11720247)