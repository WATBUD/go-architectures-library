package main

/*1.進程 Process 進程是指被執行且載入記憶體的 program) 調度者:CPU 內核
 Process LIFECYCLE
 new (新產生)：該進程正在產生中
 ready (就緒)：該進程正在等待 CPU 分配資源，只要一拿到資源就可以馬上執行
 running (執行)：該進程取得 CPU 資源並且執行中
 waiting (等待)：該進程在等待某個事件的發生，可能是等待 I/O 設備輸入輸出完成或者是接收到一個信號，也可以想成是被 block (阻塞) 住
 exit (結束)：該進程完成工作，將資源釋放掉
2.線程 Thread(線程存在於 process 裡面,thread 是作業系統進行運算排程的最小單位,進程像一個工廠，線程則是工廠裡面的工人) 調度者:CPU 內核
3.協程 Coroutine(也會有自己的 registers、context、stack) 調度者:使用者控制


非同步(Asynchronous)
同步 (Synchronous)
台灣與對岸的名詞對照表
process:
台灣：程序、處理程序
對岸：進程
thread:

台灣：執行緒
對岸：線程
concurrent:

台灣：並行
大陸：並發
parallel:

台灣：平行
大陸：並行*/
