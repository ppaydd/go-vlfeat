# go-vlfeat
[vlfeat](https://github.com/vlfeat/vlfeat) 的 go 语言版本实现

目前版本支持的 vlfeat v0.9.21

建议通过此项目学习 cgo 的基本用法

### example:

```
dsift := vlfeat.NewDsift(imgWidth, imgHeight)
defer dsift.Delete()
dsift.Process(imgData)
keypoints := dsift.GetKeypoints()
descriptors := dsift.GetDescriptors()
```

### 不足之处
- cgo 依赖于 vlfeat 的头文件与动态扩展库，需要指定路径，所以把相关文件暂时放在项目中 
- 需要使用的话建议 fork 此项后自己做一些改动
