# PlutoChan
启动Pal4重构计划，本次不只要实现先前python2所实现的全部功能，还要优化sql，增加功能。

本次不使用框架，但未来在其它项目中，在数据模式（数据库）规范化的。

## 进阶
分阶段，扩容，改进，以研究在前后端，微服务化等模式下，如果实现各功能模块，

* 第一阶段已经完成了基本上全部的查询功能，页面也做完了
* 第二是前后端分离，RESTful化，用curl进行全部的测试，**不开发专用前端**
* 第二阶段还要完成另一项功能，就是利用redis，**实现查询结果缓存，外带配置yaml化**
* 第三阶段，利用micro模块，实现项目的微服务化，学习consul的使用，这一段，在额外的分支中实现
* 第四阶段，利用test模块，实现业务逻辑与性能的双重测试，并搭建promethus，进行性能监控

## 数据库
为便于分析，特对各表的字段名进行如下改造

* 数据不动，字段的关联关系也不动，原就为ID的字段名，仍不变
* 将所有ID与名称一一对应的表，此两字段均改为ID和NAME
* 激进一些，将所有的字段名，全部精简，一律改为小写
* 反正在数据库里，每个字段只要在自己的表里没冲突就行
* 每个字段基本都有中文名，也不可能不知道其含义

按照数据库的设计规范，以mysql为例，在windows下忽略全部表名/列名/库名等的大小写，
即使在linux下，也会无视列名的大小写，内置的关键字，全部忽略大小写，无视系统平台，
name,type等又不推荐作为字段名，因为，未来自己写项目，**仍要使用全小写+下划线的模式**

## 文件说明

* 主程序入口：main.go
* 模板目录：templates
* 静态文件目录：static

## 页面
为简单起见，不用登录，首页纯静态，罗列全部可用链接，但仍保留嵌套模板设计，具体页面如下：

bootstrap4的默认字体大小为16，过大，导致页面显示表格内容过少，目前改成10，与本机类似

* equip.html: 装备模板
* index.html: 首页
* layout.html: 总布局页
* magic.html: 仙术模板
* mission.html: 任务模板
* navbar.html: 导航条
* prescription.html: 配方模板
* property.html: 道具模板
* question.html: 问答模板
* stunt.html: 特技模板
* upgrade.html: 升级模板

## 道具细分
道具共分五类：恢复、攻击、辅助、材料、剧情。由于道具数量众多，特将以下几类道具再细分

* 恢复类分为食物与其它类，其中食物以SW开头，其它以DH开头
* 辅助类分为香料与其它辅助类，香料以CX开头
* 材料类分矿石、尸块、其它材料三类
* 矿石，全部以CK开头，同时属性说明为“熔铸、锻冶的材料”
* 尸块，大部分以CS开头，同时全部是“注灵的材料”
* 其它材料，CQ开头，或者属性说明为空

## 道具购买场景

* 无法购买的道具子类：尸块、剧情类
* 无挂载技能类：香料、材料类、剧情类
* 没有价格的：剧情类

## 配方字段分析

* 熔铸没有属性、潜力、灵容量、技能、粒子特效字段
* 锻冶没有合成装备、灵容量、技能、粒子特效字段
* 注灵没有合成装备、潜力字段

## 任务分析
区分三类任务的方法：

* depended_id:1~129主线，200~243委托，300~321支线，估计规划是199内主线，299内委托，399内支线
* trunk:0主线，19~37委托，1~8支线，委托与支线编号看不出明显的规律
* quest_id:无法区分委托与支线，但主线编号均大于10000
* picture:Z开头主线，D开头委托，V开头支线，但有个例外，主线的最后任务没有图片
* story_per:委托没有剧情完成度，但无法区分主线与支线

任务内容的额外说明：

* 主线任务中任务变量与任务变灰无用
* 委托任务没有剧情完成度


## 升级数据
除精、神、武、防、速、运、灵几个基本属性外，只有格挡率是有成长的

* 天河有额外的20%暴击率加成
* 菱纱的格挡率初始值高
* 梦璃的格挡率初始值低
* 紫英有额外的五灵伤害加成，仅7点，另外还有15%的命中率加成
* 格挡率的成长为每4级1%


## 敌人明细
因敌人属性过多，为简单起见，同时也遵循查询优化的思路，特采取以下策略：

