# ecsDemo
利用dogactor，基于ecs，实现的分布式地图、实体移动、战斗的Demo
##直接编译运行

#

###*点击左上角红点，每个区域随机生成N个随机移动的实体，相互靠近的实体会发生攻击并且扣除生命，生命为0，删除实体*
- 实体进入不同的地图区域，会用不同的颜色区分
- 一个地图代表一个actor,实体会在地图之间进行无缝切换
- 实体在两区域边界，会产生副本，利用副本的移动、战斗行为和主体一致
# <img align="right" src="https://github.com/wwj31/ecsDemo/raw/master/assets/demo.jpg" alt="map Demo" title="map Demo" />
