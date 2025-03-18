为了平时工作的方便写了个(用了wails框架)
功能展示:
![image](https://github.com/user-attachments/assets/5ec9afb4-0e5b-403a-b8e2-44d06f8f3741)

实现mysql的会话管理,这样不用每次都去show processlist,而且如果是同样的sql,还能折叠起来,更方便查看问题.还可以选中直接杀掉
![image](https://github.com/user-attachments/assets/38fe80b2-94a8-49e0-8734-f610f74ae55d)
检测锁的界面,可以看到一个锁会这样被展示,可以清晰的看到是谁阻塞了谁,可以直接选中干掉阻塞的,或者干掉被阻塞的.后续会考虑加上锁的关键参数,可以方便dba处理故障.
![image](https://github.com/user-attachments/assets/f35758c9-ebea-44de-84db-1b36efc8266b)
![image](https://github.com/user-attachments/assets/ba91c67c-8a9b-494f-94ce-136778942648)
![image](https://github.com/user-attachments/assets/5ada8d24-c5db-4bb7-8855-310187018b35)
这个是redis的功能模块.可以分析redis的rdb文件.更方便查看问题.而且分析完后分析结果还能根据分析的日期默认保存起来,可以和前几天的分析结果进行比较.
而且可以点击柱状图,查看具体是哪一类的key或者某一些具体的key占用了内存
