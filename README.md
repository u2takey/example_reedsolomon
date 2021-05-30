# example for reedsolomon


```
# 1. 首先我们试试存储的效果
➜ cat ./file//testfile.txt                                                                                                           
this
is
just
a
test
file
content
is
meaningless% 

➜ ./main --file=./file//testfile.txt -encode_path=./encode                 
split data into 5+3=8 shards success.
encode data success.%  

# 这里我们把文件 testfile.txt 分成了 5份，加上3份校验数据进行存储，encode 文件夹为存储的数据内容，我们假设在现实的场景中 7 份数据存储在不同的节点上
➜ ls ./encode           
encode_0  encode_1  encode_2  encode_3  encode_4  encode_5  encode_6  encode_7

# 2. 现在我们(因为节点的磁盘损坏)损失了其中的 任意三份数据
➜ rm ./encode/encode_3 ./encode/encode_5 ./encode/encode_7   

# 3. 尝试恢复数据, 可以看出，数据恢复很成功, encode file 都恢复了，而且 recover 出来的原始数据也没问题
➜ ./main --file=./file//testfile_recover.txt -encode_path=./encode --decode
reconstruct data done
recover file success%    

➜ ls ./encode
encode_0  encode_1  encode_2  encode_3  encode_4  encode_5  encode_6  encode_7

➜ cat ./file/testfile_recover.txt 
this
is
just
a
test
file
content
is
meaningless%


# 4. 删除更多的数据, 再次尝试恢复, 可以发现恢复失败: too few shards given
➜ rm ./encode/encode_1 ./encode/encode_3 ./encode/encode_5 ./encode/encode_7   
➜ ./main --file=./file//testfile_recover.txt -encode_path=./encode --decode    
panic: too few shards given
```