* 精、气、神、武、防、速、运、灵八基本属性基本均不为0值，同组，且不合并
* 水、火、雷、风、土五灵属性合并价值极高，因为极少敌人超过一个属性
* 物理伤害追加保留原样；六种伤害吸收也不予合并
* 反弹要合并，绝大多数怪不反弹；暴击、格挡、命中、反击四项合并
* 技能的前五项，多数不全有，因此也需合并
* 偷窃要么只能偷到道具，要么只有钱，必须合并
* 由于绝大部分敌人会掉落3～4件道具，因此掉落部分不予合并
* 尽管掉落金钱上下限大多数不同，但仍合并此项；受伤音效不予合并
* 死亡特效仅3种，且仅不足2成的敌人有，因此予以忽略；存盘统计保留；
* 免疫：淮南王及之后的boss均全免，昆仑四圣免疫性稍特殊，其它普通怪无免疫
* 鉴于免疫性规律明显，差异化不突出，因此全部忽略
* 关于其它属性：怪物ID及名称出现在所有明细部分
* 图标、模型、描述、Boss、经验、等级、攻击目标、攻击方式原样保留

## 物品拾取
使用SceneItem表，是自己另行挖数据的得意之作

* 场景与区块，换成对应区域的中文名称
* 模型的路径文件字段忽略，使用模型就够了
* 贴图与坐标等保留；从装备ID到最后的的金钱数量，全部合并
* 对于城镇类，追加室外场景字段，对于夜景，追加(夜)作为最后标识

### 外围场景
室内场景所在的室外场景的获取，有以下一些规律：

* 如果只有一个室外场景，则区块与场景编号相同，如果有夜景，区块结尾会有个Y
* 如果只有一个室外场景，那么室内场景全部以N开头，如有多个，则以室外场景的结束字母开头，如寿阳城编号Q03S，而其中的如意铺则为SN01；柳府的外景编号Q03X，其经堂则为XN05
* 室内场景的类型为1,室外则为1，迷宫场景的类型固定为2
* 夜景的Y标识，通常在最后，但也不尽然，室外全在最后，但室内则有可能不是
* 室内场景编号规则为其所在室外场景的末字母为首字母，后接N，如果只一个室外场景，则首字母直接为N，接下来是数字部分，通常是01开头，对于有小差异的区块，则用小写字母区分，如N03a,N03b,N03c

## 收集详情
此功能按装备、道具、配方的明细分类，以列表形式提供查询入口，查询单品的详细获得途径

* 商店则罗列其店铺所在城镇名称给你开放条件的组合
* 武器店卖武器与相关图谱，杂货店卖防具与相关图谱，锻冶与注灵图谱也在此两类店铺中购买
* 药店卖药与香料，食品摊只卖食品，其它能买到的东西，也全在杂货店购买
* 由于每个城镇，每类店铺均只有一个，因此汇总时，不详细罗列
* 偷窃则列出敌人ID与名称组合的列表
* 掉落与偷窃类似，此处不列出掉落率
* 场景拾取则必须列出主场景与子场景，包括迷宫在内

## 杂烩

* 装备、道具、配方三类物品默认列出购买途径，除剧情类因无法购买不列外
* 武器与防具如果买不到，参考配方页面自己打造
* 佩戴仅两件买不到，是委托任务的奖励
* 道具中有少部分攻击性道具是无处购买的
* 游戏中除剧情道具外，基本都可出售，少数剧情道具也可出售
* 部分剧情涉及的辅助道具，也不可出售
* 恢复类平时均可使用，攻击及辅助类只可战时使用
* 剧情类基本是不可用的
* 额外获得途径，以场景拾取为主，可能会兼任务获取的说明，支持直接输入物品的方式来查询
* 属性合并按游戏中显示的顺序来，但值的0的属性，固定不显示，因此佩饰将不显示灵蕴与潜力
* 部分属性如刀光特效等，依据物品类别来，该类均无，则最终模板不显示该列
* 鉴于菱纱武器虽是双手，偶有左右手不同的情况，但极少，额外显示该字段，意义极微，故删除
* 部分属性虽有值，但值的种类少，且有固定规律，因此也未显示，例如敌人的气增加量、法术消耗神数量等

* 普通敌人攻击我方命中，一次加5点气，与我方人员类似，BOSS一次加20点    
* 物品的获取途径：购买、偷窃、掉落、拾取、剧情得到
* 大部分的物品都可以买到，该途径合并到了各类的主页面中，部分商品的开放变量涉及任务
* 怪物只会掉落道具，偷窃也是一样
* 场景拾取可以得到部分道具与装备
* 支线与委托可以直接得到部分道具、装备、配方，但无法通过数据库进行分析
* 极少数剧情物品通过敌人掉落得到
* 道具等级仅对恢复类准确，其它类要么无价值，要么不准，因此仅此类显示此字段

go模板的强大功能，尽在**text/template**中体现无疑，要注意体会其精髓