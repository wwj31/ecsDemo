# ecsDemo
利用actor模型，基于ecs，实现的分布式地图，完成实体基础行为：移动、战斗

# Run(go version ≥ 1.17)
```sh
$ go run github.com/wwj31/ecsDemo@v1.0.1
```

#
### *点击左上角红点，每个区域随机生成移动实体，相互靠近的实体会发生攻击并且扣除生命

- 实体进入不同的地图区域，会用不同的颜色区分
- 一个地图代表一个actor,实体会在地图之间进行无缝切换
- 实体在两区域边界，会产生副本，利用副本的移动、战斗行为和主体一致
# <img align="right" src="https://github.com/wwj31/ecsDemo/raw/master/assets/demo.jpg" alt="map Demo" title="map Demo" />


author: 229482191@qq.com