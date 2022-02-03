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

