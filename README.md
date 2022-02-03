# go-vlfeat
[vlfeat](https://github.com/vlfeat/vlfeat) 的 go 语言版本实现

目前实现的版本为 vlfeat v0.9.21

### example:

```
dsift := vlfeat.NewDsift(imgWidth, imgHeight)
defer dsift.Delete()
dsift.Process(imgData)
keypoints := dsift.GetKeypoints()
descriptors := dsift.GetDescriptors()
```

### windows 安装
1. 下载 vlfeat 0.9.21 版本的二进制包，[点此下载](https://www.vlfeat.org/download/vlfeat-0.9.21-bin.tar.gz)
2. 解压 `vlfeat-0.9.21` 文件夹到 c 盘即可
